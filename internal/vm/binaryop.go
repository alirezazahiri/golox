package vm

import "golox/pkg/common"

func (v *VM) BinaryOperation(op rune) {
	a := v.Stack.Pop()

	switch op {
	case '+':
		v.Stack.UpdateTop(func(v common.Value) common.Value { return a + v })
		return
	case '-':
		v.Stack.UpdateTop(func(v common.Value) common.Value { return a - v })
		return
	case '*':
		v.Stack.UpdateTop(func(v common.Value) common.Value { return a * v })
		return
	case '/':
		v.Stack.UpdateTop(func(v common.Value) common.Value { return a / v })
		return
	default:
		return
	}
}
