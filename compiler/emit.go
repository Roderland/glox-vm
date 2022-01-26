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

func (compiler *Compiler) emit(bytes ...byte) {
	for _, bt := range bytes {
		compiler.currentChunk().AddBytecode(bt, compiler.parser.previous.line)
	}
}
