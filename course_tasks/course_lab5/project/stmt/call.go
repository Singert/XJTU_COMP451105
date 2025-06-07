package stmt

import (
	"fmt"
	"project/expr"
)

// 解析函数调用 foo(x, y + 1, a * b)
func GenerateFunctionCall(tokens []string) []string {
	if len(tokens) < 4 || tokens[1] != "(" || tokens[len(tokens)-1] != ";" {
		panic("Invalid function call syntax")
	}

	funcName := tokens[0]
	argTokens := tokens[2 : len(tokens)-2] // 去除 foo ( ... );

	args := [][]string{}
	current := []string{}
	for _, tok := range argTokens {
		if tok == "," {
			args = append(args, current)
			current = []string{}
		} else {
			current = append(current, tok)
		}
	}
	if len(current) > 0 {
		args = append(args, current)
	}

	code := []string{}
	argVars := []string{}
	for _, argExpr := range args {
		// 转为形如 fakeAssign: __arg = expr;
		assign := append([]string{"__arg", "="}, argExpr...)
		tac := expr.GenerateAssignExpr(assign)
		code = append(code, tac[:len(tac)-1]...)      // 保留中间计算
		result := tac[len(tac)-1][len("__arg = "):]   // 提取临时变量
		argVars = append(argVars, result)
	}

	// 逆序传参
	for i := len(argVars) - 1; i >= 0; i-- {
		code = append(code, fmt.Sprintf("PAR %s", argVars[i]))
	}

	t := expr.NewTemp()
	code = append(code, fmt.Sprintf("%s = CALL %s, %d", t, funcName, len(argVars)))
	return code
}
