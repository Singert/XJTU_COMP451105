package semantic

import (
	"fmt"
	"os"
)

// ASTNode 是抽象语法树的节点定义
type ASTNode struct {
	Type  string     // 节点类型，例如 "Decl", "+", "=", "id", "num"
	Value string     // 字面值，例如 "x", "3"
	Left  *ASTNode   // 左子节点
	Right *ASTNode   // 右子节点
	Args  []*ASTNode // 可选的参数列表，用于函数调用或数组索引等
}

func PrintASTPretty(node *ASTNode, prefix string, isTail bool, file string) {
	if node == nil {
		return
	}

	connector := "├─ "
	if isTail {
		connector = "└─ "
	}

	// 打印当前节点类型和值
	output := fmt.Sprintf("%s%s%s", prefix, connector, node.Type)
	if node.Value != "" {
		output += fmt.Sprintf(" (%s)", node.Value)
	}
	output += "\n"

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

	// 计算子节点总数
	children := []*ASTNode{}
	if node.Left != nil {
		children = append(children, node.Left)
	}
	if node.Right != nil {
		children = append(children, node.Right)
	}
	children = append(children, node.Args...)

	// 递归打印子节点
	for i, child := range children {
		last := i == len(children)-1
		newPrefix := prefix
		if isTail {
			newPrefix += "    "
		} else {
			newPrefix += "│   "
		}

		// 对 Args 子节点加索引标注
		if i >= (len(children)-len(node.Args)) && len(node.Args) > 0 {
			argIndex := i - (len(children) - len(node.Args))
			output := fmt.Sprintf("%s%sArg[%d]\n", newPrefix, map[bool]string{true: "└─ ", false: "├─ "}[last], argIndex)

			// 写入文件
			if file != "" {
				f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					fmt.Printf("无法打开文件 %s: %v\n", file, err)
					return
				}
				defer f.Close()

				_, err = f.WriteString(output)
				if err != nil {
					fmt.Printf("无法写入文件 %s: %v\n", file, err)
					return
				}
			}

			// 打印到终端
			fmt.Print(output)

			PrintASTPretty(child, newPrefix+"    ", true, file)
		} else {
			PrintASTPretty(child, newPrefix, last, file)
		}
	}
}
