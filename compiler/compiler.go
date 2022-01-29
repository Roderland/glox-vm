package compiler

import (
	. "glox-vm"
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
	} else if compiler.match(TOKEN_LEFT_BRACE) {
		compiler.scope.begin()
		compiler.block()
		compiler.scope.end()
	} else {
		compiler.expression()
	}
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
