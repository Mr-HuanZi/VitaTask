package logic

import (
	"VitaTaskGo/app/constant"
	"VitaTaskGo/app/exception"
	"VitaTaskGo/app/extend"
	"VitaTaskGo/app/model"
	"VitaTaskGo/app/response"
	"VitaTaskGo/app/types"
	"VitaTaskGo/library/db"
	"VitaTaskGo/library/state"
	"errors"
	"github.com/duke-git/lancet/v2/maputil"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProjectMemberLogic struct {
	Db  *gorm.DB
	ctx *gin.Context
}

func NewProjectMemberLogic(ctx *gin.Context) *ProjectMemberLogic {
	return &ProjectMemberLogic{
		Db:  db.Db, // 赋予ORM实例
		ctx: ctx,   // 传递上下文
	}
}

// ProjectRelationLeader 关联项目负责人
func (receiver *ProjectMemberLogic) ProjectRelationLeader(data types.RelationLeaderForm) (*model.ProjectMember, error) {
	leaderModel := &model.ProjectMember{
		ProjectId: data.ProjectId,
		UserId:    data.UserId,
		Role:      constant.ProjectLeader,
	}
	return leaderModel, receiver.Db.Create(leaderModel).Error
}

// ProjectStared 用户已收藏项目
func (receiver *ProjectMemberLogic) ProjectStared(projectId uint, userId uint64) bool {
	// 获取所有状态
	roleWhereIn := extend.SliceOperator(maputil.Keys(constant.GetProjectRoles()), constant.ProjectStar, "|")

	var count int64
	receiver.Db.Model(&model.ProjectMember{}).
		Where("project_id = ?", projectId).
		Where("user_id", userId).
		Where("role IN ?", roleWhereIn).
		Count(&count)
	return count > 0
}

// GetMember 返回一个成员
func (receiver *ProjectMemberLogic) GetMember(projectId uint, userId uint64) *model.ProjectMember {
	var member *model.ProjectMember
	receiver.Db.Where("project_id = ?", projectId).Where("user_id", userId).Limit(1).Find(&member)
	return member
}

// ProjectStar 收藏项目
func (receiver *ProjectMemberLogic) ProjectStar(projectId uint, userId uint64) error {
	// 项目是否存在
	if !NewProjectLogic(receiver.ctx).ProjectExist(projectId) {
		return exception.NewException(response.ProjectNotExist)
	}
	// 获取成员信息
	member := receiver.GetMember(projectId, userId)
	if member == nil {
		// 不是成员，直接新增
		return receiver.Db.Create(&model.ProjectMember{
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
func (receiver *ProjectMemberLogic) ProjectUnStar(projectId uint, userId uint64) error {
	// 项目是否存在
	if !NewProjectLogic(receiver.ctx).ProjectExist(projectId) {
		return exception.NewException(response.ProjectNotExist)
	}
	member := receiver.GetMember(projectId, userId)
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
func (receiver *ProjectMemberLogic) GetMembers(query types.ProjectMemberListQuery) types.PagedResult[types.ProjectMemberVO] {
	var (
		projectMembers   []model.ProjectMember
		count            int64
		roleName         []string
		projectMembersVo []types.ProjectMemberVO
	)

	tx := receiver.QueryHandle(query)

	// 计算总记录数
	tx.Model(&model.ProjectMember{}).Count(&count)

	if err := tx.Model(&model.ProjectMember{}).Scopes(db.Paginate(&query.Page, &query.PageSize)).Order("role ASC").Find(&projectMembers).Error; err != nil {
		_ = exception.ErrorHandle(err, response.ProjectMemberQueryFail, "项目成员列表查询失败: ")
	}

	projectMembersVo = make([]types.ProjectMemberVO, 0)
	for _, member := range projectMembers {
		roleName = make([]string, 0)
		stateModifier := state.NewModifier(int(member.Role))
		for j, item := range constant.GetProjectRoles() {
			if stateModifier.Exist(j) {
				// 如果存在
				roleName = append(roleName, item)
			}
		}
		projectMembersVo = append(projectMembersVo, types.ProjectMemberVO{RoleName: roleName, Value: member.UserInfo.ID, Label: member.UserInfo.UserNickname, ProjectMember: member})
	}

	return types.PagedResult[types.ProjectMemberVO]{
		Items: projectMembersVo,
		Total: count,
		Page:  int64(query.Page),
	}
}

// GetAll 获取项目所有用户
func (receiver *ProjectMemberLogic) GetAll(projectId uint) ([]model.ProjectMember, error) {
	var members []model.ProjectMember
	where := model.ProjectMember{
		ProjectId: projectId,
	}
	// 查询用户
	if err := receiver.Db.Model(&model.ProjectMember{}).Where(&where).Find(&members).Error; err != nil {
		return nil, err
	}

	return members, nil
}

// QueryHandle 查询条件处理
func (receiver *ProjectMemberLogic) QueryHandle(query types.ProjectMemberListQuery) *gorm.DB {
	tx := receiver.Db.Where("project_id = ?", query.ProjectId)
	if query.Role > 0 {
		tx = tx.Where("role", query.Role)
	}

	if query.Username != "" || query.Nickname != "" {
		// Joins连接
		tx = tx.Joins("UserInfo")

		if query.Username != "" {
			tx = tx.Clauses(clause.Like{Column: "UserInfo.user_login", Value: "%" + query.Username + "%"})
		}
		if query.Nickname != "" {
			tx = tx.Clauses(clause.Like{Column: "UserInfo.user_nickname", Value: "%" + query.Nickname + "%"})
		}
		return tx
	} else {
		// 预加载
		return tx.Preload("UserInfo")
	}
}

// Remove 移除成员
// 如果不提供 role或者role小于等于0，则直接删除成员
func (receiver *ProjectMemberLogic) Remove(projectId uint, userIds []uint64, role int) error {
	// 项目是否存在
	if !NewProjectLogic(receiver.ctx).ProjectExist(projectId) {
		return exception.NewException(response.ProjectNotExist)
	}
	var lists []model.ProjectMember
	where := model.ProjectMember{
		ProjectId: projectId,
	}
	// 获取所有角色
	roles := receiver.ShouldProjectRoles(role)
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
				if err := tx.Model(&item).Update("role", item.Role).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

// Bind 绑定普通成员
// 不检查项目是否存在
func (receiver *ProjectMemberLogic) Bind(projectId uint, userIds []uint64, role int) error {
	var lists []model.ProjectMember
	where := model.ProjectMember{
		ProjectId: projectId,
	}
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
	userIds = extend.SliceUnique(userIds)
	// 查找数据库中是否存在这些用户
	if err := receiver.Db.Model(&model.ProjectMember{}).Where(&where).Where("user_id IN ?", userIds).Find(&lists).Error; err != nil {
		return err
	}

	// 获取当前任务管理员
	leader, err := receiver.GetLeader(projectId)
	if err != nil {
		return err
	}

	return receiver.Db.Transaction(func(tx *gorm.DB) error {
		// 排除列表
		excludeIds := make([]uint64, 0)
		// 此处为修改已在项目的成员权限
		for _, item := range lists {
			// 将用户id加入排除列表
			excludeIds = append(excludeIds, item.UserId)
			// 负责人不能成为普通成员
			if role == constant.ProjectMember && leader != nil && leader.UserId == item.UserId {
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
		addLists := make([]model.ProjectMember, 0)
		// 此处为新增成员
		for _, uid := range userIds {
			// 跳过上面已经处理过的
			if slice.Contain(excludeIds, uid) {
				continue
			}
			// 创建新的成员信息
			addLists = append(addLists, model.ProjectMember{
				ProjectId: projectId,
				UserId:    uid,
				Role:      int8(role), // 需要转换成int8
			})
		}
		if len(addLists) > 0 {
			// 插入新成员
			return tx.Model(&model.ProjectMember{}).Create(addLists).Error
		}
		return nil
	})
}

// InProject 是否在项目内
func (receiver *ProjectMemberLogic) InProject(projectId uint, userId uint64, role int) bool {
	tx := receiver.Db.Model(&model.ProjectMember{}).Where(&model.ProjectMember{ProjectId: projectId, UserId: userId})
	// 如果带入角色
	if role > 0 {
		// 获取所有角色
		roles := receiver.ShouldProjectRoles(role)
		if roles == nil {
			// 角色不存在，返回错误
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
		_ = exception.ErrorHandle(err, response.DbQueryError)
		return false
	}
	return count > 0
}

// GetMembersByRole 根据角色获取成员
func (receiver *ProjectMemberLogic) GetMembersByRole(projectId uint, role int) ([]model.ProjectMember, error) {
	var members []model.ProjectMember
	tx := receiver.Db.Model(&model.ProjectMember{}).Where(&model.ProjectMember{ProjectId: projectId})
	// 如果带入角色
	if role > 0 {
		// 获取所有角色
		roles := constant.GetProjectRoles()
		// 判断role参数是否符合在角色列表内
		if _, ok := roles[role]; !ok {
			// 角色不存在，返回错误
			return nil, exception.NewException(response.ProjectRoleNonExistent)
		}
		// 获取所有状态
		roleWhereIn := extend.SliceOperator(maputil.Keys(roles), role, "|")
		tx.Where("role IN ?", roleWhereIn)
	}
	// 查询用户
	if err := tx.Find(&members).Error; err != nil {
		return nil, exception.ErrorHandle(err, response.ProjectMemberQueryFail)
	}

	return members, nil
}

// ShouldProjectRoles 获取角色列表并且判断role是否在列表内
func (receiver *ProjectMemberLogic) ShouldProjectRoles(role int) map[int]string {
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
func (receiver *ProjectMemberLogic) Transfer(projectId uint, transferor, recipient uint64) error {
	// 如果两个id相同，不执行操作
	if transferor == recipient {
		return nil
	}

	var transferorMember model.ProjectMember // 移交人
	// 查询旧项目负责人
	err := receiver.Db.Where(&model.ProjectMember{
		ProjectId: projectId,
		UserId:    transferor,
	}).First(&transferorMember).Error

	// 没有查询到记录
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 返回负责人不存在的错误
		return exception.NewException(response.ProjectLeaderNotExist)
	}
	// 初始化状态修改器
	stateModifier := state.NewModifier(int(transferorMember.Role))
	// 判断该用户是否为负责人
	if stateModifier.NotExist(constant.ProjectLeader) {
		return exception.NewException(response.MemberNotProjectLeader)
	}

	return receiver.Db.Transaction(func(tx *gorm.DB) error {
		var recipientMember model.ProjectMember // 接收人
		/* 删除角色 Start */
		// 去除当前角色
		transferorMember.Role = int8(stateModifier.Detach(constant.ProjectLeader))
		if transferorMember.Role <= 0 {
			// 无角色，删除
			if err := tx.Delete(&transferorMember).Error; err != nil {
				return err
			}
		} else {
			// 修改成员角色
			if err := tx.Save(&transferorMember).Error; err != nil {
				return err
			}
		}
		/* 删除角色 End */

		/* 新增或给成员附加负责人角色 Start */
		// 先判断用户是否存在
		if !NewUserLogic(receiver.ctx).UserExist(recipient) {
			return exception.NewException(response.UserNotFound)
		}
		err := tx.Where(&model.ProjectMember{
			ProjectId: projectId,
			UserId:    recipient,
		}).First(&recipientMember).Error
		// 没有查询到记录
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 该用户不在项目内，新增一个
			return tx.Model(&model.ProjectMember{}).Create(&model.ProjectMember{
				ProjectId: projectId,
				UserId:    recipient,
				Role:      int8(constant.ProjectLeader), // 需要转换成int8
			}).Error
		} else {
			// 成员存在，给ta附加一个负责人角色
			recipientMember.Role = int8(state.NewModifier(int(recipientMember.Role)).Attach(constant.ProjectLeader))
			return tx.Save(&recipientMember).Error
		}
		/* 新增或给成员附加负责人角色 End */
	})
}

// ExistMember 项目是否存在指定角色的成员
func (receiver *ProjectMemberLogic) ExistMember(projectId uint, role int) bool {
	var count int64
	// 获取所有角色
	roles := constant.GetTaskRoles()
	// 获取所有可能的状态
	roleWhereIn := state.NewModifier(role).Contained(maputil.Keys(roles))
	// 查询
	receiver.Db.Model(&model.ProjectMember{}).
		Where("project_id = ?", projectId).
		Where("role IN ?", roleWhereIn).
		Count(&count)

	return count > 0
}

// GetLeader 获取负责人
func (receiver *ProjectMemberLogic) GetLeader(projectId uint) (*model.ProjectMember, error) {
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
