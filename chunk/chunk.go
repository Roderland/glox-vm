package chunk

const (
	OP_RETURN byte = iota
	OP_CONSTANT
	OP_NEGATE
	OP_ADD
	OP_SUBTRACT
	OP_MULTIPLY
	OP_DIVIDE
	OP_NIL
	OP_FALSE
	OP_TRUE
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

func (ck *Chunk) AddConstant(constant Value) int {
	ck.Constants = append(ck.Constants, constant)
	return len(ck.Constants) - 1
}
