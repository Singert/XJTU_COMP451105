// ===== analyzer/util.go =====
package analyzer

import "fmt"

func PrintSetMap(sets map[string]map[string]bool) {
	for sym, set := range sets {
		fmt.Printf("%s: { ", sym)
		for tok := range set {
			fmt.Printf("%s ", tok)
		}
		fmt.Println("}")
	}
}
