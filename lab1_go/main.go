package main

import (
	"flag"
	"fmt"
	"lab1/dfa"
	"strconv"
)

func main() {
	// define command line arguments
	enableEnum := flag.Bool("e", false, "Enable enumeration of valid strings")
	dotPath := flag.String("d", "./dot", "Path to save dot file")
	flag.StringVar(dotPath, "dot", "./dot", "Path to save dot file")
	flag.BoolVar(enableEnum, "enum", false, "Enable enumeration of valid strings")
	flag.Parse()
	maxLength := 3

	args := flag.Args()
	if len(args) > 0 {
		var err error
		maxLength, err = strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Invalid max length argument. Using default value of 3.")
		}
	}

	dfa, err := dfa.LoadDFAFromJson("./json/dfa.json")
	if err != nil {
		fmt.Println("Error loading DFA:", err)
		return
	}

	if *enableEnum {
		fmt.Println("A")
		validStrings := dfa.EnumValidStrings(maxLength)
		fmt.Printf("[Enum] Valid strings of length <= %d:\n", maxLength)
		for _, str := range validStrings {
			fmt.Println(str)
		}
		fmt.Println("[Enum] Done")
	}

	for {
		var (
			input   string
			dotName string
		)
		fmt.Print("Enter a string to match (or 'quit' to quit): ")
		fmt.Scanln(&input)

		if input == "quit" {
			break
		}
		Matched, trace := dfa.MatchDFA(input)
		if Matched {
			fmt.Println("✅ Accepted")
			dotName = *dotPath+"/"+input + ".dot"
			dfa.ExportToDot(dotName, trace)
		} else {
			fmt.Println("❌ Not Accepted")
		}
	}
}
