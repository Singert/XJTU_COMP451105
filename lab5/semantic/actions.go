package semantic

import (
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
			Args: children[0].([]*ASTNode),
		}
	},

	// 2: Program -> StmtList EOF
	2: func(children []interface{}) interface{} {
		return &ASTNode{
			Type: "Program",
			Args: children[0].([]*ASTNode),
		}
	},

	// 3: Stmt -> Func
	3: func(children []interface{}) interface{} {
		return children[0]
	},

	// 4: Func -> Type id ( Args ) Block
	4: func(children []interface{}) interface{} {
		return &ASTNode{
			Type:  "Func",
			Value: children[0].(*ASTNode).Value, // return type
			Left: &ASTNode{
				Type:  "id",
				Value: children[1].(lexer.Token).Lexeme,
			},
			Args:  children[3].([]*ASTNode), // function params
			Right: children[5].(*ASTNode),   // body block
		}
	},

	// 5: Stmt -> Decl
	5: func(children []interface{}) interface{} {
		return children[0]
	},

	// 6: Decl -> Type id = Expr ;
	6: func(children []interface{}) interface{} {
		return &ASTNode{
			Type:  "Decl",
			Value: children[0].(*ASTNode).Value, // type
			Left:  &ASTNode{Type: "id", Value: children[1].(lexer.Token).Lexeme},
			Right: children[3].(*ASTNode), // RHS Expr
		}
	},

	// 7: Type -> type_kw
	7: func(children []interface{}) interface{} {
		return &ASTNode{Type: "Type", Value: children[0].(lexer.Token).Lexeme}
	},

	// 8: Stmt -> id = Expr ;
	8: func(children []interface{}) interface{} {
		return &ASTNode{
			Type:  "=",
			Left:  &ASTNode{Type: "id", Value: children[0].(lexer.Token).Lexeme},
			Right: children[2].(*ASTNode),
		}
	},

	// 9: Stmt -> return Expr ;
	9: func(children []interface{}) interface{} {
		return &ASTNode{
			Type: "return",
			Left: children[1].(*ASTNode),
		}
	},

	// 10: Stmt -> Block
	10: func(children []interface{}) interface{} {
		return children[0]
	},

	// 11: Stmt -> if ( Cond ) Stmt
	11: func(children []interface{}) interface{} {
		return &ASTNode{
			Type:  "if",
			Left:  children[2].(*ASTNode),
			Right: children[4].(*ASTNode),
		}
	},

	// 12: Stmt -> if ( Cond ) Stmt else Stmt
	12: func(children []interface{}) interface{} {
		return &ASTNode{
			Type: "ifelse",
			Left: children[2].(*ASTNode),
			Right: &ASTNode{
				Type:  "else",
				Left:  children[4].(*ASTNode),
				Right: children[6].(*ASTNode),
			},
		}
	},

	// 13: Stmt -> while ( Cond ) Stmt
	13: func(children []interface{}) interface{} {
		return &ASTNode{
			Type:  "while",
			Left:  children[2].(*ASTNode),
			Right: children[4].(*ASTNode),
		}
	},

	// 14: Stmt -> id MultiIndex = Expr ;
	14: func(children []interface{}) interface{} {
		return &ASTNode{
			Type:  "arr_assign",
			Value: children[0].(lexer.Token).Lexeme,
			Args:  children[1].([]*ASTNode), // 多维索引数组
			Right: children[3].(*ASTNode),   // 右值表达式
		}
	},

	// 15: Stmt -> id ( Args ) ;
	15: func(children []interface{}) interface{} {
		return &ASTNode{
			Type:  "call_stmt",
			Value: children[0].(lexer.Token).Lexeme,
			Args:  children[2].([]*ASTNode),
		}
	},

	// 16: Block -> { StmtList }
	16: func(children []interface{}) interface{} {
		return &ASTNode{
			Type: "block",
			Args: children[1].([]*ASTNode),
		}
	},

	// 17: StmtList ->
	17: func(children []interface{}) interface{} {
		return []*ASTNode{}
	},

	// 18: StmtList -> Stmt
	18: func(children []interface{}) interface{} {
		return []*ASTNode{children[0].(*ASTNode)}
	},

	// 19: StmtList -> StmtList Stmt
	19: func(children []interface{}) interface{} {
		return append(children[0].([]*ASTNode), children[1].(*ASTNode))
	},

	// 20: Expr -> Expr + Term
	20: func(children []interface{}) interface{} {
		return &ASTNode{
			Type:  "+",
			Left:  children[0].(*ASTNode),
			Right: children[2].(*ASTNode),
		}
	},

	// 21: Expr -> Expr - Term
	21: func(children []interface{}) interface{} {
		return &ASTNode{
			Type:  "-",
			Left:  children[0].(*ASTNode),
			Right: children[2].(*ASTNode),
		}
	},

	// 22: Expr -> Term
	22: func(children []interface{}) interface{} {
		return children[0]
	},

	// 23: Term -> Term * CastExpr
	23: func(children []interface{}) interface{} {
		return &ASTNode{
			Type:  "*",
			Left:  children[0].(*ASTNode),
			Right: children[2].(*ASTNode),
		}
	},

	// 24: Term -> Term / CastExpr
	24: func(children []interface{}) interface{} {
		return &ASTNode{
			Type:  "/",
			Left:  children[0].(*ASTNode),
			Right: children[2].(*ASTNode),
		}
	},

	// 25: Term -> CastExpr
	25: func(children []interface{}) interface{} {
		return children[0]
	},

	// 26: CastExpr -> CastPrefix Factor
	26: func(children []interface{}) interface{} {
		return &ASTNode{
			Type:  "cast",
			Value: children[0].(*ASTNode).Value,
			Right: children[1].(*ASTNode),
		}
	},

	// 27: CastExpr -> Factor
	27: func(children []interface{}) interface{} {
		return children[0]
	},

	// 28: CastPrefix -> ( Type )
	28: func(children []interface{}) interface{} {
		return &ASTNode{
			Type:  "CastPrefix",
			Value: children[1].(*ASTNode).Value,
		}
	},

	// 29: Factor -> id ( Args )
	29: func(children []interface{}) interface{} {
		return &ASTNode{
			Type:  "call",
			Value: children[0].(lexer.Token).Lexeme,
			Args:  children[2].([]*ASTNode),
		}
	},

	// 30: Factor -> num
	30: func(children []interface{}) interface{} {
		return &ASTNode{Type: "num", Value: children[0].(lexer.Token).Lexeme}
	},
	// 31: Factor -> float
	31: func(children []interface{}) interface{} {
		return &ASTNode{Type: "float", Value: children[0].(lexer.Token).Lexeme}
	},
	// 32: Factor -> char
	32: func(children []interface{}) interface{} {
		return &ASTNode{Type: "char", Value: children[0].(lexer.Token).Lexeme}
	},
	// 33: Factor -> string
	33: func(children []interface{}) interface{} {
		return &ASTNode{Type: "string", Value: children[0].(lexer.Token).Lexeme}
	},
	// 31: Factor -> id
	34: func(children []interface{}) interface{} {
		return &ASTNode{Type: "id", Value: children[0].(lexer.Token).Lexeme}
	},

	// 32: Factor -> ( Expr )
	35: func(children []interface{}) interface{} {
		return children[1]
	},

	// 33: Factor -> id MultiIndex
	36: func(children []interface{}) interface{} {
		idToken := children[0].(lexer.Token)
		multiIndices := children[1].([]*ASTNode)
		return &ASTNode{
			Type:  "array_access",
			Value: idToken.Lexeme,
			Args:  multiIndices,
		}
	},

	// 34: Args -> NonEmptyArgs
	37: func(children []interface{}) interface{} {
		return children[0].([]*ASTNode)
	},

	// 35: Args ->
	38: func(children []interface{}) interface{} {
		return []*ASTNode{}
	},

	// 36: NonEmptyArgs -> Expr
	39: func(children []interface{}) interface{} {
		return []*ASTNode{children[0].(*ASTNode)}
	},

	// 37: NonEmptyArgs -> NonEmptyArgs , Expr
	40: func(children []interface{}) interface{} {
		return append(children[0].([]*ASTNode), children[2].(*ASTNode))
	},

	// 38: NonEmptyArgs -> Type id
	41: func(children []interface{}) interface{} {
		return []*ASTNode{
			{
				Type:  "Arg",
				Value: children[0].(*ASTNode).Value,
				Left:  &ASTNode{Type: "id", Value: children[1].(lexer.Token).Lexeme},
			},
		}
	},

	// 39: NonEmptyArgs -> Type id = Expr
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

	// 40: NonEmptyArgs -> NonEmptyArgs , Type id
	43: func(children []interface{}) interface{} {
		args := children[0].([]*ASTNode)
		args = append(args, &ASTNode{
			Type:  "Arg",
			Value: children[2].(*ASTNode).Value,
			Left:  &ASTNode{Type: "id", Value: children[3].(lexer.Token).Lexeme},
		})
		return args
	},

	// 41: NonEmptyArgs -> NonEmptyArgs , Type id = Expr
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

	// 42: IndexList -> Expr
	45: func(children []interface{}) interface{} {
		return []*ASTNode{children[0].(*ASTNode)}
	},

	// 43: IndexList -> IndexList , Expr
	46: func(children []interface{}) interface{} {
		return append(children[0].([]*ASTNode), children[2].(*ASTNode))
	},

	// 44: Cond -> Cond && Cond
	47: func(children []interface{}) interface{} {
		return &ASTNode{
			Type:  "&&",
			Left:  children[0].(*ASTNode),
			Right: children[2].(*ASTNode),
		}
	},

	// 45: Cond -> Cond || Cond
	48: func(children []interface{}) interface{} {
		return &ASTNode{
			Type:  "||",
			Left:  children[0].(*ASTNode),
			Right: children[2].(*ASTNode),
		}
	},

	// 46: Cond -> ! Cond
	49: func(children []interface{}) interface{} {
		return &ASTNode{
			Type:  "!",
			Right: children[1].(*ASTNode),
		}
	},

	// 47: Cond -> Expr < Expr
	50: func(children []interface{}) interface{} {
		return &ASTNode{Type: "<", Left: children[0].(*ASTNode), Right: children[2].(*ASTNode)}
	},

	// 48: Cond -> Expr > Expr
	51: func(children []interface{}) interface{} {
		return &ASTNode{Type: ">", Left: children[0].(*ASTNode), Right: children[2].(*ASTNode)}
	},

	// 49: Cond -> Expr <= Expr
	52: func(children []interface{}) interface{} {
		return &ASTNode{Type: "<=", Left: children[0].(*ASTNode), Right: children[2].(*ASTNode)}
	},

	// 50: Cond -> Expr >= Expr
	53: func(children []interface{}) interface{} {
		return &ASTNode{Type: ">=", Left: children[0].(*ASTNode), Right: children[2].(*ASTNode)}
	},

	// 51: Cond -> Expr != Expr
	54: func(children []interface{}) interface{} {
		return &ASTNode{Type: "!=", Left: children[0].(*ASTNode), Right: children[2].(*ASTNode)}
	},

	// 52: Cond -> Expr == Expr
	55: func(children []interface{}) interface{} {
		return &ASTNode{Type: "==", Left: children[0].(*ASTNode), Right: children[2].(*ASTNode)}
	},

	// 53: Cond -> ( Cond )
	56: func(children []interface{}) interface{} {
		return children[1]
	},

	// 54: Cond -> Expr
	57: func(children []interface{}) interface{} {
		return children[0]
	},
	// // 支持数组声明：Type id [ IndexList ] ;
	// 55: func(children []interface{}) interface{} {
	// 	typNode := children[0].(*ASTNode)
	// 	idToken := children[1].(lexer.Token)
	// 	indices := children[3].([]*ASTNode)
	// 	return &ASTNode{
	// 		Type:  "ArrayDecl",
	// 		Value: typNode.Value,
	// 		Left:  &ASTNode{Type: "id", Value: idToken.Lexeme},
	// 		Args:  indices,
	// 	}
	// },

	// // 支持数组声明带初始化：Type id [ IndexList ] = Expr ;
	// 56: func(children []interface{}) interface{} {
	// 	typNode := children[0].(*ASTNode)
	// 	idToken := children[1].(lexer.Token)
	// 	indices := children[3].([]*ASTNode)
	// 	expr := children[6].(*ASTNode)
	// 	return &ASTNode{
	// 		Type:  "ArrayDeclInit",
	// 		Value: typNode.Value,
	// 		Left:  &ASTNode{Type: "id", Value: idToken.Lexeme},
	// 		Args:  indices,
	// 		Right: expr,
	// 	}
	// },
	// 无初始化数组声明
	58: func(children []interface{}) interface{} {
		typNode := children[0].(*ASTNode)
		idToken := children[1].(lexer.Token)
		multiIndices := children[2].([]*ASTNode)
		return &ASTNode{
			Type:  "ArrayDeclMulti",
			Value: typNode.Value,
			Left:  &ASTNode{Type: "id", Value: idToken.Lexeme},
			Args:  multiIndices,
		}
	},

	// 带初始化数组声明（初始化列表或字符串字面量）
	59: func(children []interface{}) interface{} {
		typNode := children[0].(*ASTNode)
		idToken := children[1].(lexer.Token)
		multiIndices := children[2].([]*ASTNode)
		initExpr := children[4].(*ASTNode)
		return &ASTNode{
			Type:  "ArrayDeclMultiInit",
			Value: typNode.Value,
			Left:  &ASTNode{Type: "id", Value: idToken.Lexeme},
			Args:  multiIndices,
			Right: initExpr,
		}
	},
	// 59: MultiIndex -> [ IndexList ] MultiIndex
	60: func(children []interface{}) interface{} {
		indexList := children[1].([]*ASTNode)
		multiIndexNext := children[3].([]*ASTNode)
		return append(indexList, multiIndexNext...)
	},
	// 60: MultiIndex -> ε
	61: func(children []interface{}) interface{} {
		return []*ASTNode{}
	},

	// Decl -> Type id MultiIndex = InitList ;
	62: func(children []interface{}) interface{} {
		typNode := children[0].(*ASTNode)
		idToken := children[1].(lexer.Token)
		multiIndices := children[2].([]*ASTNode)
		initList := children[4].(*ASTNode)
		return &ASTNode{
			Type:  "ArrayDeclMultiInitList",
			Value: typNode.Value,
			Left:  &ASTNode{Type: "id", Value: idToken.Lexeme},
			Args:  multiIndices,
			Right: initList,
		}
	},
	// InitList -> { NonEmptyInitList }
	63: func(children []interface{}) interface{} {
		return &ASTNode{
			Type: "InitList",
			Args: children[1].([]*ASTNode),
		}
	},
	// InitList -> {}
	64: func(children []interface{}) interface{} {
		return &ASTNode{
			Type: "InitList",
			Args: []*ASTNode{},
		}
	},

	// NonEmptyInitList -> Expr
	65: func(children []interface{}) interface{} {
		return []*ASTNode{children[0].(*ASTNode)}
	},

	// NonEmptyInitList -> NonEmptyInitList , Expr
	66: func(children []interface{}) interface{} {
		return append(children[0].([]*ASTNode), children[2].(*ASTNode))
	},
	// Expr -> InitList
	67: func(children []interface{}) interface{} {
		return children[0]
	},
}
