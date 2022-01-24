package main

type Value float64

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
	bytecodes []byte
	lines     []int
	constants []Value
}

func _Chunk() *Chunk {
	return &Chunk{
		bytecodes: []byte{},
		lines:     []int{},
		constants: []Value{},
	}
}

func (chunk *Chunk) addBytecode(bytecode byte, line int) {
	chunk.bytecodes = append(chunk.bytecodes, bytecode)
	chunk.lines = append(chunk.lines, line)
}

func (chunk *Chunk) addConstant(constant Value) int {
	chunk.constants = append(chunk.constants, constant)
	return len(chunk.constants) - 1
}
