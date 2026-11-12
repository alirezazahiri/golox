package common

type ObjType byte

const (
	ObjStringType ObjType = iota + 1
	ObjFunctionType
	ObjNativeFnType
	ObjClosureType
	ObjUpvalueType
)

// Object is the GC header embedded in every heap-allocated Obj.
type Object struct {
	Next   Obj
	Marked bool
}

type Obj interface {
	Type() ObjType
	Header() *Object
}

func (t ObjType) String() string {
	switch t {
	case ObjStringType:
		return "string"
	case ObjFunctionType:
		return "function"
	case ObjNativeFnType:
		return "native"
	case ObjClosureType:
		return "closure"
	case ObjUpvalueType:
		return "upvalue"
	default:
		return ""
	}
}
