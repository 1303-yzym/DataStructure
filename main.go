package main

import (
	"fmt"

	"github.com/1303-yzym/DataStructure/ds"
)

func main() {
	stack := ds.NewStack[int]()
	for i := 1; i <= 5; i++ {
		stack.Push(i)
	}
	fmt.Printf("栈内容（底 -> 顶）：%s\n", stack)
	for !stack.IsEmpty() {
		v, _ := stack.Pop()
		fmt.Printf("弹出：%d\n", v)
	}

	queue := ds.NewQueue[string]()
	for _, s := range []string{"Go", "数据结构", "算法"} {
		queue.Enqueue(s)
	}
	fmt.Printf("队列长度：%d\n", queue.Len())
	for !queue.IsEmpty() {
		v, _ := queue.Dequeue()
		fmt.Printf("出队：%s\n", v)
	}
}
