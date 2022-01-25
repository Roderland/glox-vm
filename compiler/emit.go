package compiler

import (
	. "glox-vm"
	"math"
)

func (compiler *Compiler) emitConstant(value Value) {
	index := compiler.currentChunk().AddConstant(value)
	if index > math.MaxUint8 {
		compiler.parser.errorAtPrevious("Too many constants in one chunk.")
	}
	compiler.emit(OP_CONSTANT, uint8(index))
}

func (compiler *Compiler) emit(bytes ...byte) {
	for _, bt := range bytes {
		compiler.currentChunk().AddBytecode(bt, compiler.parser.previous.line)
	}
}
