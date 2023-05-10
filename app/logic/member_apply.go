package logic

import (
	"VitaTaskGo/app/constant"
	"VitaTaskGo/app/exception"
	"VitaTaskGo/app/extend"
	"VitaTaskGo/app/model"
	"VitaTaskGo/app/response"
	"VitaTaskGo/library/db"
	"VitaTaskGo/library/state"
	"github.com/duke-git/lancet/v2/maputil"
	"github.com/duke-git/lancet/v2/slice"
	"gorm.io/gorm"
)

type MemberApplyMode interface {
	model.TaskMember | model.ProjectMember
}

type MemberApplyLogic[T MemberApplyMode] struct {
	Db   *gorm.DB
	mode string
}

func NewMemberApplyLogic[T MemberApplyMode](mode string) *MemberApplyLogic[T] {
	return &MemberApplyLogic[T]{
		Db:   db.Db, // 赋予ORM实例
		mode: mode,
	}
}

// Bind 绑定任务成员
// 不检查任务是否存在
func (receiver MemberApplyLogic[T]) Bind(taskId uint, userIds []uint64, role int) error {
	var lists []model.TaskMember // todo 如果直接用 T 来定义lists，那么会导致下面的for range直接报错
	where := model.TaskMember{
		TaskId: taskId,
	}
	// 获取所有角色
	roles := receiver.ShouldRoles(role)
	if roles == nil {
		// 角色不存在，返回错误
		return exception.NewException(response.TaskRoleNonExistent)
	}
	// 不允许多个负责人或创建人
	if role == constant.TaskLeader || role == constant.TaskCreator {
		if receiver.ExistMember(taskId, role) {
			return exception.NewException(response.TaskMultipleSpecialMember)
		}
	}
	// 对userIds去重
	userIds = extend.SliceUnique(userIds)
	// 查找数据库中是否存在这些用户
	if err := receiver.Db.Model(new(T)).Where(&where).Where("user_id IN ?", userIds).Find(&lists).Error; err != nil {
		return err
	}

	return receiver.Db.Transaction(func(tx *gorm.DB) error {
		// 排除列表
		excludeIds := make([]uint64, 0)
		// 此处为修改已在任务的成员权限
		for _, item := range lists {
			// 将用户id加入排除列表
			excludeIds = append(excludeIds, item.UserId)
			// 初始化状态修改器
			stateModifier := state.NewModifier(int(item.Role))
			// 将当前角色附加给该用户
			item.Role = int8(stateModifier.Attach(role))
			// 保存
			if err := tx.Save(&item).Error; err != nil {
				return err
			}
		}
		addLists := make([]model.TaskMember, 0)
		// 此处为新增成员
		for _, uid := range userIds {
			// 跳过上面已经处理过的
			if slice.Contain(excludeIds, uid) {
				continue
			}
			// 创建新的成员信息
			addLists = append(addLists, model.TaskMember{
				TaskId: taskId,
				UserId: uid,
				Role:   int8(role), // 需要转换成int8
			})
		}
		if len(addLists) > 0 {
			// 插入新成员
			return tx.Model(&model.TaskMember{}).Create(addLists).Error
		}
		return nil
	})
}

// Remove 移除成员
// 如果不提供 role或者role小于等于0，则直接删除成员
func (receiver MemberApplyLogic[T]) Remove(taskId uint, userIds []uint64, role int) error {
	var lists []model.TaskMember
	where := model.TaskMember{
		TaskId: taskId,
	}
	// 获取所有角色
	roles := receiver.ShouldRoles(role)
	if roles == nil {
		// 角色不存在，返回错误
		return exception.NewException(response.ProjectRoleNonExistent)
	}
	if role == constant.ProjectLeader {
		// 不得移除项目负责人
		return exception.NewException(response.ProjectLeaderRemove)
	}

	if err := receiver.Db.Model(&model.ProjectMember{}).Where(&where).Where("user_id IN ?", userIds).Find(&lists).Error; err != nil {
		return err
	}

	return receiver.Db.Transaction(func(tx *gorm.DB) error {
		for _, item := range lists {
			// 初始化状态修改器
			stateModifier := state.NewModifier(int(item.Role))
			if role <= 0 {
				// 未传入角色参数，直接删除成员
				if stateModifier.Exist(constant.ProjectLeader) {
					// 不得移除项目负责人
					return exception.NewException(response.ProjectLeaderRemove)
				}
				if err := tx.Delete(&item).Error; err != nil {
					return err
				}
				continue
			}
			if stateModifier.NotExist(role) {
				// 不是该角色，跳过
				continue
			}
			// 去除当前角色
			item.Role = int8(stateModifier.Detach(role))
			if item.Role <= 0 {
				// 无角色，删除
				if err := tx.Delete(&item).Error; err != nil {
					return err
				}
			} else {
				// 修改成员角色
				if err := tx.Save(&item).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

// ShouldRoles 获取角色列表并且判断role是否在列表内
func (receiver MemberApplyLogic[T]) ShouldRoles(role int) map[int]string {
	// 获取所有角色
	roles := constant.GetTaskRoles()
	// 判断role参数是否符合在角色列表内
	if _, ok := roles[role]; !ok {
		return nil
	}
	return roles
}

// ExistMember 项目是否存在指定角色的成员
func (receiver MemberApplyLogic[T]) ExistMember(taskId uint, role int) bool {
	var count int64
	// 获取所有角色
	roles := constant.GetTaskRoles()
	// 获取所有状态
	roleWhereIn := state.NewModifier(role).Contained(maputil.Keys(roles))
	// 查询
	receiver.Db.Model(&model.TaskMember{}).
		Where("task_id = ?", taskId).
		Where("role IN ?", roleWhereIn).
		Count(&count)

	return count > 0
}

// InTask 是否在任务内
func (receiver MemberApplyLogic[T]) InTask(taskId uint, userId uint64, role int) bool {
	tx := receiver.Db.Model(&model.TaskMember{}).Where(&model.TaskMember{TaskId: taskId, UserId: userId})
	// 如果带入角色
	if role > 0 {
		// 获取所有角色
		roles := receiver.ShouldRoles(role)
		if roles == nil {
			// 角色不存在
			return false
		}
		// 获取所有状态
		roleWhereIn := state.NewModifier(role).Contained(maputil.Keys(roles))
		// 加入到查询条件中
		tx.Where("role IN ?", roleWhereIn)
	}
	// 统计数量
	var count int64
	if err := tx.Count(&count).Error; err != nil {
		// 记录日志，但不返回错误
		_ = exception.ErrorHandle(err, response.DbQueryError)
		return false
	}
	return count > 0
}
