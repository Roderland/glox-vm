package vm

import (
	"fmt"
	"github.com/Roderland/glox-vm/chunk"
	"github.com/Roderland/glox-vm/utils"
)

const MAX_FRAME = 64

type VM struct {
	frames     [MAX_FRAME]CallFrame
	frameCount int
	stack      []chunk.Value
	globals    map[string]chunk.Value
}

type CallFrame struct {
	function  *chunk.ObjFunction
	ip        int
	stackSlot int
}

func Do(function *chunk.ObjFunction, debugMode bool) bool {
	vm := initVM()
	vm.stackPush(chunk.NewObject(chunk.NewFunction(*function)))
	vm.call(*function, 0)
	vm.frames[vm.frameCount].function = function
	vm.frames[vm.frameCount].ip = 0
	vm.frames[vm.frameCount].stackSlot = 0
	return vm.Run(debugMode)
}

func initVM() *VM {
	return &VM{
		frames:     [MAX_FRAME]CallFrame{},
		frameCount: 0,
		stack:      []chunk.Value{},
		globals:    map[string]chunk.Value{},
	}
}

func (vm *VM) Run(debugMode bool) bool {
	for {
		// if debug mode is turned on, trace program execution
		if debugMode {
			vm.stackInfo()
			chunk.DisAsmInstruction(&(vm.frames[vm.frameCount-1].function.Ck), vm.frames[vm.frameCount-1].ip)
		}
		instruction := vm.readByte()
		switch instruction {
		case chunk.OP_RETURN:
			result := vm.stackPop()
			vm.frameCount--
			if vm.frameCount == 0 {
				vm.stackPop()
				return true
			}

			stackLength := vm.frames[vm.frameCount].stackSlot
			vm.stack = vm.stack[:stackLength]
			vm.stackPush(result)
			// frame = &vm.frames[vm.frameCount-1]

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
			slot := vm.frames[vm.frameCount-1].stackSlot + int(vm.readByte())
			vm.stackPush(vm.stack[slot])

		case chunk.OP_SET_LOCAL:
			slot := vm.frames[vm.frameCount-1].stackSlot + int(vm.readByte())
			vm.stack[slot] = vm.stackPeek(0)

		case chunk.OP_JUMP_IF_FALSE:
			offset := vm.readShort()
			if vm.stackPeek(0).IsFalse() {
				vm.frames[vm.frameCount-1].ip += int(offset)
			}

		case chunk.OP_JUMP:
			offset := vm.readShort()
			vm.frames[vm.frameCount-1].ip += int(offset)

		case chunk.OP_LOOP:
			offset := vm.readShort()
			vm.frames[vm.frameCount-1].ip -= int(offset)

		case chunk.OP_CALL:
			argCount := int(vm.readByte())
			if !vm.callValue(vm.stackPeek(argCount), argCount) {
				return false
			}
			// frame = &vm.frames[vm.frameCount - 1];
		}
	}
}

func (vm *VM) callValue(callee chunk.Value, argCount int) bool {
	if callee.IsObject() {
		obj := callee.AsObject()
		if obj.IsFunction() {
			return vm.call(obj.AsFunction(), argCount)
		}
	}
	vm.runtimeError("Can only call functions and classes.")
	return false
}

func (vm *VM) call(f chunk.ObjFunction, argCount int) bool {
	if argCount != f.Arity {
		vm.runtimeError("Expected %d arguments but got %d.", f.Arity, argCount)
		return false
	}

	if vm.frameCount == MAX_FRAME {
		vm.runtimeError("Stack overflow.")
		return false
	}

	frame := &vm.frames[vm.frameCount]
	vm.frameCount++
	frame.function = &f
	frame.ip = 0
	frame.stackSlot = vm.stackSize() - argCount - 1
	return true
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
	bt := vm.frames[vm.frameCount-1].function.Ck.Codes[vm.frames[vm.frameCount-1].ip]
	vm.frames[vm.frameCount-1].ip++
	return bt
}

func (vm *VM) readShort() uint16 {
	bt1 := uint16(vm.frames[vm.frameCount-1].function.Ck.Codes[vm.frames[vm.frameCount-1].ip])
	vm.frames[vm.frameCount-1].ip++
	bt2 := uint16(vm.frames[vm.frameCount-1].function.Ck.Codes[vm.frames[vm.frameCount-1].ip])
	vm.frames[vm.frameCount-1].ip++
	return (bt1 << 8) | bt2
}

func (vm *VM) readConstant() chunk.Value {
	return vm.frames[vm.frameCount-1].function.Ck.Constants[vm.readByte()]
}

func (vm *VM) runtimeError(format string, a ...interface{}) {
	utils.PrintfErr(format+"\n", a...)

	for i := vm.frameCount - 1; i >= 0; i-- {
		frame := &vm.frames[i]
		fun := frame.function
		utils.PrintfErr("[line %d] in ", fun.Ck.Lines[frame.ip-1])

		if fun.Name == "" {
			utils.PrintfErr("script\n")
		} else {
			utils.PrintfErr("%s()\n", fun.Name)
		}
	}

	vm.stackReset()
}
