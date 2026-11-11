package common

type ObjType byte

const (
	ObjStringType ObjType = iota + 1
	ObjFunctionType
	ObjNativeFnType
	ObjClosureType
	ObjUpvalueType
)

type Obj interface {
	Type() ObjType
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
