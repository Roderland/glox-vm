package chunk

import "fmt"

type value float64

func (val value) print() {
	fmt.Printf("%g", val)
}
