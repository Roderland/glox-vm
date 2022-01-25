package glox_vm

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
	Bytecodes []byte
	lines     []int
	Constants []Value
}

func CreateChunk() *Chunk {
	return &Chunk{
		Bytecodes: []byte{},
		lines:     []int{},
		Constants: []Value{},
	}
}

func (chunk *Chunk) AddBytecode(bytecode byte, line int) {
	chunk.Bytecodes = append(chunk.Bytecodes, bytecode)
	chunk.lines = append(chunk.lines, line)
}

func (chunk *Chunk) AddConstant(constant Value) int {
	chunk.Constants = append(chunk.Constants, constant)
	return len(chunk.Constants) - 1
}
