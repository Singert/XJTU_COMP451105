package semantic

import (
	"fmt"
	"lab5/lexer"
)

// ActionFuncs maps production rule indices to semantic actions.
// Each function takes the list of attributes (interface{}) popped from the attribute stack
// during a reduce action, and returns a new ASTNode (or any semantic value).
// ActionFuncs maps production index to semantic action function.

var ActionFuncs = map[int]func([]interface{}) interface{}{
	// 0: S' -> Program
	0: func(children []interface{}) interface{} {
		return children[0]
	},
	// 1: Program -> StmtList
	1: func(children []interface{}) interface{} {
		return &ASTNode{
			Type: "Program",
			Args: children[0].([]*ASTNode), // StmtList
		}
	},
	// 2: Program -> StmtList EOF
	2: func(children []interface{}) interface{} {
		return &ASTNode{
			Type: "Program",
			Args: children[0].([]*ASTNode), // StmtList
		}
	},
	// 3: Program -> FuncList
	3: func(children []interface{}) interface{} {
		return &ASTNode{
			Type: "Program",
			Args: children[0].([]*ASTNode), // FuncList
		}
	},
	// 4: FuncList -> Func
	4: func(children []interface{}) interface{} {
		return []*ASTNode{children[0].(*ASTNode)}
	},
	// 5: FuncList -> FuncList Func
	5: func(children []interface{}) interface{} {
		list := children[0].([]*ASTNode)   // FuncList
		funcNode := children[1].(*ASTNode) // Func
		return append(list, funcNode)      // 返回新的 FuncList
	},
	// 4: Func -> Type id ( Args ) Block
	6: func(children []interface{}) interface{} {
		fmt.Printf("children[5] type: %T\n", children[5]) // 打印类型信息
		typNode := children[0].(*ASTNode)                 // Type 节点
		idToken := children[1].(lexer.Token)
		args := children[3].([]*ASTNode) // Args
		block := children[5].(*ASTNode)  // Block

		return &ASTNode{
			Type: "Func",
			Left: &ASTNode{
				Type:  "id",
				Value: idToken.Lexeme,
			},
			Right: block,
			Value: typNode.Value, // 存储函数返回类型
			Args:  args,          // 函数参数列表
		}
	},
	// 1: Stmt -> Decl
	7: func(children []interface{}) interface{} {
		return children[0]
	},

	// 2: Decl -> Type id = num ;
	8: func(children []interface{}) interface{} {
		typNode := children[0].(*ASTNode) // Type 节点
		idToken := children[1].(lexer.Token)
		expr := children[3].(*ASTNode)

		return &ASTNode{
			Type: "Decl",
			Left: &ASTNode{
				Type:  "id",
				Value: idToken.Lexeme,
			},
			Right: expr,
			// 额外存储类型信息
			Value: typNode.Value,
		}
	},

	//3. Type -> type_kw
	9: func(children []interface{}) interface{} {
		tok := children[0].(lexer.Token)
		return &ASTNode{Type: "Type", Value: tok.Lexeme} // 直接使用 type_kw 的词素作为类型值
	},

	// 3: Stmt -> id = Expr ;
	10: func(children []interface{}) interface{} {
		id := children[0].(lexer.Token)
		expr := children[2].(*ASTNode)
		return &ASTNode{Type: "=", Left: &ASTNode{Type: "id", Value: id.Lexeme}, Right: expr}
	},

	// 4: Stmt -> return Expr ;
	11: func(children []interface{}) interface{} {
		expr := children[1].(*ASTNode)
		return &ASTNode{Type: "return", Left: expr}
	},

	// 5: Stmt -> Block
	12: func(children []interface{}) interface{} {
		return children[0]
	},

	// 6: Stmt -> if ( Cond ) Stmt
	13: func(children []interface{}) interface{} {
		cond := children[2].(*ASTNode)
		stmt := children[4].(*ASTNode)
		return &ASTNode{Type: "if", Left: cond, Right: stmt}
	},

	// 7: Stmt -> if ( Cond ) Stmt else Stmt
	14: func(children []interface{}) interface{} {
		cond := children[2].(*ASTNode)
		thenStmt := children[4].(*ASTNode)
		elseStmt := children[6].(*ASTNode)
		return &ASTNode{Type: "ifelse", Left: cond, Right: &ASTNode{Type: "else", Left: thenStmt, Right: elseStmt}}
	},

	// 8: Stmt -> while ( Cond ) Stmt
	15: func(children []interface{}) interface{} {
		cond := children[2].(*ASTNode)
		body := children[4].(*ASTNode)
		return &ASTNode{Type: "while", Left: cond, Right: body}
	},

	// 9: Stmt -> id [ IndexList ] = Expr ;
	16: func(children []interface{}) interface{} {
		id := children[0].(lexer.Token)
		indices := children[2].([]*ASTNode)
		expr := children[5].(*ASTNode)
		return &ASTNode{Type: "arr_assign", Value: id.Lexeme, Args: indices, Right: expr}
	},
	// 10: Stmt -> id ( Args ) ;
	17: func(children []interface{}) interface{} {
		id := children[0].(lexer.Token)
		args := children[2].([]*ASTNode)
		return &ASTNode{
			Type:  "call_stmt",
			Value: id.Lexeme,
			Args:  args,
		}
	},

	// 10: Block -> { StmtList }
	18: func(children []interface{}) interface{} {
		fmt.Printf("Block children: %+v\n", children)
		return &ASTNode{Type: "block", Args: children[1].([]*ASTNode)}
	},
	// 11: StmtList -> empty
	19: func(children []interface{}) interface{} {
		return []*ASTNode{} // 返回空的语句列表
	},
	// 11: StmtList -> Stmt
	20: func(children []interface{}) interface{} {
		return []*ASTNode{children[0].(*ASTNode)}
	},

	// 12: StmtList -> StmtList Stmt
	21: func(children []interface{}) interface{} {
		list := children[0].([]*ASTNode)
		return append(list, children[1].(*ASTNode))
	},

	// Expressions:

	22: func(c []interface{}) interface{} {
		return &ASTNode{Type: "+", Left: c[0].(*ASTNode), Right: c[2].(*ASTNode)}
	},
	23: func(c []interface{}) interface{} {
		return &ASTNode{Type: "-", Left: c[0].(*ASTNode), Right: c[2].(*ASTNode)}
	},
	24: func(c []interface{}) interface{} { return c[0] },

	25: func(c []interface{}) interface{} {
		return &ASTNode{Type: "*", Left: c[0].(*ASTNode), Right: c[2].(*ASTNode)}
	},
	26: func(c []interface{}) interface{} {
		return &ASTNode{Type: "/", Left: c[0].(*ASTNode), Right: c[2].(*ASTNode)}
	},
	27: func(c []interface{}) interface{} { return c[0] },

	// CastExpr:
	// CastExpr -> CastPrefix + CastExpr
	28: func(c []interface{}) interface{} {
		castPrefix := c[0].(*ASTNode) // CastPrefix 节点
		factor := c[1].(*ASTNode)     // CastExpr 节点
		return &ASTNode{
			Type:  "Cast",
			Value: castPrefix.Value, // 使用 CastPrefix 的值作为类型
			Left:  factor,           // CastExpr 作为右子节点
		}
	},
	// CastExpr -> Factor
	29: func(c []interface{}) interface{} {
		return c[0]
	},
	// CastPrefix -> "Type"
	30: func(c []interface{}) interface{} {
		typNode := c[1].(*ASTNode) // Type 节点
		return &ASTNode{
			Type:  "CastPrefix",
			Value: typNode.Value, // 使用 Type 的值作为 CastPrefix 的值
		}
	},
	// Factor:
	31: func(c []interface{}) interface{} {
		id := c[0].(lexer.Token)
		args := c[2].([]*ASTNode)
		return &ASTNode{Type: "call", Value: id.Lexeme, Args: args}
	},
	32: func(c []interface{}) interface{} {
		tok := c[0].(lexer.Token)
		return &ASTNode{Type: "num", Value: tok.Lexeme}
	},
	33: func(c []interface{}) interface{} {
		tok := c[0].(lexer.Token)
		return &ASTNode{Type: "id", Value: tok.Lexeme}
	},
	34: func(c []interface{}) interface{} { return c[1] },
	// Factor array access:
	35: func(children []interface{}) interface{} {
		id := children[0].(lexer.Token)
		indices := children[2].([]*ASTNode)
		return &ASTNode{Type: "arr_access", Value: id.Lexeme, Args: indices}
	},
	// Args -> NonEmptyArgs
	36: func(c []interface{}) interface{} {
		return c[0]
	},

	// Args -> ε
	37: func(c []interface{}) interface{} {
		return []*ASTNode{}
	},

	// NonEmptyArgs -> Expr
	38: func(c []interface{}) interface{} {
		if node, ok := c[0].(*ASTNode); ok {
			return []*ASTNode{node}
		}
		return []*ASTNode{}
	},

	// NonEmptyArgs -> NonEmptyArgs , Expr
	39: func(c []interface{}) interface{} {
		list, ok := c[0].([]*ASTNode)
		if !ok {
			list = []*ASTNode{}
		}
		if node, ok := c[2].(*ASTNode); ok {
			list = append(list, node)
		}
		return list
	},
	// NonEmptyArgs -> Type id
	40: func(c []interface{}) interface{} {
		typNode := c[0].(*ASTNode) // Type 节点
		idToken := c[1].(lexer.Token)
		return []*ASTNode{
			{
				Type:  "NonEmptyArg",
				Value: idToken.Lexeme,      // 使用 id 的词素作为值
				Args:  []*ASTNode{typNode}, // 将 Type 节点作为参数
			},
		}
	},
	// NonEmptyArgs -> Type id = Expr
	41: func(c []interface{}) interface{} {
		typNode := c[0].(*ASTNode) // Type 节点
		idToken := c[1].(lexer.Token)
		expr := c[3].(*ASTNode)
		return []*ASTNode{
			{
				Type:  "NonEmptyArg",
				Value: idToken.Lexeme,            // 使用 id 的词素作为值
				Args:  []*ASTNode{typNode, expr}, // 将 Type 节点和 Expr 节点作为参数
			},
		}
	},
	// NonEmptyArgs -> NonEmptyArgs , Type id
	42: func(c []interface{}) interface{} {
		list, ok := c[0].([]*ASTNode)
		if !ok {
			list = []*ASTNode{}
		}
		typNode := c[2].(*ASTNode) // Type 节点
		idToken := c[3].(lexer.Token)
		list = append(list, &ASTNode{
			Type:  "NonEmptyArg",
			Value: idToken.Lexeme,      // 使用 id 的词素作为值
			Args:  []*ASTNode{typNode}, // 将 Type 节点作为参数
		})
		return list
	},
	// NonEmptyArgs -> NonEmptyArgs , Type id = Expr
	43: func(c []interface{}) interface{} {
		list, ok := c[0].([]*ASTNode)
		if !ok {
			list = []*ASTNode{}
		}
		typNode := c[2].(*ASTNode) // Type 节点
		idToken := c[3].(lexer.Token)
		expr := c[5].(*ASTNode)
		list = append(list, &ASTNode{
			Type:  "NonEmptyArg",
			Value: idToken.Lexeme,            // 使用 id 的词素作为值
			Args:  []*ASTNode{typNode, expr}, // 将 Type 节点和 Expr 节点作为参数
		})
		return list
	},

	// IndexList -> Expr
	44: func(c []interface{}) interface{} { return []*ASTNode{c[0].(*ASTNode)} },

	// IndexList -> IndexList , Expr
	45: func(c []interface{}) interface{} {
		list := c[0].([]*ASTNode)
		return append(list, c[2].(*ASTNode))
	},

	// Cond:
	46: func(c []interface{}) interface{} {
		return &ASTNode{Type: "&&", Left: c[0].(*ASTNode), Right: c[2].(*ASTNode)}
	},
	47: func(c []interface{}) interface{} {
		return &ASTNode{Type: "||", Left: c[0].(*ASTNode), Right: c[2].(*ASTNode)}
	},
	48: func(c []interface{}) interface{} { return &ASTNode{Type: "!", Left: c[1].(*ASTNode)} },
	49: func(c []interface{}) interface{} {
		return &ASTNode{Type: "<", Left: c[0].(*ASTNode), Right: c[2].(*ASTNode)}
	},
	50: func(c []interface{}) interface{} {
		return &ASTNode{Type: ">", Left: c[0].(*ASTNode), Right: c[2].(*ASTNode)}
	},
	51: func(c []interface{}) interface{} {
		return &ASTNode{Type: "<=", Left: c[0].(*ASTNode), Right: c[2].(*ASTNode)}
	},
	52: func(c []interface{}) interface{} {
		return &ASTNode{Type: ">=", Left: c[0].(*ASTNode), Right: c[2].(*ASTNode)}
	},
	53: func(c []interface{}) interface{} {
		return &ASTNode{Type: "!=", Left: c[0].(*ASTNode), Right: c[2].(*ASTNode)}
	},
	54: func(c []interface{}) interface{} {
		return &ASTNode{Type: "==", Left: c[0].(*ASTNode), Right: c[2].(*ASTNode)}
	},
	55: func(c []interface{}) interface{} { return c[1] },
	56: func(c []interface{}) interface{} {
		return c[0].(*ASTNode)
	},
}

