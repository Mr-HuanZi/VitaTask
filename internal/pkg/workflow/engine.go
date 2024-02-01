package workflow

import (
	"VitaTaskGo/internal/api/data"
	"VitaTaskGo/internal/pkg/auth"
	"VitaTaskGo/internal/repo"
	"VitaTaskGo/pkg/db"
	"VitaTaskGo/pkg/exception"
	"VitaTaskGo/pkg/response"
	"errors"
	"github.com/duke-git/lancet/v2/random"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/duke-git/lancet/v2/strutil"
	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon/v2"
	"github.com/valyala/fastjson"
	"gorm.io/gorm"
	"strconv"
)

type Engine struct {
	Orm            *gorm.DB
	TransactionOrm *gorm.DB // 事务Orm，事务结束后记得清除
	ctx            *gin.Context

	typeId     uint
	typeData   *repo.WorkflowType
	workflowId uint
	workflow   *repo.Workflow
	operator   []repo.WorkflowOperator
	nodeInfo   *repo.WorkflowNode
	Repo       EngineRepo
	// 是否初始化
	initialized bool
	// 表单数据
	formData map[string]interface{}
}

type EngineRepo struct {
	workflowTypeRepo     repo.WorkflowTypeRepo
	workflowRepo         repo.WorkflowRepo
	workflowOperatorRepo repo.WorkflowOperatorRepo
	workflowNodeRepo     repo.WorkflowNodeRepo
}

// Open 打开一个工作流
func Open(tx *gorm.DB, ctx *gin.Context, workflowId uint) (*Engine, error) {
	workflowTypeRepo := data.NewWorkflowTypeRepo(tx, ctx)
	workflowRepo := data.NewWorkflowRepo(tx, ctx)
	workflowOperatorRepo := data.NewWorkflowOperatorRepo(tx, ctx)
	workflowNodeRepo := data.NewWorkflowNodeRepo(tx, ctx)
	// 查询工作流信息
	workflow, err := workflowRepo.Get(workflowId)
	if err != nil {
		return nil, err
	}

	// 查询工作流模板数据
	typeData, err := workflowTypeRepo.Get(workflow.TypeId)
	if err != nil {
		return nil, err
	}

	// 获取当前操作人(会有多个的情况，所以这里是Slice)
	operator, err := workflowOperatorRepo.GetWorkflowOperatorByNode(workflowId, workflow.Node)
	if err != nil {
		return nil, err
	}

	// 获取当前节点信息
	node, err := workflowNodeRepo.GetAppointNode(workflowId, workflow.Node)

	// todo 获取工作流数据

	// 设置属性
	engine := &Engine{
		Orm:         tx,
		ctx:         ctx,
		typeId:      typeData.ID,
		typeData:    typeData,
		workflowId:  workflow.ID,
		workflow:    workflow,
		operator:    operator,
		nodeInfo:    node,
		Repo:        EngineRepo{workflowTypeRepo: workflowTypeRepo, workflowRepo: workflowRepo, workflowOperatorRepo: workflowOperatorRepo, workflowNodeRepo: workflowNodeRepo},
		initialized: true,
		formData:    make(map[string]interface{}),
	}
	return engine, nil
}

// Create 创建一个工作流
func Create(tx *gorm.DB, ctx *gin.Context, typeId uint) (*Engine, error) {
	workflowTypeRepo := data.NewWorkflowTypeRepo(tx, ctx)
	workflowRepo := data.NewWorkflowRepo(tx, ctx)
	workflowOperatorRepo := data.NewWorkflowOperatorRepo(tx, ctx)
	workflowNodeRepo := data.NewWorkflowNodeRepo(tx, ctx)
	// 查询工作流模板数据
	typeData, err := workflowTypeRepo.Get(typeId)
	if err != nil {
		return nil, err
	}

	// 设置属性
	engine := &Engine{
		Orm:         tx,
		ctx:         ctx,
		typeId:      typeId,
		typeData:    typeData,
		workflowId:  0,
		workflow:    nil,
		Repo:        EngineRepo{workflowTypeRepo: workflowTypeRepo, workflowRepo: workflowRepo, workflowOperatorRepo: workflowOperatorRepo, workflowNodeRepo: workflowNodeRepo},
		initialized: true,
		formData:    make(map[string]interface{}),
	}

	return engine, nil
}

