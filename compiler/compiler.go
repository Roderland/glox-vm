package compiler

import (
	. "glox-vm"
	"math"
)

type Compiler struct {
	scanner *Scanner
	parser  *Parser
	chunk   *Chunk
	scope   *Scope
}

func (compiler *Compiler) Compile(source []byte) (*Chunk, error) {
	compiler.scanner = &Scanner{source: source, line: 1}
	compiler.parser = &Parser{}
	compiler.chunk = CreateChunk()
	compiler.scope = &Scope{compiler: compiler}
	compiler.advance()
	for !compiler.match(TOKEN_EOF) {
		compiler.declaration()
	}
	compiler.endCompile()
	return compiler.currentChunk(), compiler.parser.err
}

func (compiler *Compiler) endCompile() {
	compiler.emit(OP_RETURN)
	if compiler.parser.err == nil {
		DisassembleChunk(compiler.currentChunk(), "code")
	}
}

func (compiler *Compiler) currentChunk() *Chunk {
	return compiler.chunk
}

func (compiler *Compiler) advance() {
	compiler.parser.previous = compiler.parser.current
	for {
		compiler.parser.current = compiler.scanner.scanToken()
		if compiler.parser.current.tokenType != TOKEN_ERROR {
			break
		}
		compiler.parser.errorAtCurrent(string(compiler.parser.current.lexeme))
	}
}

func (compiler *Compiler) match(typ uint8) bool {
	if !compiler.check(typ) {
		return false
	}
	compiler.advance()
	return true
}

func (compiler *Compiler) check(typ uint8) bool {
	return compiler.parser.current.tokenType == typ
}

func (compiler *Compiler) consume(tokenType uint8, msg string) {
	if compiler.parser.current.tokenType == tokenType {
		compiler.advance()
		return
	}
	compiler.parser.errorAtCurrent(msg)
}

func (compiler *Compiler) synchronize() {
	compiler.parser.panicMode = false
	for compiler.parser.current.tokenType != TOKEN_EOF {
		if compiler.parser.current.tokenType == TOKEN_SEMICOLON {
			return
		}
		switch compiler.parser.current.tokenType {
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
		compiler.advance()
	}
}

func (compiler *Compiler) declaration() {
	if compiler.match(TOKEN_VAR) {
		compiler.varDeclaration()
	} else {
		compiler.statement()
	}
}

func (compiler *Compiler) statement() {
	if compiler.match(TOKEN_PRINT) {
		compiler.printStatement()
	} else if compiler.match(TOKEN_IF) {
		compiler.ifStatement()
	} else if compiler.match(TOKEN_LEFT_BRACE) {
		compiler.scope.begin()
		compiler.block()
		compiler.scope.end()
	} else {
		compiler.expression()
	}
}

// ifStatement if跳转的字节码实现：
// 这三个byte用于跳过OP_POP,then语句,OP_JUMP		OP_JUMP_IF_FALSE 	high-8bit low-8bit
// 将if条件值从VM栈中弹出							OP_POP
//												... then ...
// 这三个byte用于跳过OP_POP,else语句				OP_JUMP				high-8bit low-8bit
// 将if条件值从VM栈中弹出							OP_POP
// 												... else ...
func (compiler *Compiler) ifStatement() {
	compiler.consume(TOKEN_LEFT_PAREN, "Expect '(' after 'if'.")
	compiler.expression()
	compiler.consume(TOKEN_RIGHT_PAREN, "Expect ')' after condition.")
	// OP_JUMP_IF_FALSE
	thenJump := compiler.emitJump(OP_JUMP_IF_FALSE)
	// 条件为真，将if条件值从VM栈中弹出
	compiler.emit(OP_POP)
	// ... then ...
	compiler.statement()
	// OP_JUMP
	elseJump := compiler.emitJump(OP_JUMP)
	// 补充OP_JUMP_IF_FALSE的16位跳转地址
	compiler.patchJump(thenJump)
	// 条件为假，将if条件值从VM栈中弹出
	compiler.emit(OP_POP)
	// ... else ...
	if compiler.match(TOKEN_ELSE) {
		compiler.statement()
	}
	// 补充OP_JUMP的16位跳转地址
	compiler.patchJump(elseJump)
}

func (compiler *Compiler) patchJump(offset int) {
	jump := len(compiler.currentChunk().Bytecodes) - offset - 2
	if jump > math.MaxUint8 {
		compiler.parser.errorAtPrevious("Too much code to jump over.")
	}
	compiler.currentChunk().Bytecodes[offset] = uint8(jump>>8) & 0xff
	compiler.currentChunk().Bytecodes[offset+1] = uint8(jump) & 0xff
}

func (compiler *Compiler) block() {
	for !compiler.check(TOKEN_RIGHT_BRACE) && !compiler.check(TOKEN_EOF) {
		compiler.declaration()
	}
	compiler.consume(TOKEN_RIGHT_BRACE, "Expect '}' after block.")
}

func (compiler *Compiler) printStatement() {
	compiler.expression()
	compiler.consume(TOKEN_SEMICOLON, "Expect ';' after value.")
	compiler.emit(OP_PRINT)
}

func (compiler *Compiler) expressionStatement() {
	compiler.expression()
	compiler.consume(TOKEN_SEMICOLON, "Expect ';' after expression.")
	compiler.emit(OP_POP)
}
