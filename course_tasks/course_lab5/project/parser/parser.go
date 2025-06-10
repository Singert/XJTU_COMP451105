// parser/parser.go
package parser

import (
	"project/stmt"
)

// ParseProgram：顺序解析完整程序，支持多函数 + 主程序语句
func ParseProgram(tokens []string) []string {
	var code []string
	i := 0
	for i < len(tokens) {
		// 对每个完整结构（函数定义或语句）交由 Dispatch 处理
		end := findStatementEnd(tokens, i)
		stmtTokens := tokens[i:end]
		stmtCode := stmt.Dispatch(stmtTokens)
		code = append(code, stmtCode...)
		i = end
	}
	return code
}

// findStatementEnd：定位语句或函数定义结束位置（包含 ; 或大括号块）
func findStatementEnd(tokens []string, start int) int {
	// 判断是否为函数定义：形如 int foo(...) {
	if start+2 < len(tokens) && tokens[start+2] == "(" {
		parenEnd := stmt.FindCloseParen(tokens, start+2)
		if parenEnd+1 < len(tokens) && tokens[parenEnd+1] == "{" {
			braceEnd := stmt.FindCloseBrace(tokens, parenEnd+1)
			return braceEnd + 1
		}
	}

	// 控制语句或表达式语句：按 ; 结束
	for i := start; i < len(tokens); i++ {
		if tokens[i] == ";" {
			return i + 1
		}
	}
	panic("无法确定语句结束位置")
}
