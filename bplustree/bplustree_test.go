package bplustree

import (
	"fmt"
	"math/rand"
	"testing"
)

// validate 检查 B+ 树的各项结构约束是否仍然成立。
func validate[K Ordered, V any](tb testing.TB, t *Tree[K, V]) {
	tb.Helper()

	var leafDepth = -1
	var prevKey K
	first := true
	count := 0

	var walk func(n *node[K, V], depth int)
	walk = func(n *node[K, V], depth int) {
		for i := 1; i < len(n.keys); i++ {
			if !(n.keys[i-1] < n.keys[i]) {
				tb.Fatalf("keys not strictly increasing: %v", n.keys)
			}
		}
		if n != t.root && len(n.keys) < t.minKeys {
			tb.Fatalf("node underflow: %d keys < %d", len(n.keys), t.minKeys)
		}
		if len(n.keys) > t.maxKeys {
			tb.Fatalf("node overflow: %d keys > %d", len(n.keys), t.maxKeys)
		}

		if n.leaf {
			if leafDepth == -1 {
				leafDepth = depth
			} else if leafDepth != depth {
				tb.Fatalf("leaves at different depths: %d and %d", leafDepth, depth)
			}
			for i := range n.keys {
				if !first && !(prevKey < n.keys[i]) {
					tb.Fatalf("leaf chain out of order: %v then %v", prevKey, n.keys[i])
				}
				prevKey = n.keys[i]
				first = false
				count++
			}
			return
		}

		if len(n.children) != len(n.keys)+1 {
			tb.Fatalf("internal node has %d keys but %d children", len(n.keys), len(n.children))
		}
		for _, c := range n.children {
			walk(c, depth+1)
		}
	}

	walk(t.root, 0)
	if count != t.size {
		tb.Fatalf("size mismatch: counted %d, tree says %d", count, t.size)
	}
}

func TestPutGet(t *testing.T) {
	tr := New[int, string](4)
	for i := 0; i < 100; i++ {
		if !tr.Put(i, fmt.Sprint(i)) {
			t.Fatalf("put %d should report created", i)
		}
	}
	if tr.Len() != 100 {
		t.Fatalf("want len 100, got %d", tr.Len())
	}
	for i := 0; i < 100; i++ {
		if v, ok := tr.Get(i); !ok || v != fmt.Sprint(i) {
			t.Fatalf("get %d = %q, %v", i, v, ok)
		}
	}
	if _, ok := tr.Get(1000); ok {
		t.Fatal("get on missing key should report false")
	}
	if tr.Put(5, "five") {
		t.Fatal("put on existing key should report not created")
	}
	if v, _ := tr.Get(5); v != "five" {
		t.Fatalf("value should be updated, got %q", v)
	}
	validate(t, tr)
}

func TestDelete(t *testing.T) {
	tr := New[int, int](3)
	for i := 0; i < 200; i++ {
		tr.Put(i, i*10)
	}
	if !tr.Delete(0) || tr.Delete(0) {
		t.Fatal("delete should succeed once and fail afterwards")
	}
	for i := 2; i < 200; i += 2 {
		if !tr.Delete(i) {
			t.Fatalf("delete %d failed", i)
		}
	}
	validate(t, tr)
	for i := 0; i < 200; i++ {
		v, ok := tr.Get(i)
		if i%2 == 0 {
			if ok {
				t.Fatalf("key %d should be gone", i)
			}
			continue
		}
		if !ok || v != i*10 {
			t.Fatalf("key %d should remain, got %v %v", i, v, ok)
		}
	}

	for i := 0; i < 200; i++ {
		tr.Delete(i)
		validate(t, tr)
	}
	if !tr.Empty() || tr.Len() != 0 {
		t.Fatalf("tree should be empty, got len %d", tr.Len())
	}
	if tr.Height() != 1 {
		t.Fatalf("empty tree height should be 1, got %d", tr.Height())
	}
}

func TestAscendAndRange(t *testing.T) {
	tr := New[int, int](5)
	for _, k := range []int{9, 3, 7, 1, 5, 8, 2, 6, 4, 0} {
		tr.Put(k, k*k)
	}
	var got []int
	tr.Ascend(func(k, v int) bool {
		got = append(got, k)
		return true
	})
	for i, k := range got {
		if k != i {
			t.Fatalf("ascend order wrong at %d: %v", i, got)
		}
	}

	got = nil
	tr.Range(3, 7, func(k, v int) bool {
		got = append(got, k)
		return true
	})
	want := []int{3, 4, 5, 6}
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("range want %v, got %v", want, got)
	}

	// fn 返回 false 提前终止。
	calls := 0
	tr.Ascend(func(k, v int) bool {
		calls++
		return calls < 3
	})
	if calls != 3 {
		t.Fatalf("early stop want 3 calls, got %d", calls)
	}
}

func TestRandomAgainstMap(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	tr := New[int, int](4)
	ref := map[int]int{}

	for step := 0; step < 5000; step++ {
		k := rng.Intn(500)
		if rng.Intn(3) == 0 {
			_, want := ref[k]
			if got := tr.Delete(k); got != want {
				t.Fatalf("step %d: delete %d = %v, want %v", step, k, got, want)
			}
			delete(ref, k)
		} else {
			_, want := ref[k]
			if got := tr.Put(k, k*3); got == want {
				t.Fatalf("step %d: put %d created=%v, want %v", step, k, got, !want)
			}
			ref[k] = k * 3
		}
		if step%97 == 0 {
			validate(t, tr)
			if tr.Len() != len(ref) {
				t.Fatalf("step %d: len %d, want %d", step, tr.Len(), len(ref))
			}
		}
	}

	validate(t, tr)
	tr.Ascend(func(k, v int) bool {
		if ref[k] != v {
			t.Fatalf("key %d value mismatch", k)
		}
		delete(ref, k)
		return true
	})
	if len(ref) != 0 {
		t.Fatalf("%d keys missing from tree", len(ref))
	}
}

func TestStringKeys(t *testing.T) {
	tr := New[string, int](3)
	words := []string{"go", "b+tree", "算法", "data", "structure"}
	for i, w := range words {
		tr.Put(w, i)
	}
	var got []string
	tr.Ascend(func(k string, v int) bool {
		got = append(got, k)
		return true
	})
	for i := 1; i < len(got); i++ {
		if !(got[i-1] < got[i]) {
			t.Fatalf("strings not sorted: %v", got)
		}
	}
	validate(t, tr)
}

func ExampleTree() {
	tr := New[int, string](4)
	for _, k := range []int{5, 2, 8, 1, 9, 3} {
		tr.Put(k, fmt.Sprint("v", k))
	}
	tr.Delete(8)
	tr.Range(2, 6, func(k int, v string) bool {
		fmt.Println(k, v)
		return true
	})
	// Output:
	// 2 v2
	// 3 v3
	// 5 v5
}