// var ActionFuncs = map[int]func([]interface{}) interface{}{
// 	// 0: S' -> Stmt
// 	0: func(children []interface{}) interface{} {
// 		return children[0]
// 	},

// 	// 1: Stmt -> Decl
// 	1: func(children []interface{}) interface{} {
// 		return children[0]
// 	},

// 	// 2: Decl -> int id = num ;
// 	2: func(children []interface{}) interface{} {
// 		id := children[1].(lexer.Token)
// 		num := children[3].(lexer.Token)
// 		return &ASTNode{
// 			Type: "Decl",
// 			Left: &ASTNode{
// 				Type:  "=",
// 				Left:  &ASTNode{Type: "id", Value: id.Lexeme},
// 				Right: &ASTNode{Type: "num", Value: num.Lexeme},
// 			},
// 		}
// 	},

// 	// 3: Stmt -> id = Expr ;
// 	3: func(children []interface{}) interface{} {
// 		id := children[0].(lexer.Token)
// 		expr := children[2].(*ASTNode)
// 		return &ASTNode{Type: "=", Left: &ASTNode{Type: "id", Value: id.Lexeme}, Right: expr}
// 	},

// 	// 4: Stmt -> return Expr ;
// 	4: func(children []interface{}) interface{} {
// 		expr := children[1].(*ASTNode)
// 		return &ASTNode{Type: "return", Left: expr}
// 	},

