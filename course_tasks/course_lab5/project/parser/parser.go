// parser/parser.go
package parser

import (
	"project/stmt" // 新增
)

func ParseAndGenerateTAC(tokens []string) []string {
	return stmt.Dispatch(tokens)
}
