package logic

import (
	"VitaTaskGo/app/constant"
	"VitaTaskGo/app/exception"
	"VitaTaskGo/app/extend"
	"VitaTaskGo/app/extend/time_tool"
	"VitaTaskGo/app/extend/user"
	"VitaTaskGo/app/model"
	"VitaTaskGo/app/response"
	"VitaTaskGo/app/types"
	"VitaTaskGo/library/db"
	"VitaTaskGo/library/state"
	"fmt"
	"github.com/duke-git/lancet/v2/maputil"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"sort"
	"time"
)

type TaskLogic struct {
	Db  *gorm.DB
	ctx *gin.Context
}

func NewTaskLogic(ctx *gin.Context) *TaskLogic {
	return &TaskLogic{
		Db:  db.Db, // 赋予ORM实例
		ctx: ctx,   // 传递上下文
	}
}

// Lists 任务列表
func (receiver TaskLogic) Lists(query types.TaskListQuery) (*types.PagedResult[model.Task], error) {
	var (
		tasks []model.Task
		count int64
	)

	tx, queryErr := receiver.QueryHandle(query)
	if queryErr != nil {
		return nil, queryErr
	}
	tx.Model(&model.Task{}).Count(&count)

	err := tx.Model(&model.Task{}).
		Scopes(db.Paginate(&query.Page, &query.PageSize)).
		Preload("Project").
		Preload("Group").
		Preload("Member.UserInfo").
		Order("status ASC").Order("level DESC").Order("create_time DESC").
		Find(&tasks).Error
	if err != nil {
		tasks = make([]model.Task, 0)
		_ = exception.ErrorHandle(err, response.DbQueryError, "任务列表查询失败: ")
	}
	// 数据合并
	for i, task := range tasks {
		tasks[i].PlanTime = []int64{task.StartDate, task.EndDate}
		// 获取负责人
		for _, member := range task.Member {
			stateModifier := state.NewModifier(int(member.Role))
			if stateModifier.Exist(constant.TaskLeader) {
				tasks[i].Leader = member
				break
			}
		}
	}

	return extend.PagedResult[model.Task](tasks, count, int64(query.Page)), nil
}

// QueryHandle 查询处理
func (receiver TaskLogic) QueryHandle(query types.TaskListQuery) (*gorm.DB, error) {
	var tx *gorm.DB
	if query.Deleted {
		tx = receiver.Db.Unscoped() // 查询已删除的记录
	} else {
		tx = receiver.Db
	}

	if query.ID > 0 {
		tx = tx.Where("id = ?", query.ID)
	}

	// 列出当前用户所在的所有项目
	projectIds, err := NewProjectLogic(receiver.ctx).MyProjectIds()
	if err != nil {
		return nil, exception.ErrorHandle(err, response.DbQueryError)
	}
	if query.Project > 0 {
		if !slice.Contain(projectIds, query.Project) {
			// 当前用户不在项目内
			return nil, exception.NewException(response.MemberNotInProject)
		}
		tx = tx.Where("project_id = ?", query.Project)
	} else {
		// 只检索当前用户所属的项目列表
		tx = tx.Where("project_id IN ?", projectIds)
	}
	// 任务组搜索
	if query.GroupId > 0 {
		tx = tx.Where("group_id = ?", query.GroupId)
	}
	// 负责人
	if query.Leader > 0 {
		taskIds, err := NewTaskMemberLogic(receiver.ctx).GetTaskIdsByUsers([]uint64{query.Leader}, constant.TaskLeader)
		if err == nil {
			tx = tx.Where("id IN ?", taskIds)
		}
	}
	// 协助人
	if len(query.Collaborator) > 0 {
		taskIds, err := NewTaskMemberLogic(receiver.ctx).GetTaskIdsByUsers(query.Collaborator, constant.TaskMember)
		if err == nil {
			tx = tx.Where("id IN ?", taskIds)
		}
	}
	// 时间范围
	if len(query.CreateTime) >= 2 {
		createTimeRange, err := time_tool.ParseStartEndTimeToUnix(query.CreateTime, time.DateOnly, "milli")
		if err == nil {
			tx = tx.Where(
				"create_time BETWEEN ? AND ?",
				createTimeRange[0],
				createTimeRange[1],
			)
		}
	}
	// 计划时间范围
	if len(query.PlanTime) >= 2 {
		planTimeRange, err := time_tool.ParseStartEndTimeToUnix(query.PlanTime, time.DateOnly, "milli")
		if err == nil {
			// todo 改为 "(start_date BETWEEN ? AND ?) OR (end_date BETWEEN ? AND ?)"是否会更合理一些？
			tx = tx.Where(
				"start_date >= ? AND end_date <= ?",
				planTimeRange[0],
				planTimeRange[1],
			)
		}
	}
	// 标题搜索
	titleQuery := ""
	if query.Title != "" {
		titleQuery = query.Title
	} else if query.Name != "" {
		titleQuery = query.Name
	}
	if titleQuery != "" {
		tx = tx.Where("title LIKE ?", "%"+titleQuery+"%")
	}
	return tx, nil
}

