package utils

import "fmt"

func PrintfDbg(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	fmt.Printf("%c[%d;%d;%dm%s%c[0m", 0x1B, 0, 0, 36, msg, 0x1B)
}

func PrintfError(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	fmt.Printf("%c[%d;%d;%dm%s%c[0m", 0x1B, 0, 0, 30, msg, 0x1B)
}
