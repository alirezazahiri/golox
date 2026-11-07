package vm

import (
	"golox/pkg/common"
	"golox/pkg/errors"
	vmval "golox/pkg/vm"
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
		v.Stack.Push(vmval.BoolValue(left.As.Number > right.As.Number))
	case '<':
		v.Stack.Push(vmval.BoolValue(left.As.Number < right.As.Number))
	case '-':
		v.Stack.Push(vmval.NumberValue(left.As.Number - right.As.Number))
	case '*':
		v.Stack.Push(vmval.NumberValue(left.As.Number * right.As.Number))
	case '/':
		v.Stack.Push(vmval.NumberValue(left.As.Number / right.As.Number))
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
	case vmval.ValNumber:
		v.Stack.Push(vmval.NumberValue(left.As.Number + right.As.Number))
	case vmval.ValObj:
		if left.IsString() && right.IsString() {
			str := v.InternString(left.AsString().Content + right.AsString().Content)
			v.Stack.Push(vmval.StringObjValue(str))
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

	if left.Type == vmval.ValObj {
		if left.As.Obj.Type() != right.As.Obj.Type() {
			v.runtimeError(errors.ExpectedOperandToBeOfType, left.As.Obj.Type(), right.As.Obj.Type())
			return InterpretRuntimeError
		}
	}

	switch left.Type {
	case vmval.ValBool:
		v.Stack.Push(vmval.BoolValue(left.As.Bool == right.As.Bool))
	case vmval.ValNil:
		v.Stack.Push(vmval.BoolValue(true))
	case vmval.ValNumber:
		v.Stack.Push(vmval.BoolValue(left.As.Number == right.As.Number))
	case vmval.ValObj:
		v.Stack.Push(vmval.BoolValue(left.As.Obj == right.As.Obj))
	}
	return InterpretOk
}