// Initiate 发起工作流
func (engine *Engine) Initiate() error {
	// 检查是否初始化
	if !engine.initialized {
		return exception.NewException(response.WorkflowEngineNotInitialized)
	}

	// 取当前用户
	user, err := auth.CurrUser(engine.ctx)
	if err != nil {
		return err
	}

	// todo 调用Hook

	// 启动事务
	transactionErr := engine.Orm.Transaction(func(tx *gorm.DB) error {
		var err error

		engine.TransactionOrm = tx
		defer func() {
			engine.TransactionOrm = nil
			// 还原所有Repo的Orm实例
			engine.Repo.workflowRepo.SetDbInstance(engine.Orm)
			engine.Repo.workflowNodeRepo.SetDbInstance(engine.Orm)
			engine.Repo.workflowOperatorRepo.SetDbInstance(engine.Orm)
			engine.Repo.workflowTypeRepo.SetDbInstance(engine.Orm)
		}()

		// 给所有Repo设置新的Orm实例
		engine.Repo.workflowRepo.SetDbInstance(tx)
		engine.Repo.workflowNodeRepo.SetDbInstance(tx)
		engine.Repo.workflowOperatorRepo.SetDbInstance(tx)
		engine.Repo.workflowTypeRepo.SetDbInstance(tx)

		// 生成序列号
		serials, err := engine.GenerateSerials()
		if err != nil {
			return err
		}

		// 保存工作流
		workflow := &repo.Workflow{
			TypeId:    engine.typeId,
			TypeName:  engine.typeData.Name,
			OrgId:     0,       //组织ID，先设为0
			Serials:   serials, // 生成编号
			Promoter:  user.ID,
			Nickname:  user.UserNickname,
			SubmitNum: 1, // 提交次数设置为1
		}

		// 设置工作流标题
		title, ok := engine.formData["title"]
		if ok {
			workflow.Title = title.(string)
		} else {
			// 默认标题
			workflow.Title = "工作流[" + engine.typeData.Name + "]审批"
		}

		// 获取下一个节点配置
		node, NextNodeErr := engine.NextNode()
		if NextNodeErr != nil {
			return err
		}
		if node == nil {
			// 没有下一个节点了，直接设定工作流为结束状态
			workflow.Node = 0
			workflow.Status = StatusCompleted
		} else {
			// 设置当前节点序号
			workflow.Node = node.Node
			workflow.Status = StatusRunning
		}

		// 保存工作流
		err = engine.Repo.workflowRepo.Create(workflow)
		if err != nil {
			return exception.NewException(response.WorkflowEngineSaveMainDataFail)
		}

		// 保存操作人
		if workflow.Status == StatusRunning {
			// 获取下一个节点操作人
			operators, _ := engine.GetOperator(node)
			for _, operator := range operators {
				wo := repo.WorkflowOperator{
					UserId:     operator.ID,
					Nickname:   operator.UserNickname,
					Node:       workflow.Node,
					WorkflowId: workflow.ID,
				}
				err = engine.Repo.workflowOperatorRepo.Create(&wo)
				if err != nil {
					return exception.NewException(response.WorkflowEngineSaveOperatorFail)
				}
			}
		}

		// 尝试写入工作流附加数据

		// 记录日志
		return nil
	})
	return transactionErr
}

