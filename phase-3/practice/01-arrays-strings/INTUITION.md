# Module 1 — Arrays & Strings: Intuition

> **Read this before solving any problem in this module. Re-read once a week until you can teach it.**

---

## What problem do arrays solve?

You have a bunch of values. You need:
- Constant-time access by position.
- Predictable memory layout (fast for CPUs to scan).
- Cheap iteration in order.

That's it. Arrays are the most boring data structure on purpose — they're the one everything else is built on.

The reason arrays show up in ~40% of interview questions isn't because arrays are interesting. It's because **most real problems start with "I have a list of things"** and the question is whether you can manipulate that list intelligently.

---

## The mental model

A Python `list` is a contiguous block of memory holding pointers to objects.

```
nums = [10, 20, 30, 40, 50]

Index:    0    1    2    3    4
Memory: [10] [20] [30] [40] [50]
```

Three things follow from this:

1. **`nums[i]` is O(1)** — the computer multiplies `i * pointer_size` and jumps directly to that address.
2. **`nums.append(x)` is O(1) amortized** — Python over-allocates, so most appends are free. Occasionally it has to copy everything to a bigger block.
3. **`nums.insert(0, x)` is O(n)** — every existing element has to shift right by one slot.

If you remember nothing else: **front operations are expensive, back operations are cheap, indexed access is free.**

---

## When do you reach for an array?

You don't, usually. You reach for an array because the input *is* an array. The real question is: **what do you do with it?**

There are five "moves" you can make on an array. Almost every array problem is one of these:

### Move 1: Scan once and remember things

The brute force is "for each element, look at all other elements" → O(n²). The fix is almost always "scan once, remember useful things in a hash map or variable."

**Trigger phrases:**
- "find two numbers that sum to..."
- "find the first/longest/shortest..."
- "has any pair of...?"

### Move 2: Two pointers

The array is sorted, OR symmetry matters (palindromes), OR you're comparing pairs from opposite ends.

**Trigger phrases:**
- "is it a palindrome"
- "find a pair in a sorted array such that..."
- "find triplets..."
- "remove duplicates from a sorted array in place"

### Move 3: Sliding window

You care about a contiguous chunk (subarray, substring) and the answer is "the best chunk satisfying some condition."

**Trigger phrases:**
- "longest substring..."
- "smallest subarray with sum >= k"
- "max sum of k consecutive elements"

### Move 4: Prefix sums

You're asked range-sum-style questions repeatedly, OR the problem reduces to "how many times does the running sum equal X."

**Trigger phrases:**
- "sum of elements between i and j" (multiple queries)
- "how many subarrays sum to k"
- "find the equilibrium index"

### Move 5: Sort first, then scan

The problem becomes trivial once the array is sorted. You're paying O(n log n) up front to make the rest O(n).

**Trigger phrases:**
- "merge intervals"
- "find duplicates" (when memory is constrained)
- "kth smallest"

---

## The 30-second pattern recognition test

When you see an array problem, ask in order:

1. **Is the input sorted?** → two pointers is on the table.
2. **Am I looking for a contiguous chunk?** → sliding window or prefix sum.
3. **Am I asked about pairs, sums, or "have I seen this before"?** → hash map.
4. **Is the brute force O(n²) and I'm scanning every pair?** → can I scan once and remember?
5. **Would sorting make this trivial?** → sort first, then scan.

If none of these click in 30 seconds, the problem might not be a pure array problem — it might require a specific DS (heap, stack, etc.).

---

## Strings are arrays. With one twist.

In Python, strings are **immutable**. Three consequences:

1. `s[i] = 'x'` doesn't work. You'd need `list(s)`, modify, then `''.join(...)`.
2. Building a string with `+=` in a loop is O(n²). Use `''.join(list_of_strings)` instead.
3. String slicing `s[1:5]` creates a new string. Be aware in tight loops.

Beyond that, strings behave exactly like arrays of characters. Every array technique works on them.

---

## The traps that waste interview time

These are the ways smart people lose array problems:

1. **Forgetting Python list slicing copies.** `nums[1:4]` is a new list, costing O(k). If you slice inside a loop, your "O(n) solution" is secretly O(n²).

2. **Using `list.pop(0)` as a queue.** It's O(n). Use `collections.deque`.

3. **Modifying a list while iterating it.** Mutate by index or build a new list.

4. **Off-by-one errors with two pointers.** Always ask: do my pointers cross? Should `left < right` or `left <= right`? When does the loop terminate?

5. **Forgetting that `int / int` returns a float in Python 3.** Use `//` for integer division. This bites people on binary search.

6. **Confusing "subarray" (contiguous) with "subsequence" (not necessarily contiguous).** Subarrays are sliding-window territory. Subsequences are usually DP territory.

---

## Complexity instincts you should have

After this module, these should feel automatic:

| Operation | Cost | Why |
|---|---|---|
| `nums[i]` | O(1) | Direct memory jump |
| `nums.append(x)` | O(1) amortized | Over-allocation |
| `nums.insert(0, x)` | O(n) | Shift everything right |
| `nums.pop()` | O(1) | Remove from end |
| `nums.pop(0)` | O(n) | Shift everything left |
| `x in nums` | O(n) | Linear scan |
| `x in set` | O(1) avg | Hash lookup |
| `nums[1:5]` | O(k) | Copies k elements |
| `sorted(nums)` | O(n log n) | Comparison sort |
| `nums.sort()` | O(n log n) | In-place |
| `min(nums)` | O(n) | Scan once |
| `''.join(parts)` | O(total length) | Allocates once |

If an interviewer asks "what's the complexity of `' '.join(words)`," the answer is O(total length of all words combined), NOT O(number of words). This trips people up.

---

## How to use this module

1. Read this doc fully. Don't skim.
2. Take the quiz in `QUIZ.md`. Don't peek at solutions.
3. Work the 8 problems in order. Each one is teaching a specific move.
4. Try the interview simulation in `interview-sim.md` last.
5. Only then look at `SOLUTIONS.md`.

The order of the problems matters — problem 3 builds on problem 1. Don't shuffle.
