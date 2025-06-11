package main

import (
	"fmt"
	"os"
	"project/backend"
	"project/lexer"
	"project/parser"
)

func readSourceFile(filename string) string {
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(fmt.Sprintf("读取源代码文件失败: %v", err))
	}
	return string(data)
}

func main() {
	fmt.Println("===== 中间代码生成器（main.src） =====")

	source := readSourceFile("main.src")
	tokens := lexer.Tokenize(source)

	tac := parser.ParseProgram(tokens)

	fmt.Println("\n💡 完整三地址代码输出:")
	for _, line := range tac {
		fmt.Println("    " + line)
	}

	fmt.Println("===== 生成完成 =====")
	fmt.Println("\n🛠 生成 MIPS 汇编:")
	mips := backend.GenerateMIPS(tac)
	for _, line := range mips {
		fmt.Println(line)
	}
	fmt.Println("===== 生成完成 =====")
}
