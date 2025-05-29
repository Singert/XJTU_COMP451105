package syntax

// Symbol represents a terminal or non-terminal symbol in a grammar.
type Symbol string

// Production represents a production rule.
type Production struct {
	Left  Symbol
	Right []Symbol
}

// Grammar holds all productions, terminal/non-terminal sets, and start symbol.
type Grammar struct {
	Productions []Production
	StartSymbol Symbol
	Terminals   map[Symbol]bool
	NonTerms    map[Symbol]bool
}

// NewGrammar creates a new Grammar instance.
func NewGrammar(start Symbol) *Grammar {
	return &Grammar{
		StartSymbol: start,
		Terminals:   make(map[Symbol]bool),
		NonTerms:    make(map[Symbol]bool),
	}
}

// AddProduction adds a production rule to the grammar and classifies its symbols.
func (g *Grammar) AddProduction(left Symbol, right []Symbol) {
	g.Productions = append(g.Productions, Production{Left: left, Right: right})
	g.NonTerms[left] = true
	for _, symb := range right {
		if isTerminal(symb) {
			g.Terminals[symb] = true
		} else {
			g.NonTerms[symb] = true
		}
	}
}

// isTerminal checks if a symbol is a terminal based on known literals/keywords.
func isTerminal(symb Symbol) bool {
	switch symb {
	case "id", "num", "type_kw", "return", "if", "else", "while",
		"=", "+", "-", "*", "/", "==", "<", "!", "&&", "||", ">", "!=", ">=", "<=",
		"(", ")", "{", "}", ";", ",", "[", "]":
		return true
	default:
		return false
	}
}

