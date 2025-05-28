package utils

import (
	"lab5/lexer"
	"lab5/syntax"
)

// TokenToSymbol 将 Token 映射为语法分析器的 Symbol（终结符）
func TokenToSymbol(tok lexer.Token) syntax.Symbol {
	switch tok.Type {
	case lexer.TokenID:
		return "id"
	case lexer.TokenNUM, lexer.TokenFLO:
		return "num"
	case lexer.TokenOP, lexer.TokenDELIM, lexer.TokenKW:
		return syntax.Symbol(tok.Lexeme)
	default:
		return "?" // 识别失败
	}
}

// TokensToSymbols 将 Token 列表映射为 Symbol 列表
func TokensToSymbols(tokens []lexer.Token) []syntax.Symbol {
	var symbols []syntax.Symbol
	for _, tok := range tokens {
		symbols = append(symbols, TokenToSymbol(tok))
	}
	return symbols
}
