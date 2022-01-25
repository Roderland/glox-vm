package compiler

import (
	. "glox-vm"
)

type Compiler struct {
	scanner *Scanner
	parser  *Parser
	chunk   *Chunk
}

func (compiler *Compiler) Compile(source []byte) (*Chunk, error) {
	defer compiler.endCompile()
	compiler.scanner = &Scanner{source: source, line: 1}
	compiler.parser = &Parser{}
	compiler.chunk = CreateChunk()
	compiler.advance()
	compiler.expression()
	compiler.consume(TOKEN_EOF, "Expect end of expression.");
	return compiler.currentChunk(), nil
}

func (compiler *Compiler) endCompile() {
	compiler.emit(OP_RETURN)
	if !compiler.parser.isError {
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

func (compiler *Compiler) consume(tokenType uint8, msg string) {
	if compiler.parser.current.tokenType == tokenType {
		compiler.advance()
		return
	}
	compiler.parser.errorAtCurrent(msg)
}


