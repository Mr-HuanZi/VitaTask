package service

import (
	"VitaTaskGo/internal/api/data"
	"VitaTaskGo/internal/api/model/dto"
	"VitaTaskGo/internal/api/model/vo"
	"VitaTaskGo/internal/pkg"
	"VitaTaskGo/internal/pkg/workflow"
	"VitaTaskGo/internal/repo"
	"VitaTaskGo/pkg/db"
	"VitaTaskGo/pkg/exception"
	"VitaTaskGo/pkg/response"
	"VitaTaskGo/pkg/time_tool"
	"errors"
	"github.com/duke-git/lancet/v2/convertor"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
)

type WorkflowService struct {
	Db   *gorm.DB
	ctx  *gin.Context
	repo repo.WorkflowRepo
}

func NewWorkflowService(tx *gorm.DB, ctx *gin.Context) *WorkflowService {
	return &WorkflowService{
		Db:   tx,
		ctx:  ctx,
		repo: data.NewWorkflowRepo(tx, ctx),
	}
}

func (r *WorkflowService) Initiate(post dto.WorkflowInitiateDto) (*repo.Workflow, error) {
	var (
		engine *workflow.Engine
		err    error
	)

	// 创建引擎对象
	engine, err = workflow.Create(r.Db, r.ctx, post.TypeId)
	if err != nil {
		if errors.Is(err, workflow.ErrWorkflowTypeNotExist) {
			return nil, exception.NewException(response.WorkflowTypeNotExist)
		}
		return nil, exception.ErrorHandle(err, response.SystemFail)
	}

	// 将struct转换成map
	toMap, err := convertor.StructToMap(post)
	if err != nil {
		return nil, exception.ErrorHandle(err, response.SystemFail)
	}

	// 设置表单数据
	engine.SetFormData(toMap)
	// 发起工作流
	err = engine.Initiate()
	return engine.GetWorkflowInfo(), exception.ErrorHandle(err, response.SystemFail)
}

func (r *WorkflowService) ExamineApprove(post dto.WorkflowExamineApproveDto) (*repo.Workflow, error) {
	var (
		engine *workflow.Engine
		err    error
	)

	// 创建引擎对象
	engine, err = workflow.Open(r.Db, r.ctx, post.Id)
	if err != nil {
		if errors.Is(err, workflow.ErrWorkflowTypeNotExist) {
			return nil, exception.NewException(response.WorkflowTypeNotExist)
		}
		if errors.Is(err, workflow.ErrWorkflowNotExist) {
			return nil, exception.NewException(response.WorkflowNotExist)
		}
		return nil, exception.ErrorHandle(err, response.SystemFail)
	}

	// 将struct转换成map
	toMap, err := convertor.StructToMap(post)
	if err != nil {
		return nil, exception.ErrorHandle(err, response.SystemFail)
	}

	// 设置表单数据
	engine.SetFormData(toMap)
	// 执行审批
	err = engine.ExamineApprove()
	return engine.GetWorkflowInfo(), exception.ErrorHandle(err, response.SystemFail)
}

// PageList 分页列表
func (r *WorkflowService) PageList(query dto.WorkflowListQueryDto) (*dto.PagedResult[repo.Workflow], error) {
	workflowRepo := data.NewWorkflowRepo(r.Db, r.ctx)
	workflowNodeRepo := data.NewWorkflowNodeRepo(r.Db, r.ctx)
	workflowTypeRepo := data.NewWorkflowTypeRepo(r.Db, r.ctx)

	// 获取所有非系统级的工作流类型
	if !query.System {
		typeIds, err := workflowTypeRepo.GetNotSystemIds()
		if err != nil || len(typeIds) <= 0 {
			return pkg.PagedResult[repo.Workflow](nil, 0, int64(query.Page)), exception.ErrorHandle(err, response.DbQueryError, "列表查询失败: ")
		}
		if len(query.TypeId) > 0 && slice.ContainSubSlice(typeIds, query.TypeId) {
			// 使用 query.TypeId 的值，所以这里不执行任何操作
		} else {
			query.TypeId = typeIds
		}
	}

	l, total, err := workflowRepo.PageList(query)
	if err != nil {
		return pkg.PagedResult[repo.Workflow](nil, 0, int64(query.Page)), exception.ErrorHandle(err, response.DbQueryError, "列表查询失败: ")
	}

	for i, item := range l {
		// 获取节点数据
		node, err := workflowNodeRepo.GetAppointNode(item.TypeId, item.Node)
		if err == nil {
			l[i].NodeInfo = node
		}

		// 给状态赋值
		for ii, s := range workflow.StatusMap {
			if item.Status == s {
				l[i].StatusText = ii
			}
		}
	}

	return pkg.PagedResult(l, total, int64(query.Page)), exception.ErrorHandle(err, response.DbQueryError, "列表查询失败: ")
}

