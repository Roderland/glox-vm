package chunk

import (
	"github.com/Roderland/glox-vm/utils"
)

func DisAsmChunk(ck *Chunk, name string) {
	utils.PrintfDbg("=== %s ===\n", name)

	for offset := 0; offset < len(ck.Codes); {
		offset = DisAsmInstruction(ck, offset)
	}
}

func DisAsmInstruction(ck *Chunk, offset int) int {
	utils.PrintfDbg("%04d ", offset)
	if offset > 0 && ck.Lines[offset] == ck.Lines[offset-1] {
		utils.PrintfDbg("   | ")
	} else {
		utils.PrintfDbg("%4d ", ck.Lines[offset])
	}

	instruction := ck.Codes[offset]
	switch instruction {
	case OP_RETURN:
		return simpleInstruction("OP_RETURN", offset)
	case OP_CONSTANT:
		return constantInstruction("OP_CONSTANT", ck, offset)
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
	default:
		utils.PrintfDbg("Unknown opcode %d\n", instruction)
		return offset + 1
	}
}

func simpleInstruction(name string, offset int) int {
	utils.PrintfDbg("%s\n", name)
	return offset + 1
}

func constantInstruction(name string, ck *Chunk, offset int) int {
	idx := ck.Codes[offset+1]
	utils.PrintfDbg("%-16s %4d '", name, idx)
	utils.PrintfDbg(ck.Constants[idx].String())
	utils.PrintfDbg("\n")
	return offset + 2
}
