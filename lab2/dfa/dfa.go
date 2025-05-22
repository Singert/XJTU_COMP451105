package dfa

import (
	"fmt"
	"os"
)

type DFA struct {
	Alphabet     []string                     `json:"alphabet"`
	States       []string                     `json:"states"`
	StartState   string                       `json:"start_state"`
	AcceptStates []string                     `json:"accept_states"`
	Transitions  map[string]map[string]string `json:"transitions"`

	//运行时变量
	acceptMap map[string]bool
}

type TransitionTrace struct {
	From   string
	Symbol string
	To     string
}

func (d *DFA) buildAcceptMap() {
	d.acceptMap = make(map[string]bool)
	for _, s := range d.AcceptStates {
		d.acceptMap[s] = true
	}
}

func (d *DFA) ExportDFAtoDot(filename string) error {
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    fmt.Fprintln(file, "digraph DFA {")
    fmt.Fprintln(file, "  rankdir=LR;")

    // 初始状态箭头（空节点到起始状态）
    fmt.Fprintf(file, `  "" -> %s;`+"\n", d.StartState)

    // 输出所有状态节点
    for _, s := range d.States {
        shape := "circle"
        if d.acceptMap[s] {
            shape = "doublecircle"
        }
        fmt.Fprintf(file, `  %s [shape=%s];`+"\n", s, shape)
    }

    // 输出转移边
    for from, transitions := range d.Transitions {
        for symbol, to := range transitions {
            fmt.Fprintf(file, `  %s -> %s [label="%s"];`+"\n", from, to, symbol)
        }
    }

    fmt.Fprintln(file, "}")
    return nil
}


func (d *DFA) MatchDFA(input string) (bool, []TransitionTrace) {
	currentState := d.StartState
	trace := []TransitionTrace{}
	for i, ch := range input {
		symbol := string(ch)
		if !contains(d.Alphabet, symbol) {
			fmt.Printf("Step %d: %s --%s--> ❌invalid symbol\n", i+1, currentState, symbol)
			return false, nil
		}
		nextState, ok := d.Transitions[currentState][symbol]
		if !ok {
			fmt.Printf("Step %d: %s --%s--> ❌invalid transition\n", i+1, currentState, symbol)
			return false, nil
		}

		trace = append(trace, TransitionTrace{
			From:   currentState,
			Symbol: symbol,
			To:     nextState,
		})

		fmt.Printf("Step %d: %s --%s--> %s\n", i+1, currentState, symbol, nextState)
		currentState = nextState
	}
	return d.acceptMap[currentState], trace
}

func (d *DFA) EnumValidStrings(maxLength int) []string {
	var validStrings []string
	var currentString string
	var currentState string

	// BFS queue
	queue := make([]struct {
		state  string
		string string
		length int
	}, 0)

	// Start with the initial state and an empty string
	queue = append(queue, struct {
		state  string
		string string
		length int
	}{d.StartState, "", 0})

	for len(queue) > 0 {
		item := queue[0]
		queue = queue[1:]

		currentState = item.state
		currentString = item.string

		if item.length > maxLength {
			continue
		}

		if d.acceptMap[currentState] && item.length <= maxLength {
			validStrings = append(validStrings, currentString)
		}

		for _, symbol := range d.Alphabet {
			nextState, ok := d.Transitions[currentState][symbol]
			if ok {
				queue = append(queue, struct {
					state  string
					string string
					length int
				}{nextState, currentString + symbol, item.length + 1})
			}
		}
	}
	return validStrings
}

