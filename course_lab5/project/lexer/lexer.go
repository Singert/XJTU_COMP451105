// lexer/lexer.go
package lexer

import (
	"strings"
	"unicode"
)

func Tokenize(input string) []string {
	var tokens []string
	current := ""
	for _, r := range input {
		if unicode.IsSpace(r) {
			if current != "" {
				tokens = append(tokens, current)
				current = ""
			}
		} else if strings.ContainsRune("[],;+*-=/()", r) {
			if current != "" {
				tokens = append(tokens, current)
				current = ""
			}
			tokens = append(tokens, string(r))
		} else {
			current += string(r)
		}
	}
	if current != "" {
		tokens = append(tokens, current)
	}
	return tokens
}