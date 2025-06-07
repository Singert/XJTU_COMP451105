package test

import (
	"fmt"
	"lexer/scanner"
)

func TestManual() {
	code := `
int fact(int n, int a) {
    if (n <= 1) return a;
    else return fact(n - 1, n * a);
}`
	sc := scanner.NewScanner(code)
	tokens := sc.Scan()

	for _, t := range tokens {
		fmt.Printf("(%s, \"%s\") at [%d:%d]\n", t.Type, t.Lexeme, t.Line, t.Column)
	}
}
