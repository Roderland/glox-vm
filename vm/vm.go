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
	globals   map[string]chunk.Value
}

func NewVM(ck *chunk.Chunk, debugMode bool) *VM {
	return &VM{
		ck:        ck,
		debugMode: debugMode,
		stack:     []chunk.Value{},
		globals:   map[string]chunk.Value{},
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
			// Exit interpreter.
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
			if vm.stackPeek(0).IsNumber() && vm.stackPeek(1).IsNumber() {
				b := vm.stackPop().AsNumber()
				a := vm.stackPop().AsNumber()
				vm.stackPush(chunk.NewNumber(a + b))
				break
			}
			if vm.stackPeek(0).IsString() && vm.stackPeek(1).IsString() {
				b := vm.stackPop().AsString()
				a := vm.stackPop().AsString()
				vm.stackPush(chunk.NewString(a + b))
				break
			}
			vm.runtimeError("Operands must be numbers or strings.")
			return false

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

		case chunk.OP_PRINT:
			fmt.Println(vm.stackPop().String())

		case chunk.OP_POP:
			vm.stackPop()

		case chunk.OP_DEFINE_GLOBAL:
			name := vm.readConstant().AsString()
			vm.globals[name] = vm.stackPeek(0)
			vm.stackPop()

		case chunk.OP_GET_GLOBAL:
			name := vm.readConstant().AsString()
			val, ok := vm.globals[name]
			if !ok {
				vm.runtimeError("Undefined variable '%s'.", name)
				return false
			}
			vm.stackPush(val)

		case chunk.OP_SET_GLOBAL:
			name := vm.readConstant().AsString()
			_, ok := vm.globals[name]
			if !ok {
				vm.runtimeError("Undefined variable '%s'.", name)
				return false
			}
			vm.globals[name] = vm.stackPeek(0)

		case chunk.OP_GET_LOCAL:
			slot := vm.readByte()
			vm.stackPush(vm.stack[slot])

		case chunk.OP_SET_LOCAL:
			slot := vm.readByte()
			vm.stack[slot] = vm.stackPeek(0)

		case chunk.OP_JUMP_IF_FALSE:
			offset := vm.readShort()
			if vm.stackPeek(0).IsFalse() {
				vm.ip += int(offset)
			}

		case chunk.OP_JUMP:
			offset := vm.readShort()
			vm.ip += int(offset)

		case chunk.OP_LOOP:
			offset := vm.readShort()
			vm.ip -= int(offset)
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

func (vm *VM) popBinaryString() (string, string, bool) {
	if !vm.stackPeek(0).IsString() || !vm.stackPeek(1).IsString() {
		vm.runtimeError("Operands must be strings.")
		return "", "", false
	}
	b := vm.stackPop().AsString()
	a := vm.stackPop().AsString()
	return a, b, true
}

func (vm *VM) readByte() byte {
	bt := vm.ck.Codes[vm.ip]
	vm.ip++
	return bt
}

func (vm *VM) readShort() uint16 {
	bt1 := uint16(vm.ck.Codes[vm.ip])
	vm.ip++
	bt2 := uint16(vm.ck.Codes[vm.ip])
	vm.ip++
	return (bt1 << 8) | bt2
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