// ExamineApprove 审批工作流
func (engine *Engine) ExamineApprove() error {
	// 检查是否初始化
	if !engine.initialized {
		return exception.NewException(response.WorkflowEngineNotInitialized)
	}

	// 取当前用户
	user, err := auth.CurrUser(engine.ctx)
	if err != nil {
		return err
	}

	if engine.IsEnd() {
		return exception.NewException(response.WorkflowEngineEnded)
	}

	// todo 调用Hook

	var nextNode *repo.WorkflowNode
	action, ok := engine.formData["action"]
	if !ok || action == "next" {
		/* 工作流正常流转 */
		// 如果当前工作流是 已驳回 状态
		if engine.workflow.Status == StatusOverrule {
			// 提交次数+1
			engine.workflow.SubmitNum += 1
		}

		// 是否还有其他人没操作
		if !engine.MultipleOperator(user.ID) {
			// 获取下一个节点配置
			node, NextNodeErr := engine.NextNode()
			if NextNodeErr != nil {
				return err
			}
			if node == nil {
				// 没有下一个节点了，直接设定工作流为结束状态
				engine.workflow.Node = 0
				engine.workflow.Status = StatusCompleted
			} else {
				// 设置当前节点序号
				engine.workflow.Node = node.Node
				engine.workflow.Status = StatusRunning
				nextNode = node
			}
		}
	} else if action == "overrule" {
		/* 驳回工作流 */
		// 是否跳转到指定节点
		jumpNode, ok := engine.formData["jump_node"]
		if ok {
			node, err := engine.Repo.workflowNodeRepo.GetAppointNode(engine.typeId, jumpNode.(int))
			if err != nil {
				return db.FirstQueryErrorHandle(err, response.WorkflowNodeNotExist)
			}
			// 跳转的步骤不能大于或等于当前步骤
			if node.Node >= engine.nodeInfo.Node {
				return exception.NewException(response.WorkflowEngineNodeJumpErr)
			}

			engine.workflow.Node = node.Node
			nextNode = node
		} else {
			// 查询第一个节点
			node, err := engine.Repo.workflowNodeRepo.FirstNode(engine.typeId)
			if err != nil {
				return db.FirstQueryErrorHandle(err, response.WorkflowEngineNoFirstNodeSet)
			}

			engine.workflow.Node = node.Node
			nextNode = node
		}
		// 设置工作流状态
		engine.workflow.Status = StatusOverrule
	} else if action == "cancel" {
		/* 作废工作流 */
		// 此操作不更改工作流节点
		// 设置工作流状态
		engine.workflow.Status = StatusVoided
	}

	// 启动事务
	transactionErr := engine.Orm.Transaction(func(tx *gorm.DB) error {
		var err error

		engine.TransactionOrm = tx
		defer func() {
			engine.TransactionOrm = nil
			// 还原所有Repo的Orm实例
			engine.Repo.workflowRepo.SetDbInstance(engine.Orm)
			engine.Repo.workflowNodeRepo.SetDbInstance(engine.Orm)
			engine.Repo.workflowOperatorRepo.SetDbInstance(engine.Orm)
			engine.Repo.workflowTypeRepo.SetDbInstance(engine.Orm)
		}()

		// 给所有Repo设置新的Orm实例
		engine.Repo.workflowRepo.SetDbInstance(tx)
		engine.Repo.workflowNodeRepo.SetDbInstance(tx)
		engine.Repo.workflowOperatorRepo.SetDbInstance(tx)
		engine.Repo.workflowTypeRepo.SetDbInstance(tx)

		// todo 执行钩子

		// 下一个节点的操作人
		var operators []repo.User
		/* 判断工作流状态 Start */
		if engine.workflow.Status == StatusCompleted {
			/* 工作流已完成，删除该工作流的所有操作人 */
			err := engine.Repo.workflowOperatorRepo.RemoveWorkflowAllOperator(engine.workflowId)
			if err != nil {
				return exception.NewException(response.WorkflowEngineRemoveOperatorFail)
			}
		} else if engine.workflow.Status == StatusOverrule {
			/* 被驳回 */
			// 删除该工作流的所有操作人
			err := engine.Repo.workflowOperatorRepo.RemoveWorkflowAllOperator(engine.workflowId)
			if err != nil {
				return exception.NewException(response.WorkflowEngineRemoveOperatorFail)
			}
			// 获取节点操作人
			operators, _ = engine.GetOperator(nextNode)
		} else if engine.workflow.Node == engine.nodeInfo.Node {
			/* 如果工作流步骤没有改变，则尝试将当前操作人改为已确认 */
			err := engine.Repo.workflowOperatorRepo.SetHandled(engine.workflowId, engine.nodeInfo.Node, user.ID)
			if err != nil {
				return exception.NewException(response.WorkflowEngineOperatorHandleFail)
			}
		} else {
			/* 进入下一步 */
			// 获取节点操作人
			operators, _ = engine.GetOperator(nextNode)
		}
		/* 判断工作流状态 End */

		// 更新主表数据
		err = engine.Repo.workflowRepo.Save(engine.workflow)
		if err != nil {
			return exception.NewException(response.WorkflowEngineSaveMainDataFail)
		}

		// 保存操作人
		if len(operators) > 0 {
			for _, operator := range operators {
				wo := repo.WorkflowOperator{
					UserId:     operator.ID,
					Nickname:   operator.UserNickname,
					Node:       engine.workflow.Node,
					WorkflowId: engine.workflowId,
				}
				err = engine.Repo.workflowOperatorRepo.Create(&wo)
				if err != nil {
					return exception.NewException(response.WorkflowEngineSaveOperatorFail)
				}
			}
		}

		// 尝试写入工作流附加数据

		// 记录日志
		return nil
	})
	return transactionErr
}

