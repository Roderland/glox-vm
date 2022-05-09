package main

import (
	"github.com/Roderland/glox-vm/utils"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		utils.PrintfError("Usage: glox [script]\n")
		os.Exit(64)
	}

}