// 	// 5: Stmt -> Block
// 	5: func(children []interface{}) interface{} {
// 		return children[0]
// 	},

// 	// 6: Stmt -> if ( Cond ) Stmt
// 	6: func(children []interface{}) interface{} {
// 		cond := children[2].(*ASTNode)
// 		stmt := children[4].(*ASTNode)
// 		return &ASTNode{Type: "if", Left: cond, Right: stmt}
// 	},

// 	// 7: Stmt -> if ( Cond ) Stmt else Stmt
// 	7: func(children []interface{}) interface{} {
// 		cond := children[2].(*ASTNode)
// 		thenStmt := children[4].(*ASTNode)
// 		elseStmt := children[6].(*ASTNode)
// 		return &ASTNode{Type: "ifelse", Left: cond, Right: &ASTNode{Type: "else", Left: thenStmt, Right: elseStmt}}
// 	},

// 	// 8: Stmt -> while ( Cond ) Stmt
// 	8: func(children []interface{}) interface{} {
// 		cond := children[2].(*ASTNode)
// 		body := children[4].(*ASTNode)
// 		return &ASTNode{Type: "while", Left: cond, Right: body}
// 	},

// 	// 9: Stmt -> id [ IndexList ] = Expr ;
// 	9: func(children []interface{}) interface{} {
// 		id := children[0].(lexer.Token)
// 		indices := children[2].([]*ASTNode)
// 		expr := children[5].(*ASTNode)
// 		return &ASTNode{Type: "arr_assign", Value: id.Lexeme, Args: indices, Right: expr}
// 	},

