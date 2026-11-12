package common

type ObjString struct {
	Object
	Content string
	Length  int
	Hash    uint32
}

func (s *ObjString) Type() ObjType {
	return ObjStringType
}

func (s *ObjString) Header() *Object {
	return &s.Object
}

func HashString(key string) uint32 {
	var hash uint32 = 2166136261

	for i := range len(key) {
		hash ^= uint32(key[i])
		hash *= 16777619
	}

	return hash
}
