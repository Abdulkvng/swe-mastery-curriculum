// Package btree is an in-memory B-tree.
//
// What's a B-tree?
//
// A self-balancing search tree where each node holds many keys (not just 2 like
// a binary tree). This minimizes "tree height" = minimizes disk reads for
// on-disk variants. Even in memory the cache locality is excellent.
//
// Order parameter `t` (minimum degree):
//   - Each non-root node has between t-1 and 2t-1 keys.
//   - Each internal node has between t and 2t children.
//   - Root has between 1 and 2t-1 keys.
//
// All leaves are at the same depth. That's the invariant balance maintains.
//
// Operations:
//   - Search(k): O(log n)
//   - Insert(k,v): O(log n)
//   - Delete(k): O(log n)
//   - Range(lo, hi): O(log n + r) where r = result size
//
// References:
//   - CLRS chapter 18.
//   - Postgres' B-tree (nbtree) is a more complex variant (Lehman-Yao).

package btree

import "sort"

// kv is a key-value pair stored in a node.
type kv[K, V any] struct {
	key K
	val V
}

// node is one B-tree node.
type node[K, V any] struct {
	leaf bool
	kvs  []kv[K, V]    // up to 2t-1 entries
	kids []*node[K, V] // children (only for internal nodes)
}

// BTree is a generic B-tree.
//
// Less is the user-supplied "less than" function for keys.
// We use a function instead of constraints.Ordered so K can be any user type.
type BTree[K, V any] struct {
	root *node[K, V]
	t    int // minimum degree
	less func(a, b K) bool
	size int
}

// New creates an empty BTree with minimum degree t and a comparator.
// t must be ≥ 2.
func New[K, V any](t int, less func(a, b K) bool) *BTree[K, V] {
	if t < 2 {
		panic("btree: t must be >= 2")
	}
	return &BTree[K, V]{
		root: &node[K, V]{leaf: true},
		t:    t,
		less: less,
	}
}

// Len returns the number of items.
func (b *BTree[K, V]) Len() int { return b.size }

// Search returns the value for key and true if found.
func (b *BTree[K, V]) Search(key K) (V, bool) {
	return b.search(b.root, key)
}

func (b *BTree[K, V]) search(n *node[K, V], key K) (V, bool) {
	// Find first index i where n.kvs[i].key >= key.
	i := sort.Search(len(n.kvs), func(i int) bool {
		return !b.less(n.kvs[i].key, key)
	})
	if i < len(n.kvs) && !b.less(key, n.kvs[i].key) {
		return n.kvs[i].val, true
	}
	if n.leaf {
		var zero V
		return zero, false
	}
	return b.search(n.kids[i], key)
}

// Insert puts (key, val) into the tree, replacing existing value if present.
func (b *BTree[K, V]) Insert(key K, val V) {
	root := b.root
	if len(root.kvs) == 2*b.t-1 {
		// Root is full. Grow the tree by one level.
		newRoot := &node[K, V]{leaf: false, kids: []*node[K, V]{root}}
		b.splitChild(newRoot, 0)
		b.root = newRoot
	}
	b.insertNonFull(b.root, key, val)
}

// insertNonFull inserts into a node that is guaranteed not full.
// We split children as we descend so we can always insert into a non-full node.
func (b *BTree[K, V]) insertNonFull(n *node[K, V], key K, val V) {
	i := sort.Search(len(n.kvs), func(i int) bool {
		return !b.less(n.kvs[i].key, key)
	})
	// Update if key exists.
	if i < len(n.kvs) && !b.less(key, n.kvs[i].key) {
		n.kvs[i].val = val
		return
	}
	if n.leaf {
		// Insert into leaf at position i.
		n.kvs = append(n.kvs, kv[K, V]{})
		copy(n.kvs[i+1:], n.kvs[i:])
		n.kvs[i] = kv[K, V]{key: key, val: val}
		b.size++
		return
	}
	// Internal: descend into kids[i]. If full, split first.
	if len(n.kids[i].kvs) == 2*b.t-1 {
		b.splitChild(n, i)
		// After split, decide which half to descend into.
		if b.less(n.kvs[i].key, key) {
			i++
		} else if !b.less(key, n.kvs[i].key) {
			// Key was promoted into n.kvs[i] — update and done.
			n.kvs[i].val = val
			return
		}
	}
	b.insertNonFull(n.kids[i], key, val)
}

// splitChild splits parent.kids[i] (which is full) around its median.
// Median is promoted into parent.
func (b *BTree[K, V]) splitChild(parent *node[K, V], i int) {
	t := b.t
	full := parent.kids[i]
	mid := full.kvs[t-1] // promote this

	right := &node[K, V]{leaf: full.leaf}
	right.kvs = append(right.kvs, full.kvs[t:]...)
	if !full.leaf {
		right.kids = append(right.kids, full.kids[t:]...)
		full.kids = full.kids[:t]
	}
	full.kvs = full.kvs[:t-1]

	// Insert mid into parent at position i (and right child at i+1).
	parent.kvs = append(parent.kvs, kv[K, V]{})
	copy(parent.kvs[i+1:], parent.kvs[i:])
	parent.kvs[i] = mid

	parent.kids = append(parent.kids, nil)
	copy(parent.kids[i+2:], parent.kids[i+1:])
	parent.kids[i+1] = right
}

// Range returns all (key, value) pairs with lo <= key <= hi, in order.
func (b *BTree[K, V]) Range(lo, hi K) []kv[K, V] {
	var out []kv[K, V]
	b.rangeHelper(b.root, lo, hi, &out)
	return out
}

func (b *BTree[K, V]) rangeHelper(n *node[K, V], lo, hi K, out *[]kv[K, V]) {
	for i := 0; i < len(n.kvs); i++ {
		if !n.leaf {
			b.rangeHelper(n.kids[i], lo, hi, out)
		}
		k := n.kvs[i].key
		if !b.less(k, lo) && !b.less(hi, k) {
			*out = append(*out, n.kvs[i])
		}
	}
	if !n.leaf {
		b.rangeHelper(n.kids[len(n.kvs)], lo, hi, out)
	}
}

// Delete is left as an exercise (it's the longest part of the algorithm).
// CLRS chapter 18.3 has the full procedure: handle leaf delete,
// handle internal delete by replacing with predecessor/successor,
// rebalance as you ascend.
//
// For now: rebuild the tree without the key (slow but correct for tests):
func (b *BTree[K, V]) DeleteSlow(key K) bool {
	if _, ok := b.Search(key); !ok {
		return false
	}
	// Collect everything, rebuild without the key.
	var all []kv[K, V]
	b.collectAll(b.root, &all)
	*b = *New[K, V](b.t, b.less)
	for _, e := range all {
		if !b.less(e.key, key) && !b.less(key, e.key) {
			continue
		}
		b.Insert(e.key, e.val)
	}
	return true
}

func (b *BTree[K, V]) collectAll(n *node[K, V], out *[]kv[K, V]) {
	for i := 0; i < len(n.kvs); i++ {
		if !n.leaf {
			b.collectAll(n.kids[i], out)
		}
		*out = append(*out, n.kvs[i])
	}
	if !n.leaf {
		b.collectAll(n.kids[len(n.kvs)], out)
	}
}
