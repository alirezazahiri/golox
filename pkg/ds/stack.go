package ds

import "fmt"

type Stack[T any] struct {
	values []T
	top    int
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{
		values: make([]T, 0, StackMax),
		top:    0,
	}
}

func (s *Stack[T]) Push(element T) {
	if s.top >= StackMax {
		panic(fmt.Sprintf("maximum stack size is: %d", cap(s.values)))
	}
	s.values = append(s.values[:s.top], element)
	s.top++
}

func (s *Stack[T]) Pop() T {
	if s.top == 0 {
		panic("stack is empty")
	}
	s.top--
	last := s.values[s.top]
	return last
}

func (s *Stack[T]) GetAt(index int) T {
	return s.values[index]
}

func (s *Stack[T]) UpdateTop(mapper func(v T) T) {
	lastIndex := s.top - 1
	s.values[lastIndex] = mapper(s.values[lastIndex])
}

func (s *Stack[T]) Top() int {
	return s.top
}

const (
	StackMax = 256
)
