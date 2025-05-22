package stmt

import (
	"fmt"
	"project/expr"
)

func GenerateReturn(tokens []string) []string {
	// 原始 tokens: ["return", "a", "+", "b", "*", "c", ";"]
	exprTokens := tokens[1 : len(tokens)-1] // 截取表达式部分
	if len(exprTokens) == 0 {
		panic("Empty return expression")
	}

	// 构造假的赋值：__return__ = expr
	fakeAssign := append([]string{"__return__", "="}, exprTokens...)
	fmt.Println("[DEBUG] fakeAssign:", fakeAssign) // ✅ debug 输出确认
	code := expr.GenerateAssignExpr(fakeAssign)

	// 取出最后一行获取 return 的临时变量名
	last := code[len(code)-1]
	retVar := last[len("__return__ = "):]
	code = code[:len(code)-1]
	code = append(code, fmt.Sprintf("RETURN %s", retVar))
	return code
}
