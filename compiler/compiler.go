package compiler

import "fmt"

type Compiler struct {
	scanner *Scanner
}

func (compiler *Compiler) Compile(source []byte) {
	compiler.scanner = &Scanner{source: source, line: 1}
	line := -1
	for {
		token := compiler.scanner.scanToken()
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
