// Package ds 提供常用的数据结构实现。
package ds

import "fmt"

// Stack 是一个基于切片的泛型栈（后进先出）。
type Stack[T any] struct {
	items []T
}

// NewStack 创建一个空栈。
func NewStack[T any]() *Stack[T] {
	return &Stack[T]{}
}

// Push 将元素压入栈顶。
func (s *Stack[T]) Push(v T) {
	s.items = append(s.items, v)
}

// Pop 弹出栈顶元素；栈为空时返回零值和 false。
func (s *Stack[T]) Pop() (T, bool) {
	var zero T
	if len(s.items) == 0 {
		return zero, false
	}
	last := len(s.items) - 1
	v := s.items[last]
	s.items = s.items[:last]
	return v, true
}

// Peek 查看栈顶元素但不弹出；栈为空时返回零值和 false。
func (s *Stack[T]) Peek() (T, bool) {
	var zero T
	if len(s.items) == 0 {
		return zero, false
	}
	return s.items[len(s.items)-1], true
}

// Len 返回栈中元素个数。
func (s *Stack[T]) Len() int {
	return len(s.items)
}

// IsEmpty 判断栈是否为空。
func (s *Stack[T]) IsEmpty() bool {
	return len(s.items) == 0
}

// String 返回栈的可读形式，栈底在左、栈顶在右。
func (s *Stack[T]) String() string {
	return fmt.Sprint(s.items)
}
