package compiler

import (
	. "glox-vm"
	"math"
)

func (compiler *Compiler) emitOpConstant(value Value) {
	compiler.emit(OP_CONSTANT, compiler.emitConstant(value))
}

func (compiler *Compiler) emitConstant(value Value) uint8 {
	index := compiler.currentChunk().AddConstant(value)
	if index > math.MaxUint8 {
		compiler.parser.errorAtPrevious("Too many constants in one chunk.")
	}
	return uint8(index)
}

func (compiler *Compiler) emitJump(jump byte) int {
	compiler.emit(jump, uint8(0xff), uint8(0xff))
	return len(compiler.currentChunk().Bytecodes) - 2
}

func (compiler *Compiler) emitLoop(loopStart int) {
	compiler.emit(OP_LOOP)
	offset := len(compiler.currentChunk().Bytecodes) - loopStart + 2
	if offset > math.MaxUint8 {
		compiler.parser.errorAtPrevious("Loop body too large.")
	}
	compiler.emit(uint8(offset >> 8) & 0xff, uint8(offset) & 0xff)
}

func (compiler *Compiler) emit(bytes ...byte) {
	for _, bt := range bytes {
		compiler.currentChunk().AddBytecode(bt, compiler.parser.previous.line)
	}
}
