package parser

import (
	"fmt"
	"strings"
)

type Item struct {
	Left  string
	Right []string
	Dot   int // 点的位置
}

func (it Item) String() string {
	right := append([]string{}, it.Right...)
	if it.Dot < len(right) {
		right = append(right[:it.Dot], append([]string{"·"}, right[it.Dot:]...)...)
	} else {
		right = append(right, "·")
	}
	return fmt.Sprintf("%s → %s", it.Left, strings.Join(right, " "))
}

type State struct {
	Items []Item
	Index int
}

type ActionEntry struct {
	Action string // shift / reduce / accept
	Value  int    // 状态号 / 产生式编号
}

type ParsingTable struct {
	Action map[int]map[string]ActionEntry
	Goto   map[int]map[string]int
}
func GenerateSLRTable(grammar Grammar) ParsingTable {
	states := buildLR0States(grammar)
	first := computeFirstSets(grammar)
	follow := computeFollowSets(grammar, first)

	return buildParsingTable(grammar, states, follow)
}

func closure(items []Item, grammar Grammar) []Item {
	closureSet := append([]Item{}, items...)
	changed := true

	for changed {
		changed = false
		newItems := []Item{}

		for _, item := range closureSet {
			if item.Dot < len(item.Right) {
				B := item.Right[item.Dot]
				if grammar.NonTerminals[B] {
					for _, prod := range grammar.Productions {
						if prod.Left == B {
							newItem := Item{Left: B, Right: prod.Right, Dot: 0}
							if !itemInSet(newItem, closureSet) && !itemInSet(newItem, newItems) {
								newItems = append(newItems, newItem)
								changed = true
							}
						}
					}
				}
			}
		}

		closureSet = append(closureSet, newItems...)
	}
	return closureSet
}

func goTo(items []Item, symbol string, grammar Grammar) []Item {
	var moved []Item
	for _, item := range items {
		if item.Dot < len(item.Right) && item.Right[item.Dot] == symbol {
			moved = append(moved, Item{
				Left:  item.Left,
				Right: item.Right,
				Dot:   item.Dot + 1,
			})
		}
	}
	return closure(moved, grammar)
}

func buildLR0States(grammar Grammar) []State {
	startProd := Production{Left: "S'", Right: []string{grammar.StartSymbol}}
	grammar.Productions = append([]Production{startProd}, grammar.Productions...)

	I0 := closure([]Item{{Left: "S'", Right: []string{grammar.StartSymbol}, Dot: 0}}, grammar)

	states := []State{{Items: I0, Index: 0}}
	// transitions := map[string]map[int]int{}

	seen := []string{itemsToString(I0)}
	visited := 0

	for visited < len(states) {
		curr := states[visited]
		symbolSet := map[string]bool{}

		// 收集可能转移的符号
		for _, item := range curr.Items {
			if item.Dot < len(item.Right) {
				symbolSet[item.Right[item.Dot]] = true
			}
		}

		for sym := range symbolSet {
			nextItems := goTo(curr.Items, sym, grammar)
			if len(nextItems) == 0 {
				continue
			}
			key := itemsToString(nextItems)

			index := indexOfState(key, seen)
			if index == -1 {
				index = len(states)
				states = append(states, State{Items: nextItems, Index: index})
				seen = append(seen, key)
			}
		}
		visited++
	}
	return states
}

func itemInSet(item Item, set []Item) bool {
	for _, i := range set {
		if i.Left == item.Left && equal(i.Right, item.Right) && i.Dot == item.Dot {
			return true
		}
	}
	return false
}

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func itemsToString(items []Item) string {
	var s []string
	for _, i := range items {
		s = append(s, i.String())
	}
	return strings.Join(s, "|")
}

func indexOfState(key string, keys []string) int {
	for i, k := range keys {
		if k == key {
			return i
		}
	}
	return -1
}

