package parser

import (
	"fmt"
	"lab5/syntax"
)

// DFA 状态结构
type State struct {
	Items []Item
	Index int
}

// DFA 转移边
type Edge struct {
	From   int
	Symbol syntax.Symbol
	To     int
}

// DFA 自动机
type DFA struct {
	States []State
	Edges  []Edge
}

// 构建 LR(0) DFA：规范集族
func BuildDFA(g *syntax.Grammar) *DFA {
	var dfa DFA

	// 初始状态：closure({S' → ·S})
	initial := Closure(g, []Item{{ProdIndex: 0, DotPos: 0}})
	state0 := State{Items: initial, Index: 0}
	dfa.States = append(dfa.States, state0)

	stateMap := map[string]int{
		hashItems(initial): 0,
	}

	queue := []State{state0}

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		symbols := collectSymbols(g, curr.Items)
		for sym := range symbols {
			gotoSet := Goto(g, curr.Items, sym)
			if len(gotoSet) == 0 {
				continue
			}
			hash := hashItems(gotoSet)
			if idx, ok := stateMap[hash]; ok {
				dfa.Edges = append(dfa.Edges, Edge{From: curr.Index, Symbol: sym, To: idx})
			} else {
				newState := State{Items: gotoSet, Index: len(dfa.States)}
				dfa.States = append(dfa.States, newState)
				stateMap[hash] = newState.Index
				queue = append(queue, newState)
				dfa.Edges = append(dfa.Edges, Edge{From: curr.Index, Symbol: sym, To: newState.Index})
			}
		}
	}

	return &dfa
}

// 计算点后所有出现过的符号（用于尝试 GOTO）
func collectSymbols(g *syntax.Grammar, items []Item) map[syntax.Symbol]bool {
	result := make(map[syntax.Symbol]bool)
	for _, it := range items {
		p := g.Productions[it.ProdIndex]
		if it.DotPos < len(p.Right) {
			s := p.Right[it.DotPos]
			result[s] = true
		}
	}
	return result
}

// 哈希项目集（字符串化后去重）
func hashItems(items []Item) string {
	result := ""
	for _, it := range items {
		result += fmt.Sprintf("[%d.%d]", it.ProdIndex, it.DotPos)
	}
	return result
}
