package stmt

import (
	"fmt"
	"project/boolean"
	"project/expr"
)

// 支持 if (B) S; else S; 或 块语句 { ... }
func GenerateIfElse(tokens []string) []string {
	condEnd := findCloseParen(tokens, 1)
	condTokens := tokens[2:condEnd]
	bexpr := boolean.GenerateBoolExpr(condTokens)

	var thenCode, elseCode []string
	thenStart := condEnd + 1

	// THEN 分支
	if tokens[thenStart] == "{" {
		blockEnd := findCloseBrace(tokens, thenStart)
		thenBlock := tokens[thenStart : blockEnd+1]
		thenCode = ParseStmtList(thenBlock)
		elseStart := blockEnd + 2

		if tokens[elseStart] == "{" {
			blockEnd2 := findCloseBrace(tokens, elseStart)
			elseBlock := tokens[elseStart : blockEnd2+1]
			elseCode = ParseStmtList(elseBlock)
		} else {
			// else 后单条语句
			elseEnd := elseStart
			for ; elseEnd < len(tokens); elseEnd++ {
				if tokens[elseEnd] == ";" {
					break
				}
			}
			elseCode = Dispatch(tokens[elseStart : elseEnd+1])
		}
	} else {
		// then 后单条语句
		thenEnd := thenStart
		for ; thenEnd < len(tokens); thenEnd++ {
			if tokens[thenEnd] == ";" {
				break
			}
		}
		thenCode = Dispatch(tokens[thenStart : thenEnd+1])

		elseStart := thenEnd + 2
		elseEnd := elseStart
		for ; elseEnd < len(tokens); elseEnd++ {
			if tokens[elseEnd] == ";" {
				break
			}
		}
		elseCode = Dispatch(tokens[elseStart : elseEnd+1])
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

// 小括号匹配（用于 if(...)）
func findCloseParen(tokens []string, start int) int {
	level := 0
	for i := start; i < len(tokens); i++ {
		if tokens[i] == "(" {
			level++
		} else if tokens[i] == ")" {
			level--
			if level == 0 {
				return i
			}
		}
	}
	panic("Unmatched parenthesis")
}

// 大括号匹配（用于块语句）
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
