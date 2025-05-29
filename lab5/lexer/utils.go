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
func (s *Scanner) scanCharOrString(runes []rune) (Token, int) {
	if len(runes) == 0 {
		return Token{Type: TokenERROR, Lexeme: ""}, 0
	}
	if runes[0] == '\'' {
		if len(runes) < 3 {
			return Token{Type: TokenERROR, Lexeme: string(runes)}, len(runes)
		}
		if runes[1] == '\\' {
			if len(runes) < 4 || runes[3] != '\'' {
				return Token{Type: TokenERROR, Lexeme: string(runes[:min(4, len(runes))])}, min(4, len(runes))
			}
			return Token{Type: TokenCHAR, Lexeme: string(runes[:4])}, 4
		} else {
			if runes[2] != '\'' {
				return Token{Type: TokenERROR, Lexeme: string(runes[:min(3, len(runes))])}, min(3, len(runes))
			}
			return Token{Type: TokenCHAR, Lexeme: string(runes[:3])}, 3
		}
	} else if runes[0] == '"' {
		i := 1
		for i < len(runes) {
			if runes[i] == '\\' {
				i += 2
			} else if runes[i] == '"' {
				return Token{Type: TokenSTRING, Lexeme: string(runes[:i+1])}, i + 1
			} else {
				i++
			}
		}
		return Token{Type: TokenERROR, Lexeme: string(runes)}, len(runes)
	}
	return Token{Type: TokenERROR, Lexeme: string(runes[0])}, 1
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// scanComment 扫描注释，输入是runes切片，返回识别到的注释Token和长度
func (s *Scanner) scanComment(runes []rune) (Token, int) {
	if len(runes) < 2 {
		return Token{Type: TokenERROR, Lexeme: string(runes)}, len(runes)
	}

	if runes[0] == '/' && runes[1] == '/' {
		// 单行注释扫描，扫描到行尾或文件尾
		i := 2
		for i < len(runes) && runes[i] != '\n' {
			i++
		}
		return Token{Type: TokenCOMMENT_SINGLE, Lexeme: string(runes[:i])}, i
	}

	if runes[0] == '/' && runes[1] == '*' {
		// 多行注释扫描，扫描到匹配的 */
		i := 2
		for i < len(runes)-1 {
			if runes[i] == '*' && runes[i+1] == '/' {
				i += 2
				return Token{Type: TokenCOMMENT_MULTI, Lexeme: string(runes[:i])}, i
			}
			i++
		}
		// 未找到结束符，返回错误Token，长度为输入长度（到文件尾）
		return Token{Type: TokenERROR, Lexeme: string(runes)}, len(runes)
	}

	// 不是注释起始符
	return Token{Type: TokenERROR, Lexeme: string(runes[:1])}, 1
}
