package lru

import (
	"strconv"
	"sync"
	"testing"
)

func TestBasic(t *testing.T) {
	c := New[string, int](2)

	c.Put("a", 1)
	c.Put("b", 2)
	if v, ok := c.Get("a"); !ok || v != 1 {
		t.Errorf("Get a: got (%v, %v), want (1, true)", v, ok)
	}
	c.Put("c", 3)                  // should evict "b" (lru since "a" was just touched)
	if _, ok := c.Get("b"); ok {
		t.Error("expected b to be evicted")
	}
	if v, ok := c.Get("c"); !ok || v != 3 {
		t.Errorf("Get c: got (%v, %v), want (3, true)", v, ok)
	}
}

func TestUpdateExisting(t *testing.T) {
	c := New[string, int](2)
	c.Put("a", 1)
	c.Put("b", 2)
	c.Put("a", 99) // update + bump to front
	c.Put("c", 3)  // should evict "b", not "a"
	if v, ok := c.Get("a"); !ok || v != 99 {
		t.Errorf("got (%v, %v), want (99, true)", v, ok)
	}
	if _, ok := c.Get("b"); ok {
		t.Error("expected b to be evicted")
	}
}

func TestDelete(t *testing.T) {
	c := New[string, int](3)
	c.Put("a", 1); c.Put("b", 2); c.Put("c", 3)
	if !c.Delete("b") {
		t.Error("Delete b returned false")
	}
	if c.Len() != 2 {
		t.Errorf("Len = %d, want 2", c.Len())
	}
	if _, ok := c.Get("b"); ok {
		t.Error("b still present after Delete")
	}
}

func TestZeroCapacityPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic on capacity 0")
		}
	}()
	New[string, int](0)
}

func TestConcurrent(t *testing.T) {
	c := New[int, int](100)
	var wg sync.WaitGroup
	for g := 0; g < 50; g++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for i := 0; i < 1000; i++ {
				k := (id*1000 + i) % 200
				c.Put(k, k*10)
				c.Get(k)
			}
		}(g)
	}
	wg.Wait()
	// No specific assertion — just that we didn't panic / race.
}

func BenchmarkPut(b *testing.B) {
	c := New[string, int](10000)
	for i := 0; i < b.N; i++ {
		c.Put(strconv.Itoa(i), i)
	}
}

func BenchmarkGetHit(b *testing.B) {
	c := New[string, int](10000)
	for i := 0; i < 10000; i++ {
		c.Put(strconv.Itoa(i), i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Get(strconv.Itoa(i % 10000))
	}
}
