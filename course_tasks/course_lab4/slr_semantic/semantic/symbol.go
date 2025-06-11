package semantic

type Symbol struct {
	Name   string   // 标识符名称
	Type   string   // int 或 void
	Kind   string   // var, array, function
	Dim    []int    // 如果是数组，存储维度
	Params []Symbol // 如果是函数，记录形参列表
}
