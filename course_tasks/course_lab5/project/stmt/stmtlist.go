package stmt

// ParseStmtList 拆分 { ... } 中多个语句（支持控制结构、嵌套函数）
func ParseStmtList(tokens []string) []string {
	if tokens[0] != "{" || tokens[len(tokens)-1] != "}" {
		panic("Block must start with '{' and end with '}'")
	}
	inner := tokens[1 : len(tokens)-1]

	code := []string{}
	start := 0
	for start < len(inner) {
		end := findStmtEnd(inner, start)
		stmtTokens := inner[start:end]

		if len(stmtTokens) > 3 && (stmtTokens[0] == "int" || stmtTokens[0] == "void") && stmtTokens[2] == "(" {
			stmtCode := GenerateFunctionDef(stmtTokens[1:]) // 递归函数定义
			code = append(code, stmtCode...)
		} else {
			stmtCode := Dispatch(stmtTokens)
			code = append(code, stmtCode...)
		}

		start = end
	}
	return code
}

// findStmtEnd 定位从 start 开始的完整语句结束位置（包含 ; 或结构块）
func findStmtEnd(tokens []string, start int) int {
	tok := tokens[start]

	// if (...) then ... else ...
	if tok == "if" {
		condEnd := FindCloseParen(tokens, start+1)
		thenStart := condEnd + 1
		var thenEnd int
		if tokens[thenStart] == "{" {
			thenEnd = FindCloseBrace(tokens, thenStart) + 1
		} else {
			thenEnd = findSemicolon(tokens, thenStart) + 1
		}

		// 检查是否有 else
		if thenEnd < len(tokens) && tokens[thenEnd] == "else" {
			elseStart := thenEnd + 1
			var elseEnd int
			if tokens[elseStart] == "{" {
				elseEnd = FindCloseBrace(tokens, elseStart) + 1
			} else {
				elseEnd = findSemicolon(tokens, elseStart) + 1
			}
			return elseEnd
		}
		return thenEnd
	}

	// while (...) body
	if tok == "while" {
		condEnd := FindCloseParen(tokens, start+1)
		bodyStart := condEnd + 1
		if tokens[bodyStart] == "{" {
			return FindCloseBrace(tokens, bodyStart) + 1
		}
		return findSemicolon(tokens, bodyStart) + 1
	}

	// 函数定义 int foo(...) 或 void bar(...)
	if (tok == "int" || tok == "void") && start+1 < len(tokens) {
		for i := start + 1; i < len(tokens); i++ {
			if tokens[i] == "(" {
				parenEnd := FindCloseParen(tokens, i)
				if parenEnd+1 < len(tokens) && tokens[parenEnd+1] == "{" {
					braceEnd := FindCloseBrace(tokens, parenEnd+1)
					return braceEnd + 1
				}
				break
			}
		}
	}

	// 普通语句
	return findSemicolon(tokens, start) + 1
}

// findSemicolon 定位从 start 开始的第一个分号
func findSemicolon(tokens []string, start int) int {
	for i := start; i < len(tokens); i++ {
		if tokens[i] == ";" {
			return i
		}
	}
	panic("Missing semicolon")
}
