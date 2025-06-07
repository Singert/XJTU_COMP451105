package lexer

var keywords = map[string]TokenType{
	"int":    TokenTYPE_KW,
	"float":  TokenTYPE_KW,
	"double": TokenTYPE_KW,
	"char":   TokenTYPE_KW,
	"string": TokenTYPE_KW,
	"bool":   TokenTYPE_KW,
	"void":   TokenTYPE_KW,

	"return":   TokenKW,
	"if":       TokenKW,
	"else":     TokenKW,
	"for":      TokenKW,
	"while":    TokenKW,
	"break":    TokenKW,
	"continue": TokenKW,
}

type Scanner struct {
	DFAList []struct {
		DFA       *DFA
		TokenType TokenType
	}
}
