package vm

import (
	"fmt"
	"github.com/Roderland/glox-vm/chunk"
)

type Result uint8

const (
	OK Result = iota
	ERROR
)

type VM struct {
	ck        *chunk.Chunk
	ip        int
	stack     []chunk.Value
	debugMode bool
}

func NewVM(ck *chunk.Chunk, debugMode bool) *VM {
	return &VM{
		ck:        ck,
		debugMode: debugMode,
		stack:     []chunk.Value{},
	}
}

func (vm *VM) Run() Result {
	for {
		// if debug mode is turned on, trace program execution
		if vm.debugMode {
			vm.stackInfo()
			chunk.DisAsmInstruction(vm.ck, vm.ip)
		}
		instruction := vm.readByte()
		switch instruction {
		case chunk.OP_RETURN:
			fmt.Println(vm.stackPop().String())
			return OK
		case chunk.OP_CONSTANT:
			constant := vm.readConstant()
			vm.stackPush(constant)
		case chunk.OP_NEGATE:
			vm.stackPush(-vm.stackPop())
		case chunk.OP_ADD:
			b := vm.stackPop()
			a := vm.stackPop()
			vm.stackPush(a + b)
		case chunk.OP_SUBTRACT:
			b := vm.stackPop()
			a := vm.stackPop()
			vm.stackPush(a - b)
		case chunk.OP_MULTIPLY:
			b := vm.stackPop()
			a := vm.stackPop()
			vm.stackPush(a * b)
		case chunk.OP_DIVIDE:
			b := vm.stackPop()
			a := vm.stackPop()
			vm.stackPush(a / b)
		}
	}
}

func (vm *VM) readByte() byte {
	bt := vm.ck.Codes[vm.ip]
	vm.ip++
	return bt
}

func (vm *VM) readConstant() chunk.Value {
	return vm.ck.Constants[vm.readByte()]
}
