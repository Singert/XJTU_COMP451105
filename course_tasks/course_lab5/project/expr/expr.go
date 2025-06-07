package expr

import (
	"fmt"
)

var tempCounter = 0
var labelCounter = 0

func NewTemp() string {
	tempCounter++
	return fmt.Sprintf("t%d", tempCounter)
}

func NewLabel() string {
	labelCounter++
	return fmt.Sprintf("L%d", labelCounter)
}

func GenerateAssignExpr(tokens []string) []string {
	code := []string{}
	if len(tokens) < 3 || tokens[1] != "=" {
		panic("Not a valid assignment expression")
	}

	target := tokens[0]
	exprTokens := tokens[2:] // FIXED: 不再丢失最后一个 token
	var stack []string
	for i := 0; i < len(exprTokens); i++ {
		tok := exprTokens[i]

		if tok == "+" || tok == "-" || tok == "*" || tok == "/" {
			if len(stack) < 1 || i+1 >= len(exprTokens) {
				panic("Invalid expression: insufficient operands")
			}
			left := stack[len(stack)-1]
			right := exprTokens[i+1]
			i++

			t := NewTemp()
			code = append(code, fmt.Sprintf("%s = %s %s %s", t, left, tok, right))
			stack[len(stack)-1] = t
		} else {
			stack = append(stack, tok)
		}
	}

	if len(stack) == 0 {
		panic("Empty expression result")
	}

	code = append(code, fmt.Sprintf("%s = %s", target, stack[0]))
	return code
}
