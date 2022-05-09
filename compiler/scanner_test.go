package compiler

import (
	"github.com/Roderland/glox-vm/utils"
	"testing"
)

func TestScanner(t *testing.T) {
	source := "// fib 返回一个闭包函数，\n// 返回的函数每次调用都会返回下一个斐波那契（Fibonacci）数。\nfun fib() {\n\tvar a = 0;" +
		"\n\tvar b = 1;\n\n\tfun calc() {\n\t\tvar c = b;\n\t\tb = a+b;\n\t\ta = c;\n\n\t\treturn a;\n\t}\n\n\treturn " +
		"calc;\n}\n\nvar f = fib();\n\nprint f();\nprint f();\nprint f();\nprint f();\nprint f();\nprint f();"
	bs := []byte(source)
	scn := newScanner(&bs)
	line := -1
	for {
		token := scn.scanToken()
		if token.line != line {
			utils.PrintfDbg("%4d ", token.line)
			line = token.line
		} else {
			utils.PrintfDbg("   | ")
		}
		utils.PrintfDbg("%2d '%s'\n", token.tp, token.lexeme)

		if token.tp == TOKEN_EOF {
			break
		}
	}
}
