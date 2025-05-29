package semantic

import (
	"fmt"
	"strings"
)

// ASTNode 是抽象语法树的节点定义
type ASTNode struct {
	Type  string     // 节点类型，例如 "Decl", "+", "=", "id", "num"
	Value string     // 字面值，例如 "x", "3"
	Left  *ASTNode   // 左子节点
	Right *ASTNode   // 右子节点
	Args  []*ASTNode // 可选的参数列表，用于函数调用或数组索引等
}

// PrintAST 递归打印抽象语法树，带缩进
//
//	func PrintAST(node *ASTNode, indent int) {
//		if node == nil {
//			return
//		}
//		prefix := strings.Repeat("  ", indent)
//		fmt.Printf("%s%s", prefix, node.Type)
//		if node.Value != "" {
//			fmt.Printf(" (%s)", node.Value)
//		}
//		fmt.Println()
//		PrintAST(node.Left, indent+1)
//		PrintAST(node.Right, indent+1)
//	}


func PrintASTPretty(node *ASTNode, prefix string, isTail bool) {
	if node == nil {
		return
	}

	connector := "├─ "
	if isTail {
		connector = "└─ "
	}

	// 打印当前节点类型和值
	fmt.Printf("%s%s%s", prefix, connector, node.Type)
	if node.Value != "" {
		fmt.Printf(" (%s)", node.Value)
	}
	fmt.Println()

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
		if i >= (len(children) - len(node.Args)) && len(node.Args) > 0 {
			argIndex := i - (len(children) - len(node.Args))
			fmt.Printf("%s%sArg[%d]\n", newPrefix, map[bool]string{true: "└─ ", false: "├─ "}[last], argIndex)
			PrintASTPretty(child, newPrefix+"    ", true)
		} else {
			PrintASTPretty(child, newPrefix, last)
		}
	}
}

