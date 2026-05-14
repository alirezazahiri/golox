package vm

import (
	"golox/internal/chunk"
	"golox/pkg/common"
	"golox/pkg/ds"
)

type VM struct {
	Chunk *chunk.Chunk

	// Instruction Pointer
	IP int

	Stack *ds.Stack[common.Value]

	DebugMode bool
}

func New() *VM {
	return &VM{
		Chunk:     nil,
		Stack:     ds.NewStack[common.Value](),
		DebugMode: false,
	}
}

func (v *VM) WriteChunk(c *chunk.Chunk) {
	v.Chunk = c
}

func (v *VM) Free() {
	v.Chunk = nil
}

func (v *VM) ActivateDebugMode() {
	v.DebugMode = true
}
