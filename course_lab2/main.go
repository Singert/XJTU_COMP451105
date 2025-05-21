package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"
)

type Production struct {
	Head   string
	Body   []string
	Action func(stack *ParserStack)
}

// 用于属性求值
type Attribute struct {
	Type  VarType
	Name  string
	Dim   []int
	Place []string
}

type VarType int

const (
	INT VarType = iota
	FLO
	ARRAY
	FUNC
)

type SymbolEntry struct {
	Name   string
	Type   VarType
	EType  VarType // ARRAY/FUNC 专用
	Dims   int
	Dim    []int
	Offset int
	Args   []string
	Inner  *SymbolTable
}

type SymbolTable struct {
	Outer   *SymbolTable
	Width   int
	Entries map[string]*SymbolEntry
}

func NewSymbolTable(outer *SymbolTable) *SymbolTable {
	return &SymbolTable{
		Outer:   outer,
		Width:   0,
		Entries: make(map[string]*SymbolEntry),
	}
}

func (t *SymbolTable) Bind(name string, entry *SymbolEntry) {
	if _, exists := t.Entries[name]; exists {
		panic(fmt.Sprintf("symbol %s already declared", name))
	}
	t.Entries[name] = entry
}

func (t *SymbolTable) Lookup(name string) *SymbolEntry {
	for table := t; table != nil; table = table.Outer {
		if entry, ok := table.Entries[name]; ok {
			return entry
		}
	}
	return nil
}

type ParserStack struct {
	StateStack    []int
	SymbolStack   []string
	AttrStack     []Attribute
	Tokens        []Token
	Pos           int
	SymbolTable   *SymbolTable
	GlobalSymbols *SymbolTable
}

func NewParser(tokens []Token) *ParserStack {
	global := NewSymbolTable(nil)
	return &ParserStack{
		StateStack:    []int{0},
		SymbolStack:   []string{"#"},
		AttrStack:     []Attribute{{}},
		Tokens:        tokens,
		Pos:           0,
		GlobalSymbols: global,
		SymbolTable:   global,
	}
}

// 模拟产生式：D → T ID ;
func (p *ParserStack) ReduceVariableDecl(tType VarType) {
	if len(p.AttrStack) < 1 {
		panic("ReduceVariableDecl: AttrStack too short")
	}

	ident := p.AttrStack[len(p.AttrStack)-1].Name

	p.SymbolTable.Bind(ident, &SymbolEntry{
		Name:   ident,
		Type:   tType,
		Offset: p.SymbolTable.Width,
	})
	size := 4
	if tType == FLO {
		size = 8
	}
	p.SymbolTable.Width += size

	// 仅弹出变量名一个属性
	p.SymbolStack = p.SymbolStack[:len(p.SymbolStack)-1]
	p.AttrStack = p.AttrStack[:len(p.AttrStack)-1]

	// 写入归约结果
	p.SymbolStack = append(p.SymbolStack, "D")
	p.AttrStack = append(p.AttrStack, Attribute{
		Type:  tType,
		Name:  ident,
		Place: []string{ident},
	})
}

func (p *ParserStack) ReduceArrayDecl(tType VarType, size int) {
	ident := p.AttrStack[len(p.AttrStack)-1].Name // ✅ 修改

	p.SymbolTable.Bind(ident, &SymbolEntry{
		Name:   ident,
		Type:   ARRAY,
		EType:  tType,
		Dims:   1,
		Dim:    []int{size},
		Offset: p.SymbolTable.Width,
	})
	elemSize := 4
	if tType == FLO {
		elemSize = 8
	}
	p.SymbolTable.Width += elemSize * size

	// 退栈变量名
	p.SymbolStack = p.SymbolStack[:len(p.SymbolStack)-1]
	p.AttrStack = p.AttrStack[:len(p.AttrStack)-1]

	p.SymbolStack = append(p.SymbolStack, "D")
	p.AttrStack = append(p.AttrStack, Attribute{
		Type:  ARRAY,
		Name:  ident,
		Place: []string{ident},
	})
}