// 	// 10: Block -> { StmtList }
// 	10: func(children []interface{}) interface{} {
// 		return &ASTNode{Type: "block", Args: children[1].([]*ASTNode)}
// 	},

// 	// 11: StmtList -> Stmt
// 	11: func(children []interface{}) interface{} {
// 		return []*ASTNode{children[0].(*ASTNode)}
// 	},
// 	// 12: StmtList -> StmtList Stmt
// 	12: func(children []interface{}) interface{} {
// 		list := children[0].([]*ASTNode)
// 		return append(list, children[1].(*ASTNode))
// 	},

// 	// Expressions:
// 	13: func(c []interface{}) interface{} {
// 		return &ASTNode{Type: "+", Left: c[0].(*ASTNode), Right: c[2].(*ASTNode)}
// 	},
// 	14: func(c []interface{}) interface{} {
// 		return &ASTNode{Type: "-", Left: c[0].(*ASTNode), Right: c[2].(*ASTNode)}
// 	},
// 	15: func(c []interface{}) interface{} { return c[0] },
// 	16: func(c []interface{}) interface{} {
// 		return &ASTNode{Type: "*", Left: c[0].(*ASTNode), Right: c[2].(*ASTNode)}
// 	},
// 	17: func(c []interface{}) interface{} {
// 		return &ASTNode{Type: "/", Left: c[0].(*ASTNode), Right: c[2].(*ASTNode)}
// 	},
// 	18: func(c []interface{}) interface{} { return c[0] },