// Detail 工作流详情
func (r *WorkflowService) Detail(id uint) (*vo.WorkflowDetailVo, error) {
	// 实例化VO
	workflowDetailVo := new(vo.WorkflowDetailVo)
	// 实例化repo
	workflowRepo := data.NewWorkflowRepo(r.Db, r.ctx)
	workflowNodeRepo := data.NewWorkflowNodeRepo(r.Db, r.ctx)
	workflowOperatorRepo := data.NewWorkflowOperatorRepo(r.Db, r.ctx)
	workflowTypeRepo := data.NewWorkflowTypeRepo(r.Db, r.ctx)

	// 查询工作流详情
	workflowInfo, err := workflowRepo.Get(id)
	if err != nil {
		return nil, exception.ErrorHandle(err, response.DbQueryError, "详情查询失败: ")
	}
	workflowDetailVo.Workflow = workflowInfo

	// 给状态赋值
	for i, v := range workflow.StatusMap {
		if workflowInfo.Status == v {
			workflowInfo.StatusText = i
		}
	}

	// 查询工作流类型详情
	workflowType, typeErr := workflowTypeRepo.Get(workflowInfo.TypeId)
	if typeErr != nil {
		return nil, exception.ErrorHandle(typeErr, response.DbQueryError, "查询类型失败: ")
	}
	workflowDetailVo.WorkflowType = workflowType

	// 查询该工作流当前节点
	if workflowInfo.Node > 0 {
		node, nodeErr := workflowNodeRepo.GetAppointNode(workflowInfo.TypeId, workflowInfo.Node)
		if nodeErr != nil {
			return nil, exception.ErrorHandle(nodeErr, response.DbQueryError, "查询节点失败: ")
		}
		// 节点
		workflowDetailVo.Node = new(vo.WorkflowNodeVo)
		workflowDetailVo.Node.Node = node.Node
		workflowDetailVo.Node.Name = node.Name
		workflowDetailVo.Node.Action = node.Action
		workflowDetailVo.Node.ActionValue = node.ActionValue
		workflowDetailVo.Node.Everyone = node.Everyone
	}

	// 查询当前节点操作人
	operators, operatorsErr := workflowOperatorRepo.GetWorkflowOperatorByNode(workflowInfo.ID, workflowInfo.Node)
	if operatorsErr != nil {
		if errors.Is(operatorsErr, gorm.ErrRecordNotFound) {
			// 操作人
			workflowDetailVo.Operators = nil
		} else {
			return nil, exception.ErrorHandle(operatorsErr, response.DbQueryError, "查询节点操作人失败: ")
		}
	} else {
		workflowDetailVo.Operators = operators
	}

	return workflowDetailVo, nil
}

func (r *WorkflowService) TypeAdd(post dto.WorkflowTypeDto) (*repo.WorkflowType, error) {
	if len(post.OnlyName) <= 0 {
		// 唯一标志为必填项
		return nil, exception.NewException(response.WorkflowTypeOnlyNameEmpty)
	}

	workflowTypeRepo := data.NewWorkflowTypeRepo(r.Db, r.ctx)
	// 检查 OnlyName 是否有重复
	if workflowTypeRepo.ExistByOnlyName(post.OnlyName) {
		// 有记录，说明存在相同的
		return nil, exception.NewException(response.WorkflowTypeOnlyNameRepeat)
	}

	// 创建新对象
	newData := &repo.WorkflowType{
		Name:       post.Name,
		OrgId:      post.OrgId,
		OnlyName:   post.OnlyName,
		Illustrate: post.Illustrate,
	}
	// 是否系统级
	if post.System {
		newData.System = 1
	} else {
		newData.System = 0
	}

	saveErr := workflowTypeRepo.Create(newData)
	return newData, exception.ErrorHandle(saveErr, response.WorkflowTypeCreateFail)
}

func (r *WorkflowService) TypeUpdate(post dto.WorkflowTypeDto) (*repo.WorkflowType, error) {
	workflowTypeRepo := data.NewWorkflowTypeRepo(r.Db, r.ctx)

	// 获取记录
	one, err := workflowTypeRepo.Get(post.ID)
	if err != nil {
		return nil, db.FirstQueryErrorHandle(err, response.WorkflowTypeNotExist)
	}

	// 不允许修改OnlyName
	one.Name = post.Name
	one.OrgId = post.OrgId
	one.Illustrate = post.Illustrate
	// 是否系统级
	if post.System {
		one.System = 1
	} else {
		one.System = 0
	}

	saveErr := workflowTypeRepo.Save(one)
	return one, exception.ErrorHandle(saveErr, response.WorkflowTypeUpdateFail)
}

