package main

import (
	"fmt"
	"lexer/scanner"
	"lexer/util"
	"os"

	"lexer/test"
)

func main() {
	test.TestManual()
	if len(os.Args) < 2 {
		fmt.Println("Usage: lexer <source-file>")
		return
	}

	srcPath := os.Args[1]
	content, err := util.ReadFile(srcPath)
	if err != nil {
		fmt.Println("Failed to read file:", err)
		return
	}

	sc := scanner.NewScanner(content)
	tokens := sc.Scan()

	for _, t := range tokens {
		fmt.Printf("(%s, \"%s\") at [%d:%d]\n", t.Type, t.Lexeme, t.Line, t.Column)
	}
}
