package handle

import (
	"VitaTaskGo/internal/api/model/dto"
	"VitaTaskGo/internal/api/service"
	"VitaTaskGo/pkg/db"
	"VitaTaskGo/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type WorkflowApi struct {
}

func NewWorkflowApi() *WorkflowApi {
	return &WorkflowApi{}
}

func (r WorkflowApi) Initiate(ctx *gin.Context) {
	var post dto.WorkflowInitiateDto
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(
		http.StatusOK,
		response.Auto(service.NewWorkflowService(db.Db, ctx).Initiate(post)),
	)
}

func (r WorkflowApi) ExamineApprove(ctx *gin.Context) {
	var post dto.WorkflowExamineApproveDto
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(
		http.StatusOK,
		response.Auto(service.NewWorkflowService(db.Db, ctx).ExamineApprove(post)),
	)
}

func (r WorkflowApi) All(ctx *gin.Context) {
	var query dto.WorkflowListQueryDto
	if err := ctx.ShouldBindJSON(&query); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(
		http.StatusOK,
		response.Auto(service.NewWorkflowService(db.Db, ctx).AllPageList(query)),
	)
}

func (r WorkflowApi) ToDo(ctx *gin.Context) {
	var query dto.WorkflowListQueryDto
	if err := ctx.ShouldBindJSON(&query); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(
		http.StatusOK,
		response.Auto(service.NewWorkflowService(db.Db, ctx).ToDoList(query)),
	)
}

// Handled 我的已办工作流分页列表
func (r WorkflowApi) Handled(ctx *gin.Context) {
	var query dto.WorkflowListQueryDto
	if err := ctx.ShouldBindJSON(&query); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(
		http.StatusOK,
		response.Auto(service.NewWorkflowService(db.Db, ctx).HandledList(query)),
	)
}

// List 我发起的工作流分页列表
func (r WorkflowApi) List(ctx *gin.Context) {
	var query dto.WorkflowListQueryDto
	if err := ctx.ShouldBindJSON(&query); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(
		http.StatusOK,
		response.Auto(service.NewWorkflowService(db.Db, ctx).WorkflowList(query)),
	)
}

// Detail 工作流实例明细
func (r WorkflowApi) Detail(ctx *gin.Context) {
	id := ctx.Query("id")

	// 转换成uint
	idConv, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, response.Custom("缺少工作流实例ID参数", response.FormVerificationFailed, nil))
		return
	}

	ctx.JSON(
		http.StatusOK,
		response.Auto(service.NewWorkflowService(db.Db, ctx).Detail(uint(idConv))),
	)
}

func (r WorkflowApi) TypeAdd(ctx *gin.Context) {
	var post dto.WorkflowTypeDto
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(
		http.StatusOK,
		response.Auto(service.NewWorkflowService(db.Db, ctx).TypeAdd(post)),
	)
}

func (r WorkflowApi) TypeUpdate(ctx *gin.Context) {
	var post dto.WorkflowTypeDto

	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(
		http.StatusOK,
		response.Auto(service.NewWorkflowService(db.Db, ctx).TypeUpdate(post)),
	)
}

func (r WorkflowApi) TypeList(ctx *gin.Context) {
	var query dto.WorkflowTypeQueryDto
	if err := ctx.ShouldBindJSON(&query); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	query.System = false
	ctx.JSON(
		http.StatusOK,
		response.Auto(service.NewWorkflowService(db.Db, ctx).TypeList(query)),
	)
}

// TypeOrdinaryList 工作流类型普通的列表，即不列出系统内置的工作流类型
func (r WorkflowApi) TypeOrdinaryList(ctx *gin.Context) {
	var query dto.WorkflowTypeQueryDto
	if err := ctx.ShouldBindJSON(&query); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	query.System = true
	ctx.JSON(
		http.StatusOK,
		response.Auto(service.NewWorkflowService(db.Db, ctx).TypeList(query)),
	)
}

func (r WorkflowApi) TypeDelete(ctx *gin.Context) {
	var post dto.SingleUintRequired
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(
		http.StatusOK,
		response.Auto(nil, service.NewWorkflowService(db.Db, ctx).TypeDelete(post.ID)),
	)
}

func (r WorkflowApi) TypeDetail(ctx *gin.Context) {
	var post dto.SingleUintRequired
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(
		http.StatusOK,
		response.Auto(service.NewWorkflowService(db.Db, ctx).TypeDetail(post.ID)),
	)
}

func (r WorkflowApi) TypeDetailByOnlyName(ctx *gin.Context) {
	var post dto.SingleStringRequired
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(
		http.StatusOK,
		response.Auto(service.NewWorkflowService(db.Db, ctx).TypeDetailByOnlyName(post.ID)),
	)
}

