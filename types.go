package glox_vm

const (
	VAL_NIL ValueType = iota
	VAL_BOOL
	VAL_NUMBER
	VAL_STRING
	VAL_FUNCTION
)

type (
	ValueType uint8
	Value struct {
		typ  ValueType
		data interface{}
	}
)

func NewNil() Value {
	return Value{VAL_NIL, nil}
}

func NewBool(b bool) Value {
	return Value{VAL_BOOL, b}
}

func NewNumber(num float64) Value {
	return Value{VAL_NUMBER, num}
}

func NewString(s string) Value {
	return Value{VAL_STRING, s}
}

func NewFunction(fd FuncData) Value {
	return Value{VAL_FUNCTION, fd}
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

func (value Value) IsString() bool {
	return value.typ == VAL_STRING
}

func (value Value) IsFunction() bool {
	return value.typ == VAL_FUNCTION
}

func (value Value) AsBool() bool {
	return value.data.(bool)
}

func (value Value) AsNumber() float64 {
	return value.data.(float64)
}

func (value Value) AsString() string {
	return value.data.(string)
}

func (value Value) AsFunction() FuncData {
	return value.data.(FuncData)
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
	case VAL_STRING:
		return value.AsString() == target.AsString()
	default:
		return false
	}
}