// Create 创建任务
func (receiver TaskLogic) Create(post types.TaskCreateForm) (*model.Task, error) {
	// 获取当前用户
	currUser, err := user.CurrUser(receiver.ctx)
	if err != nil {
		return nil, err
	}
	// 项目是否存在
	projectLogic := NewProjectLogic(receiver.ctx)
	_, projectErr := projectLogic.GetOneProject(post.ProjectId)
	if projectErr != nil {
		return nil, projectErr
	}
	// 是否属于项目成员
	if !NewProjectMemberLogic(receiver.ctx).InProject(post.ProjectId, currUser.ID, 0) {
		return nil, exception.NewException(response.MemberNotInProject, "您不属于项目成员")
	}

	// 项目是否归档
	if projectLogic.Archived(post.ProjectId) {
		return nil, exception.NewException(response.ProjectArchived)
	}
	// 创建任务模型
	task, err := receiver.NewTask(post)
	if err != nil {
		return nil, err
	}

	transactionErr := receiver.Db.Transaction(func(tx *gorm.DB) error {
		// 保存任务数据
		if err := tx.Create(task).Error; err != nil {
			return err
		}
		taskMemberLogic := NewTaskMemberLogic(receiver.ctx)
		taskMemberLogic.Db = tx // DB拷贝
		if post.Leader > 0 {
			// 判断负责人是否存在
			if !NewUserLogic(receiver.ctx).UserExist(post.Leader) {
				return exception.NewException(response.UserNotFound)
			}
			// 保存负责人
			err = taskMemberLogic.Bind(task.ID, []uint64{post.Leader}, constant.TaskLeader)
		}
		// 保存创建人
		if err := taskMemberLogic.Bind(task.ID, []uint64{currUser.ID}, constant.TaskCreator); err != nil {
			return err
		}

		// 保存协助人
		if len(post.Collaborator) > 0 {
			err := taskMemberLogic.Bind(task.ID, post.Collaborator, constant.TaskMember)
			if err != nil {
				return err
			}
		}
		/* 创建任务对话 Start */
		dialogLogic := NewDialogLogic(receiver.ctx)
		dialogLogic.Db = tx // DB拷贝
		// 获取成员
		members, err := taskMemberLogic.GetAllMember(task.ID)
		if err != nil {
			return err
		}
		// 获取成员UID
		var memberIds = make([]uint64, len(members))
		for i, member := range members {
			memberIds[i] = member.UserId
		}
		dialog, err := dialogLogic.Create(fmt.Sprintf("任务[%s]聊天", task.Title), constant.DialogTypeTask, memberIds)
		if err != nil {
			return err
		}
		/* 创建任务对话 End */

		// 关联对话ID到任务
		err = tx.Model(&task).Update("dialog_id", dialog.ID).Error
		if err != nil {
			return err
		}

		// 记录日志
		taskLogLogic := NewTaskLogLogic(receiver.ctx)
		taskLogLogic.Orm = tx // DB拷贝
		_, err = taskLogLogic.Add(types.TaskLogForm{
			TaskId:      task.ID,
			OperateType: constant.TaskOperatorCreate,
			Message:     "创建了任务",
		})
		return err
	})

	if err := exception.ErrorHandle(transactionErr, response.TaskCreateFail, "创建任务失败: "); err != nil {
		return nil, err
	}

	return task, nil
}

