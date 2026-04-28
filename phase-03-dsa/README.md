# Phase 3 — Data Structures & Algorithms

> Every interview at Apple, Datadog, and frankly anywhere will hit DSA. But DSA isn't just an interview tax — it's how you think about *cost*. When a senior engineer reads a code review and says "this is O(n²), can we make it O(n log n)?" — they're not flexing. They've internalized that some operations are 1000x cheaper than others, and they care.
>
> This phase isn't 800 LeetCode problems. It's the **right** ~80, organized by data structure, with full Go implementations of every primitive, full Rust implementations of two key data structures (LRU cache, B-tree node), and patterns mapped to Apple's actual interview question bank.

**Time:** 8–14 days, plus continuous practice afterward.

**You'll know you're done when:** you can recognize a problem's pattern in 30 seconds, code the canonical solution in 15 minutes, and articulate the time/space complexity without thinking.

---

## Table of contents

1. [What does this even mean? — "Algorithms"](#what-algorithms-mean)
2. [Module 3.1 — Big-O notation, properly](#module-31--big-o)
3. [Module 3.2 — Arrays and strings](#module-32--arrays)
4. [Module 3.3 — Hash maps and hash sets](#module-33--hashmaps)
5. [Module 3.4 — Linked lists](#module-34--linked-lists)
6. [Module 3.5 — Stacks and queues](#module-35--stacks-queues)
7. [Module 3.6 — Trees: BST, AVL, B-tree](#module-36--trees)
8. [Module 3.7 — Heaps and priority queues](#module-37--heaps)
9. [Module 3.8 — Tries](#module-38--tries)
10. [Module 3.9 — Graphs: BFS, DFS, Dijkstra, topological sort](#module-39--graphs)
11. [Module 3.10 — Two pointers, sliding window, prefix sums](#module-310--patterns)
12. [Module 3.11 — Recursion and backtracking](#module-311--recursion)
13. [Module 3.12 — Dynamic programming](#module-312--dp)
14. [Module 3.13 — Bit manipulation](#module-313--bits)
15. [🛠️ Project: LRU cache in Rust + Go](#project-lru)
16. [The 80 problems (mapped)](#the-80)
17. [Apple-style interview tactics](#apple-tactics)
18. [What you should now know](#what-you-should-now-know)

---

<a name="what-algorithms-mean"></a>
## 🧠 What does this even mean? — "Algorithms"

An **algorithm** is just a recipe — a finite sequence of steps to solve a problem.

A **data structure** is a way of organizing data so that certain operations are cheap.

The interplay matters: the right data structure makes hard algorithms easy. Searching an unsorted list = O(n). Searching a sorted array = O(log n). Searching a hash map = O(1) on average. Same problem, three different costs depending on how you stored the data.

**Why this is harder for self-taught engineers:** you can ship a lot of code without ever consciously choosing a data structure. Then a system slows down at scale and you have no framework for "why." DSA gives you that framework.

---

<a name="module-31--big-o"></a>
## Module 3.1 — Big-O notation, properly

> 📖 **Definition — Big-O:** A way to describe how an algorithm's **runtime or memory** grows as the input size grows. It's about *trends*, not stopwatch timings.

### The cheat sheet

| Class | Name | Example |
|---|---|---|
| O(1) | Constant | hash map lookup, array index |
| O(log n) | Logarithmic | binary search, BST lookup (balanced) |
| O(n) | Linear | iterate an array, linked list search |
| O(n log n) | Linearithmic | merge sort, quicksort (avg), heapsort |
| O(n²) | Quadratic | nested loops, bubble sort |
| O(n³) | Cubic | naive matrix multiply |
| O(2ⁿ) | Exponential | naive recursive Fibonacci, brute-force subset sum |
| O(n!) | Factorial | brute-force traveling salesman, permutations |

### Read it like this

When someone says "this is O(n)," they mean: as the input gets 10x bigger, the time gets ~10x longer. As it gets 1000x bigger, time gets ~1000x longer.

When someone says O(n log n), they mean: 10x bigger input → ~33x longer (for n around 1000). Looks similar to linear at small n but grows slower than n².

When someone says O(n²): 10x bigger → 100x longer. **This is where things start to hurt at scale.**

For n = 1,000,000:
- O(n) = 1M operations, microseconds
- O(n log n) = ~20M operations, fractions of a second
- O(n²) = 1 trillion operations, minutes to hours
- O(2ⁿ) = literally beyond the heat death of the universe

### Big-O of common operations

| Data structure | Access | Search | Insert | Delete |
|---|---|---|---|---|
| Array | O(1) | O(n) | O(n) | O(n) |
| Sorted array | O(1) | O(log n) | O(n) | O(n) |
| Linked list | O(n) | O(n) | O(1)* | O(1)* |
| Hash map | — | O(1) avg, O(n) worst | O(1) avg | O(1) avg |
| BST (balanced) | — | O(log n) | O(log n) | O(log n) |
| BST (unbalanced) | — | O(n) | O(n) | O(n) |
| Heap | O(1) min/max | O(n) | O(log n) | O(log n) |
| Trie | — | O(L) | O(L) | O(L) |

*at a known position

### Space complexity

Same notation for memory. An algorithm can be O(n) time but O(1) space (in-place reverse) or O(n) time AND O(n) space (recursive solution allocating a stack).

> 🎯 **Interview tip:** Always state both time AND space complexity. "This is O(n) time, O(1) space — we modify in place."

### Amortized vs worst case

> 📖 **Definition — Amortized:** Average cost over many operations, even if any one is expensive. Example: appending to a dynamic array. Most appends are O(1). Occasionally the array doubles (O(n) copy). Amortized: O(1).

If you don't say "amortized," interviewers assume worst-case. So either say "amortized O(1)" or be ready to defend the worst case.

---

<a name="module-32--arrays"></a>
## Module 3.2 — Arrays and strings

The most fundamental data structure. Contiguous memory, indexed access.

### In Go: slices

A Go slice is a struct: `(pointer, length, capacity)`. The underlying array is shared. Capacity ≥ length. Append doubles capacity when full (amortized O(1)).

```go
nums := []int{1, 2, 3, 4, 5}
fmt.Println(nums[2])           // 3
nums = append(nums, 6)         // amortized O(1)
sub := nums[1:4]               // view of nums[1..4), shares memory
sub[0] = 99                    // also changes nums[1]!
fmt.Println(nums)              // [1, 99, 3, 4, 5, 6]
```

> 🎯 The "shared memory" gotcha is interview gold. If you do `slice2 := slice1[1:3]`, mutations leak. To detach: `slice2 := append([]int(nil), slice1[1:3]...)`.

### Common patterns

#### Reverse in place — O(n) time, O(1) space

```go
func reverse(a []int) {
    i, j := 0, len(a)-1
    for i < j {
        a[i], a[j] = a[j], a[i]
        i++
        j--
    }
}
```

#### Find duplicates with a set

```go
func hasDup(a []int) bool {
    seen := map[int]struct{}{}
    for _, v := range a {
        if _, ok := seen[v]; ok {
            return true
        }
        seen[v] = struct{}{}
    }
    return false
}
```

`struct{}` is the zero-byte type — empty value when you only care about keys.

#### Two-sum — classic warm-up

```go
// Given nums and target, return indices of two numbers that sum to target.
func twoSum(nums []int, target int) [2]int {
    idx := map[int]int{} // value -> index
    for i, v := range nums {
        if j, ok := idx[target-v]; ok {
            return [2]int{j, i}
        }
        idx[v] = i
    }
    return [2]int{-1, -1}
}
```

O(n) time, O(n) space. The naive nested-loop is O(n²).

#### Strings are arrays of bytes (or runes)

```go
s := "héllo"
fmt.Println(len(s))         // 6 — BYTES, because é is 2 bytes in UTF-8
fmt.Println(len([]rune(s))) // 5 — runes (code points)
```

For ASCII you can index `s[i]` (returns a byte). For Unicode, convert to `[]rune` first or iterate with `for i, r := range s`.

---

<a name="module-33--hashmaps"></a>
## Module 3.3 — Hash maps and hash sets

> 📖 **Definition — Hash map:** A data structure that maps keys to values using a hash function. Internally it's an array of "buckets" — the hash function picks which bucket; collisions are handled by chaining (linked list per bucket) or open addressing.

Average ops: O(1). Worst case: O(n) if every key collides — but with a good hash function, this never happens in practice.

```go
m := map[string]int{}
m["foo"] = 42
v, ok := m["foo"]   // v=42, ok=true
delete(m, "foo")

// Iterate. Note: Go DELIBERATELY randomizes iteration order to prevent
// you from depending on it.
for k, v := range m {
    fmt.Println(k, v)
}
```

### When to reach for a hash map

Whenever the brute force is "for each item, scan the rest" — try a hash map. Examples:
- Two-sum
- Group anagrams
- Detect cycles (visited set)
- Count frequencies
- Check if a word exists in a dictionary

### Sets in Go (no built-in type)

```go
type Set[T comparable] struct {
    m map[T]struct{}
}
func NewSet[T comparable]() *Set[T] { return &Set[T]{m: map[T]struct{}{}} }
func (s *Set[T]) Add(v T)            { s.m[v] = struct{}{} }
func (s *Set[T]) Has(v T) bool       { _, ok := s.m[v]; return ok }
func (s *Set[T]) Remove(v T)         { delete(s.m, v) }
func (s *Set[T]) Len() int           { return len(s.m) }
```

`comparable` is a Go generic constraint meaning "anything that supports `==`."

---

<a name="module-34--linked-lists"></a>
## Module 3.4 — Linked lists

A sequence of nodes, each pointing to the next. No contiguous memory.

```go
type Node struct {
    Val  int
    Next *Node
}
```

Pros: O(1) insert/delete at known positions.
Cons: O(n) random access, bad cache locality.

In real software you almost never use a linked list — slices are faster for nearly everything because of CPU caches. **But interviews love them** because they teach pointer manipulation.

### Reverse a linked list — the classic

```go
func reverse(head *Node) *Node {
    var prev *Node
    curr := head
    for curr != nil {
        next := curr.Next   // save
        curr.Next = prev    // reverse the link
        prev = curr         // step
        curr = next
    }
    return prev             // new head
}
```

Visualize:
```
Initial:  1 -> 2 -> 3 -> nil
After:    nil <- 1 <- 2 <- 3   (head now points to 3)
```

### Detect cycle (Floyd's tortoise and hare)

```go
func hasCycle(head *Node) bool {
    slow, fast := head, head
    for fast != nil && fast.Next != nil {
        slow = slow.Next
        fast = fast.Next.Next
        if slow == fast {
            return true
        }
    }
    return false
}
```

Two pointers, fast moves twice as quickly. If there's a cycle, fast eventually laps slow.

### Find middle node

```go
func middle(head *Node) *Node {
    slow, fast := head, head
    for fast != nil && fast.Next != nil {
        slow = slow.Next
        fast = fast.Next.Next
    }
    return slow
}
```

When fast reaches the end, slow is at the middle. Fewer passes than length-then-half.

---

<a name="module-35--stacks-queues"></a>
## Module 3.5 — Stacks and queues

**Stack:** LIFO (last in, first out). Push, Pop, Peek. O(1) all.
**Queue:** FIFO (first in, first out). Enqueue, Dequeue. O(1) all.

In Go, both can be slices but a queue with slice prefix-removal is O(n). Use a `container/list` or a circular buffer for O(1) queue.

```go
// Stack on a slice — clean and idiomatic
type Stack[T any] struct{ s []T }
func (s *Stack[T]) Push(v T)        { s.s = append(s.s, v) }
func (s *Stack[T]) Pop() (T, bool) {
    if len(s.s) == 0 { var z T; return z, false }
    v := s.s[len(s.s)-1]
    s.s = s.s[:len(s.s)-1]
    return v, true
}
func (s *Stack[T]) Peek() (T, bool) {
    if len(s.s) == 0 { var z T; return z, false }
    return s.s[len(s.s)-1], true
}
func (s *Stack[T]) Len() int { return len(s.s) }
```

### Classic uses

- **Stack:** parsing balanced parens, evaluating expressions, undo, DFS, function call stack itself
- **Queue:** BFS, scheduling, request buffering

### Valid parentheses (interview classic)

```go
func validParens(s string) bool {
    stack := []byte{}
    pairs := map[byte]byte{')': '(', ']': '[', '}': '{'}
    for i := 0; i < len(s); i++ {
        c := s[i]
        switch c {
        case '(', '[', '{':
            stack = append(stack, c)
        case ')', ']', '}':
            if len(stack) == 0 || stack[len(stack)-1] != pairs[c] {
                return false
            }
            stack = stack[:len(stack)-1]
        }
    }
    return len(stack) == 0
}
```

---

<a name="module-36--trees"></a>
## Module 3.6 — Trees: BST, AVL, B-tree

### Binary tree basics

```go
type TreeNode struct {
    Val         int
    Left, Right *TreeNode
}
```

### Three traversal orders

For tree:
```
      1
     / \
    2   3
   / \
  4   5
```

```go
func inorder(n *TreeNode, out *[]int) {
    if n == nil { return }
    inorder(n.Left, out)
    *out = append(*out, n.Val)
    inorder(n.Right, out)
}
// inorder(root) → [4, 2, 5, 1, 3]
// preorder      → [1, 2, 4, 5, 3]   (root, left, right)
// postorder     → [4, 5, 2, 3, 1]   (left, right, root)
```

### Level-order (BFS) traversal

```go
func levelOrder(root *TreeNode) [][]int {
    if root == nil { return nil }
    var result [][]int
    queue := []*TreeNode{root}
    for len(queue) > 0 {
        size := len(queue)
        level := make([]int, 0, size)
        for i := 0; i < size; i++ {
            n := queue[0]
            queue = queue[1:]
            level = append(level, n.Val)
            if n.Left != nil  { queue = append(queue, n.Left) }
            if n.Right != nil { queue = append(queue, n.Right) }
        }
        result = append(result, level)
    }
    return result
}
```

### Binary Search Tree (BST)

> 📖 **Definition — BST:** A binary tree where every left descendant ≤ node ≤ every right descendant.

Search/insert/delete: O(log n) if balanced, O(n) if degenerate (sorted insert builds a linked list).

```go
func search(n *TreeNode, val int) *TreeNode {
    for n != nil {
        switch {
        case val == n.Val:
            return n
        case val < n.Val:
            n = n.Left
        default:
            n = n.Right
        }
    }
    return nil
}

func insert(n *TreeNode, val int) *TreeNode {
    if n == nil { return &TreeNode{Val: val} }
    if val < n.Val {
        n.Left = insert(n.Left, val)
    } else if val > n.Val {
        n.Right = insert(n.Right, val)
    }
    return n
}
```

### AVL trees and Red-Black trees (don't implement — understand)

Self-balancing BSTs. Insertions trigger rotations to keep depth ≈ log n. AVL is stricter (faster lookups, slower mutations); RB is looser (used by Linux kernel scheduler, Java's TreeMap, etc.).

You won't be asked to implement one in an interview, but you should know:
- Both guarantee O(log n) operations
- They differ by how strict the balance invariant is
- Java's `TreeMap` is RB; C++ STL's `map` is RB

### B-tree (you'll see this in databases)

Multi-way tree (each node has many children, not just 2). Designed for *disk* — minimize disk seeks by packing many keys per node. Used in:
- Postgres indexes (and most relational DBs)
- Filesystems (ext4, NTFS, HFS+)

We'll implement a B-tree node in Rust later in this phase — see `projects/btree-rust/`.

---

<a name="module-37--heaps"></a>
## Module 3.7 — Heaps and priority queues

> 📖 **Definition — Heap:** A complete binary tree where every parent ≥ children (max-heap) or ≤ children (min-heap). Stored in an array using the parent/child index trick.

Indexed array form: parent of `i` is `(i-1)/2`; children of `i` are `2i+1` and `2i+2`.

Operations:
- Peek (top): O(1)
- Insert: O(log n) — bubble up
- Pop: O(log n) — swap last with first, sift down

Go provides `container/heap`:

```go
import "container/heap"

type IntHeap []int
func (h IntHeap) Len() int            { return len(h) }
func (h IntHeap) Less(i, j int) bool  { return h[i] < h[j] }  // min-heap
func (h IntHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *IntHeap) Push(x any)         { *h = append(*h, x.(int)) }
func (h *IntHeap) Pop() any {
    old := *h
    n := len(old)
    x := old[n-1]
    *h = old[:n-1]
    return x
}

h := &IntHeap{2, 1, 5, 3}
heap.Init(h)
heap.Push(h, 0)
fmt.Println(heap.Pop(h))  // 0 (smallest)
```

### When to use a heap

- "Top K" problems (top K largest/smallest)
- Median of a stream (two heaps trick)
- Task scheduling / event simulation
- Dijkstra's algorithm
- Merging k sorted lists

---

<a name="module-38--tries"></a>
## Module 3.8 — Tries

> 📖 **Definition — Trie (prefix tree):** A tree where each node represents a character, and paths from root to nodes spell words. Lookups O(L) where L = key length, regardless of dictionary size.

```go
type Trie struct {
    children [26]*Trie  // a-z
    isEnd    bool
}

func (t *Trie) Insert(word string) {
    node := t
    for i := 0; i < len(word); i++ {
        idx := word[i] - 'a'
        if node.children[idx] == nil {
            node.children[idx] = &Trie{}
        }
        node = node.children[idx]
    }
    node.isEnd = true
}

func (t *Trie) Search(word string) bool {
    node := t.find(word)
    return node != nil && node.isEnd
}

func (t *Trie) StartsWith(prefix string) bool {
    return t.find(prefix) != nil
}

func (t *Trie) find(s string) *Trie {
    node := t
    for i := 0; i < len(s); i++ {
        idx := s[i] - 'a'
        if node.children[idx] == nil { return nil }
        node = node.children[idx]
    }
    return node
}
```

Used in: autocomplete, IDE intellisense, spell checkers, IP routing tables.

---

<a name="module-39--graphs"></a>
## Module 3.9 — Graphs: BFS, DFS, Dijkstra, topological sort

### Representations

```go
// Adjacency list — most common
type Graph struct {
    adj map[int][]int   // node -> neighbors
}

// Adjacency matrix — when graph is dense and small
matrix := [][]int{
    {0, 1, 0},
    {1, 0, 1},
    {0, 1, 0},
}
```

### BFS — Breadth-first search

Layer by layer. Uses a queue. Good for **shortest path in unweighted graphs.**

```go
func bfs(g *Graph, start int) []int {
    visited := map[int]bool{start: true}
    queue := []int{start}
    order := []int{}
    for len(queue) > 0 {
        node := queue[0]
        queue = queue[1:]
        order = append(order, node)
        for _, n := range g.adj[node] {
            if !visited[n] {
                visited[n] = true
                queue = append(queue, n)
            }
        }
    }
    return order
}
```

### DFS — Depth-first search

Uses recursion (or an explicit stack). Useful for: cycle detection, topological sort, finding connected components.

```go
func dfs(g *Graph, start int) []int {
    visited := map[int]bool{}
    var order []int
    var rec func(int)
    rec = func(n int) {
        if visited[n] { return }
        visited[n] = true
        order = append(order, n)
        for _, nb := range g.adj[n] {
            rec(nb)
        }
    }
    rec(start)
    return order
}
```

### Dijkstra's — shortest path with weights

Used for weighted graphs (no negative weights). Uses a min-heap.

```go
import "container/heap"

type Edge struct{ to, weight int }
type Item struct{ node, dist int }
type PQ []Item

func (p PQ) Len() int            { return len(p) }
func (p PQ) Less(i, j int) bool  { return p[i].dist < p[j].dist }
func (p PQ) Swap(i, j int)       { p[i], p[j] = p[j], p[i] }
func (p *PQ) Push(x any)         { *p = append(*p, x.(Item)) }
func (p *PQ) Pop() any           { old := *p; v := old[len(old)-1]; *p = old[:len(old)-1]; return v }

func dijkstra(graph map[int][]Edge, source int) map[int]int {
    dist := map[int]int{source: 0}
    pq := &PQ{{node: source, dist: 0}}
    heap.Init(pq)
    for pq.Len() > 0 {
        it := heap.Pop(pq).(Item)
        if it.dist > dist[it.node] { continue }    // stale entry
        for _, e := range graph[it.node] {
            nd := it.dist + e.weight
            if d, ok := dist[e.to]; !ok || nd < d {
                dist[e.to] = nd
                heap.Push(pq, Item{node: e.to, dist: nd})
            }
        }
    }
    return dist
}
```

### Topological sort

For DAGs (directed acyclic graphs). Order nodes so every edge goes from earlier to later. Used in build systems, course prerequisites, task scheduling.

Kahn's algorithm (BFS-based):

```go
func topoSort(numNodes int, edges [][2]int) []int {
    indeg := make([]int, numNodes)
    adj := make([][]int, numNodes)
    for _, e := range edges {
        adj[e[0]] = append(adj[e[0]], e[1])
        indeg[e[1]]++
    }
    queue := []int{}
    for i := 0; i < numNodes; i++ {
        if indeg[i] == 0 { queue = append(queue, i) }
    }
    var order []int
    for len(queue) > 0 {
        n := queue[0]; queue = queue[1:]
        order = append(order, n)
        for _, nb := range adj[n] {
            indeg[nb]--
            if indeg[nb] == 0 { queue = append(queue, nb) }
        }
    }
    if len(order) != numNodes {
        return nil   // cycle exists
    }
    return order
}
```

---

<a name="module-310--patterns"></a>
## Module 3.10 — Two pointers, sliding window, prefix sums

These are *patterns*, not data structures. Recognizing them turns hard problems into easy ones.

### Two pointers

Two indices walking through an array, often from opposite ends or at different speeds.

```go
// Is s a palindrome?
func isPalindrome(s string) bool {
    i, j := 0, len(s)-1
    for i < j {
        if s[i] != s[j] { return false }
        i++; j--
    }
    return true
}

// Three-sum (find triplets summing to 0)
func threeSum(nums []int) [][]int {
    sort.Ints(nums)
    var result [][]int
    for i := 0; i < len(nums)-2; i++ {
        if i > 0 && nums[i] == nums[i-1] { continue }
        l, r := i+1, len(nums)-1
        for l < r {
            sum := nums[i] + nums[l] + nums[r]
            switch {
            case sum < 0: l++
            case sum > 0: r--
            default:
                result = append(result, []int{nums[i], nums[l], nums[r]})
                for l < r && nums[l] == nums[l+1] { l++ }
                for l < r && nums[r] == nums[r-1] { r-- }
                l++; r--
            }
        }
    }
    return result
}
```

### Sliding window

Maintain a window `[l, r]` over an array, expand and contract based on conditions.

```go
// Longest substring without repeating characters
func lengthOfLongestSubstring(s string) int {
    seen := map[byte]int{}
    best, l := 0, 0
    for r := 0; r < len(s); r++ {
        if i, ok := seen[s[r]]; ok && i >= l {
            l = i + 1
        }
        seen[s[r]] = r
        if r-l+1 > best { best = r - l + 1 }
    }
    return best
}
```

### Prefix sums

Precompute `prefix[i] = sum(arr[0..i])`. Range sum becomes O(1): `sum(arr[i..j]) = prefix[j] - prefix[i-1]`.

```go
type NumArray struct{ pre []int }

func New(nums []int) *NumArray {
    pre := make([]int, len(nums)+1)
    for i, v := range nums { pre[i+1] = pre[i] + v }
    return &NumArray{pre: pre}
}

func (n *NumArray) SumRange(i, j int) int {
    return n.pre[j+1] - n.pre[i]
}
```

---

<a name="module-311--recursion"></a>
## Module 3.11 — Recursion and backtracking

Recursion: a function that calls itself.

Two ingredients: **base case** (when to stop) and **recursive case** (how to reduce to a smaller problem).

```go
func factorial(n int) int {
    if n <= 1 { return 1 }       // base case
    return n * factorial(n-1)    // recursive case
}
```

### Backtracking

Explore all possibilities, undo when a path fails. Generic skeleton:

```go
func solve(state State) {
    if isComplete(state) {
        record(state)
        return
    }
    for _, choice := range choices(state) {
        apply(choice, &state)
        solve(state)
        undo(choice, &state)   // backtrack
    }
}
```

#### Generate all permutations

```go
func permute(nums []int) [][]int {
    var result [][]int
    var rec func(curr []int, used []bool)
    rec = func(curr []int, used []bool) {
        if len(curr) == len(nums) {
            cp := append([]int(nil), curr...)
            result = append(result, cp)
            return
        }
        for i, v := range nums {
            if used[i] { continue }
            used[i] = true
            curr = append(curr, v)
            rec(curr, used)
            curr = curr[:len(curr)-1]   // undo
            used[i] = false             // undo
        }
    }
    rec(nil, make([]bool, len(nums)))
    return result
}
```

#### N-queens, Sudoku, word ladder

All have the same shape: try, recurse, undo.

---

<a name="module-312--dp"></a>
## Module 3.12 — Dynamic programming

> 📖 **Definition — DP:** Solve a problem by combining solutions to overlapping subproblems, caching results so we don't redo work.

Two flavors:
- **Top-down (memoized recursion):** write the recursion, cache the answer.
- **Bottom-up (tabulation):** fill a table iteratively.

### Fibonacci — the canonical example

Naive recursion: O(2ⁿ). Memoized: O(n).

```go
// Top-down with memo
func fib(n int) int {
    memo := make(map[int]int)
    var rec func(int) int
    rec = func(k int) int {
        if k < 2 { return k }
        if v, ok := memo[k]; ok { return v }
        v := rec(k-1) + rec(k-2)
        memo[k] = v
        return v
    }
    return rec(n)
}

// Bottom-up, O(1) space
func fibBU(n int) int {
    if n < 2 { return n }
    a, b := 0, 1
    for i := 2; i <= n; i++ {
        a, b = b, a+b
    }
    return b
}
```

### Coin change — minimum coins to make `amount`

```go
func coinChange(coins []int, amount int) int {
    dp := make([]int, amount+1)
    for i := range dp { dp[i] = amount + 1 }
    dp[0] = 0
    for i := 1; i <= amount; i++ {
        for _, c := range coins {
            if i-c >= 0 && dp[i-c]+1 < dp[i] {
                dp[i] = dp[i-c] + 1
            }
        }
    }
    if dp[amount] > amount { return -1 }
    return dp[amount]
}
```

### Longest Common Subsequence

```go
func lcs(a, b string) int {
    m, n := len(a), len(b)
    dp := make([][]int, m+1)
    for i := range dp { dp[i] = make([]int, n+1) }
    for i := 1; i <= m; i++ {
        for j := 1; j <= n; j++ {
            if a[i-1] == b[j-1] {
                dp[i][j] = dp[i-1][j-1] + 1
            } else {
                dp[i][j] = max(dp[i-1][j], dp[i][j-1])
            }
        }
    }
    return dp[m][n]
}
```

### How to recognize a DP problem

Three signs:
1. "Find the optimal/min/max/count of X" (rather than "exact path")
2. Choices at each step
3. Future choices depend only on state, not history

---

<a name="module-313--bits"></a>
## Module 3.13 — Bit manipulation

Useful for: low-level optimization, certain interview questions, embedded work, networking (subnet masks).

### The basics

```go
// Bitwise operators
a & b    // AND
a | b    // OR
a ^ b    // XOR
^a       // NOT
a << n   // shift left (multiply by 2^n)
a >> n   // shift right (divide by 2^n)
```

### Common tricks

```go
// Check if i-th bit is set
isSet := (n >> i) & 1 == 1

// Set i-th bit
n |= (1 << i)

// Clear i-th bit
n &^= (1 << i)         // Go's specific "and-not" operator

// Toggle i-th bit
n ^= (1 << i)

// Is power of 2? (only one bit set)
isPow2 := n > 0 && n&(n-1) == 0

// Count set bits (Brian Kernighan)
func popcount(n uint) int {
    c := 0
    for n != 0 {
        n &= n - 1   // clears lowest set bit
        c++
    }
    return c
}

// Find single number where all others appear twice (XOR magic)
func singleNumber(nums []int) int {
    result := 0
    for _, n := range nums { result ^= n }   // pairs cancel; single survives
    return result
}
```

---

<a name="project-lru"></a>
## 🛠️ Project: LRU Cache in Go AND Rust

> 📖 **Definition — LRU Cache:** Least-Recently-Used cache. Fixed capacity. When full, evict whichever entry hasn't been touched in the longest. Used everywhere — CPU caches, web browser caches, DB page caches, Redis with LRU eviction.

The classic interview question: "design an LRU cache with O(1) get and put."

The trick: **doubly-linked list + hash map**. List maintains recency order; map gives O(1) lookup of node by key.

See `projects/lru-cache-go/` and `projects/lru-cache-rust/` for full implementations.

### Why both languages?

- **Go version**: clean and short. Demonstrates pointer manipulation in a friendly language.
- **Rust version**: forces you to think about ownership. Doubly-linked lists in safe Rust are notoriously tricky — you'll learn `Rc<RefCell<T>>` (shared mutable state). This is the kind of exercise that levels you up.

---

<a name="the-80"></a>
## The 80 problems (mapped)

Bookmark [LeetCode](https://leetcode.com). Solve in order. Each problem has its pattern listed.

### Arrays & strings (10)
1. Two Sum (hash map)
2. Best Time to Buy and Sell Stock (one-pass)
3. Contains Duplicate (hash set)
4. Product of Array Except Self (prefix products)
5. Maximum Subarray (Kadane / DP)
6. Merge Intervals (sort + scan)
7. 3Sum (two pointers)
8. Group Anagrams (hash + sort key)
9. Longest Substring Without Repeating Characters (sliding window)
10. Longest Palindromic Substring (expand from center)

### Linked lists (5)
11. Reverse Linked List
12. Linked List Cycle
13. Merge Two Sorted Lists
14. Remove Nth Node From End (two pointers)
15. Add Two Numbers (carry)

### Stacks & queues (5)
16. Valid Parentheses
17. Min Stack
18. Daily Temperatures (monotonic stack)
19. Implement Queue using Stacks
20. Largest Rectangle in Histogram (monotonic stack)

### Trees (10)
21. Maximum Depth of Binary Tree
22. Validate Binary Search Tree
23. Invert Binary Tree
24. Symmetric Tree
25. Binary Tree Level Order Traversal (BFS)
26. Binary Tree Right Side View
27. Lowest Common Ancestor of a BST
28. Lowest Common Ancestor of a Binary Tree
29. Serialize and Deserialize Binary Tree
30. Kth Smallest Element in a BST (in-order)

### Heaps (5)
31. Kth Largest Element in an Array
32. Top K Frequent Elements
33. Find Median from Data Stream (two heaps)
34. Merge K Sorted Lists
35. Task Scheduler

### Graphs (10)
36. Number of Islands (DFS or BFS)
37. Clone Graph
38. Course Schedule (cycle detection / topo sort)
39. Pacific Atlantic Water Flow
40. Word Ladder (BFS)
41. Shortest Path in Binary Matrix (BFS)
42. Network Delay Time (Dijkstra)
43. Cheapest Flights Within K Stops (Bellman-Ford or modified Dijkstra)
44. Alien Dictionary (topo sort)
45. Redundant Connection (Union-Find)

### Tries (3)
46. Implement Trie
47. Word Search II
48. Design Add and Search Words

### Sliding window / two pointers (5)
49. Minimum Window Substring
50. Longest Repeating Character Replacement
51. Container With Most Water
52. Trapping Rain Water
53. Permutation in String

### Backtracking (5)
54. Subsets
55. Permutations
56. Combination Sum
57. Word Search
58. N-Queens

### DP — 1D (8)
59. Climbing Stairs
60. House Robber
61. Coin Change
62. Longest Increasing Subsequence
63. Word Break
64. Decode Ways
65. Maximum Product Subarray
66. Partition Equal Subset Sum

### DP — 2D (6)
67. Longest Common Subsequence
68. Edit Distance
69. Unique Paths
70. Longest Palindromic Subsequence
71. 0/1 Knapsack
72. Best Time to Buy and Sell Stock with Cooldown

### Bit manipulation (4)
73. Single Number
74. Number of 1 Bits
75. Counting Bits
76. Sum of Two Integers (without + or -)

### Design (4)
77. LRU Cache
78. LFU Cache
79. Design Twitter
80. Design In-Memory File System

---

<a name="apple-tactics"></a>
## 🎯 Apple-style interview tactics

Apple's bar at the new-grad level is roughly:
- 1–2 medium LeetCode-style questions per coding round
- Focus on *correctness, communication, edge cases, complexity*
- They like classics: tree traversal, graph BFS/DFS, DP basics
- Often a follow-up "what if the input is huge / streaming / distributed?" — connects DSA to system design

### The 4-step framework for any DSA question

1. **Restate.** "So I have an array of integers and I need to return the indices of two numbers that sum to target?" — confirms understanding, buys you 30 seconds to think.

2. **Discuss.** "The brute force is O(n²). Can we do better with a hash map for O(n)?" — shows complexity awareness BEFORE coding.

3. **Code.** Talk through what each line does as you write. Use clear names.

4. **Test.** "Let me trace through an example." Pick a test case, walk through it. Then talk about edge cases: empty input? single element? duplicates? negatives?

### What separates good from great

- You verbalize trade-offs: "We could also use a sorted set here, which would give us O(n log n) but with O(1) extra space."
- You don't pretend. If you don't see the optimal, code the brute force, then say "I think we can do better — let me think about caching subresults."
- You read your own code. Run through it line by line at the end. Catches off-by-one errors.

---

<a name="what-you-should-now-know"></a>
## ✅ What you should now know

- [ ] Big-O notation, common complexities, amortized analysis
- [ ] Arrays, hash maps, linked lists, stacks, queues, heaps, tries
- [ ] BST operations, balanced trees conceptually
- [ ] BFS, DFS, Dijkstra, topological sort
- [ ] Two pointers, sliding window, prefix sums
- [ ] Recursion + backtracking template
- [ ] DP — top-down memo and bottom-up table
- [ ] Bit tricks
- [ ] LRU cache implemented in Go AND Rust
- [ ] At least 30 of the 80 problems solved (target 60+ before Apple interview)

---

**Next:** [Phase 4 — Databases](../phase-04-databases/README.md)
