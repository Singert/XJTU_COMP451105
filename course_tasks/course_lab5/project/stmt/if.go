package stmt

import (
	"fmt"
	"project/boolean"
	"project/expr"
)

// 支持 if (B) S; else S;
func GenerateIfElse(tokens []string) []string {
	condEnd := findCloseParen(tokens, 1)
	condTokens := tokens[2:condEnd]
	bexpr := boolean.GenerateBoolExpr(condTokens)

	var thenCode, elseCode []string
	thenStart := condEnd + 1

	// 块语句 then
	if tokens[thenStart] == "{" {
		blockEnd := findCloseBrace(tokens, thenStart)
		thenAssign := tokens[thenStart : blockEnd+1]
		thenCode = ParseStmtList(thenAssign)

		elseStart := blockEnd + 2
		blockEnd2 := findCloseBrace(tokens, elseStart)
		elseAssign := tokens[elseStart : blockEnd2+1]
		elseCode = ParseStmtList(elseAssign)
	} else {
		// 单条语句 then
		thenEnd := thenStart
		for ; thenEnd < len(tokens); thenEnd++ {
			if tokens[thenEnd] == ";" {
				break
			}
		}
		thenAssign := tokens[thenStart : thenEnd+1]
		thenCode = expr.GenerateAssignExpr(thenAssign)

		elseStart := thenEnd + 2
		elseAssign := tokens[elseStart:]
		elseCode = expr.GenerateAssignExpr(elseAssign)
	}

	endLabel := expr.NewLabel()
	code := []string{}
	code = append(code, bexpr.Code...)
	code = append(code, fmt.Sprintf("LABEL %s", bexpr.TC))
	code = append(code, thenCode...)
	code = append(code, fmt.Sprintf("GOTO %s", endLabel))
	code = append(code, fmt.Sprintf("LABEL %s", bexpr.FC))
	code = append(code, elseCode...)
	code = append(code, fmt.Sprintf("LABEL %s", endLabel))
	return code
}

func findCloseBrace(tokens []string, start int) int {
	level := 0
	for i := start; i < len(tokens); i++ {
		if tokens[i] == "{" {
			level++
		} else if tokens[i] == "}" {
			level--
			if level == 0 {
				return i
			}
		}
	}
	panic("Unmatched brace")
}
