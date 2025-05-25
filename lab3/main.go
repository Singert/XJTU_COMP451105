package main

import (
	"fmt"
	"lab3/driver"
	"lab3/grammar"
	"lab3/parser"
)

func main() {
	// === Step 1: 构造文法 ===
	g := grammar.NewGrammar("S'")
	g.AddProduction("S'", []grammar.Symbol{"E"}) // 增广产生式
	g.AddProduction("E", []grammar.Symbol{"E", "+", "T"})
	g.AddProduction("E", []grammar.Symbol{"T"})
	g.AddProduction("T", []grammar.Symbol{"T", "*", "F"})
	g.AddProduction("T", []grammar.Symbol{"F"})
	g.AddProduction("F", []grammar.Symbol{"(", "E", ")"})
	g.AddProduction("F", []grammar.Symbol{"id"})

	// === Step 2: 构建 LR(0) 项目集 DFA ===
	dfa := parser.BuildDFA(g)

	// 打印状态集（可选）
	fmt.Println("======= LR(0) 项目集规范族（状态集） =======")
	for _, state := range dfa.States {
		fmt.Printf("状态 I%d:\n", state.Index)
		for _, it := range state.Items {
			fmt.Println("  ", it.String(g))
		}
	}

	// === Step 3: 构建 ACTION / GOTO 表 ===
	table := parser.BuildParseTable(g, dfa)

	// === Step 4: 输入 token 流（来自词法分析器或手动输入）===
	input := []grammar.Symbol{"id", "+", "id", "*", "id"}

	fmt.Println("\n======= 分析过程 =======")
	driver.Run(input, g, dfa, table)
	parser.ExportDFAtoDOT(dfa, g, "dfa.dot")
	fmt.Println("\nDFA 已导出到 dfa.dot 文件。请使用 Graphviz 工具查看。")
}
