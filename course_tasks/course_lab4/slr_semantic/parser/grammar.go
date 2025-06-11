package parser

type Production struct {
	Left  string   // 左部非终结符
	Right []string // 右部符号序列
}

type Grammar struct {
	Terminals    map[string]bool
	NonTerminals map[string]bool
	Productions  []Production
	StartSymbol  string
}

// 构造文法对象（硬编码版本）
func LoadGrammar() Grammar {
	grammar := Grammar{
		Terminals:    make(map[string]bool),
		NonTerminals: make(map[string]bool),
		StartSymbol:  "P",
	}

	// 定义产生式
	add := func(lhs string, rhs ...string) {
		grammar.Productions = append(grammar.Productions, Production{
			Left:  lhs,
			Right: rhs,
		})
		grammar.NonTerminals[lhs] = true
		for _, sym := range rhs {
			if !isNonTerminal(sym) && sym != "ε" {
				grammar.Terminals[sym] = true
			}
		}
	}

	// 文法规则（简化自图示）
	add("P", "DeclList")
	add("DeclList", "DeclList", "Decl")
	add("DeclList", "Decl")
	add("Decl", "T", "id", ";")
	add("Decl", "T", "id", "[", "num", "]", ";")
	add("Decl", "T", "id", "(", "ParamList", ")", ";")
	add("T", "int")
	add("T", "void")
	add("ParamList", "ParamList", ",", "Param")
	add("ParamList", "Param")
	add("Param", "T", "id")

	return grammar
}

func isNonTerminal(s string) bool {
	return s == "P" || s == "DeclList" || s == "Decl" || s == "T" || s == "ParamList" || s == "Param"
}
