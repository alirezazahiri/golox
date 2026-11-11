package vm

import (
	"golox/pkg/common"
	"golox/pkg/ds"
	"golox/pkg/table"
)

const (
	FramesMax  = 64
	Uint8Count = 256
	StackMax   = FramesMax * Uint8Count
)

type VM struct {
	Frames     [FramesMax]CallFrame
	FrameCount int

	Stack *ds.Stack[common.Value]

	Strings *table.Table
	Globals *table.Table

	DebugMode bool
}

func New() *VM {
	vm := &VM{
		Stack:      ds.NewStack[common.Value](),
		DebugMode:  false,
		Strings:    table.New(),
		Globals:    table.New(),
		FrameCount: 0,
	}
	vm.defineNatives()
	return vm
}

func (v *VM) InternString(chars string) *common.ObjString {
	hash := common.HashString(chars)

	if existing := v.Strings.FindString(chars, hash); existing != nil {
		return existing
	}

	str := &common.ObjString{
		Content: chars,
		Length:  len(chars),
		Hash:    hash,
	}

	v.Strings.Set(str, common.NilValue())
	return str
}

func (v *VM) Free() {
	v.FrameCount = 0
	v.Stack.Reset()
}

func (v *VM) ActivateDebugMode() {
	v.DebugMode = true
}
