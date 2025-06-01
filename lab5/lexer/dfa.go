package lexer

import (
	"fmt"
	"os"
)

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

func (d *DFA) MatchDFA(input string, verbose bool) (bool, []TransitionTrace) {
	currentState := d.StartState
	trace := []TransitionTrace{}
	for i, ch := range input {
		symbol := string(ch)
		if !Contains(d.Alphabet, symbol) {
			if verbose {
				fmt.Printf("Step %d: %s --%s--> ❌invalid symbol\n", i+1, currentState, symbol)
			}
			return false, nil
		}
		nextState, ok := d.Transitions[currentState][symbol]
		if !ok {
			if verbose {
				fmt.Printf("Step %d: %s --%s--> ❌invalid transition\n", i+1, currentState, symbol)
			}
			return false, nil
		}

		trace = append(trace, TransitionTrace{
			From:   currentState,
			Symbol: symbol,
			To:     nextState,
		})

		if verbose {
			print("AAA")
			fmt.Printf("Step %d: %s --%s--> %s\n", i+1, currentState, symbol, nextState)
		}
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

func (d *DFA) CheckValidity(verbose bool, tokenType TokenType) bool {
	fmt.Printf("[Checking DFA Validity]: <%s>\n", tokenType)

	// 1️⃣ check start_state not nil
	if verbose {
		fmt.Print("checking start_state not nil: ")
	}
	if d.StartState == "" {
		if verbose {
			fmt.Print("❌ Start state is nil\n")
		}
		fmt.Println("[DFA Invalid]")
		return false
	} else if verbose {
		fmt.Println("PASS")
	}

	// 2️⃣ check start_state in states
	if verbose {
		fmt.Print("checking start_state in states: ")
	}
	if !Contains(d.States, d.StartState) {
		if verbose {
			fmt.Printf("❌ Start state %s not in states\n", d.StartState)
		}
		fmt.Println("[DFA Invalid]")
		return false
	} else if verbose {
		fmt.Println("PASS")
	}

	// 3️⃣ check accept_states not nil
	if verbose {
		fmt.Print("checking accept_states not nil: ")
	}
	if len(d.AcceptStates) == 0 {
		if verbose {
			fmt.Print("❌ Accept states is nil\n")
		}
		fmt.Println("[DFA Invalid]")
		return false
	} else if verbose {
		fmt.Println("PASS")
	}

	// 4️⃣ check accept_states in states
	if verbose {
		fmt.Print("checking accept_states in states: ")
	}
	for _, s := range d.AcceptStates {
		if !Contains(d.States, s) {
			if verbose {
				fmt.Printf("❌ Accept state %s not in states\n", s)
			}
			fmt.Println("[DFA Invalid]")
			return false
		}
	}
	if verbose {
		fmt.Println("PASS")
	}

	// 5️⃣ check transitions not nil
	if verbose {
		fmt.Print("checking transitions not nil: ")
	}
	if len(d.Transitions) == 0 {
		if verbose {
			fmt.Print("❌ Transitions is nil\n")
		}
		fmt.Println("[DFA Invalid]")
		return false
	} else if verbose {
		fmt.Println("PASS")
	}

	// 6️⃣ check all transitions refer to valid states and symbols
	if verbose {
		fmt.Print("checking transitions for valid states and symbols: ")
	}
	for from, trans := range d.Transitions {
		if !Contains(d.States, from) {
			if verbose {
				fmt.Printf("❌ Transition state %s not in states\n", from)
			}
			fmt.Println("[DFA Invalid]")
			return false
		}
		for symbol, to := range trans {
			if !Contains(d.Alphabet, symbol) {
				if verbose {
					fmt.Printf("❌ Transition symbol %s not in alphabet\n", symbol)
				}
				fmt.Println("[DFA Invalid]")
				return false
			}
			if !Contains(d.States, to) {
				if verbose {
					fmt.Printf("❌ Transition destination %s not in states\n", to)
				}
				fmt.Println("[DFA Invalid]")
				return false
			}
		}
	}
	if verbose {
		fmt.Println("PASS")
	}
	for _, s := range d.States {
		for _, sym := range d.Alphabet {
			if _, ok := d.Transitions[s][sym]; !ok {
				if verbose {
					fmt.Printf("[CheckValidity] Warning: State '%s' missing transition on symbol '%s'\n", s, sym)
				}
			}
		}
	}
	fmt.Println("[DFA Valid]")
	return true
}

func (d *DFA) ExportToDot(filename string, trace []TransitionTrace) error {
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




