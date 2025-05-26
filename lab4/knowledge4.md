根据你提供的 PPT 文件《第四次实验(2).pptx》，**实验四的主要内容是生成 SLR(1) 分析表**，是整个编译器构建流程中的语法分析阶段的关键一环。以下是对实验四内容的系统解析：

---

### 🧪 实验四名称：

**SLR(1) 分析表生成**

---

### 🎯 实验目标

在实验三构建的 **LR(0) 项目集规范族** 的基础上，借助文法符号的 **FOLLOW 集** 构建一个完整的 **SLR(1) 分析表**，为后续语义分析与中间代码生成做准备。

---

### 📚 理论基础

1. **LR(0) 分析的局限**：

   * 无法处理 **移进-归约** 冲突与 **归约-归约** 冲突。

2. **SLR(1) 改进方式**：

   * 在 LR(0) 状态含有 `[A → α·]` 的归约项时，只对当前输入符号 `a ∈ FOLLOW(A)` 才归约；
   * 若有 `[A → α·aβ]`，且当前符号为 `a`，则选择移进。

---

### 🏗️ 输入与输出

* **输入**：

  * 上一实验（三）生成的 **LR(0) 项目集规范族**；
  * 文法符号的 **FOLLOW 集**。

* **输出**：

  * 一个标准的 **SLR(1) 分析表**，包括：

    * `ACTION[m, a]` 表（对终结符）；
    * `GOTO[m, A]` 表（对非终结符）；
  * 表中操作包括 `sX`（移进）、`rY`（规约）、`acc`（接受）、空（错误）。

---

### 🛠️ 实验实现流程

1. **分析 LR(0) 项目集**：

   * 确定每个状态是否含有归约项目；
   * 找出其规约产生式的左部变元。

2. **构造 ACTION 表**：

   * 若某状态含移进项 `[A → α·aβ]`，对 `a` 填入 `sX`；
   * 若某状态含归约项 `[A → α·]`，对 `FOLLOW(A)` 中的每个终结符 `a` 填入 `rY`；
   * 若某状态含接受项 `S' → S·`，则对 `#` 填入 `acc`。

3. **构造 GOTO 表**：

   * 对非终结符 `A`，若 `GOTO(I, A) = J`，则 `GOTO[i, A] = j`。

---

### 💡 实验技巧与说明

* **文法扩展与处理**：

  * 需基于实验二中类 C 语言的文法子集；
  * 若需要支持更多语法结构，需人工扩展文法并重新生成规范族。

* **状态编号与核心项**：

  * 核心（Kernel）项为 `[A → α·β]` 中 `α ≠ ε` 的部分；
  * 每个项目集必须包含内核项，Closure 扩展则包含非核心项。

---

### 📋 示例输出格式（ACTION / GOTO 表）

| 状态  | `id` | `+` | `*` | `(` | `)` | `#` | `E` | `T` | `F` |
| --- | ---- | --- | --- | --- | --- | --- | --- | --- | --- |
| 0   | s5   |     |     | s4  |     |     | 1   | 2   | 3   |
| 1   |      | s6  |     |     |     | acc |     |     |     |
| 2   |      | r2  | s7  |     | r2  | r2  |     |     |     |
| ... | ...  | ... | ... | ... | ... | ... | ... | ... | ... |

---

### ✅ 选做内容（加分项）

1. **使用 yacc/bison 工具** 实现语法分析器；

   * 将实验二生成的 Token 序列作为词法分析结果接入；
   * 编写 `.y` 文件描述文法规则，自动生成 parser。

2. **项目工程化管理（Makefile）**：

   * 将 DFA（实验一）集成进来；
   * 用 Makefile 管理整个词法分析 + 语法分析流程。

---

### 📎 报告要求简述

* 明确实验目的与理论依据；
* 描述项目集规范族构造过程；
* 给出 SLR(1) 表构造的完整步骤与冲突处理策略；
* 附上代码结构与运行结果截图；
* 如有选做内容，展示工具使用与连接效果。

