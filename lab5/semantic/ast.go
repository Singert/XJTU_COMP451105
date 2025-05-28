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
func PrintAST(node *ASTNode, indent int) {
	if node == nil {
		return
	}
	prefix := strings.Repeat("  ", indent)
	fmt.Printf("%s%s", prefix, node.Type)
	if node.Value != "" {
		fmt.Printf(" (%s)", node.Value)
	}
	fmt.Println()

	// 递归打印 Left 和 Right
	PrintAST(node.Left, indent+1)
	PrintAST(node.Right, indent+1)

	// 递归打印 Args 中的所有子节点
	for i, child := range node.Args {
		fmt.Printf("%sArg[%d]:\n", prefix+"  ", i)
		PrintAST(child, indent+2)
	}

}
