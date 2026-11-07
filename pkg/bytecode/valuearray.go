package bytecode

import "golox/pkg/vm"

type ValueArray struct {
	Values []vm.Value
}

func NewValueArray() *ValueArray {
	return &ValueArray{
		Values: make([]vm.Value, 0, 8),
	}
}

func (v *ValueArray) Write(value vm.Value) {
	v.Values = append(v.Values, value)
}

func (v *ValueArray) Free() {
	v.Values = nil
}

func (v *ValueArray) Size() int {
	return len(v.Values)
}
