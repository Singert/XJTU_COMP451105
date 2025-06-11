// ===== analyzer/follow.go =====
package analyzer

import (
	"ll1-analyzer/grammar"
)

func ComputeFollowSets(g *grammar.Grammar, first map[string]map[string]bool) map[string]map[string]bool {
	follow := map[string]map[string]bool{}
	for nt := range g.NonTerminals {
		follow[nt] = map[string]bool{}
	}
	follow[g.Start]["#"] = true

	changed := true
	for changed {
		changed = false
		for _, prod := range g.Productions {
			lhs := prod.Left.Name
			rhs := prod.Right
			for i := 0; i < len(rhs); i++ {
				sym := rhs[i]
				if sym.IsNonTerminal() {
					trailer := map[string]bool{}
					j := i + 1
					for ; j < len(rhs); j++ {
						next := rhs[j]
						for tok := range first[next.Name] {
							if tok != "ε" {
								trailer[tok] = true
							}
						}
						if !first[next.Name]["ε"] {
							break
						}
					}
					if j == len(rhs) || allNullable(rhs[i+1:], first) {
						for tok := range follow[lhs] {
							trailer[tok] = true
						}
					}
					for tok := range trailer {
						if !follow[sym.Name][tok] {
							follow[sym.Name][tok] = true
							changed = true
						}
					}
				}
			}
		}
	}
	return follow
}

func allNullable(seq []grammar.Symbol, first map[string]map[string]bool) bool {
	for _, sym := range seq {
		if !first[sym.Name]["ε"] {
			return false
		}
	}
	return true
}
