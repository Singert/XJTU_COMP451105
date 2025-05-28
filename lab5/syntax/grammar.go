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
	case "id", "num", "int", "return", "if", "else", "while",
		"=", "+", "-", "*", "/", "==", "<", "!", "&&", "||",
		"(", ")", "{", "}", ";", ",", "[", "]":
		return true
	default:
		return false
	}
}

// DefineGrammar defines the grammar and all its productions.
func DefineGrammar() *Grammar {
	g := NewGrammar("S'")
	//0. 定义起始符号
	g.AddProduction("S'", []Symbol{"Stmt"})

	// ==== 声明语句 ====

	//1. 声明语句提取为 Decl，保留语义分类
	g.AddProduction("Stmt", []Symbol{"Decl"})
	//2. Decl 语句定义，支持 int、float 等类型
	g.AddProduction("Decl", []Symbol{"int", "id", "=", "num", ";"}) // 可扩展更多类型如 float 等

	// ==== 普通语句 ====

	//3. 变量赋值
	g.AddProduction("Stmt", []Symbol{"id", "=", "Expr", ";"})
	//4. return 语句
	g.AddProduction("Stmt", []Symbol{"return", "Expr", ";"})
	//5. 块语句
	g.AddProduction("Stmt", []Symbol{"Block"})
	//6. if 语句
	g.AddProduction("Stmt", []Symbol{"if", "(", "Cond", ")", "Stmt"})
	//7. if-else 语句
	g.AddProduction("Stmt", []Symbol{"if", "(", "Cond", ")", "Stmt", "else", "Stmt"})
	//8. while 语句
	g.AddProduction("Stmt", []Symbol{"while", "(", "Cond", ")", "Stmt"})
	//9. 数组赋值
	g.AddProduction("Stmt", []Symbol{"id", "[", "IndexList", "]", "=", "Expr", ";"})
	//10. 函数调用语句
	g.AddProduction("Stmt", []Symbol{"id", "(", "Args", ")", ";"})

	// === 块与语句序列 ===

	//11. 块语句定义
	g.AddProduction("Block", []Symbol{"{", "StmtList", "}"})
	//12. 语句列表定义
	g.AddProduction("StmtList", []Symbol{"Stmt"})
	//13. 语句列表递归定义
	g.AddProduction("StmtList", []Symbol{"StmtList", "Stmt"})

	// === 表达式结构 ===
	//14. 表达式加法
	g.AddProduction("Expr", []Symbol{"Expr", "+", "Term"})
	//15. 表达式减法
	g.AddProduction("Expr", []Symbol{"Expr", "-", "Term"})
	//16. 表达式单一项
	g.AddProduction("Expr", []Symbol{"Term"})
	//17. 乘法项
	g.AddProduction("Term", []Symbol{"Term", "*", "Factor"})
	//18. 除法项
	g.AddProduction("Term", []Symbol{"Term", "/", "Factor"})
	//19. 基本因子
	g.AddProduction("Term", []Symbol{"Factor"})

	// === 基本因子 + 函数调用 ===
	
	//20. 函数调用
	g.AddProduction("Factor", []Symbol{"id", "(", "Args", ")"}) // 函数调用
	//21. 数字因子
	g.AddProduction("Factor", []Symbol{"num"})
	//22. 标识符因子
	g.AddProduction("Factor", []Symbol{"id"})
	//23. 括号表达式
	g.AddProduction("Factor", []Symbol{"(", "Expr", ")"})

	// === 函数参数列表 ===
	//24. 函数参数列表非空
	g.AddProduction("Args", []Symbol{"NonEmptyArgs"})                  
	//25. 函数参数列表空
	g.AddProduction("Args", []Symbol{})                                    
	//26. 非空参数单个表达式
	g.AddProduction("NonEmptyArgs", []Symbol{"Expr"})                      
	//27. 非空参数递归定义
	g.AddProduction("NonEmptyArgs", []Symbol{"NonEmptyArgs", ",", "Expr"}) 

	// === 多维数组索引 ===

	//28. 数组索引单个表达式
	g.AddProduction("IndexList", []Symbol{"Expr"})
	//29. 数组索引递归定义
	g.AddProduction("IndexList", []Symbol{"IndexList", ",", "Expr"})

	// === 条件表达式 ===
	//30. 条件表达式与运算
	g.AddProduction("Cond", []Symbol{"Cond", "&&", "Cond"})
	//31. 条件表达式或运算
	g.AddProduction("Cond", []Symbol{"Cond", "||", "Cond"})
	//32. 条件表达式非运算
	g.AddProduction("Cond", []Symbol{"!", "Cond"})
	//33. 条件表达式小于比较
	g.AddProduction("Cond", []Symbol{"Expr", "<", "Expr"})
	//34. 条件表达式等于比较
	g.AddProduction("Cond", []Symbol{"Expr", "==", "Expr"})
	//35. 条件表达式括号
	g.AddProduction("Cond", []Symbol{"(", "Cond", ")"})
	//36. 条件表达式单一表达式
	g.AddProduction("Cond", []Symbol{"Expr"})

	return g
}

