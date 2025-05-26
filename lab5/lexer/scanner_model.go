package lexer

var keywords = map[string]TokenType{
	"int":      TokenKW,
	"return":   TokenKW,
	"if":       TokenKW,
	"else":     TokenKW,
	"for":      TokenKW,
	"while":    TokenKW,
	"break":    TokenKW,
	"continue": TokenKW,
	"void":     TokenKW,
	"char":     TokenKW,
	"float":    TokenKW,
	"double":   TokenKW,
	// 按需继续添加其他C语言关键字
}

type Scanner struct {
	DFAList []struct {
		DFA       *DFA
		TokenType TokenType
	}
}
