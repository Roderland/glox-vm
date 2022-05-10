package compiler

import (
	"github.com/Roderland/glox-vm/chunk"
	"testing"
)

func TestCompile(t *testing.T) {
	source := "-1 +2 * 3"
	Compile(chunk.NewChunk(), []byte(source), true)
}
