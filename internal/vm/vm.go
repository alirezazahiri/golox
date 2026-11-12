package vm

import (
	"golox/internal/parser"
	"golox/pkg/bytecode"
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

	Strings      *table.Table
	Globals      *table.Table
	OpenUpvalues *bytecode.Upvalue
	InitString   *common.ObjString

	objects        common.Obj
	bytesAllocated int
	nextGC         int
	grayStack      []common.Obj
	compilerRoot   **parser.Compiler

	DebugMode bool
}

func New() *VM {
	vm := &VM{
		Stack:          ds.NewStack[common.Value](),
		DebugMode:      false,
		Strings:        table.New(),
		Globals:        table.New(),
		FrameCount:     0,
		nextGC:         initialGCThreshold,
		bytesAllocated: 0,
	}
	vm.InitString = vm.InternString("init")
	vm.defineNatives()
	return vm
}

func (v *VM) Free() {
	v.FrameCount = 0
	v.Stack.Reset()
	v.OpenUpvalues = nil
	v.compilerRoot = nil
	v.InitString = nil
	v.freeObjects()
	v.Strings = table.New()
	v.Globals = table.New()
	v.nextGC = initialGCThreshold
}

func (v *VM) ActivateDebugMode() {
	v.DebugMode = true
}

// SetNextGC sets the GC threshold (tests only).
func (v *VM) SetNextGC(n int) {
	v.nextGC = n
}

// BytesAllocated returns the tracked heap size (tests only).
func (v *VM) BytesAllocated() int {
	return v.bytesAllocated
}

// ObjectCount returns the number of tracked heap objects (tests only).
func (v *VM) ObjectCount() int {
	return v.objectCount()
}
