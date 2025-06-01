package semantic

import (
	"fmt"
	"os"
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

// PrintQuadruples 打印四元组
func PrintQuadruples(file string) {
	for _, quad := range quadruples {
		output := fmt.Sprintf("(%s, %s, %s, %s)\n", quad.Op, quad.Left, quad.Right, quad.Result)

		// 如果文件路径不为空，则同时输出到文件
		if file != "" {
			f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Printf("无法打开文件 %s: %v\n", file, err)
				return
			}
			defer f.Close()

			// 写入文件
			_, err = f.WriteString(output)
			if err != nil {
				fmt.Printf("无法写入文件 %s: %v\n", file, err)
				return
			}
		}

		// 打印到终端
		fmt.Print(output)
	}
}
