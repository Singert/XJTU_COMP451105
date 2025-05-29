package parser

import (
	"fmt"
	"lab5/lexer"
	"lab5/semantic"
	"lab5/syntax"
	"os"
	"strings"
)

// 表示每一个分析步骤
type ParseStep struct {
	ID          int
	StateStack  []int
	SymbolStack []syntax.Symbol
	Input       syntax.Symbol
	Action      string
}

// 主分析函数
func Run(input []syntax.Symbol, g *syntax.Grammar, dfa *DFA, table *ParseTable, tokenStream []lexer.Token, verbose bool) *ParseError {
	stateStack := []int{0}
	symbolStack := []syntax.Symbol{"#"}
	attrStack := []interface{}{"#"}
	steps := []ParseStep{}
	input = append(input, "#")

	i := 0      // 当前输入符号指针
	tokIdx := 0 // 当前 tokenStream 的指针（用于语义动作）

	stepID := 0

	for {
		currState := stateStack[len(stateStack)-1]
		// 跳过注释Token
		for tokIdx < len(tokenStream) &&
			(tokenStream[tokIdx].Type == lexer.TokenCOMMENT_SINGLE || tokenStream[tokIdx].Type == lexer.TokenCOMMENT_MULTI) {
			i++
			tokIdx++
		}

		// 如果已读到输入末尾
		if i >= len(input) {
			// 检查当前状态对终结符 '#' 是否有动作
			action, ok := table.Action[currState]["#"]
			if ok && action.Typ == Accept {
				// 接受状态，语法分析成功
				fmt.Printf("状态栈: %v\t符号栈: %v\t当前输入: %s\t动作: 接受 ✅\n", stateStack, symbolStack, "#")
				err := ExportParseFlowDOT(steps, "parse_flow.dot")
				if err == nil {
					fmt.Println("✔ 分析流程图已导出为 parse_flow.dot（可用 dot -Tpng 查看）")
				}

				root := attrStack[len(attrStack)-1]
				fmt.Println("======= 抽象语法树 AST =======")
				semantic.PrintASTPretty(root.(*semantic.ASTNode), "", true)
				return nil
			} else {
				// 没有接受动作，报错
				fmt.Printf("状态栈: %v\t符号栈: %v\t当前输入: EOF\t动作: ERROR\n", stateStack, symbolStack)
				return CatchParseError(currState, "#", tokenStream, tokIdx, table)
			}
		}
		currToken := input[i]
		action, ok := table.Action[currState][currToken]

		if !ok {
			fmt.Printf("状态栈: %v\t符号栈: %v\t当前输入: %s\t动作: ERROR\n", stateStack, symbolStack, currToken)
			fmt.Println("111")
			return CatchParseError(currState, currToken, tokenStream, tokIdx, table)
		}

		var actionStr string
		switch action.Typ {
		case Shift:
			actionStr = fmt.Sprintf("shift → 状态 %d", action.Value)
		case Reduce:
			prod := g.Productions[action.Value]
			actionStr = fmt.Sprintf("reduce → %s → %s", prod.Left, joinSymbols(prod.Right))
		case Accept:
			actionStr = "接受 ✅"
		}

		steps = append(steps, ParseStep{
			ID:          stepID,
			StateStack:  append([]int(nil), stateStack...),
			SymbolStack: append([]syntax.Symbol(nil), symbolStack...),
			Input:       currToken,
			Action:      actionStr,
		})
		if verbose {
			fmt.Printf("状态栈: %v\t符号栈: %v\t当前输入: %s\t动作: %s\n",
				stateStack, symbolStack, currToken, actionStr)
		}
		stepID++

		switch action.Typ {
		case Shift:
			stateStack = append(stateStack, action.Value)
			symbolStack = append(symbolStack, currToken)

			// 从 tokenStream 中提取 token 作为属性
			if tokIdx < len(tokenStream) {
				attrStack = append(attrStack, tokenStream[tokIdx])
				tokIdx++
			} else {
				attrStack = append(attrStack, nil) // 安全兜底
			}
			i++

		case Reduce:
			prod := g.Productions[action.Value]
			rhsLen := len(prod.Right)
			if rhsLen > len(symbolStack) {
				fmt.Println("❌ 归约失败：符号栈不足")
				fmt.Println("222")
				return CatchParseError(currState, currToken, tokenStream, tokIdx, table)
			}
			stateStack = stateStack[:len(stateStack)-rhsLen]
			symbolStack = symbolStack[:len(symbolStack)-rhsLen]
			symbolStack = append(symbolStack, prod.Left)

			top := stateStack[len(stateStack)-1]
			newState, ok := table.Goto[top][prod.Left]
			if !ok {
				fmt.Println("❌ GOTO失败")
				fmt.Println("333")
				return CatchParseError(currState, currToken, tokenStream, tokIdx, table)
			}
			stateStack = append(stateStack, newState)

			// 🔧 语义动作：取出 RHS 属性，执行 action
			children := attrStack[len(attrStack)-rhsLen:]
			attrStack = attrStack[:len(attrStack)-rhsLen]
			actionFunc, exists := semantic.ActionFuncs[action.Value]
			if !exists {
				fmt.Printf("⚠ 未定义语义动作: 产生式编号 %d\n", action.Value)
				attrStack = append(attrStack, nil)
			} else {
				result := actionFunc(children)
				attrStack = append(attrStack, result)
			}

		case Accept:
			fmt.Printf("状态栈: %v\t符号栈: %v\t当前输入: %s\t动作: 接受 ✅\n", stateStack, symbolStack, currToken)
			err := ExportParseFlowDOT(steps, "parse_flow.dot")
			if err == nil {
				fmt.Println("✔ 分析流程图已导出为 parse_flow.dot（可用 dot -Tpng 查看）")
			}

			root := attrStack[len(attrStack)-1]
			fmt.Println("======= 抽象语法树 AST =======")
			semantic.PrintASTPretty(root.(*semantic.ASTNode), "", true)
			// semantic.PrintAST(root.(*semantic.ASTNode), 0)
			return nil
		default:
			fmt.Println("动作: ERROR")
			fmt.Println("555")
			return CatchParseError(currState, currToken, tokenStream, tokIdx, table)
		}
	}
}

