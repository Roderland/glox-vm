package glox_vm

const (
	OP_RETURN byte = iota
	OP_CONSTANT
	OP_NEGATE
	OP_ADD
	OP_SUBTRACT
	OP_MULTIPLY
	OP_DIVIDE
	OP_NIL
	OP_TRUE
	OP_FALSE
	OP_NOT
	OP_EQUAL
	OP_GREATER
	OP_LESS
	OP_PRINT
	OP_POP
	OP_DEFINE_GLOBAL
	OP_GET_GLOBAL
	OP_SET_GLOBAL
	OP_GET_LOCAL
	OP_SET_LOCAL
)

type Chunk struct {
	Bytecodes []byte
	Lines     []int
	Constants []Value
}

func CreateChunk() *Chunk {
	return &Chunk{
		Bytecodes: []byte{},
		Lines:     []int{},
		Constants: []Value{},
	}
}

func (chunk *Chunk) AddBytecode(bytecode byte, line int) {
	chunk.Bytecodes = append(chunk.Bytecodes, bytecode)
	chunk.Lines = append(chunk.Lines, line)
}

func (chunk *Chunk) AddConstant(constant Value) int {
	chunk.Constants = append(chunk.Constants, constant)
	return len(chunk.Constants) - 1
}
