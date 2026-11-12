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
├── main.go          # 示例入口，演示栈与队列
├── ds/              # 数据结构实现
│   ├── stack.go     # 泛型栈（后进先出）
│   ├── stack_test.go
│   ├── queue.go     # 泛型队列（先进先出）
│   └── queue_test.go
├── go.mod
└── README.md
```

## 已实现

| 结构 | 位置 | 说明 |
| --- | --- | --- |
| Stack | `ds/stack.go` | 泛型栈，`Push` / `Pop` / `Peek` / `Len` / `IsEmpty` |
| Queue | `ds/queue.go` | 泛型队列，`Enqueue` / `Dequeue` / `Front` / `Len` / `IsEmpty` |

## 待办

- [ ] 链表（单链表 / 双链表）
- [ ] 二叉搜索树与遍历
- [ ] 哈希表
- [ ] 常见排序算法
- [ ] 图与最短路径

## 许可

暂未指定开源许可证。
