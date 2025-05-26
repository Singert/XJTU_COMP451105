package main

import (
	"fmt"
	"lab4/driver"
	"lab4/grammar"
	"lab4/parser"
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
	// 构建 FOLLOW 集
	follow := grammar.ComputeFollow(g)
	// === Step 3: 构建 ACTION / GOTO 表 ===
	table := parser.BuildParseTable(g, dfa, follow)
	// 打印 ACTION 和 GOTO 表
	parser.PrintParseTable(table, g)
	err := parser.ExportParseTableDOT(table, "parse_table.dot")
	if err != nil {
		fmt.Println("导出 parse_table.dot 失败:", err)
	} else {
		fmt.Println("✔ ACTION/GOTO 表已导出为 parse_table.dot（可用 dot -Tpng 查看）")
	}
	// === Step 4: 输入 token 流（来自词法分析器或手动输入）===
	input := []grammar.Symbol{"id", "+", "id", "*", "id"}

	fmt.Println("\n======= 分析过程 =======")
	driver.Run(input, g, dfa, table)
	parser.ExportDFAtoDOT(dfa, g, "dfa.dot")
	fmt.Println("\nDFA 已导出到 dfa.dot 文件。请使用 Graphviz 工具查看。")
}
