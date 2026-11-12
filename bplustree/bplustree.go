// Package bplustree 实现了一棵泛型 B+ 树，支持插入、查找、删除、
// 有序遍历与范围查询。
//
// B+ 树的所有键值对都存放在叶节点中，叶节点之间通过 next 指针按序相连，
// 因此顺序遍历和范围查询可以直接在叶节点链表上完成，无需回溯内部节点。
package bplustree

import "sort"

// Ordered 约束了 B+ 树键可用的类型。
type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 | ~string
}

// node 是 B+ 树的节点。叶节点使用 keys/values/next，
// 内部节点使用 keys/children，两者由 leaf 区分。
type node[K Ordered, V any] struct {
	leaf     bool
	keys     []K
	values   []V
	children []*node[K, V]
	next     *node[K, V]
}

// Tree 是一棵泛型 B+ 树。零值不可用，请通过 New 创建。
type Tree[K Ordered, V any] struct {
	order   int // 每个节点最多的子节点数
	maxKeys int // 每个节点最多的键数
	minKeys int // 除根节点外每个节点最少的键数
	root    *node[K, V]
	size    int
}

// New 创建一棵 B+ 树。order 表示每个节点的最大子节点数，
// 小于 3 时会被提升为 3（此时每个节点最多 2 个键）。
func New[K Ordered, V any](order int) *Tree[K, V] {
	if order < 3 {
		order = 3
	}
	return &Tree[K, V]{
		order:   order,
		maxKeys: order - 1,
		minKeys: (order+1)/2 - 1,
		root:    &node[K, V]{leaf: true},
	}
}

// Len 返回树中键值对的数量。
func (t *Tree[K, V]) Len() int { return t.size }

// Empty 判断树是否为空。
func (t *Tree[K, V]) Empty() bool { return t.size == 0 }

// Height 返回树的层数，空树或只有根叶节点时为 1。
func (t *Tree[K, V]) Height() int {
	h := 1
	for n := t.root; !n.leaf; n = n.children[0] {
		h++
	}
	return h
}

// Get 查找键对应的值。
func (t *Tree[K, V]) Get(key K) (V, bool) {
	n := t.root
	for !n.leaf {
		n = n.children[childIndex(n.keys, key)]
	}
	if i := lowerBound(n.keys, key); i < len(n.keys) && n.keys[i] == key {
		return n.values[i], true
	}
	var zero V
	return zero, false
}

// Put 插入或更新一个键值对，键不存在时返回 true。
func (t *Tree[K, V]) Put(key K, val V) bool {
	splitKey, right, split, created := t.insert(t.root, key, val)
	if split {
		t.root = &node[K, V]{
			keys:     []K{splitKey},
			children: []*node[K, V]{t.root, right},
		}
	}
	return created
}

// Delete 删除键及其对应的值，键不存在时返回 false。
func (t *Tree[K, V]) Delete(key K) bool {
	if !t.deleteFrom(t.root, key) {
		return false
	}
	// 根节点合并后可能只剩一个孩子，需要降低树高。
	for !t.root.leaf && len(t.root.keys) == 0 {
		t.root = t.root.children[0]
	}
	return true
}

// Ascend 按顺序遍历所有键值对，fn 返回 false 时提前终止。
func (t *Tree[K, V]) Ascend(fn func(key K, val V) bool) {
	if t.size == 0 {
		return
	}
	n := t.root
	for !n.leaf {
		n = n.children[0]
	}
	for ; n != nil; n = n.next {
		for i := range n.keys {
			if !fn(n.keys[i], n.values[i]) {
				return
			}
		}
	}
}

// Range 按顺序遍历键落在 [start, end) 内的键值对，fn 返回 false 时提前终止。
func (t *Tree[K, V]) Range(start, end K, fn func(key K, val V) bool) {
	if t.size == 0 || !(start < end) {
		return
	}
	n := t.root
	for !n.leaf {
		n = n.children[childIndex(n.keys, start)]
	}
	for ; n != nil; n = n.next {
		for i := range n.keys {
			if n.keys[i] < start {
				continue
			}
			if !(n.keys[i] < end) {
				return
			}
			if !fn(n.keys[i], n.values[i]) {
				return
			}
		}
	}
}

// insert 向以 n 为根的子树插入键值对。发生分裂时返回上移的分隔键和新兄弟节点。
func (t *Tree[K, V]) insert(n *node[K, V], key K, val V) (splitKey K, right *node[K, V], split bool, created bool) {
	if n.leaf {
		i := lowerBound(n.keys, key)
		if i < len(n.keys) && n.keys[i] == key {
			n.values[i] = val
			return splitKey, nil, false, false
		}
		n.keys = insertAt(n.keys, i, key)
		n.values = insertAt(n.values, i, val)
		t.size++
		if len(n.keys) <= t.maxKeys {
			return splitKey, nil, false, true
		}
		k, r := t.splitLeaf(n)
		return k, r, true, true
	}

	i := childIndex(n.keys, key)
	splitKey, right, split, created = t.insert(n.children[i], key, val)
	if !split {
		return splitKey, nil, false, created
	}
	n.keys = insertAt(n.keys, i, splitKey)
	n.children = insertAt(n.children, i+1, right)
	if len(n.keys) <= t.maxKeys {
		return splitKey, nil, false, created
	}
	k, r := t.splitInternal(n)
	return k, r, true, created
}

