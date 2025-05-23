package scanner

import (
	"lab2/core/dfa"
)

var keywords = map[string]dfa.TokenType{
	"int":      dfa.TokenKW,
	"return":   dfa.TokenKW,
	"if":       dfa.TokenKW,
	"else":     dfa.TokenKW,
	"for":      dfa.TokenKW,
	"while":    dfa.TokenKW,
	"break":    dfa.TokenKW,
	"continue": dfa.TokenKW,
	"void":     dfa.TokenKW,
	"char":     dfa.TokenKW,
	"float":    dfa.TokenKW,
	"double":   dfa.TokenKW,
	// 按需继续添加其他C语言关键字
}

type Scanner struct {
	DFAList []struct {
		DFA       *dfa.DFA
		TokenType dfa.TokenType
	}
}
