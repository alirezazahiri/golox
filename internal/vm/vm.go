package vm

import "golox/internal/chunk"

type VM struct {
	Chunk *chunk.Chunk
	IP     int
}

func New() *VM {
	return &VM{
		Chunk: nil,
	}
}

func (v *VM) Free() {
	v.Chunk = nil
}
