package main

import "fmt"

func add(chunk *Chunk) {
	index1 := chunk.addConstant(1.2)
	chunk.addBytecode(OP_CONSTANT, 123)
	chunk.addBytecode(uint8(index1), 123)

	index2 := chunk.addConstant(3.4)
	chunk.addBytecode(OP_CONSTANT, 123)
	chunk.addBytecode(uint8(index2), 123)

	chunk.addBytecode(OP_ADD, 123)

	index3 := chunk.addConstant(5.6)
	chunk.addBytecode(OP_CONSTANT, 123)
	chunk.addBytecode(uint8(index3), 123)

	chunk.addBytecode(OP_DIVIDE, 123)

	chunk.addBytecode(OP_NEGATE, 123)

	chunk.addBytecode(OP_RETURN, 123)
}

func main() {
	chunk := _Chunk()
	add(chunk)
	vm := _VM()
	vm.interpret(chunk)
	vm.run()
}

func disassembleChunk(chunk *Chunk, title string) {
	fmt.Printf("== %s ==\n", title)
	for offset := 0; offset < len(chunk.bytecodes); offset = disassembleInstruction(chunk, offset) {}
}

func disassembleInstruction(chunk *Chunk, offset int) int {
	fmt.Printf("%04d ", offset)
	if offset > 0 && chunk.lines[offset] == chunk.lines[offset-1] {
		fmt.Printf("   | ");
	} else {
		fmt.Printf("%4d ", chunk.lines[offset])
	}
	instruction := chunk.bytecodes[offset]
	switch chunk.bytecodes[offset] {
	case OP_RETURN:
		return simpleInstruction("OP_RETURN", offset)
	case OP_CONSTANT:
		return constantInstruction("OP_CONSTANT", chunk, offset)
	case OP_NEGATE:
		return simpleInstruction("OP_NEGATE", offset)
	case OP_ADD:
		return simpleInstruction("OP_ADD", offset)
	case OP_SUBTRACT:
		return simpleInstruction("OP_SUBTRACT", offset)
	case OP_MULTIPLY:
		return simpleInstruction("OP_MULTIPLY", offset)
	case OP_DIVIDE:
		return simpleInstruction("OP_DIVIDE", offset)
	default:
		fmt.Printf("Unknown opcode %d\n", instruction)
		return offset + 1
	}
}

func simpleInstruction(name string, offset int) int {
	fmt.Printf("%s\n", name)
	return offset + 1
}

func constantInstruction(name string, chunk *Chunk, offset int) int {
	index := chunk.bytecodes[offset+1]
	fmt.Printf("%-16s %4d '", name, index)
	printValue(chunk.constants[index])
	fmt.Println()
	return offset + 2
}

func printValue(value Value) {
	fmt.Printf("%g", value)
}
