package vo

// OptionItem 前端Select组件选项类型
type OptionItem[T comparable] struct {
	Value string `json:"value"`
	Label T      `json:"label"`
}
