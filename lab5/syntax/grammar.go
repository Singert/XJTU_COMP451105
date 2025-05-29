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
	case "id", "num", "float", "double", "string", "char", "type_kw", "return", "if", "for", "else", "while",
		"=", "+", "-", "*", "/", "==", "<", "!", "&&", "||", ">", "!=", ">=", "<=",
		"(", ")", "{", "}", ";", ",", "[", "]":
		return true
	default:
		return false
	}
}

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
	g.AddProduction("Factor", []Symbol{"float"})
	// 文法 32
	g.AddProduction("Factor", []Symbol{"char"})
	// 文法 33
	g.AddProduction("Factor", []Symbol{"string"})
	// 文法 34
	g.AddProduction("Factor", []Symbol{"id"})
	// 文法 35
	g.AddProduction("Factor", []Symbol{"(", "Expr", ")"})
	// 文法 36
	g.AddProduction("Factor", []Symbol{"id", "MultiIndex"})

	// ==== 函数参数列表 ====
	// 文法 37
	g.AddProduction("Args", []Symbol{"NonEmptyArgs"})
	// 文法 38
	g.AddProduction("Args", []Symbol{})

	// 文法 39
	g.AddProduction("NonEmptyArgs", []Symbol{"Expr"})
	// 文法 40
	g.AddProduction("NonEmptyArgs", []Symbol{"NonEmptyArgs", ",", "Expr"})
	// 文法 41
	g.AddProduction("NonEmptyArgs", []Symbol{"Type", "id"})
	// 文法 42
	g.AddProduction("NonEmptyArgs", []Symbol{"Type", "id", "=", "Expr"})
	// 文法 43
	g.AddProduction("NonEmptyArgs", []Symbol{"NonEmptyArgs", ",", "Type", "id"})
	// 文法 44
	g.AddProduction("NonEmptyArgs", []Symbol{"NonEmptyArgs", ",", "Type", "id", "=", "Expr"})

	// ==== 多维数组索引 ====
	// 文法 45
	g.AddProduction("IndexList", []Symbol{"Expr"})
	// 文法 46
	g.AddProduction("IndexList", []Symbol{"IndexList", ",", "Expr"})

	// ==== 条件表达式 ====
	// 文法 47
	g.AddProduction("Cond", []Symbol{"Cond", "&&", "Cond"})
	// 文法 48
	g.AddProduction("Cond", []Symbol{"Cond", "||", "Cond"})
	// 文法 49
	g.AddProduction("Cond", []Symbol{"!", "Cond"})
	// 文法 50
	g.AddProduction("Cond", []Symbol{"Expr", "<", "Expr"})
	// 文法 51
	g.AddProduction("Cond", []Symbol{"Expr", ">", "Expr"})
	// 文法 52
	g.AddProduction("Cond", []Symbol{"Expr", "<=", "Expr"})
	// 文法 53
	g.AddProduction("Cond", []Symbol{"Expr", ">=", "Expr"})
	// 文法 54
	g.AddProduction("Cond", []Symbol{"Expr", "!=", "Expr"})
	// 文法 55
	g.AddProduction("Cond", []Symbol{"Expr", "==", "Expr"})
	// 文法 56
	g.AddProduction("Cond", []Symbol{"(", "Cond", ")"})
	// 文法 57
	g.AddProduction("Cond", []Symbol{"Expr"})

	// ==== 数组声明 ====
	// 文法 58
	g.AddProduction("Decl", []Symbol{"Type", "id", "MultiIndex", ";"})
	// 文法 59
	g.AddProduction("Decl", []Symbol{"Type", "id", "MultiIndex", "=", "Expr", ";"})
	// 文法 60
	g.AddProduction("MultiIndex", []Symbol{"[", "IndexList", "]", "MultiIndex"})
	// 文法 61
	g.AddProduction("MultiIndex", []Symbol{})

	// 文法 62
	g.AddProduction("Decl", []Symbol{"Type", "id", "MultiIndex", "=", "InitList", ";"})

	// 文法 63
	g.AddProduction("InitList", []Symbol{"{", "NonEmptyInitList", "}"})
	// 文法 64
	g.AddProduction("InitList", []Symbol{"{", "}"})

	// 文法 65
	g.AddProduction("NonEmptyInitList", []Symbol{"Expr"})
	// 文法 66
	g.AddProduction("NonEmptyInitList", []Symbol{"NonEmptyInitList", ",", "Expr"})

	// 文法 67
	g.AddProduction("Expr", []Symbol{"InitList"})

	// ==== 一元负号 ====
	// 文法 68
	g.AddProduction("Factor", []Symbol{"-", "Factor"})

	// ==== 支持带数组下标的形参声明 ====
	// 文法 69
	g.AddProduction("NonEmptyArgs", []Symbol{"Type", "id", "MultiIndex"})
	// 文法 70
	g.AddProduction("NonEmptyArgs", []Symbol{"Type", "id", "MultiIndex", "=", "Expr"})

	// ==== 支持空的IndexList ====
	// 文法 71
	g.AddProduction("IndexList", []Symbol{})

	// ==== 支持 for 循环 ====
	// 文法 72
	g.AddProduction("Stmt", []Symbol{"for", "(", "ForInit", ";", "Cond", ";", "Expr", ")", "Stmt"})
	// 文法 73
	g.AddProduction("ForInit", []Symbol{"Decl"})
	// 文法 74
	g.AddProduction("ForInit", []Symbol{"Expr"})
	// 文法 75
	g.AddProduction("ForInit", []Symbol{})

	// ==== 赋值表达式支持 ====
	// 文法 76
	g.AddProduction("Expr", []Symbol{"id", "=", "Expr"})

	return g
}
