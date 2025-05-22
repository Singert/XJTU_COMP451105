// main.go
package main

import (
	"fmt"
	"os"
	"project/lexer"
	"project/parser"
)

func main() {
	input, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}

	tokens := lexer.Tokenize(string(input))
	tac := parser.ParseAndGenerateTAC(tokens)

	for _, line := range tac {
		fmt.Println(line)
	}
}
