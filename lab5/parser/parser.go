package parser

import (
	"fmt"
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
func Run(input []syntax.Symbol, g *syntax.Grammar, dfa *DFA, table *ParseTable) {
	stateStack := []int{0}
	symbolStack := []syntax.Symbol{"#"}
	steps := []ParseStep{}
	input = append(input, "#") // 加结束符

	i := 0
	stepID := 0

	for {
		currState := stateStack[len(stateStack)-1]
		currToken := input[i]
		action, ok := table.Action[currState][currToken]

		if !ok {
			fmt.Printf("状态栈: %v\t符号栈: %v\t当前输入: %s\t动作: ERROR\n", stateStack, symbolStack, currToken)
			break
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

		// 记录当前步骤
		steps = append(steps, ParseStep{
			ID:          stepID,
			StateStack:  append([]int(nil), stateStack...),
			SymbolStack: append([]syntax.Symbol(nil), symbolStack...),
			Input:       currToken,
			Action:      actionStr,
		})
		// 同步输出分析过程到终端
		fmt.Printf("状态栈: %v\t符号栈: %v\t当前输入: %s\t动作: %s\n",
			stateStack, symbolStack, currToken, actionStr)
		stepID++

		// 执行动作
		switch action.Typ {
		case Shift:
			stateStack = append(stateStack, action.Value)
			symbolStack = append(symbolStack, currToken)
			i++

		case Reduce:
			prod := g.Productions[action.Value]
			rhsLen := len(prod.Right)
			if rhsLen > len(symbolStack) {
				fmt.Println("❌ 归约失败：符号栈不足")
				return
			}
			stateStack = stateStack[:len(stateStack)-rhsLen]
			symbolStack = symbolStack[:len(symbolStack)-rhsLen]
			symbolStack = append(symbolStack, prod.Left)

			top := stateStack[len(stateStack)-1]
			newState, ok := table.Goto[top][prod.Left]
			if !ok {
				fmt.Println("❌ GOTO失败")
				return
			}
			stateStack = append(stateStack, newState)

		case Accept:
			fmt.Printf("状态栈: %v\t符号栈: %v\t当前输入: %s\t动作: 接受 ✅\n", stateStack, symbolStack, currToken)
			err := ExportParseFlowDOT(steps, "parse_flow.dot")
			if err == nil {
				fmt.Println("✔ 分析流程图已导出为 parse_flow.dot（可用 dot -Tpng 查看）")
			}

			return

		default:
			fmt.Println("动作: ERROR")
			return
		}
	}
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

// func ExportParseFlowDOT(steps []ParseStep, filename string) error {
// 	file, err := os.Create(filename)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	fmt.Fprintln(file, "digraph ParseFlow {")
// 	fmt.Fprintln(file, `  rankdir=LR;`)
// 	fmt.Fprintln(file, `  node [shape=box, fontname="monospace"];`)

// 	for _, step := range steps {
// 		label := fmt.Sprintf("栈: %v\\n符号: %v\\n输入: %s\\n动作: %s",
// 			step.StateStack, step.SymbolStack, step.Input, step.Action)
// 		fmt.Fprintf(file, `  step%d [label="%s"];`+"\n", step.ID, label)
// 	}

// 	for i := 0; i < len(steps)-1; i++ {
// 		fmt.Fprintf(file, "  step%d -> step%d;\n", steps[i].ID, steps[i+1].ID)
// 	}

// 	fmt.Fprintln(file, "}")
// 	return nil
// }
