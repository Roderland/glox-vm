package chunk

import "fmt"

type Value float64

func (val Value) String() string {
	return fmt.Sprintf("%g", val)
}
