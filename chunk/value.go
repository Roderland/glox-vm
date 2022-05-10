package chunk

import "fmt"

type LoxType uint8

const (
	VAL_BOOL LoxType = iota
	VAL_NIL
	VAL_NUMBER
	VAL_STRING
)

type Value struct {
	lt LoxType
	v  interface{}
}

var Nil = Value{VAL_NIL, nil}
var False = Value{VAL_BOOL, false}
var True = Value{VAL_BOOL, true}

func (val Value) String() string {
	var str string
	switch val.lt {
	case VAL_BOOL:
		if val.AsBool() {
			str = "true"
		} else {
			str = "false"
		}
	case VAL_NIL:
		str = "nil"
	case VAL_NUMBER:
		str = fmt.Sprintf("%g", val.v)
	}
	return str
}

func NewNumber(f float64) Value {
	return Value{
		lt: VAL_NUMBER,
		v:  f,
	}
}

func NewBool(b bool) Value {
	if b {
		return True
	}
	return False
}

func (val Value) AsNumber() float64 {
	return val.v.(float64)
}

func (val Value) AsBool() bool {
	return val.v.(bool)
}

func (val Value) IsNumber() bool {
	return val.lt == VAL_NUMBER
}

func (val Value) IsNil() bool {
	return val.lt == VAL_NIL
}

func (val Value) IsBool() bool {
	return val.lt == VAL_BOOL
}

func (val Value) IsFalse() bool {
	return val.IsNil() || (val.IsBool() && !val.AsBool())
}

func Equal(a, b Value) bool {
	if a.lt == b.lt {
		switch a.lt {
		case VAL_NIL:
			return true
		case VAL_BOOL:
			return a.AsBool() == b.AsBool()
		case VAL_NUMBER:
			return a.AsNumber() == b.AsNumber()
		default:
			return false
		}
	}
	return false
}
