package main

import (
	"encoding/json"
	"fmt"
	"os"
	"project/lexer"
	"project/parser"
)

type TestCase struct {
	Purpose string `json:"purpose"`
	Code    string `json:"code"`
}

func main() {
	// 读取 test_case.json
	data, err := os.ReadFile("test_case.json")
	if err != nil {
		panic(fmt.Sprintf("读取 test_case.json 失败: %v", err))
	}

	var tests []TestCase
	if err := json.Unmarshal(data, &tests); err != nil {
		panic(fmt.Sprintf("解析 JSON 失败: %v", err))
	}

	for idx, test := range tests {
		fmt.Println("==============================================")
		fmt.Printf("⭐ 测试 #%d: %s\n", idx+1, test.Purpose)
		fmt.Println("----------------------------------------------")
		fmt.Printf("🧪 输入代码: %s\n", test.Code)

		tokens := lexer.Tokenize(test.Code)
		tac := parser.ParseAndGenerateTAC(tokens)

		fmt.Println("💡 生成三地址代码:")
		for _, line := range tac {
			fmt.Printf("    %s\n", line)
		}
		fmt.Print("==============================================\n\n")
	}
}
