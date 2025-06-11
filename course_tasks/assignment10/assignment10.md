# 作业10：


# 1、 题目11.2
解：
##  一、符号表（静态层级结构）

```text
@table (level 0, outer: NIL)
├── x : int
├── y : int
├── q : FUNC, 参数: s(), int x；局部: int y
│   └── @q_table (level 1, outer: @table)
│       ├── s : FUNPTT
│       ├── x : int
│       └── y : int
├── p : FUNC（无参数）
│   └── @p_table (level 1, outer: @table)
│       ├── r : FUNC, 参数: int x；局部: int z
│       │   └── @r_table (level 2, outer: @p_table)
│       │       ├── x : int
│       │       └── z : int
```

---

##  二、活动树（运行时过程调用结构）

```text
ε()
└── p()
    └── q(r(), 45)
        └── r(55)
```

* `ε()` 是最外层环境
* `p()` 调用 `q(...)`，其中实参是函数 `r` 和表达式 `x*3=15*3=45`
* `q(...)` 在调用 `s(x+10)` 即 `r(55)` 时进入 `r` 的环境

---

##  三、各活动记录（活动生存期中的栈帧内容）

### 1. ε()@frame

```text
<访问链>: NIL
<控制链>: NIL
<返址>: -
x: 15
y: 21
q[1]: q@label
q[0]: _
p[1]: p@label
p[0]: _
```

### 2. p()@frame

```text
<访问链>: ε()
<控制链>: ε()
<返址>: -
r[1]: ε()
r[0]: r@label
t4: 45   // x*3 的值
t5: _    // 暂未使用
```

### 3. q(r(), 45)@frame

```text
<参数2>: 45
<参数1>: &r@label
<访问链>: ε()
<控制链>: p()
<返址>: -
s[1]: p()
s[0]: r@label
x: 45
y: 76    // 来自 r(55) 的返回值
t1: 55   // 计算 x+10
t2: 76   // 函数返回值赋给 y
```

### 4. r(55)@frame

```text
<参数1>: 55
<访问链>: p()
<控制链>: q(r(), 45)
<返址>: -
x: 55
z: 76   // 计算 x + y = 55 + 21
t3: 76
```
---

##  四、栈快照（执行到 `return z;` 前，地址从高到低）

假定当前执行的是 `r(55)`，即第八行，`return z;` 语句**即将返回之前**，`z = x + y = 55 + 21 = 76`。

###  当前活动记录为 `r(55)@frame`，其主调是 `q(r(), 45)`

| 地址  | 内容               | 所属栈帧         | 注释                |
| --- | ---------------- | ------------ | ----------------- |
| 515 | t3 = 76          | r(55)@frame  | 临时变量 t3           |
| 514 | z = 76           | r(55)@frame  | 局部变量 z            |
| 513 | x = 55           | r(55)@frame  | 形参 x              |
| 512 | <返址>             | r(55)@frame  | 返回地址（略）           |
| 511 | <控制链> = 507      | r(55)@frame  | 控制链指向 q 栈帧        |
| 510 | <访问链> = 504      | r(55)@frame  | 访问链指向 p 栈帧        |
| 509 | t1 = 55          | q(...)@frame | x + 10            |
| 508 | t2 = 76          | q(...)@frame | 保存 r 的返回值         |
| 507 | y = 76           | q(...)@frame | 局部变量 y            |
| 506 | x = 45           | q(...)@frame | 形参 x              |
| 505 | s\[1] = 504      | q(...)@frame | s 的访问链 = p\@frame |
| 504 | s\[0] = r\@label | q(...)@frame | s 的代码地址           |
| 503 | <返址>             | q(...)@frame | 返回地址（略）           |
| 502 | <控制链> = 498      | q(...)@frame | 控制链指向 p 栈帧        |
| 501 | <访问链> = 498      | q(...)@frame | 访问链指向 p 栈帧        |
| 500 | ...              | ε()@frame    | 静态全局环境            |

---


# 2 ARM32 指令模板（模式替换对）



###  一、内存访问和赋值（Load/Store）

| 三地址码模式         | ARM32 指令                             |
| -------------- | ------------------------------------ |
| `t = rs + k`   | `ADD rt, rs, #k`                     |
| `rt = M[rs]`   | `LDR rt, [rs]`                       |
| `rt = M[k]`    | `LDR rt, =k` <br> `LDR rt, [rt]`     |
| `M[rs] = rt`   | `STR rt, [rs]`                       |
| `M[k] = rt`    | `LDR rtmp, =k` <br> `STR rt, [rtmp]` |
| `rd = rs + rt` | `ADD rd, rs, rt`                     |
| `rd = rs`      | `MOV rd, rs`                         |
| `rd = k`       | `MOV rd, #k` （如果 k ∈ 8-bit）或用 `LDR`  |
| `rd = rs + k`  | `ADD rd, rs, #k`                     |

