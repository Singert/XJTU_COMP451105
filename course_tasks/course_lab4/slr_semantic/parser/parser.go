package parser

import (
	"fmt"
	"slr_semantic/lexer"
	"slr_semantic/semantic"
	"strconv"
)

func RunParser(tokens []string, table ParsingTable, grammar Grammar, symtab *semantic.SymbolTable, values []lexer.Token) {
	tokens = append(tokens, "$") // 终止符
	stateStack := []int{0}
	symbolStack := []string{}

	attrStack := []interface{}{} // 属性栈，对应每个符号

	i := 0 // token 位置

	for {
		state := stateStack[len(stateStack)-1]
		lookahead := tokens[i]

		action, ok := table.Action[state][lookahead]
		if !ok {
			panic(fmt.Sprintf("语法错误：在状态 %d 遇到符号 '%s'", state, lookahead))
		}

		switch action.Action {
		case "shift":
			stateStack = append(stateStack, action.Value)
			symbolStack = append(symbolStack, lookahead)
			attrStack = append(attrStack, values[i].Value) // 使用真实值
			i++
		case "reduce":
			prod := grammar.Productions[action.Value]
			n := len(prod.Right)

			// 弹出栈
			stateStack = stateStack[:len(stateStack)-n]
			symbolStack = symbolStack[:len(symbolStack)-n]
			children := attrStack[len(attrStack)-n:]
			attrStack = attrStack[:len(attrStack)-n]

			// 语义动作
			node := applySemanticAction(prod.Left, prod.Right, children, symtab)

			// 压入非终结符
			symbolStack = append(symbolStack, prod.Left)
			attrStack = append(attrStack, node)
			top := stateStack[len(stateStack)-1]
			nextState := table.Goto[top][prod.Left]
			stateStack = append(stateStack, nextState)

		case "accept":
			fmt.Println("✅ 分析成功！")
			return
		}
	}
}
func applySemanticAction(lhs string, rhs []string, children []interface{}, symtab *semantic.SymbolTable) interface{} {
	switch lhs {
	case "Decl":
		if len(rhs) == 3 && rhs[2] == ";" {
			// T id ;
			typ := children[0].(string)
			name := children[1].(string)
			symtab.Add(semantic.Symbol{
				Name: name, Type: typ, Kind: "var",
			})
		} else if len(rhs) == 6 && rhs[2] == "[" {
			// T id [ num ] ;
			typ := children[0].(string)
			name := children[1].(string)
			size := children[3].(string)
			symtab.Add(semantic.Symbol{
				Name: name, Type: typ, Kind: "array", Dim: []int{atoi(size)},
			})
		} else if len(rhs) == 6 && rhs[2] == "(" {
			// T id ( ParamList ) ;
			typ := children[0].(string)
			name := children[1].(string)
			params := children[3].([]semantic.Symbol)
			symtab.Add(semantic.Symbol{
				Name: name, Type: typ, Kind: "function", Params: params,
			})
		}
		return nil
	case "T":
		return children[0].(string)
	case "Param":
		return semantic.Symbol{
			Name: children[1].(string),
			Type: children[0].(string),
			Kind: "var",
		}
	case "ParamList":
		if len(rhs) == 1 {
			return []semantic.Symbol{children[0].(semantic.Symbol)}
		} else {
			return append(children[0].([]semantic.Symbol), children[2].(semantic.Symbol))
		}
	case "DeclList":
		return nil
	default:
		if len(children) > 0 {
			return children[0] // 透传
		}
	}
	return nil
}

func atoi(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}
