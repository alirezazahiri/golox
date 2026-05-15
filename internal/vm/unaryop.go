package vm

import (
	"golox/pkg/common"
	"golox/pkg/errors"
)

// UnaryOperation handles unery operations like negation (-) or not (!)
func (v *VM) UnaryOperation(op rune) InterpretResult {
	switch op {
	case '-':
		if !v.Stack.GetAt(-1).IsNumber() {
			v.runtimeError(errors.OperandMustBeNumber)
			return InterpretRuntimeError
		}
		v.Stack.UpdateTop(func(v common.Value) common.Value { return common.NumberValue(-v.As.Number) })
		break
	case '!':
		v.Stack.UpdateTop(func(v common.Value) common.Value { return common.BoolValue(!v.As.Bool) })
		break
	}

	return InterpretOk
}