// NextNode 获取当前工作流的下一个节点
func (engine *Engine) NextNode() (*repo.WorkflowNode, error) {
	// 假设当前节点是0
	var currNode = 0
	// 先获取当前工作流的节点
	if engine.workflow != nil {
		currNode = engine.workflow.Node
	}
	// 如果节点序号小于等于0，表示该工作流还没有正式发起
	if currNode <= 0 {
		// 获取第一个节点
		node, err := engine.Repo.workflowNodeRepo.GetNextNode(engine.typeId, 0)
		if err != nil {
			// 第一个节点未设置
			return nil, db.FirstQueryErrorHandle(err, response.WorkflowEngineNoFirstNodeSet)
		}
		// 把第一个节点的序号赋值给 currNode
		currNode = node.Node
	}

	// 获取下一个节点配置
	node, err := engine.Repo.workflowNodeRepo.GetNextNode(engine.typeId, currNode)
	if err != nil {
		// 如果只是没有记录，两个参数都返回nil
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	// todo 节点自定义条件判断，以后再说吧
	return node, nil
}

func (engine *Engine) GetOperator(workflowNode *repo.WorkflowNode) ([]repo.User, error) {
	var (
		userList = make([]repo.User, 0)
		err      error
	)

	// 实例化UserRepo
	userRepo := data.NewUserRepo(engine.GetCorrectOrm(), engine.ctx)

	// 如果直接指定了用户
	if len(workflowNode.ActionValue) > 0 {
		// 只能是json格式字符串
		var fp fastjson.Parser
		fpv, err := fp.Parse(workflowNode.ActionValue)
		if err != nil {
			return nil, err
		}

		for _, item := range fpv.GetArray() {
			u, err := userRepo.GetUser(uint64(item.GetInt()))
			if err != nil {
				return nil, err
			}
			userList = append(userList, *u)
		}

		return userList, err
	}

	// 处理Action
	if len(workflowNode.Action) > 0 {
		nodeAction, err := GetAction(workflowNode.Action)
		if err != nil {
			return nil, err
		}

		userList, err = nodeAction.Handle(engine)
	}

	// 如果用户列表还是空的，取当前登录人
	if len(userList) <= 0 {
		var u *repo.User
		u, err = auth.CurrUser(engine.ctx)
		if err != nil {
			return nil, err
		}

		userList = append(userList, *u)
	}

	return userList, nil
}

func (engine *Engine) GenerateSerials() (string, error) {
	start := carbon.Now().StartOfDay().TimestampMilli()
	end := carbon.Now().EndOfDay().TimestampMilli()
	// 获取当天的开始与结束
	// 这里暂时不考虑并发的情况
	total, err := engine.Repo.workflowRepo.GetDayTotal(start, end)
	if err != nil {
		return "", exception.ErrorHandle(err, response.WorkflowEngineSerialGenerationFailed)
	}

	// 用 0 填充total为 4 位长度的字符串
	// total 要先 +1
	index := strutil.PadStart(strconv.FormatInt(total+1, 10), 4, "0")
	// 例子: 20230807 + <random number>*3 + 0001
	return carbon.Now().ToShortDateString() + strconv.Itoa(random.RandInt(100, 999)) + index, nil
}

func (engine *Engine) GetCorrectOrm() *gorm.DB {
	if engine.TransactionOrm != nil {
		return engine.TransactionOrm
	} else {
		return engine.Orm
	}
}

func (engine *Engine) SetFormData(in map[string]interface{}) {
	engine.formData = in
}

func (engine *Engine) SetFormDataField(key string, value interface{}) {
	engine.formData[key] = value
}

// IsEnd 工作流是否已结束
func (engine *Engine) IsEnd() bool {
	if engine.workflow == nil {
		return false
	}

	if slice.Contain([]int{StatusVoided, StatusCompleted}, engine.workflow.Status) {
		return true
	}

	return false
}

// MultipleOperator 是否还有其它操作人
func (engine *Engine) MultipleOperator(userId uint64) bool {
	if engine.nodeInfo.Everyone == 1 {
		b, err := engine.Repo.workflowOperatorRepo.OtherOperator(engine.workflowId, engine.nodeInfo.Node, userId)
		if err != nil {
			_ = exception.ErrorHandle(err, response.DbQueryError)
			return false
		}
		return b
	}

	return false
}
