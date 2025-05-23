package stmt

// 找到匹配的右括号位置（tokens[start] == "("）
func findCloseParen(tokens []string, start int) int {
	level := 0
	for i := start; i < len(tokens); i++ {
		if tokens[i] == "(" {
			level++
		} else if tokens[i] == ")" {
			level--
			if level == 0 {
				return i
			}
		}
	}
	panic("Unmatched parenthesis")
}
