# Problem 8 — Subarray Sum Equals K

> **Pattern:** Prefix sums + hash map (running sum trick)
> **Difficulty:** Medium
> **Why it's last:** This problem looks like sliding window but ISN'T — sliding window requires monotonicity (window only grows or shrinks predictably with the condition). Negative numbers break that. The prefix sum + hash map trick is one of the most powerful patterns in array problems, and it's the one that separates intermediate from advanced.

---

## Problem

Given an array of integers `nums` and an integer `k`, return the **total number of contiguous subarrays** whose sum equals `k`.

### Examples

```
nums = [1, 1, 1],     k = 2   →  2     (subarrays [1,1] starting at index 0, and [1,1] starting at index 1)
nums = [1, 2, 3],     k = 3   →  2     ([3] and [1,2])
nums = [1, -1, 1],    k = 0   →  1     ([1, -1])
nums = [1, 2, 1, 2, 1], k = 3 →  4
```

### Constraints

- `1 <= len(nums) <= 2 * 10^4`
- `-1000 <= nums[i] <= 1000`
- `-10^7 <= k <= 10^7`

**Note:** nums can contain negative numbers. This is why sliding window does NOT work here.

---

## Think before coding

1. Brute force: try every subarray, compute its sum, count matches. What's the complexity? (Be careful — there are two ways to do brute force, with different complexities.)
2. Why doesn't sliding window work here? Concrete example: `nums = [1, -1, 1, 1]`, `k = 1`. Walk through and convince yourself that "shrink window when sum > k" doesn't work.
3. Definition: `prefix[i]` = sum of nums[0..i-1]. Then the sum of `nums[i..j-1]` is `prefix[j] - prefix[i]`. **Reframe the problem in terms of prefix sums.**

---

## Hints

<details>
<summary>Hint 1 — restate the problem</summary>

The sum of subarray from index `i` to `j-1` (inclusive of i, exclusive of j) is `prefix[j] - prefix[i]`.

We want this to equal `k`:

```
prefix[j] - prefix[i] = k
prefix[i] = prefix[j] - k
```

So as we walk through and compute `prefix[j]`, we want to know: **how many earlier prefixes had value `prefix[j] - k`?** Each one corresponds to a valid subarray ending at index j-1.

</details>

<details>
<summary>Hint 2 — the data structure</summary>

A hash map from `prefix value → count of how many times this prefix value has occurred`. Initialize with `{0: 1}` to handle subarrays starting at index 0.

</details>

<details>
<summary>Hint 3 — the loop</summary>

```python
from collections import defaultdict

count = 0
running_sum = 0
seen = defaultdict(int)
seen[0] = 1   # the empty prefix has sum 0; this lets us count subarrays starting at index 0

for x in nums:
    running_sum += x
    count += seen[running_sum - k]
    seen[running_sum] += 1

return count
```

</details>

<details>
<summary>Hint 4 — why the {0: 1} initialization?</summary>

If `running_sum == k` at some point, that means the entire prefix from index 0 sums to k. We need `seen[0]` to be 1 in that case so we count it. Without the initialization, we'd miss every subarray that starts at index 0.

Trace through `nums=[1,1,1]`, `k=2` to convince yourself.

</details>

<details>
<summary>Hint 5 — order of operations matters</summary>

Inside the loop:
1. Compute `running_sum`.
2. Check `seen[running_sum - k]` and add to count.
3. THEN increment `seen[running_sum]`.

If you increment first, you might count `running_sum - k = running_sum` (i.e., k=0) using the very index you're at — which would correspond to an empty subarray. Wrong.

</details>

---

## Your solution

Write your solution in `08-subarray-sum-k.py`.

After finishing:
1. Trace through `[1, 1, 1], k=2` showing the contents of `seen` after each iteration.
2. Confirm: time O(n), space O(n) for the hash map.
3. Convince yourself this works with negative numbers using `[1, -1, 1], k=0`.

---

## Follow-up questions

1. What if you needed to return the **longest** subarray with sum k, instead of counting all of them? (Hint: store first-seen-index, not count.)
2. What if you needed the count of subarrays with sum *divisible* by k, not equal to k? (Same trick, but key the hash map by `running_sum % k`.)
3. What if all numbers were positive and you needed the count of subarrays with sum AT MOST k? (Sliding window now works! Why?)
