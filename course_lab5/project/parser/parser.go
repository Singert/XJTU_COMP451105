// parser/parser.go
package parser

import (
	"project/expr"
	"project/generator"
	"project/stmt" // 新增
)

func ParseAndGenerateTAC(tokens []string) []string {
	if len(tokens) > 0 && tokens[0] == "if" {
		return generator.GenerateIfStatement()
	} else if len(tokens) > 0 && tokens[0] == "while" {
		return generator.GenerateWhileStatement()
	} else if len(tokens) > 0 && tokens[0] == "return" {
		return stmt.GenerateReturn(tokens)
	} else if len(tokens) > 1 && tokens[1] == "(" {
		return stmt.GenerateFunctionCall(tokens) 
	} else if len(tokens) >= 3 && tokens[1] == "=" {
		return expr.GenerateAssignExpr(tokens)
	} else if len(tokens) >= 4 && tokens[1] == "[" {
		return stmt.GenerateArrayAssignment(tokens)
	} else {
		return generator.GenerateExampleArrayAssignment()
	}
}
