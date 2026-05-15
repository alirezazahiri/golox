package common

type ObjType byte

const (
	ObjStringType ObjType = iota + 1
)

type Obj interface {
	Type() ObjType
}

func (t ObjType) String() string {
	switch t {
	case ObjStringType:
		return "string"
	default:
		return ""
	}
}
