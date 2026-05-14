package vm

func (v *VM) BinaryOperation(op rune) {
	a := v.Stack.Pop()
	b := v.Stack.Pop()

	switch op {
	case '+':
		v.Stack.Push(a + b)
		return
	case '-':
		v.Stack.Push(a - b)
		return
	case '*':
		v.Stack.Push(a * b)
		return
	case '/':
		v.Stack.Push(a / b)
		return
	default:
		return
	}
}
