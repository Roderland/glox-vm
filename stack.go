package main

import (
	"fmt"
	"unsafe"
)

const STACK_MAX = 64

type Stack struct {
	top    *Value
	frames [STACK_MAX]Value
}

func (stack *Stack) reset() {
	stack.top = (*Value)(unsafe.Pointer(&stack.frames))
}

func (stack *Stack) push(value Value) {
	*stack.top = value
	stack.top = (*Value)(unsafe.Pointer(uintptr(unsafe.Pointer(stack.top)) + unsafe.Sizeof(value)))
}

func (stack *Stack) pop() (value Value) {
	stack.top = (*Value)(unsafe.Pointer(uintptr(unsafe.Pointer(stack.top)) - unsafe.Sizeof(value)))
	return *stack.top
}

func (stack *Stack) print(value Value) {
	fmt.Print("          ")
	n := (uintptr(unsafe.Pointer(stack.top)) - uintptr(unsafe.Pointer(&stack.frames))) / unsafe.Sizeof(value)
	for i := 0; i < int(n); i ++  {
		fmt.Print("[ ")
		printValue(*(*Value)(unsafe.Pointer(uintptr(unsafe.Pointer(&stack.frames)) + uintptr(i)*unsafe.Sizeof(value))))
		fmt.Print(" ]")
	}
	fmt.Println()
}
