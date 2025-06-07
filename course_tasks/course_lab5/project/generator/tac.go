// generator/tac.go
package generator

import "fmt"

var tempCounter = 0
var labelCounter = 0

func newTemp() string {
	tempCounter++
	return fmt.Sprintf("t%d", tempCounter)
}

func newLabel() string {
	labelCounter++
	return fmt.Sprintf("L%d", labelCounter)
}

func GenerateExampleArrayAssignment() []string {
	code := []string{}
	t1 := newTemp()
	code = append(code, fmt.Sprintf("%s = i + 1", t1))
	t2 := newTemp()
	code = append(code, fmt.Sprintf("%s = j * 2", t2))
	t3 := newTemp()
	code = append(code, fmt.Sprintf("%s = 4", t3))
	t4 := newTemp()
	code = append(code, fmt.Sprintf("%s = %s * 5", t4, t1))
	t5 := newTemp()
	code = append(code, fmt.Sprintf("%s = %s + %s", t5, t4, t2))
	t6 := newTemp()
	code = append(code, fmt.Sprintf("%s = %s * 20", t6, t5))
	t7 := newTemp()
	code = append(code, fmt.Sprintf("%s = %s + %s", t7, t6, t3))
	t8 := newTemp()
	code = append(code, fmt.Sprintf("%s = %s * 4", t8, t7))
	t9 := newTemp()
	code = append(code, fmt.Sprintf("%s = 66", t9))
	code = append(code, fmt.Sprintf("a[%s] = %s", t8, t9))
	return code
}

func GenerateIfStatement() []string {
	code := []string{}
	t1 := newTemp()
	code = append(code, fmt.Sprintf("%s = 1", t1))
	l1 := newLabel()
	l2 := newLabel()
	code = append(code, fmt.Sprintf("IF x < %s THEN %s ELSE %s", t1, l1, l2))
	code = append(code, fmt.Sprintf("LABEL %s", l1))
	code = append(code, "y = 1")
	code = append(code, fmt.Sprintf("GOTO %s", newLabel()))
	code = append(code, fmt.Sprintf("LABEL %s", l2))
	code = append(code, "y = k")
	return code
}

func GenerateWhileStatement() []string {
	code := []string{}
	entry := newLabel()
	trueLabel := newLabel()
	falseLabel := newLabel()
	code = append(code, fmt.Sprintf("LABEL %s", entry))
	code = append(code, fmt.Sprintf("IF z != 0 THEN %s ELSE %s", trueLabel, falseLabel))
	code = append(code, fmt.Sprintf("LABEL %s", trueLabel))
	code = append(code, "x = 1")
	code = append(code, "y = b")
	code = append(code, fmt.Sprintf("GOTO %s", entry))
	code = append(code, fmt.Sprintf("LABEL %s", falseLabel))
	return code
}