// NewTask 获取一个新对象
func (receiver TaskLogic) NewTask(data types.TaskCreateForm) (*model.Task, error) {
	task := &model.Task{
		ProjectId: data.ProjectId,
		GroupId:   data.GroupId,
		Title:     data.Title,
		Describe:  data.Describe,
		Status:    0, // 新任务是未完成的
		Level:     data.Level,
	}

	// 时间范围
	if len(data.PlanTime) >= 2 {
		planTime, err := time_tool.ParseTimeRangeToUnix(data.PlanTime, time.DateOnly, "milli")
		if err != nil {
			return nil, err
		}
		task.StartDate = planTime[0] // 开始时间-毫秒
		// 解析时间戳
		t := time.Unix(planTime[1]/1e3, 0)
		// 把结束时间调整为当天最后1秒
		task.EndDate = time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, t.Location()).UnixMilli() // 结束时间-毫秒
	}

	return task, nil
}

// Detail 获取任务详情
func (receiver TaskLogic) Detail(taskId uint) (*model.Task, error) {
	task, err := receiver.GetOne(taskId)
	if err != nil {
		return nil, err
	}
	task.PlanTime = []int64{task.StartDate, task.EndDate}
	// 获取负责人
	for _, member := range task.Member {
		stateModifier := state.NewModifier(int(member.Role))
		// 取出任务负责人
		if stateModifier.Exist(constant.TaskLeader) {
			task.Leader = member
		}
		// 取出任务创建者
		if stateModifier.Exist(constant.TaskCreator) {
			task.Creator = member
		}
		// 取出普通成员/协助者
		if stateModifier.Exist(constant.TaskMember) {
			// Collaborator 是指针类型，可以直接append
			task.Collaborator = append(task.Collaborator, member)
		}
	}
	return task, nil
}

// GetOne 获取任务的简单详情
func (receiver TaskLogic) GetOne(taskId uint) (*model.Task, error) {
	var task *model.Task
	err := receiver.Db.Preload("Project").
		Preload("Group").
		Preload("Member.UserInfo").
		Order("status ASC").Order("level DESC").Order("create_time DESC").
		First(&task, taskId).Error

	if err != nil {
		return nil, exception.ErrorHandle(err, response.TaskNotExist, "获取单条项目记录失败: ")
	}
	if task.ID <= 0 {
		return nil, exception.NewException(response.TaskNotExist)
	}

	return task, nil
}

// Roles 获取所有角色
func (receiver TaskLogic) Roles() map[int]string {
	roles := constant.GetTaskRoles()
	return roles
}

// Status 获取所有状态
func (receiver TaskLogic) Status() []types.TaskStatusVo {
	statusMap := constant.GetTaskStatus()
	// 获取Keys
	statusMapKeys := maputil.Keys(statusMap)
	// 对Keys升序
	sort.Ints(statusMapKeys)
	// 创建Vo切片
	taskStatusVO := make([]types.TaskStatusVo, len(statusMap))
	i := 0
	// 遍历排序好的Keys
	for _, t := range statusMapKeys {
		status, ok := statusMap[t]
		if !ok {
			continue
		}
		statusDrop := "processing" // 状态点
		switch t {
		case constant.TaskStatusProcessing:
			statusDrop = "processing"
			break
		case constant.TaskStatusCompleted:
			statusDrop = "success"
			break
		case constant.TaskStatusArchived:
			statusDrop = "default"
			break
		}
		// 转成VO
		taskStatusVO[i] = types.TaskStatusVo{
			Label:  status,
			Value:  t,
			Status: statusDrop,
		}
		i++
	}

	return taskStatusVO
}

