package dfa

import (
	"encoding/json"
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


