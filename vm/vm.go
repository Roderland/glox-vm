package vm

import (
	"fmt"
	. "glox-vm"
	"os"
	"reflect"
	"unsafe"
)

type VM struct {
	ip      *byte
	chunk   *Chunk
	stack   *Stack
	globals map[string]Value
}

func InitVM(chunk *Chunk) *VM {
	stack := new(Stack)
	stack.reset()
	var vm VM
	vm.stack = stack
	vm.chunk = chunk
	vm.ip = (*byte)(unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&vm.chunk.Bytecodes)).Data))
	vm.globals = map[string]Value{}
	return &vm
}

func (vm *VM) Run() InterpretResult {
	for {
		vm.stack.print(NewNil())
		offset := uintptr(unsafe.Pointer(vm.ip)) - (*reflect.SliceHeader)(unsafe.Pointer(&vm.chunk.Bytecodes)).Data
		DisassembleInstruction(vm.chunk, int(offset))

		switch vm.readBytecode() {
		case OP_RETURN:
			return OK
		case OP_CONSTANT:
			vm.stack.push(vm.readConstant())
		case OP_NEGATE:
			if !vm.stack.peek(0).IsNumber() {
				vm.runtimeError("Operand must be a number.")
				return RUNTIME_ERROR
			}
			vm.stack.push(NewNumber(-vm.stack.pop().AsNumber()))
		case OP_ADD:
			if vm.stack.peek(0).IsString() && vm.stack.peek(1).IsString() {
				b := vm.stack.pop().AsString()
				a := vm.stack.pop().AsString()
				vm.stack.push(NewString(a + b))
			} else if vm.stack.peek(0).IsNumber() && vm.stack.peek(1).IsNumber() {
				b := vm.stack.pop().AsNumber()
				a := vm.stack.pop().AsNumber()
				vm.stack.push(NewNumber(a + b))
			} else {
				vm.runtimeError("Operands must be numbers.")
				return RUNTIME_ERROR
			}
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
		case OP_PRINT:
			PrintValue(vm.stack.pop())
			fmt.Println()
		case OP_POP:
			vm.stack.pop()
		case OP_DEFINE_GLOBAL:
			name := vm.readConstant().AsString()
			value := vm.stack.pop()
			vm.globals[name] = value
		case OP_GET_GLOBAL:
			name := vm.readConstant().AsString()
			value, ok := vm.globals[name]
			if !ok {
				vm.runtimeError("Undefined variable '%s'.", name)
				return RUNTIME_ERROR
			}
			vm.stack.push(value)
		case OP_SET_GLOBAL:
			name := vm.readConstant().AsString()
			if _, ok := vm.globals[name]; !ok {
				vm.runtimeError("Undefined variable '%s'.", name)
				return RUNTIME_ERROR
			}
			vm.globals[name] = vm.stack.peek(0)
		case OP_GET_LOCAL:
			index := vm.readBytecode()
			value := vm.stack.frames[index]
			vm.stack.push(value)
		case OP_SET_LOCAL:
			index := vm.readBytecode()
			value := vm.stack.peek(0)
			vm.stack.frames[index] = value
		case OP_JUMP_IF_FALSE:
			jump := vm.readShort()
			if vm.stack.peek(0).IsFalse() {
				vm.jump(jump)
			}
		case OP_JUMP:
			jump := vm.readShort()
			vm.jump(jump)
		case OP_LOOP:
			loop := vm.readShort()
			vm.loop(loop)
		}
	}
}

// readConstant 以下一个字节码为索引从常量池中读取一个常量
func (vm *VM) readConstant() Value {
	index := vm.readBytecode()
	return vm.chunk.Constants[index]
}

// readBytecode 读取下一个字节码
func (vm *VM) readBytecode() (bt byte) {
	bt = *vm.ip
	vm.next()
	return
}

// readShort 读取两个字节码
func (vm *VM) readShort() uint16 {
	high := *vm.ip
	vm.next()
	low := *vm.ip
	vm.next()
	return (uint16(high) << 8) | uint16(low)
}

func (vm *VM) next() {
	vm.ip = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(vm.ip)) + 1))
}

func (vm *VM) jump(offset uint16) {
	vm.ip = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(vm.ip)) + uintptr(offset)))
}

func (vm *VM) loop(offset uint16) {
	vm.ip = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(vm.ip)) - uintptr(offset)))
}

func (vm *VM) runtimeError(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format+"\n", a...)
	offset := uintptr(unsafe.Pointer(vm.ip)) - (*reflect.SliceHeader)(unsafe.Pointer(&vm.chunk.Bytecodes)).Data - 1
	line := vm.chunk.Lines[offset]
	_, _ = fmt.Fprintf(os.Stderr, "[line %d] in script\n", line)
	vm.stack.reset()
}
