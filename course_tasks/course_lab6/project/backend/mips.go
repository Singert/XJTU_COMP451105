package backend

import (
	"fmt"
	"strings"
)

type Function struct {
	Name      string
	TAC       []string
	Locals    map[string]int
	StackSize int
}

func GroupFunctions(tac []string) []Function {
	var functions []Function
	var cur *Function
	for _, line := range tac {
		if strings.HasPrefix(line, "LABEL FUNC_") {
			name := strings.TrimPrefix(line, "LABEL FUNC_")
			cur = &Function{
				Name:   name,
				TAC:    []string{},
				Locals: map[string]int{},
			}
			functions = append(functions, *cur)
		} else if strings.HasPrefix(line, "ENDFUNC") {
			// ok
		} else {
			functions[len(functions)-1].TAC = append(functions[len(functions)-1].TAC, line)
		}
	}
	return functions
}

func AllocateStackSlots(fn *Function) {
	offset := -4
	seen := map[string]bool{}
	for _, line := range fn.TAC {
		tokens := strings.Fields(line)
		for _, tok := range tokens {
			if isKeyword(tok) || isNumber(tok) || strings.HasPrefix(tok, "L") {
				continue
			}
			if !seen[tok] {
				fn.Locals[tok] = offset
				offset -= 4
				seen[tok] = true
			}
		}
	}
	fn.StackSize = -offset
}

func isKeyword(s string) bool {
	keywords := []string{"POP", "PAR", "RETURN", "CALL", "GOTO", "LABEL", "IF", "PRINT", "="}
	for _, k := range keywords {
		if s == k {
			return true
		}
	}
	return false
}

func isNumber(s string) bool {
	_, err := fmt.Sscanf(s, "%d", new(int))
	return err == nil
}

func GenerateMIPS(tac []string) []string {
	funcs := GroupFunctions(tac)
	var output []string

	for i := range funcs {
		fn := &funcs[i]
		AllocateStackSlots(fn)
		output = append(output, fmt.Sprintf("%s:", fn.Name))
		output = append(output, fmt.Sprintf("  addi $sp, $sp, -%d", fn.StackSize+8))
		output = append(output, fmt.Sprintf("  sw $ra, %d($sp)", fn.StackSize+4))
		output = append(output, fmt.Sprintf("  sw $fp, %d($sp)", fn.StackSize))
		output = append(output, fmt.Sprintf("  move $fp, $sp"))

		hasReturn := false
		for _, line := range fn.TAC {
			if strings.HasPrefix(line, "RETURN") {
				hasReturn = true
			}
			output = append(output, translateLine(fn, line)...)
		}
		if !hasReturn {
			output = append(output,
				fmt.Sprintf("  lw $ra, %d($fp)", fn.StackSize+4),
				fmt.Sprintf("  lw $fp, %d($fp)", fn.StackSize),
				fmt.Sprintf("  addi $sp, $sp, %d", fn.StackSize+8),
				"  jr $ra",
			)
		}

	}
	return output
}

