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
	case '>':
		v.Stack.Push(common.BoolValue(left.As.Number > right.As.Number))
		break
	case '<':
		v.Stack.Push(common.BoolValue(left.As.Number < right.As.Number))
		break
	case '+':
		v.Stack.Push(common.NumberValue(left.As.Number + right.As.Number))
		break
	case '-':
		v.Stack.Push(common.NumberValue(left.As.Number - right.As.Number))
		break
	case '*':
		v.Stack.Push(common.NumberValue(left.As.Number * right.As.Number))
		break
	case '/':
		v.Stack.Push(common.NumberValue(left.As.Number / right.As.Number))
		break
	default:
		return InterpretOk
	}

	return InterpretOk
}

func (v *VM) ValuesEqual() InterpretResult {
	right := v.Stack.Pop()
	left := v.Stack.Pop()

	if left.Type != right.Type {
		v.runtimeError(errors.ExpectedOperandToBeOfType, left.Type, right.Type)
		return InterpretRuntimeError
	}

	switch left.Type {
	case common.ValBool:
		v.Stack.Push(common.BoolValue(left.As.Bool == right.As.Bool))
		break
	case common.ValNil:
		v.Stack.Push(common.BoolValue(true))
		break
	case common.ValNumber:
		v.Stack.Push(common.BoolValue(true))
		break
	}
	return InterpretOk
}
