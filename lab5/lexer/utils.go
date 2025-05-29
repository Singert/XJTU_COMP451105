package lexer


// util function to check if a string exists in a slice
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
// 辅助识别char和string
func (s *Scanner) scanCharOrString(input string) (Token, int) {
	if len(input) == 0 {
		return Token{Type: TokenERROR, Lexeme: ""}, 0
	}
	runes := []rune(input)
	if runes[0] == '\'' {
		// 处理 char 字符常量
		// 格式示例：'a', '\n', '\'', '\\'
		if len(runes) < 3 {
			return Token{Type: TokenERROR, Lexeme: input}, len(runes)
		}
		// 简单处理转义字符情况
		if runes[1] == '\\' {
			if len(runes) < 4 || runes[3] != '\'' {
				return Token{Type: TokenERROR, Lexeme: string(runes[:min(4, len(runes))])}, min(4, len(runes))
			}
			return Token{Type: "CHAR", Lexeme: string(runes[:4])}, 4
		} else {
			if runes[2] != '\'' {
				return Token{Type: TokenERROR, Lexeme: string(runes[:min(3, len(runes))])}, min(3, len(runes))
			}
			return Token{Type: "CHAR", Lexeme: string(runes[:3])}, 3
		}
	} else if runes[0] == '"' {
		// 处理 string 字符串常量
		i := 1
		for i < len(runes) {
			if runes[i] == '\\' {
				// 跳过转义符及下一个字符
				i += 2
			} else if runes[i] == '"' {
				// 结束引号
				return Token{Type: "STRING", Lexeme: string(runes[:i+1])}, i + 1
			} else {
				i++
			}
		}
		// 没有找到结束引号，错误
		return Token{Type: TokenERROR, Lexeme: input}, len(runes)
	}
	return Token{Type: TokenERROR, Lexeme: string(runes[0])}, 1
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}