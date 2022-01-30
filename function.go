package glox_vm

const (
	TYPE_FUNCTION FuncType = iota
	TYPE_SCRIPT
)

type (
	FuncType uint8
	FuncData struct {
		Ft       FuncType
		Arity    int
		Name     string
		FunChunk Chunk
	}
)

func NewFuncData(ft FuncType, name string) *FuncData {
	return &FuncData{
		Ft:       ft,
		Name:     name,
		FunChunk: CreateChunk(),
	}
}
