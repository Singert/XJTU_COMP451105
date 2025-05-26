package lexer

import (
	"encoding/json"
	"fmt"
	"os"
)

func LoadDFAFromJson(fileName string, verbose bool) (*DFA, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	var dfa DFA
	err = json.Unmarshal(data, &dfa)
	if err != nil {
		return nil, err
	}
	dfa.CheckValidity(verbose, "")
	dfa.buildAcceptMap()
	dfa.ExportDFAtoDot("./dot/dfa.dot")
	return &dfa, nil
}

func LoadMultiDFAFromJson(fileName string, dotPath string, verbose bool) (*[]DFAWithTokenType, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	var dfas []DFAWithTokenType
	err = json.Unmarshal(data, &dfas)
	if err != nil {
		return nil, err
	}
	for i := range dfas {
		if verbose {
			fmt.Printf("Loaded DFA from %s, q0 transitions: %+v\n", fileName, dfas[i].DFA.Transitions["q0"])
		}
		dfas[i].DFA.CheckValidity(verbose, dfas[i].TokenType)
		dfas[i].DFA.buildAcceptMap()
		dotPath := dotPath + "/" + string(dfas[i].TokenType) + ".dot"
		dfas[i].DFA.ExportDFAtoDot(dotPath)
	}
	return &dfas, nil
}
