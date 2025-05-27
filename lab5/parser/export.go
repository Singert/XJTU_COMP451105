package parser

import (
	"fmt"
	"lab5/syntax"
	"os"
)

// 导出 DFA 为 DOT 文件
func ExportDFAtoDOT(dfa *DFA, g *syntax.Grammar, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Fprintln(file, "digraph DFA {")
	fmt.Fprintln(file, `  rankdir=LR;`)
	fmt.Fprintln(file, `  node [shape=box, fontname="monospace"];`)

	for _, state := range dfa.States {
		label := fmt.Sprintf("I%d\\n", state.Index)
		for _, it := range state.Items {
			label += escape(it.String(g)) + "\\l"
		}
		fmt.Fprintf(file, "  I%d [label=\"%s\"];\n", state.Index, label)
	}

	for _, edge := range dfa.Edges {
		fmt.Fprintf(file, "  I%d -> I%d [label=\"%s\"];\n", edge.From, edge.To, edge.Symbol)
	}

	fmt.Fprintln(file, "}")
	return nil
}

// PrintParseTable 将 ACTION 和 GOTO 表输出到终端
func PrintParseTable(table *ParseTable, g *syntax.Grammar) {
	fmt.Println("======= ACTION 表 =======")
	fmt.Printf("%-8s", "State")
	for t := range g.Terminals {
		fmt.Printf("%-8s", t)
	}
	fmt.Printf("%-8s\n", "#")

	for state, acts := range table.Action {
		fmt.Printf("%-8d", state)
		for t := range g.Terminals {
			fmt.Printf("%-8s", acts[t].String(g))
		}
		fmt.Printf("%-8s\n", acts["#"].String(g))
	}

	fmt.Println("\n======= GOTO 表 =======")
	fmt.Printf("%-8s", "State")
	for nt := range g.NonTerms {
		fmt.Printf("%-8s", nt)
	}
	fmt.Println()
	for state, gotos := range table.Goto {
		fmt.Printf("%-8d", state)
		for nt := range g.NonTerms {
			if to, ok := gotos[nt]; ok {
				fmt.Printf("%-8d", to)
			} else {
				fmt.Printf("%-8s", "")
			}
		}
		fmt.Println()
	}
}

// ExportParseTableDOT 导出分析表为 DOT 图
func ExportParseTableDOT(table *ParseTable, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Fprintln(file, "digraph ParseTable {")
	fmt.Fprintln(file, `  rankdir=LR;`)
	fmt.Fprintln(file, `  node [shape=ellipse];`)

	for from, actions := range table.Action {
		for symbol, act := range actions {
			label := act.String(nil)
			if label != "" {
				to := act.Value
				fmt.Fprintf(file, `  S%d -> S%d [label="%s/%s"];`+"\n", from, to, symbol, label)
			}
		}
	}
	for from, gotos := range table.Goto {
		for nt, to := range gotos {
			fmt.Fprintf(file, `  S%d -> S%d [label="%s"];`+"\n", from, to, nt)
		}
	}

	fmt.Fprintln(file, "}")
	return nil
}

// 转义 DOT 格式用的字符串
func escape(s string) string {
	return s // 可添加更复杂的清理逻辑
}
