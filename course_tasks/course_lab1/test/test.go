package test

import (
	"fmt"
	"lexer/scanner"
)

func TestManual() {
	code := `
int raw(int x) {
    y = x + 5;
    return y;
}

void foo(int y) {
    int z;
    void bar(int x, int soo()) {
        if (x > 3)
            bar(x / 3, soo);
        else
            z = soo(x);
        print z;
    }
    bar(y, raw);
}

foo(6);`
	sc := scanner.NewScanner(code)
	tokens := sc.Scan()

	for _, t := range tokens {
		fmt.Printf("(%s, \"%s\") at [%d:%d]\n", t.Type, t.Lexeme, t.Line, t.Column)
	}
}
