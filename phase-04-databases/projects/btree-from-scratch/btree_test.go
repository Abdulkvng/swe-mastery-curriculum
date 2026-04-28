package btree

import (
	"math/rand"
	"sort"
	"testing"
)

func intLess(a, b int) bool { return a < b }

func TestInsertSearch(t *testing.T) {
	bt := New[int, string](3, intLess)
	keys := []int{10, 20, 5, 15, 30, 25, 1, 7, 12, 17}
	for _, k := range keys {
		bt.Insert(k, "v")
	}
	if got := bt.Len(); got != len(keys) {
		t.Errorf("Len = %d, want %d", got, len(keys))
	}
	for _, k := range keys {
		if _, ok := bt.Search(k); !ok {
			t.Errorf("missing key %d", k)
		}
	}
	if _, ok := bt.Search(999); ok {
		t.Error("found key that shouldn't exist")
	}
}

func TestUpdate(t *testing.T) {
	bt := New[int, string](2, intLess)
	bt.Insert(1, "a")
	bt.Insert(1, "b")
	if bt.Len() != 1 {
		t.Errorf("Len = %d, want 1", bt.Len())
	}
	if v, _ := bt.Search(1); v != "b" {
		t.Errorf("got %v, want b", v)
	}
}

func TestRange(t *testing.T) {
	bt := New[int, int](3, intLess)
	for _, k := range []int{1, 5, 10, 15, 20, 25, 30, 35, 40} {
		bt.Insert(k, k*10)
	}
	out := bt.Range(10, 30)
	if len(out) != 5 {
		t.Errorf("got %d, want 5: %v", len(out), out)
	}
	for i := 1; i < len(out); i++ {
		if out[i-1].key > out[i].key {
			t.Error("range not sorted")
		}
	}
}

func TestStress(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	bt := New[int, int](4, intLess)
	want := map[int]int{}
	for i := 0; i < 5000; i++ {
		k := rng.Intn(2000)
		v := rng.Int()
		bt.Insert(k, v)
		want[k] = v
	}
	if bt.Len() != len(want) {
		t.Errorf("Len = %d, want %d", bt.Len(), len(want))
	}
	for k, v := range want {
		got, ok := bt.Search(k)
		if !ok || got != v {
			t.Errorf("k=%d: got (%d,%v), want (%d,true)", k, got, ok, v)
		}
	}

	// Range scan all → should be sorted.
	keys := make([]int, 0, len(want))
	for k := range want {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	out := bt.Range(0, 1<<30)
	if len(out) != len(keys) {
		t.Errorf("range len %d, want %d", len(out), len(keys))
	}
	for i := range out {
		if out[i].key != keys[i] {
			t.Errorf("idx %d: got %d, want %d", i, out[i].key, keys[i])
			break
		}
	}
}
