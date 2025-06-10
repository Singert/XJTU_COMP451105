package stmt

import (
	"fmt"
	"project/expr"
	"project/generator"
	"strings"
)

// Dispatch 根据 token 判别语句类型并生成三地址码
func Dispatch(tokens []string) []string {

	fmt.Println("🔍 Dispatch tokens:", strings.Join(tokens, " "))
	if len(tokens) == 0 {
		return nil
	}

	// 函数定义（int foo(...) 或 void foo(...)）
	if len(tokens) > 3 && (tokens[0] == "int" || tokens[0] == "void") && tokens[2] == "(" {
		return GenerateFunctionDef(tokens[1:]) // 跳过类型
	}

	// if-else 分支
	if tokens[0] == "if" && contains(tokens, "else") {
		return GenerateIfElse(tokens)
	}

	// while 循环
	if tokens[0] == "while" {
		return GenerateWhile(tokens)
	}

	// return 表达式
	if tokens[0] == "return" {
		return GenerateReturn(tokens)
	}

	// print 语句
	if tokens[0] == "print" && len(tokens) >= 2 && tokens[len(tokens)-1] == ";" {
		return []string{fmt.Sprintf("PRINT %s", tokens[1])}
	}

	// 语句块 { ... }
	if tokens[0] == "{" {
		return ParseStmtList(tokens)
	}

	// 数组赋值 a[...,...] = ...
	if len(tokens) >= 4 && tokens[1] == "[" {
		return GenerateArrayAssignment(tokens)
	}

	// 函数调用 foo(...);
	if len(tokens) > 1 && tokens[1] == "(" && tokens[len(tokens)-1] == ";" {
		return GenerateFunctionCall(tokens)
	}

	// 赋值语句 a = ...
	if len(tokens) >= 3 && tokens[1] == "=" {
		return expr.GenerateAssignExpr(tokens)
	}

	// fallback 示例
	return generator.GenerateExampleArrayAssignment()
}

// 判断 tokens 中是否包含指定字符串
func contains(tokens []string, s string) bool {
	for _, tok := range tokens {
		if tok == s {
			return true
		}
	}
	return false
}
