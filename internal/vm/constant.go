package vm

import "golox/pkg/common"

func (v *VM) ReadConstant() (common.Value, error) {
	i, err := v.ReadByte()
	if err != nil {
		return 0, err
	}
	return v.Chunk.Constants.Values[i], nil
}

func (v *VM) ReadConstantLong() (common.Value, error) {
	b1, err := v.ReadByte()
	if err != nil {
		return 0, err
	}

	b2, err := v.ReadByte()
	if err != nil {
		return 0, err
	}

	b3, err := v.ReadByte()
	if err != nil {
		return 0, err
	}

	index := (int(b1) << 16) | (int(b2) << 8) | int(b3)

	return v.Chunk.Constants.Values[index], nil
}
