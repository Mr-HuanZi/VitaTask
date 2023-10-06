package workflow

import (
	"VitaTaskGo/internal/biz"
	"VitaTaskGo/internal/pkg/exception"
	"VitaTaskGo/internal/pkg/response"
)

type NodeAction interface {
	ActionName() string
	Handle(engine *Engine) ([]biz.User, error)
}

var ActionPool = make(map[string]NodeAction)

// RegisterAction 注册节点动作
func RegisterAction(name string, na NodeAction) {
	// 重复的Key直接覆盖
	ActionPool[name] = na
}

// GetAllActionName 获取已注册的动作名称
// key对应的是 RegisterAction 方法的 name 参数
// value对应的是 NodeAction 的 ActionName 方法
func GetAllActionName() map[string]string {
	kv := make(map[string]string)
	for key, action := range ActionPool {
		kv[key] = action.ActionName()
	}
	return kv
}

func GetAction(name string) (NodeAction, error) {
	na, ok := ActionPool[name]
	if !ok {
		return nil, exception.NewException(response.WorkflowEngineActionNotRegistered)
	}

	return na, nil
}
