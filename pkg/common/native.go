package common

type NativeFn func(argCount int, args []Value) Value

type ObjNative struct {
	Function NativeFn
}

func NewNative(fn NativeFn) *ObjNative {
	return &ObjNative{
		Function: fn,
	}
}

func (*NativeFn) Type() ObjType {
	return ObjNativeFnType
}
