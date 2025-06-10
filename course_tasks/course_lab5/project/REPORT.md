
---

# 中间代码生成器作业报告

## 一、实验背景

本实验旨在通过实现一个中间代码生成器，将源程序的结构化源代码转化为中间代码（三地址代码）。三地址代码是编译过程中间的一种表达方式，通常用于后续的优化和目标代码生成。实验主要包括以下任务：

* **词法分析**：将源代码分割为一系列的标记（tokens）。
* **语法分析**：将标记转换为抽象语法树（AST）或类似结构。
* **中间代码生成**：根据语法分析结果，生成三地址代码（TAC）。

本实验中，主要解析了具有函数定义、条件语句、循环语句及函数调用等特征的C风格源代码。

## 二、项目概述

该项目的核心目标是实现一个中间代码生成器，能够从源代码（例如 `main.src` 文件）生成相应的三地址代码。具体步骤如下：

1. **词法分析**：解析源代码并将其转化为标记（tokens）。
2. **语法分析**：根据标记解析出程序结构，并将其组织成适当的语句或表达式。
3. **中间代码生成**：根据解析结果生成三地址代码，用于后续的优化和目标代码生成。

### 项目目录结构

```
├── assets
├── boolean
│   └── expr.go          // 布尔表达式生成
├── expr
│   └── expr.go          // 表达式生成工具
├── generator
│   └── tac.go           // 三地址代码生成
├── lexer
│   └── lexer.go         // 词法分析器
├── parser
│   └── parser.go        // 语法分析器
├── stmt
│   ├── array_assign.go  // 数组赋值语句处理
│   ├── call.go          // 函数调用处理
│   ├── dispatch.go      // 语句分发处理
│   ├── function.go      // 函数定义处理
│   ├── if.go            // if语句处理
│   ├── return.go        // return语句处理
│   ├── stmtlist.go      // 语句块处理
│   └── while.go         // while语句处理
└── main.go              // 主入口文件
```

## 三、核心技术实现

### 1. **词法分析（Tokenizer）**

词法分析器的任务是将源代码字符串转化为一系列的标记。主要使用了字符流解析技术，处理了包括数字、标识符、运算符等在内的各种语言元素。

```go
func Tokenize(input string) []string {
    var tokens []string
    current := ""
    runes := []rune(input)
    ...
    // 处理标识符、运算符和括号
}
```

### 2. **语法分析（Parser）**

语法分析器负责解析标记并将其组织成语法结构。具体使用了递归下降法（Recursive Descent Parsing）来处理不同的语句和表达式。

```go
func ParseProgram(tokens []string) []string {
    var code []string
    i := 0
    for i < len(tokens) {
        end := findStmtEnd(tokens, i)
        stmtTokens := tokens[i:end]
        stmtCode := stmt.Dispatch(stmtTokens)
        code = append(code, stmtCode...)
        i = end
    }
    return code
}
```

### 3. **三地址代码生成**

三地址代码生成的核心是将解析后的语法结构转换为三地址代码。主要操作包括处理算术运算、条件语句、循环结构等。三地址代码的生成函数会为每个语句生成一系列中间操作，如加法、赋值、条件跳转等。

```go
func GenerateFunctionDef(tokens []string) []string {
    ...
    // 处理函数参数、函数体以及局部变量
}
```

### 4. **语句分发（Dispatch）**

根据不同的语法结构（如函数定义、赋值语句、条件语句等），语句分发函数将调用相应的生成函数进行处理。

```go
func Dispatch(tokens []string) []string {
    ...
    // 判断并调用不同的生成函数（如GenerateFunctionDef，GenerateIfElse等）
}
```

## 四、实验中的关键问题与修复

### 1. **函数参数解析问题**

在处理函数定义时，遇到函数指针作为参数（如 `int soo()`）时，解析存在问题。为解决此问题，修改了参数提取逻辑，正确跳过了函数指针类型。

修复后的代码：

```go
for ; tokens[i] != ")"; i++ {
    if tokens[i] == "," || tokens[i] == "int" || tokens[i] == "void" {
        continue
    }
    if tokens[i+1] == "(" {
        params = append(params, tokens[i])
        for tokens[i] != ")" {
            i++
        }
    } else {
        params = append(params, tokens[i])
    }
}
```

### 2. **嵌套函数解析顺序错误**

在解析嵌套函数时，顺序没有正确处理，导致部分函数未被正确生成。通过改进 `ParseStmtList()` 中对函数体的递归处理，解决了该问题。

修复后的代码：

```go
stmtTokens := inner[start : braceEnd+1]
stmtCode := GenerateFunctionDef(stmtTokens[1:])
code = append(code, stmtCode...)
start = braceEnd + 1
continue // 确保函数解析顺序正确
```

### 3. **三地址代码生成错误**

生成的三地址代码中出现了无关的 `t1 = i + 1` 等多余操作，原因是 fallback 机制未能准确判断是否是合法的三地址代码生成结构。通过增强函数调用条件，确保只有在正确的情况下生成代码。

修复后的代码：

```go
if len(tokens) >= 4 && tokens[1] == "(" && tokens[len(tokens)-1] == ";" {
    paren := 0
    for i := 1; i < len(tokens)-1; i++ {
        if tokens[i] == "(" {
            paren++
        } else if tokens[i] == ")" {
            paren--
        }
    }
    if paren == 0 {
        return GenerateFunctionCall(tokens)
    }
}
```

## 五、实验结果

通过修复上述问题，程序成功生成了对于题目11.2给出的源代码符合预期的三地址代码，完整的输出如下：

```txt
LABEL FUNC_raw
POP x
t1 = x + 5
y = t1
RETURN y
ENDFUNC raw
LABEL FUNC_foo
POP y
LABEL FUNC_bar
POP x
POP soo
t2 = x > 3
IF t2 != 0 THEN L1 ELSE L2
LABEL L1
t3 = x / 3
PAR soo
PAR t3
t4 = CALL bar, 2
GOTO L3
LABEL L2
z = soo
LABEL L3
PRINT z
ENDFUNC bar
PAR raw
PAR y
t5 = CALL bar, 2
ENDFUNC foo
PAR 6
t6 = CALL foo, 1
```

## 六、总结

本实验成功实现了一个中间代码生成器，能够从源代码中解析出三地址代码。通过对词法分析、语法分析以及三地址代码生成过程的逐步改进，我们解决了嵌套函数解析、函数指针参数处理以及代码生成中的冗余问题。未来可以继续扩展对更多语言特性的支持，如数组操作、结构体支持等。
