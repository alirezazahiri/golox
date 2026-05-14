package common

type Value float64

type ValueArray struct {
	Values []Value
}

func NewValueArray() *ValueArray {
	return &ValueArray{
		Values: make([]Value, 0, 8),
	}
}

func (v *ValueArray) Write(value Value) {
	v.Values = append(v.Values, value)
}

func (v *ValueArray) Free() {
	v.Values = nil
}

func (v *ValueArray) Size() int {
	return len(v.Values)
}