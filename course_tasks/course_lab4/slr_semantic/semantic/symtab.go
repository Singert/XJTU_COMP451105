package semantic

import (
	"fmt"
)

type SymbolTable struct {
	symbols map[string]Symbol
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		symbols: make(map[string]Symbol),
	}
}

func (st *SymbolTable) Add(sym Symbol) error {
	if _, exists := st.symbols[sym.Name]; exists {
		return fmt.Errorf("符号 %s 重复定义", sym.Name)
	}
	st.symbols[sym.Name] = sym
	return nil
}

func (st *SymbolTable) Get(name string) (Symbol, bool) {
	sym, ok := st.symbols[name]
	return sym, ok
}

func (st *SymbolTable) Dump() {
	fmt.Println("🔎 符号表：")
	fmt.Printf("%-10s %-6s %-8s %-10s %s\n", "Name", "Type", "Kind", "Dim", "Params")
	for _, sym := range st.symbols {
		dimStr := fmt.Sprintf("%v", sym.Dim)
		paramStr := "["
		for i, p := range sym.Params {
			if i > 0 {
				paramStr += ", "
			}
			paramStr += fmt.Sprintf("%s %s", p.Type, p.Name)
		}
		paramStr += "]"
		if sym.Kind != "function" {
			paramStr = "-"
		}
		if sym.Kind != "array" {
			dimStr = "-"
		}
		fmt.Printf("%-10s %-6s %-8s %-10s %s\n",
			sym.Name, sym.Type, sym.Kind, dimStr, paramStr)
	}
}
