package main

import (
	"lab5/lexer"
	"lab5/parser"
	"lab5/syntax"
	"lab5/utils"
	"os"
)

func main() {
	// 1. 加载 DFA
	dfas, _ := lexer.LoadMultiDFAFromJson("./assets/all_dfa.json", "./dot", false)
	scanner := lexer.NewScanner()
	for _, d := range *dfas {
		scanner.RegisterDFA(d.DFA, d.TokenType)
	}

	// 2. 读取源代码
	codeBytes, _ := os.ReadFile("./assets/source.c")
	code := string(codeBytes)

	// 3. Tokenize → Symbol
	tokens := scanner.Tokenize(code)
	symbols := utils.TokensToSymbols(tokens)

	// 4. 构造文法
	g := syntax.DefineGrammar()

	// 5. 构造分析表
	dfa := parser.BuildDFA(g)
	follow := syntax.ComputeFollow(g)
	table := parser.BuildParseTable(g, dfa, follow)

	// 6. 调用 parser + 语义分析
	parser.Run(symbols, g, dfa, table, tokens)
}

// ---

// ## ✅ 说明

// * 所有语句（Stmt）统一入口；
// * 声明语句提取为 `Decl`，保留语义分类；
// * 可在 `semantic/actions.go` 中为 `Decl` 单独构建 AST 子树；
// * 便于支持如 `Decl → int id;`、`float id = 1.0;` 等扩展。

// ---

// 下一步你可以：

// 1. 用 `DefineGrammar()` 替代你原来手动添加产生式部分；
// 2. 在 `main.go` 中添加：`g := syntax.DefineGrammar()`；
// 3. 用你已有的 `tokens := Tokenize(...)` 流程执行测试。

// 是否需要我继续输出对应语义动作 `semantic/actions.go` 扩展？比如为这个文法中每一类句型构建 AST？
// 来你可以：

// 在 parser.Run() 中引用并调用 semantic.ActionFuncs[prodIndex](children)；

// 在最终 Accept 阶段输出 AST，可视化验证语义构造正确性；

// 使用 PrintAST 或我可以为你生成 .dot 导出图。

// 是否需要我帮你生成一个 ExportASTToDot() 函数来将 AST 渲染为 Graphviz 图？是否想先测试 int x = 3; 和 return x + 1; 等语句 AST 输出？我可以提供测试样例
