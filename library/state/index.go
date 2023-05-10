package state

// todo 建议使用 github.com/bits-and-blooms/bitset 库代替

// Modifier 状态修改器
type Modifier struct {
	State int // 当前状态
}

func NewModifier(state int) *Modifier {
	return &Modifier{State: state}
}

// Attach 附加
func (receiver *Modifier) Attach(in int) int {
	receiver.State = receiver.State | in
	return receiver.State
}

// Detach 分离
func (receiver *Modifier) Detach(in int) int {
	receiver.State = receiver.State ^ in // 按位异或
	return receiver.State
}

// Exist 存在返回True
func (receiver *Modifier) Exist(in int) bool {
	return receiver.State&in == in
}

// NotExist 与 Exist 相反
func (receiver *Modifier) NotExist(in int) bool {
	return receiver.State&in == 0
}

// InCombination State是否在in内
func (receiver *Modifier) InCombination(in []int) bool {
	return len(receiver.Contained(in)) > 0
}

// Contained 获取 in 中包含 State 的元素
func (receiver *Modifier) Contained(in []int) []int {
	if len(in) <= 0 {
		return nil
	}
	out := make([]int, 0)
	// 获取输入值的累加值
	max := receiver.LoopAttach(in)
	for i := 1; i <= max; i++ {
		// 是否包含
		if receiver.State&i == receiver.State {
			out = append(out, i)
		}
	}
	return out
}

// LoopAttach 自我附加
func (receiver *Modifier) LoopAttach(in []int) int {
	if len(in) <= 0 {
		return 0
	}
	var num = 0
	for _, i := range in {
		num |= i
	}
	return num
}
