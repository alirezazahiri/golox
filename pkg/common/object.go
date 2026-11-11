package common

type ObjType byte

const (
	ObjStringType ObjType = iota + 1
	ObjFunctionType
	ObjNativeFnType
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
	default:
		return ""
	}
}
