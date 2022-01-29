package compiler

import (
	"fmt"
	"testing"
)

func TestScanner(t *testing.T) {
	source := "var a = 0;\nwhile (a < 3) {\n    print \"a\";\n    a = a + 1;\n}"
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