func (r *WorkflowService) TypeList(query dto.WorkflowTypeQueryDto) (*dto.PagedResult[repo.WorkflowType], error) {
	var (
		queryBo dto.WorkflowTypeQueryBo
	)

	workflowTypeRepo := data.NewWorkflowTypeRepo(r.Db, r.ctx)

	// 搜索处理
	queryBo.UintId = query.UintId
	queryBo.PagingQuery = query.PagingQuery
	queryBo.DeletedQuery = query.DeletedQuery
	queryBo.OnlyName = query.OnlyName
	if len(query.Name) > 0 {
		queryBo.Name = query.Name
	} else if len(query.Title) > 0 {
		queryBo.Name = query.Title
	}
	if len(query.CreateTime) >= 2 {
		createTimeRange, err := time_tool.ParseStartEndTimeToUnix(query.CreateTime, time.DateOnly, "milli")
		if err != nil {
			return pkg.PagedResult[repo.WorkflowType](nil, 0, int64(query.Page)), exception.ErrorHandle(err, response.TimeParseFail)
		}

		queryBo.CreateTime = createTimeRange
	}

	l, total, err := workflowTypeRepo.PageList(queryBo)

	return pkg.PagedResult(l, total, int64(query.Page)), exception.ErrorHandle(err, response.DbQueryError, "列表查询失败: ")
}

func (r *WorkflowService) TypeDelete(id uint) error {
	workflowTypeRepo := data.NewWorkflowTypeRepo(r.Db, r.ctx)
	_, err := workflowTypeRepo.Get(id)
	if err != nil {
		return db.FirstQueryErrorHandle(err, response.WorkflowTypeNotExist)
	}

	return exception.ErrorHandle(workflowTypeRepo.Delete(id), response.WorkflowTypeDeleteFail)
}

// TypeDetail 获取工作流类型详情
func (r *WorkflowService) TypeDetail(id uint) (*repo.WorkflowType, error) {
	one, err := data.NewWorkflowTypeRepo(r.Db, r.ctx).Get(id)
	return one, db.FirstQueryErrorHandle(err, response.WorkflowTypeNotExist)
}

// TypeDetailByOnlyName 获取工作流类型详情
func (r *WorkflowService) TypeDetailByOnlyName(onlyName string) (*repo.WorkflowType, error) {
	one, err := data.NewWorkflowTypeRepo(r.Db, r.ctx).GetByOnlyName(onlyName)
	return one, db.FirstQueryErrorHandle(err, response.WorkflowTypeNotExist)
}

// TypeOptions 获取Label+Value格式的工作流类型列表
func (r *WorkflowService) TypeOptions(keyWords string, system bool) ([]dto.UniversalSimpleList[uint], error) {
	workflowTypeRepo := data.NewWorkflowTypeRepo(r.Db, r.ctx)
	l, err := workflowTypeRepo.GetOptions(keyWords, system)
	if err != nil {
		return nil, db.FirstQueryErrorHandle(err, response.WorkflowTypeNotExist)
	}

	s := make([]dto.UniversalSimpleList[uint], len(l))
	// 因为不是数字键值，所以这里需要另外的变量来做下标
	i := 0
	for _, v := range l {
		s[i] = dto.UniversalSimpleList[uint]{
			Label: v.Name,
			Value: v.ID,
		}
		i++
	}
	return s, nil
}

func (r *WorkflowService) NodeAdd(post dto.WorkflowNodeDto) (*repo.WorkflowNode, error) {
	workflowTypeRepo := data.NewWorkflowTypeRepo(r.Db, r.ctx)
	workflowNodeRepo := data.NewWorkflowNodeRepo(r.Db, r.ctx)

	if post.TypeId <= 0 {
		// 缺少工作流类型ID
		return nil, exception.NewException(response.WorkflowTypeNotExist)
	}
	// 获取记录
	typeData, err := workflowTypeRepo.Get(post.TypeId)
	if err != nil {
		return nil, db.FirstQueryErrorHandle(err, response.WorkflowTypeNotExist)
	}

	// 如果设置的节点序号小于等于0则将其设置为1
	if post.Node <= 0 {
		post.Node = 1
	}

	// 创建新对象
	saveData := &repo.WorkflowNode{
		TypeId:      typeData.ID,
		Node:        post.Node,
		Name:        post.Name,
		Action:      post.Action,
		ActionValue: post.ActionValue,
	}
	createErr := workflowNodeRepo.Create(saveData)
	return nil, exception.ErrorHandle(createErr, response.WorkflowNodeCreateFail)
}

