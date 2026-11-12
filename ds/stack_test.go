package ds

import "testing"

func TestStack(t *testing.T) {
	s := NewStack[int]()
	if !s.IsEmpty() || s.Len() != 0 {
		t.Fatalf("new stack should be empty, got len=%d", s.Len())
	}

	for i := 1; i <= 3; i++ {
		s.Push(i)
	}
	if s.Len() != 3 {
		t.Fatalf("want len 3, got %d", s.Len())
	}
	if got, ok := s.Peek(); !ok || got != 3 {
		t.Fatalf("peek want 3, got %v ok=%v", got, ok)
	}

	for want := 3; want >= 1; want-- {
		got, ok := s.Pop()
		if !ok || got != want {
			t.Fatalf("pop want %d, got %v ok=%v", want, got, ok)
		}
	}

	if _, ok := s.Pop(); ok {
		t.Fatal("pop on empty stack should report false")
	}
}
