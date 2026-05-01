# Problem 1 — Two Sum

> **Pattern:** Hash map (scan once, remember things)
> **Difficulty:** Easy
> **Why it's first:** Every array interview starts here. If you can't articulate why a hash map turns this from O(n²) to O(n), nothing else will click.

---

## Problem

Given an array of integers `nums` and an integer `target`, return the indices of the two numbers that add up to `target`.

You may assume each input has exactly one solution, and you may not use the same element twice.

You can return the answer in any order.

### Examples

```
nums = [2, 7, 11, 15], target = 9      →  [0, 1]   (because 2 + 7 = 9)
nums = [3, 2, 4],     target = 6       →  [1, 2]   (because 2 + 4 = 6)
nums = [3, 3],        target = 6       →  [0, 1]
```

### Constraints

- `2 <= len(nums) <= 10^4`
- `-10^9 <= nums[i] <= 10^9`
- `-10^9 <= target <= 10^9`

---

## Think before coding

Before you read any hint, answer these:

1. What's the brute-force approach? What's its time complexity?
2. As you walk through the array left-to-right, what would you *want to know* at each position to find the answer in one pass?
3. What data structure lets you answer "have I seen value X before, and at what index?" in O(1)?

---

## Hints

<details>
<summary>Hint 1 — reframe the problem</summary>

Instead of asking "does there exist a pair (i, j) such that nums[i] + nums[j] = target," ask: **"as I look at nums[i], have I already seen the value (target - nums[i]) earlier in the array?"**

That reframe turns a pair-search into a single-element search, which is what hash maps are good at.

</details>

<details>
<summary>Hint 2 — what to store</summary>

You need to map **value → index** (not the other way around). When you see `nums[i]`, you compute `need = target - nums[i]` and check if `need` is in your map. If yes, you have your two indices.

</details>

<details>
<summary>Hint 3 — order of operations</summary>

Be careful: when you process `nums[i]`, do you check the map *before* or *after* inserting `nums[i]`?

Answer: check first, *then* insert. Otherwise if `target = 2 * nums[i]` and you only have one such element, you'll match yourself.

</details>

---

## Your solution

Write your solution in `01-two-sum.py`. When you're done:

1. Trace through it on the example `nums = [3, 3], target = 6` to make sure it handles duplicates.
2. State the time and space complexity in a comment at the top.
3. Then check `SOLUTIONS.md` for the canonical version.

---

## Follow-up questions to think about

(These are the kinds of things an Apple/Datadog interviewer would ask after you finish.)

1. What if the array were sorted? Could you do better than O(n) space?
2. What if there could be multiple valid pairs and you had to return all of them?
3. What if you needed to find three numbers that sum to target instead of two?
4. What if the input was a stream — numbers arriving one at a time and you had to answer "is there any pair seen so far that sums to target" after each insert?
