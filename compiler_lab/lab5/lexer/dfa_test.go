package lexer

import (
	"fmt"
	"os"
	"testing"
	
)

func TestDFA(t *testing.T) {
	dfaWithTokenType, err := LoadMultiDFAFromJson("../json/all_dfa.json", "./dot/test", true)
	if err != nil {
		fmt.Println("Error loading DFA:", err)
		os.Exit(1)
	}

	for _, entry := range *dfaWithTokenType {
		fmt.Printf("\nTesting DFA for token type: %s\n", entry.TokenType)
		fmt.Printf("Alphabet: %v\n", entry.DFA.Alphabet)

		// 测试单字符匹配
		for _, symbol := range entry.DFA.Alphabet {
			ok, _ := entry.DFA.MatchDFA(symbol, true)
			fmt.Printf("Single symbol '%s': matched=%v\n", symbol, ok)
		}

		// 测试复杂词素，针对不同token类型
		testWords := []string{}
		switch entry.TokenType {
		case TokenID:
			testWords = []string{"int", "x", "var1", "_temp"}
		case TokenNUM:
			testWords = []string{"0", "123", "4567"}
		case TokenFLO:
			testWords = []string{"3.1e5", "0.123", ".5", "1e10", "6.022E23"}
		case TokenOP:
			testWords = []string{"=", "==", "+", "+=", "!"}
		case TokenDELIM:
			testWords = []string{"(", ")", "{", "}", ";", ","}
		}

		for _, w := range testWords {
			ok, trace := entry.DFA.MatchDFA(w, true)
			fmt.Printf("Word '%s': matched=%v\n", w, ok)
			for _, step := range trace {
				fmt.Printf("  %s --%s--> %s\n", step.From, step.Symbol, step.To)
			}
		}
	}
}