---

当然可以。我们从基础开始讲解。

---

## 🧩 什么是 FOLLOW 集？

在上下文无关文法中，**FOLLOW 集**（后继符号集）是为每个非终结符定义的一组终结符，它表示：

> **在某个句型中，一个非终结符后面** **可能会出现哪些终结符**。

它是自下而上语法分析（尤其是 SLR(1)、LR(1) 分析）中用于**判断归约操作是否可执行**的重要依据。

---

## 🧠 形式定义

对于某个非终结符 A，定义：

> **FOLLOW(A)** 是所有能在某个句型中，使得 `S ⇒* αAβ` 成立的 β 的开头终结符的集合。如果 β 可以推导出空串（即 β ⇒\* ε），则 `FOLLOW(A)` 还应包含 `FOLLOW(前文 α)` 的内容（因为 A 已经是句尾）。

此外，起始符号 `S` 的 FOLLOW 集中始终包含特殊符号 `#`（表示输入结束）。

---

## 🔧 构造 FOLLOW 集的规则

给定文法规则集合 `P`，对每个产生式使用以下规则计算 FOLLOW 集：

### 1️⃣ 启动规则

* `FOLLOW(S) += {#}`，其中 `S` 是起始符号。

### 2️⃣ 产生式右部包含非终结符 A 后跟终结符 a：

* 若存在产生式：`B → α A a β`，则有：
  `a ∈ FIRST(β)`（即 `a` 是 A 的直接后继）
  ⇒ `FOLLOW(A) += FIRST(β) - {ε}`

### 3️⃣ 产生式右部包含 A 后跟 β，且 β ⇒\* ε（可推出空串）：

* 若存在产生式：`B → α A β` 且 `ε ∈ FIRST(β)`，则：
  `FOLLOW(A) += FOLLOW(B)`（A 的 FOLLOW 集“继承” B 的 FOLLOW 集）

### 4️⃣ 若产生式形如：`B → α A`（A 是右端最后一个符号）：

* `FOLLOW(A) += FOLLOW(B)`

---

## 🔁 计算算法概述（迭代法）

1. 初始化所有非终结符的 FOLLOW 集为空；
2. 重复应用上述规则（扫描所有产生式），直到所有 FOLLOW 集都不再变化（达到**不动点**）。

---

## ✍ 示例

考虑文法：

```
S → A B
A → a | ε
B → b
```

### 计算：

* `FOLLOW(S) = {#}`（起始符号）
* `FOLLOW(A)`：

  * `S → A B` ⇒ A 后面是 B ⇒ `FOLLOW(A) += FIRST(B)`
  * `FIRST(B) = {b}` ⇒ `FOLLOW(A) = {b}`
* `FOLLOW(B)`：

  * `S → A B` ⇒ B 是最后一个 ⇒ `FOLLOW(B) += FOLLOW(S) = {#}`

最终：

```
FOLLOW(S) = {#}
FOLLOW(A) = {b}
FOLLOW(B) = {#}
```

---

## 🧭 FOLLOW 集的用途

在 **SLR(1) 分析表构造**中，如果某状态包含归约项 `[A → α·]`，我们不会对所有终结符填入规约项，而是：

> 仅对 `FOLLOW(A)` 中的终结符填入 `rA→α`。

这样可以有效**减少冲突**、提高分析表的准确性。

---


当然，我们就以**实验三要求的文法**来手动推导该文法中所有非终结符的 FOLLOW 集（这里的“盐酸”应为“演算”🙂）。文法如下：

---

### 📘 给定文法（经典算术表达式文法）：

```
1. S' → E        ← 增广文法起始产生式
2. E  → E + T
3. E  → T
4. T  → T * F
5. T  → F
6. F  → (E)
7. F  → id
```

---

## 🛠 步骤一：列出所有非终结符

