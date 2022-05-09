package compiler

import "github.com/Roderland/glox-vm/utils"

type Compiler struct {
	source []byte
	scn    *scanner
}

func (cpl *Compiler) compile() {
	cpl.scn = newScanner(&cpl.source)
	line := -1
	for {
		token := cpl.scn.scanToken()
		if token.line != line {
			utils.PrintfDbg("%4d ", token.line)
			line = token.line
		} else {
			utils.PrintfDbg("   | ")
		}
		utils.PrintfDbg("%2d '%s'\n", token.tp, token.lexeme)

		if token.tp == TOKEN_EOF {
			break
		}
	}
}
