package main

import (
	"fmt"
	. "glox-vm"
	"glox-vm/compiler"
	"glox-vm/vm"
	"io/ioutil"
	"os"
)

func main() {
	interpret(readFile(), &compiler.Compiler{})
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

func interpret(source []byte, compiler *compiler.Compiler) InterpretResult {
	script, err := compiler.Compile(source)
	if err != nil {
		return COMPILE_ERROR
	}
	return vm.InitVM(script).Run()
}