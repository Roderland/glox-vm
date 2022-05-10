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

func (vm *VM) stackPeek(distance int) chunk.Value {
	return vm.stack[vm.stackSize()-distance-1]
}

func (vm *VM) stackSize() int {
	return len(vm.stack)
}

func (vm *VM) stackReset() {
	vm.stack = []chunk.Value{}
}

func (vm *VM) stackInfo() {
	utils.PrintfDbg("          ")
	for idx := 0; idx < vm.stackSize(); idx++ {
		utils.PrintfDbg("[ ")
		val := vm.stack[idx]
		str := val.AsString()
		if val.IsString() {
			str = "\"" + str + "\""
		}
		utils.PrintfDbg(str)
		utils.PrintfDbg(" ]")
	}
	utils.PrintfDbg("\n")
}
