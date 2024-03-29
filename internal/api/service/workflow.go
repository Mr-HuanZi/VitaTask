package service

import (
	"VitaTaskGo/internal/api/data"
	"VitaTaskGo/internal/api/model/dto"
	"VitaTaskGo/internal/pkg"
	"VitaTaskGo/internal/pkg/workflow"
	"VitaTaskGo/internal/repo"
	"VitaTaskGo/pkg/db"
	"VitaTaskGo/pkg/exception"
	"VitaTaskGo/pkg/response"
	"VitaTaskGo/pkg/time_tool"
	"github.com/duke-git/lancet/v2/convertor"
	"github.com/duke-git/lancet/v2/maputil"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
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

func (r *WorkflowService) Initiate(post dto.WorkflowInitiateDto) error {
	var (
		engine *workflow.Engine
		err    error
	)

	// 创建引擎对象
	engine, err = workflow.Create(r.Db, r.ctx, post.TypeId)
	if err != nil {
		return exception.ErrorHandle(err, response.SystemFail)
	}

	// 将struct转换成map
	toMap, err := convertor.StructToMap(post)
	if err != nil {
		return exception.ErrorHandle(err, response.SystemFail)
	}

	// 设置表单数据
	engine.SetFormData(toMap)
	// 发起工作流
	err = engine.Initiate()
	return exception.ErrorHandle(err, response.SystemFail)
}

func (r *WorkflowService) ExamineApprove(post dto.WorkflowExamineApproveDto) error {
	var (
		engine *workflow.Engine
		err    error
	)

	// 创建引擎对象
	engine, err = workflow.Open(r.Db, r.ctx, post.WorkflowId)
	if err != nil {
		return exception.ErrorHandle(err, response.SystemFail)
	}

	// 将struct转换成map
	toMap, err := convertor.StructToMap(post)
	if err != nil {
		return exception.ErrorHandle(err, response.SystemFail)
	}

	// 设置表单数据
	engine.SetFormData(toMap)
	// 执行审批
	err = engine.ExamineApprove()
	return exception.ErrorHandle(err, response.SystemFail)
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

	// 工作流状态 转换为 数字字符串
	if status, ok := workflow.StatusMap[query.Status]; ok {
		query.Status = strconv.Itoa(status)
	}

	l, total, err := workflowRepo.PageList(query)
	if err != nil {
		return pkg.PagedResult[repo.Workflow](nil, 0, int64(query.Page)), exception.ErrorHandle(err, response.DbQueryError, "列表查询失败: ")
	}

	statusNames := maputil.Keys(workflow.StatusMap)

	for i, item := range l {
		// 获取节点数据
		node, err := workflowNodeRepo.GetAppointNode(item.TypeId, item.Node)
		if err == nil {
			l[i].NodeInfo = node
		}

		// 给状态赋值
		if len(statusNames)-1 > item.Status {
			l[i].StatusText = statusNames[item.Status]
		}
	}

	return pkg.PagedResult(l, total, int64(query.Page)), exception.ErrorHandle(err, response.DbQueryError, "列表查询失败: ")
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

func (r *WorkflowService) TypeDetail(id uint) (*repo.WorkflowType, error) {
	one, err := data.NewWorkflowTypeRepo(r.Db, r.ctx).Get(id)
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
func (r *WorkflowService) StatusList() map[string]map[string]string {
	return workflow.StatusEnum
}
