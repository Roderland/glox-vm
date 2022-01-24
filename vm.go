package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

type VM struct {
	ip    *byte
	chunk *Chunk
	stack *Stack
}

func _VM() *VM {
	stack := new(Stack)
	stack.reset()
	return &VM{stack: stack}
}

func (vm *VM) next() {
	vm.ip = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(vm.ip)) + 1))
}

func (vm *VM) interpret(source string) InterpretResult {
	vm.chunk = nil
	vm.ip = (*byte)(unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&vm.chunk.bytecodes)).Data))
	return OK
}

func (vm *VM) run() InterpretResult {
	for {
		vm.stack.print(0)
		offset := uintptr(unsafe.Pointer(vm.ip)) - (*reflect.SliceHeader)(unsafe.Pointer(&vm.chunk.bytecodes)).Data
		disassembleInstruction(vm.chunk, int(offset))
		instruction := *vm.ip
		vm.next()
		switch instruction {
		case OP_RETURN:
			printValue(vm.stack.pop())
			fmt.Println()
			return OK
		case OP_CONSTANT:
			index := *vm.ip
			vm.next()
			constant := vm.chunk.constants[int(index)]
			vm.stack.push(constant)
		case OP_NEGATE:
			vm.stack.push(-vm.stack.pop())
		case OP_ADD:
			b := vm.stack.pop()
			a := vm.stack.pop()
			vm.stack.push(a + b)
		case OP_SUBTRACT:
			b := vm.stack.pop()
			a := vm.stack.pop()
			vm.stack.push(a - b)
		case OP_MULTIPLY:
			b := vm.stack.pop()
			a := vm.stack.pop()
			vm.stack.push(a * b)
		case OP_DIVIDE:
			b := vm.stack.pop()
			a := vm.stack.pop()
			vm.stack.push(a / b)
		}
	}
}
