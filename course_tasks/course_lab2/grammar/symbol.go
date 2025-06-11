// ===== grammar/symbol.go =====
package grammar

type SymbolType int

const (
	Terminal SymbolType = iota
	NonTerminal
)

type Symbol struct {
	Name string
	Type SymbolType
}

func (s Symbol) IsTerminal() bool {
	return s.Type == Terminal
}

func (s Symbol) IsNonTerminal() bool {
	return s.Type == NonTerminal
}
