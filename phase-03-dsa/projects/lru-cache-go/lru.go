// Package lru provides a thread-safe, generic LRU cache.
//
// Design:
//   - A doubly-linked list maintains recency: front = most recent, back = least.
//   - A hash map maps keys to list nodes for O(1) lookup.
//   - On Get: move node to front. O(1).
//   - On Put: if key exists, update value + move to front. Else insert at front;
//     if over capacity, evict the back node.
//
// Time complexity: O(1) for Get and Put.
// Space: O(capacity).
package lru

import "sync"

// node is a doubly-linked list node holding a key/value pair.
// Lowercase = unexported. Encapsulation: callers can't touch our internals.
type node[K comparable, V any] struct {
	key        K
	value      V
	prev, next *node[K, V]
}

// Cache is the public type. Generic over K (comparable) and V (any).
type Cache[K comparable, V any] struct {
	mu       sync.Mutex
	capacity int
	items    map[K]*node[K, V]
	head     *node[K, V] // sentinel; head.next = most recent
	tail     *node[K, V] // sentinel; tail.prev = least recent
}

// New returns an empty Cache with the given capacity.
// Capacity must be positive.
func New[K comparable, V any](capacity int) *Cache[K, V] {
	if capacity <= 0 {
		panic("lru: capacity must be > 0")
	}
	c := &Cache[K, V]{
		capacity: capacity,
		items:    make(map[K]*node[K, V], capacity),
		head:     &node[K, V]{},
		tail:     &node[K, V]{},
	}
	// Sentinel nodes simplify edge cases — we never need to nil-check
	// "is this the first/last node" inside operations.
	c.head.next = c.tail
	c.tail.prev = c.head
	return c
}

// Len returns the current number of items.
func (c *Cache[K, V]) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.items)
}

// Get returns the value for key and true if present.
// On hit, the entry is marked most-recent.
func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	n, ok := c.items[key]
	if !ok {
		var zero V
		return zero, false
	}
	c.moveToFront(n)
	return n.value, true
}

// Put inserts or updates the value for key.
// If the cache is at capacity, the least-recent entry is evicted.
func (c *Cache[K, V]) Put(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if n, ok := c.items[key]; ok {
		n.value = value
		c.moveToFront(n)
		return
	}

	n := &node[K, V]{key: key, value: value}
	c.addToFront(n)
	c.items[key] = n

	if len(c.items) > c.capacity {
		// Evict the LRU node (the one before tail).
		evict := c.tail.prev
		c.unlink(evict)
		delete(c.items, evict.key)
	}
}

// Delete removes a key. Returns true if it was present.
func (c *Cache[K, V]) Delete(key K) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	n, ok := c.items[key]
	if !ok {
		return false
	}
	c.unlink(n)
	delete(c.items, key)
	return true
}

// === internal list helpers ===

// addToFront inserts n right after head.
func (c *Cache[K, V]) addToFront(n *node[K, V]) {
	n.prev = c.head
	n.next = c.head.next
	c.head.next.prev = n
	c.head.next = n
}

// unlink removes n from the list.
func (c *Cache[K, V]) unlink(n *node[K, V]) {
	n.prev.next = n.next
	n.next.prev = n.prev
	n.prev = nil
	n.next = nil
}

// moveToFront unlinks n and re-inserts at front.
func (c *Cache[K, V]) moveToFront(n *node[K, V]) {
	c.unlink(n)
	c.addToFront(n)
}
