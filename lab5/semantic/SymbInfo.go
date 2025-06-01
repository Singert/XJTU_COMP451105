package semantic

import (
    "fmt"
    "lab5/syntax"
)

// SymbolInfo 表示符号在符号表中的详细信息
type SymbolInfo struct {
    Type      string // 符号类型（例如 int, float）
    IsInit    bool   // 是否已初始化
    IsConstant bool  // 是否为常量
}

// Scope 表示一个作用域
type Scope struct {
    Symbols map[syntax.Symbol]*SymbolInfo // 当前作用域的符号集合，使用Symbol类型作为符号名
}

// SymbolTable 管理所有作用域
type SymbolTable struct {
    scopes []*Scope // 作用域栈
}

// NewSymbolTable 创建一个新的符号表
func NewSymbolTable() *SymbolTable {
    return &SymbolTable{
        scopes: make([]*Scope, 0),
    }
}

// NewScope 创建一个新的作用域
func NewScope() *Scope {
    return &Scope{
        Symbols: make(map[syntax.Symbol]*SymbolInfo),
    }
}

// EnterScope 进入一个新的作用域
func (st *SymbolTable) EnterScope() {
    st.scopes = append(st.scopes, NewScope())
}

// ExitScope 退出当前作用域
func (st *SymbolTable) ExitScope() {
    if len(st.scopes) > 0 {
        st.scopes = st.scopes[:len(st.scopes)-1]
    }
}

// DeclareSymbol 声明一个符号并插入到当前作用域
func (st *SymbolTable) DeclareSymbol(name syntax.Symbol, typ string, isConstant bool) error {
    if len(st.scopes) == 0 {
        return fmt.Errorf("no scope to declare symbol")
    }

    currentScope := st.scopes[len(st.scopes)-1]

    // 检查符号是否已声明
    if _, exists := currentScope.Symbols[name]; exists {
        return fmt.Errorf("symbol %s already declared in current scope", name)
    }

    // 添加符号到当前作用域
    currentScope.Symbols[name] = &SymbolInfo{
        Type:      typ,
        IsInit:    false, // 初始时，变量未初始化
        IsConstant: isConstant,
    }
    return nil
}

// LookupSymbol 查找符号，支持多层作用域查找
func (st *SymbolTable) LookupSymbol(name syntax.Symbol) (*SymbolInfo, error) {
    for i := len(st.scopes) - 1; i >= 0; i-- {
        scope := st.scopes[i]
        if symbolInfo, exists := scope.Symbols[name]; exists {
            return symbolInfo, nil
        }
    }
    return nil, fmt.Errorf("symbol %s not found", name)
}

// MarkSymbolInitialized 标记符号已初始化
func (st *SymbolTable) MarkSymbolInitialized(name syntax.Symbol) error {
    symbolInfo, err := st.LookupSymbol(name)
    if err != nil {
        return err
    }
    symbolInfo.IsInit = true
    return nil
}



