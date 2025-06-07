package grammar

// ComputeFirst 计算每个符号的 FIRST 集
func ComputeFirst(g *Grammar) map[Symbol]map[Symbol]bool {
	first := make(map[Symbol]map[Symbol]bool)

	// 初始化所有符号的 FIRST 集
	for t := range g.Terminals {
		first[t] = map[Symbol]bool{t: true} // 终结符的 FIRST 集是其本身
	}
	for nt := range g.NonTerms {
		first[nt] = make(map[Symbol]bool)
	}

	changed := true
	for changed {
		changed = false
		for _, prod := range g.Productions {
			lhs := prod.Left
			rhs := prod.Right

			nullable := true
			for _, sym := range rhs {
				for f := range first[sym] {
					if !first[lhs][f] {
						first[lhs][f] = true
						changed = true
					}
				}
				if !first[sym]["ε"] {
					nullable = false
					break
				}
			}
			if nullable {
				if !first[lhs]["ε"] {
					first[lhs]["ε"] = true
					changed = true
				}
			}
		}
	}
	return first
}

// ComputeFollow 计算每个非终结符的 FOLLOW 集，符合教材规则
func ComputeFollow(g *Grammar) map[Symbol]map[Symbol]bool {
	first := ComputeFirst(g)
	follow := make(map[Symbol]map[Symbol]bool)

	// 初始化 FOLLOW 集为空集合
	for nt := range g.NonTerms {
		follow[nt] = make(map[Symbol]bool)
	}
	// 起始符号的 FOLLOW 集包含 #（输入终止符）
	follow[g.StartSymbol]["#"] = true

	changed := true
	for changed {
		changed = false

		for _, prod := range g.Productions {
			lhs := prod.Left
			rhs := prod.Right
			for i := 0; i < len(rhs); i++ {
				symb := rhs[i]
				if !g.NonTerms[symb] {
					continue
				}

				// 计算 β = rhs[i+1:]
				beta := rhs[i+1:]
				firstBeta := make(map[Symbol]bool)
				nullable := true
				for _, b := range beta {
					for f := range first[b] {
						if f != "ε" {
							firstBeta[f] = true
						}
					}
					if !first[b]["ε"] {
						nullable = false
						break
					}
				}
				// R2: FOLLOW(B) ⊇ FIRST(β) \ {ε}
				for f := range firstBeta {
					if !follow[symb][f] {
						follow[symb][f] = true
						changed = true
					}
				}
				// R1: 若 β ⇒* ε，则 FOLLOW(B) ⊇ FOLLOW(A)
				if len(beta) == 0 || nullable {
					for f := range follow[lhs] {
						if !follow[symb][f] {
							follow[symb][f] = true
							changed = true
						}
					}
				}
			}
		}
	}
	return follow
}