func (r *WorkflowService) NodeUpdate(post dto.WorkflowNodeDto) (*repo.WorkflowNode, error) {
	workflowNodeRepo := data.NewWorkflowNodeRepo(r.Db, r.ctx)

	// 获取节点记录
	nodeData, err := workflowNodeRepo.Get(post.ID)
	if err != nil {
		return nil, db.FirstQueryErrorHandle(err, response.WorkflowNodeNotExist)
	}

	// 如果设置的节点序号小于等于0则视为不修改
	if post.Node <= 0 {
		post.Node = nodeData.Node
	}

	// 修改数据
	// 不允许修改 TypeId
	nodeData.Name = post.Name
	nodeData.Node = post.Node
	nodeData.Action = post.Action
	nodeData.ActionValue = post.ActionValue

	saveErr := workflowNodeRepo.Save(nodeData)
	return nil, exception.ErrorHandle(saveErr, response.WorkflowNodeUpdateFail)
}

func (r *WorkflowService) NodeDelete(id uint) error {
	workflowNodeRepo := data.NewWorkflowNodeRepo(r.Db, r.ctx)
	_, err := workflowNodeRepo.Get(id)
	if err != nil {
		return db.FirstQueryErrorHandle(err, response.WorkflowNodeNotExist)
	}

	return exception.ErrorHandle(workflowNodeRepo.Delete(id), response.WorkflowNodeDeleteFail)
}

func (r *WorkflowService) NodeList(query dto.WorkflowNodeQueryDto) (*dto.PagedResult[repo.WorkflowNode], error) {
	var (
		queryBo dto.WorkflowNodeQueryBo
	)

	workflowNodeRepo := data.NewWorkflowNodeRepo(r.Db, r.ctx)

	// 搜索处理
	queryBo.UintId = query.UintId
	queryBo.PagingQuery = query.PagingQuery
	queryBo.DeletedQuery = query.DeletedQuery
	queryBo.Action = query.Action
	queryBo.TypeId = query.TypeId
	if len(query.Name) > 0 {
		queryBo.Name = query.Name
	} else if len(query.Title) > 0 {
		queryBo.Name = query.Title
	}
	if len(query.CreateTime) >= 2 {
		createTimeRange, err := time_tool.ParseStartEndTimeToUnix(query.CreateTime, time.DateOnly, "milli")
		if err != nil {
			return pkg.PagedResult[repo.WorkflowNode](nil, 0, int64(query.Page)), exception.ErrorHandle(err, response.TimeParseFail)
		}

		queryBo.CreateTime = createTimeRange
	}

	l, total, err := workflowNodeRepo.PageList(queryBo)

	return pkg.PagedResult(l, total, int64(query.Page)), exception.ErrorHandle(err, response.DbQueryError, "列表查询失败: ")
}

// NodeTypeAll 获取指定工作流模板的所有节点(无分页)
func (r *WorkflowService) NodeTypeAll(id uint) ([]repo.WorkflowNode, error) {
	workflowNodeRepo := data.NewWorkflowNodeRepo(r.Db, r.ctx)
	// 获取该工作流类型的所有节点配置
	workflowNodes, nodeErr := workflowNodeRepo.GetTypeAll(id)
	if nodeErr != nil {
		return nil, exception.ErrorHandle(nodeErr, response.DbQueryError, "查询节点失败: ")
	}
	return workflowNodes, nil
}

func (r *WorkflowService) Actions() []dto.UniversalSimpleList[string] {
	kv := workflow.GetAllActionName()
	s := make([]dto.UniversalSimpleList[string], len(kv))
	// 因为不是数字键值，所以这里需要另外的变量来做下标
	i := 0
	for k, v := range kv {
		s[i] = dto.UniversalSimpleList[string]{
			Label: v,
			Value: k,
		}
		i++
	}
	return s
}

// StatusList 工作流状态列表
// 适配Antd Pro表格格式
func (r *WorkflowService) StatusList() map[int]map[string]string {
	return workflow.StatusEnum
}

