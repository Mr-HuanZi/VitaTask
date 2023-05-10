package logic

import (
	"VitaTaskGo/app/constant"
	"VitaTaskGo/app/exception"
	"VitaTaskGo/app/extend"
	"VitaTaskGo/app/model"
	"VitaTaskGo/app/response"
	"VitaTaskGo/library/db"
	"VitaTaskGo/library/state"
	"errors"
	"github.com/duke-git/lancet/v2/maputil"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	ErrSameTaskLeader = errors.New("same task leader")
)

type TaskMemberLogic struct {
	Db  *gorm.DB
	ctx *gin.Context
}

func NewTaskMemberLogic(ctx *gin.Context) *TaskMemberLogic {
	return &TaskMemberLogic{
		Db:  db.Db, // 赋予ORM实例
		ctx: ctx,   // 传递上下文
	}
}

// Bind 绑定任务成员
// 不检查任务是否存在
func (receiver TaskMemberLogic) Bind(taskId uint, userIds []uint64, role int) error {
	var lists []model.TaskMember
	where := model.TaskMember{
		TaskId: taskId,
	}
	// 获取所有角色
	roles := receiver.ShouldRoles(role)
	if roles == nil {
		// 角色不存在
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
	if err := receiver.Db.Model(&model.TaskMember{}).Where(&where).Where("user_id IN ?", userIds).Find(&lists).Error; err != nil {
		return err
	}

	// 获取当前任务负责人
	leader, err := receiver.GetLeader(taskId)
	if err != nil {
		return err
	}

	return receiver.Db.Transaction(func(tx *gorm.DB) error {
		// 排除列表
		excludeIds := make([]uint64, 0)
		// 此处为修改已在任务的成员权限
		for _, item := range lists {
			// 将用户id加入排除列表
			excludeIds = append(excludeIds, item.UserId)
			// 负责人不能成为普通成员
			if role == constant.TaskMember && leader != nil && leader.UserId == item.UserId {
				continue
			}
			// 初始化状态修改器
			stateModifier := state.NewModifier(int(item.Role))
			// 将当前角色附加给该用户
			// 保存
			if err := tx.Model(&item).Update("role", stateModifier.Attach(role)).Error; err != nil {
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
// 允许移除负责人
func (receiver TaskMemberLogic) Remove(taskId uint, userIds []uint64, role int) error {
	var lists []model.TaskMember
	where := model.TaskMember{
		TaskId: taskId,
	}

	if role > 0 {
		// 传入指定角色时生效
		if receiver.ShouldRoles(role) == nil {
			// 角色不存在
			return exception.NewException(response.TaskRoleNonExistent)
		}
		if role == constant.TaskCreator {
			// 不得移除创建人
			return exception.NewException(response.TaskCreatorRemove)
		}
	}

	// 提供了用户Id参数
	if len(userIds) > 0 {
		if err := receiver.Db.Model(&model.TaskMember{}).Where(&where).Where("user_id IN ?", userIds).Find(&lists).Error; err != nil {
			return err
		}
	} else {
		// 查询指定角色的人员
		// 获取所有角色
		roles := constant.GetTaskRoles()
		// 获取所有可能的状态
		roleWhereIn := state.NewModifier(role).Contained(maputil.Keys(roles))
		if err := receiver.Db.Model(&model.TaskMember{}).Where(&where).Where("role IN ?", roleWhereIn).Find(&lists).Error; err != nil {
			return err
		}
	}

	return receiver.Db.Transaction(func(tx *gorm.DB) error {
		for _, item := range lists {
			// 初始化状态修改器
			stateModifier := state.NewModifier(int(item.Role))
			// 未传入角色参数，直接删除成员
			if role <= 0 {
				// 不得移除创建人
				if stateModifier.Exist(constant.TaskCreator) {
					return exception.NewException(response.TaskCreatorRemove)
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
				if err := tx.Model(&item).Update("role", item.Role).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

// ShouldRoles 获取角色列表并且判断role是否在列表内
func (receiver TaskMemberLogic) ShouldRoles(role int) map[int]string {
	// 获取所有角色
	roles := constant.GetTaskRoles()
	// 判断role参数是否符合在角色列表内
	if _, ok := roles[role]; !ok {
		return nil
	}
	return roles
}

// ExistMember 任务是否存在指定角色的成员
func (receiver TaskMemberLogic) ExistMember(taskId uint, role int) bool {
	var count int64
	// 获取所有角色
	roles := constant.GetTaskRoles()
	// 获取所有可能的状态
	roleWhereIn := state.NewModifier(role).Contained(maputil.Keys(roles))
	// 查询
	receiver.Db.Model(&model.TaskMember{}).
		Where("task_id = ?", taskId).
		Where("role IN ?", roleWhereIn).
		Count(&count)

	return count > 0
}

// InTask 是否在任务内
func (receiver TaskMemberLogic) InTask(taskId uint, userId uint64, role int) bool {
	tx := receiver.Db.Model(&model.TaskMember{}).Where(&model.TaskMember{TaskId: taskId, UserId: userId})
	// 如果带入角色
	if role > 0 {
		// 获取所有角色
		roles := receiver.ShouldRoles(role)
		if roles == nil {
			// 角色不存在
			return false
		}
		// 获取所有可能的状态
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

// GetMembersByRole 获取指定角色的成员
func (receiver TaskMemberLogic) GetMembersByRole(taskId uint, role int) ([]model.TaskMember, error) {
	var lists []model.TaskMember
	where := map[string]interface{}{"task_id": taskId}
	// 获取所有角色
	roles := receiver.ShouldRoles(role)
	if roles == nil {
		// 角色不存在
		return nil, exception.NewException(response.TaskRoleNonExistent)
	}
	// 获取所有可能的状态
	roleWhereIn := state.NewModifier(role).Contained(maputil.Keys(roles))
	err := receiver.Db.Model(&model.TaskMember{}).Where(where).Where("role IN ?", roleWhereIn).Find(&lists).Error
	return lists, err
}

// GetLeader 获取负责人
func (receiver TaskMemberLogic) GetLeader(taskId uint) (*model.TaskMember, error) {
	members, err := receiver.GetMembersByRole(taskId, constant.TaskLeader)
	if err != nil {
		return nil, err
	}
	// 如果没有负责人
	if len(members) <= 0 {
		return nil, nil
	}
	return &members[0], nil
}

// GetAllMember 获取所有用户
func (receiver TaskMemberLogic) GetAllMember(taskId uint) ([]model.TaskMember, error) {
	var members []model.TaskMember
	where := model.TaskMember{
		TaskId: taskId,
	}
	// 查询用户
	if err := receiver.Db.Model(&model.TaskMember{}).Where(&where).Find(&members).Error; err != nil {
		return nil, err
	}
	return members, nil
}

// GetTaskIdsByUsers 获取指定角色用户所在的任务id
func (receiver TaskMemberLogic) GetTaskIdsByUsers(userIds []uint64, role int) ([]uint, error) {
	var taskIds []uint
	// 获取所有角色
	roles := receiver.ShouldRoles(role)
	if roles == nil {
		// 角色不存在
		return nil, exception.NewException(response.TaskRoleNonExistent)
	}
	// 获取所有可能的状态
	roleWhereIn := state.NewModifier(role).Contained(maputil.Keys(roles))
	err := receiver.Db.Model(&model.TaskMember{}).Select("task_id").Where("user_id IN ?", userIds).Where("role IN ?", roleWhereIn).Find(&taskIds).Error
	return taskIds, err
}

// Transfer 转移任务负责人
// 接收人与负责人相同时返回 ErrSameTaskLeader 错误
func (receiver TaskMemberLogic) Transfer(taskId uint, recipient uint64) error {
	// 获取所有角色
	roles := constant.GetTaskRoles()
	// 获取所有可能的状态
	roleWhereIn := state.NewModifier(constant.TaskLeader).Contained(maputil.Keys(roles))
	// 遍历需要移交的任务
	var leaders []model.TaskMember
	where := map[string]interface{}{"task_id": taskId}
	// 查询所有负责人角色，避免某些意外情况下存在多个负责人
	err := receiver.Db.Model(&model.TaskMember{}).Where(where).Where("role IN ?", roleWhereIn).Find(&leaders).Error
	if err != nil {
		return err
	}
	// 定义跳过标记
	skip := false
	// 遍历
	for _, leader := range leaders {
		// 如果接收人与当前负责人是同一个uid，则跳过
		// 但是仍会删除其他的负责人
		if leader.UserId == recipient {
			skip = true
			continue
		}
		// 初始化状态修改器
		stateModifier := state.NewModifier(int(leader.Role))
		// 移除当前人员的负责人角色
		leader.Role = int8(stateModifier.Detach(constant.TaskLeader))
		if leader.Role <= 0 {
			// 无角色，删除
			if err := receiver.Db.Delete(&leader).Error; err != nil {
				return err
			}
		} else {
			// 修改成员角色
			if err := receiver.Db.Model(&leader).Update("role", leader.Role).Error; err != nil {
				return err
			}
		}
	}
	// 是否跳过
	if skip {
		return ErrSameTaskLeader
	}
	/* 重新绑定负责人 Start */
	var member model.TaskMember
	err = receiver.Db.Model(&model.TaskMember{}).Where(where).Where("user_id = ?", recipient).First(&member).Error
	if err == gorm.ErrRecordNotFound {
		// 该成员不在任务内
		newMember := model.TaskMember{
			TaskId: taskId,
			UserId: recipient,
			Role:   int8(constant.TaskLeader), // 需要转换成int8
		}
		if err := receiver.Db.Model(&model.TaskMember{}).Create(&newMember).Error; err != nil {
			return err
		}
	} else if err != nil {
		// 存在其它错误，返回
		return err
	} else {
		// 已存在该成员
		stateModifier := state.NewModifier(int(member.Role))
		if err := receiver.Db.Model(&member).Update("role", stateModifier.Attach(constant.TaskLeader)).Error; err != nil {
			return err
		}
	}
	/* 重新绑定负责人 End */
	return nil
}
