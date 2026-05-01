# Problem 5 — Three Sum

> **Pattern:** Sort + two pointers (with duplicate skipping)
> **Difficulty:** Medium
> **Why it's #5:** This is the canonical "they sound similar but Three Sum is way harder than Two Sum" interview problem. The reason: dedup. Most candidates get the algorithm right and lose points on duplicate handling.

---

## Problem

Given an integer array `nums`, return all the triplets `[nums[i], nums[j], nums[k]]` such that `i != j`, `i != k`, `j != k`, and `nums[i] + nums[j] + nums[k] == 0`.

The solution set **must not contain duplicate triplets**.

### Examples

```
nums = [-1, 0, 1, 2, -1, -4]   →  [[-1, -1, 2], [-1, 0, 1]]
nums = [0, 1, 1]               →  []
nums = [0, 0, 0]               →  [[0, 0, 0]]
nums = [-2, 0, 0, 2, 2]        →  [[-2, 0, 2]]    (only one, even though there are two 2's)
```

### Constraints

- `3 <= len(nums) <= 3000`
- `-10^5 <= nums[i] <= 10^5`

---

## Think before coding

1. Brute force is three nested loops — O(n³). What's the optimal target? (Hint: if Two Sum was O(n) and this is "fix one element, then Two Sum the rest," what does that suggest?)
2. Two Sum on a sorted array can be done with two pointers (no hash map). Why does sorting help here even more — beyond just enabling two pointers?
3. The hardest part is **avoiding duplicate triplets**. What kinds of duplicates can arise? Think of three places where you might generate a duplicate.

---

## Hints

<details>
<summary>Hint 1 — high-level structure</summary>

Sort the array. Then for each index `i`, fix `nums[i]` as the first element of the triplet, and use two pointers (`left = i+1`, `right = n-1`) to find pairs that sum to `-nums[i]`.

That's O(n²) overall — outer loop is O(n), each inner two-pointer scan is O(n).

</details>

<details>
<summary>Hint 2 — three places duplicates can come from</summary>

1. **Outer loop:** if `nums[i] == nums[i-1]`, skip — same first element would generate the same triplets.
2. **After finding a valid triplet:** advance `left` past all duplicate values of `nums[left]`, and advance `right` past all duplicate values of `nums[right]`.
3. **Don't process i = 0 specially** — only skip when `i > 0`. Otherwise you'd skip the very first valid first-element.

</details>

<details>
<summary>Hint 3 — the inner loop, carefully</summary>

```
while left < right:
    s = nums[i] + nums[left] + nums[right]
    if s < 0:
        left += 1
    elif s > 0:
        right -= 1
    else:
        # found a triplet — record it
        result.append([nums[i], nums[left], nums[right]])

        # NOW skip duplicates on both sides
        while left < right and nums[left] == nums[left + 1]:
            left += 1
        while left < right and nums[right] == nums[right - 1]:
            right -= 1

        # finally, move past the matched pair
        left += 1
        right -= 1
```

</details>

<details>
<summary>Hint 4 — a small optimization</summary>

After sorting, if `nums[i] > 0`, you can stop the outer loop entirely. Why? Because all later elements are >= nums[i] > 0, so three positives can't sum to 0.

This isn't required for correctness, but interviewers love when you mention it.

</details>

---

## Your solution

Write your solution in `05-three-sum.py`.

After finishing:
1. Trace through `[-2, 0, 0, 2, 2]` and convince yourself you only emit `[-2, 0, 2]` once.
2. Confirm what happens on `[0, 0, 0]` — your code should return `[[0, 0, 0]]`, not crash.

---

## Follow-up questions

1. What about Four Sum? (Same pattern, one more outer loop, O(n³).)
2. **k-Sum?** Can you write a generic solver that takes `k` and `target`?
3. What if you only needed to *count* the number of valid triplets, not list them? Could you do better than O(n²)?
