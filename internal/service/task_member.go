package service

import (
	"VitaTaskGo/internal/biz"
	"VitaTaskGo/internal/constant"
	"VitaTaskGo/internal/data"
	"VitaTaskGo/internal/pkg"
	"VitaTaskGo/internal/pkg/auth"
	"VitaTaskGo/internal/pkg/exception"
	"VitaTaskGo/internal/pkg/response"
	"VitaTaskGo/internal/pkg/state"
	"errors"
	"github.com/duke-git/lancet/v2/maputil"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	ErrSameTaskLeader = errors.New("same task leader")
)

type TaskMemberService struct {
	Db   *gorm.DB
	ctx  *gin.Context
	repo biz.TaskMemberRepo
}

func NewTaskMemberService(tx *gorm.DB, ctx *gin.Context) *TaskMemberService {
	return &TaskMemberService{
		Db:   tx,
		ctx:  ctx,
		repo: data.NewTaskMemberRepo(tx, ctx),
	}
}

// Bind 绑定任务成员
// 不检查任务是否存在
func (receiver TaskMemberService) Bind(taskId uint, userIds []uint64, role int) error {
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
	userIds = pkg.SliceUnique(userIds)
	// 获取任务成员
	members, err := receiver.repo.GetTaskMembers(taskId, userIds)
	if err != nil {
		return exception.ErrorHandle(err, response.DbQueryError)
	}

	// 获取当前任务负责人
	leader, err := receiver.GetLeader(taskId)
	if err != nil {
		return err
	}

	return receiver.Db.Transaction(func(tx *gorm.DB) error {
		// 重新实例化Repo
		taskMemberRepo := data.NewTaskMemberRepo(tx, receiver.ctx)
		// todo 需要把成员同步到项目

		// 排除列表
		excludeIds := make([]uint64, 0)
		// 此处为修改已在任务的成员权限
		for _, item := range members {
			// 将用户id加入排除列表
			excludeIds = append(excludeIds, item.UserId)
			// 负责人不能成为普通成员
			if role == constant.TaskMember && (leader != nil && leader.UserId == item.UserId) {
				continue
			}
			// 初始化状态修改器
			stateModifier := state.NewModifier(int(item.Role))
			// 保存
			if err := taskMemberRepo.UpdateField(item.ID, "role", stateModifier.Attach(role)); err != nil {
				return err
			}
		}

		// 计算两个切片的差集，得到的结果就是需要新增的成员ID切片
		userIds = slice.Difference(userIds, excludeIds)
		// 此处为新增成员
		for _, uid := range userIds {
			// 跳过上面已经处理过的
			if slice.Contain(excludeIds, uid) {
				continue
			}
			// 创建新的成员信息
			err := taskMemberRepo.Create(&biz.TaskMember{
				TaskId: taskId,
				UserId: uid,
				Role:   int8(role), // 需要转换成int8
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// Remove 移除成员
// 无视角色，直接删除成员记录
// 允许移除负责人
func (receiver TaskMemberService) Remove(taskId uint, userIds []uint64) error {
	// 获取成员记录
	members, err := receiver.repo.GetTaskMembers(taskId, userIds)
	if err != nil {
		return exception.ErrorHandle(err, response.DbQueryError)
	}

	// 获取当前用户
	currUser, err := auth.CurrUser(receiver.ctx)
	if err != nil {
		return err
	}

	return receiver.Db.Transaction(func(tx *gorm.DB) error {
		// 重新实例化Repo
		taskMemberRepo := data.NewTaskMemberRepo(tx, receiver.ctx)

		for _, item := range members {
			// 不得删除自己
			if item.UserId == currUser.ID {
				continue
			}

			// 直接删除成员记录
			if err := taskMemberRepo.Delete(item.ID); err != nil {
				return err
			}
		}
		return nil
	})
}

// RemoveRole 按角色删除项目成员
// 允许移除负责人
func (receiver TaskMemberService) RemoveRole(taskId uint, userIds []uint64, role int) error {
	var (
		err     error
		members []biz.TaskMember
	)

	if receiver.ShouldRoles(role) == nil {
		// 角色不存在
		return exception.NewException(response.TaskRoleNonExistent)
	}
	if role == constant.TaskCreator {
		// 不得移除创建人
		return exception.NewException(response.TaskCreatorRemove)
	}

	if len(userIds) > 0 {
		members, err = receiver.repo.GetTaskMembers(taskId, userIds)
		if err != nil {
			return exception.ErrorHandle(err, response.DbQueryError)
		}
	} else {
		// 没有提供用户ID就获取当前任务所有的成员
		members, err = receiver.repo.GetTaskAllMember(taskId)
		if err != nil {
			return exception.ErrorHandle(err, response.DbQueryError)
		}
	}

	return receiver.Db.Transaction(func(tx *gorm.DB) error {
		// 重新实例化Repo
		taskMemberRepo := data.NewTaskMemberRepo(tx, receiver.ctx)

		for _, item := range members {
			// 初始化状态修改器
			stateModifier := state.NewModifier(int(item.Role))
			if stateModifier.NotExist(role) {
				// 不是该角色，跳过
				continue
			}

			// 去除当前角色
			item.Role = int8(stateModifier.Detach(role))
			if item.Role <= 0 {
				// 无角色，删除
				if err := taskMemberRepo.Delete(item.ID); err != nil {
					return err
				}
			} else {
				// 修改成员角色
				if err := taskMemberRepo.UpdateField(item.ID, "role", item.Role); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

// ShouldRoles 获取角色列表并且判断role是否在列表内
func (receiver TaskMemberService) ShouldRoles(role int) map[int]string {
	// 获取所有角色
	roles := constant.GetTaskRoles()
	// 判断role参数是否符合在角色列表内
	if _, ok := roles[role]; !ok {
		return nil
	}
	return roles
}

// ExistMember 任务是否存在指定角色的成员
func (receiver TaskMemberService) ExistMember(taskId uint, role int) bool {
	// 获取所有角色
	roles := receiver.ShouldRoles(role)
	if roles == nil {
		// 角色不存在，返回错误
		return false
	}
	// 获取所有可能的状态
	roleWhereIn := state.NewModifier(role).Contained(maputil.Keys(roles))
	// 查询
	l, _ := receiver.repo.GetMembersByRole(taskId, roleWhereIn)

	return len(l) > 0
}

// InTask 是否在任务内
func (receiver TaskMemberService) InTask(taskId uint, userId uint64, role int) bool {
	var roleWhereIn []int
	// 如果带入角色
	if role > 0 {
		// 获取所有角色
		roles := receiver.ShouldRoles(role)
		if roles == nil {
			// 角色不存在
			return false
		}
		// 获取所有可能的状态
		roleWhereIn = state.NewModifier(role).Contained(maputil.Keys(roles))
	}
	return receiver.repo.InTask(taskId, userId, roleWhereIn)
}

// GetMembersByRole 获取指定角色的成员
func (receiver TaskMemberService) GetMembersByRole(taskId uint, role int) ([]biz.TaskMember, error) {
	// 获取所有角色
	roles := receiver.ShouldRoles(role)
	if roles == nil {
		// 角色不存在
		return nil, exception.NewException(response.TaskRoleNonExistent)
	}
	// 获取所有可能的状态
	roleWhereIn := state.NewModifier(role).Contained(maputil.Keys(roles))
	return receiver.repo.GetMembersByRole(taskId, roleWhereIn)
}

// GetLeader 获取负责人
func (receiver TaskMemberService) GetLeader(taskId uint) (*biz.TaskMember, error) {
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

// GetTaskIdsByUsers 获取指定角色用户所在的任务id
func (receiver TaskMemberService) GetTaskIdsByUsers(userIds []uint64, role int) ([]uint, error) {
	// 获取所有角色
	roles := receiver.ShouldRoles(role)
	if roles == nil {
		// 角色不存在
		return nil, exception.NewException(response.TaskRoleNonExistent)
	}
	// 获取所有可能的状态
	roleWhereIn := state.NewModifier(role).Contained(maputil.Keys(roles))
	return receiver.repo.GetTaskIdsByUsers(userIds, roleWhereIn)
}
