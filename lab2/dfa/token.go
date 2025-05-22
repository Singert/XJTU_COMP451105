package dfa

type TokenType string

const (
	TokenID    TokenType = "ID"  // Identifier
	TokenNUM   TokenType = "NUM" // Number
	TokenFLO   TokenType = "FLO"
	TokenOP    TokenType = "OP"
	TokenDELIM TokenType = "DELIM"
	TokenKW    TokenType = "KEYWORD"
	TokenERROR TokenType = "ERROR"
)

type Token struct {
	Type   TokenType
	Lexeme string
}
