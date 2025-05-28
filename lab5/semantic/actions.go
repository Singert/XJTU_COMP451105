package semantic

import (
	"lab5/lexer"
)

// ActionFuncs maps production rule indices to semantic actions.
// Each function takes the list of attributes (interface{}) popped from the attribute stack
// during a reduce action, and returns a new ASTNode (or any semantic value).
// ActionFuncs maps production index to semantic action function.

var ActionFuncs = map[int]func([]interface{}) interface{}{
	// 0: S' -> Stmt
	0: func(children []interface{}) interface{} {
		return children[0]
	},

	// 1: Stmt -> Decl
	1: func(children []interface{}) interface{} {
		return children[0]
	},

	// 2: Decl -> int id = num ;
	2: func(children []interface{}) interface{} {
		id := children[1].(lexer.Token)
		num := children[3].(lexer.Token)
		return &ASTNode{
			Type: "Decl",
			Left: &ASTNode{
				Type:  "=",
				Left:  &ASTNode{Type: "id", Value: id.Lexeme},
				Right: &ASTNode{Type: "num", Value: num.Lexeme},
			},
		}
	},

	// 3: Stmt -> id = Expr ;
	3: func(children []interface{}) interface{} {
		id := children[0].(lexer.Token)
		expr := children[2].(*ASTNode)
		return &ASTNode{Type: "=", Left: &ASTNode{Type: "id", Value: id.Lexeme}, Right: expr}
	},

	// 4: Stmt -> return Expr ;
	4: func(children []interface{}) interface{} {
		expr := children[1].(*ASTNode)
		return &ASTNode{Type: "return", Left: expr}
	},

	// 5: Stmt -> Block
	5: func(children []interface{}) interface{} {
		return children[0]
	},

	// 6: Stmt -> if ( Cond ) Stmt
	6: func(children []interface{}) interface{} {
		cond := children[2].(*ASTNode)
		stmt := children[4].(*ASTNode)
		return &ASTNode{Type: "if", Left: cond, Right: stmt}
	},

	// 7: Stmt -> if ( Cond ) Stmt else Stmt
	7: func(children []interface{}) interface{} {
		cond := children[2].(*ASTNode)
		thenStmt := children[4].(*ASTNode)
		elseStmt := children[6].(*ASTNode)
		return &ASTNode{Type: "ifelse", Left: cond, Right: &ASTNode{Type: "else", Left: thenStmt, Right: elseStmt}}
	},

	// 8: Stmt -> while ( Cond ) Stmt
	8: func(children []interface{}) interface{} {
		cond := children[2].(*ASTNode)
		body := children[4].(*ASTNode)
		return &ASTNode{Type: "while", Left: cond, Right: body}
	},

	// 9: Stmt -> id [ IndexList ] = Expr ;
	9: func(children []interface{}) interface{} {
		id := children[0].(lexer.Token)
		indices := children[2].([]*ASTNode)
		expr := children[5].(*ASTNode)
		return &ASTNode{Type: "arr_assign", Value: id.Lexeme, Args: indices, Right: expr}
	},
	// 10: Stmt -> id ( Args ) ;
	10: func(children []interface{}) interface{} {
		id := children[0].(lexer.Token)
		args := children[2].([]*ASTNode)
		return &ASTNode{
			Type:  "call_stmt",
			Value: id.Lexeme,
			Args:  args,
		}
	},

	// 10: Block -> { StmtList }
	11: func(children []interface{}) interface{} {
		return &ASTNode{Type: "block", Args: children[1].([]*ASTNode)}
	},

	// 11: StmtList -> Stmt
	12: func(children []interface{}) interface{} {
		return []*ASTNode{children[0].(*ASTNode)}
	},

	// 12: StmtList -> StmtList Stmt
	13: func(children []interface{}) interface{} {
		list := children[0].([]*ASTNode)
		return append(list, children[1].(*ASTNode))
	},

	// Expressions:
	14: func(c []interface{}) interface{} {
		return &ASTNode{Type: "+", Left: c[0].(*ASTNode), Right: c[2].(*ASTNode)}
	},
	15: func(c []interface{}) interface{} {
		return &ASTNode{Type: "-", Left: c[0].(*ASTNode), Right: c[2].(*ASTNode)}
	},
	16: func(c []interface{}) interface{} { return c[0] },
	17: func(c []interface{}) interface{} {
		return &ASTNode{Type: "*", Left: c[0].(*ASTNode), Right: c[2].(*ASTNode)}
	},
	18: func(c []interface{}) interface{} {
		return &ASTNode{Type: "/", Left: c[0].(*ASTNode), Right: c[2].(*ASTNode)}
	},
	19: func(c []interface{}) interface{} { return c[0] },

	// Factor:
	20: func(c []interface{}) interface{} {
		id := c[0].(lexer.Token)
		args := c[2].([]*ASTNode)
		return &ASTNode{Type: "call", Value: id.Lexeme, Args: args}
	},
	21: func(c []interface{}) interface{} {
		tok := c[0].(lexer.Token)
		return &ASTNode{Type: "num", Value: tok.Lexeme}
	},
	22: func(c []interface{}) interface{} {
		tok := c[0].(lexer.Token)
		return &ASTNode{Type: "id", Value: tok.Lexeme}
	},
	23: func(c []interface{}) interface{} { return c[1] },

	// Args -> NonEmptyArgs
	24: func(c []interface{}) interface{} {
		return c[0]
	},

	// Args -> ε
	25: func(c []interface{}) interface{} {
		return []*ASTNode{}
	},

	// NonEmptyArgs -> Expr
	26: func(c []interface{}) interface{} {
		if node, ok := c[0].(*ASTNode); ok {
			return []*ASTNode{node}
		}
		return []*ASTNode{}
	},

	// NonEmptyArgs -> NonEmptyArgs , Expr
	27: func(c []interface{}) interface{} {
		list, ok := c[0].([]*ASTNode)
		if !ok {
			list = []*ASTNode{}
		}
		if node, ok := c[2].(*ASTNode); ok {
			list = append(list, node)
		}
		return list
	},

	// IndexList -> Expr
	28: func(c []interface{}) interface{} { return []*ASTNode{c[0].(*ASTNode)} },

	// IndexList -> IndexList , Expr
	29: func(c []interface{}) interface{} {
		list := c[0].([]*ASTNode)
		return append(list, c[2].(*ASTNode))
	},

	// Cond:
	30: func(c []interface{}) interface{} {
		return &ASTNode{Type: "&&", Left: c[0].(*ASTNode), Right: c[2].(*ASTNode)}
	},
	31: func(c []interface{}) interface{} {
		return &ASTNode{Type: "||", Left: c[0].(*ASTNode), Right: c[2].(*ASTNode)}
	},
	32: func(c []interface{}) interface{} { return &ASTNode{Type: "!", Left: c[1].(*ASTNode)} },
	33: func(c []interface{}) interface{} {
		return &ASTNode{Type: "<", Left: c[0].(*ASTNode), Right: c[2].(*ASTNode)}
	},
	34: func(c []interface{}) interface{} {
		return &ASTNode{Type: "==", Left: c[0].(*ASTNode), Right: c[2].(*ASTNode)}
	},
	35: func(c []interface{}) interface{} { return c[1] },
	36: func(c []interface{}) interface{} {
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