func (r WorkflowApi) TypeOptions(ctx *gin.Context) {
	keyWords := ctx.DefaultQuery("keyWords", "")
	system := ctx.DefaultQuery("system", "")

	systemQuery := false

	if len(system) > 0 {
		systemQuery = true
	}

	ctx.JSON(
		http.StatusOK,
		response.Auto(service.NewWorkflowService(db.Db, ctx).TypeOptions(keyWords, systemQuery)),
	)
}

func (r WorkflowApi) NodeAdd(ctx *gin.Context) {
	var post dto.WorkflowNodeDto

	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(
		http.StatusOK,
		response.Auto(service.NewWorkflowService(db.Db, ctx).NodeAdd(post)),
	)
}

func (r WorkflowApi) NodeUpdate(ctx *gin.Context) {
	var post dto.WorkflowNodeDto

	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(
		http.StatusOK,
		response.Auto(service.NewWorkflowService(db.Db, ctx).NodeUpdate(post)),
	)
}

func (r WorkflowApi) NodeList(ctx *gin.Context) {
	var query dto.WorkflowNodeQueryDto
	if err := ctx.ShouldBindJSON(&query); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(
		http.StatusOK,
		response.Auto(service.NewWorkflowService(db.Db, ctx).NodeList(query)),
	)
}

func (r WorkflowApi) NodeDelete(ctx *gin.Context) {
	var post dto.SingleUintRequired
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(
		http.StatusOK,
		response.Auto(nil, service.NewWorkflowService(db.Db, ctx).NodeDelete(post.ID)),
	)
}

// NodeTypeAll 获取指定工作流模板的所有节点(无分页)
func (r WorkflowApi) NodeTypeAll(ctx *gin.Context) {
	var post dto.SingleUintRequired
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(
		http.StatusOK,
		response.Auto(service.NewWorkflowService(db.Db, ctx).NodeTypeAll(post.ID, true, false)),
	)
}

// NodeFirst 获取指定工作流模板的第一个节点
func (r WorkflowApi) NodeFirst(ctx *gin.Context) {
	var post dto.SingleUintRequired
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(
		http.StatusOK,
		response.Auto(service.NewWorkflowService(db.Db, ctx).NodeFirst(post.ID)),
	)
}

// NodeSaveSchema 保存节点表单配置
func (r WorkflowApi) NodeSaveSchema(ctx *gin.Context) {
	var post dto.WorkflowNodeSaveFormDto
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(
		http.StatusOK,
		response.Auto(nil, service.NewWorkflowService(db.Db, ctx).NodeSaveSchema(post)),
	)
}

// NodeGetSchema 获取节点表单配置
func (r WorkflowApi) NodeGetSchema(ctx *gin.Context) {
	var post dto.SingleUintRequired
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(
		http.StatusOK,
		response.Auto(service.NewWorkflowService(db.Db, ctx).NodeGetSchema(post.ID)),
	)
}

// NodeSaveCirculation 保存节点表单配置
func (r WorkflowApi) NodeSaveCirculation(ctx *gin.Context) {
	var post dto.WorkflowNodeSaveCirculationDto
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(
		http.StatusOK,
		response.Auto(nil, service.NewWorkflowService(db.Db, ctx).NodeSaveCirculation(post)),
	)
}

func (r WorkflowApi) Actions(ctx *gin.Context) {
	ctx.JSON(
		http.StatusOK,
		response.SuccessData(service.NewWorkflowService(db.Db, ctx).Actions()),
	)
}

func (r WorkflowApi) StatusList(ctx *gin.Context) {
	ctx.JSON(
		http.StatusOK,
		response.SuccessData(service.NewWorkflowService(db.Db, ctx).StatusList()),
	)
}

func (r WorkflowApi) LogPageLists(ctx *gin.Context) {
	var query dto.WorkflowLogQueryDto
	if err := ctx.ShouldBindJSON(&query); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(
		http.StatusOK,
		response.Auto(service.NewWorkflowService(db.Db, ctx).LogPageLists(query)),
	)
}

// Footprint 工作流足迹
func (r WorkflowApi) Footprint(ctx *gin.Context) {
	var post dto.SingleUintRequired
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(
		http.StatusOK,
		response.Auto(service.NewWorkflowService(db.Db, ctx).Footprint(post.ID)),
	)
}

// NewWorkflow 创建新的工作流
// 根据工作流唯一标识获取工作流模板信息并返回给前端
// 再由前端请求另外的接口创建工作流实例
func (r WorkflowApi) NewWorkflow(ctx *gin.Context) {
	var post dto.SingleStringRequired
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(
		http.StatusOK,
		response.Auto(service.NewWorkflowService(db.Db, ctx).NewWorkflow(post.ID)),
	)
}
