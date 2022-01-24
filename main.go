package main

import (
	"fmt"
	"glox-vm/compiler"
	"io/ioutil"
	"os"
)

type InterpretResult uint8

const (
	OK InterpretResult = iota
	COMPILE_ERROR
	RUNTIME_ERROR
)

func main() {
	interpret(readFile(), &compiler.Compiler{}, _VM())
}

func readFile() []byte {
	if len(os.Args) != 2 {
		fmt.Println("Usage: glox [InputFile]")
		os.Exit(64)
	}
	bytes, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Printf("Failed to read file '%s'.\n", os.Args[1])
		os.Exit(65)
	}
	return bytes
}

func interpret(source []byte, compiler *compiler.Compiler, vm *VM) InterpretResult {
	compiler.Compile(source)
	return OK
}