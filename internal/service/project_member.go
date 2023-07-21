package service

import (
	"VitaTaskGo/internal/biz"
	"VitaTaskGo/internal/constant"
	"VitaTaskGo/internal/data"
	"VitaTaskGo/internal/dto"
	"VitaTaskGo/internal/pkg"
	"VitaTaskGo/internal/pkg/auth"
	"VitaTaskGo/internal/pkg/exception"
	"VitaTaskGo/internal/pkg/response"
	"VitaTaskGo/internal/pkg/state"
	"errors"
	"github.com/duke-git/lancet/v2/maputil"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/gin-gonic/gin"
	"github.com/gotidy/copy"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ProjectMemberService struct {
	Db   *gorm.DB
	ctx  *gin.Context
	repo biz.ProjectMemberRepo
}

func NewProjectMemberService(tx *gorm.DB, ctx *gin.Context) *ProjectMemberService {
	return &ProjectMemberService{
		Db:   tx,
		ctx:  ctx,
		repo: data.NewProjectMemberRepo(tx, ctx),
	}
}

// GetMember 返回一个成员
func (receiver *ProjectMemberService) GetMember(projectId uint, userId uint64) *biz.ProjectMember {
	var member *biz.ProjectMember
	receiver.Db.Where("project_id = ?", projectId).Where("user_id", userId).Limit(1).Find(&member)
	return member
}

// ProjectStar 收藏项目
func (receiver *ProjectMemberService) ProjectStar(projectId uint, userId uint64) error {
	// 项目是否存在
	if !data.NewProjectRepo(receiver.Db, receiver.ctx).Exist(projectId) {
		return exception.NewException(response.ProjectNotExist)
	}

	// 获取成员信息
	member, err := receiver.repo.GetProjectMember(projectId, userId)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return exception.ErrorHandle(err, response.DbQueryError)
	}

	// todo 是否可以使用Bing方法？
	if member == nil {
		// 不是成员，直接新增
		return receiver.Db.Create(&biz.ProjectMember{
			ProjectId: projectId,
			UserId:    userId,
			Role:      constant.ProjectStar,
		}).Error
	} else {
		// 已经是项目成员
		stateModifier := state.NewModifier(int(member.Role))
		// 是否收藏过项目
		if stateModifier.Exist(constant.ProjectStar) {
			return exception.NewException(response.ProjectStared)
		}
		// 附加上收藏者角色
		member.Role = int8(stateModifier.Attach(constant.ProjectStar))
		return receiver.Db.Save(&member).Error
	}
}

// ProjectUnStar 取消收藏项目
func (receiver *ProjectMemberService) ProjectUnStar(projectId uint, userId uint64) error {
	// 项目是否存在
	if !data.NewProjectRepo(receiver.Db, receiver.ctx).Exist(projectId) {
		return exception.NewException(response.ProjectNotExist)
	}
	member, err := receiver.repo.GetProjectMember(projectId, userId)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return exception.ErrorHandle(err, response.DbQueryError)
	}

	// todo 是否可以使用 Remove 方法
	if member != nil {
		stateModifier := state.NewModifier(int(member.Role))
		// 如果没有收藏过项目，直接返回nil
		if stateModifier.NotExist(constant.ProjectStar) {
			return nil
		}
		// 分离收藏者角色
		role := stateModifier.Detach(constant.ProjectStar)
		if role <= 0 {
			// 该成员在项目内无任何角色，删除
			return receiver.Db.Delete(&member).Error
		}
		// 保存数据
		member.Role = int8(role)
		return receiver.Db.Save(&member).Error
	}
	return nil
}

// GetMembers 获取成员
func (receiver *ProjectMemberService) GetMembers(query dto.ProjectMemberListQuery) *dto.PagedResult[dto.ProjectMemberVO] {
	// 获取列表记录
	projectMembers, total, err := receiver.repo.PageListProjectMember(query)
	if err != nil {
		logrus.Errorln("项目成员列表查询失败: ", err)
		return pkg.PagedResult[dto.ProjectMemberVO](nil, total, int64(query.Page))
	}

	projectMembersVo := make([]dto.ProjectMemberVO, len(projectMembers))
	for i, member := range projectMembers {
		roleName := make([]string, 0)
		stateModifier := state.NewModifier(int(member.Role))
		for j, item := range constant.GetProjectRoles() {
			if stateModifier.Exist(j) {
				// 如果存在
				roleName = append(roleName, item)
			}
		}

		// 拷贝数据
		vo := dto.ProjectMemberVO{}
		copiers := copy.New(func(c *copy.Options) {
			c.Skip = true
		})
		copiers.Copy(&vo, &member)
		vo.RoleName = roleName
		if member.UserInfo != nil {
			vo.UserInfo = member.UserInfo
			vo.Value = member.UserInfo.ID
			vo.Label = member.UserInfo.UserNickname
		}
		vo.Project = member.Project
		projectMembersVo[i] = vo
	}

	return pkg.PagedResult(projectMembersVo, total, int64(query.Page))
}

