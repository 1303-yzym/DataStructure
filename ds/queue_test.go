package ds

import "testing"

func TestQueue(t *testing.T) {
	q := NewQueue[string]()
	if !q.IsEmpty() {
		t.Fatal("new queue should be empty")
	}

	q.Enqueue("a")
	q.Enqueue("b")
	q.Enqueue("c")

	if got, ok := q.Front(); !ok || got != "a" {
		t.Fatalf("front want a, got %v ok=%v", got, ok)
	}

	for _, want := range []string{"a", "b", "c"} {
		got, ok := q.Dequeue()
		if !ok || got != want {
			t.Fatalf("dequeue want %s, got %v ok=%v", want, got, ok)
		}
	}

	if _, ok := q.Dequeue(); ok {
		t.Fatal("dequeue on empty queue should report false")
	}
}
