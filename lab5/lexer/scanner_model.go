package lexer

var keywords = map[string]TokenType{
	"int":      TokenTYPE_KW,
	"float":    TokenTYPE_KW,
	"double":   TokenTYPE_KW,
	"char":     TokenTYPE_KW,
	"string":   TokenTYPE_KW,
	"bool":     TokenTYPE_KW,

	"return":   TokenKW,
	"if":       TokenKW,
	"else":     TokenKW,
	"for":      TokenKW,
	"while":    TokenKW,
	"break":    TokenKW,
	"continue": TokenKW,
	"void":     TokenKW,


	// 按需继续添加其他C语言关键字
}

type Scanner struct {
	DFAList []struct {
		DFA       *DFA
		TokenType TokenType
	}
}
