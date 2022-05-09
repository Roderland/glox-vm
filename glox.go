package main

import (
	"github.com/Roderland/glox-vm/chunk"
	"github.com/Roderland/glox-vm/vm"
)

func main() {
	ck := chunk.NewChunk()

	// 1.2
	ck.WriteConstant(1.2, 1)

	// 3.4
	ck.WriteConstant(3.4, 1)

	// 1.2 + 3.4
	ck.Write(chunk.OP_ADD, 1)

	// 5.6
	ck.WriteConstant(5.6, 1)

	// 4.6 / 5.6
	ck.Write(chunk.OP_DIVIDE, 1)

	// -0.8214285714285714
	ck.Write(chunk.OP_NEGATE, 1)

	// return
	ck.Write(chunk.OP_RETURN, 1)

	vm.NewVM(ck, true).Run()
}
