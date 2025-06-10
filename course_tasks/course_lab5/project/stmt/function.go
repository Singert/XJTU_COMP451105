// package stmt

// import (
// 	"fmt"
// )

// func GenerateFunctionDef(tokens []string) []string {
// 	// 支持格式：int foo(int x) { ... }
// 	// 提取函数名和参数
// 	if len(tokens) < 6 || tokens[1] != "(" {
// 		panic("Unsupported function definition syntax")
// 	}
// 	name := tokens[0]
// 	params := []string{}
// 	i := 2
// 	for ; tokens[i] != ")"; i++ {
// 		if tokens[i] != "," {
// 			params = append(params, tokens[i])
// 		}
// 	}
// 	bodyStart := i + 1
// 	if tokens[bodyStart] != "{" {
// 		panic("Function body must be a block")
// 	}
// 	bodyEnd := findCloseBrace(tokens, bodyStart)
// 	bodyTokens := tokens[bodyStart : bodyEnd+1]

// 	// 函数入口标签
// 	code := []string{}
// 	funcLabel := fmt.Sprintf("FUNC_%s", name)
// 	code = append(code, fmt.Sprintf("LABEL %s", funcLabel))

// 	// 参数声明（可忽略类型）
// 	for idx := len(params) - 1; idx >= 0; idx-- {
// 		code = append(code, fmt.Sprintf("POP %s", params[idx]))
// 	}

//		// 生成函数体代码
//		bodyCode := Dispatch(bodyTokens)
//		code = append(code, bodyCode...)
//		code = append(code, fmt.Sprintf("ENDFUNC %s", name))
//		return code
//	}
package stmt

import (
	"fmt"
)

// GenerateFunctionDef 支持函数定义：int foo(int x) { ... } 或 void foo(...) { ... }
// 参数支持形如：int x, int (*f)(), 不检查类型合法性，只处理变量名。

func GenerateFunctionDef(tokens []string) []string {
	if len(tokens) < 6 || tokens[1] != "(" {
		panic("Unsupported function definition syntax")
	}

	funcName := tokens[0]

	// 提取参数
	params := []string{}
	i := 2
	for ; tokens[i] != ")"; i++ {
		if tokens[i] == "," || tokens[i] == "int" || tokens[i] == "void" {
			continue
		}
		if tokens[i+1] == "(" {
			// 跳过函数指针参数 e.g. int soo()
			params = append(params, tokens[i])
			for tokens[i] != ")" {
				i++
			}
		} else {
			params = append(params, tokens[i])
		}
	}
	bodyStart := i + 1

	funcLabel := fmt.Sprintf("FUNC_%s", funcName)
	code := []string{fmt.Sprintf("LABEL %s", funcLabel)}

	// 逆序 POP 参数
	for j := len(params) - 1; j >= 0; j-- {
		code = append(code, fmt.Sprintf("POP %s", params[j]))
	}

	// 处理函数体：递归调用 Dispatch 处理语句块
	var bodyTokens []string
	if bodyStart < len(tokens) && tokens[bodyStart] == "{" {
		bodyEnd := FindCloseBrace(tokens, bodyStart)
		bodyTokens = tokens[bodyStart : bodyEnd+1]
	} else {
		// 非块语句体
		bodyEnd := bodyStart
		for bodyEnd < len(tokens) && tokens[bodyEnd] != ";" {
			bodyEnd++
		}
		bodyTokens = tokens[bodyStart : bodyEnd+1]
	}

	// 递归调用 Dispatch 处理函数体中的语句（支持嵌套调用）
	bodyCode := Dispatch(bodyTokens)
	code = append(code, bodyCode...)
	code = append(code, fmt.Sprintf("ENDFUNC %s", funcName))
	return code
}