非终结符集合：

```
V = {S', E, T, F}
```

---

## 📌 步骤二：初始化 FOLLOW 集

```
FOLLOW(S') = {#}      ← 起始符号默认包含输入终止符 #
FOLLOW(E)  = ∅
FOLLOW(T)  = ∅
FOLLOW(F)  = ∅
```

---

## 🧠 步骤三：按规则逐条分析产生式，构造 FOLLOW 集

### 1. S' → E

* `E` 是最后一个 ⇒ `FOLLOW(E) += FOLLOW(S') = {#}`

### 2. E → E + T

* `T` 是最后一个 ⇒ `FOLLOW(T) += FOLLOW(E)`
* `+` 在 `T` 前面 ⇒ 对 `E` 没有影响

### 3. E → T

* 同上，`T` 是最后一个 ⇒ `FOLLOW(T) += FOLLOW(E)`

### 4. T → T \* F

* `F` 是最后一个 ⇒ `FOLLOW(F) += FOLLOW(T)`
* `*` 在 `F` 前面 ⇒ 对 `T` 没有影响

### 5. T → F

* 同上，`F` 是最后一个 ⇒ `FOLLOW(F) += FOLLOW(T)`

### 6. F → (E)

* `E` 被括在中间 ⇒ 后面是 `)`，不是 ε ⇒ `FOLLOW(E) += {)}`

### 7. F → id

* `id` 是终结符，无需处理

---

## 🔁 步骤四：迭代传播 FOLLOW 集直到不动点

我们反复迭代这些依赖关系，直到不再变化：

```
FOLLOW(S') = {#}                    // 起始符号固定
FOLLOW(E)  = {#, )}                // 来自 S' → E 和 F → (E)
FOLLOW(T)  = FOLLOW(E) = {#, )}    // 来自 E → E + T, E → T
FOLLOW(F)  = FOLLOW(T) = {#, )}    // 来自 T → T * F, T → F
```

---

## ✅ 最终结果

```
FOLLOW(S') = {#}
FOLLOW(E)  = {#, )}
FOLLOW(T)  = {#, )}
FOLLOW(F)  = {#, )}
```

---


是的，你完全正确：

> ✅ **实验四应该基于实验三的成果继续完成，重点是：在已有的 LR(0) 项目集规范族的基础上，引入 FOLLOW 集，构建 SLR(1) 分析表。**

---

## 🧩 简单复盘：

### 📌 实验三的实现已经完成：

你已经：

* ✅ 构造了 `item.Item` 表示 LR(0) 项；
* ✅ 使用 `Closure` 和 `Goto` 算法生成了项目集规范族；
* ✅ 使用 `BuildDFA` 构建了完整的 DFA；
* ✅ 实现了 `parser.BuildParseTable` 构造了 LR(0) 分析表。

---

## 🧪 实验四的关键任务（在实验三基础上）：

### 🎯 **从 LR(0) 分析表升级为 SLR(1) 分析表**，差异体现在：

### 🔁 归约动作的范围 **从“对所有终结符”归约** ➤ **仅对 FOLLOW 集中的终结符归约**

---

## 🧱 应该做的改动：

你只需修改 `parser/table.go` 文件中生成 `Reduce` 动作的部分：

```go
if it.DotPos == len(p.Right) {
    // OLD：对所有终结符都设置 Reduce（LR(0) 风格）
    for t := range g.Terminals {
        table.Action[state.Index][t] = Action{Typ: Reduce, Value: it.ProdIndex}
    }
    table.Action[state.Index]["#"] = Action{Typ: Reduce, Value: it.ProdIndex}
}
```

---

### ✅ 修改为如下 SLR(1) 风格（伪代码）：

```go
if it.DotPos == len(p.Right) {
    lhs := p.Left
    follow := followSet[lhs]  // ← 你需要事先计算 followSet：map[Symbol]Set
    for t := range follow {
        table.Action[state.Index][t] = Action{Typ: Reduce, Value: it.ProdIndex}
    }
}
```