// DefineGrammar defines the grammar and all its productions.
func DefineGrammar() *Grammar {
	g := NewGrammar("S'")

	// 文法 0
	g.AddProduction("S'", []Symbol{"Program"})

	// ==== 程序入口 ====
	// 文法 1
	g.AddProduction("Program", []Symbol{"StmtList"})
	// 文法 2
	g.AddProduction("Program", []Symbol{"StmtList", "EOF"})

	// ==== 函数定义作为语句支持 ====
	// 文法 3
	g.AddProduction("Stmt", []Symbol{"Func"})

	// ==== 函数定义 ====
	// 文法 4
	g.AddProduction("Func", []Symbol{"Type", "id", "(", "Args", ")", "Block"})

	// ==== 声明语句 ====
	// 文法 5
	g.AddProduction("Stmt", []Symbol{"Decl"})
	// 文法 6
	g.AddProduction("Decl", []Symbol{"Type", "id", "=", "Expr", ";"})

	// ==== 类型定义 ====
	// 文法 7
	g.AddProduction("Type", []Symbol{"type_kw"})

	// ==== 普通语句 ====
	// 文法 8
	g.AddProduction("Stmt", []Symbol{"id", "=", "Expr", ";"})
	// 文法 9
	g.AddProduction("Stmt", []Symbol{"return", "Expr", ";"})
	// 文法 10
	g.AddProduction("Stmt", []Symbol{"Block"})
	// 文法 11
	g.AddProduction("Stmt", []Symbol{"if", "(", "Cond", ")", "Stmt"})
	// 文法 12
	g.AddProduction("Stmt", []Symbol{"if", "(", "Cond", ")", "Stmt", "else", "Stmt"})
	// 文法 13
	g.AddProduction("Stmt", []Symbol{"while", "(", "Cond", ")", "Stmt"})
	// 文法 14
	g.AddProduction("Stmt", []Symbol{"id", "MultiIndex", "=", "Expr", ";"})
	// 文法 15
	g.AddProduction("Stmt", []Symbol{"id", "(", "Args", ")", ";"})

	// ==== 块与语句序列 ====
	// 文法 16
	g.AddProduction("Block", []Symbol{"{", "StmtList", "}"})
	// 文法 17
	g.AddProduction("StmtList", []Symbol{})
	// 文法 18
	g.AddProduction("StmtList", []Symbol{"Stmt"})
	// 文法 19
	g.AddProduction("StmtList", []Symbol{"StmtList", "Stmt"})

	// ==== 表达式结构 ====
	// 文法 20
	g.AddProduction("Expr", []Symbol{"Expr", "+", "Term"})
	// 文法 21
	g.AddProduction("Expr", []Symbol{"Expr", "-", "Term"})
	// 文法 22
	g.AddProduction("Expr", []Symbol{"Term"})
	// 文法 23
	g.AddProduction("Term", []Symbol{"Term", "*", "CastExpr"})
	// 文法 24
	g.AddProduction("Term", []Symbol{"Term", "/", "CastExpr"})
	// 文法 25
	g.AddProduction("Term", []Symbol{"CastExpr"})

	// ==== 强制类型转换 ====
	// 文法 26
	g.AddProduction("CastExpr", []Symbol{"CastPrefix", "Factor"})
	// 文法 27
	g.AddProduction("CastExpr", []Symbol{"Factor"})
	// 文法 28
	g.AddProduction("CastPrefix", []Symbol{"(", "Type", ")"})

	// ==== 基本因子 + 函数调用 ====
	// 文法 29
	g.AddProduction("Factor", []Symbol{"id", "(", "Args", ")"})
	// 文法 30
	g.AddProduction("Factor", []Symbol{"num"})
	// 文法 31
	g.AddProduction("Factor", []Symbol{"id"})
	// 文法 32
	g.AddProduction("Factor", []Symbol{"(", "Expr", ")"})
	// 文法 33
	g.AddProduction("Factor", []Symbol{"id",  "MultiIndex"})

	// ==== 函数参数列表 ====
	// 文法 34
	g.AddProduction("Args", []Symbol{"NonEmptyArgs"})
	// 文法 35
	g.AddProduction("Args", []Symbol{})
	// 文法 36
	g.AddProduction("NonEmptyArgs", []Symbol{"Expr"})
	// 文法 37
	g.AddProduction("NonEmptyArgs", []Symbol{"NonEmptyArgs", ",", "Expr"})
	// 文法 38
	g.AddProduction("NonEmptyArgs", []Symbol{"Type", "id"})
	// 文法 39
	g.AddProduction("NonEmptyArgs", []Symbol{"Type", "id", "=", "Expr"})
	// 文法 40
	g.AddProduction("NonEmptyArgs", []Symbol{"NonEmptyArgs", ",", "Type", "id"})
	// 文法 41
	g.AddProduction("NonEmptyArgs", []Symbol{"NonEmptyArgs", ",", "Type", "id", "=", "Expr"})

	// ==== 多维数组索引 ====
	// 文法 42
	g.AddProduction("IndexList", []Symbol{"Expr"})
	// 文法 43
	g.AddProduction("IndexList", []Symbol{"IndexList", ",", "Expr"})

	// ==== 条件表达式 ====
	// 文法 44
	g.AddProduction("Cond", []Symbol{"Cond", "&&", "Cond"})
	// 文法 45
	g.AddProduction("Cond", []Symbol{"Cond", "||", "Cond"})
	// 文法 46
	g.AddProduction("Cond", []Symbol{"!", "Cond"})
	// 文法 47
	g.AddProduction("Cond", []Symbol{"Expr", "<", "Expr"})
	// 文法 48
	g.AddProduction("Cond", []Symbol{"Expr", ">", "Expr"})
	// 文法 49
	g.AddProduction("Cond", []Symbol{"Expr", "<=", "Expr"})
	// 文法 50
	g.AddProduction("Cond", []Symbol{"Expr", ">=", "Expr"})
	// 文法 51
	g.AddProduction("Cond", []Symbol{"Expr", "!=", "Expr"})
	// 文法 52
	g.AddProduction("Cond", []Symbol{"Expr", "==", "Expr"})
	// 文法 53
	g.AddProduction("Cond", []Symbol{"(", "Cond", ")"})
	// 文法 54
	g.AddProduction("Cond", []Symbol{"Expr"})

	// 额外添加数组声明产生式
	// // 文法 55
	// g.AddProduction("Decl", []Symbol{"Type", "id", "[", "IndexList", "]", ";"})
	// // 文法 56
	// g.AddProduction("Decl", []Symbol{"Type", "id", "[", "IndexList", "]", "=", "Expr", ";"})
	g.AddProduction("Decl", []Symbol{"Type", "id", "MultiIndex", ";"})
	g.AddProduction("Decl", []Symbol{"Type", "id", "MultiIndex", "=", "Expr", ";"})

	g.AddProduction("MultiIndex", []Symbol{"[", "IndexList", "]","MultiIndex"})
	g.AddProduction("MultiIndex", []Symbol{}) // 终止符号

	return g
}
