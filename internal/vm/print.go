package vm

import (
	"fmt"

	"golox/pkg/debug"
	vmval "golox/pkg/vm"
)

func (v *VM) PrintValue() InterpretResult {
	var value vmval.Value

	if v.DebugMode {
		value = v.Stack.GetAt(-1)
	} else {
		value = v.Stack.Pop()
	}

	fmt.Print(debug.PrintValue(value))
	return InterpretOk
}
