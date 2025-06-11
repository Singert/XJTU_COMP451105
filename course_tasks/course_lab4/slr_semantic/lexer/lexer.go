package lexer

import (
	"strings"
	"unicode"
)

// Token 表示一个词法单元，包含词法类别和原始文本
type Token struct {
	Kind  string // "id", "num", "int" ...
	Value string // 原始值，如 "a", "10"
}

func Tokenize(input string) []string {
	var tokens []string
	current := ""
	runes := []rune(input)
	i := 0

	for i < len(runes) {
		r := runes[i]

		// 跳过空白
		if unicode.IsSpace(r) {
			if current != "" {
				tokens = append(tokens, current)
				current = ""
			}
			i++
			continue
		}

		// 处理数字
		if unicode.IsDigit(r) {
			current += string(r)
			i++
			for i < len(runes) && unicode.IsDigit(runes[i]) {
				current += string(runes[i])
				i++
			}
			tokens = append(tokens, current)
			current = ""
			continue
		}

		// 处理标识符
		if unicode.IsLetter(r) {
			current += string(r)
			i++
			for i < len(runes) && (unicode.IsLetter(runes[i]) || unicode.IsDigit(runes[i])) {
				current += string(runes[i])
				i++
			}
			tokens = append(tokens, current)
			current = ""
			continue
		}

		// 处理符号
		if strings.ContainsRune("[]();,", r) {
			if current != "" {
				tokens = append(tokens, current)
				current = ""
			}
			tokens = append(tokens, string(r))
			i++
			continue
		}

		// 默认跳过其他字符
		i++
	}

	if current != "" {
		tokens = append(tokens, current)
	}
	return tokens
}
