package vm

import (
	"fmt"
	"golox/pkg/common"
	"golox/pkg/debug"
)

func (v *VM) ReadConstant() (common.Value, error) {
	i, err := v.ReadByte()
	if err != nil {
		return common.NilValue(), err
	}
	return v.Chunk.Constants.Values[i], nil
}

func (v *VM) ReadConstantLong() (common.Value, error) {
	b1, err := v.ReadByte()
	if err != nil {
		return common.NilValue(), err
	}

	b2, err := v.ReadByte()
	if err != nil {
		return common.NilValue(), err
	}

	b3, err := v.ReadByte()
	if err != nil {
		return common.NilValue(), err
	}

	index := (int(b1) << 16) | (int(b2) << 8) | int(b3)

	return v.Chunk.Constants.Values[index], nil
}

func (v *VM) ConstantOperation(op byte) InterpretResult {
	var c common.Value
	var err error
	
	switch op {
	case byte(common.OpConstant):
		c, err = v.ReadConstant()
		break
	case byte(common.OpConstantLong):
		c, err = v.ReadConstantLong()
		break
	}

	if err != nil {
		return InterpretRuntimeError
	}

	v.Stack.Push(c)
	if v.DebugMode {
		fmt.Println(debug.PrintValue(c))
	}

	return InterpretOk
}
