// Project: lab5
// file: ./parser/error.go
package parser

import (
	"fmt"
	"lab5/lexer"
)

type ParseError struct {
	Line, Column int
	Token        lexer.Token //出错的 Token
	Expected     []lexer.TokenType
	Msg          string
}

func (e *ParseError) Error() string {
	if len(e.Expected) > 0 {
		return fmt.Sprintf("Syntax error at <line %d,column %d>:Unexpected token <%s,%s>,expected one of %v.", e.Line, e.Column, e.Token.Type, e.Token.Lexeme, e.Expected)
	}
	return fmt.Sprintf("Syntax error at <line %d,column %d>:Unexpected token <%s,%s>.", e.Line, e.Column, e.Token.Type, e.Token.Lexeme)
}
