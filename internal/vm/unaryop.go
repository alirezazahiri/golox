package vm

import (
	"golox/pkg/errors"
	vmval "golox/pkg/vm"
)

// UnaryOperation handles unery operations like negation (-) or not (!)
func (v *VM) UnaryOperation(op rune) InterpretResult {
	switch op {
	case '-':
		if !v.Stack.GetAt(-1).IsNumber() {
			v.runtimeError(errors.OperandMustBeNumber)
			return InterpretRuntimeError
		}
		v.Stack.UpdateTop(func(v vmval.Value) vmval.Value { return vmval.NumberValue(-v.As.Number) })
	case '!':
		v.Stack.UpdateTop(func(v vmval.Value) vmval.Value { return vmval.BoolValue(!v.As.Bool) })
	}

	return InterpretOk
}
