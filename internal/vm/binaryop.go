package vm

func (v *VM) BinaryOperation(op rune) {
	right := v.Stack.Pop()
	left := v.Stack.Pop()

	switch op {
	case '+':
		v.Stack.Push(left + right)
		return
	case '-':
		v.Stack.Push(left - right)
		return
	case '*':
		v.Stack.Push(left * right)
		return
	case '/':
		v.Stack.Push(left / right)
		return
	default:
		return
	}
}
