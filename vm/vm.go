package vm

import (
	"fmt"
	. "glox-vm"
	"os"
	"reflect"
	"unsafe"
)

const CALL_FRAMES_MAX = 64

var emptyValue = Value{}

type (
	CallFrame struct {
		function *FuncData
		ip       *byte
		slots    *Value
	}
	VM struct {
		//ip      *byte
		//chunk   *Chunk
		stack          *Stack
		globals        map[string]Value
		callFrames     [CALL_FRAMES_MAX]CallFrame
		callFrameCount int
	}
)

func InitVM(function *FuncData) *VM {
	stack := new(Stack)
	stack.reset()
	var vm VM
	vm.stack = stack
	vm.globals = map[string]Value{}
	//vm.callFrameCount = 1
	//vm.callFrames[0] = CallFrame{
	//	function: function,
	//	ip:       (*byte)(unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&function.FunChunk.Bytecodes)).Data)),
	//	slots:    stack.top,
	//}
	vm.stack.push(NewFunction(*function))
	vm.callValue(vm.stack.peek(0), 0)
	return &vm
}

// scriptPoint 脚本起点
func (vm *VM) scriptPoint() *byte {
	bytes := vm.callFrames[0].function.FunChunk.Bytecodes
	return (*byte)(unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&bytes)).Data))
}

// functionPoint 当前函数起点
func (vm *VM) functionPoint() *byte {
	bytes := vm.callFrame().function.FunChunk.Bytecodes
	return (*byte)(unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&bytes)).Data))
}

// currentPoint 当前地址
func (vm *VM) currentPoint() *byte {
	return vm.callFrame().ip
}

// returnPoint 返回调用者地址
func (vm *VM) returnPoint() *byte {
	if vm.callFrameCount < 1 {
		panic("<script> can't return.")
	}
	return vm.callFrames[vm.callFrameCount-2].ip
}

// callFrame 当前函数调用栈帧
func (vm *VM) callFrame() *CallFrame {
	return &vm.callFrames[vm.callFrameCount-1]
}

// getFrameSlot 从当前函数栈帧中获取value
func (vm *VM) getFrameSlot(index uint8) Value {
	return *(*Value)(unsafe.Pointer(uintptr(unsafe.Pointer(vm.callFrame().slots)) + uintptr(index)*unsafe.Sizeof(emptyValue)))
}

// setFrameSlot 向当前函数栈帧中设置value
func (vm *VM) setFrameSlot(index uint8, value Value) {
	*(*Value)(unsafe.Pointer(uintptr(unsafe.Pointer(vm.callFrame().slots)) + uintptr(index)*unsafe.Sizeof(emptyValue))) = value
}

func (vm *VM) Run() InterpretResult {
	for {
		vm.stack.print()
		offset := uintptr(unsafe.Pointer(vm.currentPoint())) - uintptr(unsafe.Pointer(vm.functionPoint()))
		frame := vm.callFrame()
		DisassembleInstruction(&frame.function.FunChunk, int(offset))

		switch vm.readBytecode() {
		case OP_RETURN:
			result := vm.stack.pop()
			vm.callFrameCount --
			if vm.callFrameCount == 0 {
				vm.stack.pop()
				return OK
			}
			vm.stack.top = vm.callFrames[vm.callFrameCount].slots
			vm.stack.pop()
			vm.stack.push(result)
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
				vm.runtimeError("Operands must be numbers or strings.")
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
			value := vm.getFrameSlot(index)
			vm.stack.push(value)
		case OP_SET_LOCAL:
			index := vm.readBytecode()
			value := vm.stack.peek(0)
			vm.setFrameSlot(index, value)
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
		case OP_CALL:
			argCount := vm.readBytecode()
			funcValue := vm.stack.peek(int(argCount))
			if !vm.callValue(funcValue, argCount) {
				return RUNTIME_ERROR
			}
		}
	}
}

// readConstant 以下一个字节码为索引从常量池中读取一个常量
func (vm *VM) readConstant() Value {
	index := vm.readBytecode()
	return vm.callFrame().function.FunChunk.Constants[index]
}

// readBytecode 读取下一个字节码
func (vm *VM) readBytecode() (bt byte) {
	bt = *vm.currentPoint()
	vm.next()
	return
}

// readShort 读取两个字节码
func (vm *VM) readShort() uint16 {
	high := *vm.currentPoint()
	vm.next()
	low := *vm.currentPoint()
	vm.next()
	return (uint16(high) << 8) | uint16(low)
}

func (vm *VM) next() {
	vm.callFrame().ip = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(vm.currentPoint())) + 1))
}

func (vm *VM) jump(offset uint16) {
	vm.callFrame().ip = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(vm.currentPoint())) + uintptr(offset)))
}

func (vm *VM) loop(offset uint16) {
	vm.callFrame().ip = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(vm.currentPoint())) - uintptr(offset)))
}

func (vm *VM) callValue(callee Value, argCount uint8) bool {
	if callee.IsFunction() {
		funcData := callee.AsFunction()
		return vm.callFunction(&funcData, argCount)
	}
	vm.runtimeError("Can only call functions and classes.")
	return false
}

func (vm *VM) callFunction(fd *FuncData, argCount uint8) bool {
	if argCount != uint8(fd.Arity) {
		vm.runtimeError("Expected %d arguments but got %d.", fd.Arity, argCount)
		return false
	}
	if vm.callFrameCount == CALL_FRAMES_MAX {
		vm.runtimeError("Stack overflow.")
		return false
	}
	//callerAddr := frame.ip
	vm.callFrameCount++
	frame := vm.callFrame()
	frame.function = fd
	bytes := fd.FunChunk.Bytecodes
	frame.ip = (*byte)(unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&bytes)).Data))
	frame.slots = (*Value)(unsafe.Pointer(uintptr(unsafe.Pointer(vm.stack.top)) - uintptr(argCount)*unsafe.Sizeof(emptyValue)))
	return true
}

func (vm *VM) runtimeError(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format+"\n", a...)
	for i := vm.callFrameCount-1; i>=0 ; i-- {
		frame := vm.callFrames[i]
		function := frame.function
		offset := uintptr(unsafe.Pointer(frame.ip)) - (*reflect.SliceHeader)(unsafe.Pointer(&function.FunChunk.Bytecodes)).Data - 1
		line := function.FunChunk.Lines[offset]
		_, _ = fmt.Fprintf(os.Stderr, "[line %d] in ", line)
		if function.Name == "" {
			_, _ = fmt.Fprintf(os.Stderr, "script\n")
		} else {
			_, _ = fmt.Fprintf(os.Stderr, "%s()\n", function.Name)
		}
	}
	vm.stack.reset()
}
