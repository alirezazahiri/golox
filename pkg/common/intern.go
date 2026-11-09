package common

// StringInterner deduplicates string objects during compile and runtime.
type StringInterner interface {
	InternString(chars string) *ObjString
}
