package stmt

// 将 { ... } 中的语句块分解为多个 S 并递归调用
func ParseStmtList(tokens []string) []string {
	code := []string{}

	// 去除首尾大括号
	if tokens[0] != "{" || tokens[len(tokens)-1] != "}" {
		panic("Block must start with '{' and end with '}'")
	}
	inner := tokens[1 : len(tokens)-1]

	// 拆分成若干语句（以分号为界）
	start := 0
	for i := 0; i < len(inner); i++ {
		if inner[i] == ";" {
			stmtTokens := inner[start : i+1]
			stmtCode := Dispatch(stmtTokens)
			code = append(code, stmtCode...)
			start = i + 1
		}
	}
	return code
}
