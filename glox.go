package main

import (
	"fmt"
	"github.com/Roderland/glox-vm/compiler"
	"github.com/Roderland/glox-vm/utils"
	"github.com/Roderland/glox-vm/vm"
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

func interpret(source []byte) InterpretResult {
	function, ok := compiler.Compile(source, true)
	if !ok {
		return COMPILE_ERROR
	}

	fmt.Println("====================== output ======================")
	if !vm.Do(function, true) {
		return RUNTIME_ERROR
	}

	return OK
}