// func DefineGrammar() *Grammar {
// 	g := NewGrammar("S'")
// 	g.AddProduction("S'", []Symbol{"Stmt"})

// 	// === 声明语句 ===
// 	g.AddProduction("Stmt", []Symbol{"Decl"})
// 	g.AddProduction("Decl", []Symbol{"int", "id", "=", "num", ";"}) // 可扩展更多类型如 float 等

// 	// === 普通语句 ===
// 	g.AddProduction("Stmt", []Symbol{"id", "=", "Expr", ";"})                         // 变量赋值
// 	g.AddProduction("Stmt", []Symbol{"return", "Expr", ";"})                          // return
// 	g.AddProduction("Stmt", []Symbol{"Block"})                                        // 块语句
// 	g.AddProduction("Stmt", []Symbol{"if", "(", "Cond", ")", "Stmt"})                 // if
// 	g.AddProduction("Stmt", []Symbol{"if", "(", "Cond", ")", "Stmt", "else", "Stmt"}) // if-else
// 	g.AddProduction("Stmt", []Symbol{"while", "(", "Cond", ")", "Stmt"})              // while
// 	g.AddProduction("Stmt", []Symbol{"id", "[", "IndexList", "]", "=", "Expr", ";"})  // 数组赋值

// 	// === 块与语句序列 ===
// 	g.AddProduction("Block", []Symbol{"{", "StmtList", "}"})
// 	g.AddProduction("StmtList", []Symbol{"Stmt"})
// 	g.AddProduction("StmtList", []Symbol{"StmtList", "Stmt"})

// 	// === 表达式结构 ===
// 	g.AddProduction("Expr", []Symbol{"Expr", "+", "Term"})
// 	g.AddProduction("Expr", []Symbol{"Expr", "-", "Term"})
// 	g.AddProduction("Expr", []Symbol{"Term"})
// 	g.AddProduction("Term", []Symbol{"Term", "*", "Factor"})
// 	g.AddProduction("Term", []Symbol{"Term", "/", "Factor"})
// 	g.AddProduction("Term", []Symbol{"Factor"})

// 	// === 基本因子 + 函数调用 ===
// 	g.AddProduction("Factor", []Symbol{"id", "(", "Args", ")"}) // 函数调用
// 	g.AddProduction("Factor", []Symbol{"num"})
// 	g.AddProduction("Factor", []Symbol{"id"})
// 	g.AddProduction("Factor", []Symbol{"(", "Expr", ")"})

// 	// === 函数参数列表 ===
// 	g.AddProduction("Args", []Symbol{}) // 支持无参数调用
// 	g.AddProduction("Args", []Symbol{"Expr"})
// 	g.AddProduction("Args", []Symbol{"Args", ",", "Expr"})

// 	// === 多维数组索引 ===
// 	g.AddProduction("IndexList", []Symbol{"Expr"})
// 	g.AddProduction("IndexList", []Symbol{"IndexList", ",", "Expr"})

// 	// === 条件表达式 ===
// 	g.AddProduction("Cond", []Symbol{"Cond", "&&", "Cond"})
// 	g.AddProduction("Cond", []Symbol{"Cond", "||", "Cond"})
// 	g.AddProduction("Cond", []Symbol{"!", "Cond"})
// 	g.AddProduction("Cond", []Symbol{"Expr", "<", "Expr"})
// 	g.AddProduction("Cond", []Symbol{"Expr", "==", "Expr"})
// 	g.AddProduction("Cond", []Symbol{"(", "Cond", ")"})

// 	return g
// }
