package vm

import (
	"fmt"
	"os"
)

func (v *VM) runtimeError(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format, args...)
	fmt.Fprintln(os.Stderr)

	for i := v.FrameCount - 1; i >= 0; i-- {
		frame := &v.Frames[i]
		function := frame.Function
		instruction := frame.IP - 1
		line := function.Chunk.GetLine(instruction)

		fmt.Fprintf(os.Stderr, "[line %d] in ", line)
		if function.Name == nil {
			fmt.Fprintln(os.Stderr, "script")
		} else {
			fmt.Fprintf(os.Stderr, "%s()\n", function.Name.Content)
		}
	}

	v.resetStack()
}

func (v *VM) resetStack() {
	v.Stack.Reset()
	v.FrameCount = 0
}