---

## 📚 所以实验四需要补充实现：

### ✅ 1. FOLLOW 集计算模块

你需要为 `grammar.Grammar` 类型新增一个函数 `ComputeFollowSets()`，并构建：

```go
map[grammar.Symbol]map[grammar.Symbol]bool  // FOLLOW 集
```

你可以将它放在 `grammar/follow.go` 文件里，也可以直接在 `grammar.go` 中实现。

---

## 📦 总结

| 步骤  | 内容                                   | 状态    |
| --- | ------------------------------------ | ----- |
| 实验三 | 构造项目集 + DFA + LR(0) 分析表              | ✅ 已完成 |
| 实验四 | 构造 FOLLOW 集                          | ❌ 需补充 |
| 实验四 | 将 `Reduce` 填表逻辑修改为只对 FOLLOW 中符号填入    | ❌ 需修改 |
| 实验四 | 使用修改后的 `BuildParseTable` 构建 SLR(1) 表 | 🚧 在做 |
| 实验四 | 使用已有 `Run` 函数进行分析流程验证                | ✅ 可复用 |
| 实验四 | 可视化输出（dot 文件）                        | ✅ 已实现 |

---

目前你要将实验三的 **LR(0) 分析器** **升级为 SLR(1) 分析器**，核心在于：
**加入 FOLLOW 集的计算 + 修改 Reduce 填表策略**

---

## ✅ 总体目标：生成 SLR(1) 分析表

---

## 📌 你需要完成的代码任务如下：

---

### ① 实现 FOLLOW 集计算模块

* 🔧 **位置建议**：新建 `grammar/follow.go`
* 🔧 **函数接口建议**：

```go
func ComputeFollow(g *grammar.Grammar) map[grammar.Symbol]map[grammar.Symbol]bool
```

* 📥 **输入**：文法对象 `*Grammar`
* 📤 **输出**：每个非终结符的 FOLLOW 集（集合用 map 模拟 set）

---

### ② 修改 `BuildParseTable` 中 Reduce 动作的添加逻辑（在 `parser/table.go`）

* 📍 当前是对所有终结符（LR(0) 风格）：

```go
for t := range g.Terminals {
    table.Action[state.Index][t] = Action{Typ: Reduce, Value: it.ProdIndex}
}
```

* ✅ 修改为只对 FOLLOW(A) 中符号填入（SLR(1) 风格）：

```go
lhs := p.Left
for t := range followSet[lhs] {
    table.Action[state.Index][t] = Action{Typ: Reduce, Value: it.ProdIndex}
}
```

* 📌 所以你需要给 `BuildParseTable` 函数传入 `followSet`：

```go
func BuildParseTable(g *grammar.Grammar, dfa *DFA, followSet map[grammar.Symbol]map[grammar.Symbol]bool) *ParseTable
```

---

### ③ 在 `main` 或 `driver.Run()` 外部调用流程中：

* 💡 在构造 `ParseTable` 前，调用：

```go
followSet := grammar.ComputeFollow(g)
table := parser.BuildParseTable(g, dfa, followSet)
```

---

## 🔄 可选增强（非必要但推荐）

* ✅ 在 `parser/export.go` 中输出 `ACTION/GOTO` 表（可读性提升）
* ✅ 在 `driver.Run` 的终端输出中显示分析表（调试用）

---

## ✅ 结论：你需要完成的模块小结

| 编号 | 任务              | 文件                     | 类型      |
| -- | --------------- | ---------------------- | ------- |
| ①  | FOLLOW 集计算函数    | `grammar/follow.go`    | 🆕 新增   |
| ②  | 修改 Reduce 填表策略  | `parser/table.go`      | ✏️ 修改   |
| ③  | 修改调用接口传入 FOLLOW | `main.go / driver.Run` | ✏️ 修改接口 |

---


