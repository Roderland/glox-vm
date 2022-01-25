package glox_vm

import (
	"unsafe"
)

const (
	VAL_NIL ValueType = iota
	VAL_BOOL
	VAL_NUMBER
)

type (
	ValueType uint8
	Value struct {
		typ  ValueType
		data [8]byte
	}
)

var (
	nilVal = [8]byte{}
	falseVal = [8]byte{}
	trueVal = [8]byte{1}
)

func NewNil() Value {
	return Value{VAL_NIL, nilVal}
}

func NewBool(b bool) Value {
	if b {
		return Value{VAL_BOOL, trueVal}
	}
	return Value{VAL_BOOL, falseVal}
}

func NewNumber(num float64) Value {
	return Value{VAL_NUMBER, *(*[8]byte)(unsafe.Pointer(&num))}
}

func (value Value) IsNil() bool {
	return value.typ == VAL_NIL
}

func (value Value) IsBool() bool {
	return value.typ == VAL_BOOL
}

func (value Value) IsNumber() bool {
	return value.typ == VAL_NUMBER
}

func (value Value) AsBool() bool {
	return value.data != falseVal
}

func (value Value) AsNumber() float64 {
	return *(*float64)(unsafe.Pointer(&value.data))
}

func (value Value) IsFalse() bool {
	return value.IsNil() || (value.IsBool() && !value.AsBool())
}

func (value Value) Equal(target Value) bool {
	if value.typ != target.typ {
		return false
	}
	switch value.typ {
	case VAL_NIL: return true
	case VAL_BOOL:
		return value.AsBool() == target.AsBool()
	case VAL_NUMBER:
		return value.AsNumber() == target.AsNumber()
	default:
		return false
	}
}

