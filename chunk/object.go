package chunk

type ObjType uint8
type FunType uint8

const (
	OBJ_FUNCTION ObjType = iota
	OBJ_NATIVE
)

const (
	FUNCTION FunType = iota
	SCRIPT
)

type Object struct {
	ot      ObjType
	next    *Object
	content interface{}
}

func (val Value) AsObject() Object {
	return val.v.(Object)
}

func (val Value) IsObject() bool {
	return val.lt == VAL_OBJECT
}

func NewObject(obj Object) Value {
	return Value{
		lt: VAL_OBJECT,
		v:  obj,
	}
}

func NewFunction(function ObjFunction) Object {
	return Object{
		ot:      OBJ_FUNCTION,
		content: function,
	}
}

func NewNative(native NativeFunction) Object {
	return Object{
		ot:      OBJ_NATIVE,
		content: native,
	}
}

type ObjFunction struct {
	Name  string
	Arity int
	Ck    Chunk
}

type NativeFunction func(args ...Value) Value

func (obj *Object) IsFunction() bool {
	return obj.ot == OBJ_FUNCTION
}

func (obj Object) AsFunction() ObjFunction {
	return obj.content.(ObjFunction)
}

func (obj Object) IsNative() bool {
	return obj.ot == OBJ_NATIVE
}

func (obj Object) AsNative() NativeFunction {
	return obj.content.(NativeFunction)
}

func (of ObjFunction) GetName() string {
	name := of.Name
	if name == "" {
		return "script"
	}
	return name
}

func (obj Object) String() string {
	var str string
	switch obj.ot {
	case OBJ_FUNCTION:
		str = "<fn " + obj.AsFunction().GetName() + ">"
	case OBJ_NATIVE:
		str = "<native fn>"
	}
	return str
}
