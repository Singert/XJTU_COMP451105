package boolean

import (
	"fmt"
	"project/expr"
)

type BoolExprResult struct {
	Code []string // 中间代码段
	TC   string   // true label
	FC   string   // false label
}

// 支持简单布尔表达式 B → E relop E
func GenerateCondExpr(tokens []string) ([]string, string) {
	if len(tokens) != 3 {
		panic("Only E relop E format is currently supported")
	}

	left := tokens[0]
	op := tokens[1]
	right := tokens[2]

	t := expr.NewTemp()
	code := []string{fmt.Sprintf("%s = %s %s %s", t, left, op, right)}
	return code, t
}

func GenerateBoolExpr(tokens []string) BoolExprResult {
	// (B)
	if len(tokens) >= 3 && tokens[0] == "(" && tokens[len(tokens)-1] == ")" {
		return GenerateBoolExpr(tokens[1 : len(tokens)-1])
	}

	// !B
	if len(tokens) >= 2 && tokens[0] == "!" {
		sub := GenerateBoolExpr(tokens[1:])
		return BoolExprResult{
			Code: sub.Code,
			TC:   sub.FC,
			FC:   sub.TC,
		}
	}

	// B1 || B2
	if i := findSplit(tokens, "||"); i != -1 {
		left := GenerateBoolExpr(tokens[:i])
		right := GenerateBoolExpr(tokens[i+1:])
		code := append(left.Code, fmt.Sprintf("LABEL %s", left.FC))
		code = append(code, right.Code...)
		return BoolExprResult{
			Code: code,
			TC:   left.TC,
			FC:   right.FC,
		}
	}

	// B1 && B2
	if i := findSplit(tokens, "&&"); i != -1 {
		left := GenerateBoolExpr(tokens[:i])
		right := GenerateBoolExpr(tokens[i+1:])
		code := append(left.Code, fmt.Sprintf("LABEL %s", left.TC))
		code = append(code, right.Code...)
		return BoolExprResult{
			Code: code,
			TC:   right.TC,
			FC:   left.FC,
		}
	}

	// E relop E
	if len(tokens) == 3 && isRelop(tokens[1]) {
		t := expr.NewTemp()
		code := []string{fmt.Sprintf("%s = %s %s %s", t, tokens[0], tokens[1], tokens[2])}
		tc := expr.NewLabel()
		fc := expr.NewLabel()
		code = append(code, fmt.Sprintf("IF %s != 0 THEN %s ELSE %s", t, tc, fc))
		return BoolExprResult{
			Code: code,
			TC:   tc,
			FC:   fc,
		}
	}
	// 单变量条件：如 z → z != 0
	if len(tokens) == 1 {
		t := tokens[0]
		tt := expr.NewTemp()
		code := []string{fmt.Sprintf("%s = %s != 0", tt, t)}
		tc := expr.NewLabel()
		fc := expr.NewLabel()
		code = append(code, fmt.Sprintf("IF %s != 0 THEN %s ELSE %s", tt, tc, fc))
		return BoolExprResult{
			Code: code,
			TC:   tc,
			FC:   fc,
		}
	}
	panic(fmt.Sprintf("Unsupported boolean expression: %v", tokens))
}

func isRelop(op string) bool {
	switch op {
	case "<", "<=", ">", ">=", "==", "!=":
		return true
	default:
		return false
	}
}

// 辅助函数：处理括号平衡并定位逻辑运算符
func findSplit(tokens []string, op string) int {
	level := 0
	for i, tok := range tokens {
		switch tok {
		case "(":
			level++
		case ")":
			level--
		case op:
			if level == 0 {
				return i
			}
		}
	}
	return -1
}
