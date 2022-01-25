package compiler

import (
	"fmt"
	"os"
)

type Parser struct {
	current   Token
	previous  Token
	isError   bool
	panicMode bool
}

func (parser *Parser) errorAtCurrent(msg string) {
	parser.errorAt(&parser.current, msg)
}

func (parser *Parser) errorAtPrevious(msg string) {
	parser.errorAt(&parser.previous, msg)
}

func (parser *Parser) errorAt(token *Token, msg string) {
	if parser.panicMode {
		return
	}
	parser.panicMode = true
	fprintfError("[line %d] Error", token.line)
	if token.tokenType == TOKEN_EOF {
		fprintfError(" at end")
	} else if token.tokenType == TOKEN_ERROR {

	} else {
		fprintfError(" at '%s'", token.lexeme)
	}
	fprintfError(": %s\n", msg)
	parser.isError = true
}

func fprintfError(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format, a...)
}