// 	// Factor:
// 	19: func(c []interface{}) interface{} {
// 		id := c[0].(lexer.Token)
// 		args := c[2].([]*ASTNode)
// 		return &ASTNode{Type: "call", Value: id.Lexeme, Args: args}
// 	},
// 	20: func(c []interface{}) interface{} {
// 		tok := c[0].(lexer.Token)
// 		return &ASTNode{Type: "num", Value: tok.Lexeme}
// 	},
// 	21: func(c []interface{}) interface{} {
// 		tok := c[0].(lexer.Token)
// 		return &ASTNode{Type: "id", Value: tok.Lexeme}
// 	},
// 	22: func(c []interface{}) interface{} { return c[1] },

// 	// Args:
// 	23: func(c []interface{}) interface{} { return []*ASTNode{c[0].(*ASTNode)} },
// 	24: func(c []interface{}) interface{} {
// 		list := c[0].([]*ASTNode)
// 		return append(list, c[2].(*ASTNode))
// 	},

// 	// IndexList:
// 	25: func(c []interface{}) interface{} { return []*ASTNode{c[0].(*ASTNode)} },
// 	26: func(c []interface{}) interface{} {
// 		list := c[0].([]*ASTNode)
// 		return append(list, c[2].(*ASTNode))
// 	},

// 	// Cond:
// 	27: func(c []interface{}) interface{} {
// 		return &ASTNode{Type: "&&", Left: c[0].(*ASTNode), Right: c[2].(*ASTNode)}
// 	},
// 	28: func(c []interface{}) interface{} {
// 		return &ASTNode{Type: "||", Left: c[0].(*ASTNode), Right: c[2].(*ASTNode)}
// 	},
// 	29: func(c []interface{}) interface{} { return &ASTNode{Type: "!", Left: c[1].(*ASTNode)} },
// 	30: func(c []interface{}) interface{} {
// 		return &ASTNode{Type: "<", Left: c[0].(*ASTNode), Right: c[2].(*ASTNode)}
// 	},
// 	31: func(c []interface{}) interface{} {
// 		return &ASTNode{Type: "==", Left: c[0].(*ASTNode), Right: c[2].(*ASTNode)}
// 	},
// 	32: func(c []interface{}) interface{} { return c[1] },
// }