func computeFirstSets(grammar Grammar) map[string]map[string]bool {
	first := make(map[string]map[string]bool)

	// 初始化
	for t := range grammar.Terminals {
		first[t] = map[string]bool{t: true}
	}
	for nt := range grammar.NonTerminals {
		first[nt] = map[string]bool{}
	}

	changed := true
	for changed {
		changed = false
		for _, prod := range grammar.Productions {
			A := prod.Left
			Xs := prod.Right
			for i := 0; i < len(Xs); i++ {
				B := Xs[i]
				for t := range first[B] {
					if !first[A][t] {
						first[A][t] = true
						changed = true
					}
				}
				// 如果 B 不能推出 ε，则中断
				if !first[B]["ε"] {
					break
				}
			}
			// 所有都能推出 ε，A 也能推出 ε
			allNullable := true
			for _, B := range Xs {
				if !first[B]["ε"] {
					allNullable = false
					break
				}
			}
			if allNullable {
				if !first[A]["ε"] {
					first[A]["ε"] = true
					changed = true
				}
			}
		}
	}

	return first
}

func computeFollowSets(grammar Grammar, first map[string]map[string]bool) map[string]map[string]bool {
	follow := make(map[string]map[string]bool)
	for nt := range grammar.NonTerminals {
		follow[nt] = map[string]bool{}
	}
	follow[grammar.StartSymbol]["$"] = true

	changed := true
	for changed {
		changed = false
		for _, prod := range grammar.Productions {
			A := prod.Left
			right := prod.Right
			for i := 0; i < len(right); i++ {
				B := right[i]
				if !grammar.NonTerminals[B] {
					continue
				}
				// 处理 First(β)
				firstBeta := map[string]bool{}
				allNullable := true
				for j := i + 1; j < len(right); j++ {
					sym := right[j]
					for t := range first[sym] {
						if t != "ε" {
							firstBeta[t] = true
						}
					}
					if !first[sym]["ε"] {
						allNullable = false
						break
					}
				}
				// Follow(B) += First(β)
				for t := range firstBeta {
					if !follow[B][t] {
						follow[B][t] = true
						changed = true
					}
				}
				// 若 β 全可空，Follow(B) += Follow(A)
				if i == len(right)-1 || allNullable {
					for t := range follow[A] {
						if !follow[B][t] {
							follow[B][t] = true
							changed = true
						}
					}
				}
			}
		}
	}
	return follow
}

func buildParsingTable(grammar Grammar, states []State, follow map[string]map[string]bool) ParsingTable {
	table := ParsingTable{
		Action: make(map[int]map[string]ActionEntry),
		Goto:   make(map[int]map[string]int),
	}

	// 为每个状态初始化表项
	for _, state := range states {
		table.Action[state.Index] = make(map[string]ActionEntry)
		table.Goto[state.Index] = make(map[string]int)

		for _, item := range state.Items {
			// 1. 移进项：A → α·aβ, a为终结符
			if item.Dot < len(item.Right) {
				a := item.Right[item.Dot]
				if grammar.Terminals[a] {
					next := goTo(state.Items, a, grammar)
					for _, s := range states {
						if itemsToString(s.Items) == itemsToString(next) {
							table.Action[state.Index][a] = ActionEntry{"shift", s.Index}
						}
					}
				}
			}

			// 2. 归约项：A → α·
			if item.Dot == len(item.Right) {
				if item.Left == "S'" {
					// 接受项
					table.Action[state.Index]["$"] = ActionEntry{"accept", 0}
				} else {
					// 查找产生式编号
					for i, prod := range grammar.Productions {
						if item.Left == prod.Left && equal(item.Right, prod.Right) {
							for a := range follow[item.Left] {
								table.Action[state.Index][a] = ActionEntry{"reduce", i}
							}
						}
					}
				}
			}
		}

		// GOTO 表：对每个非终结符
		symbolSet := make(map[string]bool)
		for _, item := range state.Items {
			if item.Dot < len(item.Right) {
				symbolSet[item.Right[item.Dot]] = true
			}
		}

		for sym := range symbolSet {
			if grammar.NonTerminals[sym] {
				next := goTo(state.Items, sym, grammar)
				for _, s := range states {
					if itemsToString(s.Items) == itemsToString(next) {
						table.Goto[state.Index][sym] = s.Index
					}
				}
			}
		}
	}

	return table
}

func PrintParsingTable(table ParsingTable) {
	fmt.Println("ACTION 表:")
	for state, row := range table.Action {
		for sym, entry := range row {
			fmt.Printf("  state %d, symbol %-4s → %-6s %d\n", state, sym, entry.Action, entry.Value)
		}
	}
	fmt.Println("GOTO 表:")
	for state, row := range table.Goto {
		for sym, dest := range row {
			fmt.Printf("  state %d, non-terminal %-8s → %d\n", state, sym, dest)
		}
	}
}
