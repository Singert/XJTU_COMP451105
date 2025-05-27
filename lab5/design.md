你已经非常出色地完成了项目结构的整合 👏！目前你的项目 `lab5/` 已将：

* ✅ 实验二词法分析器整合为 `lexer/` 模块
* ✅ 实验四语法分析器结构统一为 `parser/` + `syntax/` 模块
* ✅ 项目整体文件夹布局清晰、模块独立、职责明确

---

## ✅ 当前项目状态总结：

| 模块          | 说明                                   |
| ----------- | ------------------------------------ |
| `lexer/`    | 已完成：多个 DFA 加载 + 最长匹配词法扫描器            |
| `syntax/`   | 已完成：文法表示 + FIRST/FOLLOW 集计算          |
| `parser/`   | 已完成：LR(0) 项目集 + SLR 表 + 分析驱动         |
| `main.go`   | 预计尚未接通 scanner → parser 的 Symbol 流输入 |
| `semantic/` | 空目录，尚未开始 AST 构建 / 符号表集成              |

---

## 🧭 接下来建议的开发路线图（实验五）

---

### ✅ 第 1 步：整合词法 + 语法分析

> **目标：** 让 `parser.Run()` 支持从 `lexer.Scanner` 生成的 token 流驱动语法分析

#### 🔧 子任务

* 在 `main.go` 中添加逻辑：

  * 加载所有 DFA（你已完成）；
  * 注册 DFA 至 Scanner；
  * 对源代码 `string` 进行扫描，转换为 `[]lexer.Token`；
  * 将 token 映射为 `[]syntax.Symbol`；
  * 调用 `parser.Run(symbols, g, dfa, table)`。

#### 🧩 示例代码片段（建议放入 `main.go`）：

```go
func tokenToSymbol(tok lexer.Token) syntax.Symbol {
	switch tok.Type {
	case lexer.TokenID:
		return "id"
	case lexer.TokenNUM, lexer.TokenFLO:
		return "num"
	case lexer.TokenOP, lexer.TokenDELIM, lexer.TokenKW:
		return syntax.Symbol(tok.Lexeme)
	default:
		return "?"
	}
}
```

---

### ✅ 第 2 步：属性栈 + 语义动作框架初始化

> **目标：** 在归约动作时触发语义动作，构建 AST，或维护符号表

#### 🔧 子任务

* 在 `parser/parser.go` 中：

  * 新增 `attrStack []interface{}`；
  * 每次 `shift` 操作插入 `nil`；
  * 每次 `reduce` 时根据产生式编号触发 `semantic.ActionFuncs[pid](...)`；
  * 将结果压回栈中；
* 在 `semantic/actions.go` 中定义：

  ```go
  var ActionFuncs = map[int]func([]interface{}) interface{}{
      1: func(children []interface{}) interface{} {
          return &ASTNode{Type: "+", Left: children[0].(*ASTNode), Right: children[2].(*ASTNode)}
      },
      ...
  }
  ```

---

### ✅ 第 3 步：AST 结构定义与输出（`semantic/ast.go`）

```go
type ASTNode struct {
	Type  string       // "+", "*", "num", "id", etc.
	Value string       // 用于字面量或变量名
	Left  *ASTNode
	Right *ASTNode
}
```

你可以在 `Accept` 动作后访问 `attrStack[0].(*ASTNode)`，将其打印或生成 `.dot` 可视化图。

---

### ✅ 第 4 步（可选）：添加符号表模块（`semantic/symbol_table.go`）

* 定义作用域栈：

```go
type SymbolTable struct {
	stack []map[string]string // 名称 → 类型
}
```

* 在 `int x;` 之类语句的规约时写入符号表
* 在 `x = 1` 中查表验证类型一致

---

## ✅ 项目推进优先级建议：

| 优先级   | 任务                               |
| ----- | -------------------------------- |
| 🔴 高  | 完成 `main.go` 串联 scanner → parser |
| 🔴 高  | 在 `parser.Run()` 中加上 `attrStack` |
| 🟡 中  | 实现 `semantic/actions.go` 的语义动作   |
| 🟡 中  | `ast.go` 中定义结构 + 输出 DOT          |
| 🟢 可选 | `symbol_table.go` 支持声明/查找/类型检查   |

---

### ✅ 我可以立即帮你生成：

1. `main.go` 整合模板（支持从文件扫描并驱动语法分析）
2. `semantic/actions.go` 初始语义动作函数表
3. `semantic/ast.go` AST 节点结构定义
4. DOT 输出 AST 树的函数

是否需要我逐个帮你写好这些模板文件？你想先从哪一个开始？
