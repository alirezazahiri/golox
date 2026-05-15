package vm

import (
	"fmt"
	"os"
)

func (v *VM) runtimeError(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format, args...)
	fmt.Fprintln(os.Stderr)

	instruction := v.IP - 1
	line := v.Chunk.GetLine(instruction)

	fmt.Fprintf(os.Stderr, "[line %d] in script\n", line)

	v.resetStack()
}

func (v *VM) resetStack() {
	v.Stack.Reset()
}
