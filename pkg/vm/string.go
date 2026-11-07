package vm

import "golox/pkg/common"

func StringObjValue(s *common.ObjString) Value {
	return ObjValue(s)
}

func (v Value) IsString() bool {
	if v.Type != ValObj {
		return false
	}
	_, ok := v.As.Obj.(*common.ObjString)
	return ok
}

func (v Value) AsString() *common.ObjString {
	return v.As.Obj.(*common.ObjString)
}
