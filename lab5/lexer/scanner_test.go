package lexer_test

import (
	"lab5/lexer"
	"testing"
)

func TestScanner(t *testing.T) {
	// 加载所有DFA并注册
	dfaWithTokenType, err := lexer.LoadMultiDFAFromJson("../json/all_json", "./dot/test", true)
	if err != nil {
		t.Fatalf("Failed to load DFA: %v", err)
	}

	scanner := lexer.NewScanner()
	for i := range *dfaWithTokenType {
		scanner.RegisterDFA((*dfaWithTokenType)[i].DFA, (*dfaWithTokenType)[i].TokenType)
	}

	tests := []struct {
		input    string
		expected []lexer.Token
	}{
		{
			input: "int x = 3.15;",
			expected: []lexer.Token{
				{Type: lexer.TokenID, Lexeme: "int"},
				{Type: lexer.TokenID, Lexeme: "x"},
				{Type: lexer.TokenOP, Lexeme: "="},
				{Type: lexer.TokenFLO, Lexeme: "3.15"},
				{Type: lexer.TokenDELIM, Lexeme: ";"},
			},
		},
		{
			input: "var1+=42",
			expected: []lexer.Token{
				{Type: lexer.TokenID, Lexeme: "var1"},
				{Type: lexer.TokenOP, Lexeme: "+="},
				{Type: lexer.TokenNUM, Lexeme: "42"},
			},
		},
		{
			input: "\t \n  \r",
			expected: []lexer.Token{
				{Type: lexer.TokenWithespace, Lexeme: "\t \n  \r"},
			},
		},
		{
			input: "@unknown",
			expected: []lexer.Token{
				{Type: lexer.TokenERROR, Lexeme: "@"},
			},
		},
	}

	for _, tt := range tests {
		pos := 0
		inputRunes := []rune(tt.input)
		var gotTokens []lexer.Token

		for pos < len(inputRunes) {
			token, length, _, _ := scanner.Scan(string(inputRunes[pos:]))
			if length == 0 {
				t.Fatalf("Zero length token detected, possible infinite loop for input: %s", tt.input)
			}
			pos += length
			if token.Type == lexer.TokenWithespace {
				continue // 忽略空白token
			}
			gotTokens = append(gotTokens, token)
		}

		if len(gotTokens) != len(tt.expected) {
			t.Errorf("Input %q: expected %d tokens, got %d", tt.input, len(tt.expected), len(gotTokens))
			continue
		}
		for i := range gotTokens {
			if gotTokens[i].Type != tt.expected[i].Type || gotTokens[i].Lexeme != tt.expected[i].Lexeme {
				t.Errorf("Input %q: token %d expected (%s, %q), got (%s, %q)",
					tt.input, i, tt.expected[i].Type, tt.expected[i].Lexeme, gotTokens[i].Type, gotTokens[i].Lexeme)
			}
		}
	}
}
