package vm

import (
	"fmt"
	"os"
)

func (v *VM) runtimeError(format string, args ...any) {
	frame := v.currentFrame()
	instruction := frame.IP - 1
	line := frame.Function.Chunk.GetLine(instruction)

	name := "<script>"
	if frame.Function.Name != nil && frame.Function.Name.Content != "" {
		name = frame.Function.Name.Content
	}

	fmt.Fprintf(os.Stderr, "[line %d] in %s\n", line, name)
	fmt.Fprintf(os.Stderr, format, args...)
	fmt.Fprintln(os.Stderr)

	v.resetStack()
}

func (v *VM) resetStack() {
	v.Stack.Reset()
	v.FrameCount = 0
}
