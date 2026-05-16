package common

type ObjString struct {
	Content string
	Length  int
	Hash    uint32
}

func (s *ObjString) Type() ObjType {
	return ObjStringType
}

func StringObjValue(s *ObjString) Value {
	return Value{
		Type: ValObj,
		As: Union{
			Obj: s,
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

func HashString(key string) uint32 {
	var hash uint32 = 2166136261

	for i := range len(key) {
		hash ^= uint32(key[i])
		hash *= 16777619
	}

	return hash
}
