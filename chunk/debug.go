package chunk

import (
	"fmt"
)

type Debugger struct {
	name   string
	ck     *Chunk
	offset int
}

func NewDebugger(ck *Chunk, name string) *Debugger {
	return &Debugger{
		name:   name,
		ck:     ck,
		offset: 0,
	}
}

func (dbg *Debugger) DisAsmChunk() {
	fmt.Printf("=== %s ===\n", dbg.name)

	for dbg.offset = 0; dbg.offset < len(dbg.ck.codes); {
		dbg.disAsmInstruction()
	}
}

func (dbg *Debugger) disAsmInstruction() {
	fmt.Printf("%04d ", dbg.offset)
	if dbg.offset > 0 && dbg.ck.lines[dbg.offset] == dbg.ck.lines[dbg.offset-1] {
		fmt.Printf("   | ")
	} else {
		fmt.Printf("%4d ", dbg.ck.lines[dbg.offset])
	}

	instruction := dbg.ck.codes[dbg.offset]
	switch instruction {
	case OP_RETURN:
		dbg.simpleInstruction("OP_RETURN")
	case OP_CONSTANT:
		dbg.constantInstruction("OP_CONSTANT")
	default:
		fmt.Printf("Unknown opcode %d\n", instruction)
		dbg.offset++
	}
}

func (dbg *Debugger) simpleInstruction(name string) {
	fmt.Printf("%s\n", name)
	dbg.offset++
}

func (dbg *Debugger) constantInstruction(name string) {
	idx := dbg.ck.codes[dbg.offset+1]
	fmt.Printf("%-16s %4d '", name, idx)
	dbg.ck.constants[idx].print()
	fmt.Println()
	dbg.offset += 2
}