// ChangeStatus 更改任务状态
func (receiver TaskLogic) ChangeStatus(taskId uint, status int) error {
	task, err := receiver.GetOne(taskId)
	if err != nil {
		return err
	}
	// 项目是否归档
	if NewProjectLogic(receiver.ctx).Archived(task.ProjectId) {
		return exception.NewException(response.ProjectArchived)
	}
	// 获取状态Map
	statusMap := constant.GetTaskStatus()
	if _, ok := statusMap[status]; !ok {
		return exception.NewException(response.TaskStatusNotExist)
	}
	// 创建待修改数据的Map
	updates := make(map[string]interface{})
	// 保存要修改的状态值
	updates["status"] = status
	// 如果要对任务进行归档
	if status == constant.TaskStatusArchived {
		if task.Status != 1 {
			return exception.NewException(response.TaskStatusProcessing, "未完成的任务不能归档")
		}
		// 记录归档时间
		updates["archived_date"] = time.Now().UnixMilli()
	} else if status == constant.TaskStatusCompleted {
		// 如果是标记为已完成
		// 记录完成时间
		updates["complete_date"] = time.Now().UnixMilli()
	}
	// 更改状态
	err = receiver.Db.Model(&task).Updates(updates).Error
	if err != nil {
		return err
	}

	// 记录日志
	_, err = NewTaskLogLogic(receiver.ctx).Add(types.TaskLogForm{
		TaskId:      task.ID,
		OperateType: constant.TaskOperatorStatus,
		Message:     fmt.Sprintf("修改了任务状态为[%s]", statusMap[status]),
	})
	return err
}

