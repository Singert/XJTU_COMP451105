// ===== analyzer/ll1check.go =====
package analyzer

import (
	"fmt"
	"ll1-analyzer/grammar"
)

func CheckLL1(g *grammar.Grammar, first, follow map[string]map[string]bool) []string {
	violations := []string{}
	rules := map[string][][]string{}
	for _, p := range g.Productions {
		lhs := p.Left.Name
		rhsStr := []string{}
		for _, sym := range p.Right {
			rhsStr = append(rhsStr, sym.Name)
		}
		rules[lhs] = append(rules[lhs], rhsStr)
	}

	for lhs, rhss := range rules {
		tests := make([]map[string]bool, len(rhss))
		nullable := make([]bool, len(rhss))
		for i, rhs := range rhss {
			tests[i] = map[string]bool{}
			nullable[i] = true
			for _, s := range rhs {
				for tok := range first[s] {
					if tok != "ε" {
						tests[i][tok] = true
					}
				}
				if !first[s]["ε"] {
					nullable[i] = false
					break
				}
			}
			if nullable[i] {
				for tok := range follow[lhs] {
					tests[i][tok] = true
				}
			}
		}
		// pairwise check
		for i := 0; i < len(rhss); i++ {
			for j := i + 1; j < len(rhss); j++ {
				for tok := range tests[i] {
					if tests[j][tok] {
						violations = append(violations, fmt.Sprintf("Conflict at %s on token '%s' between %v and %v", lhs, tok, rhss[i], rhss[j]))
						break
					}
				}
			}
		}
	}
	return violations
}