// 简化版语法分析主循环
func (p *ParserStack) Parse() {
	for {
		tok := p.Tokens[p.Pos]

		switch tok.Type {
		case INT_KW, FLOAT_KW:
			var vt VarType
			if tok.Type == INT_KW {
				vt = INT
			} else {
				vt = FLO
			}
			p.Pos++

			// 读取变量名
			id := p.Tokens[p.Pos]
			if id.Type != IDENT {
				panic("expect identifier after type")
			}
			p.AttrStack = append(p.AttrStack, Attribute{Name: id.Value})
			p.Pos++

			// 分支：变量 or 数组
			switch p.Tokens[p.Pos].Type {
			case SEMI:
				p.Pos++
				p.ReduceVariableDecl(vt)

			case LBRACK:
				p.Pos++
				sizeTok := p.Tokens[p.Pos]
				if sizeTok.Type != CONST {
					panic("expect constant inside array brackets")
				}
				sz, err := strconv.Atoi(sizeTok.Value)
				if err != nil {
					panic("invalid integer constant")
				}
				p.Pos++

				if p.Tokens[p.Pos].Type != RBRACK {
					panic("missing closing bracket ]")
				}
				p.Pos++

				if p.Tokens[p.Pos].Type != SEMI {
					panic("missing semicolon after array declaration")
				}
				p.Pos++

				p.ReduceArrayDecl(vt, sz)

			default:
				panic(fmt.Sprintf("unexpected token after identifier: %v", p.Tokens[p.Pos]))
			}

		case EOF:
			return

		default:
			panic(fmt.Sprintf("unexpected token: %v", tok))
		}
	}
}

type TokenType int

const (
	INT_KW TokenType = iota + 100
	FLOAT_KW
	IDENT
	LPAREN
	RPAREN
	LBRACK
	RBRACK
	COMMA
	SEMI
	EOF
	CONST
)

type Token struct {
	Type  TokenType
	Value string
}

func Lex(input string) []Token {
	var tokens []Token
	i := 0
	for i < len(input) {
		c := input[i]
		if unicode.IsSpace(rune(c)) {
			i++
			continue
		}
		switch {
		case strings.HasPrefix(input[i:], "int") && (i+3 == len(input) || !unicode.IsLetter(rune(input[i+3])) && !unicode.IsDigit(rune(input[i+3]))):
			tokens = append(tokens, Token{INT_KW, "int"})
			i += 3
		case strings.HasPrefix(input[i:], "float") && (i+5 == len(input) || !unicode.IsLetter(rune(input[i+5])) && !unicode.IsDigit(rune(input[i+5]))):
			tokens = append(tokens, Token{FLOAT_KW, "float"})
			i += 5
		case unicode.IsLetter(rune(c)):
			j := i
			for j < len(input) && (unicode.IsLetter(rune(input[j])) || unicode.IsDigit(rune(input[j]))) {
				j++
			}
			tokens = append(tokens, Token{IDENT, input[i:j]})
			i = j
		case unicode.IsDigit(rune(c)):
			j := i
			for j < len(input) && unicode.IsDigit(rune(input[j])) {
				j++
			}
			tokens = append(tokens, Token{CONST, input[i:j]})
			i = j
		case c == '[':
			tokens = append(tokens, Token{LBRACK, "["})
			i++
		case c == ']':
			tokens = append(tokens, Token{RBRACK, "]"})
			i++
		case c == '(':
			tokens = append(tokens, Token{LPAREN, "("})
			i++
		case c == ')':
			tokens = append(tokens, Token{RPAREN, ")"})
			i++
		case c == ',':
			tokens = append(tokens, Token{COMMA, ","})
			i++
		case c == ';':
			tokens = append(tokens, Token{SEMI, ";"})
			i++
		default:
			panic("Invalid character: " + string(c))
		}
	}
	tokens = append(tokens, Token{EOF, ""})
	return tokens
}

func main() {
	// 打开输入文件
	path := filepath.Join("test", "input.txt")
	content, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("无法读取输入文件: %v\n", err)
		os.Exit(1)
	}

	// 词法分析
	tokens := Lex(string(content))

	// 初始化分析器并运行
	parser := NewParser(tokens)
	parser.Parse()

	// 输出符号表
	fmt.Println("符号表构造完成：")
	for name, entry := range parser.GlobalSymbols.Entries {
		fmt.Printf("  name: %s, type: %v, offset: %d\n", name, entry.Type, entry.Offset)
		if entry.Type == ARRAY {
			fmt.Printf("    dims: %d, dim: %v, etype: %v\n", entry.Dims, entry.Dim, entry.EType)
		}
	}
}
