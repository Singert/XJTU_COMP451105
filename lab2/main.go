package main

import (
	"bufio"
	"flag"
	"fmt"
	"lab2/core"
	"lab2/core/dfa"
	"os"
	_ "os/exec"
	"strings"
)

func main() {
	// define command line arguments
	dotPath := flag.String("d", "./dot", "Path to save dot file")
	flag.StringVar(dotPath, "dot", "./dot", "Path to save dot file")
	flag.Parse()

	dfaWithTokenType, err := dfa.LoadMultiDFAFromJson("./json/all_dfa.json", *dotPath)
	if err != nil {
		fmt.Println("Error loading DFA:", err)
		os.Exit(1)
	}

	scanner := core.NewScanner()
	for i := range *dfaWithTokenType {
		scanner.RegisterDFA((*dfaWithTokenType)[i].DFA, (*dfaWithTokenType)[i].TokenType)
	}

	var reader = bufio.NewReader(os.Stdin)

	for {

		fmt.Print("Enter a string to match (or 'quit' to quit): ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}
		line = strings.TrimSpace(line)

		if line == "quit" {
			break
		}
		// Matched, trace := dfa.MatchDFA(input)
		// if Matched {
		// 	fmt.Println("✅ Accepted")
		// 	dotName = *dotPath + "/" + input + ".dot"
		// 	dfa.ExportToDot(dotName, trace)

		// } else {
		// 	fmt.Println("❌ Not Accepted")
		// }

		// pos := 0
		// for pos < len(input) {
		// 	token, length := scanner.Scan(input[pos:])
		// 	if token.Type == dfa.TokenWithespace {
		// 		pos += length
		// 		fmt.Printf("[main:when token with space]Skip %d whitespace\n", length)
		// 		continue
		// 	}
		// 	if token.Type == dfa.TokenERROR {
		// 		fmt.Printf("❌ Error: invalid token '%s' at position %d\n", token.Lexeme, pos)
		// 		pos += length
		// 		continue
		// 	}

		// 	fmt.Printf("[Token]: <%s>, [Lexeme]: '%s'\n", token.Type, token.Lexeme)
		// 	// 可选：导出匹配的 dot 文件，示例用第一个 DFA (你可扩展对应匹配 DFA)
		// 	_, trace := scanner.DFAList[0].DFA.MatchDFA(token.Lexeme,false)
		// 	dotName := fmt.Sprintf("%s/%s_%d.dot", *dotPath, token.Lexeme, pos)
		// 	err := scanner.DFAList[0].DFA.ExportToDot(dotName, trace)
		// 	if err != nil {
		// 		fmt.Println("Export dot failed:", err)
		// 	}

		// 	pos += length
		// }
		pos := 0
		inputRunes := []rune(line)
		length := len(inputRunes)

		for pos < length {
			fmt.Printf("[DEBUG] pos=%d, next char='%c'\n", pos, inputRunes[pos])

			subInput := string(inputRunes[pos:])
			token, tokenLen := scanner.Scan(subInput)
			fmt.Printf("[DEBUG] token='%s', length=%d\n", token.Lexeme, tokenLen)

			if tokenLen == 0 {
				pos++ // 防止死循环
				continue
			}

			if token.Type == dfa.TokenWithespace {
				pos += tokenLen
				fmt.Printf("[main] Skip %d whitespace characters\n", tokenLen)
				continue
			}

			if token.Type == dfa.TokenERROR {
				fmt.Printf("❌ Error: invalid token '%s' at position %d\n", token.Lexeme, pos)
				pos += tokenLen
				continue
			}

			fmt.Printf("[Token]: <%s>, [Lexeme]: '%s'\n", token.Type, token.Lexeme)
			_, trace := scanner.DFAList[0].DFA.MatchDFA(token.Lexeme, false)
			dotName := fmt.Sprintf("%s/%s_%d.dot", *dotPath, token.Lexeme, pos)
			err := scanner.DFAList[0].DFA.ExportToDot(dotName, trace)
			if err != nil {
				fmt.Println("Export dot failed:", err)
			}

			pos += tokenLen
		}

	}
}
