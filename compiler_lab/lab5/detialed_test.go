package main_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"lab5/lexer"
	"lab5/parser"
	"lab5/syntax"
	"lab5/utils"

	"testing"
)

func TestLab5Detialed(t *testing.T) {
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
	files, err := filepath.Glob("testcases/*.c")
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
		tokens := scanner.Tokenize(code, true)
		symbols := utils.TokensToSymbols(tokens)
		// 4. 初始化符号表
		// 语法 + 语义分析

		parsererr := parser.Run(symbols, g, dfa, table, tokens, true, file)
		if parsererr != nil {
			fmt.Println(parsererr.Error())
		}
	}
	fmt.Println("🤬 测试完成！")
	dotfiles, err := filepath.Glob("output/*.dot")
	if err != nil {
		fmt.Println("❌ 获取输出文件失败:", err)
		return
	}
	for _, dotfile := range dotfiles {
		//将dot文件转为 png
		pngfile := dotfile[:len(dotfile)-4] + ".png"
		cmd := exec.Command("dot", "-Tpng", dotfile, "-o", pngfile)
		err := cmd.Run()
		if err != nil {
			fmt.Printf("❌ 转换 %s 为 PNG 失败: %v\n", dotfile, err)
			continue
		}
		fmt.Printf("✔️ 转换 %s 为 PNG 成功: %s\n", dotfile, pngfile)
	}
}
