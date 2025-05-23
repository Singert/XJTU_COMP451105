package dfa

import (
	"encoding/json"
	"fmt"
	"os"
)

func LoadDFAFromJson(fileName string) (*DFA, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	var dfa DFA
	err = json.Unmarshal(data, &dfa)
	if err != nil {
		return nil, err
	}
	dfa.CheckValidity()
	dfa.buildAcceptMap()
	dfa.ExportDFAtoDot("./dot/dfa.dot")
	return &dfa, nil
}

func LoadMultiDFAFromJson(fileName string, dotPath string) (*[]DFAWithTokenType, error) {
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
		fmt.Printf("Loaded DFA from %s, q0 transitions: %+v\n", fileName, dfas[i].DFA.Transitions["q0"])
		dfas[i].DFA.CheckValidity()
		dfas[i].DFA.buildAcceptMap()
		dotPath := dotPath + "/" + string(dfas[i].TokenType) + ".dot"
		dfas[i].DFA.ExportDFAtoDot(dotPath)
	}
	return &dfas, nil
}
