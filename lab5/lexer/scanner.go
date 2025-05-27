package lexer

import (
	"fmt"
	"os"
	"unicode"
)

func NewScanner() *Scanner {
	return &Scanner{}
}

func (s *Scanner) RegisterDFA(d *DFA, t TokenType) {
	s.DFAList = append(s.DFAList, struct {
		DFA       *DFA
		TokenType TokenType
	}{DFA: d, TokenType: t})
}

func isWhitespace(r rune) bool {
	return unicode.IsSpace(r)
}

func (s *Scanner) Scan(input string) (matched Token, length int, matchedDFA *DFA, trace []TransitionTrace) {
	if len(input) == 0 {
		return Token{Type: TokenERROR, Lexeme: ""}, 0, nil, nil
	}

	runes := []rune(input)
	if unicode.IsSpace(runes[0]) {
		i := 1
		for i < len(runes) && unicode.IsSpace(runes[i]) {
			i++
		}
		return Token{Type: TokenWithespace, Lexeme: string(runes[:i])}, i, nil, nil
	}

	maxLen := 0
	var maxToken Token
	var maxDFA *DFA
	var maxTrace []TransitionTrace

	for _, entry := range s.DFAList {
		for i := 1; i <= len(runes); i++ {
			sub := string(runes[:i])
			ok, trace := entry.DFA.MatchDFA(sub, false)
			if ok && i > maxLen {
				maxLen = i
				maxToken = Token{Type: entry.TokenType, Lexeme: sub}
				maxDFA = entry.DFA
				maxTrace = trace
			}
		}
	}

	if maxLen == 0 {
		return Token{Type: TokenERROR, Lexeme: string(runes[0])}, 1, nil, nil
	}

	if maxToken.Type == TokenID {
		if keywordType, ok := keywords[maxToken.Lexeme]; ok {
			maxToken.Type = keywordType
		}
	}

	_, _ = maxDFA.MatchDFA(maxToken.Lexeme, true)
	return maxToken, maxLen, maxDFA, maxTrace
}

func ScanAndOutput(scanner *Scanner, input string, dotPath string, tok *os.File, verbose bool) {
	pos := 0
	inputRunes := []rune(input)
	length := len(inputRunes)

	for pos < length {
		if verbose {
			fmt.Printf("[DEBUG] pos=%d, next char='%c'\n", pos, inputRunes[pos])
		}

		subInput := string(inputRunes[pos:])
		token, tokenLen, matchedDFA, trace := scanner.Scan(subInput)
		if verbose {
			fmt.Printf("[DEBUG] token='%s', length=%d\n", token.Lexeme, tokenLen)
		}
		if tokenLen == 0 {
			pos++ // 防止死循环
			continue
		}

		if token.Type == TokenWithespace {
			pos += tokenLen
			fmt.Printf("[main] Skip %d whitespace characters\n", tokenLen)
			continue
		}

		if token.Type == TokenERROR {
			fmt.Printf("❌ Error: invalid token '%s' at position %d\n", token.Lexeme, pos)
			pos += tokenLen
			continue
		}

		fmt.Fprintf(tok, "%s %s\n", token.Type, token.Lexeme)
		fmt.Printf("[Token]: <%s>, [Lexeme]: '%s'\n", token.Type, token.Lexeme)

		dotName := fmt.Sprintf("%s/%s_%d.dot", dotPath, token.Lexeme, pos)
		err := matchedDFA.ExportToDot(dotName, trace)
		if err != nil {
			fmt.Println("Export dot failed:", err)
		}

		pos += tokenLen
	}
}

func ScanAndOutputWithStream(scanner *Scanner, input string, dotPath string, tok *os.File, verbose bool)  []Token {
	tokens := []Token{}
	pos := 0
	inputRunes := []rune(input)
	length := len(inputRunes)

	for pos < length {
		if verbose {
			fmt.Printf("[DEBUG] pos=%d, next char='%c'\n", pos, inputRunes[pos])
		}

		subInput := string(inputRunes[pos:])
		token, tokenLen, matchedDFA, trace := scanner.Scan(subInput)
		if verbose {
			fmt.Printf("[DEBUG] token='%s', length=%d\n", token.Lexeme, tokenLen)
		}
		if tokenLen == 0 {
			pos++ // 防止死循环
			continue
		}

		if token.Type == TokenWithespace {
			pos += tokenLen
			fmt.Printf("[main] Skip %d whitespace characters\n", tokenLen)
			continue
		}

		if token.Type == TokenERROR {
			fmt.Printf("❌ Error: invalid token '%s' at position %d\n", token.Lexeme, pos)
			pos += tokenLen
			continue
		}

		fmt.Fprintf(tok, "%s %s\n", token.Type, token.Lexeme)
		fmt.Printf("[Token]: <%s>, [Lexeme]: '%s'\n", token.Type, token.Lexeme)

		dotName := fmt.Sprintf("%s/%s_%d.dot", dotPath, token.Lexeme, pos)
		err := matchedDFA.ExportToDot(dotName, trace)
		if err != nil {
			fmt.Println("Export dot failed:", err)
		}

		pos += tokenLen
	}
}