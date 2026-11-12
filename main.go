package main

import (
	"fmt"

	"github.com/1303-yzym/DataStructure/bplustree"
	"github.com/1303-yzym/DataStructure/ds"
)

func main() {
	demoStack()
	demoQueue()
	demoBPlusTree()
}

func demoStack() {
	stack := ds.NewStack[int]()
	for i := 1; i <= 5; i++ {
		stack.Push(i)
	}
	fmt.Printf("栈内容（底 -> 顶）：%s\n", stack)
	for !stack.IsEmpty() {
		v, _ := stack.Pop()
		fmt.Printf("弹出：%d\n", v)
	}
}

func demoQueue() {
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

func demoBPlusTree() {
	tree := bplustree.New[int, string](4)
	for _, k := range []int{50, 20, 80, 10, 30, 70, 90, 25, 35, 65} {
		tree.Put(k, fmt.Sprintf("值%d", k))
	}
	fmt.Printf("B+ 树大小：%d，高度：%d\n", tree.Len(), tree.Height())

	fmt.Print("有序遍历：")
	first := true
	tree.Ascend(func(k int, v string) bool {
		if !first {
			fmt.Print(" ")
		}
		first = false
		fmt.Print(k)
		return true
	})
	fmt.Println()

	fmt.Print("范围查询 [20, 60)：")
	tree.Range(20, 60, func(k int, v string) bool {
		fmt.Printf("%d=%s ", k, v)
		return true
	})
	fmt.Println()

	tree.Delete(20)
	_, ok := tree.Get(20)
	fmt.Printf("删除 20 后查询：存在=%v，大小=%d\n", ok, tree.Len())
}
