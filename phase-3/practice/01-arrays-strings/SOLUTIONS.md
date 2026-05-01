# Module 1 — Arrays & Strings: Solutions

> **DO NOT OPEN THIS FILE** until you have a working solution to the problem in question. Each section is in a collapsible — open one at a time. The discipline of struggling first is what builds intuition.

---

## Problem 1 — Two Sum

<details>
<summary>Reveal solution + explanation</summary>

```python
from typing import List


def two_sum(nums: List[int], target: int) -> List[int]:
    """
    Time:  O(n) — single pass, each lookup/insert is O(1) average.
    Space: O(n) — the hash map can hold up to n entries.
    """
    seen = {}  # value -> index

    for i, x in enumerate(nums):
        need = target - x
        if need in seen:
            return [seen[need], i]
        seen[x] = i

    return []  # problem guarantees a solution; defensive fallback
```

### Why this works

The brute force tries every pair: O(n²). The optimization is: instead of asking "for each pair, do they sum to target?", ask "for each `x`, has `target - x` been seen before?" That's a single-element question, which a hash map answers in O(1).

### The order trap

Look closely at the loop:
1. Compute `need`.
2. **Check if `need` is in `seen`.**
3. **THEN** insert `x` into `seen`.

If you swapped 2 and 3, the test case `nums=[3, 3], target=6` would break: when you process the first 3, you'd insert it, then when you process it again you'd find a "match" with itself. Always check-then-insert.

### Common bugs

- Returning `[i, j]` in the wrong order — usually doesn't matter, but read the problem.
- Using `if seen[need]:` instead of `if need in seen:` — the former crashes on missing keys.
- Forgetting to handle "no solution" — the problem says one exists, but real code should return something safe.

</details>

---

## Problem 2 — Best Time to Buy and Sell Stock

<details>
<summary>Reveal solution + explanation</summary>

```python
from typing import List


def max_profit(prices: List[int]) -> int:
    """
    Time:  O(n) — single pass.
    Space: O(1) — two scalars.
    """
    if not prices:
        return 0

    min_so_far = prices[0]
    best = 0

    for price in prices[1:]:
        # IMPORTANT: compute profit BEFORE updating min, so we never
        # buy and sell on the same day.
        best = max(best, price - min_so_far)
        min_so_far = min(min_so_far, price)

    return best
```

### Why this works

Reframe: "What's the best profit if I sell on day `i`?" Answer: `prices[i] - min(prices[0..i-1])`. The maximum over all i is the answer.

You don't need to remember every past price. Just the cheapest one. That's the collapse to O(1) space.

### The order subtlety

`best = max(best, price - min_so_far)` happens BEFORE `min_so_far = min(min_so_far, price)`. If you flipped the order, you might compute `min_so_far` as today's price and then immediately compute "profit" as `price - price = 0`. Subtle but real bug.

You can also skip the special case for `prices[0]` by initializing `min_so_far = float('inf')` and looping from index 0.

</details>

---

## Problem 3 — Contains Duplicate

<details>
<summary>Reveal solution + explanation</summary>

```python
from typing import List


def contains_duplicate(nums: List[int]) -> bool:
    """
    Time:  O(n) — one pass.
    Space: O(n) — set holds up to n unique elements.
    """
    seen = set()

    for x in nums:
        if x in seen:
            return True
        seen.add(x)

    return False
```

### Why a set, not a dict

We don't care about *where* we saw the duplicate, only whether we saw it. A set is the precise data structure for "membership only." It uses slightly less memory than a dict (no value slot).

### Tradeoffs to verbalize in an interview

- "I used a hash set — O(n) time, O(n) space."
- "Alternative: sort first, scan adjacent pairs. O(n log n) time, O(1) extra space if sorting in place. Better when memory is tight."
- "Alternative: brute force — O(n²) time, O(1) space. Useful only for tiny inputs."

### The Pythonic shortcut

`return len(set(nums)) != len(nums)` works but doesn't short-circuit. If the duplicate is at the start, you scan the whole array unnecessarily. The explicit loop is better in interviews because it shows the algorithmic thinking.

</details>

---

## Problem 4 — Valid Palindrome

<details>
<summary>Reveal solution + explanation</summary>

```python
def is_palindrome(s: str) -> bool:
    """
    Time:  O(n) — pointers traverse the string at most once total.
    Space: O(1) — two pointers, no extra string built.
    """
    left, right = 0, len(s) - 1

    while left < right:
        # Skip non-alphanumeric on the left
        while left < right and not s[left].isalnum():
            left += 1
        # Skip non-alphanumeric on the right
        while left < right and not s[right].isalnum():
            right -= 1

        if s[left].lower() != s[right].lower():
            return False

        left += 1
        right -= 1

    return True
```

