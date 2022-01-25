package compiler

type Scanner struct {
	start   int
	current int
	line    int
	source  []byte
}

func (scanner *Scanner) scanToken() Token {
	scanner.start = scanner.current
	if scanner.isAtEnd() {
		return scanner.makeToken(TOKEN_EOF)
	}
	scanner.skipWhitespace()
	scanner.start = scanner.current
	c := scanner.advance()
	if isAlpha(c) {
		return scanner.makeIdentifier()
	}
	if isDigit(c) {
		return scanner.makeNumber()
	}
	switch c {
	case '(':
		return scanner.makeToken(TOKEN_LEFT_PAREN)
	case ')':
		return scanner.makeToken(TOKEN_RIGHT_PAREN)
	case '{':
		return scanner.makeToken(TOKEN_LEFT_BRACE)
	case '}':
		return scanner.makeToken(TOKEN_RIGHT_BRACE)
	case ';':
		return scanner.makeToken(TOKEN_SEMICOLON)
	case ',':
		return scanner.makeToken(TOKEN_COMMA)
	case '.':
		return scanner.makeToken(TOKEN_DOT)
	case '-':
		return scanner.makeToken(TOKEN_MINUS)
	case '+':
		return scanner.makeToken(TOKEN_PLUS)
	case '/':
		return scanner.makeToken(TOKEN_SLASH)
	case '*':
		return scanner.makeToken(TOKEN_STAR)
	case '!':
		if scanner.match('=') {
			scanner.makeToken(TOKEN_BANG_EQUAL)
		} else {
			scanner.makeToken(TOKEN_BANG)
		}
	case '=':
		if scanner.match('=') {
			scanner.makeToken(TOKEN_EQUAL_EQUAL)
		} else {
			scanner.makeToken(TOKEN_EQUAL)
		}
	case '<':
		if scanner.match('=') {
			scanner.makeToken(TOKEN_LESS_EQUAL)
		} else {
			scanner.makeToken(TOKEN_LESS)
		}
	case '>':
		if scanner.match('=') {
			scanner.makeToken(TOKEN_GREATER_EQUAL)
		} else {
			scanner.makeToken(TOKEN_GREATER)
		}
	case '"':
		return scanner.makeString()
	}
	return scanner.errorToken("Unexpected character.")
}

func (scanner *Scanner) isAtEnd() bool {
	return scanner.current >= len(scanner.source)
}

func (scanner *Scanner) makeToken(tokenType uint8) Token {
	return Token{
		tokenType: tokenType,
		lexeme:    scanner.source[scanner.start:scanner.current],
		line:      scanner.line,
	}
}

func (scanner *Scanner) errorToken(msg string) Token {
	return Token{
		tokenType: TOKEN_ERROR,
		lexeme:    []byte(msg),
		line:      scanner.line,
	}
}

func (scanner *Scanner) advance() (c uint8) {
	c = scanner.source[scanner.current]
	scanner.current++
	return
}

func (scanner *Scanner) match(expected uint8) bool {
	if scanner.isAtEnd() || scanner.source[scanner.current] != expected {
		return false
	}
	scanner.current++
	return true
}

func (scanner *Scanner) skipWhitespace() {
	for {
		switch scanner.peek() {
		case ' ':
			scanner.advance()
		case '\r':
			scanner.advance()
		case '\t':
			scanner.advance()
		case '\n':
			scanner.line++
			scanner.advance()
		case '/':
			if scanner.peekNext() == '/' {
				for scanner.peek() != '\n' && !scanner.isAtEnd() {
					scanner.advance()
				}
			} else {
				return
			}
		default:
			return
		}
	}
}

func (scanner *Scanner) peek() uint8 {
	return scanner.source[scanner.current]
}

func (scanner *Scanner) peekNext() uint8 {
	if scanner.isAtEnd() {
		return 0
	}
	return scanner.source[scanner.current+1]
}

func (scanner *Scanner) makeString() Token {
	for scanner.peek() != '"' && !scanner.isAtEnd() {
		if scanner.peek() == '\n' {
			scanner.line++
		}
		scanner.advance()
	}
	if scanner.isAtEnd() {
		return scanner.errorToken("Unterminated string.")
	}
	scanner.advance()
	return scanner.makeToken(TOKEN_STRING)
}

func (scanner *Scanner) makeNumber() Token {
	for isDigit(scanner.peek()) {
		scanner.advance()
	}
	if scanner.peek() == '.' && isDigit(scanner.peekNext()) {
		scanner.advance()
		for isDigit(scanner.peek()) {
			scanner.advance()
		}
	}
	return scanner.makeToken(TOKEN_NUMBER)
}

func (scanner *Scanner) makeIdentifier() Token {
	for isAlpha(scanner.peek()) || isDigit(scanner.peek()) {
		scanner.advance()
	}
	return scanner.makeToken(scanner.identifierType())
}

func (scanner *Scanner) identifierType() uint8 {
	switch scanner.source[scanner.start] {
	case 'a':
		return scanner.checkKeyword(1, 2, "nd", TOKEN_AND)
	case 'c':
		return scanner.checkKeyword(1, 4, "lass", TOKEN_CLASS)
	case 'e':
		return scanner.checkKeyword(1, 3, "lse", TOKEN_ELSE)
	case 'i':
		return scanner.checkKeyword(1, 1, "f", TOKEN_IF)
	case 'n':
		return scanner.checkKeyword(1, 2, "il", TOKEN_NIL)
	case 'o':
		return scanner.checkKeyword(1, 1, "r", TOKEN_OR)
	case 'p':
		return scanner.checkKeyword(1, 4, "rint", TOKEN_PRINT)
	case 'r':
		return scanner.checkKeyword(1, 5, "eturn", TOKEN_RETURN)
	case 's':
		return scanner.checkKeyword(1, 4, "uper", TOKEN_SUPER)
	case 'v':
		return scanner.checkKeyword(1, 2, "ar", TOKEN_VAR)
	case 'w':
		return scanner.checkKeyword(1, 4, "hile", TOKEN_WHILE)
	case 'f':
		if scanner.current-scanner.start > 1 {
			switch scanner.source[scanner.start+1] {
			case 'a':
				return scanner.checkKeyword(2, 3, "lse", TOKEN_FALSE)
			case 'o':
				return scanner.checkKeyword(2, 1, "r", TOKEN_FOR)
			case 'u':
				return scanner.checkKeyword(2, 1, "n", TOKEN_FUN)
			}
		}
	case 't':
		if scanner.current-scanner.start > 1 {
			switch scanner.source[scanner.start+1] {
			case 'h':
				return scanner.checkKeyword(2, 2, "is", TOKEN_THIS)
			case 'r':
				return scanner.checkKeyword(2, 2, "ue", TOKEN_TRUE)
			}
		}
	}

	return TOKEN_IDENTIFIER
}

func (scanner *Scanner) checkKeyword(start, length int, rest string, tokenType uint8) uint8 {
	if scanner.current-scanner.start == start+length && string(scanner.source[scanner.start+start:scanner.current]) == rest {
		return tokenType
	}
	return TOKEN_IDENTIFIER
}

func isDigit(c uint8) bool {
	return c >= '0' && c <= '9'
}

func isAlpha(c uint8) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}
