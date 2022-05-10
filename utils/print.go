package utils

import (
	"fmt"
	"os"
)

func PrintfDbg(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	fmt.Printf("%c[%d;%d;%dm%s%c[0m", 0x1B, 0, 0, 36, msg, 0x1B)
}

func PrintfErr(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	_, err := fmt.Fprintf(os.Stderr, "%c[%d;%d;%dm%s%c[0m", 0x1B, 0, 0, 31, msg, 0x1B)
	if err != nil {
		return
	}
}