// splitLeaf 把溢出的叶节点一分为二，返回右节点的首个键。
func (t *Tree[K, V]) splitLeaf(n *node[K, V]) (K, *node[K, V]) {
	mid := len(n.keys) / 2
	right := &node[K, V]{
		leaf:   true,
		keys:   append([]K(nil), n.keys[mid:]...),
		values: append([]V(nil), n.values[mid:]...),
		next:   n.next,
	}
	n.keys = n.keys[:mid:mid]
	n.values = n.values[:mid:mid]
	n.next = right
	return right.keys[0], right
}

// splitInternal 把溢出的内部节点一分为二，返回上移的分隔键。
func (t *Tree[K, V]) splitInternal(n *node[K, V]) (K, *node[K, V]) {
	mid := len(n.keys) / 2
	splitKey := n.keys[mid]
	right := &node[K, V]{
		keys:     append([]K(nil), n.keys[mid+1:]...),
		children: append([]*node[K, V](nil), n.children[mid+1:]...),
	}
	n.keys = n.keys[:mid:mid]
	n.children = n.children[:mid+1 : mid+1]
	return splitKey, right
}

// deleteFrom 删除以 n 为根的子树中的键，返回是否删除成功。
func (t *Tree[K, V]) deleteFrom(n *node[K, V], key K) bool {
	if n.leaf {
		i := lowerBound(n.keys, key)
		if i >= len(n.keys) || n.keys[i] != key {
			return false
		}
		n.keys = append(n.keys[:i], n.keys[i+1:]...)
		n.values = append(n.values[:i], n.values[i+1:]...)
		t.size--
		return true
	}

	i := childIndex(n.keys, key)
	if !t.deleteFrom(n.children[i], key) {
		return false
	}
	if len(n.children[i].keys) < t.minKeys {
		t.rebalance(n, i)
	}
	return true
}

// rebalance 修复第 i 个孩子键数不足的问题：优先向兄弟借键，否则合并。
func (t *Tree[K, V]) rebalance(n *node[K, V], i int) {
	var left, right *node[K, V]
	if i > 0 {
		left = n.children[i-1]
	}
	if i < len(n.children)-1 {
		right = n.children[i+1]
	}
	switch {
	case left != nil && len(left.keys) > t.minKeys:
		t.borrowFromLeft(n, i)
	case right != nil && len(right.keys) > t.minKeys:
		t.borrowFromRight(n, i)
	case left != nil:
		t.merge(n, i-1)
	default:
		t.merge(n, i)
	}
}

// borrowFromLeft 从左侧兄弟节点借一个键给第 i 个孩子。
func (t *Tree[K, V]) borrowFromLeft(n *node[K, V], i int) {
	child, left := n.children[i], n.children[i-1]
	if child.leaf {
		child.keys = insertAt(child.keys, 0, left.keys[len(left.keys)-1])
		child.values = insertAt(child.values, 0, left.values[len(left.values)-1])
		left.keys = left.keys[:len(left.keys)-1]
		left.values = left.values[:len(left.values)-1]
		n.keys[i-1] = child.keys[0]
		return
	}
	child.keys = insertAt(child.keys, 0, n.keys[i-1])
	child.children = insertAt(child.children, 0, left.children[len(left.children)-1])
	n.keys[i-1] = left.keys[len(left.keys)-1]
	left.keys = left.keys[:len(left.keys)-1]
	left.children = left.children[:len(left.children)-1]
}

// borrowFromRight 从右侧兄弟节点借一个键给第 i 个孩子。
func (t *Tree[K, V]) borrowFromRight(n *node[K, V], i int) {
	child, right := n.children[i], n.children[i+1]
	if child.leaf {
		child.keys = append(child.keys, right.keys[0])
		child.values = append(child.values, right.values[0])
		right.keys = right.keys[1:]
		right.values = right.values[1:]
		n.keys[i] = right.keys[0]
		return
	}
	child.keys = append(child.keys, n.keys[i])
	child.children = append(child.children, right.children[0])
	n.keys[i] = right.keys[0]
	right.keys = right.keys[1:]
	right.children = right.children[1:]
}

// merge 合并第 i 和第 i+1 个孩子，并删除它们之间的分隔键。
func (t *Tree[K, V]) merge(n *node[K, V], i int) {
	left, right := n.children[i], n.children[i+1]
	if left.leaf {
		left.keys = append(left.keys, right.keys...)
		left.values = append(left.values, right.values...)
		left.next = right.next
	} else {
		left.keys = append(left.keys, n.keys[i])
		left.keys = append(left.keys, right.keys...)
		left.children = append(left.children, right.children...)
	}
	n.keys = append(n.keys[:i], n.keys[i+1:]...)
	n.children = append(n.children[:i+1], n.children[i+2:]...)
}

// lowerBound 返回第一个不小于 key 的下标。
func lowerBound[K Ordered](keys []K, key K) int {
	return sort.Search(len(keys), func(i int) bool { return !(keys[i] < key) })
}

// childIndex 返回 key 应落入的孩子下标，即第一个大于 key 的键的位置。
func childIndex[K Ordered](keys []K, key K) int {
	return sort.Search(len(keys), func(i int) bool { return key < keys[i] })
}

// insertAt 在下标 i 处插入元素。
func insertAt[T any](s []T, i int, v T) []T {
	s = append(s, v)
	copy(s[i+1:], s[i:])
	s[i] = v
	return s
}
