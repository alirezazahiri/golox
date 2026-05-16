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

// SumOperation handles the summation of strings and numbers
func (v *VM) SumOperation() InterpretResult {
	right := v.Stack.Pop()
	left := v.Stack.Pop()

	if left.Type != right.Type {
		v.runtimeError(errors.ExpectedOperandToBeOfType, left.Type, right.Type)
		return InterpretRuntimeError
	}

	switch left.Type {
	case common.ValNumber:
		v.Stack.Push(common.NumberValue(left.As.Number + right.As.Number))
		break
	case common.ValObj:
		if left.IsString() && right.IsString() {
			str := v.InternString(left.AsString().Content + right.AsString().Content)
			v.Stack.Push(common.StringObjValue(str))
			break
		}
		v.runtimeError(errors.ExpectedOperandsToBeOfType, common.ObjStringType, left.As.Obj.Type(), right.As.Obj.Type())
		return InterpretRuntimeError
	}

	return InterpretOk
}

// ValuesEqual is responsible for checking the equality of two values of any type
func (v *VM) ValuesEqual() InterpretResult {
	right := v.Stack.Pop()
	left := v.Stack.Pop()

	if left.Type != right.Type {
		v.runtimeError(errors.ExpectedOperandToBeOfType, left.Type, right.Type)
		return InterpretRuntimeError
	}

	if left.Type == common.ValObj {
		if left.As.Obj.Type() != right.As.Obj.Type() {
			v.runtimeError(errors.ExpectedOperandToBeOfType, left.As.Obj.Type(), right.As.Obj.Type())
			return InterpretRuntimeError
		}
	}

	switch left.Type {
	case common.ValBool:
		v.Stack.Push(common.BoolValue(left.As.Bool == right.As.Bool))
		break
	case common.ValNil:
		v.Stack.Push(common.BoolValue(true))
		break
	case common.ValNumber:
		v.Stack.Push(common.BoolValue(left.As.Number == right.As.Number))
		break
	case common.ValObj:
		v.Stack.Push(common.BoolValue(left.As.Obj == right.As.Obj))
		break
	}
	return InterpretOk
}
