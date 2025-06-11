// ===== analyzer/first.go =====
package analyzer

import (
	"ll1-analyzer/grammar"
)

func ComputeFirstSets(g *grammar.Grammar) map[string]map[string]bool {
	first := map[string]map[string]bool{}

	for nt := range g.NonTerminals {
		first[nt] = map[string]bool{}
	}
	for t := range g.Terminals {
		first[t] = map[string]bool{t: true}
	}
	first["ε"] = map[string]bool{"ε": true}

	changed := true
	for changed {
		changed = false
		for _, prod := range g.Productions {
			A := prod.Left.Name
			rhs := prod.Right
			nullable := true
			for _, sym := range rhs {
				for tok := range first[sym.Name] {
					if tok != "ε" && !first[A][tok] {
						first[A][tok] = true
						changed = true
					}
				}
				if !first[sym.Name]["ε"] {
					nullable = false
					break
				}
			}
			if nullable && !first[A]["ε"] {
				first[A]["ε"] = true
				changed = true
			}
		}
	}
	return first
}