### The critical detail

The `left < right` check inside the inner skip loops. Without it, on `"....."` (all non-alphanumeric), `left` would keep advancing past `right`, potentially indexing out of bounds (or worse, comparing garbage).

### Alternative — clean then check

```python
cleaned = [c.lower() for c in s if c.isalnum()]
return cleaned == cleaned[::-1]
```

Works, but allocates O(n) extra space. The two-pointer version is what the interviewer wants.

### Why two pointers is the right tool

Palindromes are about symmetry around a center. Two pointers from opposite ends naturally express that symmetry. Any "is X symmetric?" problem on an array or string is two-pointer territory.

</details>

---

## Problem 5 — Three Sum

<details>
<summary>Reveal solution + explanation</summary>

```python
from typing import List


def three_sum(nums: List[int]) -> List[List[int]]:
    """
    Time:  O(n²) — outer loop n, inner two-pointer scan n each.
    Space: O(1) extra (output not counted); sorting is O(log n) stack.
    """
    nums.sort()
    n = len(nums)
    result = []

    for i in range(n - 2):
        # Optimization: once nums[i] > 0, three positives can't sum to 0.
        if nums[i] > 0:
            break

        # Skip duplicate first elements.
        if i > 0 and nums[i] == nums[i - 1]:
            continue

        left, right = i + 1, n - 1
        while left < right:
            total = nums[i] + nums[left] + nums[right]

            if total < 0:
                left += 1
            elif total > 0:
                right -= 1
            else:
                result.append([nums[i], nums[left], nums[right]])

                # Skip duplicate left and right values
                while left < right and nums[left] == nums[left + 1]:
                    left += 1
                while left < right and nums[right] == nums[right - 1]:
                    right -= 1

                left += 1
                right -= 1

    return result
```

### The three places duplicates arise

1. **Outer loop:** if `nums[i]` is the same as `nums[i-1]`, we'd find the same triplets again. Skip.
2. **After finding a hit:** if there are multiple equal values around `left` (or `right`), we'd find the same triplet starting from the next equal value. Skip.
3. **Be careful with `i > 0`:** we skip duplicates only after the first iteration. Otherwise we'd never even start.

### Why sort first

Sorting enables two pointers. It also makes dedup trivial (duplicates are adjacent). The O(n log n) sort is dominated by the O(n²) main loop, so it's free.

### Common bug

Forgetting the `while left < right` check inside the dedup-skip loops. Without it, you can read past array bounds.

</details>

---

## Problem 6 — Product of Array Except Self

<details>
<summary>Reveal solution + explanation</summary>

```python
from typing import List


def product_except_self(nums: List[int]) -> List[int]:
    """
    Time:  O(n) — two passes.
    Space: O(1) extra (output array doesn't count).
    """
    n = len(nums)
    answer = [1] * n

    # First pass: answer[i] = product of everything LEFT of i.
    left_product = 1
    for i in range(n):
        answer[i] = left_product
        left_product *= nums[i]

    # Second pass: multiply in the product of everything RIGHT of i.
    right_product = 1
    for i in range(n - 1, -1, -1):
        answer[i] *= right_product
        right_product *= nums[i]

    return answer
```

### Trace through `[1, 2, 3, 4]`

After pass 1 (left products): `[1, 1, 2, 6]`
- answer[0] = 1 (nothing left of index 0)
- answer[1] = 1 (just nums[0])
- answer[2] = 1 * 2 = 2
- answer[3] = 1 * 2 * 3 = 6

Pass 2 (multiply by right products), going right-to-left:
- i=3: right_product = 1, answer[3] = 6 * 1 = 6, right_product becomes 4
- i=2: answer[2] = 2 * 4 = 8, right_product becomes 12
- i=1: answer[1] = 1 * 12 = 12, right_product becomes 24
- i=0: answer[0] = 1 * 24 = 24

Final: `[24, 12, 8, 6]`. ✓

### Why no division

Division fails on zeros. With one zero, all positions except the zero's slot become 0. With two or more zeros, the entire answer is zeros. The prefix/suffix approach handles zeros naturally.

</details>

---

## Problem 7 — Longest Substring Without Repeating Characters

<details>
<summary>Reveal solution + explanation</summary>

