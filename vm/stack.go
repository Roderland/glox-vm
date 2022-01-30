package vm

import (
	"fmt"
	. "glox-vm"
	"math"
	"unsafe"
)

const STACK_MAX = CALL_FRAMES_MAX * math.MaxUint8

type Stack struct {
	top    *Value
	frames [STACK_MAX]Value
}

func (stack *Stack) reset() {
	stack.top = (*Value)(unsafe.Pointer(&stack.frames))
}

func (stack *Stack) push(value Value) {
	*stack.top = value
	stack.top = (*Value)(unsafe.Pointer(uintptr(unsafe.Pointer(stack.top)) + unsafe.Sizeof(emptyValue)))
}

func (stack *Stack) pop() Value {
	stack.top = (*Value)(unsafe.Pointer(uintptr(unsafe.Pointer(stack.top)) - unsafe.Sizeof(emptyValue)))
	return *stack.top
}

func (stack *Stack) peek(distance int) Value {
	return *(*Value)(unsafe.Pointer(uintptr(unsafe.Pointer(stack.top)) - uintptr(distance+1)*unsafe.Sizeof(emptyValue)))
}

func (stack *Stack) print() {
	fmt.Print("          ")
	n := (uintptr(unsafe.Pointer(stack.top)) - uintptr(unsafe.Pointer(&stack.frames))) / unsafe.Sizeof(emptyValue)
	for i := 0; i < int(n); i ++  {
		fmt.Print("[ ")
		PrintValue(*(*Value)(unsafe.Pointer(uintptr(unsafe.Pointer(&stack.frames)) + uintptr(i)*unsafe.Sizeof(emptyValue))))
		fmt.Print(" ]")
	}
	fmt.Println()
}
