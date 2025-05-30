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
	case lexer.TokenNUM:
		return "num"
	case lexer.TokenFLO:
		return "float"
	case lexer.TokenCHAR:
		return "char"
	case lexer.TokenSTRING:
		return "string"
	case lexer.TokenOP, lexer.TokenDELIM, lexer.TokenKW:
		return syntax.Symbol(tok.Lexeme)
	case lexer.TokenTYPE_KW:
		return "type_kw" // Type keyword
	case lexer.TokenEOF:
		return "#"
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

func SymbolToTokenType(sym syntax.Symbol) lexer.TokenType {
	switch sym {
	case "id":
		return lexer.TokenID
	case "num":
		// 这里返回 TokenNUM 作为默认，也可以根据需求返回 TokenFLO
		return lexer.TokenNUM
	// 关键字符号
	case "return", "if", "else", "for", "while", "break", "continue", "void":
		return lexer.TokenKW
	// 操作符和分隔符，简化示例可全部归为 TokenOP
	case "int", "float", "double", "char", "type_kw", "string", "bool":
		return lexer.TokenTYPE_KW // Type keyword
	case "+", "-", "*", "/", "=", "==", "<", ">", "&&", "||", "!", ";", ",", "(", ")", "[", "]", "{", "}":
		// 这里示例统一返回 TokenOP，实际可根据字符区分 OP 或 DELIM
		// 也可以使用一个 map 辅助判别
		return lexer.TokenOP
	case "#":
		return lexer.TokenEOF
	default:
		// 未知或错误符号
		return lexer.TokenERROR
	}
}

func SymbolsToTokenTypes(symbols []syntax.Symbol) []lexer.TokenType {
	var tokenTypes []lexer.TokenType
	for _, sym := range symbols {
		tokenTypes = append(tokenTypes, SymbolToTokenType(sym))
	}
	return tokenTypes
}
