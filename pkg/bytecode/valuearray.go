package bytecode

import "golox/pkg/common"

type ValueArray struct {
	Values []common.Value
}

func NewValueArray() *ValueArray {
	return &ValueArray{
		Values: make([]common.Value, 0, 8),
	}
}

func (v *ValueArray) Write(value common.Value) {
	v.Values = append(v.Values, value)
}

func (v *ValueArray) Size() int {
	return len(v.Values)
}
