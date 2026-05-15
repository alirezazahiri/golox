package common

type ValueType byte

const (
	ValBool ValueType = iota + 1
	ValNil
	ValNumber
)

type Union struct {
	Bool   bool
	Number float64
}

type Value struct {
	Type ValueType
	As   Union
}

func BoolValue(value bool) Value {
	return Value{
		Type: ValBool,
		As: Union{
			Bool: value,
		},
	}
}

func NilValue() Value {
	return Value{
		Type: ValNil,
		As: Union{
			Number: 0,
			Bool:   false,
		},
	}
}

func NumberValue(value float64) Value {
	return Value{
		Type: ValNumber,
		As: Union{
			Number: value,
			Bool: value != 0,
		},
	}
}

func (v Value) IsBool() bool {
	return v.Type == ValBool
}
func (v Value) IsNil() bool {
	return v.Type == ValNil
}
func (v Value) IsNumber() bool {
	return v.Type == ValNumber
}

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
