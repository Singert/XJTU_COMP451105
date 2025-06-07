package scanner

import (
	"lexer/dfa"
	"lexer/token"
	"unicode"
)

type Scanner struct {
	input  []rune
	pos    int
	line   int
	column int

	idDFA    *dfa.DFA
	numDFA   *dfa.DFA
	keywords map[string]token.TokenType
}

func NewScanner(src string) *Scanner {
	return &Scanner{
		input:  []rune(src),
		pos:    0,
		line:   1,
		column: 1,
		idDFA:  dfa.BuildIDOrKeywordDFA(nil),
		numDFA: dfa.BuildNumberDFA(),
		keywords: map[string]token.TokenType{
			"int": token.INT, "void": token.VOID, "if": token.IF,
			"else": token.ELSE, "while": token.WHILE, "return": token.RETURN,
		},
	}
}

// 获取下一个字符，不推进
func (s *Scanner) peek() rune {
	if s.pos >= len(s.input) {
		return 0
	}
	return s.input[s.pos]
}

// 获取下一个字符并推进
func (s *Scanner) advance() rune {
	ch := s.peek()
	if ch == 0 {
		return 0
	}
	s.pos++
	if ch == '\n' {
		s.line++
		s.column = 1
	} else {
		s.column++
	}
	return ch
}

// 回退一个字符
// func (s *Scanner) retreat() {
// 	if s.pos == 0 {
// 		return
// 	}
// 	s.pos--
// 	if s.input[s.pos] == '\n' {
// 		s.line--
// 	}
// 	s.column--
// }

// 跳过空白字符
func (s *Scanner) skipWhitespace() {
	for unicode.IsSpace(s.peek()) {
		s.advance()
	}
}

// 主扫描函数
func (s *Scanner) Scan() []token.Token {
	var tokens []token.Token

	for {
		s.skipWhitespace()
		ch := s.peek()
		if ch == 0 {
			tokens = append(tokens, token.Token{Type: token.EOF, Lexeme: "", Line: s.line, Column: s.column})
			break
		}

		if unicode.IsLetter(ch) || ch == '_' {
			tokens = append(tokens, s.scanIdentifier())
		} else if unicode.IsDigit(ch) {
			tokens = append(tokens, s.scanNumber())
		} else {
			tokens = append(tokens, s.scanSymbol())
		}
	}

	return tokens
}

// ================================
func (s *Scanner) scanIdentifier() token.Token {
	start := s.pos
	startCol := s.column
	for unicode.IsLetter(s.peek()) || unicode.IsDigit(s.peek()) || s.peek() == '_' {
		s.advance()
	}

	lexeme := string(s.input[start:s.pos])
	tokType := token.ID
	if kw, ok := s.keywords[lexeme]; ok {
		tokType = kw
	}

	return token.Token{
		Type:   tokType,
		Lexeme: lexeme,
		Line:   s.line,
		Column: startCol,
	}
}

func (s *Scanner) scanNumber() token.Token {
	start := s.pos
	startCol := s.column
	for unicode.IsDigit(s.peek()) {
		s.advance()
	}

	return token.Token{
		Type:   token.NUM,
		Lexeme: string(s.input[start:s.pos]),
		Line:   s.line,
		Column: startCol,
	}
}

func (s *Scanner) scanSymbol() token.Token {
	ch := s.advance()
	startCol := s.column - 1

	switch ch {
	case ';':
		return s.makeTok(token.SCO, string(ch), startCol)
	case ',':
		return s.makeTok(token.CMA, string(ch), startCol)
	case '{':
		return s.makeTok(token.LBR, string(ch), startCol)
	case '}':
		return s.makeTok(token.RBR, string(ch), startCol)
	case '(':
		return s.makeTok(token.LPA, string(ch), startCol)
	case ')':
		return s.makeTok(token.RPA, string(ch), startCol)
	case '+':
		return s.makeTok(token.ADD, string(ch), startCol)
	case '*':
		return s.makeTok(token.MUL, string(ch), startCol)
	case '&':
		if s.peek() == '&' {
			s.advance()
			return s.makeTok(token.AND, "&&", startCol)
		}
	case '|':
		if s.peek() == '|' {
			s.advance()
			return s.makeTok(token.OR, "||", startCol)
		}
	case '<':
		if s.peek() == '=' {
			s.advance()
			return s.makeTok(token.ROP, "<=", startCol)
		}
		return s.makeTok(token.ROP, "<", startCol)
	case '=':
		if s.peek() == '=' {
			s.advance()
			return s.makeTok(token.ROP, "==", startCol)
		}
	}

	return token.Token{Type: token.ILLEGAL, Lexeme: string(ch), Line: s.line, Column: startCol}
}

func (s *Scanner) makeTok(t token.TokenType, lexeme string, col int) token.Token {
	return token.Token{
		Type:   t,
		Lexeme: lexeme,
		Line:   s.line,
		Column: col,
	}
}
