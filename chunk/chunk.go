package chunk

import (
	"fmt"
	"math"
	"os"
)

const (
	OP_RETURN byte = iota
	OP_CONSTANT
)

type Chunk struct {
	codes     []byte
	constants []value
	lines     []int
}

func NewChunk() *Chunk {
	return &Chunk{
		codes: make([]byte, 0),
	}
}

func (ck *Chunk) Write(code byte, line int) {
	ck.codes = append(ck.codes, code)
	ck.lines = append(ck.lines, line)
}

func (ck *Chunk) WriteConstant(constant value, line int) {
	idx := len(ck.constants)
	if idx >= math.MaxUint8 {
		fmt.Println("The number of constants exceeds the limit 255.")
		os.Exit(1)
	}
	ck.constants = append(ck.constants, constant)
	ck.Write(OP_CONSTANT, line)
	ck.Write(uint8(idx), line)
}