// Remove 移除成员
// 无视角色，直接删除成员记录
// 只允许项目负责人操作
func (receiver *ProjectMemberService) Remove(projectId uint, userIds []uint64) error {
	// 项目是否存在
	if !data.NewProjectRepo(receiver.Db, receiver.ctx).Exist(projectId) {
		return exception.NewException(response.ProjectNotExist)
	}

	// 获取当前用户
	currUser, err := auth.CurrUser(receiver.ctx)
	if err != nil {
		return err
	}

	// 是否为负责人
	if !receiver.IsLeader(projectId, currUser.ID) {
		return exception.NewException(response.MemberNotProjectLeader)
	}

	// 获取成员记录
	members, err := receiver.repo.GetProjectMembers(projectId, userIds)
	if err != nil {
		return exception.ErrorHandle(err, response.DbQueryError)
	}

	return receiver.Db.Transaction(func(tx *gorm.DB) error {
		// 重新实例化Repo
		projectMemberRepo := data.NewProjectMemberRepo(tx, receiver.ctx)

		for _, item := range members {
			// 不得删除自己
			if item.UserId == currUser.ID {
				continue
			}

			// 直接删除成员记录
			if err := projectMemberRepo.DeleteProjectMember(item.ID); err != nil {
				return err
			}
		}
		return nil
	})
}

