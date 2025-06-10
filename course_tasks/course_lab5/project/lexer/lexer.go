// lexer/lexer.go
package lexer

import (
	"strings"
	"unicode"
)

// 词法分析：优化对函数定义的支持
func Tokenize(input string) []string {
	var tokens []string
	current := ""
	runes := []rune(input)

	i := 0
	for i < len(runes) {
		r := runes[i]

		// 跳过空白字符
		if unicode.IsSpace(r) {
			if current != "" {
				tokens = append(tokens, current)
				current = ""
			}
			i++
			continue
		}

		// 处理双字符运算符：==, !=, <=, >=
		if i+1 < len(runes) {
			pair := string(runes[i]) + string(runes[i+1])
			if pair == "==" || pair == "!=" || pair == "<=" || pair == ">=" {
				if current != "" {
					tokens = append(tokens, current)
					current = ""
				}
				tokens = append(tokens, pair)
				i += 2
				continue
			}
		}

		// 处理单字符符号：[], ;, +, -, *等
		if strings.ContainsRune("[],;+*-=/()<>", r) {
			if current != "" {
				tokens = append(tokens, current)
				current = ""
			}
			tokens = append(tokens, string(r))
			i++
			continue
		}

		// 识别函数参数（在括号内的部分）
		if r == '(' || r == ')' {
			if current != "" {
				tokens = append(tokens, current)
				current = ""
			}
			tokens = append(tokens, string(r))
			i++
			continue
		}

		// 否则拼接当前字符（标识符或数字）
		current += string(r)
		i++
	}

	// 处理剩余的部分
	if current != "" {
		tokens = append(tokens, current)
	}

	return tokens
}
