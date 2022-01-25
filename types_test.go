package glox_vm

import (
	"fmt"
	"testing"
)

func TestValue(t *testing.T) {
	value := NewBool(false)
	fmt.Println(value.AsNumber())
	value = NewNumber(1.23)
	fmt.Println(value.AsNumber())
	if value.IsNumber() {
		fmt.Println(value.data)
	}
}
