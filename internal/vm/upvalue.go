package vm

import "golox/pkg/bytecode"

func (v *VM) captureUpvalue(localSlot int) *bytecode.Upvalue {
	var prevUpvalue *bytecode.Upvalue
	upvalue := v.OpenUpvalues

	for upvalue != nil && upvalue.Slot > localSlot {
		prevUpvalue = upvalue
		upvalue = upvalue.Next
	}

	if upvalue != nil && upvalue.Slot == localSlot {
		return upvalue
	}

	createdUpvalue := v.allocateUpvalue(localSlot, v.Stack.Slot(localSlot))
	createdUpvalue.Next = upvalue

	if prevUpvalue == nil {
		v.OpenUpvalues = createdUpvalue
	} else {
		prevUpvalue.Next = createdUpvalue
	}

	return createdUpvalue
}

func (v *VM) closeUpvalues(lastSlot int) {
	for v.OpenUpvalues != nil && v.OpenUpvalues.Slot >= lastSlot {
		upvalue := v.OpenUpvalues
		upvalue.Closed = *upvalue.Location
		upvalue.Location = &upvalue.Closed
		v.OpenUpvalues = upvalue.Next
	}
}
