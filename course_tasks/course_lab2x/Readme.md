由**网安2201常兆鑫2226114409**完成
## 简单变量声明 `int x, y, z;`

### 产生式与属性文法：

```
D → T Ṽ {
  for (e ∈ Ṽ.in) bind(e, T.type);
  D.place = Ṽ.in
}

T → int     { T.type = INT }
T → float   { T.type = FLO }

Ṽ → Ṽ, d {
  v = getn(d);
  Ṽ[0].in = endcons(Ṽ[1].in, v)
}

Ṽ → d {
  v = getn(d);
  Ṽ.in = list(v)
}
```

---

##  一维数组声明 `int a[10];`

### 属性文法产生式：

```
D → T d[i] {
  x = getn(d);
  c = getv(i);
  bind(x, ARRAY);
  lookup(x, dims: 1);
  lookup(x, dim[0]: c);
  c *= sizeof(T.type);
  lookup(x, etype: T.type);
  update[width, ?c, add];
  lookup(x, base: lookup(width:));
  D.place = list(x)
}
```


  * `x` 是数组名
  * `c` 是维长（从 `i` 中获得）
  * 所有数组信息都登记进符号表
  * 更新符号表表头宽度（`width`）
  * `place` 记录声明的标识符集合

---

## 函数声明 `int foo(int x){int y; Š}`

### 属性文法产生式：

```
D → T d(Ǎ) { Ď Š } {
  tab = pop(symtab);
  tab->outer = top(symtab);
  x = getn(d);
  bind(x, FUNC);
  lookup(x, mytab: tab);
  update[width, ?sizeof(FUNC), add];
  lookup(x, offset: lookup(width:));
  D.place = list(x);
  tab->code = Š.code;
}
```

#### 参数列表处理：

```
Ǎ → ε {
  Ǎ.place = NIL;
  push(symtab, newtab());
}

A → T d {
  x = getn(d);
  bind(x, T.type);
  update[width, ?sizeof(T.type), add];
  lookup(x, offset: lookup(width:));
  update[arglist, ?x, endcons];
  update[argc, 1, add];
  A.place = list(x);
}

Ǎ → Ǎ A {
  Ǎ[0].place = append(A.place, Ǎ[1].place);
}
```

---



