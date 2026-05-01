# Problem 3 — Contains Duplicate

> **Pattern:** Hash set (existence checking)
> **Difficulty:** Easy
> **Why it's #3:** Sets are hash maps when you only care about presence, not values. Cementing that distinction is worth one whole problem.

---

## Problem

Given an integer array `nums`, return `True` if any value appears at least twice, and `False` if every element is distinct.

### Examples

```
nums = [1, 2, 3, 1]              →  True
nums = [1, 2, 3, 4]              →  False
nums = [1, 1, 1, 3, 3, 4, 3, 2, 4, 2]  →  True
```

### Constraints

- `1 <= len(nums) <= 10^5`
- `-10^9 <= nums[i] <= 10^9`

---

## Think before coding

1. What's the simplest possible brute force? Two nested loops? What's the complexity?
2. What if you sorted first and then scanned? What does that cost?
3. Is there a way to do it in one pass?

There are **three legitimate approaches** here, with different time/space tradeoffs. Try to think of all three before reading hints.

---

## Hints

<details>
<summary>Hint 1 — three approaches</summary>

**Approach A:** Brute force, two nested loops. O(n²) time, O(1) space.
**Approach B:** Sort, then scan adjacent pairs. O(n log n) time, O(1) extra space (if sorting in place).
**Approach C:** Walk through, adding to a set. If you see one already in the set, return True. O(n) time, O(n) space.

The default expected answer is **C** unless the interviewer constrains memory.

</details>

<details>
<summary>Hint 2 — the Pythonic shortcut</summary>

```python
return len(set(nums)) != len(nums)
```

This works but is one line — not great for an interview. It also doesn't short-circuit. The explicit loop demonstrates more.

</details>

<details>
<summary>Hint 3 — when to mention each tradeoff</summary>

In an interview, code approach C, but verbalize: "I could also sort and check adjacent pairs if memory were tight. That's O(n log n) time, O(1) space, vs. O(n) time, O(n) space for the hash set approach."

This shows you think about tradeoffs without being asked.

</details>

---

## Your solution

Write your solution in `03-contains-duplicate.py`. Implement approach C (hash set), but include a comment block listing the alternatives and their tradeoffs.

---

## Follow-up questions

1. What if you needed to find duplicates within `k` indices of each other? (i.e., `i - j` <= k where nums[i] == nums[j])
2. What if you had to return all duplicates, not just check existence?
3. What if the input were a stream and you had to answer "does this new element duplicate any past element" after each insert?
