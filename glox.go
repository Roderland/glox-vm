package main

import (
	"fmt"
	"github.com/Roderland/glox-vm/chunk"
	"github.com/Roderland/glox-vm/compiler"
	"github.com/Roderland/glox-vm/utils"
	"github.com/Roderland/glox-vm/vm"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		utils.PrintfErr("Usage: glox [script]\n")
		os.Exit(64)
	}

	bytes, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Printf("Failed to read file '%s'.\n", os.Args[1])
		os.Exit(65)
	}

	interpret(bytes)
}

type InterpretResult uint8

const (
	OK InterpretResult = iota
	RUNTIME_ERROR
	COMPILE_ERROR
)

func interpret(source []byte) InterpretResult {
	ck := chunk.NewChunk()

	if !compiler.Compile(ck, source, true) {
		return COMPILE_ERROR
	}

	if !vm.NewVM(ck, true).Run() {
		return RUNTIME_ERROR
	}

	return OK
}
