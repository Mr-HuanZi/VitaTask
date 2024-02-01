package service

import (
	"VitaTaskGo/internal/api/data"
	"VitaTaskGo/internal/api/model/dto"
	"VitaTaskGo/internal/pkg"
	"VitaTaskGo/internal/pkg/auth"
	"VitaTaskGo/internal/pkg/constant"
	"VitaTaskGo/internal/repo"
	"VitaTaskGo/pkg/exception"
	"VitaTaskGo/pkg/response"
	"github.com/duke-git/lancet/v2/maputil"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"
)

type TaskLogService struct {
	Orm  *gorm.DB
	ctx  *gin.Context
	repo repo.TaskLogRepo
}

func NewTaskLogService(tx *gorm.DB, ctx *gin.Context) *TaskLogService {
	return &TaskLogService{
		Orm:  tx,  // 赋予ORM实例
		ctx:  ctx, // 传递上下文
		repo: data.NewTaskLogRepo(tx, ctx),
	}
}

// Add 新增日志
func (receiver TaskLogService) Add(dto dto.TaskLogForm) (*repo.TaskLog, error) {
	// 检查任务是否存在
	// 此处会检索软删除的记录
	err := receiver.Orm.Model(&repo.Task{}).Unscoped().First(&repo.Task{}, dto.TaskId).Error
	if err != nil {
		// 有错误就代表不存在
		return nil, err
	}

	// 操作类型是否合法
	if !slice.Contain(maputil.Keys(constant.GetTaskLogOperatorMaps()), dto.OperateType) {
		return nil, exception.NewException(response.TaskOperatorTypeIllegal)
	}

	// 操作人
	if dto.Operator <= 0 {
		// 获取当前用户
		currUser, err := auth.CurrUser(receiver.ctx)
		if err != nil {
			return nil, err
		}

		dto.Operator = currUser.ID
	} else {
		// 检查操作人是否存在
		_, err = data.NewUserRepo(receiver.Orm, receiver.ctx).GetUser(dto.Operator)
		if err != nil {
			return nil, exception.ErrorHandle(err, response.UserNotFound)
		}
	}
	// 操作时间
	if dto.OperateTime <= 0 {
		dto.OperateTime = carbon.Now().TimestampMilli()
	}
	logData := &repo.TaskLog{
		TaskId:      dto.TaskId,
		OperateType: dto.OperateType,
		Operator:    dto.Operator,
		OperateTime: dto.OperateTime,
		Message:     dto.Message,
	}

	return logData, receiver.repo.Create(logData)
}

func (receiver TaskLogService) List(query dto.TaskLogQuery) *dto.PagedResult[repo.TaskLog] {
	logs, total, err := receiver.repo.PageListTaskLog(query)

	if err != nil {
		logs = make([]repo.TaskLog, 0)
		_ = exception.ErrorHandle(err, response.DbQueryError, "任务日志列表查询失败: ")
	}

	return pkg.PagedResult[repo.TaskLog](logs, total, int64(query.Page))
}

// Operators 获取操作类型
func (receiver TaskLogService) Operators() []map[string]string {
	operators := constant.GetTaskLogOperatorMaps()
	s := make([]map[string]string, len(operators))

	i := 0
	for k, item := range operators {
		s[i] = map[string]string{"label": item, "value": k}
		i++
	}
	return s
}
