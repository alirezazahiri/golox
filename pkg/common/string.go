package common

type ObjString struct {
	Content string
	Length  int
	Hash    uint32
}

func (s *ObjString) Type() ObjType {
	return ObjStringType
}

func HashString(key string) uint32 {
	var hash uint32 = 2166136261

	for i := range len(key) {
		hash ^= uint32(key[i])
		hash *= 16777619
	}

	return hash
}
