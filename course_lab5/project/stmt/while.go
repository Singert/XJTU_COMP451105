package stmt

import (
	"fmt"
	"project/boolean"
	"project/expr"
)

// 支持 while (B) S;
func GenerateWhile(tokens []string) []string {
	condEnd := findCloseParen(tokens, 1)
	condTokens := tokens[2:condEnd]
	bexpr := boolean.GenerateBoolExpr(condTokens)

	bodyStart := condEnd + 1
	var bodyCode []string

	if tokens[bodyStart] == "{" {
		blockEnd := findCloseBrace(tokens, bodyStart)
		bodyTokens := tokens[bodyStart : blockEnd+1]
		bodyCode = ParseStmtList(bodyTokens)
	} else {
		// 单条语句
		bodyEnd := bodyStart
		for ; bodyEnd < len(tokens); bodyEnd++ {
			if tokens[bodyEnd] == ";" {
				break
			}
		}
		bodyTokens := tokens[bodyStart : bodyEnd+1]
		bodyCode = expr.GenerateAssignExpr(bodyTokens)
	}

	entry := expr.NewLabel()
	code := []string{}
	code = append(code, fmt.Sprintf("LABEL %s", entry))
	code = append(code, bexpr.Code...)
	code = append(code, fmt.Sprintf("LABEL %s", bexpr.TC))
	code = append(code, bodyCode...)
	code = append(code, fmt.Sprintf("GOTO %s", entry))
	code = append(code, fmt.Sprintf("LABEL %s", bexpr.FC))
	return code
}
