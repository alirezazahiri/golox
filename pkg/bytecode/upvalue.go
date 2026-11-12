package bytecode

import "golox/pkg/common"

type Upvalue struct {
	common.Object
	Location *common.Value
	Closed   common.Value
	Next     *Upvalue
	Slot     int
}

func (u *Upvalue) Type() common.ObjType {
	return common.ObjUpvalueType
}

func (u *Upvalue) Header() *common.Object {
	return &u.Object
}

func NewUpvalue(slot int, location *common.Value) *Upvalue {
	return &Upvalue{
		Location: location,
		Slot:     slot,
		Closed:   common.NilValue(),
	}
}
