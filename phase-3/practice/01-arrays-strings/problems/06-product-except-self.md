# Problem 6 — Product of Array Except Self

> **Pattern:** Prefix products + suffix products
> **Difficulty:** Medium
> **Why it's #6:** This problem looks like it requires division, but the constraint says you can't use it. The fix — pre-computing prefix and suffix products — is a thinking pattern that recurs constantly (range queries, sliding stats, etc.).

---

## Problem

Given an integer array `nums`, return an array `answer` such that `answer[i]` equals the product of all elements of `nums` **except** `nums[i]`.

**Constraints:**
- You **must NOT use division**.
- You should solve it in **O(n) time**.
- Bonus: solve it in O(1) extra space (output array doesn't count).

### Examples

```
nums = [1, 2, 3, 4]       →  [24, 12, 8, 6]
nums = [-1, 1, 0, -3, 3]  →  [0, 0, 9, 0, 0]
```

### Constraints

- `2 <= len(nums) <= 10^5`
- `-30 <= nums[i] <= 30`
- The product of any prefix or suffix is guaranteed to fit in a 32-bit integer.

---

## Think before coding

1. The "obvious" solution: compute the total product, then for each `i` divide by `nums[i]`. Why is this disallowed? (Hint: what about zeros?)
2. If you can't divide, you have to **build** each `answer[i]` from scratch. But a naive O(n²) approach would multiply (n-1) elements per position. How can you avoid redoing work?
3. For each position `i`, what are you really multiplying together?

---

## Hints

<details>
<summary>Hint 1 — split the product</summary>

For position `i`:

```
answer[i] = (product of everything LEFT of i) * (product of everything RIGHT of i)
```

If you precompute both, each answer is O(1).

</details>

<details>
<summary>Hint 2 — two passes, two arrays</summary>

```
left[i]  = nums[0] * nums[1] * ... * nums[i-1]    (1 if i == 0)
right[i] = nums[i+1] * ... * nums[n-1]            (1 if i == n-1)
answer[i] = left[i] * right[i]
```

Building each array takes O(n). Then combining them is O(n). Total O(n) time, O(n) space.

</details>

<details>
<summary>Hint 3 — collapsing to O(1) extra space</summary>

You don't need a separate `left` and `right` array. Use the answer array itself:

1. First pass (left to right): fill `answer[i]` with the product of everything left of `i`.
2. Second pass (right to left): keep a running `right_product` variable, multiply each `answer[i]` by it, then update `right_product *= nums[i]`.

The answer array is being reused; no other extra space.

</details>

<details>
<summary>Hint 4 — be careful about index 0 and index n-1</summary>

`answer[0]` should NOT include nums[0] in the left product. So the left product at index 0 is `1` (the empty product). Same for `answer[n-1]` and the right side.

</details>

---

## Your solution

Write your solution in `06-product-except-self.py`.

Implement the O(1)-extra-space version. After finishing:
1. Trace through `[1, 2, 3, 4]` showing the answer array after pass 1 and after pass 2.
2. Confirm that `[0, 4, 0]` works correctly — multiple zeros is a common test case to break naive solutions.

---

## Follow-up questions

1. Does this approach work if the input contains zeros? (Yes — but think about why: the division approach would fail on a zero.)
2. What if division WERE allowed? How many zeros change the answer?
3. What if the array were huge and stored across many machines? How would you parallelize this?
