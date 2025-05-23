package main

import (
	"bufio"
	"flag"
	"fmt"
	"lab2/core/dfa"
	"lab2/core/scanner"
	"os"
	_ "os/exec"
	"strings"
)

func main() {
	// define command line arguments
	dotPath := flag.String("d", "./dot", "Path to save dot file")
	verbose := flag.Bool("v", false, "Verbose output")
	outputFile := flag.String("o", "", "Output file path for tokens (default: print to stdout)")
	sourcefile := flag.String("s", "", "Source file path (default: read from stdin)")
	flag.StringVar(sourcefile, "source", "", "Source file path (default: read from stdin)")
	flag.StringVar(outputFile, "output", "", "Output file path for tokens (default: print to stdout)")
	flag.BoolVar(verbose, "verbose", false, "Verbose output")
	flag.StringVar(dotPath, "dot", "./dot", "Path to save dot file")
	flag.Parse()

	var out *os.File
	var err error
	if *outputFile != "" {
		out, err = os.Create(*outputFile)
		if err != nil {
			fmt.Printf("Failed to create output file: %v\n", err)
			return
		}
		defer out.Close()
	} else {
		out = os.Stdout
	}

	dfaWithTokenType, err := dfa.LoadMultiDFAFromJson("./json/all_dfa.json", *dotPath, *verbose)
	if err != nil {
		fmt.Println("Error loading DFA:", err)
		os.Exit(1)
	}

	newScanner := scanner.NewScanner()
	for i := range *dfaWithTokenType {
		newScanner.RegisterDFA((*dfaWithTokenType)[i].DFA, (*dfaWithTokenType)[i].TokenType)
	}

	if *sourcefile != "" {
		fmt.Printf("parse file: %s\n", *sourcefile)
		contentBytes, err := os.ReadFile(*sourcefile)
		if err != nil {
			fmt.Printf("Failed to read source file: %v\n", err)
			os.Exit(1)
		}
		code := string(contentBytes)
		scanner.ScanAndOutput(newScanner, code, *dotPath, out)
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

		fmt.Fprintf(out, "parsing input: %s\n", line)
		scanner.ScanAndOutput(newScanner, line, *dotPath, out)
		// pos := 0
		// inputRunes := []rune(line)
		// length := len(inputRunes)

		// for pos < length {
		// 	fmt.Printf("[DEBUG] pos=%d, next char='%c'\n", pos, inputRunes[pos])

		// 	subInput := string(inputRunes[pos:])
		// 	token, tokenLen := newScanner.Scan(subInput)
		// 	fmt.Printf("[DEBUG] token='%s', length=%d\n", token.Lexeme, tokenLen)

		// 	if tokenLen == 0 {
		// 		pos++ // 防止死循环
		// 		continue
		// 	}

		// 	if token.Type == dfa.TokenWithespace {
		// 		pos += tokenLen
		// 		fmt.Printf("[main] Skip %d whitespace characters\n", tokenLen)
		// 		continue
		// 	}

		// 	if token.Type == dfa.TokenERROR {
		// 		fmt.Printf("❌ Error: invalid token '%s' at position %d\n", token.Lexeme, pos)
		// 		pos += tokenLen
		// 		continue
		// 	}
		// 	fmt.Fprintf(out, "[Token]: <%s>, [Lexeme]: '%s'\n", token.Type, token.Lexeme)
		// 	fmt.Printf("[Token]: <%s>, [Lexeme]: '%s'\n", token.Type, token.Lexeme)
		// 	_, trace := newScanner.DFAList[0].DFA.MatchDFA(token.Lexeme, false)
		// 	dotName := fmt.Sprintf("%s/%s_%d.dot", *dotPath, token.Lexeme, pos)
		// 	err := newScanner.DFAList[0].DFA.ExportToDot(dotName, trace)
		// 	if err != nil {
		// 		fmt.Println("Export dot failed:", err)
		// 	}

		// 	pos += tokenLen
		// }

	}
}
