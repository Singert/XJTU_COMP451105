package dfa

import (
    "unicode"
)

type StateID int

type Transition struct {
    From  StateID
    Input rune
    To    StateID
}

type DFA struct {
    Start     StateID
    Accepting map[StateID]string         // 状态 -> TokenType 字符串
    Trans     map[StateID]map[rune]StateID
}

// 创建空DFA
func NewDFA() *DFA {
    return &DFA{
        Accepting: make(map[StateID]string),
        Trans:     make(map[StateID]map[rune]StateID),
    }
}

// 添加一个状态转移
func (d *DFA) AddTransition(from StateID, input rune, to StateID) {
    if d.Trans[from] == nil {
        d.Trans[from] = make(map[rune]StateID)
    }
    d.Trans[from][input] = to
}

// 获取从某状态出发的转移
func (d *DFA) Next(from StateID, input rune) (StateID, bool) {
    to, ok := d.Trans[from][input]
    return to, ok
}

// 设置接受状态及对应的 token 类型
func (d *DFA) SetAccepting(state StateID, tokenType string) {
    d.Accepting[state] = tokenType
}

// 判断是否是接受状态
func (d *DFA) IsAccepting(state StateID) (string, bool) {
    tokenType, ok := d.Accepting[state]
    return tokenType, ok
}

// 工具函数：构造匹配标识符/关键字的DFA
func BuildIDOrKeywordDFA(keywords map[string]string) *DFA {
    d := NewDFA()
    var state StateID = 0
    d.Start = state

    // 状态0：字母开头
    state++
    d.AddTransition(0, '_', state)
    for ch := 'a'; ch <= 'z'; ch++ {
        d.AddTransition(0, ch, state)
        d.AddTransition(0, unicode.ToUpper(ch), state)
    }

    // 状态1+：字母数字下划线循环
    d.SetAccepting(state, "ID") // 默认认为是ID，后续在外部检查是否为关键字
    for i := state; i < state+10; i++ {
        for ch := 'a'; ch <= 'z'; ch++ {
            d.AddTransition(i, ch, i+1)
            d.AddTransition(i, unicode.ToUpper(ch), i+1)
        }
        for ch := '0'; ch <= '9'; ch++ {
            d.AddTransition(i, ch, i+1)
        }
        d.AddTransition(i, '_', i+1)
        d.SetAccepting(i+1, "ID")
    }

    return d
}

// 工具函数：构造匹配数字的DFA（仅十进制整数）
func BuildNumberDFA() *DFA {
    d := NewDFA()
    d.Start = 0
    d.AddTransition(0, '0', 1)
    d.SetAccepting(1, "NUM")
    for ch := '1'; ch <= '9'; ch++ {
        d.AddTransition(0, ch, 2)
        d.AddTransition(2, ch, 2)
        d.AddTransition(2, '0', 2)
    }
    d.SetAccepting(2, "NUM")
    return d
}