// Update 更新任务
func (receiver TaskLogic) Update(taskId uint, post types.TaskCreateForm) (*model.Task, error) {
	// 项目是否存在
	projectLogic := NewProjectLogic(receiver.ctx)
	_, projectErr := projectLogic.GetOneProject(post.ProjectId)
	if projectErr != nil {
		return nil, projectErr
	}
	// 项目是否归档
	if projectLogic.Archived(post.ProjectId) {
		return nil, exception.NewException(response.ProjectArchived)
	}
	// 任务组是否存在
	if post.GroupId > 0 {
		if _, err := NewTaskGroupLogic(receiver.ctx).GetOne(post.GroupId); err != nil {
			return nil, err
		}
	}
	// 查询任务
	task, taskErr := receiver.GetOne(taskId)
	if taskErr != nil {
		return nil, taskErr
	}
	// 更新各个字段
	taskSave := map[string]interface{}{
		"project_id": post.ProjectId,
		"group_id":   post.GroupId,
		"title":      post.Title,
		"describe":   post.Describe,
		"level":      post.Level,
	}
	// 计划时间
	if len(post.PlanTime) >= 2 {
		planTime, err := time_tool.ParseTimeRangeToUnix(post.PlanTime, time.DateOnly, "milli")
		if err != nil {
			return nil, err
		}
		taskSave["start_date"] = planTime[0] // 开始时间
		// 把结束时间调整为当天最后1秒
		t := time.Unix(planTime[1]/1e3, 0)                                                                      // 解析时间戳
		taskSave["end_date"] = time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, t.Location()).UnixMilli() // 结束时间
	} else {
		// 清空计划时间
		taskSave["start_date"] = 0 // 开始时间
		taskSave["end_date"] = 0   // 结束时间
	}
	err := receiver.Db.Transaction(func(tx *gorm.DB) error {
		// 实例化成员服务
		taskMemberLogic := NewTaskMemberLogic(receiver.ctx)
		taskMemberLogic.Db = tx // DB拷贝
		// 实例化任务日志服务
		taskLogLogic := NewTaskLogLogic(receiver.ctx)
		taskLogLogic.Orm = tx // DB拷贝
		// 保存任务数据
		if err := tx.Model(&task).Updates(taskSave).Error; err != nil {
			return err
		}
		// 记录日志
		_, err := taskLogLogic.Add(types.TaskLogForm{
			TaskId:      task.ID,
			OperateType: constant.TaskOperatorUpdate,
			Message:     "修改了任务信息",
		})
		/* 保存负责人 Start */
		if post.Leader > 0 {
			// 转移负责人
			err := taskMemberLogic.Transfer(task.ID, post.Leader)
			if err != nil {
				if err == ErrSameTaskLeader {
					// 不执行任何操作
				} else {
					return err
				}
			} else {
				// 记录日志
				_, err = taskLogLogic.Add(types.TaskLogForm{
					TaskId:      task.ID,
					OperateType: constant.TaskOperatorChangeLeader,
					Message:     "变更了负责人",
				})
				if err != nil {
					return err
				}
			}
		} else {
			// 没有提供负责人参数，移除当前负责人
			err := taskMemberLogic.Remove(task.ID, nil, constant.TaskLeader)
			if err != nil {
				return err
			}
			// 记录日志
			_, err = taskLogLogic.Add(types.TaskLogForm{
				TaskId:      task.ID,
				OperateType: constant.TaskOperatorRemoveLeader,
				Message:     "移除负责人",
			})
			if err != nil {
				return err
			}
		}
		/* 保存负责人 End */

		/* 保存协作人 Start */
		// 先移除所有协作人
		err = taskMemberLogic.Remove(task.ID, nil, constant.TaskMember)
		if err != nil {
			return err
		}
		// 如果没有提供协作人参数，就认定为移除协作人
		// 重新绑定
		if len(post.Collaborator) > 0 {
			err := taskMemberLogic.Bind(task.ID, post.Collaborator, constant.TaskMember)
			if err != nil {
				return err
			}
		}
		// 记录日志
		_, err = taskLogLogic.Add(types.TaskLogForm{
			TaskId:      task.ID,
			OperateType: constant.TaskOperatorChangeCollaborator,
			Message:     "变更协作人",
		})
		if err != nil {
			return err
		}
		/* 保存协作人 End */

		/* 对话处理 Start */
		dialogLogic := NewDialogLogic(receiver.ctx)
		dialogLogic.Db = tx // DB拷贝
		// 获取任务成员成员
		taskMembers, err := taskMemberLogic.GetAllMember(task.ID)
		if err != nil {
			return err
		}
		// 获取成员UID
		var memberIds = make([]uint64, len(taskMembers))
		for i, member := range taskMembers {
			memberIds[i] = member.UserId
		}
		// 如果还没有对话，就创建一个
		if task.DialogId <= 0 {
			dialog, err := dialogLogic.Create(fmt.Sprintf("任务[%s]聊天", task.Title), constant.DialogTypeTask, memberIds)
			if err != nil {
				return err
			}
			// 关联对话ID到任务
			err = tx.Model(&task).Update("dialog_id", dialog.ID).Error
			if err != nil {
				return err
			}
		} else {
			// 获取对话成员
			dialogMembers, err := dialogLogic.GetMembers(task.DialogId)
			if err != nil {
				return err
			}
			// 提取UID
			var (
				taskMemberIds   = make([]uint64, len(taskMembers))
				dialogMemberIds = make([]uint64, len(dialogMembers))
			)
			for i, item := range taskMembers {
				taskMemberIds[i] = item.UserId
			}
			for i, item := range dialogMembers {
				dialogMemberIds[i] = item.UserId
			}
			// 在 对话中 但 不在任务成员中 的，移出对话
			err = dialogLogic.Exit(task.DialogId, slice.Difference(dialogMemberIds, taskMemberIds))
			if err != nil {
				return err
			}
			// 在 任务成员中 但 不在对话中 的，加入对话
			err = dialogLogic.Join(task.DialogId, slice.Difference(taskMemberIds, dialogMemberIds))
			if err != nil {
				return err
			}
		}
		/* 对话处理 End */
		return nil
	})
	return task, err
}

