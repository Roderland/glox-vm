package compiler

import (
	"fmt"
	"github.com/Roderland/glox-vm/chunk"
	"github.com/Roderland/glox-vm/utils"
	"math"
	"os"
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
	//expression()
	//consume(TOKEN_EOF, "Expect end of expression.")
	for !match(TOKEN_EOF) {
		declaration()
	}
	endCompile()
	if disAsmMode {
		if !prs.hadError {
			chunk.DisAsmChunk(ck, "compile result")
		}
	}
	return !prs.hadError
}

func declaration() {
	if match(TOKEN_VAR) {
		varDeclaration()
	} else {
		statement()
	}

	if prs.panicMode {
		synchronize()
	}
}

func varDeclaration() {
	global := parseVariable("Expect variable name.")

	if match(TOKEN_EQUAL) {
		expression()
	} else {
		emitBytes(chunk.OP_NIL)
	}
	consume(TOKEN_SEMICOLON, "Expect ';' after variable declaration.")

	defineVariable(global)
}

func parseVariable(errorMessage string) byte {
	consume(TOKEN_IDENTIFIER, errorMessage)
	return identifierConstant(prs.previous)
}

func identifierConstant(varName *token) uint8 {
	return makeConstant(chunk.NewString(varName.lexeme))
}

func defineVariable(global uint8) {
	emitBytes(chunk.OP_DEFINE_GLOBAL, global)
}

func statement() {
	if match(TOKEN_PRINT) {
		printStatement()
	} else {
		expressionStatement()
	}
}

func expressionStatement() {
	expression()
	consume(TOKEN_SEMICOLON, "Expect ';' after expression.")
	emitBytes(chunk.OP_POP)
}

func printStatement() {
	expression()
	consume(TOKEN_SEMICOLON, "Expect ';' after value.")
	emitBytes(chunk.OP_PRINT)
}

func match(tp tokenType) bool {
	if check(tp) {
		advance()
		return true
	}
	return false
}

func check(tp tokenType) bool {
	return prs.current.tp == tp
}

func endCompile() {
	emitBytes(chunk.OP_RETURN)
}

func currentChunk() *chunk.Chunk {
	return cck
}

func emitConstant(value chunk.Value) uint8 {
	idx := makeConstant(value)
	emitBytes(chunk.OP_CONSTANT, idx)
	return idx
}

func makeConstant(value chunk.Value) uint8 {
	idx := currentChunk().AddConstant(value)
	if idx >= math.MaxUint8 {
		fmt.Println("The number of Constants exceeds the limit 255 of one chunk.")
		os.Exit(1)
	}
	idx8 := uint8(idx)
	return idx8
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

func synchronize() {
	prs.panicMode = false

	for prs.current.tp != TOKEN_EOF {
		if prs.previous.tp == TOKEN_SEMICOLON {
			return
		}
		switch prs.current.tp {
		case TOKEN_CLASS:
			return
		case TOKEN_FUN:
			return
		case TOKEN_VAR:
			return
		case TOKEN_FOR:
			return
		case TOKEN_IF:
			return
		case TOKEN_WHILE:
			return
		case TOKEN_PRINT:
			return
		case TOKEN_RETURN:
			return
		default:
		}

		advance()
	}

}
