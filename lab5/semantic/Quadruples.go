package semantic

import (
	"fmt"
)

// Quadruple 表示四元组结构
type Quadruple struct {
	Op     string // 操作符，例如 "=", "+", "*", 等
	Left   string // 第一个操作数
	Right  string // 第二个操作数（如果有的话）
	Result string // 结果，通常是临时变量或目标变量
}

var quadruples []Quadruple // 存储生成的四元组
var labelCounter int = 0   // 跳转标签计数器
var tempVarCounter int = 0

func PrintQuadruples() {
	for _, quad := range quadruples {
		fmt.Printf("(%s, %s, %s, %s)\n", quad.Op, quad.Left, quad.Right, quad.Result)
	}
}