// Delete 删除任务
// 2023年3月11日 暂时不删除成员以及任务组
func (receiver TaskLogic) Delete(taskId uint) error {
	task, err := receiver.GetOne(taskId)
	if err != nil {
		return err
	}
	// 项目是否归档
	if NewProjectLogic(receiver.ctx).Archived(task.ProjectId) {
		return exception.NewException(response.ProjectArchived)
	}
	// 执行删除
	err = receiver.Db.Delete(&task).Error
	if err != nil {
		return err
	}
	// 记录日志
	_, err = NewTaskLogLogic(receiver.ctx).Add(types.TaskLogForm{
		TaskId:      task.ID,
		OperateType: constant.TaskOperatorDelete,
		Message:     "删除了任务",
	})
	return err
}

// Statistics 任务数量统计
// 已完成数量、未完成数量、按时完成数量、超时完成数量
func (receiver TaskLogic) Statistics(projectId uint) types.TaskStatistics {
	taskStatistics := types.TaskStatistics{
		Completed:  receiver.TaskNumber(projectId, []int{constant.TaskStatusCompleted, constant.TaskStatusArchived}),
		Processing: receiver.TaskNumber(projectId, []int{constant.TaskStatusProcessing}),
	}
	// 任务延误数量
	taskStatistics.FinishOnTime, taskStatistics.TimeoutCompletion = receiver.TaskDelayNumber(projectId)

	return taskStatistics
}

// TaskNumber 获取任务数量
func (receiver TaskLogic) TaskNumber(projectId uint, status []int) int64 {
	var count int64

	tx := receiver.Db.Model(model.Task{}).Where("project_id = ?", projectId).Where("status IN ?", status)

	// 查询数量并处理错误
	if exception.ErrorHandle(tx.Count(&count).Error, response.DbQueryError) != nil {
		return 0
	}
	return count
}

// TaskDelayNumber 任务延误数量
// 返回1-按时完成的任务数量 返回2-超时完成的数量
func (receiver TaskLogic) TaskDelayNumber(projectId uint) (int64, int64) {
	var (
		list                    []model.Task
		finishOnTimeNumber      int64
		timeoutCompletionNumber int64
	)
	// 查询已完成的任务(包括已归档的)
	tx := receiver.Db.Model(model.Task{}).Where("project_id = ?", projectId).Where("status IN ?", []int{constant.TaskStatusCompleted, constant.TaskStatusArchived})

	if err := tx.Find(&list).Error; err != nil {
		_ = exception.ErrorHandle(err, response.DbQueryError)
		return 0, 0
	}

	// 开始统计数量
	for _, item := range list {
		// 如果任务没有指定计划结束时间，算按时完成
		if item.EndDate <= 0 {
			finishOnTimeNumber += 1
			continue
		}

		// 完成时间 小于 计划结束时间
		if item.CompleteDate > 0 && item.CompleteDate < item.EndDate {
			finishOnTimeNumber += 1
			continue
		}

		// 其它情况，都算超时
		timeoutCompletionNumber += 1
	}

	return finishOnTimeNumber, timeoutCompletionNumber
}

