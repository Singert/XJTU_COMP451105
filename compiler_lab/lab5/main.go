package main

import (
	"flag"
	"lab5/lexer"
	"lab5/parser"
	"lab5/syntax"
	"lab5/utils"
	"os"
)

func main() {
	sf := flag.String("source", "./assets/source.c", "Source code file to parse")
	flag.StringVar(sf, "source", "./assets/source.c", "Source code file to parse")
	flag.Parse()

	// 1. 加载 DFA
	dfas, _ := lexer.LoadMultiDFAFromJson("./assets/all_dfa.json", "./dot", false)
	scanner := lexer.NewScanner()
	for _, d := range *dfas {
		scanner.RegisterDFA(d.DFA, d.TokenType)
	}

	// 2. 读取源代码
	codeBytes, _ := os.ReadFile(*sf)
	code := string(codeBytes)

	// 3. Tokenize → Symbol
	tokens := scanner.Tokenize(code, true)
	symbols := utils.TokensToSymbols(tokens)

	// 4. 构造文法
	g := syntax.DefineGrammar()

	// 5. 构造分析表
	dfa := parser.BuildDFA(g)
	follow := syntax.ComputeFollow(g)
	table := parser.BuildParseTable(g, dfa, follow)

	// 6. 调用 parser + 语义分析
	parser.Run(symbols, g, dfa, table, tokens, false, (*sf))
}

// ---

// TODO:是否需要我帮你生成一个 ExportASTToDot() 函数来将 AST 渲染为 Graphviz 图？是否想先测试 int x = 3; 和 return x + 1; 等语句 AST 输出？我可以提供测试样例
