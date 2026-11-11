package vm

import (
	"golox/pkg/bytecode"
)

type CallFrame struct {
	Closure *bytecode.Closure
	IP      int
	Slots   int
}

func (v *VM) currentFrame() *CallFrame {
	return &v.Frames[v.FrameCount-1]
}

func (v *VM) currentChunk() *bytecode.Chunk {
	return v.currentFrame().Closure.Function.Chunk
}
