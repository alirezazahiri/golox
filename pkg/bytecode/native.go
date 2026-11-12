package bytecode

import "golox/pkg/common"

type NativeFn func(argCount int, args []common.Value) common.Value

type ObjNative struct {
	common.Object
	Function NativeFn
}

func NewNative(fn NativeFn) *ObjNative {
	return &ObjNative{
		Function: fn,
	}
}

func (*ObjNative) Type() common.ObjType {
	return common.ObjNativeFnType
}

func (n *ObjNative) Header() *common.Object {
	return &n.Object
}

func NativeFnValue(n *ObjNative) common.Value {
	return common.ObjValue(n)
}

func AsNative(v common.Value) *ObjNative {
	return v.As.Obj.(*ObjNative)
}
