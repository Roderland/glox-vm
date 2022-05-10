package chunk

import (
	"fmt"
	"math"
	"os"
)

const (
	OP_RETURN byte = iota
	OP_CONSTANT
	OP_NEGATE
	OP_ADD
	OP_SUBTRACT
	OP_MULTIPLY
	OP_DIVIDE
)

type Chunk struct {
	Codes     []byte
	Constants []Value
	Lines     []int
}

func NewChunk() *Chunk {
	return &Chunk{
		Codes: make([]byte, 0),
	}
}

func (ck *Chunk) Write(code byte, line int) {
	ck.Codes = append(ck.Codes, code)
	ck.Lines = append(ck.Lines, line)
}

func (ck *Chunk) WriteConstant(constant Value, line int) int {
	idx := len(ck.Constants)
	if idx >= math.MaxUint8 {
		fmt.Println("The number of Constants exceeds the limit 255 of one chunk.")
		os.Exit(1)
	}
	ck.Constants = append(ck.Constants, constant)
	ck.Write(OP_CONSTANT, line)
	ck.Write(uint8(idx), line)
	return idx
}
