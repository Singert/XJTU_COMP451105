package semantic

import (
	"fmt"
	"lab5/lexer"
	"strings"
)

// ActionFuncs maps production rule indices to semantic actions.
// Each function takes the list of attributes (interface{}) popped from the attribute stack
// during a reduce action, and returns a new ASTNode (or any semantic value).
// ActionFuncs maps production index to semantic action function.

var ActionFuncs = map[int]func([]interface{}) interface{}{
	// 文法 0: S' -> Program
	0: func(children []interface{}) interface{} {
		return children[0]
	},
	// 文法 1: Program -> StmtList
	1: func(children []interface{}) interface{} {
		return &ASTNode{
			Type: "Program",
			Args: children[0].([]*ASTNode),
		}
	},

	// 文法 2: Program -> StmtList EOF
	2: func(children []interface{}) interface{} {
		return &ASTNode{
			Type: "Program",
			Args: children[0].([]*ASTNode),
		}
	},

	// 文法 3: Stmt -> Func
	3: func(children []interface{}) interface{} {
		return children[0]
	},

	// 文法 4: Func -> Type id ( Args ) Block
	4: func(children []interface{}) interface{} {
		// 提取形参列表字符串
		args := children[3].([]*ASTNode)
		paramList := make([]string, len(args))
		for i, arg := range args {
			paramList[i] = fmt.Sprintf("%s:%s", arg.Left.Value, arg.Value) // 格式化为 "参数名:类型"
		}
		paramListStr := fmt.Sprintf("(%s)", strings.Join(paramList, ", ")) // 拼接为 "(参数1:类型1, 参数2:类型2)"

		// 生成函数声明的四元组
		quadruples = append(quadruples, Quadruple{
			Op:     "func_decl",
			Left:   children[0].(*ASTNode).Value,     // 返回类型
			Right:  paramListStr,                     // 形参列表字符串
			Result: children[1].(lexer.Token).Lexeme, // 函数名
		})

		// 构造函数节点
		return &ASTNode{
			Type:  "Func",
			Value: children[0].(*ASTNode).Value, // 返回类型
			Left: &ASTNode{
				Type:  "id",
				Value: children[1].(lexer.Token).Lexeme, // 函数名
			},
			Args:  args,                   // 形参列表
			Right: children[5].(*ASTNode), // 函数体
		}
	},

	// 文法 5: Stmt -> Decl
	5: func(children []interface{}) interface{} {
		return children[0]
	},

	// 文法 6: Decl -> Type id = Expr ;
	6: func(children []interface{}) interface{} {
		// 提取类型、变量名和表达式节点
		typNode := children[0].(*ASTNode)
		idToken := children[1].(lexer.Token)
		exprNode := children[3].(*ASTNode)

		// 生成声明的四元组
		quadruples = append(quadruples, Quadruple{
			Op:     "decl",
			Left:   typNode.Value,  // 类型
			Right:  exprNode.Value, // 表达式的值
			Result: idToken.Lexeme, // 变量名
		})

		// 构造声明的 AST 节点
		return &ASTNode{
			Type:  "Decl",
			Value: typNode.Value,                               // 类型
			Left:  &ASTNode{Type: "id", Value: idToken.Lexeme}, // 变量名
			Right: exprNode,                                    // 右侧表达式
		}
	},
	// 文法 7: Type -> type_kw
	7: func(children []interface{}) interface{} {
		return &ASTNode{Type: "Type", Value: children[0].(lexer.Token).Lexeme}
	},

	// 文法 8: Stmt -> id = Expr ;
	8: func(children []interface{}) interface{} {
		idToken := children[0].(lexer.Token) // 提取变量名
		exprNode := children[2].(*ASTNode)   // 提取表达式节点

		// 生成赋值操作的四元组
		quadruples = append(quadruples, Quadruple{
			Op:     "=",
			Left:   exprNode.Value, // 表达式的值
			Right:  "",             // 赋值操作没有右操作数
			Result: idToken.Lexeme, // 变量名
		})

		// 构造赋值语句的 AST 节点
		return &ASTNode{
			Type:  "=",
			Left:  &ASTNode{Type: "id", Value: idToken.Lexeme}, // 左侧变量名
			Right: exprNode,                                    // 右侧表达式
		}
	},

	// 文法 9: Stmt -> return Expr ;
	9: func(children []interface{}) interface{} {
		exprNode := children[1].(*ASTNode) // 提取表达式节点

		// 生成返回操作的四元组
		quadruples = append(quadruples, Quadruple{
			Op:     "return",
			Left:   exprNode.Value, // 返回的表达式值
			Right:  "",             // 返回操作没有右操作数
			Result: "",             // 返回操作没有结果变量
		})

		// 构造返回语句的 AST 节点
		return &ASTNode{
			Type: "return",
			Left: exprNode,
		}
	},

	// 文法 10: Stmt -> Block
	10: func(children []interface{}) interface{} {
		return children[0]
	},

	// 文法 11: Stmt -> if ( Cond ) Stmt
	11: func(children []interface{}) interface{} {
		labelIfTrue := fmt.Sprintf("Label_if_true_%d", labelCounter)
		labelCounter++

		// 生成条件跳转的四元组
		quadruples = append(quadruples, Quadruple{
			Op:     "if",
			Left:   children[2].(*ASTNode).Value, // 条件表达式
			Right:  "",                           // 没有右操作数
			Result: labelIfTrue,                  // 跳转标签
		})

		return &ASTNode{
			Type:  "if",
			Left:  children[2].(*ASTNode), // 条件表达式
			Right: children[4].(*ASTNode), // if 语句块
		}
	},

	// 文法 12: Stmt -> if ( Cond ) Stmt else Stmt
	12: func(children []interface{}) interface{} {
		labelElse := fmt.Sprintf("Label_else_%d", labelCounter)
		labelCounter++
		labelIfEnd := fmt.Sprintf("Label_if_end_%d", labelCounter)
		labelCounter++

		// 条件跳转到 else 的四元组
		quadruples = append(quadruples, Quadruple{
			Op:     "if",
			Left:   children[2].(*ASTNode).Value, // 条件表达式
			Right:  "",                           // 没有右操作数
			Result: labelElse,                    // 跳转到 else 的标签
		})

		// 无条件跳转到 if 结束的四元组
		quadruples = append(quadruples, Quadruple{
			Op:     "goto",
			Left:   "",         // 无条件跳转
			Right:  "",         // 没有右操作数
			Result: labelIfEnd, // 跳转到 if 结束的标签
		})

		return &ASTNode{
			Type: "ifelse",
			Left: children[2].(*ASTNode), // 条件表达式
			Right: &ASTNode{
				Type:  "else",
				Left:  children[4].(*ASTNode), // if 语句块
				Right: children[6].(*ASTNode), // else 语句块
			},
		}
	},
	// 文法 13: Stmt -> while ( Cond ) Stmt
	13: func(children []interface{}) interface{} {
		labelWhileStart := fmt.Sprintf("while_start_%d", labelCounter)
		labelCounter++
		labelWhileEnd := fmt.Sprintf("while_end_%d", labelCounter)
		labelCounter++

		// 生成循环开始标签
		quadruples = append(quadruples, Quadruple{
			Op:     "label",
			Left:   "",
			Right:  "",
			Result: labelWhileStart,
		})

		// 条件跳转到循环结束
		quadruples = append(quadruples, Quadruple{
			Op:     "if",
			Left:   children[2].(*ASTNode).Value, // 条件表达式
			Right:  "",
			Result: labelWhileEnd,
		})

		// 循环体
		quadruples = append(quadruples, Quadruple{
			Op:     "goto",
			Left:   "",
			Right:  "",
			Result: labelWhileStart,
		})

		// 循环结束标签
		quadruples = append(quadruples, Quadruple{
			Op:     "label",
			Left:   "",
			Right:  "",
			Result: labelWhileEnd,
		})

		return &ASTNode{
			Type:  "while",
			Left:  children[2].(*ASTNode), // 条件表达式
			Right: children[4].(*ASTNode), // 循环体
		}
	},

	// 文法 14: Stmt -> id MultiIndex = Expr ;
	14: func(children []interface{}) interface{} {
		quadruples = append(quadruples, Quadruple{
			Op:     "arr_assign",
			Left:   fmt.Sprintf("%s[%v]", children[0].(lexer.Token).Lexeme, children[1]), // 数组名和索引
			Right:  children[3].(*ASTNode).Value,                                         // 表达式的值
			Result: "",                                                                   // 赋值操作没有结果变量
		})

		return &ASTNode{
			Type:  "arr_assign",
			Value: children[0].(lexer.Token).Lexeme,
			Args:  children[1].([]*ASTNode),
			Right: children[3].(*ASTNode),
		}
	},

	// 文法 15: Stmt -> id ( Args ) ;
	15: func(children []interface{}) interface{} {
		funcName := children[0].(lexer.Token).Lexeme
		args := children[2].([]*ASTNode)

		// 为每个参数生成四元组
		for _, arg := range args {
			quadruples = append(quadruples, Quadruple{
				Op:     "param",
				Left:   arg.Value, // 参数值
				Right:  "",
				Result: "",
			})
		}

		// 生成函数调用四元组
		quadruples = append(quadruples, Quadruple{
			Op:     "call",
			Left:   funcName,                     // 函数名
			Right:  fmt.Sprintf("%d", len(args)), // 参数个数
			Result: "",                           // 函数调用没有直接结果
		})

		return &ASTNode{
			Type:  "call_stmt",
			Value: children[0].(lexer.Token).Lexeme,
			Args:  children[2].([]*ASTNode),
		}
	},

	// 文法 16: Block -> { StmtList }
	16: func(children []interface{}) interface{} {
		return &ASTNode{
			Type: "block",
			Args: children[1].([]*ASTNode),
		}
	},

	// 文法 17: StmtList -> ε
	17: func(children []interface{}) interface{} {
		return []*ASTNode{}
	},

	// 文法 18: StmtList -> Stmt
	18: func(children []interface{}) interface{} {
		return []*ASTNode{children[0].(*ASTNode)}
	},
	// 文法 19: StmtList -> StmtList Stmt
	19: func(children []interface{}) interface{} {
		return append(children[0].([]*ASTNode), children[1].(*ASTNode))
	},

	// 文法 20: Expr -> Expr + Term
	20: func(children []interface{}) interface{} {
		tempVar := fmt.Sprintf("t%d", tempVarCounter)
		tempVarCounter++

		quadruples = append(quadruples, Quadruple{
			Op:     "+",
			Left:   children[0].(*ASTNode).Value, // 左操作数
			Right:  children[2].(*ASTNode).Value, // 右操作数
			Result: tempVar,                      // 临时变量存储结果
		})

		return &ASTNode{
			Type:  "+",
			Value: tempVar, // 返回临时变量作为结果
			Left:  children[0].(*ASTNode),
			Right: children[2].(*ASTNode),
		}
	},

	// 文法 21: Expr -> Expr - Term
	21: func(children []interface{}) interface{} {
		tempVar := fmt.Sprintf("t%d", tempVarCounter)
		tempVarCounter++

		quadruples = append(quadruples, Quadruple{
			Op:     "-",
			Left:   children[0].(*ASTNode).Value, // 左操作数
			Right:  children[2].(*ASTNode).Value, // 右操作数
			Result: tempVar,                      // 临时变量存储结果
		})

		return &ASTNode{
			Type:  "-",
			Value: tempVar, // 返回临时变量作为结果
			Left:  children[0].(*ASTNode),
			Right: children[2].(*ASTNode),
		}
	},

	// 文法 22: Expr -> Term
	22: func(children []interface{}) interface{} {
		return children[0]
	},

	// 文法 23: Term -> Term * CastExpr
	23: func(children []interface{}) interface{} {
		tempVar := fmt.Sprintf("t%d", tempVarCounter)
		tempVarCounter++

		quadruples = append(quadruples, Quadruple{
			Op:     "*",
			Left:   children[0].(*ASTNode).Value, // 左操作数
			Right:  children[2].(*ASTNode).Value, // 右操作数
			Result: tempVar,                      // 临时变量存储结果
		})

		return &ASTNode{
			Type:  "*",
			Value: tempVar, // 返回临时变量作为结果
			Left:  children[0].(*ASTNode),
			Right: children[2].(*ASTNode),
		}
	},

	// 文法 24: Term -> Term / CastExpr
	24: func(children []interface{}) interface{} {
		tempVar := fmt.Sprintf("t%d", tempVarCounter)
		tempVarCounter++

		quadruples = append(quadruples, Quadruple{
			Op:     "/",
			Left:   children[0].(*ASTNode).Value, // 左操作数
			Right:  children[2].(*ASTNode).Value, // 右操作数
			Result: tempVar,                      // 临时变量存储结果
		})

		return &ASTNode{
			Type:  "/",
			Value: tempVar, // 返回临时变量作为结果
			Left:  children[0].(*ASTNode),
			Right: children[2].(*ASTNode),
		}
	},

	// 文法 25: Term -> CastExpr
	25: func(children []interface{}) interface{} {
		return children[0]
	},

	26: func(children []interface{}) interface{} {
		tempVar := fmt.Sprintf("t%d", tempVarCounter)
		tempVarCounter++

		quadruples = append(quadruples, Quadruple{
			Op:     "cast",
			Left:   children[1].(*ASTNode).Value, // 被转换的值
			Right:  children[0].(*ASTNode).Value, // 目标类型
			Result: tempVar,                      // 临时变量存储结果
		})

		return &ASTNode{
			Type:  "cast",
			Value: tempVar, // 返回临时变量作为结果
			Right: children[1].(*ASTNode),
		}
	},

	// 文法 27: CastExpr -> Factor
	27: func(children []interface{}) interface{} {
		return children[0]
	},

	// 文法 28: CastPrefix -> ( Type )
	28: func(children []interface{}) interface{} {
		return &ASTNode{
			Type:  "CastPrefix",
			Value: children[1].(*ASTNode).Value,
		}
	},

	// 文法 29: Factor -> id ( Args )
	29: func(children []interface{}) interface{} {
		funcName := children[0].(lexer.Token).Lexeme
		args := children[2].([]*ASTNode)

		// 为每个参数生成四元组
		for _, arg := range args {
			quadruples = append(quadruples, Quadruple{
				Op:     "param",
				Left:   arg.Value, // 参数值
				Right:  "",
				Result: "",
			})
		}

		// 生成函数调用四元组
		quadruples = append(quadruples, Quadruple{
			Op:     "call",
			Left:   funcName,                     // 函数名
			Right:  fmt.Sprintf("%d", len(args)), // 参数个数
			Result: "",                           // 函数调用没有直接结果
		})

		return &ASTNode{
			Type:  "call",
			Value: funcName,
			Args:  args,
		}
	},

	// 文法 30: Factor -> num
	30: func(children []interface{}) interface{} {
		return &ASTNode{Type: "num", Value: children[0].(lexer.Token).Lexeme}
	},

	// 文法 31: Factor -> float
	31: func(children []interface{}) interface{} {
		return &ASTNode{Type: "float", Value: children[0].(lexer.Token).Lexeme}
	},

	// 文法 32: Factor -> char
	32: func(children []interface{}) interface{} {
		return &ASTNode{Type: "char", Value: children[0].(lexer.Token).Lexeme}
	},

	// 文法 33: Factor -> string
	33: func(children []interface{}) interface{} {
		return &ASTNode{Type: "string", Value: children[0].(lexer.Token).Lexeme}
	},

	// 文法 34: Factor -> id
	34: func(children []interface{}) interface{} {
		return &ASTNode{Type: "id", Value: children[0].(lexer.Token).Lexeme}
	},

	// 文法 35: Factor -> ( Expr )
	35: func(children []interface{}) interface{} {
		return children[1]
	},

	36: func(children []interface{}) interface{} {
		idToken := children[0].(lexer.Token)
		multiIndices := children[1].([]*ASTNode)

		// 为数组访问生成四元组
		tempVar := fmt.Sprintf("t%d", tempVarCounter)
		tempVarCounter++

		quadruples = append(quadruples, Quadruple{
			Op:     "array_access",
			Left:   idToken.Lexeme,                  // 数组名
			Right:  fmt.Sprintf("%v", multiIndices), // 索引列表
			Result: tempVar,                         // 临时变量存储结果
		})

		return &ASTNode{
			Type:  "array_access",
			Value: tempVar, // 返回临时变量作为结果
			Args:  multiIndices,
		}
	},

	// 文法 37: Args -> NonEmptyArgs
	37: func(children []interface{}) interface{} {
		return children[0].([]*ASTNode)
	},

	// 文法 38: Args ->
	38: func(children []interface{}) interface{} {
		return []*ASTNode{}
	},
	// 文法 39: NonEmptyArgs -> Expr
	39: func(children []interface{}) interface{} {
		return []*ASTNode{children[0].(*ASTNode)}
	},

	// 文法 40: NonEmptyArgs -> NonEmptyArgs , Expr
	40: func(children []interface{}) interface{} {
		return append(children[0].([]*ASTNode), children[2].(*ASTNode))
	},

	// 文法 41: NonEmptyArgs -> Type id
	41: func(children []interface{}) interface{} {
		fmt.Printf("\n\n\nActionFuncs[41]: ASTNode = %v\n", &ASTNode{Type: "id", Value: children[1].(lexer.Token).Lexeme})
		return []*ASTNode{
			{
				Type:  "Arg",
				Value: children[0].(*ASTNode).Value,
				Left:  &ASTNode{Type: "id", Value: children[1].(lexer.Token).Lexeme},
			},
		}
	},

	// 文法 42: NonEmptyArgs -> Type id = Expr
	42: func(children []interface{}) interface{} {
		return []*ASTNode{
			{
				Type:  "Arg",
				Value: children[0].(*ASTNode).Value,
				Left:  &ASTNode{Type: "id", Value: children[1].(lexer.Token).Lexeme},
				Right: children[3].(*ASTNode),
			},
		}
	},

	// 文法 43: NonEmptyArgs -> NonEmptyArgs , Type id
	43: func(children []interface{}) interface{} {
		args := children[0].([]*ASTNode)
		args = append(args, &ASTNode{
			Type:  "Arg",
			Value: children[2].(*ASTNode).Value,
			Left:  &ASTNode{Type: "id", Value: children[3].(lexer.Token).Lexeme},
		})
		return args
	},

	// 文法 44: NonEmptyArgs -> NonEmptyArgs , Type id = Expr
	44: func(children []interface{}) interface{} {
		args := children[0].([]*ASTNode)
		args = append(args, &ASTNode{
			Type:  "Arg",
			Value: children[2].(*ASTNode).Value,
			Left:  &ASTNode{Type: "id", Value: children[3].(lexer.Token).Lexeme},
			Right: children[5].(*ASTNode),
		})
		return args
	},

	// 文法 45: IndexList -> Expr
	45: func(children []interface{}) interface{} {
		return []*ASTNode{children[0].(*ASTNode)}
	},

	// 文法 46: IndexList -> IndexList , Expr
	46: func(children []interface{}) interface{} {
		return append(children[0].([]*ASTNode), children[2].(*ASTNode))
	},

	// 文法 47: Cond -> Cond && Cond
	47: func(children []interface{}) interface{} {
		tempVar := fmt.Sprintf("t%d", tempVarCounter)
		tempVarCounter++

		quadruples = append(quadruples, Quadruple{
			Op:     "&&",
			Left:   children[0].(*ASTNode).Value, // 左操作数
			Right:  children[2].(*ASTNode).Value, // 右操作数
			Result: tempVar,                      // 临时变量存储结果
		})

		return &ASTNode{
			Type:  "&&",
			Value: tempVar, // 返回临时变量作为结果
			Left:  children[0].(*ASTNode),
			Right: children[2].(*ASTNode),
		}
	},

	// 文法 48: Cond -> Cond || Cond
	48: func(children []interface{}) interface{} {
		tempVar := fmt.Sprintf("t%d", tempVarCounter)
		tempVarCounter++

		quadruples = append(quadruples, Quadruple{
			Op:     "||",
			Left:   children[0].(*ASTNode).Value, // 左操作数
			Right:  children[2].(*ASTNode).Value, // 右操作数
			Result: tempVar,                      // 临时变量存储结果
		})

		return &ASTNode{
			Type:  "||",
			Value: tempVar, // 返回临时变量作为结果
			Left:  children[0].(*ASTNode),
			Right: children[2].(*ASTNode),
		}
	},

	// 文法 49: Cond -> ! Cond
	49: func(children []interface{}) interface{} {
		tempVar := fmt.Sprintf("t%d", tempVarCounter)
		tempVarCounter++

		quadruples = append(quadruples, Quadruple{
			Op:     "!",
			Left:   children[1].(*ASTNode).Value, // 操作数
			Right:  "",                           // 逻辑非没有右操作数
			Result: tempVar,                      // 临时变量存储结果
		})

		return &ASTNode{
			Type:  "!",
			Value: tempVar, // 返回临时变量作为结果
			Right: children[1].(*ASTNode),
		}
	},

	// 文法 50: Cond -> Expr < Expr
	50: func(children []interface{}) interface{} {
		tempVar := fmt.Sprintf("t%d", tempVarCounter)
		tempVarCounter++

		quadruples = append(quadruples, Quadruple{
			Op:     "<",
			Left:   children[0].(*ASTNode).Value, // 左操作数
			Right:  children[2].(*ASTNode).Value, // 右操作数
			Result: tempVar,                      // 临时变量存储结果
		})

		return &ASTNode{
			Type:  "<",
			Value: tempVar, // 返回临时变量作为结果
			Left:  children[0].(*ASTNode),
			Right: children[2].(*ASTNode),
		}
	},

	// 文法 51: Cond -> Expr > Expr
	51: func(children []interface{}) interface{} {
		tempVar := fmt.Sprintf("t%d", tempVarCounter)
		tempVarCounter++

		quadruples = append(quadruples, Quadruple{
			Op:     ">",
			Left:   children[0].(*ASTNode).Value, // 左操作数
			Right:  children[2].(*ASTNode).Value, // 右操作数
			Result: tempVar,                      // 临时变量存储结果
		})

		return &ASTNode{
			Type:  ">",
			Value: tempVar, // 返回临时变量作为结果
			Left:  children[0].(*ASTNode),
			Right: children[2].(*ASTNode),
		}
	},

	// 文法 52: Cond -> Expr <= Expr
	52: func(children []interface{}) interface{} {
		tempVar := fmt.Sprintf("t%d", tempVarCounter)
		tempVarCounter++

		quadruples = append(quadruples, Quadruple{
			Op:     "<=",
			Left:   children[0].(*ASTNode).Value, // 左操作数
			Right:  children[2].(*ASTNode).Value, // 右操作数
			Result: tempVar,                      // 临时变量存储结果
		})

		return &ASTNode{
			Type:  "<=",
			Value: tempVar, // 返回临时变量作为结果
			Left:  children[0].(*ASTNode),
			Right: children[2].(*ASTNode),
		}
	},

	// 文法 53: Cond -> Expr >= Expr
	53: func(children []interface{}) interface{} {
		tempVar := fmt.Sprintf("t%d", tempVarCounter)
		tempVarCounter++

		quadruples = append(quadruples, Quadruple{
			Op:     ">=",
			Left:   children[0].(*ASTNode).Value, // 左操作数
			Right:  children[2].(*ASTNode).Value, // 右操作数
			Result: tempVar,                      // 临时变量存储结果
		})

		return &ASTNode{
			Type:  ">=",
			Value: tempVar, // 返回临时变量作为结果
			Left:  children[0].(*ASTNode),
			Right: children[2].(*ASTNode),
		}
	},

	// 文法 54: Cond -> Expr != Expr
	54: func(children []interface{}) interface{} {
		tempVar := fmt.Sprintf("t%d", tempVarCounter)
		tempVarCounter++

		quadruples = append(quadruples, Quadruple{
			Op:     "!=",
			Left:   children[0].(*ASTNode).Value, // 左操作数
			Right:  children[2].(*ASTNode).Value, // 右操作数
			Result: tempVar,                      // 临时变量存储结果
		})

		return &ASTNode{
			Type:  "!=",
			Value: tempVar, // 返回临时变量作为结果
			Left:  children[0].(*ASTNode),
			Right: children[2].(*ASTNode),
		}
	},

	// 文法 55: Cond -> Expr == Expr
	55: func(children []interface{}) interface{} {
		tempVar := fmt.Sprintf("t%d", tempVarCounter)
		tempVarCounter++

		quadruples = append(quadruples, Quadruple{
			Op:     "==",
			Left:   children[0].(*ASTNode).Value, // 左操作数
			Right:  children[2].(*ASTNode).Value, // 右操作数
			Result: tempVar,                      // 临时变量存储结果
		})

		return &ASTNode{
			Type:  "==",
			Value: tempVar, // 返回临时变量作为结果
			Left:  children[0].(*ASTNode),
			Right: children[2].(*ASTNode),
		}
	},

	// 文法 56: Cond -> ( Cond )
	56: func(children []interface{}) interface{} {
		return children[1]
	},

	// 文法 57: Cond -> Expr
	57: func(children []interface{}) interface{} {
		return children[0]
	},

	// 文法 58: 无初始化数组声明
	58: func(children []interface{}) interface{} {
		typNode := children[0].(*ASTNode)
		idToken := children[1].(lexer.Token)
		multiIndices := children[2].([]*ASTNode)

		quadruples = append(quadruples, Quadruple{
			Op:     "array_decl",
			Left:   typNode.Value,                   // 数组类型
			Right:  fmt.Sprintf("%v", multiIndices), // 索引列表
			Result: idToken.Lexeme,                  // 数组名
		})

		return &ASTNode{
			Type:  "ArrayDeclMulti",
			Value: typNode.Value,
			Left:  &ASTNode{Type: "id", Value: idToken.Lexeme},
			Args:  multiIndices,
		}
	},
	// 文法 59: 带初始化数组声明（初始化列表或字符串字面量）
	59: func(children []interface{}) interface{} {
		typNode := children[0].(*ASTNode)
		idToken := children[1].(lexer.Token)
		multiIndices := children[2].([]*ASTNode)
		initExpr := children[4].(*ASTNode)

		quadruples = append(quadruples, Quadruple{
			Op:     "array_decl_init",
			Left:   typNode.Value,                   // 数组类型
			Right:  fmt.Sprintf("%v", multiIndices), // 索引列表
			Result: idToken.Lexeme,                  // 数组名
		})

		quadruples = append(quadruples, Quadruple{
			Op:     "init",
			Left:   idToken.Lexeme, // 数组名
			Right:  initExpr.Value, // 初始化表达式
			Result: "",             // 初始化操作没有结果变量
		})

		return &ASTNode{
			Type:  "ArrayDeclMultiInit",
			Value: typNode.Value,
			Left:  &ASTNode{Type: "id", Value: idToken.Lexeme},
			Args:  multiIndices,
			Right: initExpr,
		}
	},

	// 文法 60: MultiIndex -> [ IndexList ] MultiIndex
	60: func(children []interface{}) interface{} {
		indexList := children[1].([]*ASTNode)
		multiIndexNext := children[3].([]*ASTNode)
		return append(indexList, multiIndexNext...)
	},

	// 文法 61: MultiIndex -> ε
	61: func(children []interface{}) interface{} {
		return []*ASTNode{}
	},

	// 文法 62: Decl -> Type id MultiIndex = InitList ;
	62: func(children []interface{}) interface{} {
		typNode := children[0].(*ASTNode)
		idToken := children[1].(lexer.Token)
		multiIndices := children[2].([]*ASTNode)
		initList := children[4].(*ASTNode)
		indvidualIndices := make([]string, len(multiIndices))
		for i, idx := range multiIndices {
			indvidualIndices[i] = idx.Value // 获取每个索引的值
		}
		quadruples = append(quadruples, Quadruple{
			Op:     "array_decl_init_list",
			Left:   typNode.Value,                       // 数组类型
			Right:  fmt.Sprintf("%v", indvidualIndices), // 索引列表
			Result: idToken.Lexeme,                      // 数组名
		})
		values := make([]string, len(initList.Args))
		for i, arg := range initList.Args {
			values[i] = arg.Value // 获取初始化列表中的每个值
		}
		quadruples = append(quadruples, Quadruple{
			Op:     "init_list",
			Left:   idToken.Lexeme,            // 数组名
			Right:  fmt.Sprintf("%v", values), // 初始化列表
			Result: "",                        // 初始化操作没有结果变量
		})

		return &ASTNode{
			Type:  "ArrayDeclMultiInitList",
			Value: typNode.Value,
			Left:  &ASTNode{Type: "id", Value: idToken.Lexeme},
			Args:  multiIndices,
			Right: initList,
		}
	},

	// 文法 63: InitList -> { NonEmptyInitList }
	63: func(children []interface{}) interface{} {
		return &ASTNode{
			Type: "InitList",
			Args: children[1].([]*ASTNode),
		}
	},

	// 文法 64: InitList -> {}
	64: func(children []interface{}) interface{} {
		return &ASTNode{
			Type: "InitList",
			Args: []*ASTNode{},
		}
	},

	// 文法 65: NonEmptyInitList -> Expr
	65: func(children []interface{}) interface{} {
		return []*ASTNode{children[0].(*ASTNode)}
	},

	// 文法 66: NonEmptyInitList -> NonEmptyInitList , Expr
	66: func(children []interface{}) interface{} {
		return append(children[0].([]*ASTNode), children[2].(*ASTNode))
	},

	// 文法 67: Expr -> InitList
	67: func(children []interface{}) interface{} {
		return children[0]
	},

	// 文法 68: Factor -> - Factor (一元负号)
	68: func(children []interface{}) interface{} {
		return &ASTNode{
			Type:  "-",
			Right: children[1].(*ASTNode),
		}
	},

	// 文法 69: NonEmptyArgs -> Type id MultiIndex
	69: func(children []interface{}) interface{} {
		fmt.Printf("\n\n\nActionFuncs[69]\n\n\n")
		return []*ASTNode{
			{
				Type:  "Arg",
				Value: children[0].(*ASTNode).Value,
				Left:  &ASTNode{Type: "id", Value: children[1].(lexer.Token).Lexeme},
				Args:  children[2].([]*ASTNode), // 多维索引
			},
		}
	},

	// 文法 70: NonEmptyArgs -> Type id MultiIndex = Expr
	70: func(children []interface{}) interface{} {
		return []*ASTNode{
			{
				Type:  "Arg",
				Value: children[0].(*ASTNode).Value,
				Left:  &ASTNode{Type: "id", Value: children[1].(lexer.Token).Lexeme},
				Args:  children[2].([]*ASTNode), // 多维索引
				Right: children[4].(*ASTNode),   // 初始化表达式
			},
		}
	},

	// 文法 71: IndexList -> ε
	71: func(children []interface{}) interface{} {
		return []*ASTNode{} // 返回空的索引列表
	},

	// 文法 72: Stmt -> for ( ForInit ; Cond ; Expr ) Stmt
	72: func(children []interface{}) interface{} {
		forInit := children[2].(*ASTNode)
		cond := children[4].(*ASTNode)
		update := children[6].(*ASTNode)
		body := children[8].(*ASTNode)

		startLabel := fmt.Sprintf("for_start_%d", labelCounter)
		labelCounter++
		endLabel := fmt.Sprintf("for_end_%d", labelCounter)
		labelCounter++
		updateLabel := fmt.Sprintf("for_update_%d", labelCounter)
		labelCounter++

		// 生成循环开始标签
		quadruples = append(quadruples, Quadruple{
			Op:     "label",
			Left:   "",
			Right:  "",
			Result: startLabel,
		})

		// 生成 ForInit 的四元组
		quadruples = append(quadruples, Quadruple{
			Op:     "for_init",
			Left:   forInit.Value,
			Right:  "",
			Result: fmt.Sprintf("for_%s", startLabel),
		})

		// 条件跳转到循环结束
		quadruples = append(quadruples, Quadruple{
			Op:     "if",
			Left:   cond.Value,
			Right:  "",
			Result: endLabel,
		})

		// 循环体
		quadruples = append(quadruples, Quadruple{
			Op:     "body",
			Left:   body.Value,
			Right:  "",
			Result: fmt.Sprintf("for_%s", startLabel),
		})

		// 更新语句标签
		quadruples = append(quadruples, Quadruple{
			Op:     "label",
			Left:   "",
			Right:  "",
			Result: updateLabel,
		})

		// 更新语句
		quadruples = append(quadruples, Quadruple{
			Op:     "update",
			Left:   update.Value,
			Right:  "",
			Result: fmt.Sprintf("for_%s", updateLabel),
		})

		// 跳转回条件检查
		quadruples = append(quadruples, Quadruple{
			Op:     "goto",
			Left:   "",
			Right:  "",
			Result: startLabel,
		})

		// 循环结束标签
		quadruples = append(quadruples, Quadruple{
			Op:     "label",
			Left:   "",
			Right:  "",
			Result: endLabel,
		})

		return &ASTNode{
			Type:  "for",
			Args:  []*ASTNode{forInit, update}, // 可以存 forInit 和 更新语句
			Left:  cond,                        // 条件作为 Left
			Right: body,                        // 循环体作为 Right
		}
	},

	// 文法 73: ForInit -> Decl
	73: func(children []interface{}) interface{} {
		return children[0].(*ASTNode) // 返回声明节点
	},

	// 文法 74: ForInit -> Expr
	74: func(children []interface{}) interface{} {
		return children[0].(*ASTNode) // 返回表达式节点
	},

	// 文法 75: ForInit -> ε
	75: func(children []interface{}) interface{} {
		return &ASTNode{Type: "ForInitEmpty"} // 返回空的 ForInit
	},

	// 文法 76: Expr -> id = Expr
	76: func(children []interface{}) interface{} {
		idToken := children[0].(lexer.Token) // lexer.Token for "id"
		exprNode := children[2].(*ASTNode)   // 确保表达式是 ASTNode

		// 生成赋值操作的四元组
		quadruples = append(quadruples, Quadruple{
			Op:     "=",
			Left:   exprNode.Value, // 表达式的值
			Right:  "",             // 赋值操作没有右操作数
			Result: idToken.Lexeme, // 变量名
		})

		return &ASTNode{
			Type:  "=",
			Left:  &ASTNode{Type: "id", Value: idToken.Lexeme},
			Right: exprNode,
		}
	},
}