func translateLine(fn *Function, line string) []string {
	tokens := strings.Fields(line)
	if len(tokens) == 0 {
		return nil
	}

	switch tokens[0] {
	case "LABEL":
		return []string{fmt.Sprintf("%s:", tokens[1])}
	case "POP":
		v := tokens[1]
		return []string{fmt.Sprintf("  lw $t0, 0($sp)"), fmt.Sprintf("  sw $t0, %d($fp)", fn.Locals[v]), "  addi $sp, $sp, 4"}
	case "PAR":
		arg := tokens[1]
		if isNumber(arg) {
			return []string{
				fmt.Sprintf("  li $a0, %s", arg),
				"  addi $sp, $sp, -4",
				"  sw $a0, 0($sp)",
			}
		}
		return []string{
			fmt.Sprintf("  lw $a0, %d($fp)", fn.Locals[arg]),
			"  addi $sp, $sp, -4",
			"  sw $a0, 0($sp)",
		}
	case "RETURN":
		v := tokens[1]
		return []string{
			fmt.Sprintf("  lw $v0, %d($fp)", fn.Locals[v]),
			fmt.Sprintf("  lw $ra, %d($fp)", fn.StackSize+4),
			fmt.Sprintf("  lw $fp, %d($fp)", fn.StackSize),
			fmt.Sprintf("  addi $sp, $sp, %d", fn.StackSize+8),
			"  jr $ra",
		}
	case "PRINT":
		v := tokens[1]
		return []string{
			fmt.Sprintf("  lw $a0, %d($fp)", fn.Locals[v]),
			"  li $v0, 1",
			"  syscall",
		}
	case "GOTO":
		return []string{fmt.Sprintf("  j %s", tokens[1])}
	case "IF":
		cond := tokens[1]
		lt := tokens[5]
		lf := tokens[7]
		offset := fn.Locals[cond]
		return []string{
			fmt.Sprintf("  lw $t0, %d($fp)", offset),
			fmt.Sprintf("  bne $t0, $zero, %s", lt),
			fmt.Sprintf("  j %s", lf),
		}
	default:
		// CALL
		if len(tokens) >= 5 && tokens[2] == "CALL" {
			dst := tokens[0]
			fnName := tokens[3]
			fnName = fnName[:len(fnName)-1]
			return []string{
				fmt.Sprintf("  jal %s", fnName),
				fmt.Sprintf("  sw $v0, %d($fp)", fn.Locals[dst]),
			}
		}
		// z = tX（CALL 结果赋值）
		if len(tokens) == 3 && tokens[1] == "=" {
			dst := tokens[0]
			src := tokens[2]

			// 检查 src 是否是某个 CALL 临时变量
			if strings.HasPrefix(src, "t") && fn.Locals[src] != 0 {
				// z = tX → load + store
				return []string{
					fmt.Sprintf("  lw $t0, %d($fp)", fn.Locals[src]),
					fmt.Sprintf("  sw $t0, %d($fp)", fn.Locals[dst]),
				}
			}
			// 其他情况保留原始赋值逻辑
			if isNumber(src) {
				return []string{
					fmt.Sprintf("  li $t0, %s", src),
					fmt.Sprintf("  sw $t0, %d($fp)", fn.Locals[dst]),
				}
			}
			return []string{
				fmt.Sprintf("  lw $t0, %d($fp)", fn.Locals[src]),
				fmt.Sprintf("  sw $t0, %d($fp)", fn.Locals[dst]),
			}
		}

		// a = b + c / a = b
		if len(tokens) == 5 && tokens[1] == "=" {
			dst, a, op, b := tokens[0], tokens[2], tokens[3], tokens[4]
			lines := []string{}
			if isNumber(a) {
				lines = append(lines, fmt.Sprintf("  li $t0, %s", a))
			} else {
				lines = append(lines, fmt.Sprintf("  lw $t0, %d($fp)", fn.Locals[a]))
			}
			if isNumber(b) {
				lines = append(lines, fmt.Sprintf("  li $t1, %s", b))
			} else {
				lines = append(lines, fmt.Sprintf("  lw $t1, %d($fp)", fn.Locals[b]))
			}
			switch op {
			case "+":
				lines = append(lines, "  add $t2, $t0, $t1")
			case "-":
				lines = append(lines, "  sub $t2, $t0, $t1")
			case "*":
				lines = append(lines, "  mul $t2, $t0, $t1")
			case "/":
				lines = append(lines, "  div $t0, $t1", "  mflo $t2")
			case "<":
				lines = append(lines, "  slt $t2, $t0, $t1")
			case ">":
				lines = append(lines, "  slt $t2, $t1, $t0")
			case "==":
				lines = append(lines, "  sub $t2, $t0, $t1", "  seq $t2, $t2, $zero")
			case "!=":
				lines = append(lines, "  sub $t2, $t0, $t1", "  sne $t2, $t2, $zero")
			default:
				return []string{fmt.Sprintf("# [未支持操作符] %s", line)}
			}
			lines = append(lines, fmt.Sprintf("  sw $t2, %d($fp)", fn.Locals[dst]))
			return lines
		}

		// a = b
		if len(tokens) == 3 && tokens[1] == "=" {
			dst, src := tokens[0], tokens[2]
			if isNumber(src) {
				return []string{
					fmt.Sprintf("  li $t0, %s", src),
					fmt.Sprintf("  sw $t0, %d($fp)", fn.Locals[dst]),
				}
			}
			return []string{
				fmt.Sprintf("  lw $t0, %d($fp)", fn.Locals[src]),
				fmt.Sprintf("  sw $t0, %d($fp)", fn.Locals[dst]),
			}
		}
	}

	return []string{fmt.Sprintf("# [未翻译] %s", line)}
}

/*
✅ 下一步（可选功能）
你可以继续扩展支持：

z = f(x) 带参数函数调用完整寄存器协议（包括 $a0~$a3）

全局变量管理

a[i] = v 等数组操作
*/
