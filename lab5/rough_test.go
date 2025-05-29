package main_test

import (
	"fmt"
	"os"
	"path/filepath"

	"lab5/lexer"
	"lab5/parser"
	"lab5/syntax"
	"lab5/utils"

	"testing"
)

func TestLab5Rough(t *testing.T) {
	var errs int = 0
	// 1. 加载语法、构建 DFA 和表
	g := syntax.DefineGrammar()
	dfa := parser.BuildDFA(g)
	follow := syntax.ComputeFollow(g)
	table := parser.BuildParseTable(g, dfa, follow)
	// parser.PrintParseTable(table, g)
	fmt.Print("test start \n")

	// 2. 加载 DFA 词法器
	dfaWithType, err := lexer.LoadMultiDFAFromJson("assets/all_dfa.json", "dot", false)
	if err != nil {
		fmt.Println("❌ DFA 加载失败:", err)
		return
	}
	scanner := lexer.NewScanner()
	for _, d := range *dfaWithType {
		scanner.RegisterDFA(d.DFA, d.TokenType)
	}
	fmt.Println("Starting tests...")
	// 3. 遍历 testcases 文件夹
	files, err := filepath.Glob("testcases/finalcases/*.c")
	if err != nil {
		fmt.Println("❌ 获取测试文件失败:", err)
		return
	}
	fmt.Println("找到测试文件:", len(files))
	for _, file := range files {
		fmt.Printf("\n======== 测试文件: %s ========\n", file)
		data, err := os.ReadFile(file)
		if err != nil {
			fmt.Printf("❌ 读取失败: %v\n", err)
			continue
		}

		code := string(data)
		tokens := scanner.Tokenize(code, false)
		symbols := utils.TokensToSymbols(tokens)

		// 语法 + 语义分析
		parsererr := parser.Run(symbols, g, dfa, table, tokens, false)
		if parsererr != nil {
			errs++
			fmt.Println(parsererr.Error())
		}
	}
	if errs > 0 {
		fmt.Printf("🤬 测试完成！共 %d 个错误\n", errs)
		return
	}
	fmt.Println("🤬 测试完成！")
}
