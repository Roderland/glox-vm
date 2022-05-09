package main

import (
	"github.com/Roderland/glox-vm/chunk"
)

func main() {
	ck := chunk.NewChunk()

	ck.WriteConstant(1.2, 1)
	ck.Write(chunk.OP_RETURN, 1)

	chunk.NewDebugger(ck, "test").DisAsmChunk()
}
