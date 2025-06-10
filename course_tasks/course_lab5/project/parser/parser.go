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
		end := findStmtEnd(tokens, i)
		stmtTokens := tokens[i:end]
		stmtCode := stmt.Dispatch(stmtTokens)
		code = append(code, stmtCode...)
		i = end
	}
	return code
}

// findStatementEnd 定位从 start 开始的完整语句结束位置（包含 ; 或结构块）
func findStmtEnd(tokens []string, start int) int {
	tok := tokens[start]

	// 函数定义 int foo(...) 或 void bar(...)
	if (tok == "int" || tok == "void") && start+2 < len(tokens) && tokens[start+2] == "(" {
		// 在这里进行函数定义结束的查找
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
