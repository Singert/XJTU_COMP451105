package token

type TokenType string

const (
    // 基本类型
    ID    TokenType = "ID"
    NUM   TokenType = "NUM"

    // 关键字
    INT    TokenType = "INT"
    VOID   TokenType = "VOID"
    IF     TokenType = "IF"
    ELSE   TokenType = "ELSE"
    WHILE  TokenType = "WHILE"
    RETURN TokenType = "RETURN"

    // 分隔符
    SCO TokenType = "SCO" // ;
    CMA TokenType = "CMA" // ,
    LBR TokenType = "LBR" // {
    RBR TokenType = "RBR" // }
    LPA TokenType = "LPA" // (
    RPA TokenType = "RPA" // )

    // 运算符
    ADD TokenType = "ADD" // +
    MUL TokenType = "MUL" // *
    AND TokenType = "AND" // &&
    OR  TokenType = "OR"  // ||
    ROP TokenType = "ROP" // <, <=, == 等

    // 特殊
    ILLEGAL TokenType = "ILLEGAL"
    EOF     TokenType = "EOF"
)

type Token struct {
    Type    TokenType
    Lexeme  string
    Line    int
    Column  int
}
