package vm

import (
	"golox/internal/chunk"
	"golox/pkg/common"
	"golox/pkg/ds"
	"golox/pkg/table"
)

type VM struct {
	Chunk *chunk.Chunk

	// Instruction Pointer
	IP int

	Stack *ds.Stack[common.Value]

	Strings *table.Table

	DebugMode bool
}

func New() *VM {
	return &VM{
		Chunk:     nil,
		Stack:     ds.NewStack[common.Value](),
		DebugMode: false,
		Strings:   table.New(),
	}
}

func (vm *VM) InternString(chars string) *common.ObjString {
	hash := common.HashString(chars)

	if existing := vm.Strings.FindString(chars, hash); existing != nil {
		return existing
	}

	str := &common.ObjString{
		Content: chars,
		Hash:    hash,
	}

	vm.Strings.Set(str, common.NilValue())
	return str
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
