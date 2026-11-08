package common

func StringObjValue(s *ObjString) Value {
	return ObjValue(s)
}

func (v Value) IsString() bool {
	if v.Type != ValObj {
		return false
	}
	_, ok := v.As.Obj.(*ObjString)
	return ok
}

func (v Value) AsString() *ObjString {
	return v.As.Obj.(*ObjString)
}