// RemoveRole 按角色删除项目成员
// 只允许项目负责人操作
func (receiver *ProjectMemberService) RemoveRole(projectId uint, userIds []uint64, role int) error {
	// 项目是否存在
	if !data.NewProjectRepo(receiver.Db, receiver.ctx).Exist(projectId) {
		return exception.NewException(response.ProjectNotExist)
	}

	// 获取当前用户
	currUser, err := auth.CurrUser(receiver.ctx)
	if err != nil {
		return err
	}

	// 是否为负责人
	if !receiver.IsLeader(projectId, currUser.ID) {
		return exception.NewException(response.MemberNotProjectLeader)
	}

	// 获取所有角色
	if receiver.ShouldProjectRoles(role) == nil {
		// 角色不存在，返回错误
		return exception.NewException(response.ProjectRoleNonExistent)
	}
	if role == constant.ProjectLeader {
		// 不得移除项目负责人
		return exception.NewException(response.ProjectLeaderRemove)
	}

	members, err := receiver.repo.GetProjectMembers(projectId, userIds)
	if err != nil {
		return exception.ErrorHandle(err, response.DbQueryError)
	}

	return receiver.Db.Transaction(func(tx *gorm.DB) error {
		// 重新实例化Repo
		projectMemberRepo := data.NewProjectMemberRepo(tx, receiver.ctx)

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
				if err := projectMemberRepo.DeleteProjectMember(item.ID); err != nil {
					return err
				}
			} else {
				// 修改成员角色
				if err := projectMemberRepo.UpdateField(item.ID, "role", item.Role); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

// Bind 绑定普通成员
// 不检查项目是否存在
func (receiver *ProjectMemberService) Bind(projectId uint, userIds []uint64, role int) error {
	// 获取所有角色
	roles := receiver.ShouldProjectRoles(role)
	if roles == nil {
		// 角色不存在，返回错误
		return exception.NewException(response.ProjectRoleNonExistent)
	}

	if role == constant.ProjectLeader || role == constant.TaskCreator {
		// 不允许多个负责人或创建人
		if receiver.ExistMember(projectId, role) {
			// 不允许多个管理员
			return exception.NewException(response.ProjectMultipleSpecialMember)
		}
	}
	// 对userIds去重
	userIds = pkg.SliceUnique(userIds)
	members, err := receiver.repo.GetProjectMembers(projectId, userIds)
	if err != nil {
		return exception.ErrorHandle(err, response.DbQueryError)
	}

	// 获取当前任务管理员
	leader, err := receiver.GetLeader(projectId)
	if err != nil {
		return err
	}

	return receiver.Db.Transaction(func(tx *gorm.DB) error {
		// 重新实例化Repo
		projectMemberRepo := data.NewProjectMemberRepo(tx, receiver.ctx)
		// 排除列表
		excludeIds := make([]uint64, 0)
		// 此处为修改已在项目的成员权限
		for _, item := range members {
			// 将用户id加入排除列表
			excludeIds = append(excludeIds, item.UserId)
			// 负责人不能成为普通成员
			if role == constant.ProjectMember && (leader != nil && leader.UserId == item.UserId) {
				continue
			}
			// 初始化状态修改器
			stateModifier := state.NewModifier(int(item.Role))
			// 如果已是当前角色，跳过
			if stateModifier.Exist(role) {
				continue
			}
			// 将当前角色附加给该用户
			if err := projectMemberRepo.UpdateField(item.ID, "role", stateModifier.Attach(role)); err != nil {
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
			err := projectMemberRepo.CreateProjectMember(&biz.ProjectMember{
				ProjectId: projectId,
				UserId:    uid,
				Role:      int8(role), // 需要转换成int8
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// InProject 是否在项目内
func (receiver *ProjectMemberService) InProject(projectId uint, userId uint64, role int) bool {
	// 获取所有角色
	roles := receiver.ShouldProjectRoles(role)
	if roles == nil {
		// 角色不存在，返回错误
		return false
	}
	// 获取所有状态
	roleWhereIn := state.NewModifier(role).Contained(maputil.Keys(roles))
	return receiver.repo.InProject(projectId, userId, roleWhereIn)
}

// GetMembersByRole 根据角色获取成员
func (receiver *ProjectMemberService) GetMembersByRole(projectId uint, role int) ([]biz.ProjectMember, error) {
	// 获取所有角色
	roles := receiver.ShouldProjectRoles(role)
	if roles == nil {
		// 角色不存在，返回错误
		return nil, exception.NewException(response.ProjectRoleNonExistent)
	}

	// 获取所有状态
	roleWhereIn := state.NewModifier(role).Contained(maputil.Keys(roles))
	return receiver.repo.GetMembersByRole(projectId, roleWhereIn)
}

// ShouldProjectRoles 获取角色列表并且判断role是否在列表内
func (receiver *ProjectMemberService) ShouldProjectRoles(role int) map[int]string {
	// 获取所有角色
	roles := constant.GetProjectRoles()
	// 判断role参数是否符合在角色列表内
	if _, ok := roles[role]; !ok {
		// 角色不存在，返回错误
		return nil
	}

	return roles
}

// Transfer 移交项目
func (receiver *ProjectMemberService) Transfer(projectId uint, recipient uint64) error {
	// 先判断接收人是否存在
	if !data.NewUserRepo(receiver.Db, receiver.ctx).Exist(recipient) {
		return exception.NewException(response.UserNotFound)
	}

	// 获取旧项目负责人
	transferorMember, err := receiver.GetLeader(projectId)
	// 允许没有负责人的情况，但是要处理其它的错误
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return exception.ErrorHandle(err, response.DbQueryError)
	}

	return receiver.Db.Transaction(func(tx *gorm.DB) error {
		// 重新实例化Repo
		projectMemberRepo := data.NewProjectMemberRepo(tx, receiver.ctx)
		/* 删除角色 Start */
		if transferorMember != nil {
			// 如果负责人ID和接收人ID相同
			if transferorMember.UserId == recipient {
				// 不执行操作
				return nil
			}

			// 初始化状态修改器
			stateModifier := state.NewModifier(int(transferorMember.Role))
			// 删除负责人角色
			role := stateModifier.Detach(constant.ProjectLeader)
			if role <= 0 {
				// 无角色，删除
				if err := projectMemberRepo.DeleteProjectMember(transferorMember.ID); err != nil {
					return exception.ErrorHandle(err, response.DbExecuteError)
				}
			} else {
				// 修改
				if err := projectMemberRepo.UpdateField(transferorMember.ID, "role", role); err != nil {
					return exception.ErrorHandle(err, response.DbExecuteError)
				}
			}
		}
		/* 删除角色 End */

		/* 新增或给成员附加负责人角色 Start */
		// 获取接收人
		recipientMember, err := projectMemberRepo.GetProjectMember(projectId, recipient)
		// 没有查询到记录
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 该用户不在项目内，新增一个
				e := projectMemberRepo.CreateProjectMember(&biz.ProjectMember{
					ProjectId: projectId,
					UserId:    recipient,
					Role:      int8(constant.ProjectLeader), // 需要转换成int8
				})
				return exception.ErrorHandle(e, response.DbExecuteError)
			}

			// 有错误，返回
			return exception.ErrorHandle(err, response.DbQueryError)
		}

		// 成员存在，给ta附加一个负责人角色
		e := projectMemberRepo.UpdateField(recipientMember.ID, "role", state.NewModifier(int(recipientMember.Role)).Attach(constant.ProjectLeader))
		return exception.ErrorHandle(e, response.DbQueryError)
		/* 新增或给成员附加负责人角色 End */
	})
}

// ExistMember 项目是否存在指定角色的成员
func (receiver *ProjectMemberService) ExistMember(projectId uint, role int) bool {
	// 获取所有角色
	roles := receiver.ShouldProjectRoles(role)
	if roles == nil {
		// 角色不存在，返回错误
		return false
	}
	// 获取所有状态
	roleWhereIn := state.NewModifier(role).Contained(maputil.Keys(roles))
	// 查询记录
	list, _ := receiver.repo.GetMembersByRole(projectId, roleWhereIn)

	return len(list) > 0
}

// GetLeader 获取负责人
func (receiver *ProjectMemberService) GetLeader(projectId uint) (*biz.ProjectMember, error) {
	members, err := receiver.GetMembersByRole(projectId, constant.TaskLeader)
	if err != nil {
		return nil, err
	}
	// 如果没有负责人
	if len(members) <= 0 {
		return nil, nil
	}
	return &members[0], nil
}

// IsLeader 用户是否为项目负责人
func (receiver *ProjectMemberService) IsLeader(projectId uint, userId uint64) bool {
	// 获取项目负责人
	leader, err := receiver.GetLeader(projectId)
	if err != nil {
		_ = exception.ErrorHandle(err, response.DbQueryError)
		return false
	}

	return leader.UserId == userId
}