---

### 二、跳转（GOTO / LABEL）

| 三地址码模式    | ARM32 指令 |
| --------- | -------- |
| `GOTO l`  | `B l`    |
| `LABEL l` | `l:`     |

---

### 三、条件跳转（IF-THEN-ELSE）

#### 相等判断：`IF rs = rt THEN l1 ELSE l2`

| 三地址码模式                       | ARM32 指令                               |
| ---------------------------- | -------------------------------------- |
| `IF rs = rt THEN l1 ELSE l2` | `CMP rs, rt` <br> `BEQ l1` <br> `B l2` |
| `LABEL l2`                   | `l2:`                                  |
| `LABEL l1`                   | `l1:`                                  |

#### 小于判断：`IF rs < rt THEN l1 ELSE l2`

| 三地址码模式                       | ARM32 指令                               |
| ---------------------------- | -------------------------------------- |
| `IF rs < rt THEN l1 ELSE l2` | `CMP rs, rt` <br> `BLT l1` <br> `B l2` |
| `LABEL l1`                   | `l1:`                                  |
| `LABEL l2`                   | `l2:`                                  |

---

### 补充：立即数加载（当常数不适合 `MOV` 时）

| 逻辑       | 指令序列（ARM32）                          |
| -------- | ------------------------------------ |
| 加载大常数 k  | `LDR rt, =k`                         |
| 间接访问地址 k | `LDR rtmp, =k` <br> `LDR rt, [rtmp]` |

---


# 3、完整可执行程序

---

##  伪代码重构

```c
int x; 
float z;
int a[10][20]; // 初始化为 a[i][j] = i + j

float bar(int y) {
    float x;
    x = y * PI;
    return x;
}

float foo(int x, float (*boo)(int), int arr[]) {
    if (x == 0)
        z = boo(arr[1]);
    else
        return boo(arr[6 * x]);
}

print foo(2, bar, a);
```

---

##  全局设定

* 使用 `.data` 定义全局变量 `x`, `z`, `a`
* 使用 `.text` 实现 `bar`, `foo`, `main`
* 数组以线性方式存储：`a[i][j] = a[i*20 + j]`
* 调用协议遵循 MIPS 常规约定（参数 `$a0`-`$a3`、返回值 `$v0`、临时 `$t`，保存 `$s`）

---

##  MIPS 可执行程序

```mips
.data
x: .word 0
z: .float 0.0
a: .space 800           # 10*20*4 bytes
PI: .float 3.1415926
newline: .asciiz "\n"

.text
.globl main

############################
# float bar(int y)
# returns: float result in $f0
############################
bar:
    # 参数 y 在 $a0
    mtc1 $a0, $f12        # 整数转浮点
    cvt.s.w $f12, $f12
    l.s $f2, PI
    mul.s $f0, $f12, $f2  # f0 = y * PI
    jr $ra

############################
# float foo(int x, boo_ptr, int[] arr)
# 参数: $a0=x, $a1=boo_ptr, $a2=arr(base)
# 返回: $f0
############################
foo:
    beq $a0, $zero, foo_if
    # else 分支
    li $t0, 6
    mul $t1, $a0, $t0
    sll $t1, $t1, 2        # *4 字节偏移
    add $t1, $a2, $t1      # arr + 6*x
    lw $a0, 0($t1)
    jalr $a1               # call boo(arr[6*x])
    mov.s $f0, $f0
    jr $ra

foo_if:
    lw $a0, 4($a2)         # arr[1] 偏移 4 字节
    jalr $a1               # call boo(arr[1])
    s.s $f0, z             # z = boo(...)
    jr $ra

############################
# main 函数
############################
main:
    # 初始化数组 a[i][j] = i + j
    li $t0, 0          # i
outer_loop:
    li $t1, 0          # j
inner_loop:
    add $t2, $t0, $t1
    mul $t3, $t0, 20
    add $t3, $t3, $t1
    sll $t3, $t3, 2
    la $t4, a
    add $t5, $t4, $t3
    sw $t2, 0($t5)
    addi $t1, $t1, 1
    li $t6, 20
    blt $t1, $t6, inner_loop
    addi $t0, $t0, 1
    li $t6, 10
    blt $t0, $t6, outer_loop

    # 准备调用 foo(2, bar, a)
    li $a0, 2
    la $a1, bar
    la $a2, a
    jal foo

    # 输出结果
    mov.s $f12, $f0
    li $v0, 2        # syscall: print float
    syscall
    li $v0, 4
    la $a0, newline
    syscall

    li $v0, 10
    syscall
```

---