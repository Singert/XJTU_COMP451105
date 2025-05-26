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
	sourcefile := flag.String("s", "", "Source file path (default: read from stdin)")
	tokfile := flag.String("t", "./default.tok", "Token file path (default: ./default.tok)")
	flag.StringVar(sourcefile, "source", "", "Source file path (default: read from stdin)")
	flag.BoolVar(verbose, "verbose", false, "Verbose output")
	flag.StringVar(dotPath, "dot", "./dot", "Path to save dot file")
	flag.StringVar(tokfile, "token", "./default.tok", "Token file path (default: ./default.tok)")
	flag.Parse()

	var tok *os.File
	var err error
	if *tokfile != "" {
		tok, err = os.Create(*tokfile)
		if err != nil {
			fmt.Printf("Failed to create output file: %v\n", err)
			return
		}
		defer tok.Close()
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
		scanner.ScanAndOutput(newScanner, code, *dotPath, tok, *verbose)
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

		fmt.Printf("parsing input: %s\n", line)
		scanner.ScanAndOutput(newScanner, line, *dotPath, tok, *verbose)

	}
}
