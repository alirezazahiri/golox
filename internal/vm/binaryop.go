package vm

import (
	"golox/pkg/common"
	"golox/pkg/errors"
)

// BinaryOperation handles arithmetic operations like division, multiplication, summation and subtraction.
func (v *VM) BinaryOperation(op rune) InterpretResult {
	right := v.Stack.Pop()
	left := v.Stack.Pop()

	if !right.IsNumber() || !left.IsNumber() {
		v.runtimeError(errors.OperandsMustBeNumbers)
		return InterpretRuntimeError
	}

	switch op {
	case '+':
		v.Stack.Push(common.NumberValue(left.As.Number + right.As.Number))
		return InterpretOk
	case '-':
		v.Stack.Push(common.NumberValue(left.As.Number - right.As.Number))
		return InterpretOk
	case '*':
		v.Stack.Push(common.NumberValue(left.As.Number * right.As.Number))
		return InterpretOk
	case '/':
		v.Stack.Push(common.NumberValue(left.As.Number / right.As.Number))
		return InterpretOk
	default:
		return InterpretOk
	}
}
