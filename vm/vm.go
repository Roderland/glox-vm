package vm

import (
	"fmt"
	. "glox-vm"
	"os"
	"reflect"
	"unsafe"
)

type VM struct {
	ip    *byte
	chunk *Chunk
	stack *Stack
}

func InitVM(chunk *Chunk) *VM {
	stack := new(Stack)
	stack.reset()
	var vm VM
	vm.stack = stack
	vm.chunk = chunk
	vm.ip = (*byte)(unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&vm.chunk.Bytecodes)).Data))
	return &vm
}

func (vm *VM) next() {
	vm.ip = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(vm.ip)) + 1))
}

func (vm *VM) Run() InterpretResult {
	for {
		vm.stack.print(NewNil())
		offset := uintptr(unsafe.Pointer(vm.ip)) - (*reflect.SliceHeader)(unsafe.Pointer(&vm.chunk.Bytecodes)).Data
		DisassembleInstruction(vm.chunk, int(offset))
		instruction := *vm.ip
		vm.next()
		switch instruction {
		case OP_RETURN:
			PrintValue(vm.stack.pop())
			fmt.Println()
			return OK
		case OP_CONSTANT:
			index := *vm.ip
			vm.next()
			constant := vm.chunk.Constants[int(index)]
			vm.stack.push(constant)
		case OP_NEGATE:
			if !vm.stack.peek(0).IsNumber() {
				vm.runtimeError("Operand must be a number.")
				return RUNTIME_ERROR
			}
			vm.stack.push(NewNumber(-vm.stack.pop().AsNumber()))
		case OP_ADD:
			if !vm.stack.peek(0).IsNumber() || !vm.stack.peek(1).IsNumber() {
				vm.runtimeError("Operands must be numbers.")
				return RUNTIME_ERROR
			}
			b := vm.stack.pop().AsNumber()
			a := vm.stack.pop().AsNumber()
			vm.stack.push(NewNumber(a + b))
		case OP_SUBTRACT:
			if !vm.stack.peek(0).IsNumber() || !vm.stack.peek(1).IsNumber() {
				vm.runtimeError("Operands must be numbers.")
				return RUNTIME_ERROR
			}
			b := vm.stack.pop().AsNumber()
			a := vm.stack.pop().AsNumber()
			vm.stack.push(NewNumber(a - b))
		case OP_MULTIPLY:
			if !vm.stack.peek(0).IsNumber() || !vm.stack.peek(1).IsNumber() {
				vm.runtimeError("Operands must be numbers.")
				return RUNTIME_ERROR
			}
			b := vm.stack.pop().AsNumber()
			a := vm.stack.pop().AsNumber()
			vm.stack.push(NewNumber(a * b))
		case OP_DIVIDE:
			if !vm.stack.peek(0).IsNumber() || !vm.stack.peek(1).IsNumber() {
				vm.runtimeError("Operands must be numbers.")
				return RUNTIME_ERROR
			}
			b := vm.stack.pop().AsNumber()
			a := vm.stack.pop().AsNumber()
			vm.stack.push(NewNumber(a / b))
		case OP_NIL:
			vm.stack.push(NewNil())
		case OP_FALSE:
			vm.stack.push(NewBool(false))
		case OP_TRUE:
			vm.stack.push(NewBool(true))
		case OP_NOT:
			vm.stack.push(NewBool(vm.stack.pop().IsFalse()))
		case OP_EQUAL:
			b := vm.stack.pop()
			a := vm.stack.pop()
			vm.stack.push(NewBool(a.Equal(b)))
		case OP_GREATER:
			if !vm.stack.peek(0).IsNumber() || !vm.stack.peek(1).IsNumber() {
				vm.runtimeError("Operands must be numbers.")
				return RUNTIME_ERROR
			}
			b := vm.stack.pop().AsNumber()
			a := vm.stack.pop().AsNumber()
			vm.stack.push(NewBool(a > b))
		case OP_LESS:
			if !vm.stack.peek(0).IsNumber() || !vm.stack.peek(1).IsNumber() {
				vm.runtimeError("Operands must be numbers.")
				return RUNTIME_ERROR
			}
			b := vm.stack.pop().AsNumber()
			a := vm.stack.pop().AsNumber()
			vm.stack.push(NewBool(a < b))
		}
	}
}

func (vm *VM) runtimeError(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format + "\n", a...)
	offset := uintptr(unsafe.Pointer(vm.ip)) - (*reflect.SliceHeader)(unsafe.Pointer(&vm.chunk.Bytecodes)).Data - 1
	line := vm.chunk.Lines[offset]
	_, _ = fmt.Fprintf(os.Stderr, "[line %d] in script\n", line)
	vm.stack.reset()
}
