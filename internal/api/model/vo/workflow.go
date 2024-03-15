package vo

import "VitaTaskGo/internal/repo"

type WorkflowDetailVo struct {
	Workflow     *repo.Workflow          `json:"workflow"`
	Node         *WorkflowNodeVo         `json:"node"`
	Operators    []repo.WorkflowOperator `json:"operators"`
	WorkflowType *repo.WorkflowType      `json:"workflow_type"`
}

type WorkflowNodeVo struct {
	Node        int    `json:"node,omitempty"`
	Name        string `json:"name"`
	Action      string `json:"action"`
	ActionValue string `json:"action_value"`
	Everyone    int    `json:"everyone"`
}
