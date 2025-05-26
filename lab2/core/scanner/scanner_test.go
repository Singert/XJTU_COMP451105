package scanner_test

import (
	"lab2/core/dfa"
	"lab2/core/scanner"
	"testing"
)

func TestScanner(t *testing.T) {
	// 加载所有DFA并注册
	dfaWithTokenType, err := dfa.LoadMultiDFAFromJson("../json/all_dfa.json", "./dot/test", true)
	if err != nil {
		t.Fatalf("Failed to load DFA: %v", err)
	}

	scanner := scanner.NewScanner()
	for i := range *dfaWithTokenType {
		scanner.RegisterDFA((*dfaWithTokenType)[i].DFA, (*dfaWithTokenType)[i].TokenType)
	}

	tests := []struct {
		input    string
		expected []dfa.Token
	}{
		{
			input: "int x = 3.15;",
			expected: []dfa.Token{
				{Type: dfa.TokenID, Lexeme: "int"},
				{Type: dfa.TokenID, Lexeme: "x"},
				{Type: dfa.TokenOP, Lexeme: "="},
				{Type: dfa.TokenFLO, Lexeme: "3.15"},
				{Type: dfa.TokenDELIM, Lexeme: ";"},
			},
		},
		{
			input: "var1+=42",
			expected: []dfa.Token{
				{Type: dfa.TokenID, Lexeme: "var1"},
				{Type: dfa.TokenOP, Lexeme: "+="},
				{Type: dfa.TokenNUM, Lexeme: "42"},
			},
		},
		{
			input: "\t \n  \r",
			expected: []dfa.Token{
				{Type: dfa.TokenWithespace, Lexeme: "\t \n  \r"},
			},
		},
		{
			input: "@unknown",
			expected: []dfa.Token{
				{Type: dfa.TokenERROR, Lexeme: "@"},
			},
		},
	}

	for _, tt := range tests {
		pos := 0
		inputRunes := []rune(tt.input)
		var gotTokens []dfa.Token

		for pos < len(inputRunes) {
			token, length ,_,_:= scanner.Scan(string(inputRunes[pos:]))
			if length == 0 {
				t.Fatalf("Zero length token detected, possible infinite loop for input: %s", tt.input)
			}
			pos += length
			if token.Type == dfa.TokenWithespace {
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
