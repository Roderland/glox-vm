package glox_vm

type InterpretResult uint8

const (
	OK InterpretResult = iota
	COMPILE_ERROR
	RUNTIME_ERROR
)