// DailySituation 每日任务情况
// 统计最近7天(默认)的任务完成与新增情况
func (receiver TaskLogic) DailySituation(query types.DailySituationQuery) ([]types.DailySituationVo, error) {
	var startDate, endDate carbon.Carbon
	// 6天+1天的偏移，刚好是7天
	if query.StartDate == "" && query.EndDate == "" {
		// 没有提供开始和结束时间
		// 默认情况是6天内
		endDate = carbon.Now()
		startDate = carbon.Now().SubDays(6)
	} else if query.StartDate == "" {
		// 没有提供开始时间
		// 解析结束时间
		endDate = carbon.Parse(query.EndDate)
		// 获取6天前
		startDate = endDate.SubDays(6)
	} else if query.EndDate == "" {
		// 没有提供结束时间
		// 解析开始时间
		startDate = carbon.Parse(query.StartDate)
		// 获取6天后
		endDate = startDate.AddDays(6)
	} else {
		// 完整的开始与结束时间
		startDate = carbon.Parse(query.StartDate)
		endDate = carbon.Parse(query.EndDate)
	}

	// 当天的开始
	startDate = startDate.StartOfDay()
	// 当天的结束
	endDate = endDate.EndOfDay()
	// 计算天数差(正常情况下是正整数，两个时间如果反过来就是负整数)
	// 为什么+1？比如3月20日~3月28日，DiffInDays计算得到的时间差是8天(即20日~27日)，但实际应用中我们需要把28日也算上，所以要+1
	dayDiff := startDate.DiffInDays(endDate) + 1

	if dayDiff <= 0 {
		// 差值小于0表示开始时间大于结束时间
		return nil, exception.NewException(response.StartTimeGtEndTime)
	} else if dayDiff > 30 {
		// 时间跨度太长
		return nil, exception.NewException(response.TimeSpanTooLong)
	}
	logrus.Debugln(dayDiff, startDate.ToDateString(), endDate.ToDateString())
	// 遍历时间差
	var (
		i              int64 = 0
		statusWhere          = []int{constant.TaskStatusCompleted, constant.TaskStatusArchived}
		start, end     carbon.Carbon
		dailySituation = make([]types.DailySituationVo, 0)
	)

	for ; i < dayDiff; i++ {
		var (
			addQuantity        int64 = 0
			completedQuantity  int64 = 0
			incompleteQuantity int64 = 0
			err                error
		)

		if i <= 0 {
			start = startDate
			end = startDate.EndOfDay()
		} else {
			// + i天
			start = startDate.AddDays(int(i)).StartOfDay()
			end = start.EndOfDay()
		}
		logrus.Debugln(start.ToDateString(), end.ToDateString())
		// 当天新增的任务
		err = receiver.Db.Model(model.Task{}).
			Where("project_id = ?", query.ProjectId).
			Where("create_time BETWEEN ? AND ?", start.TimestampMilli(), end.TimestampMilli()).
			Count(&addQuantity).Error
		if err != nil {
			return nil, exception.ErrorHandle(err, response.DbQueryError)
		}
		// 当天已完成的任务
		err = receiver.Db.Model(model.Task{}).
			Where("project_id = ?", query.ProjectId).
			Where("status IN ?", statusWhere).
			Where("complete_date BETWEEN ? AND ?", start.TimestampMilli(), end.TimestampMilli()).
			Count(&completedQuantity).Error
		if err != nil {
			return nil, exception.ErrorHandle(err, response.DbQueryError)
		}
		// 当天未完成的任务
		err = receiver.Db.Model(model.Task{}).
			Where("project_id = ?", query.ProjectId).
			Where("status = ?", constant.TaskStatusProcessing).
			Count(&incompleteQuantity).Error
		if err != nil {
			return nil, exception.ErrorHandle(err, response.DbQueryError)
		}
		dailySituation = append(dailySituation,
			types.DailySituationVo{
				Label: "已完成",
				Date:  start.ToDateString(),
				Value: completedQuantity,
			},
			types.DailySituationVo{
				Label: "未完成",
				Date:  start.ToDateString(),
				Value: incompleteQuantity,
			},
			types.DailySituationVo{
				Label: "新增",
				Date:  start.ToDateString(),
				Value: addQuantity,
			},
		)
	}

	return dailySituation, nil
}
