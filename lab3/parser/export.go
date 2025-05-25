package parser

import (
	"fmt"
	"lab3/grammar"
	"os"
)

// 导出 DFA 为 DOT 文件
func ExportDFAtoDOT(dfa *DFA, g *grammar.Grammar, filename string) error {
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

// 转义 DOT 格式用的字符串
func escape(s string) string {
	return s // 可添加更复杂的清理逻辑
}
