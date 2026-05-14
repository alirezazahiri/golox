package ds

import "testing"

func TestStackPushAndTop(t *testing.T) {
	s := NewStack[int]()

	s.Push(1)
	s.Push(2)
	s.Push(3)

	if s.Top() != 3 {
		t.Fatalf("expected top to be 3, got %d", s.Top())
	}
}

func TestStackPopOrder(t *testing.T) {
	s := NewStack[int]()

	s.Push(10)
	s.Push(20)
	s.Push(30)

	if v := s.Pop(); v != 30 {
		t.Fatalf("expected 30, got %d", v)
	}

	if v := s.Pop(); v != 20 {
		t.Fatalf("expected 20, got %d", v)
	}

	if v := s.Pop(); v != 10 {
		t.Fatalf("expected 10, got %d", v)
	}
}

func TestStackPushPopMixed(t *testing.T) {
	s := NewStack[int]()

	s.Push(1)
	s.Push(2)

	if v := s.Pop(); v != 2 {
		t.Fatalf("expected 2, got %d", v)
	}

	s.Push(3)

	if v := s.Pop(); v != 3 {
		t.Fatalf("expected 3, got %d", v)
	}

	if v := s.Pop(); v != 1 {
		t.Fatalf("expected 1, got %d", v)
	}
}

func TestStackTopAfterPop(t *testing.T) {
	s := NewStack[int]()

	s.Push(1)
	s.Push(2)
	s.Push(3)

	s.Pop()
	s.Pop()

	if s.Top() != 1 {
		t.Fatalf("expected top index to be 1, got %d", s.Top())
	}
}

func TestStackReuseAfterPop(t *testing.T) {
	s := NewStack[int]()

	s.Push(1)
	s.Push(2)
	s.Push(3)

	s.Pop()
	s.Pop()

	s.Push(4)

	if v := s.Pop(); v != 4 {
		t.Fatalf("expected 4, got %d", v)
	}

	if v := s.Pop(); v != 1 {
		t.Fatalf("expected 1, got %d", v)
	}
}

func TestStackOverflow(t *testing.T) {
	s := NewStack[int]()

	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic on stack overflow, got none, current stack size: %d", len(s.values))
		}
	}()

	for i := range StackMax + 1 {
		s.Push(i)
	}
}
