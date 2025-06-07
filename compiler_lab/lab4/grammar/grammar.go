package grammar

//Sybol represents a symbol in a grammar, which can be a terminal or non-terminal.
type Symbol string 

// Production represents a production rule in a grammar.
type Production struct {
	Left  Symbol
	Right []Symbol
}

type Grammar struct {
	Productions []Production
	StartSymbol Symbol // StartSymbol is the symbol from which the grammar starts.
	Terminals map[Symbol]bool // Terminals are the symbols that represent actual tokens in the language.
	NonTerms map[Symbol]bool // NonTerms are the symbols that represent non-terminal constructs in the language.
}


// NewGrammar creates a new Grammar instance with the specified start symbol.
func NewGrammar(start Symbol) *Grammar {
	return &Grammar{
		StartSymbol: start,
		Terminals:   make(map[Symbol]bool),
		NonTerms:    make(map[Symbol]bool),
	}
}

// AddProduction adds a production rule to the grammar.
func (g *Grammar) AddProduction(left Symbol, right []Symbol) {
	g.Productions = append(g.Productions, Production{
		Left:  left,
		Right: right,
	})
	// Mark the left side as a non-terminal
	g.NonTerms[left] = true
	// Mark the right side symbols as terminals or non-terminals
	for _, symb := range right {
		if isTerminal(symb) {
			g.Terminals[symb] = true
		} else {
			g.NonTerms[symb] = true
		}
	}	
}

// isTerminal checks if a symbol is a terminal. Simply by if lower letter or identifier.
func isTerminal(symb Symbol) bool {
	// A terminal is typically a symbol that does not start with an uppercase letter.
	// This is a simple heuristic; adjust as necessary for your grammar.
		if symb == "id" || symb == "+" || symb == "*" || symb == "(" || symb == ")" {
		return true
	}
	return false
}