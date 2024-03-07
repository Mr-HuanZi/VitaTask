package workflow

const (
	// StatusVoided 已作废
	StatusVoided = iota

	// StatusCompleted 已完成
	StatusCompleted

	// StatusRunning 进行中
	StatusRunning

	// StatusOverrule 驳回
	StatusOverrule
)

// StatusMap 状态Map
// 请严格按照常量定义的顺序来
var StatusMap = map[string]int{
	"voided":    StatusVoided,
	"completed": StatusCompleted,
	"running":   StatusRunning,
	"overrule":  StatusOverrule,
}

// StatusEnum 状态枚举，兼容Antd Pro
// 请严格按照常量定义的顺序来
var StatusEnum = map[int]map[string]string{
	StatusVoided:    {"text": "已作废", "status": "Error"},
	StatusCompleted: {"text": "已完成", "status": "Success"},
	StatusRunning:   {"text": "进行中", "status": "Processing"},
	StatusOverrule:  {"text": "驳回", "status": "Warning"},
}
