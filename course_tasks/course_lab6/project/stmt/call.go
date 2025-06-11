package stmt

import (
	"fmt"
	"project/expr"
	"strings"
)

// // 解析函数调用 foo(x, y + 1, a * b)
// func GenerateFunctionCall(tokens []string) []string {
// 	if len(tokens) < 4 || tokens[1] != "(" || tokens[len(tokens)-1] != ";" {
// 		panic("Invalid function call syntax")
// 	}

// 	funcName := tokens[0]
// 	argTokens := tokens[2 : len(tokens)-2] // 去除 foo ( ... );

// 	args := [][]string{}
// 	current := []string{}
// 	for _, tok := range argTokens {
// 		if tok == "," {
// 			args = append(args, current)
// 			current = []string{}
// 		} else {
// 			current = append(current, tok)
// 		}
// 	}
// 	if len(current) > 0 {
// 		args = append(args, current)
// 	}

// 	code := []string{}
// 	argVars := []string{}
// 	for _, argExpr := range args {
// 		// 转为形如 fakeAssign: __arg = expr;
// 		assign := append([]string{"__arg", "="}, argExpr...)
// 		tac := expr.GenerateAssignExpr(assign)
// 		code = append(code, tac[:len(tac)-1]...)      // 保留中间计算
// 		result := tac[len(tac)-1][len("__arg = "):]   // 提取临时变量
// 		argVars = append(argVars, result)
// 	}

// 	// 逆序传参
// 	for i := len(argVars) - 1; i >= 0; i-- {
// 		code = append(code, fmt.Sprintf("PAR %s", argVars[i]))
// 	}

//		t := expr.NewTemp()
//		code = append(code, fmt.Sprintf("%s = CALL %s, %d", t, funcName, len(argVars)))
//		return code
//	}
func GenerateFunctionCall(tokens []string) []string {
	if len(tokens) < 4 || tokens[1] != "(" || tokens[len(tokens)-1] != ";" {
		panic("Invalid function call syntax")
	}
	// 兼容现有处理
	return generateCallCore(tokens[0], tokens[2:len(tokens)-2])
}
func GenerateFunctionCallWithAssign(lhs string, fullTokens []string) []string {
	if len(fullTokens) < 4 || fullTokens[1] != "(" || fullTokens[len(fullTokens)-1] != ";" {
		panic("Invalid assigned function call syntax")
	}
	funcName := fullTokens[0]
	argTokens := fullTokens[2 : len(fullTokens)-1]
	returnVars := generateCallCore(funcName, argTokens)

	lastLine := returnVars[len(returnVars)-1]
	equalIdx := strings.Index(lastLine, "=")
	if equalIdx == -1 {
		panic("Expected CALL assignment in last line")
	}
	lastTemp := strings.TrimSpace(lastLine[:equalIdx])

	return append(returnVars, fmt.Sprintf("%s = %s", lhs, lastTemp))
}

func generateCallCore(funcName string, argTokens []string) []string {
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
		assign := append([]string{"__arg", "="}, argExpr...)
		tac := expr.GenerateAssignExpr(assign)
		code = append(code, tac[:len(tac)-1]...)
		result := tac[len(tac)-1][len("__arg = "):]
		argVars = append(argVars, result)
	}

	for i := len(argVars) - 1; i >= 0; i-- {
		code = append(code, fmt.Sprintf("PAR %s", argVars[i]))
	}

	t := expr.NewTemp()
	code = append(code, fmt.Sprintf("%s = CALL %s, %d", t, funcName, len(argVars)))
	return code
}
