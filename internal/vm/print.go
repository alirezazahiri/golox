package vm

import (
	"fmt"

	"golox/pkg/common"
	"golox/pkg/debug"
)

func (v *VM) PrintValue() InterpretResult {
	var value common.Value
	
	if v.DebugMode {
		value = v.Stack.GetAt(-1)
	} else {
		value = v.Stack.Pop()
	}

	fmt.Print(debug.PrintValue(value))
	return InterpretOk
}
