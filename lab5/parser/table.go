package parser

import (
	"fmt"
	"lab5/syntax"
)

type ActionType int

const (
	Shift ActionType = iota
	Reduce
	Accept
	Error
)

type Action struct {
	Typ   ActionType
	Value int // Shift: 目标状态，Reduce: 产生式编号
}

// ACTION[state][terminal] = Action
// GOTO[state][nonterminal] = int
type ParseTable struct {
	Action map[int]map[syntax.Symbol]Action
	Goto   map[int]map[syntax.Symbol]int
}

// 构造 ACTION 和 GOTO 表(支持follow set)
func BuildParseTable(g *syntax.Grammar, dfa *DFA, follow map[syntax.Symbol]map[syntax.Symbol]bool) *ParseTable {
	table := &ParseTable{
		Action: make(map[int]map[syntax.Symbol]Action),
		Goto:   make(map[int]map[syntax.Symbol]int),
	}

	for _, state := range dfa.States {
		table.Action[state.Index] = make(map[syntax.Symbol]Action)
		table.Goto[state.Index] = make(map[syntax.Symbol]int)

		for _, it := range state.Items {
			p := g.Productions[it.ProdIndex]

			// 接受项：S' → S ·
			if it.ProdIndex == 0 && it.DotPos == len(p.Right) {
				table.Action[state.Index]["#"] = Action{Typ: Accept}
				continue
			}

			// 归约项：A → α ·
			if it.DotPos == len(p.Right) {
				lhs := p.Left
				for a := range follow[lhs] {
					table.Action[state.Index][a] = Action{Typ: Reduce, Value: it.ProdIndex}
				}
			}
		}
	}

	// 处理转移边：移进或 GOTO
	for _, e := range dfa.Edges {
		if g.Terminals[e.Symbol] || e.Symbol == "#" {
			table.Action[e.From][e.Symbol] = Action{Typ: Shift, Value: e.To}
		} else if g.NonTerms[e.Symbol] {
			table.Goto[e.From][e.Symbol] = e.To
		}
	}

	return table
}

// 可打印表示
func (a Action) String(g *syntax.Grammar) string {
	switch a.Typ {
	case Shift:
		return fmt.Sprintf("s%d", a.Value)
	case Reduce:
		return fmt.Sprintf("r%d", a.Value)
	case Accept:
		return "acc"
	default:
		return ""
	}
}