```python
def length_of_longest_substring(s: str) -> int:
    """
    Time:  O(n) — each character enters and leaves the window at most once.
    Space: O(min(n, alphabet_size)) for the hash map.
    """
    last_seen = {}  # char -> most recent index
    left = 0
    best = 0

    for right, ch in enumerate(s):
        # If ch is in the current window, jump left past its previous index.
        if ch in last_seen and last_seen[ch] >= left:
            left = last_seen[ch] + 1

        last_seen[ch] = right
        best = max(best, right - left + 1)

    return best
```

### The `>= left` check

Consider `"abba"`:
- right=0 'a': window [0,0], last_seen = {a:0}
- right=1 'b': window [0,1], last_seen = {a:0, b:1}
- right=2 'b': last_seen[b]=1, which is >= left (0). Jump left to 2. Window [2,2]. last_seen = {a:0, b:2}.
- right=3 'a': last_seen[a]=0, but 0 < left (2). The 'a' at index 0 is NOT in the current window. Don't move left. Window [2,3]. last_seen = {a:3, b:2}. best = 2.

Without the `>= left` check, you'd move `left` backward to 1, breaking everything.

### Why O(n)

`right` advances exactly n times. `left` only moves forward — when it moves, it can jump multiple positions, but the total distance it travels across the whole loop is at most n. Total work: 2n = O(n).

### Alternative — set-based

```python
seen = set()
left = 0
best = 0

for right, ch in enumerate(s):
    while ch in seen:
        seen.remove(s[left])
        left += 1
    seen.add(ch)
    best = max(best, right - left + 1)

return best
```

Cleaner to read, same complexity. The `last_seen` version is slightly faster in practice because it jumps `left` directly instead of stepping it forward.

</details>

---

## Problem 8 — Subarray Sum Equals K

<details>
<summary>Reveal solution + explanation</summary>

```python
from collections import defaultdict
from typing import List


def subarray_sum(nums: List[int], k: int) -> int:
    """
    Time:  O(n) — single pass.
    Space: O(n) — hash map of seen prefix sums.
    """
    count = 0
    running_sum = 0
    seen = defaultdict(int)
    seen[0] = 1   # empty prefix has sum 0; lets us count subarrays from index 0.

    for x in nums:
        running_sum += x
        # If (running_sum - k) was a previous prefix, the subarray between
        # there and here sums to k.
        count += seen[running_sum - k]
        seen[running_sum] += 1

    return count
```

### Trace through `[1, 1, 1], k=2`

Initial: seen = {0: 1}, running_sum = 0, count = 0

- x=1: running_sum=1. count += seen[1-2=-1] = 0. seen = {0:1, 1:1}.
- x=1: running_sum=2. count += seen[2-2=0] = 1. count=1. seen = {0:1, 1:1, 2:1}.
- x=1: running_sum=3. count += seen[3-2=1] = 1. count=2. seen = {0:1, 1:1, 2:1, 3:1}.

Final count: 2. ✓ (The two valid subarrays are nums[0:2] and nums[1:3].)

### Why sliding window doesn't work

Sliding window relies on: "if a window's condition is broken, shrinking it will help." With negative numbers, that's false. Adding an element might *decrease* the running sum, so "shrink when sum is too big" is no longer a valid rule.

Prefix sums sidestep this entirely. Instead of tracking a window, you track running totals and look for matching past totals.

### Why `{0: 1}`

If `running_sum == k`, that means the subarray from index 0 to the current index sums to k. We want this to count. To make `seen[running_sum - k] = seen[0]` evaluate to 1, we initialize `seen[0] = 1`.

Trace `[3], k=3` to convince yourself:
- x=3: running_sum=3. count += seen[3-3=0] = 1. count=1.

Without the `{0: 1}` init, we'd return 0 — wrong.

### Why `[0,0,0,0]` with k=0 returns 10

Every subarray of zeros sums to 0. There are 4+3+2+1 = 10 subarrays of a length-4 array, and all sum to 0. The trick: each prefix sum of 0 we've seen contributes to the count, and we keep accumulating them. The math works out to "n choose 2 + n" for an all-zeros array.

</details>

---

## Pattern recap

After working through all 8, your brain should have these grooves cut:

| Problem | The move | Why |
|---|---|---|
| Two Sum | Hash map | Pair lookup by complement |
| Best Time Stock | Single-scan + scalar | Collapse history to one variable |
| Contains Duplicate | Hash set | Existence-only |
| Valid Palindrome | Two pointers | Symmetry around center |
| Three Sum | Sort + two pointers | Sortedness enables direction; dedup is the trick |
| Product Except Self | Prefix + suffix | Range info from preprocessing |
| Longest Substring | Sliding window (variable) | Contiguous chunk + condition |
| Subarray Sum K | Prefix sum + hash map | Range queries → matching past values |

The pattern table on the back of this card is the entire reason for this module.
