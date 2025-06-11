// main.go
package main

import (
	"fmt"
	"ll1-analyzer/analyzer"
	"ll1-analyzer/grammar"
)

func main() {
	g := grammar.NewGrammar("P")

	// 手动添加文法产生式
	g.AddProduction("P", []string{"Ď", "Š"})
	g.AddProduction("Ď", []string{"ε"})
	g.AddProduction("Ď", []string{"Ď", "D", ";"})
	g.AddProduction("D", []string{"T", "d"})
	g.AddProduction("D", []string{"T", "d", "[", "]"})
	g.AddProduction("D", []string{"T", "d", "(", "Ǎ", ")", "{", "Ď", "Š", "}"})
	g.AddProduction("T", []string{"int"})
	g.AddProduction("T", []string{"void"})
	g.AddProduction("Ǎ", []string{"ε"})
	g.AddProduction("Ǎ", []string{"Ǎ", "A", ";"})
	g.AddProduction("A", []string{"T", "d"})
	g.AddProduction("A", []string{"d", "[", "]"})
	g.AddProduction("A", []string{"T", "d", "(", ")"})
	g.AddProduction("Š", []string{"S"})
	g.AddProduction("Š", []string{"Š", ";", "S"})
	g.AddProduction("S", []string{"d", "=", "E"})
	g.AddProduction("S", []string{"if", "(", "B", ")", "S"})
	g.AddProduction("S", []string{"if", "(", "B", ")", "S", "else", "S"})
	g.AddProduction("S", []string{"while", "(", "B", ")", "S"})
	g.AddProduction("S", []string{"return", "E"})
	g.AddProduction("S", []string{"{", "Š", "}"})
	g.AddProduction("S", []string{"d", "(", "Ř", ")"})
	g.AddProduction("B", []string{"B", "∧", "B"})
	g.AddProduction("B", []string{"B", "∨", "B"})
	g.AddProduction("B", []string{"E", "r", "E"})
	g.AddProduction("B", []string{"E"})
	g.AddProduction("E", []string{"d", "=", "E"})
	g.AddProduction("E", []string{"i"})
	g.AddProduction("E", []string{"d"})
	g.AddProduction("E", []string{"d", "(", "Ř", ")"})
	g.AddProduction("E", []string{"E", "+", "E"})
	g.AddProduction("E", []string{"E", "*", "E"})
	g.AddProduction("E", []string{"(", "E", ")"})
	g.AddProduction("Ř", []string{"ε"})
	g.AddProduction("Ř", []string{"Ř", "R", ","})
	g.AddProduction("R", []string{"E"})
	g.AddProduction("R", []string{"d", "[", "]"})
	g.AddProduction("R", []string{"d", "(", ")"})

	first := analyzer.ComputeFirstSets(&g)
	follow := analyzer.ComputeFollowSets(&g, first)
	violations := analyzer.CheckLL1(&g, first, follow)

	fmt.Println("== FIRST ==")
	analyzer.PrintSetMap(first)
	fmt.Println("== FOLLOW ==")
	analyzer.PrintSetMap(follow)
	fmt.Println("== LL(1) Violations ==")
	for _, v := range violations {
		fmt.Println(v)
	}
}
