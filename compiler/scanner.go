package compiler

type scanner struct {
	source  []byte
	start   int
	current int
	line    int
}

func (scn *scanner) init(source []byte) {
	scn.source = append(source, ' ')
	scn.line = 1
}

func (scn *scanner) scanToken() *token {
	scn.skipWhite()

	scn.start = scn.current

	if scn.isAtEnd() {
		return scn.makeToken(TOKEN_EOF)
	}

	c := scn.advance()

	if isAlpha(c) {
		return scn.ident()
	}

	if isDigit(c) {
		return scn.number()
	}

	switch c {
	case '(':
		return scn.makeToken(TOKEN_LEFT_PAREN)
	case ')':
		return scn.makeToken(TOKEN_RIGHT_PAREN)
	case '{':
		return scn.makeToken(TOKEN_LEFT_BRACE)
	case '}':
		return scn.makeToken(TOKEN_RIGHT_BRACE)
	case ';':
		return scn.makeToken(TOKEN_SEMICOLON)
	case ',':
		return scn.makeToken(TOKEN_COMMA)
	case '.':
		return scn.makeToken(TOKEN_DOT)
	case '-':
		return scn.makeToken(TOKEN_MINUS)
	case '+':
		return scn.makeToken(TOKEN_PLUS)
	case '*':
		return scn.makeToken(TOKEN_STAR)
	case '/':
		return scn.makeToken(TOKEN_SLASH)
	case '!':
		if scn.match('=') {
			return scn.makeToken(TOKEN_BANG_EQUAL)
		} else {
			return scn.makeToken(TOKEN_BANG)
		}
	case '=':
		if scn.match('=') {
			return scn.makeToken(TOKEN_EQUAL_EQUAL)
		} else {
			return scn.makeToken(TOKEN_EQUAL)
		}
	case '<':
		if scn.match('=') {
			return scn.makeToken(TOKEN_LESS_EQUAL)
		} else {
			return scn.makeToken(TOKEN_LESS)
		}
	case '>':
		if scn.match('=') {
			return scn.makeToken(TOKEN_GREATER_EQUAL)
		} else {
			return scn.makeToken(TOKEN_GREATER)
		}
	case '"':
		return scn.string()
	}

	return scn.errorToken("Unexpected character.")
}

func (scn *scanner) skipWhite() {
	for !scn.isAtEnd() {
		c := scn.peek()
		switch c {
		case ' ':
			scn.advance()
		case '\r':
			scn.advance()
		case '\t':
			scn.advance()
		case '\n':
			scn.line++
			scn.advance()
		case '/':
			if scn.peekNext() == '/' {
				// A comment goes until the end of the line.
				for !scn.isAtEnd() && scn.peek() != '\n' {
					scn.advance()
				}
			} else {
				return
			}
		default:
			return
		}
	}
}

func (scn *scanner) string() *token {
	for !scn.isAtEnd() && scn.peek() != '"' {
		if scn.peek() == '\n' {
			scn.line++
		}
		scn.advance()
	}

	if scn.isAtEnd() {
		return scn.errorToken("Unterminated string.")
	}

	scn.advance()
	return scn.makeToken(TOKEN_STRING)
}

func (scn *scanner) number() *token {
	for isDigit(scn.peek()) {
		scn.advance()
	}

	if scn.peek() == '.' && isDigit(scn.peekNext()) {
		scn.advance()
		for isDigit(scn.peek()) {
			scn.advance()
		}
	}

	return scn.makeToken(TOKEN_NUMBER)
}

func (scn *scanner) peek() byte {
	return scn.source[scn.current]
}

func (scn *scanner) peekNext() byte {
	if scn.isAtEnd() {
		return 0
	}
	return scn.source[scn.current+1]
}

func (scn *scanner) match(expected byte) bool {
	if scn.isAtEnd() || scn.peek() != expected {
		return false
	}
	scn.current++
	return true
}

func (scn *scanner) advance() byte {
	scn.current++
	return scn.source[scn.current-1]
}

func (scn *scanner) isAtEnd() bool {
	return scn.current >= len(scn.source)
}

func (scn *scanner) makeToken(tp tokenType) *token {
	var tk token
	tk.tp = tp
	tk.line = scn.line
	tk.lexeme = string(scn.source[scn.start:scn.current])
	return &tk
}

func (scn *scanner) errorToken(msg string) *token {
	var tk token
	tk.tp = TOKEN_ERROR
	tk.line = scn.line
	tk.lexeme = msg
	return &tk
}

func (scn *scanner) ident() *token {
	for isAlpha(scn.peek()) || isDigit(scn.peek()) {
		scn.advance()
	}
	return scn.makeToken(scn.identifierType())
}

func (scn *scanner) identifierType() tokenType {
	switch scn.source[scn.start] {
	case 'a':
		return scn.checkKeyword(1, 2, "nd", TOKEN_AND)
	case 'c':
		return scn.checkKeyword(1, 4, "lass", TOKEN_CLASS)
	case 'e':
		return scn.checkKeyword(1, 3, "lse", TOKEN_ELSE)
	case 'i':
		return scn.checkKeyword(1, 1, "f", TOKEN_IF)
	case 'n':
		return scn.checkKeyword(1, 2, "il", TOKEN_NIL)
	case 'o':
		return scn.checkKeyword(1, 1, "r", TOKEN_OR)
	case 'p':
		return scn.checkKeyword(1, 4, "rint", TOKEN_PRINT)
	case 'r':
		return scn.checkKeyword(1, 5, "eturn", TOKEN_RETURN)
	case 's':
		return scn.checkKeyword(1, 4, "uper", TOKEN_SUPER)
	case 'v':
		return scn.checkKeyword(1, 2, "ar", TOKEN_VAR)
	case 'w':
		return scn.checkKeyword(1, 4, "hile", TOKEN_WHILE)
	case 'f':
		if scn.current-scn.start > 1 {
			switch scn.source[scn.start+1] {
			case 'a':
				return scn.checkKeyword(2, 3, "lse", TOKEN_FALSE)
			case 'o':
				return scn.checkKeyword(2, 1, "r", TOKEN_FOR)
			case 'u':
				return scn.checkKeyword(2, 1, "n", TOKEN_FUN)
			}
		}
	case 't':
		if scn.current-scn.start > 1 {
			switch scn.source[scn.start+1] {
			case 'h':
				return scn.checkKeyword(2, 2, "is", TOKEN_THIS)
			case 'r':
				return scn.checkKeyword(2, 2, "ue", TOKEN_TRUE)
			}
		}
	}

	return TOKEN_IDENTIFIER
}

func (scn *scanner) checkKeyword(start, length int, rest string, tp tokenType) tokenType {
	if scn.current-scn.start == start+length && string(scn.source[scn.start+start:scn.current]) == rest {
		return tp
	}
	return TOKEN_IDENTIFIER
}

func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c == '_')
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}
