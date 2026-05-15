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
			v.Stack.Push(common.StringValue(left.AsString().Content + right.AsString().Content))
			break
		}
		v.runtimeError(errors.ExpectedOperandToBeOfType, left.As.Obj.Type(), right.As.Obj.Type())
		return InterpretRuntimeError
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
		v.Stack.Push(common.BoolValue(v.objsEqual(left, right)))
		break
	}
	return InterpretOk
}

func (v *VM) objsEqual(a, b common.Value) bool {
	if a.IsString() {
		return a.AsString().Content == b.AsString().Content
	}
	return false
}