func (r *WorkflowService) LogPageLists(query dto.WorkflowLogQueryDto) (*dto.PagedResult[vo.WorkflowLogVo], error) {
	var (
		queryBo dto.WorkflowLogQueryBo
	)

	workflowLogRepo := data.NewWorkflowLogRepo(r.Db, r.ctx)
	workflowNodeRepo := data.NewWorkflowNodeRepo(r.Db, r.ctx)
	workflowRepo := data.NewWorkflowRepo(r.Db, r.ctx)

	// 搜索处理
	queryBo.UintId = query.UintId
	queryBo.PagingQuery = query.PagingQuery
	queryBo.WorkflowId = query.WorkflowId
	queryBo.Node = query.Node
	queryBo.Action = query.Action

	// 操作人
	if query.Operator > 0 {
		queryBo.Operator = []uint64{query.Operator}
	}

	if len(query.CreateTime) >= 2 {
		createTimeRange, err := time_tool.ParseStartEndTimeToUnix(query.CreateTime, time.DateOnly, "milli")
		if err != nil {
			return pkg.PagedResult[vo.WorkflowLogVo](nil, 0, int64(query.Page)), exception.ErrorHandle(err, response.TimeParseFail)
		}

		queryBo.CreateTime = createTimeRange
	}

	l, total, err := workflowLogRepo.PageList(queryBo)

	workflowLogVo := make([]vo.WorkflowLogVo, len(l))
	for i, item := range l {
		itemVo := vo.WorkflowLogVo{}
		err := convertor.CopyProperties(&itemVo, item)
		if err != nil {
			return pkg.PagedResult[vo.WorkflowLogVo](nil, 0, int64(query.Page)), exception.ErrorHandle(err, response.SystemFail, "Vo数据转换失败")
		}

		// 查询工作流详情
		workflowInfo, err := workflowRepo.Get(item.WorkflowId)
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return pkg.PagedResult[vo.WorkflowLogVo](nil, 0, int64(query.Page)), exception.ErrorHandle(err, response.DbQueryError, "获取工作流Info失败")
			}
		}

		if workflowInfo != nil {
			// 获取节点数据
			node, err := workflowNodeRepo.GetAppointNode(workflowInfo.TypeId, item.Node)
			if err == nil {
				itemVo.NodeInfo = node
			}
		}

		workflowLogVo[i] = itemVo
	}

	return pkg.PagedResult(workflowLogVo, total, int64(query.Page)), exception.ErrorHandle(err, response.DbQueryError, "列表查询失败: ")
}

func (r *WorkflowService) Footprint(id uint) ([]vo.WorkflowFootprintVo, error) {
	workflowLogRepo := data.NewWorkflowLogRepo(r.Db, r.ctx)
	workflowNodeRepo := data.NewWorkflowNodeRepo(r.Db, r.ctx)
	workflowRepo := data.NewWorkflowRepo(r.Db, r.ctx)

	// 查询工作流详情
	workflowInfo, err := workflowRepo.Get(id)
	if err != nil {
		return nil, exception.ErrorHandle(err, response.WorkflowNotExist)
	}

	// 获取该工作流类型的所有节点配置
	workflowNodes, nodeErr := workflowNodeRepo.GetTypeAll(workflowInfo.TypeId)
	if nodeErr != nil {
		return nil, exception.ErrorHandle(nodeErr, response.DbQueryError, "查询节点失败: ")
	}

	// 获取该工作流的所有操作记录
	workflowLogs, logErr := workflowLogRepo.GetWorkflowAll(workflowInfo.ID)
	if logErr != nil {
		return nil, exception.ErrorHandle(nodeErr, response.DbQueryError, "查询日志失败: ")
	}

	// 遍历日志，找出第一个 action 是 initiate 的
	// 因为 workflowLogs 在查询时已经按 create_time 倒序了，这里直接遍历即可
	filteredLog := make([]repo.WorkflowLog, 99)
	for _, log := range workflowLogs {
		filteredLog = append(filteredLog, log)
		if log.Action == workflow.Initiate {
			break
		}
	}

	// 创建Vo数据
	footprintVo := make([]vo.WorkflowFootprintVo, len(workflowNodes))
	// 遍历节点数据
	for i, node := range workflowNodes {
		item := vo.WorkflowFootprintVo{
			Node: node.Node,
			Name: node.Name,
			// 是否是当前节点
			Curr: workflowInfo.Node == node.Node,
		}
		// 遍历操作记录
		for _, log := range filteredLog {
			if log.Node == node.Node {
				// 操作说明(多条只保留最后一个)
				item.Explain = log.Explain
				item.Time = log.CreateTime
				// 记录操作人
				if item.Operators == nil {
					// 如果没有初始化，就初始化一下
					item.Operators = make([]vo.WorkflowFootprintOperatorVo, 0)
				}

				item.Operators = append(item.Operators, vo.WorkflowFootprintOperatorVo{
					Uid:      log.Operator,
					Nickname: log.Nickname,
				})
			}
		}

		footprintVo[i] = item
	}

	return footprintVo, nil
}
