// ===== grammar/grammar.go =====
package grammar

import "strings"

type Grammar struct {
	Start        string
	Productions  []Production
	Terminals    map[string]Symbol
	NonTerminals map[string]Symbol
}

func NewGrammar(start string) Grammar {
	return Grammar{
		Start:        start,
		Productions:  []Production{},
		Terminals:    make(map[string]Symbol),
		NonTerminals: make(map[string]Symbol),
	}
}

func (g *Grammar) symbol(name string) Symbol {
	if name == "ε" {
		return Symbol{Name: "ε", Type: Terminal}
	}
	if strings.ToLower(name) == name && name != "ε" {
		if s, ok := g.Terminals[name]; ok {
			return s
		}
		s := Symbol{Name: name, Type: Terminal}
		g.Terminals[name] = s
		return s
	} else {
		if s, ok := g.NonTerminals[name]; ok {
			return s
		}
		s := Symbol{Name: name, Type: NonTerminal}
		g.NonTerminals[name] = s
		return s
	}
}

func (g *Grammar) AddProduction(lhs string, rhs []string) {
	left := g.symbol(lhs)
	right := []Symbol{}
	for _, name := range rhs {
		right = append(right, g.symbol(name))
	}
	g.Productions = append(g.Productions, Production{Left: left, Right: right})
}
