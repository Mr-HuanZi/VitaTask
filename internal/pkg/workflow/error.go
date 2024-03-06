package workflow

import (
	"errors"
)

var (
	// ErrWorkflowTypeNotExist 工作流类型(模板)不存在
	ErrWorkflowTypeNotExist = errors.New("workflow type not exist")
	// ErrWorkflowNotExist 工作流不存在
	ErrWorkflowNotExist = errors.New("workflow not exist")
)
