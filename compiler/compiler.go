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
	enclosing    *compiler
	function     *chunk.ObjFunction
	functionType chunk.FunType
	locals       [MAX_LOCAL_COUNT]local
	localCount   int
	scopeDepth   int
}

func newCompiler(functionType chunk.FunType) *compiler {
	cpl := compiler{
		enclosing:    cpl,
		function:     &chunk.ObjFunction{},
		functionType: functionType,
		locals:       [MAX_LOCAL_COUNT]local{},
		localCount:   0,
		scopeDepth:   0,
	}

	if functionType != chunk.SCRIPT {
		cpl.function.Name = prs.previous.lexeme
	}

	cpl.locals[cpl.localCount].depth = 0
	cpl.locals[cpl.localCount].name.lexeme = ""
	cpl.localCount++
	return &cpl
}

type local struct {
	name  token
	depth int
}

var scn scanner
var prs parser

// var cck *chunk.Chunk
var cpl *compiler

func Compile(source []byte, disAsmMode bool) (*chunk.ObjFunction, bool) {
	scn.init(source)
	// cck = chunk.NewChunk()
	cpl = newCompiler(chunk.SCRIPT)
	prs.hadError = false
	prs.panicMode = false
	advance()

	for !match(TOKEN_EOF) {
		declaration()
	}

	return endCompile(disAsmMode), !prs.hadError
}

func declaration() {
	if match(TOKEN_VAR) {
		varDeclaration()
	} else if match(TOKEN_FUN) {
		funDeclaration()
	} else {
		statement()
	}

	if prs.panicMode {
		synchronize()
	}
}

func funDeclaration() {
	global := parseVariable("Expect function name.")
	markInitialized()
	function(chunk.FUNCTION)
	defineVariable(global)
}

func function(ft chunk.FunType) {
	cpl = newCompiler(ft)
	beginScope()
	consume(TOKEN_LEFT_PAREN, "Expect '(' after function name.")

	if !check(TOKEN_RIGHT_PAREN) {
		cpl.function.Arity++
		if cpl.function.Arity > 255 {
			errorAtCurrent("Can't have more than 255 parameters.")
		}
		constant := parseVariable("Expect parameter name.")
		defineVariable(constant)

		for match(TOKEN_COMMA) {
			cpl.function.Arity++
			if cpl.function.Arity > 255 {
				errorAtCurrent("Can't have more than 255 parameters.")
			}
			constant := parseVariable("Expect parameter name.")
			defineVariable(constant)
		}
	}

	consume(TOKEN_RIGHT_PAREN, "Expect ')' after parameters.")
	consume(TOKEN_LEFT_BRACE, "Expect '{' before function body.")
	block()

	// endScope()
	fun := endCompile(true)
	val := chunk.NewObject(chunk.NewFunction(*fun))
	emitBytes(chunk.OP_CONSTANT, makeConstant(val))
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
	if cpl.scopeDepth == 0 {
		return
	}
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
	} else if match(TOKEN_RETURN) {
		returnStatement()
	} else {
		expressionStatement()
	}
}

func returnStatement() {
	if cpl.functionType == chunk.SCRIPT {
		errorAtPrevious("Can't return from top-level code.")
	}

	if match(TOKEN_SEMICOLON) {
		emitReturn()
	} else {
		expression()
		consume(TOKEN_SEMICOLON, "Expect ';' after return value.")
		emitBytes(chunk.OP_RETURN)
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

func endCompile(disAsmMode bool) *chunk.ObjFunction {
	emitReturn()
	function := cpl.function
	if disAsmMode {
		if !prs.hadError {
			chunk.DisAsmChunk(currentChunk(), function.GetName())
		}
	}
	cpl = cpl.enclosing
	return function
}

func emitReturn() {
	emitBytes(chunk.OP_NIL)
	emitBytes(chunk.OP_RETURN)
}

func currentChunk() *chunk.Chunk {
	return &cpl.function.Ck
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
