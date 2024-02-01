package workflow

import (
	"VitaTaskGo/internal/api/data"
	"VitaTaskGo/internal/repo"
)

type AdministratorNodeAction struct {
}

type InitiatorNodeAction struct {
}

func (r *AdministratorNodeAction) ActionName() string {
	return "管理员操作"
}

func (r *AdministratorNodeAction) Handle(engine *Engine) ([]repo.User, error) {
	return data.NewUserRepo(engine.GetCorrectOrm(), engine.ctx).GetAdministrators()
}

func (r *InitiatorNodeAction) ActionName() string {
	return "发起人操作"
}

func (r *InitiatorNodeAction) Handle(engine *Engine) ([]repo.User, error) {
	// 发起人操作
	u, err := data.NewUserRepo(engine.GetCorrectOrm(), engine.ctx).GetUser(engine.workflow.Promoter)
	if err != nil {
		return nil, err
	}

	return []repo.User{*u}, nil
}

// Init 工作流模块初始化
func Init() {
	// 注册节点动作-管理员操作
	RegisterAction("Administrator", &AdministratorNodeAction{})
	// 注册节点动作-发起人操作
	RegisterAction("Initiator", &InitiatorNodeAction{})
}