// 辅助函数getAction
func getAction(table *ParseTable, state int, token syntax.Symbol) (Action, bool) {
	actions, ok := table.Action[state]
	if !ok {
		return Action{}, false
	}
	action, ok := actions[token]
	return action, ok
}

// 辅助符号拼接
func joinSymbols(syms []syntax.Symbol) string {
	if len(syms) == 0 {
		return "ε"
	}
	res := ""
	for _, s := range syms {
		res += string(s) + " "
	}
	return res
}

// 输出 DOT 图
func ExportParseFlowDOT(steps []ParseStep, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// 使用 top-to-bottom 排列
	fmt.Fprintln(file, "digraph ParseFlow {")
	fmt.Fprintln(file, `  rankdir=TB;`)
	fmt.Fprintln(file, `  node [shape=box, fontname="monospace", fontsize=10];`)
	fmt.Fprintln(file, `  edge [fontname="monospace"];`)

	for _, step := range steps {
		// 为不同类型动作设置不同颜色
		var color string
		switch {
		case strings.HasPrefix(step.Action, "shift"):
			color = "lightblue"
		case strings.HasPrefix(step.Action, "reduce"):
			color = "palegreen"
		case strings.HasPrefix(step.Action, "接受"):
			color = "gold"
		default:
			color = "white"
		}

		label := fmt.Sprintf(
			"Step %d\\n栈: %v\\n符号: %v\\n输入: %s\\n动作: %s",
			step.ID, step.StateStack, step.SymbolStack, step.Input, step.Action,
		)

		fmt.Fprintf(file,
			`  step%d [label="%s", style=filled, fillcolor=%s];`+"\n",
			step.ID, label, color,
		)
	}

	// 连边
	for i := 0; i < len(steps)-1; i++ {
		fmt.Fprintf(file, "  step%d -> step%d;\n", steps[i].ID, steps[i+1].ID)
	}

	fmt.Fprintln(file, "}")
	return nil
}
