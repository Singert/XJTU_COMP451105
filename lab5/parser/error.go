// Project: lab5
// file: ./parser/error.go
package parser

import (
	"fmt"
	"lab5/lexer"
	"lab5/syntax"
)

type ParseError struct {
	Line, Column int
	Token        lexer.Token //出错的 Token
	Expected     []syntax.Symbol
	Msg          string
}

func (e *ParseError) Error() string {
	if len(e.Expected) > 0 {
		return fmt.Sprintf("Syntax error at <line %d,column %d>:Unexpected token <%s,%s>,expected one of %v.", e.Line, e.Column, e.Token.Type, e.Token.Lexeme, e.Expected)
	}
	return fmt.Sprintf("Syntax error at <line %d,column %d>:Unexpected token <%s,%s>.", e.Line, e.Column, e.Token.Type, e.Token.Lexeme)
}

func CatchParseError(currState int, currToken syntax.Symbol, tokenStream []lexer.Token, tokIdx int, table *ParseTable) *ParseError {
	// 收集期望符号
	expectedTokens := []syntax.Symbol{}
	for symb := range table.Action[currState] {
		expectedTokens = append(expectedTokens, symb)
	}

	var errTok lexer.Token
	if tokIdx < len(tokenStream) {
		errTok = tokenStream[tokIdx]
	} else {
		errTok = lexer.Token{Type: "EOF", Lexeme: "EOF", Line: -1, Column: -1}
	}

	fmt.Printf("语法错误！在第 %d 行，第 %d 列，发现非法符号: <%s, %s>\n", errTok.Line, errTok.Column, errTok.Type, errTok.Lexeme)
	fmt.Printf("期望符号有: %v\n", expectedTokens)

	return &ParseError{
		Line:     errTok.Line,
		Column:   errTok.Column,
		Token:    errTok,
		Expected: expectedTokens,
		Msg:      "语法分析失败，输入不符合文法规则。",
	}
}
