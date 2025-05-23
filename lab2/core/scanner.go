package core

import (
	"lab2/core/dfa"
	"unicode"
)

type Scanner struct {
	DFAList []struct {
		DFA       *dfa.DFA
		TokenType dfa.TokenType
	}
}

func NewScanner() *Scanner {
	return &Scanner{}
}

func (s *Scanner) RegisterDFA(d *dfa.DFA, t dfa.TokenType) {
	s.DFAList = append(s.DFAList, struct {
		DFA       *dfa.DFA
		TokenType dfa.TokenType
	}{DFA: d, TokenType: t})
}

func isWhitespace(r rune) bool {
	return unicode.IsSpace(r)
}

func (s *Scanner) Scan(input string) (matched dfa.Token, length int) {
	if len(input) == 0 {
		return dfa.Token{Type: dfa.TokenERROR, Lexeme: ""}, 0
	}

	runes := []rune(input)
	if unicode.IsSpace(runes[0]) {
		i := 1
		for i < len(runes) && unicode.IsSpace(runes[i]) {
			i++
		}
		return dfa.Token{Type: dfa.TokenWithespace, Lexeme: string(runes[:i])}, i
	}

	maxLen := 0
	var maxToken dfa.Token
	var maxDFA *dfa.DFA

	for _, entry := range s.DFAList {
		for i := 1; i <= len(runes); i++ {
			sub := string(runes[:i])
			ok, _ := entry.DFA.MatchDFA(sub, false)
			if ok && i > maxLen {
				maxLen = i
				maxToken = dfa.Token{Type: entry.TokenType, Lexeme: sub}
				maxDFA = entry.DFA
			}
		}
	}

	if maxLen == 0 {
		return dfa.Token{Type: dfa.TokenERROR, Lexeme: string(runes[0])}, 1
	}
	_, _ = maxDFA.MatchDFA(maxToken.Lexeme, true)
	return maxToken, maxLen
}

// func (s *Scanner) Scan(input string) (matched dfa.Token, length int) {
// 	fmt.Printf("[Scan] input segment: '%s'\n", input)

// 	maxLen := 0
// 	var maxToken dfa.Token

// 	// 从输入的第一个字符开始尝试匹配
// 	for _, entry := range s.DFAList {
// 		// 这里我们尝试从字符长度为 1 到最大长度进行匹配
// 		for i := 1; i <= len(input); i++ {
// 			sub := input[:i] // 获取从开头到第i个字符的子字符串
// 			ok, _ := entry.DFA.MatchDFA(sub)
// 			fmt.Printf("[Scan] trying DFA token %s on input '%s'\n", entry.TokenType, sub)

// 			if ok {
// 				if i > maxLen { // 如果当前匹配的长度大于之前的最大匹配长度，更新 maxToken
// 					maxLen = i
// 					maxToken = dfa.Token{Type: entry.TokenType, Lexeme: sub}
// 				}
// 			}
// 		}
// 	}

// 	// 如果没有匹配，检查是否为空白符或者错误字符
// 	if maxLen == 0 {
// 		runes := []rune(input)
// 		if len(runes) > 0 && isWhitespace(runes[0]) {
// 			return dfa.Token{Type: dfa.TokenWithespace, Lexeme: string(runes[0])}, 1
// 		}
// 		return dfa.Token{Type: dfa.TokenERROR, Lexeme: string(input[0])}, 1
// 	}

// 	return maxToken, maxLen
// }
