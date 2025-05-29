package lexer

type DFA struct {
	Alphabet     []string                     `json:"alphabet"`
	States       []string                     `json:"states"`
	StartState   string                       `json:"start_state"`
	AcceptStates []string                     `json:"accept_states"`
	Transitions  map[string]map[string]string `json:"transitions"`

	//运行时变量
	acceptMap map[string]bool
}

type TokenType string

const (
	TokenID    TokenType = "ID"  // Identifier
	TokenNUM   TokenType = "NUM" // Number
	TokenFLO   TokenType = "FLO"
	TokenOP    TokenType = "OP"
	TokenDELIM TokenType = "DELIM"
	TokenKW    TokenType = "KEYWORD"
	TokenERROR TokenType = "ERROR"
	TokenWithespace TokenType = "WHITESPACE"
	TokenChar TokenType = "CHAR" // Character literal
	TokenString TokenType = "STRING" // String literal
	TokenEOF  TokenType = "EOF" // End of File
	TokenTYPE_KW TokenType ="TYPE_KW" // Type keyword (e.g., int, float, char, etc.)
)

type Token struct {
	Type   TokenType
	Lexeme string 
	Line  int
	Column int
}


type DFAWithTokenType struct{
	TokenType TokenType `json:"token_type"`
	DFA *DFA `json:"dfa"`
}

type TransitionTrace struct {
	From   string
	Symbol string
	To     string
}