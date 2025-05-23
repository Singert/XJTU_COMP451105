// lexer/lexer.go
package lexer

import (
	"strings"
	"unicode"
)

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

		// 双字符运算符：==, !=, <=, >=
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

		// 单字符符号
		if strings.ContainsRune("[],;+*-=/()<>", r) {
			if current != "" {
				tokens = append(tokens, current)
				current = ""
			}
			tokens = append(tokens, string(r))
			i++
			continue
		}

		// 否则拼接当前字符
		current += string(r)
		i++
	}

	if current != "" {
		tokens = append(tokens, current)
	}

	return tokens
}
