package scanner

import (
	"fmt"
	"lab2/core/dfa"
	"os"
	"unicode"
)

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

	if maxToken.Type == dfa.TokenID {
		if keywordType, ok := keywords[maxToken.Lexeme]; ok {
			maxToken.Type = keywordType
		}
	}

	_, _ = maxDFA.MatchDFA(maxToken.Lexeme, true)
	return maxToken, maxLen
}

func ScanAndOutput(scanner *Scanner, input string, dotPath string, out *os.File) {
	pos := 0
	inputRunes := []rune(input)
	length := len(inputRunes)

	for pos < length {
		fmt.Printf("[DEBUG] pos=%d, next char='%c'\n", pos, inputRunes[pos])

		subInput := string(inputRunes[pos:])
		token, tokenLen := scanner.Scan(subInput)
		fmt.Printf("[DEBUG] token='%s', length=%d\n", token.Lexeme, tokenLen)

		if tokenLen == 0 {
			pos++ // 防止死循环
			continue
		}

		if token.Type == dfa.TokenWithespace {
			pos += tokenLen
			fmt.Printf("[main] Skip %d whitespace characters\n", tokenLen)
			continue
		}

		if token.Type == dfa.TokenERROR {
			fmt.Printf("❌ Error: invalid token '%s' at position %d\n", token.Lexeme, pos)
			pos += tokenLen
			continue
		}

		fmt.Fprintf(out, "[Token]: <%s>, [Lexeme]: '%s'\n", token.Type, token.Lexeme)
		fmt.Printf("[Token]: <%s>, [Lexeme]: '%s'\n", token.Type, token.Lexeme)

		_, trace := scanner.DFAList[0].DFA.MatchDFA(token.Lexeme, false)
		dotName := fmt.Sprintf("%s/%s_%d.dot", dotPath, token.Lexeme, pos)
		err := scanner.DFAList[0].DFA.ExportToDot(dotName, trace)
		if err != nil {
			fmt.Println("Export dot failed:", err)
		}

		pos += tokenLen
	}
}
