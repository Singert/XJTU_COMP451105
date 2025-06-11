package main

import (
	"fmt"
	"os"
	"slr_semantic/lexer"
	"slr_semantic/parser"
	"slr_semantic/semantic"
)


func main() {
	input := readSource("data/input.txt")
	raw := lexer.Tokenize(input)

	richTokens := tokenizeRich(raw)
	kinds := make([]string, len(richTokens))
	for i, t := range richTokens {
		kinds[i] = t.Kind
	}

	grammar := parser.LoadGrammar()
	table := parser.GenerateSLRTable(grammar)
	symtab := semantic.NewSymbolTable()

	parser.RunParser(kinds, table, grammar, symtab, richTokens)

	fmt.Println()
	symtab.Dump()
}

func readSource(filename string) string {
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func tokenizeRich(tokens []string) []lexer.Token {
	var result []lexer.Token
	for _, tok := range tokens {
		if isKeyword(tok) || isSymbol(tok) {
			result = append(result, lexer.Token{Kind: tok, Value: tok})
		} else if isNumber(tok) {
			result = append(result, lexer.Token{Kind: "num", Value: tok})
		} else {
			result = append(result, lexer.Token{Kind: "id", Value: tok})
		}
	}
	return result
}

func isKeyword(tok string) bool {
	return tok == "int" || tok == "void"
}

func isSymbol(tok string) bool {
	symbols := map[string]bool{
		"[": true, "]": true, "(": true, ")": true, ";": true, ",": true,
	}
	return symbols[tok]
}

func isNumber(tok string) bool {
	for _, r := range tok {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
