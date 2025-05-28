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
	//0. 定义起始符号
	g.AddProduction("S'", []Symbol{"Program"})
	// ==== 程序入口 ====

	//1.
	g.AddProduction("Program", []Symbol{"StmtList"})
	//2.
	g.AddProduction("Program", []Symbol{"StmtList", "EOF"}) // 程序可以以 EOF 结束
	//3.程序可以包括函数定义
	g.AddProduction("Program", []Symbol{"FuncList"})

	//4. 函数列表
	g.AddProduction("FuncList", []Symbol{"Func"})
	//5. 函数列表递归定义
	g.AddProduction("FuncList", []Symbol{"FuncList", "Func"}) // 函数列表可以有多个函数定义

	// ==== 函数定义 ====
	//6. 函数定义
	g.AddProduction("Func", []Symbol{"Type", "id", "(", "Args", ")", "Block"}) // 函数定义

	// ==== 声明语句 ====

	//7. 声明语句提取为 Decl，保留语义分类
	g.AddProduction("Stmt", []Symbol{"Decl"})
	//8. Decl 语句定义，支持 int、float 等类型
	g.AddProduction("Decl", []Symbol{"Type", "id", "=", "Expr", ";"}) // 可扩展更多类型如 float 等

	//9. 类型定义

	g.AddProduction("Type", []Symbol{"type_kw"}) // 支持多种类型，如 int, float, double 等

	// ==== 普通语句 ====

	//10. 变量赋值
	g.AddProduction("Stmt", []Symbol{"id", "=", "Expr", ";"})
	//11. return 语句
	g.AddProduction("Stmt", []Symbol{"return", "Expr", ";"})
	//12. 块语句
	g.AddProduction("Stmt", []Symbol{"Block"})
	//13. if 语句
	g.AddProduction("Stmt", []Symbol{"if", "(", "Cond", ")", "Stmt"})
	//14. if-else 语句
	g.AddProduction("Stmt", []Symbol{"if", "(", "Cond", ")", "Stmt", "else", "Stmt"})
	//15. while 语句
	g.AddProduction("Stmt", []Symbol{"while", "(", "Cond", ")", "Stmt"})
	//16. 数组赋值
	g.AddProduction("Stmt", []Symbol{"id", "[", "IndexList", "]", "=", "Expr", ";"})
	//17. 函数调用语句
	g.AddProduction("Stmt", []Symbol{"id", "(", "Args", ")", ";"})

	// === 块与语句序列 ===

	//18. 块语句定义
	g.AddProduction("Block", []Symbol{"{", "StmtList", "}"})
	//19. 语句列表空定义
	g.AddProduction("StmtList", []Symbol{})
	//20. 语句列表定义
	g.AddProduction("StmtList", []Symbol{"Stmt"})
	//21. 语句列表递归定义
	g.AddProduction("StmtList", []Symbol{"StmtList", "Stmt"})

	// === 表达式结构 ===

	//22. 表达式加法
	g.AddProduction("Expr", []Symbol{"Expr", "+", "Term"})
	//23. 表达式减法
	g.AddProduction("Expr", []Symbol{"Expr", "-", "Term"})
	//24. 表达式单一项
	g.AddProduction("Expr", []Symbol{"Term"})
	//25. 乘法项
	g.AddProduction("Term", []Symbol{"Term", "*", "CastExpr"}) // 乘法项可以是强制类型转换表达式
	//26. 除法项
	g.AddProduction("Term", []Symbol{"Term", "/", "CastExpr"}) // 除法项可以是强制类型转换表达式
	//27. 基本因子
	g.AddProduction("Term", []Symbol{"CastExpr"}) // 基本因子可以是强制类型转换表达式

	// === 强制类型转换 ===
	//28.
	g.AddProduction("CastExpr", []Symbol{"CastPrefix", "Factor"}) // 强制类型转换
	//29.
	g.AddProduction("CastExpr", []Symbol{"Factor"}) // 基本因子作为强制类型转换的基础
	//30.
	g.AddProduction("CastPrefix", []Symbol{"(", "Type", ")"})
	// === 基本因子 + 函数调用 ===

	//31. 函数调用
	g.AddProduction("Factor", []Symbol{"id", "(", "Args", ")"}) // 函数调用
	//32. 数字因子
	g.AddProduction("Factor", []Symbol{"num"})
	//33. 标识符因子
	g.AddProduction("Factor", []Symbol{"id"})
	//34. 括号表达式
	g.AddProduction("Factor", []Symbol{"(", "Expr", ")"})
	//35. 数组索引
	g.AddProduction("Factor", []Symbol{"id", "[", "IndexList", "]"}) // 数组索引

	// === 函数参数列表 ===
	//36. 函数参数列表非空
	g.AddProduction("Args", []Symbol{"NonEmptyArgs"})
	//37. 函数参数列表空
	g.AddProduction("Args", []Symbol{})
	//38. 非空参数单个表达式
	g.AddProduction("NonEmptyArgs", []Symbol{"Expr"})
	//39. 非空参数递归定义
	g.AddProduction("NonEmptyArgs", []Symbol{"NonEmptyArgs", ",", "Expr"})
	//40. type+id
	g.AddProduction("NonEmptyArgs", []Symbol{"Type", "id"})
	//41. 支持带默认值的参数
	g.AddProduction("NonEmptyArgs", []Symbol{"Type", "id", "=", "Expr"}) //
	//42. 非空参数可以是多个类型+标识符
	g.AddProduction("NonEmptyArgs", []Symbol{"NonEmptyArgs", ",", "Type", "id"}) // 支持多个参数
	//43. 非空参数可以是多个类型+标识符+默认值
	g.AddProduction("NonEmptyArgs", []Symbol{"NonEmptyArgs", ",", "Type", "id", "=", "Expr"}) // 支持多个参数
	// === 多维数组索引 ===

	//44. 数组索引单个表达式
	g.AddProduction("IndexList", []Symbol{"Expr"})
	//45. 数组索引递归定义
	g.AddProduction("IndexList", []Symbol{"IndexList", ",", "Expr"})

	// === 条件表达式 ===
	//46. 条件表达式与运算
	g.AddProduction("Cond", []Symbol{"Cond", "&&", "Cond"})
	//47. 条件表达式或运算
	g.AddProduction("Cond", []Symbol{"Cond", "||", "Cond"})
	//48. 条件表达式非运算
	g.AddProduction("Cond", []Symbol{"!", "Cond"})
	//49. 条件表达式小于比较
	g.AddProduction("Cond", []Symbol{"Expr", "<", "Expr"})
	//50.  条件表达式大于比较
	g.AddProduction("Cond", []Symbol{"Expr", ">", "Expr"})
	//51. 条件表达式小于等于比较
	g.AddProduction("Cond", []Symbol{"Expr", "<=", "Expr"})
	//52. 条件表达式大于等于比较
	g.AddProduction("Cond", []Symbol{"Expr", ">=", "Expr"})
	//53. 条件表达式等于比较
	g.AddProduction("Cond", []Symbol{"Expr", "!=", "Expr"})
	//54. 条件表达式等于比较
	g.AddProduction("Cond", []Symbol{"Expr", "==", "Expr"})

	//55. 条件表达式括号
	g.AddProduction("Cond", []Symbol{"(", "Cond", ")"})
	//56. 条件表达式单一表达式
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
