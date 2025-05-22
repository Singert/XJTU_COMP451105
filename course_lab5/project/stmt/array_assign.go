package stmt

import (
	"fmt"
	"project/expr"
)

func GenerateArrayAssignment(tokens []string) []string {
	// 解析形如：a[i + 1, j * 2, 4] = 66;

	// 提取数组名
	name := tokens[0]
	// 提取下标部分：tokens[2: ?]，直到遇到 "]"
	var idxExprs [][]string
	i := 2
	current := []string{}
	for ; i < len(tokens); i++ {
		tok := tokens[i]
		if tok == "]" {
			if len(current) > 0 {
				idxExprs = append(idxExprs, current)
			}
			break
		} else if tok == "," {
			idxExprs = append(idxExprs, current)
			current = []string{}
		} else {
			current = append(current, tok)
		}
	}

	// 下标维度默认：d1=5, d2=20, d3=4
	dims := []int{5, 20, 4}
	code := []string{}
	vars := []string{}

	for _, e := range idxExprs {
		assign := append([]string{"__tmp", "="}, e...)
		tac := expr.GenerateAssignExpr(assign)
		code = append(code, tac[:len(tac)-1]...)
		res := tac[len(tac)-1][len("__tmp = "):]
		vars = append(vars, res)
	}

	// 地址计算：(v1 * d2 + v2) * d3 + v3 → *4
	t := vars[0]
	for i := 1; i < len(vars); i++ {
		t1 := expr.NewTemp()
		code = append(code, fmt.Sprintf("%s = %s * %d", t1, t, dims[i-1]))
		t2 := expr.NewTemp()
		code = append(code, fmt.Sprintf("%s = %s + %s", t2, t1, vars[i]))
		t = t2
	}

	// ×4（元素大小）
	offset := expr.NewTemp()
	code = append(code, fmt.Sprintf("%s = %s * 4", offset, t))

	// 处理右侧表达式
	rightExpr := tokens[i+2 : len(tokens)-1]
	rightAssign := append([]string{"__val", "="}, rightExpr...)
	rightCode := expr.GenerateAssignExpr(rightAssign)
	code = append(code, rightCode...)

	val := rightCode[len(rightCode)-1][len("__val = "):]

	code = append(code, fmt.Sprintf("%s[%s] = %s", name, offset, val))
	return code
}
