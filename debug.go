package glox_vm

import "fmt"

func DisassembleChunk(chunk *Chunk, title string) {
	fmt.Printf("== %s ==\n", title)
	for offset := 0; offset < len(chunk.Bytecodes); offset = DisassembleInstruction(chunk, offset) {}
}

func DisassembleInstruction(chunk *Chunk, offset int) int {
	fmt.Printf("%04d ", offset)
	if offset > 0 && chunk.lines[offset] == chunk.lines[offset-1] {
		fmt.Printf("   | ");
	} else {
		fmt.Printf("%4d ", chunk.lines[offset])
	}
	instruction := chunk.Bytecodes[offset]
	switch chunk.Bytecodes[offset] {
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
	index := chunk.Bytecodes[offset+1]
	fmt.Printf("%-16s %4d '", name, index)
	PrintValue(chunk.Constants[index])
	fmt.Println()
	return offset + 2
}

func PrintValue(value Value) {
	fmt.Printf("%g", value)
}