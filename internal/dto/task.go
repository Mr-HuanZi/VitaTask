package dto

type TaskListQuery struct {
	PagingQuery
	UintId
	QueryParams
	DeletedQuery
	Project      uint     `json:"project"`
	PlanTime     []string `json:"plan_time"`
	Level        uint     `json:"level"`
	Leader       uint64   `json:"leader"`       // 负责人
	Collaborator []uint64 `json:"collaborator"` // 协助人
	GroupId      uint     `json:"group"`
}

type TaskListQueryBO struct {
	PagingQuery
	UintId
	QueryParams
	DeletedQuery
	ProjectIds          []uint
	PlanTime            []string
	Level               uint
	LeaderTaskIds       []uint
	CollaboratorTaskIds []uint
	GroupId             uint
}

type TaskCreateForm struct {
	ProjectId    uint     `json:"project,omitempty" binding:"required"`
	GroupId      uint     `json:"group,omitempty"`
	Title        string   `json:"title,omitempty" binding:"required"`
	Describe     string   `json:"describe,omitempty"`
	Level        uint     `json:"level,omitempty"`
	PlanTime     []string `json:"plan_time"`
	Leader       uint64   `json:"leader" binding:"required"` // 负责人
	Collaborator []uint64 `json:"collaborator"`              // 协助人
}

type TaskStatusVo struct {
	Label  string `json:"label"`
	Value  int    `json:"value"`
	Status string `json:"status"`
}

type TaskChangeStatus struct {
	SingleUintRequired
	Status int `json:"status"`
}

type TaskGroupForm struct {
	UintId
	ProjectId uint   `json:"project" binding:"required"`
	Name      string `json:"name" binding:"required"`
}

type TaskGroupQuery struct {
	PagingQuery
	QueryParams
	ProjectId uint `json:"project" binding:"required"`
}

type TaskStatistics struct {
	Completed         int64 `json:"completed,omitempty"`
	Processing        int64 `json:"processing,omitempty"`
	FinishOnTime      int64 `json:"finish_on_time,omitempty"`
	TimeoutCompletion int64 `json:"timeout_completion,omitempty"`
}

type DailySituationQuery struct {
	ProjectId uint   `json:"project" binding:"required"`
	StartDate string `json:"start_date,omitempty"`
	EndDate   string `json:"end_date,omitempty"`
}

type DailySituationVo struct {
	Label string `json:"label,omitempty"`
	Date  string `json:"date,omitempty"`
	Value int64  `json:"value"`
}

type TaskLogForm struct {
	TaskId      uint   `json:"task_id"`
	OperateType string `json:"operate_type"`
	Operator    uint64 `json:"operator"`
	OperateTime int64  `json:"operate_time"`
	Message     string `json:"message"`
}

type TaskLogQuery struct {
	PagingQuery
	ProjectId   uint     `json:"project_id"`
	TaskIds     []uint   `json:"task_id"`
	OperateType string   `json:"operate_type"`
	Operator    uint64   `json:"operator"`
	OperateTime []string `json:"operate_time"`
	CreateTime  []string `json:"create_time"`
}
