package glox_vm

import "fmt"

func DisassembleChunk(chunk *Chunk, title string) {
	fmt.Printf("== %s ==\n", title)
	for offset := 0; offset < len(chunk.Bytecodes); offset = DisassembleInstruction(chunk, offset) {
	}
}

func DisassembleInstruction(chunk *Chunk, offset int) int {
	fmt.Printf("%04d ", offset)
	if offset > 0 && chunk.Lines[offset] == chunk.Lines[offset-1] {
		fmt.Printf("   | ")
	} else {
		fmt.Printf("%4d ", chunk.Lines[offset])
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
	case OP_NIL:
		return simpleInstruction("OP_NIL", offset)
	case OP_FALSE:
		return simpleInstruction("OP_FALSE", offset)
	case OP_TRUE:
		return simpleInstruction("OP_TRUE", offset)
	case OP_NOT:
		return simpleInstruction("OP_NOT", offset)
	case OP_EQUAL:
		return simpleInstruction("OP_EQUAL", offset)
	case OP_GREATER:
		return simpleInstruction("OP_GREATER", offset)
	case OP_LESS:
		return simpleInstruction("OP_LESS", offset)
	case OP_PRINT:
		return simpleInstruction("OP_PRINT", offset)
	case OP_POP:
		return simpleInstruction("OP_POP", offset)
	case OP_DEFINE_GLOBAL:
		return constantInstruction("OP_DEFINE_GLOBAL", chunk, offset)
	case OP_GET_GLOBAL:
		return constantInstruction("OP_GET_GLOBAL", chunk, offset)
	case OP_SET_GLOBAL:
		return constantInstruction("OP_SET_GLOBAL", chunk, offset)
	case OP_GET_LOCAL:
		return byteInstruction("OP_GET_LOCAL", chunk, offset)
	case OP_SET_LOCAL:
		return byteInstruction("OP_SET_LOCAL", chunk, offset)
	case OP_JUMP_IF_FALSE:
		return jumpInstruction("OP_JUMP_IF_FALSE", 1, chunk, offset)
	case OP_JUMP:
		return jumpInstruction("OP_JUMP", 1, chunk, offset)
	case OP_LOOP:
		return jumpInstruction("OP_LOOP", -1, chunk, offset)
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

func byteInstruction(name string, chunk *Chunk, offset int) int {
	index := chunk.Bytecodes[offset+1]
	fmt.Printf("%-16s %4d '", name, index)
	fmt.Println()
	return offset + 2
}

func jumpInstruction(name string, sign int, chunk *Chunk, offset int) int {
	jump := uint16(chunk.Bytecodes[offset+1]) << 8
	jump |= uint16(chunk.Bytecodes[offset+2])
	fmt.Printf("%-16s %4d -> %d\n", name, offset, offset+3+sign*int(jump))
	return offset + 3
}

func PrintValue(value Value) {
	switch value.typ {
	case VAL_BOOL:
		fmt.Print(value.AsBool())
	case VAL_NIL:
		fmt.Print("nil")
	case VAL_NUMBER:
		fmt.Printf("%g", value.AsNumber())
	case VAL_STRING:
		fmt.Printf("%s", value.AsString())
	}
}
