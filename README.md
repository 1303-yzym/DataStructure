# DataStructure

用 Go 实现与练习常用数据结构，并配套单元测试。

## 环境要求

- Go 1.25 及以上

## 快速开始

```bash
# 运行示例
go run .

# 运行全部测试
go test ./...

# 静态检查
go vet ./...
```

## 目录结构

```
.
├── main.go              # 示例入口，演示栈、队列与 B+ 树
├── ds/                  # 基础数据结构
│   ├── stack.go         # 泛型栈（后进先出）
│   ├── stack_test.go
│   ├── queue.go         # 泛型队列（先进先出）
│   └── queue_test.go
├── bplustree/           # B+ 树
│   ├── bplustree.go     # 泛型 B+ 树实现
│   └── bplustree_test.go
├── go.mod
└── README.md
```

## 已实现

| 结构 | 位置 | 说明 |
| --- | --- | --- |
| Stack | `ds/stack.go` | 泛型栈，`Push` / `Pop` / `Peek` / `Len` / `IsEmpty` |
| Queue | `ds/queue.go` | 泛型队列，`Enqueue` / `Dequeue` / `Front` / `Len` / `IsEmpty` |
| B+ 树 | `bplustree/bplustree.go` | 泛型 B+ 树，`Put` / `Get` / `Delete` / `Ascend` / `Range` / `Height` |

## B+ 树

```go
tree := bplustree.New[int, string](4) // 4 阶：每个节点最多 3 个键
tree.Put(10, "值10")
tree.Put(20, "值20")

v, ok := tree.Get(20)                 // 查找
tree.Delete(10)                        // 删除

tree.Ascend(func(k int, v string) bool { // 升序遍历
    return true // 返回 false 提前结束
})

tree.Range(10, 30, func(k int, v string) bool { // 范围查询 [10, 30)
    return true
})
```

实现要点：所有键值对存放在叶节点，叶节点用 `next` 指针串成有序链表，
因此遍历与范围查询只需沿叶子链扫描；插入满则分裂并把分隔键上移，
删除后若节点键数不足，优先向兄弟借键，借不到再与兄弟合并。

## 待办

- [ ] 链表（单链表 / 双链表）
- [ ] 二叉搜索树与遍历
- [ ] 哈希表
- [ ] 常见排序算法
- [ ] 图与最短路径

## 许可

暂未指定开源许可证。
