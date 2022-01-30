package compiler

import (
	"fmt"
	"testing"
)

func TestScanner(t *testing.T) {
	source := "fun a() {\n    b();\n}"
	scanner := &Scanner{source: []byte(source), line: 1}
	line := -1
	for {
		token := scanner.scanToken()
		if token.line != line {
			fmt.Printf("%4d ", token.line)
			line = token.line
		} else {
			fmt.Print("   | ")
		}
		fmt.Printf("%2d '%.*s'\n", token.tokenType, len(token.lexeme), token.lexeme)
		if token.tokenType == TOKEN_EOF {
			break
		}
	}
}
