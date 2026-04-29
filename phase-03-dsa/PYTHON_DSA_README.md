# Phase 3 — Data Structures & Algorithms in Python

> This is the Python version of the DSA phase. It mirrors the Go README, but uses Pythonic tools like `list`, `dict`, `set`, `deque`, `heapq`, `Counter`, and `defaultdict`.
>
> The goal is not to memorize every problem. The goal is to recognize the pattern, explain the intuition, write the clean template, and state the complexity.

**Time:** 8–14 days, plus continuous practice afterward.

**You'll know you're done when:** you can recognize a problem's pattern in 30 seconds, code the canonical Python solution in 15 minutes, and explain time/space complexity clearly.

---

## Table of contents

1. [What does this even mean? — Algorithms](#what-algorithms-mean)
2. [Module 3.1 — Big-O notation](#module-31--big-o)
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
15. [Python DSA cheat sheet](#python-dsa-cheat-sheet)
16. [Practice order](#practice-order)

---

<a name="what-algorithms-mean"></a>
## What does this even mean? — Algorithms

An **algorithm** is a recipe. It is a finite set of steps to solve a problem.

A **data structure** is a way to organize data so certain operations are cheap.

The same problem can be easy or hard depending on the data structure:

- Searching an unsorted list: O(n)
- Searching a sorted list with binary search: O(log n)
- Searching a dictionary/set: O(1) average case

DSA teaches you how to think about cost.

When an interviewer asks for a better solution, they usually mean:

- Can you avoid nested loops?
- Can you remember useful information as you scan?
- Can you avoid recalculating the same thing?
- Can you use the structure of the input?

---

<a name="module-31--big-o"></a>
## Module 3.1 — Big-O notation

Big-O describes how runtime or memory grows as input size grows.

| Class | Name | Example |
|---|---|---|
| O(1) | Constant | list index, dict lookup average case |
| O(log n) | Logarithmic | binary search |
| O(n) | Linear | scan once |
| O(n log n) | Linearithmic | sorting |
| O(n²) | Quadratic | nested loops |
| O(2ⁿ) | Exponential | brute-force subsets |
| O(n!) | Factorial | permutations |

### Python operation costs

| Structure | Operation | Cost |
|---|---|---|
| `list` | index access | O(1) |
| `list` | append | O(1) amortized |
| `list` | insert/delete at front | O(n) |
| `dict` | lookup/insert/delete | O(1) average |
| `set` | lookup/insert/delete | O(1) average |
| `deque` | append/pop both ends | O(1) |
| `heapq` | push/pop | O(log n) |
| sorting | `list.sort()` / `sorted()` | O(n log n) |

Interview habit:

```text
Time: O(n)
Space: O(k), where k is the number of unique items stored.
```

Always state both time and space.

---

<a name="module-32--arrays"></a>
## Module 3.2 — Arrays and strings

In Python, the main array-like structure is a `list`.

```python
nums = [1, 2, 3, 4, 5]
print(nums[2])      # 3
nums.append(6)      # O(1) amortized
sub = nums[1:4]     # copy: [2, 3, 4]
```

Important Python difference from Go:

- Go slices are views over the same underlying array.
- Python slices create a new list copy.

### Reverse in place

```python
def reverse(nums):
    left = 0
    right = len(nums) - 1

    while left < right:
        nums[left], nums[right] = nums[right], nums[left]
        left += 1
        right -= 1
```

Time: O(n)

Space: O(1)

### Find duplicates with a set

```python
def has_duplicate(nums):
    seen = set()

    for x in nums:
        if x in seen:
            return True
        seen.add(x)

    return False
```

Time: O(n)

Space: O(n)

### Two Sum

```python
def two_sum(nums, target):
    index = {}  # value -> index

    for i, x in enumerate(nums):
        need = target - x

        if need in index:
            return [index[need], i]

        index[x] = i

    return [-1, -1]
```

Intuition:

Instead of scanning the whole array for the matching number, store numbers you already passed.

### Strings

Python strings are immutable.

```python
s = "hello"
print(s[0])      # h
print(len(s))    # 5
```

If you need to build a string repeatedly, do not do this:

```python
s = ""
for word in words:
    s += word     # can become slow
```

Do this:

```python
result = []
for word in words:
    result.append(word)

s = "".join(result)
```

---

<a name="module-33--hashmaps"></a>
## Module 3.3 — Hash maps and hash sets

Python uses:

- `dict` for hash maps
- `set` for hash sets
- `Counter` for frequency counts
- `defaultdict` for default values

### Hash map basics

```python
m = {}
m["foo"] = 42

if "foo" in m:
    print(m["foo"])

del m["foo"]
```

### Frequency map

```python
from collections import defaultdict

freq = defaultdict(int)

for x in nums:
    freq[x] += 1
```

Or:

```python
from collections import Counter

freq = Counter(nums)
```

### Group anagrams

```python
from collections import defaultdict


def group_anagrams(words):
    groups = defaultdict(list)

    for word in words:
        key = tuple(sorted(word))
        groups[key].append(word)

    return list(groups.values())
```

Intuition:

Words that are anagrams have the same sorted form.

### Set template

```python
seen = set()

for x in nums:
    if x in seen:
        return True
    seen.add(x)
```

Use a set when you care about existence only.

---

<a name="module-34--linked-lists"></a>
## Module 3.4 — Linked lists

A linked list is a chain of nodes.

```python
class ListNode:
    def __init__(self, val=0, next=None):
        self.val = val
        self.next = next
```

Pros:

- O(1) insert/delete if you already have the node

Cons:

- O(n) search
- O(n) random access
- worse cache locality than arrays

### Reverse a linked list

```python
def reverse_list(head):
    prev = None
    curr = head

    while curr:
        nxt = curr.next
        curr.next = prev
        prev = curr
        curr = nxt

    return prev
```

Intuition:

Walk node by node and flip each pointer backward.

### Detect cycle

```python
def has_cycle(head):
    slow = head
    fast = head

    while fast and fast.next:
        slow = slow.next
        fast = fast.next.next

        if slow is fast:
            return True

    return False
```

Intuition:

If there is a cycle, the fast pointer eventually catches the slow pointer.

### Find middle node

```python
def middle_node(head):
    slow = head
    fast = head

    while fast and fast.next:
        slow = slow.next
        fast = fast.next.next

    return slow
```

When fast reaches the end, slow is in the middle.

---

<a name="module-35--stacks-queues"></a>
## Module 3.5 — Stacks and queues

### Stack

A stack is LIFO: last in, first out.

Use Python `list`:

```python
stack = []
stack.append(10)
stack.append(20)
print(stack.pop())  # 20
```

Common uses:

- valid parentheses
- undo
- DFS
- monotonic stack
- expression parsing

### Valid parentheses

```python
def is_valid(s):
    stack = []
    pairs = {
        ')': '(',
        ']': '[',
        '}': '{',
    }

    for char in s:
        if char in "([{":
            stack.append(char)
        elif char in pairs:
            if not stack or stack[-1] != pairs[char]:
                return False
            stack.pop()

    return len(stack) == 0
```

### Queue

A queue is FIFO: first in, first out.

Do not use `list.pop(0)` for a queue because it is O(n).

Use `deque`:

```python
from collections import deque

queue = deque()
queue.append(10)
queue.append(20)
print(queue.popleft())  # 10
```

Common uses:

- BFS
- scheduling
- level-order traversal
- shortest path in unweighted graphs

---

<a name="module-36--trees"></a>
## Module 3.6 — Trees: BST, AVL, B-tree

### Binary tree basics

```python
class TreeNode:
    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right
```

### DFS traversals

```python
def inorder(root):
    result = []

    def dfs(node):
        if not node:
            return
        dfs(node.left)
        result.append(node.val)
        dfs(node.right)

    dfs(root)
    return result
```

Orders:

- Inorder: left, root, right
- Preorder: root, left, right
- Postorder: left, right, root

For BSTs, inorder traversal returns values in sorted order.

### Level-order traversal

```python
from collections import deque


def level_order(root):
    if not root:
        return []

    result = []
    queue = deque([root])

    while queue:
        level = []

        for _ in range(len(queue)):
            node = queue.popleft()
            level.append(node.val)

            if node.left:
                queue.append(node.left)
            if node.right:
                queue.append(node.right)

        result.append(level)

    return result
```

### Binary Search Tree search

```python
def search_bst(root, val):
    curr = root

    while curr:
        if val == curr.val:
            return curr
        elif val < curr.val:
            curr = curr.left
        else:
            curr = curr.right

    return None
```

### Binary Search Tree insert

```python
def insert_bst(root, val):
    if not root:
        return TreeNode(val)

    if val < root.val:
        root.left = insert_bst(root.left, val)
    elif val > root.val:
        root.right = insert_bst(root.right, val)

    return root
```

### AVL and Red-Black trees

You usually do not implement these in interviews.

Know the idea:

- They are self-balancing binary search trees.
- They keep operations O(log n).
- AVL is stricter with balance.
- Red-Black is more common in standard libraries.

### B-tree intuition

A B-tree is a multi-way search tree.

Instead of each node having only 2 children, a B-tree node can store many keys and many children.

Why databases use B-trees:

- Disk access is slow.
- B-trees keep the tree shallow.
- Each node can store many keys, reducing disk reads.

Used in:

- database indexes
- filesystems
- storage engines

---

<a name="module-37--heaps"></a>
## Module 3.7 — Heaps and priority queues

A heap is used when you repeatedly need the smallest or largest item.

Python's `heapq` is a min heap.

```python
import heapq

heap = []
heapq.heappush(heap, 3)
heapq.heappush(heap, 1)
heapq.heappush(heap, 2)

print(heapq.heappop(heap))  # 1
```

### Max heap trick

```python
import heapq

heap = []

for x in nums:
    heapq.heappush(heap, -x)

largest = -heapq.heappop(heap)
```

### Top K largest

```python
import heapq


def top_k_largest(nums, k):
    heap = []

    for x in nums:
        heapq.heappush(heap, x)

        if len(heap) > k:
            heapq.heappop(heap)

    return heap
```

Intuition:

Keep a min heap of size k. If it grows larger than k, remove the smallest. The remaining values are the k largest.

Time: O(n log k)

Space: O(k)

### Merge k sorted lists

```python
import heapq


def merge_k_lists(lists):
    heap = []

    for i, node in enumerate(lists):
        if node:
            heapq.heappush(heap, (node.val, i, node))

    dummy = ListNode(0)
    curr = dummy

    while heap:
        _, i, node = heapq.heappop(heap)
        curr.next = node
        curr = curr.next

        if node.next:
            heapq.heappush(heap, (node.next.val, i, node.next))

    return dummy.next
```

---

<a name="module-38--tries"></a>
## Module 3.8 — Tries

A trie is a prefix tree.

Use it when you need fast prefix lookup.

Common uses:

- autocomplete
- spell check
- word search
- prefix matching

```python
class TrieNode:
    def __init__(self):
        self.children = {}
        self.is_end = False


class Trie:
    def __init__(self):
        self.root = TrieNode()

    def insert(self, word):
        node = self.root

        for char in word:
            if char not in node.children:
                node.children[char] = TrieNode()
            node = node.children[char]

        node.is_end = True

    def search(self, word):
        node = self._find(word)
        return node is not None and node.is_end

    def starts_with(self, prefix):
        return self._find(prefix) is not None

    def _find(self, text):
        node = self.root

        for char in text:
            if char not in node.children:
                return None
            node = node.children[char]

        return node
```

Time:

- Insert: O(L)
- Search: O(L)
- Prefix search: O(L)

Where L is the word length.

---

<a name="module-39--graphs"></a>
## Module 3.9 — Graphs: BFS, DFS, Dijkstra, topological sort

### Graph representation

Adjacency list is most common:

```python
graph = {
    0: [1, 2],
    1: [0, 3],
    2: [0],
    3: [1],
}
```

Or build one:

```python
from collections import defaultdict


def build_graph(edges):
    graph = defaultdict(list)

    for a, b in edges:
        graph[a].append(b)
        graph[b].append(a)

    return graph
```

### BFS

```python
from collections import deque


def bfs(graph, start):
    visited = {start}
    queue = deque([start])
    order = []

    while queue:
        node = queue.popleft()
        order.append(node)

        for nei in graph[node]:
            if nei not in visited:
                visited.add(nei)
                queue.append(nei)

    return order
```

Intuition:

BFS explores all nodes at distance 1, then distance 2, then distance 3.

Use BFS for shortest path in unweighted graphs.

### DFS

```python
def dfs(graph, start):
    visited = set()
    order = []

    def visit(node):
        if node in visited:
            return

        visited.add(node)
        order.append(node)

        for nei in graph[node]:
            visit(nei)

    visit(start)
    return order
```

Use DFS for:

- connected components
- cycle detection
- path existence
- tree recursion
- grid islands

### Dijkstra's algorithm

Use for shortest path in a weighted graph with no negative weights.

```python
import heapq


def dijkstra(graph, source):
    dist = {source: 0}
    heap = [(0, source)]

    while heap:
        curr_dist, node = heapq.heappop(heap)

        if curr_dist > dist[node]:
            continue

        for nei, weight in graph[node]:
            new_dist = curr_dist + weight

            if nei not in dist or new_dist < dist[nei]:
                dist[nei] = new_dist
                heapq.heappush(heap, (new_dist, nei))

    return dist
```

Graph shape:

```python
graph = {
    0: [(1, 4), (2, 1)],
    1: [(3, 1)],
    2: [(1, 2), (3, 5)],
    3: [],
}
```

### Topological sort

Use for directed acyclic graphs.

Examples:

- course schedule
- build systems
- dependency ordering

```python
from collections import defaultdict, deque


def topo_sort(num_nodes, edges):
    graph = defaultdict(list)
    indegree = [0] * num_nodes

    for a, b in edges:
        graph[a].append(b)
        indegree[b] += 1

    queue = deque()

    for i in range(num_nodes):
        if indegree[i] == 0:
            queue.append(i)

    order = []

    while queue:
        node = queue.popleft()
        order.append(node)

        for nei in graph[node]:
            indegree[nei] -= 1
            if indegree[nei] == 0:
                queue.append(nei)

    if len(order) != num_nodes:
        return []  # cycle exists

    return order
```

---

<a name="module-310--patterns"></a>
## Module 3.10 — Two pointers, sliding window, prefix sums

These are patterns, not data structures.

Recognizing them turns hard problems into easy ones.

---

### Two pointers

Two pointers means using two indexes to avoid extra loops.

Use it when:

- the input is sorted
- you compare pairs
- you need to reverse something
- you need left/right boundaries

```python
def is_palindrome(s):
    left = 0
    right = len(s) - 1

    while left < right:
        if s[left] != s[right]:
            return False
        left += 1
        right -= 1

    return True
```

### Three Sum

```python
def three_sum(nums):
    nums.sort()
    result = []

    for i in range(len(nums) - 2):
        if i > 0 and nums[i] == nums[i - 1]:
            continue

        left = i + 1
        right = len(nums) - 1

        while left < right:
            total = nums[i] + nums[left] + nums[right]

            if total < 0:
                left += 1
            elif total > 0:
                right -= 1
            else:
                result.append([nums[i], nums[left], nums[right]])

                while left < right and nums[left] == nums[left + 1]:
                    left += 1
                while left < right and nums[right] == nums[right - 1]:
                    right -= 1

                left += 1
                right -= 1

    return result
```

---

### Sliding window

Sliding window is for contiguous subarrays or substrings.

Intuition:

> Keep a live window from `left` to `right`. Expand right to include more. If the window becomes invalid, move left until it is valid again.

This avoids checking every possible window.

Brute force checks many repeated windows and often becomes O(n²).

Sliding window usually makes it O(n) because each element enters and leaves the window at most once.

### Fixed-size sliding window

Use when the window size is always `k`.

```python
def max_sum_size_k(nums, k):
    window_sum = 0
    best = float('-inf')

    for right in range(len(nums)):
        window_sum += nums[right]

        if right >= k:
            window_sum -= nums[right - k]

        if right >= k - 1:
            best = max(best, window_sum)

    return best
```

### Variable-size sliding window template

```python
left = 0
answer = 0
state = {}

for right in range(len(nums)):
    # Add nums[right] to the window/state

    while window_is_invalid:
        # Remove nums[left] from the window/state
        left += 1

    # Update answer using the valid window
```

### Longest substring without repeating characters

```python
def length_of_longest_substring(s):
    seen = {}
    left = 0
    best = 0

    for right, char in enumerate(s):
        if char in seen and seen[char] >= left:
            left = seen[char] + 1

        seen[char] = right
        best = max(best, right - left + 1)

    return best
```

Why it works:

- `seen` stores the latest index of each character.
- If the same character appears inside the current window, the window is invalid.
- Move `left` right after the previous copy.
- Never move `left` backward.

### Minimum size subarray sum

```python
def min_subarray_len(target, nums):
    left = 0
    total = 0
    best = float('inf')

    for right in range(len(nums)):
        total += nums[right]

        while total >= target:
            best = min(best, right - left + 1)
            total -= nums[left]
            left += 1

    return 0 if best == float('inf') else best
```

---

### Prefix sums

Prefix sum is for fast range sums.

```text
prefix[i] = sum of nums before index i
sum nums[left:right+1] = prefix[right + 1] - prefix[left]
```

```python
def build_prefix(nums):
    prefix = [0]

    for x in nums:
        prefix.append(prefix[-1] + x)

    return prefix


def range_sum(prefix, left, right):
    return prefix[right + 1] - prefix[left]
```

### Subarray sum equals k

```python
from collections import defaultdict


def subarray_sum(nums, k):
    count = 0
    curr = 0
    seen = defaultdict(int)
    seen[0] = 1

    for x in nums:
        curr += x
        count += seen[curr - k]
        seen[curr] += 1

    return count
```

---

<a name="module-311--recursion"></a>
## Module 3.11 — Recursion and backtracking

Recursion means a function calls itself.

Every recursive solution needs:

1. Base case: when to stop
2. Recursive case: how to reduce the problem

```python
def factorial(n):
    if n <= 1:
        return 1
    return n * factorial(n - 1)
```

### Backtracking intuition

Backtracking means:

1. Choose
2. Explore
3. Undo

Use it when you need all possible answers.

### Generic backtracking template

```python
def backtrack(path, choices):
    if done(path):
        result.append(path[:])
        return

    for choice in choices:
        if not valid(choice):
            continue

        path.append(choice)
        backtrack(path, choices)
        path.pop()
```

### Generate permutations

```python
def permute(nums):
    result = []
    path = []
    used = [False] * len(nums)

    def backtrack():
        if len(path) == len(nums):
            result.append(path[:])
            return

        for i in range(len(nums)):
            if used[i]:
                continue

            used[i] = True
            path.append(nums[i])

            backtrack()

            path.pop()
            used[i] = False

    backtrack()
    return result
```

### Generate subsets

```python
def subsets(nums):
    result = []
    path = []

    def backtrack(start):
        result.append(path[:])

        for i in range(start, len(nums)):
            path.append(nums[i])
            backtrack(i + 1)
            path.pop()

    backtrack(0)
    return result
```

---

<a name="module-312--dp"></a>
## Module 3.12 — Dynamic programming

Dynamic programming means solving overlapping subproblems once and reusing the answers.

Use DP when:

- the same subproblem appears repeatedly
- the answer depends on smaller answers
- brute force recursion repeats work

The most important question:

```text
What does dp[i] mean?
```

If you cannot define `dp[i]`, you do not have a DP solution yet.

### Top-down memoization

```python
def fib(n):
    memo = {}

    def solve(k):
        if k < 2:
            return k
        if k in memo:
            return memo[k]

        memo[k] = solve(k - 1) + solve(k - 2)
        return memo[k]

    return solve(n)
```

### Bottom-up tabulation

```python
def fib(n):
    if n < 2:
        return n

    dp = [0] * (n + 1)
    dp[1] = 1

    for i in range(2, n + 1):
        dp[i] = dp[i - 1] + dp[i - 2]

    return dp[n]
```

### Space-optimized DP

```python
def fib(n):
    if n < 2:
        return n

    prev2 = 0
    prev1 = 1

    for _ in range(2, n + 1):
        curr = prev1 + prev2
        prev2 = prev1
        prev1 = curr

    return prev1
```

### Climbing stairs

```python
def climb_stairs(n):
    if n <= 2:
        return n

    prev2 = 1
    prev1 = 2

    for _ in range(3, n + 1):
        curr = prev1 + prev2
        prev2 = prev1
        prev1 = curr

    return prev1
```

Intuition:

To reach step `i`, you either came from `i - 1` or `i - 2`.

So:

```text
dp[i] = dp[i - 1] + dp[i - 2]
```

---

<a name="module-313--bits"></a>
## Module 3.13 — Bit manipulation

Bit manipulation works with binary representation.

Common operations:

| Operation | Meaning |
|---|---|
| `x & 1` | check if odd |
| `x >> 1` | divide by 2 |
| `x << 1` | multiply by 2 |
| `x ^ y` | XOR |
| `x & (x - 1)` | remove lowest set bit |

### Check if a number is odd

```python
def is_odd(x):
    return (x & 1) == 1
```

### Count set bits

```python
def count_bits(x):
    count = 0

    while x:
        x &= x - 1
        count += 1

    return count
```

Intuition:

`x & (x - 1)` removes the lowest set bit.

### Single number

Every number appears twice except one.

```python
def single_number(nums):
    ans = 0

    for x in nums:
        ans ^= x

    return ans
```

Why it works:

- `x ^ x = 0`
- `x ^ 0 = x`
- XOR cancels duplicates

---

<a name="python-dsa-cheat-sheet"></a>
## Python DSA cheat sheet

### Imports

```python
from collections import defaultdict, Counter, deque
import heapq
import bisect
```

### Common patterns

| Problem wording | Pattern |
|---|---|
| two numbers sum to target | hash map / two pointers if sorted |
| longest substring/subarray | sliding window |
| shortest subarray with condition | sliding window |
| range sum | prefix sum |
| top k | heap |
| valid parentheses | stack |
| shortest path unweighted | BFS |
| all combinations/permutations | backtracking |
| same subproblem repeats | dynamic programming |
| sorted array | binary search / two pointers |
| tree level order | BFS |
| tree path/depth | DFS |
| prefix lookup | trie |
| dependency order | topological sort |

### Explanation template

```text
The brute-force approach would be ____.
That repeats work because ____.
The key observation is ____.
I can store/track ____ using ____.
Each element is processed ____ times.
So the time complexity is ____ and the space complexity is ____.
```

Sliding window example:

```text
The brute-force approach checks every substring, which is O(n²).
The key observation is that I only need one active window.
I expand the right side as I scan.
When the window becomes invalid, I move the left side forward.
Each character enters and leaves the window at most once, so the time is O(n).
```

---

<a name="practice-order"></a>
## Practice order

1. Arrays and strings
2. Hash maps and sets
3. Linked lists
4. Stacks and queues
5. Trees
6. Heaps
7. Tries
8. Graphs
9. Two pointers
10. Sliding window
11. Prefix sums
12. Recursion and backtracking
13. Dynamic programming
14. Bit manipulation

Do not try to master everything at once. Pattern recognition comes first. Speed comes after repetition.