func (d *DFA) CheckValidity() bool {
	fmt.Printf("[Checking DFA Validity]\n")

	// 1️⃣ check start_state not nil
	fmt.Printf("checking start_state not nil: ")
	if d.StartState == "" {
		fmt.Print("❌ Start state is nil\n")
		fmt.Println("[DFA Invalid]")
		return false
	} else {
		fmt.Println("PASS")
	}

	// 2️⃣ check start_state in states
	fmt.Print("checking start_state in states: ")
	if !contains(d.States, d.StartState) {
		fmt.Printf("❌ Start state %s not in states\n", d.StartState)
		fmt.Println("[DFA Invalid]")
		return false
	} else {
		fmt.Println("PASS")
	}

	// 3️⃣ check accept_states not nil
	fmt.Print("checking accept_states not nil: ")
	if len(d.AcceptStates) == 0 {
		fmt.Print("❌ Accept states is nil\n")
		fmt.Println("[DFA Invalid]")
		return false
	} else {
		fmt.Println("PASS")
	}

	// 4️⃣ check accept_states in states
	fmt.Print("checking accept_states in states: ")
	for _, s := range d.AcceptStates {
		if !contains(d.States, s) {
			fmt.Printf("❌ Accept state %s not in states\n", s)
			fmt.Println("[DFA Invalid]")
			return false
		}
	}
	fmt.Println("PASS")

	// 5️⃣ check transitions not nil
	fmt.Print("checking transitions not nil: ")
	if len(d.Transitions) == 0 {
		fmt.Print("❌ Transitions is nil\n")
		fmt.Println("[DFA Invalid]")
		return false
	} else {
		fmt.Println("PASS")
	}

	// 6️⃣ check all transitions refer to valid states and symbols
	fmt.Print("checking transitions for valid states and symbols: ")
	for from, trans := range d.Transitions {
		if !contains(d.States, from) {
			fmt.Printf("❌ Transition state %s not in states\n", from)
			fmt.Println("[DFA Invalid]")
			return false
		}
		for symbol, to := range trans {
			if !contains(d.Alphabet, symbol) {
				fmt.Printf("❌ Transition symbol %s not in alphabet\n", symbol)
				fmt.Println("[DFA Invalid]")
				return false
			}
			if !contains(d.States, to) {
				fmt.Printf("❌ Transition destination %s not in states\n", to)
				fmt.Println("[DFA Invalid]")
				return false
			}
		}
	}
	fmt.Println("PASS")

	fmt.Println("[DFA Valid]")
	return true
}

func (d *DFA) ExportToDot(filename string, trace []TransitionTrace) error{
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Fprintln(file, "digraph DFA {")
	fmt.Fprintln(file, "  rankdir=LR;")
	fmt.Fprintln(file, `  "" -> `+d.StartState+`;`)

	// accept states double circle
	for _, s := range d.States {
		shape := "circle"
		if d.acceptMap[s] {
			shape = "doublecircle"
		}
		fmt.Fprintf(file, "  %s [shape=%s];\n", s, shape)
	}

	// mark all matched paths with highlight color
	edgeMark := make(map[string]bool)
	for _, t := range trace {
		key := fmt.Sprintf("%s_%s_%s", t.From, t.Symbol, t.To)
		edgeMark[key] = true
	}

	// transition-edges
	for from, trans := range d.Transitions {
		for symbol, to := range trans {
			key := fmt.Sprintf("%s_%s_%s", from, symbol, to)
			if edgeMark[key] {
				fmt.Fprintf(file, `  %s -> %s [label=%s, color=red, penwidth=4];`+"\n", from, to, symbol)
			} else {
				fmt.Fprintf(file, `  %s -> %s [label=%s];`+"\n", from, to, symbol)
			}
		}
	}
	fmt.Fprintf(file, "}")
	return nil
}

// util function to check if a string exists in a slice
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

/*
TODO:
[] 这个 DFA 的状态过程生成 图形可视化（如 Graphviz .dot 文件
[] “合法性检查”和“所有长度 ≤ N 的合法字符串输出”
添加合法性检查（是否所有状态符号转移完备）

枚举所有规则串（长度 ≤ N）：使用 BFS 枚举输入串并判断是否被接受
重构为模块化包（方便后续实验复用）
将 dfa.Match() 改成可输出状态轨迹
*/
