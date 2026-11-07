package vm

import (
	"golox/pkg/bytecode"
	"golox/pkg/common"
	"golox/pkg/ds"
	"golox/pkg/table"
	vmval "golox/pkg/vm"
)

type VM struct {
	Chunk *bytecode.Chunk

	// Instruction Pointer
	IP int

	Stack *ds.Stack[vmval.Value]

	Strings *table.Table
	Globals *table.Table

	DebugMode bool
}

func New() *VM {
	return &VM{
		Chunk:     nil,
		Stack:     ds.NewStack[vmval.Value](),
		DebugMode: false,
		Strings:   table.New(),
		Globals:   table.New(),
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

	vm.Strings.Set(str, vmval.NilValue())
	return str
}

func (v *VM) WriteChunk(c *bytecode.Chunk) {
	v.Chunk = c
}

func (v *VM) Free() {
	v.Chunk = nil
}

func (v *VM) ActivateDebugMode() {
	v.DebugMode = true
}
