package vm

import (
	"github.com/Roderland/glox-vm/chunk"
	"github.com/Roderland/glox-vm/utils"
	"math"
)

const STACK_MAX = math.MaxUint8

func (vm *VM) stackPush(value chunk.Value) {
	if vm.stackSize() >= STACK_MAX {
		// todo
	}
	vm.stack = append(vm.stack, value)
}

func (vm *VM) stackPop() chunk.Value {
	idx := len(vm.stack) - 1
	value := vm.stack[idx]
	vm.stack = vm.stack[:idx]
	return value
}

func (vm *VM) stackSize() int {
	return len(vm.stack)
}

func (vm *VM) stackInfo() {
	utils.PrintfDbg("          ")
	for idx := 0; idx < vm.stackSize(); idx++ {
		utils.PrintfDbg("[ ")
		utils.PrintfDbg(vm.stack[idx].String())
		utils.PrintfDbg(" ]")
	}
	utils.PrintfDbg("\n")
}
