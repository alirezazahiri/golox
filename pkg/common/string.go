package common

type ObjString struct {
	Content string
	Length  int
}

func (s *ObjString) Type() ObjType {
	return ObjStringType
}

func StringValue(s string) Value {
	return Value{
		Type: ValObj,
		As: Union{
			Obj: &ObjString{
				Content: s,
				Length:  len(s),
			},
			Bool: s != "",
		},
	}
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
