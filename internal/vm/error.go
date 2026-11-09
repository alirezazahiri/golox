package vm

import (
	"fmt"
	"os"
)

func (v *VM) runtimeError(format string, args ...any) {
	frame := v.currentFrame()
	instruction := frame.IP - 1
	line := frame.Function.Chunk.GetLine(instruction)

	fmt.Fprintf(os.Stderr, "[line %d] in script\n", line)
	fmt.Fprintf(os.Stderr, format, args...)
	fmt.Fprintln(os.Stderr)

	v.resetStack()
}

func (v *VM) resetStack() {
	v.Stack.Reset()
	v.FrameCount = 0
}
