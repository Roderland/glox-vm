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

const MAX_LOCAL_COUNT = math.MaxUint8 + 1

type compiler struct {
	locals     [MAX_LOCAL_COUNT]local
	localCount int
	scopeDepth int
}

func newCompiler() *compiler {
	return &compiler{
		locals:     [MAX_LOCAL_COUNT]local{},
		localCount: 0,
		scopeDepth: 0,
	}
}

type local struct {
	name  token
	depth int
}

var scn scanner
var prs parser
var cck *chunk.Chunk
var cpl *compiler

func Compile(ck *chunk.Chunk, source []byte, disAsmMode bool) bool {
	scn.init(source)
	cck = ck
	cpl = newCompiler()
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
	declareVariable()
	if cpl.scopeDepth > 0 {
		return 0
	}
	return identifierConstant(prs.previous)
}

func declareVariable() {
	if cpl.scopeDepth == 0 {
		return
	}

	name := *prs.previous

	for i := cpl.localCount - 1; i >= 0; i-- {
		lc := &cpl.locals[i]
		if lc.depth != -1 && lc.depth < cpl.scopeDepth {
			break
		}
		if lc.name.lexeme == name.lexeme {
			errorAtPrevious("Already a variable with this name in this scope.")
		}
	}
	addLocal(name)
}

func addLocal(tk token) {
	if cpl.localCount == MAX_LOCAL_COUNT {
		errorAtPrevious("Too many local variables in function.")
	}
	cpl.locals[cpl.localCount].name = tk
	cpl.locals[cpl.localCount].depth = -1
	cpl.localCount++
}

func identifierConstant(varName *token) uint8 {
	return makeConstant(chunk.NewString(varName.lexeme))
}

func defineVariable(global uint8) {
	if cpl.scopeDepth > 0 {
		markInitialized()
		return
	}
	emitBytes(chunk.OP_DEFINE_GLOBAL, global)
}

func markInitialized() {
	cpl.locals[cpl.localCount-1].depth = cpl.scopeDepth
}

func statement() {
	if match(TOKEN_PRINT) {
		printStatement()
	} else if match(TOKEN_LEFT_BRACE) {
		beginScope()
		block()
		endScope()
	} else if match(TOKEN_IF) {
		ifStatement()
	} else if match(TOKEN_WHILE) {
		whileStatement()
	} else if match(TOKEN_FOR) {
		forStatement()
	} else {
		expressionStatement()
	}
}

func forStatement() {
	beginScope()
	consume(TOKEN_LEFT_PAREN, "Expect '(' after 'for'.")
	if match(TOKEN_SEMICOLON) {
		// No initializer.
	} else if match(TOKEN_VAR) {
		varDeclaration()
	} else {
		expressionStatement()
	}

	loopStart := len(currentChunk().Codes)
	incrStart := loopStart
	exitJump := -1
	if !match(TOKEN_SEMICOLON) {
		expression()
		consume(TOKEN_SEMICOLON, "Expect ';'.")

		exitJump = emitJump(chunk.OP_JUMP_IF_FALSE)
		emitBytes(chunk.OP_POP)
	}

	if !match(TOKEN_RIGHT_PAREN) {
		bodyJump := emitJump(chunk.OP_JUMP)
		incrStart = len(currentChunk().Codes)
		expression()
		emitBytes(chunk.OP_POP)
		consume(TOKEN_RIGHT_PAREN, "Expect ')' after for clauses.")

		emitLoop(loopStart)
		loopStart = incrStart
		patchJump(bodyJump)
	}

	statement()
	emitLoop(incrStart)

	if exitJump != -1 {
		patchJump(exitJump)
		emitBytes(chunk.OP_POP)
	}

	endScope()
}

func whileStatement() {
	consume(TOKEN_LEFT_PAREN, "Expect '(' after 'while'.")
	loopStart := len(currentChunk().Codes)
	expression()
	consume(TOKEN_RIGHT_PAREN, "Expect ')' after condition.")

	exitJump := emitJump(chunk.OP_JUMP_IF_FALSE)
	emitBytes(chunk.OP_POP)
	statement()
	emitLoop(loopStart)

	patchJump(exitJump)
	emitBytes(chunk.OP_POP)
}

func emitLoop(loopStart int) {
	emitBytes(chunk.OP_LOOP)

	offset := len(currentChunk().Codes) - loopStart + 2
	if offset > math.MaxUint16 {
		errorAtPrevious("Loop body too large.")
	}

	emitBytes(uint8(offset>>8) & 0xff)
	emitBytes(uint8(offset) & 0xff)
}

func ifStatement() {
	consume(TOKEN_LEFT_PAREN, "Expect '(' after 'if'.")
	expression()
	consume(TOKEN_RIGHT_PAREN, "Expect ')' after condition.")

	thenJump := emitJump(chunk.OP_JUMP_IF_FALSE)
	emitBytes(chunk.OP_POP)
	statement()

	elseJump := emitJump(chunk.OP_JUMP)
	patchJump(thenJump)
	emitBytes(chunk.OP_POP)

	if match(TOKEN_ELSE) {
		statement()
	}
	patchJump(elseJump)
}

func emitJump(instruction byte) int {
	emitBytes(instruction, 0xff, 0xff)
	return len(currentChunk().Codes) - 2
}

func patchJump(offset int) {
	jump := len(currentChunk().Codes) - offset - 2
	if jump > math.MaxUint16 {
		errorAtPrevious("Too much code to jump over.")
	}

	currentChunk().Codes[offset] = byte((jump >> 8) & 0xff)
	currentChunk().Codes[offset+1] = byte(jump & 0xff)
}

func beginScope() {
	cpl.scopeDepth++
}
func endScope() {
	cpl.scopeDepth--

	for cpl.localCount > 0 && cpl.locals[cpl.localCount-1].depth > cpl.scopeDepth {
		emitBytes(chunk.OP_POP)
		cpl.localCount--
	}
}
func block() {
	for !check(TOKEN_EOF) && !check(TOKEN_RIGHT_BRACE) {
		declaration()
	}
	consume(TOKEN_RIGHT_BRACE, "Expect '}' after block.")
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
