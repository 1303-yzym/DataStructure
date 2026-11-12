package ds

// Queue 是一个基于切片的泛型队列（先进先出）。
type Queue[T any] struct {
	items []T
}

// NewQueue 创建一个空队列。
func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{}
}

// Enqueue 将元素加入队尾。
func (q *Queue[T]) Enqueue(v T) {
	q.items = append(q.items, v)
}

// Dequeue 从队头取出元素；队列为空时返回零值和 false。
func (q *Queue[T]) Dequeue() (T, bool) {
	var zero T
	if len(q.items) == 0 {
		return zero, false
	}
	v := q.items[0]
	q.items = q.items[1:]
	return v, true
}

// Front 查看队头元素但不取出；队列为空时返回零值和 false。
func (q *Queue[T]) Front() (T, bool) {
	var zero T
	if len(q.items) == 0 {
		return zero, false
	}
	return q.items[0], true
}

// Len 返回队列中元素个数。
func (q *Queue[T]) Len() int {
	return len(q.items)
}

// IsEmpty 判断队列是否为空。
func (q *Queue[T]) IsEmpty() bool {
	return len(q.items) == 0
}
