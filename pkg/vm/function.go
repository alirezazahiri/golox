package vm

import "golox/pkg/common"

func FunctionObjValue(fn *common.ObjFunction) Value {
	return ObjValue(fn)
}

func (v Value) IsFunction() bool {
	if v.Type != ValObj {
		return false
	}
	_, ok := v.As.Obj.(*common.ObjFunction)
	return ok
}

func (v Value) AsFunction() *common.ObjFunction {
	return v.As.Obj.(*common.ObjFunction)
}
