package vm

import (
	"fmt"
	"github.com/Roderland/glox-vm/chunk"
	"github.com/Roderland/glox-vm/utils"
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

func (vm *VM) Run() bool {
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
			return true

		case chunk.OP_CONSTANT:
			constant := vm.readConstant()
			vm.stackPush(constant)

		case chunk.OP_NEGATE:
			if !vm.stackPeek(0).IsNumber() {
				vm.runtimeError("Operand must be a number.")
				return false
			}
			vm.stackPush(chunk.NewNumber(-vm.stackPop().AsNumber()))

		case chunk.OP_ADD:
			a, b, ok := vm.popBinaryNumber()
			if !ok {
				return false
			}
			vm.stackPush(chunk.NewNumber(a + b))

		case chunk.OP_SUBTRACT:
			a, b, ok := vm.popBinaryNumber()
			if !ok {
				return false
			}
			vm.stackPush(chunk.NewNumber(a - b))

		case chunk.OP_MULTIPLY:
			a, b, ok := vm.popBinaryNumber()
			if !ok {
				return false
			}
			vm.stackPush(chunk.NewNumber(a * b))

		case chunk.OP_DIVIDE:
			a, b, ok := vm.popBinaryNumber()
			if !ok {
				return false
			}
			vm.stackPush(chunk.NewNumber(a / b))

		case chunk.OP_NIL:
			vm.stackPush(chunk.Nil)
		case chunk.OP_FALSE:
			vm.stackPush(chunk.False)
		case chunk.OP_TRUE:
			vm.stackPush(chunk.True)

		case chunk.OP_NOT:
			vm.stackPush(chunk.NewBool(vm.stackPop().IsFalse()))

		case chunk.OP_EQUAL:
			b := vm.stackPop()
			a := vm.stackPop()
			vm.stackPush(chunk.NewBool(chunk.Equal(a, b)))

		case chunk.OP_GREATER:
			a, b, ok := vm.popBinaryNumber()
			if !ok {
				return false
			}
			vm.stackPush(chunk.NewBool(a > b))

		case chunk.OP_LESS:
			a, b, ok := vm.popBinaryNumber()
			if !ok {
				return false
			}
			vm.stackPush(chunk.NewBool(a < b))
		}
	}
}

func (vm *VM) popBinaryNumber() (float64, float64, bool) {
	if !vm.stackPeek(0).IsNumber() || !vm.stackPeek(1).IsNumber() {
		vm.runtimeError("Operands must be numbers.")
		return 0, 0, false
	}
	b := vm.stackPop().AsNumber()
	a := vm.stackPop().AsNumber()
	return a, b, true
}

func (vm *VM) readByte() byte {
	bt := vm.ck.Codes[vm.ip]
	vm.ip++
	return bt
}

func (vm *VM) readConstant() chunk.Value {
	return vm.ck.Constants[vm.readByte()]
}

func (vm *VM) runtimeError(format string, a ...interface{}) {
	utils.PrintfErr(format, a...)
	line := vm.ck.Lines[vm.ip-1]
	utils.PrintfErr("[line %d] in script\n", line)
	vm.stackReset()
}
