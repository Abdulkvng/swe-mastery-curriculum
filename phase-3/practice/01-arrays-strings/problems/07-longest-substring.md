# Problem 7 — Longest Substring Without Repeating Characters

> **Pattern:** Variable-size sliding window
> **Difficulty:** Medium
> **Why it's #7:** This is THE sliding window problem. Once you understand this one, every other sliding window problem is a variation.

---

## Problem

Given a string `s`, find the length of the longest substring without repeating characters.

### Examples

```
"abcabcbb"   →  3   ("abc")
"bbbbb"      →  1   ("b")
"pwwkew"     →  3   ("wke" — note: "pwke" is a subsequence, not a substring)
""           →  0
"dvdf"       →  3   ("vdf")
```

### Constraints

- `0 <= len(s) <= 5 * 10^4`
- `s` consists of English letters, digits, symbols, and spaces.

---

## Think before coding

This is the most important "think before coding" of the module. Spend real time here.

1. The brute force is: for each start position `i`, walk forward until you hit a repeat, track the max. What's the complexity? (It's O(n²) for the outer loop, but checking "is this a repeat" is another factor — what makes it potentially O(n³)?)
2. The key insight of sliding window: once you know `s[i..j]` is valid, you might not need to start over from `i+1`. Why?
3. When you encounter a repeat, where should the new window start? (Hint: think about *which* repeat you hit, not just that you hit one.)

---

## Hints

<details>
<summary>Hint 1 — the window invariant</summary>

Maintain a window `[left, right]` that is **always valid** (no repeats). Walk `right` forward one step at a time. When adding `s[right]` would break the invariant, move `left` forward until it's valid again.

```
left = 0
for right in range(len(s)):
    while s[right] is already in window:
        remove s[left]; left += 1
    add s[right] to window
    update answer
```

</details>

<details>
<summary>Hint 2 — what data structure tracks "is this in the window"?</summary>

A `set` works for "is char in window." But for the smarter version, use a `dict` mapping `char → most recent index seen`. This lets you jump `left` directly past the previous occurrence in one step instead of stepping it forward one character at a time.

</details>

<details>
<summary>Hint 3 — the index-jumping version</summary>

```python
last_seen = {}
left = 0
best = 0

for right, ch in enumerate(s):
    if ch in last_seen and last_seen[ch] >= left:
        left = last_seen[ch] + 1   # jump past the previous occurrence
    last_seen[ch] = right
    best = max(best, right - left + 1)

return best
```

The `last_seen[ch] >= left` check is critical. Why? Consider `"abba"`:
- right=2, s[right]='b', last_seen['b']=1, left becomes 2
- right=3, s[right]='a', last_seen['a']=0... but 0 < left (= 2), so the 'a' at index 0 isn't actually in the window anymore. Without this check, we'd move `left` backward, breaking everything.

</details>

<details>
<summary>Hint 4 — why this is O(n)</summary>

`right` advances n times. `left` only advances forward, never back, total at most n times. Total work: O(n). This is the sliding window guarantee — each element enters and leaves the window once.

</details>

---

## Your solution

Write your solution in `07-longest-substring.py`.

Implement the index-jumping version. Trace through `"dvdf"`:
- right=0 ('d'): window = "d", best = 1
- right=1 ('v'): window = "dv", best = 2
- right=2 ('d'): jump left to 1, window = "vd", best = 2
- right=3 ('f'): window = "vdf", best = 3

Confirm your code does exactly this.

---

## Follow-up questions

1. What if you were allowed at most K repeating characters? (Lookup: "Longest Substring with At Most K Distinct Characters.")
2. What if the string contained Unicode? Does your solution still work? (It should — Python handles this natively, but the underlying point is hash maps key by value, not by byte.)
3. What if you needed to RETURN the substring itself, not just its length?
