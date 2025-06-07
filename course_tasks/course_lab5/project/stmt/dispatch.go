package stmt

import (
	"project/expr"
	"project/generator"
)

// 核心函数：根据 token 判别语句类型并生成三地址码
func Dispatch(tokens []string) []string {
	if len(tokens) > 0 && tokens[0] == "if" && contains(tokens, "else") {
		return GenerateIfElse(tokens)
	} else if len(tokens) > 0 && tokens[0] == "while" {
		return GenerateWhile(tokens)
	} else if len(tokens) > 0 && tokens[0] == "return" {
		return GenerateReturn(tokens)
	} else if len(tokens) > 1 && tokens[1] == "(" {
		return GenerateFunctionCall(tokens)
	} else if len(tokens) >= 3 && tokens[1] == "=" {
		if tokens[2] == "{" { // 支持数组赋值右值为语句块的情况（可选扩展）
			return ParseStmtList(tokens[2:])
		}
		return expr.GenerateAssignExpr(tokens)
	} else if len(tokens) >= 4 && tokens[1] == "[" {
		return GenerateArrayAssignment(tokens)
	} else if len(tokens) > 0 && tokens[0] == "{" {
		return ParseStmtList(tokens)
	} else if len(tokens) > 0 && tokens[0] == "if" {
		return generator.GenerateIfStatement() // 老的测试分支
	} else {
		return generator.GenerateExampleArrayAssignment() // fallback 示例
	}
}

// 工具：判断 token 序列中是否包含某个字符串
func contains(tokens []string, s string) bool {
	for _, tok := range tokens {
		if tok == s {
			return true
		}
	}
	return false
}
