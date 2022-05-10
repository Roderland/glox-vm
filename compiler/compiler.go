package compiler

import (
	"github.com/Roderland/glox-vm/chunk"
	"github.com/Roderland/glox-vm/utils"
)

type parser struct {
	current   *token
	previous  *token
	hadError  bool
	panicMode bool
}

var scn scanner
var prs parser
var cck *chunk.Chunk

func Compile(ck *chunk.Chunk, source []byte, disAsmMode bool) bool {
	scn.init(source)
	cck = ck
	prs.hadError = false
	prs.panicMode = false
	advance()
	expression()
	consume(TOKEN_EOF, "Expect end of expression.")
	endCompile()
	if disAsmMode {
		if !prs.hadError {
			chunk.DisAsmChunk(ck, "compile result")
		}
	}
	return !prs.hadError
}

func endCompile() {
	emitBytes(chunk.OP_RETURN)
}

func currentChunk() *chunk.Chunk {
	return cck
}

func emitConstant(value chunk.Value) {
	makeConstant(value)
}

func makeConstant(value chunk.Value) {
	currentChunk().WriteConstant(value, prs.previous.line)
}

func emitBytes(bts ...byte) {
	c := currentChunk()
	for _, bt := range bts {
		c.Write(bt, prs.previous.line)
	}
}

func advance() {
	prs.previous = prs.current

	for {
		prs.current = scn.scanToken()
		if prs.current.tp != TOKEN_ERROR {
			break
		}
		errorAtCurrent(prs.current.lexeme)
	}
}

func consume(tp tokenType, msg string) {
	if prs.current.tp == tp {
		advance()
		return
	}

	errorAtCurrent(msg)
}

func errorAtCurrent(msg string) {
	errorAt(prs.current, msg)
}

func errorAtPrevious(msg string) {
	errorAt(prs.previous, msg)
}

func errorAt(tk *token, msg string) {
	if prs.panicMode {
		return
	} else {
		prs.panicMode = true
	}

	utils.PrintfErr("[line %d] Error", tk.line)

	if tk.tp == TOKEN_EOF {
		utils.PrintfErr(" at end")
	} else if tk.tp == TOKEN_ERROR {
		// Nothing.
	} else {
		utils.PrintfErr(" at '%s'", tk.lexeme)
	}

	utils.PrintfErr(": %s\n", msg)
	prs.hadError = true
}
