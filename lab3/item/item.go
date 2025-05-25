package item

import (
	"fmt"
	"lab3/grammar"
)

// Item represents an item in a grammar, which is a production with a dot indicating the position of the parser.
type Item struct {
	ProdIndex int // Index of the production in the grammar
	DotPos    int // Position of the dot in the production
}

// Equal checks if two items are equal.
func (it Item) Equal(other Item) bool {
	return it.ProdIndex == other.ProdIndex && it.DotPos == other.DotPos
}

// String returns a string representation of the item.
func (it Item) String(g *grammar.Grammar) string {
	prod := g.Productions[it.ProdIndex]
	right := make([]grammar.Symbol, len(prod.Right)+1)
	copy(right, prod.Right[:it.DotPos])
	right[it.DotPos] = "."
	copy(right[it.DotPos+1:], prod.Right[it.DotPos:])
	return fmt.Sprintf("%s -> %s", prod.Left, jion(right))
}

// Jion help print
func jion(syms []grammar.Symbol) string {
	str := ""
	for _, s := range syms {
		str += string(s) + " "
	}
	return str
}

// Closure computes the closure of a set of items in the context of a grammar.
func Closure(g *grammar.Grammar, items []Item) []Item {
	closure := make([]Item, len(items))
	copy(closure, items)

	changed := true
	for changed {
		changed = false
		for _, item := range closure {
			prod := g.Productions[item.ProdIndex]
			if item.DotPos >= len(prod.Right) {
				continue // Dot is at the end of the production, so no further processing needed
			}
			symb := prod.Right[item.DotPos]
			if !g.NonTerms[symb] {
				continue // Only process non-terminals
			}
			// Find all productions for the non-terminal symbol
			for i, prod := range g.Productions {
				if prod.Left == symb {
					newItem := Item{ProdIndex: i, DotPos: 0}
					// Check if the new item is already in the closure
					if !containsItem(closure, newItem) {
						closure = append(closure, newItem)
						changed = true // We added a new item, so we need to continue processing
					}
				}
			}
		}
	}
	return closure
}

// containsItem checks if a slice of items contains a specific item.
func containsItem(items []Item, target Item) bool {
	for _, it := range items {
		if it.Equal(target) {
			return true
		}
	}
	return false
}

// Goto computes the goto set for a given set of items and a symbol in the context of a grammar.

func Goto(g *grammar.Grammar, items []Item, symb grammar.Symbol) []Item {
	var moved []Item
	for _, item := range items {
		prod := g.Productions[item.ProdIndex]
		if item.DotPos < len(prod.Right) && prod.Right[item.DotPos] == symb {
			moved = append(moved, Item{ProdIndex: item.ProdIndex, DotPos: item.DotPos + 1})
		}
	}
	return Closure(g, moved)
}
