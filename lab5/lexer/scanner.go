package lexer

import (
	"fmt"
	"lab5/syntax"
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

	// 优先检测注释
	if len(runes) >= 2 && runes[0] == '/' && (runes[1] == '/' || runes[1] == '*') {
		token, length := s.scanComment(runes)
		return token, length, nil, nil
	}
	// 检测空白字符
	if unicode.IsSpace(runes[0]) {
		i := 1
		for i < len(runes) && unicode.IsSpace(runes[i]) {
			i++
		}
		return Token{Type: TokenWithespace, Lexeme: string(runes[:i])}, i, nil, nil
	}
	// 检测字符串和字符字面量
	if runes[0] == '\'' || runes[0] == '"' {
		token, len := s.scanCharOrString(runes)
		return token, len, nil, nil
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

	_, _ = maxDFA.MatchDFA(maxToken.Lexeme, false)
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

func ScanAndOutputWithStream(scanner *Scanner, input string, dotPath string, tok *os.File, verbose bool) []Token {
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
		tokens = append(tokens, token)
		dotName := fmt.Sprintf("%s/%s_%d.dot", dotPath, token.Lexeme, pos)
		err := matchedDFA.ExportToDot(dotName, trace)
		if err != nil {
			fmt.Println("Export dot failed:", err)
		}

		pos += tokenLen
	}
	return tokens
}

// Tokenize 扫描输入字符串并生成 Token 列表
func (s *Scanner) Tokenize(input string, verbose bool) []Token {
	tokens := []Token{}
	pos := 0
	inputRunes := []rune(input)
	length := len(inputRunes)

	line := 1
	col := 1

	fmt.Println("===TOKENIZE starting===")
	for pos < length {
		subInput := string(inputRunes[pos:])
		token, tokenLen, _, _ := s.Scan(subInput)

		// 赋值行列给token
		token.Line = line
		token.Column = col

		if tokenLen == 0 {
			pos++
			col++
			continue
		}

		if token.Type == TokenWithespace {
			// 统计空白中的换行数量，更新行列
			for i := 0; i < tokenLen; i++ {
				if inputRunes[pos+i] == '\n' {
					line++
					col = 1
				} else {
					col++
				}
			}
			pos += tokenLen
			continue
		}

		if token.Type == TokenERROR {
			fmt.Printf("❌ Error: invalid token '%s' at line %d, column %d\n", token.Lexeme, line, col)
			// 同样更新行列
			for i := 0; i < tokenLen; i++ {
				if inputRunes[pos+i] == '\n' {
					line++
					col = 1
				} else {
					col++
				}
			}
			pos += tokenLen
			continue
		}

		if verbose {
			fmt.Printf("[Token]: <%s>, [Lexeme]: <%s> [symbol]:<%s> at line %d, column %d\n",
				token.Type, token.Lexeme, tokenToSymbol(token), token.Line, token.Column)
		}
		tokens = append(tokens, token)

		// 更新行列，统计token中的换行
		for i := 0; i < tokenLen; i++ {
			if inputRunes[pos+i] == '\n' {
				line++
				col = 1
			} else {
				col++
			}
		}
		pos += tokenLen
	}

	// 添加EOF token，行列取最后一个token结束处
	if len(tokens) > 0 {
		lastTok := tokens[len(tokens)-1]
		eofToken := Token{
			Type:   "EOF",
			Lexeme: "EOF",
			Line:   lastTok.Line,
			Column: lastTok.Column + len([]rune(lastTok.Lexeme)),
		}
		tokens = append(tokens, eofToken)
	} else {
		eofToken := Token{
			Type:   "EOF",
			Lexeme: "EOF",
			Line:   1,
			Column: 1,
		}
		tokens = append(tokens, eofToken)
	}
	fmt.Println("===TOKENIZE completed <manually add EOF token>===")
	return tokens
}

// TokenToSymbol 将 Token 映射为语法分析器的 Symbol（终结符）
func tokenToSymbol(tok Token) syntax.Symbol {
	switch tok.Type {
	case TokenID:
		return "id"
	case TokenNUM:
		return "num"
	case TokenFLO:
		return "float"
	case TokenCHAR:
		return "char"
	case TokenSTRING:
		return "string"
	case TokenOP, TokenDELIM, TokenKW:
		return syntax.Symbol(tok.Lexeme)
	case TokenCOMMENT_SINGLE, TokenCOMMENT_MULTI:
		return "comment" // 注释不作为语法分析的终结符
	case TokenTYPE_KW:
		return "type_kw" // Type keyword
	case TokenEOF:
		return "#"
	default:
		return "?" // 识别失败
	}
}
