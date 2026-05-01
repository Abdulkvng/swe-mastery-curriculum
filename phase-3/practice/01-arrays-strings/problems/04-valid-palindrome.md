# Problem 4 — Valid Palindrome

> **Pattern:** Two pointers (opposite ends moving inward)
> **Difficulty:** Easy
> **Why it's #4:** This is the cleanest possible introduction to two pointers. After this, you'll never need to ask "is two-pointer applicable here" — you'll just feel it.

---

## Problem

A phrase is a **palindrome** if, after converting all uppercase letters to lowercase and removing all non-alphanumeric characters, it reads the same forward and backward.

Given a string `s`, return `True` if it's a palindrome, `False` otherwise.

### Examples

```
"A man, a plan, a canal: Panama"   →  True
"race a car"                       →  False
" "                                →  True   (empty string after cleaning)
"0P"                               →  False  ('0' != 'p')
```

### Constraints

- `1 <= len(s) <= 2 * 10^5`
- `s` consists only of printable ASCII characters.

---

## Think before coding

1. The "easy" approach is: clean the string (remove non-alphanumeric, lowercase), then check `cleaned == cleaned[::-1]`. What's the time complexity? The space complexity?
2. Can you do it without building a cleaned copy? If so, what would your loop look like?
3. What edge cases could break a naive two-pointer solution? (Hint: what if the character at `left` is a comma?)

---

## Hints

<details>
<summary>Hint 1 — two pointers, opposite ends</summary>

Set `left = 0` and `right = len(s) - 1`. Walk them toward each other. If at any point `s[left] != s[right]`, return False. If they meet or cross without mismatch, return True.

</details>

<details>
<summary>Hint 2 — handling junk characters</summary>

When `s[left]` isn't alphanumeric, advance `left` without checking. Same for `right`. Use Python's `str.isalnum()`. Compare with `.lower()`.

</details>

<details>
<summary>Hint 3 — careful with the inner skip loops</summary>

Watch out:

```python
while left < right and not s[left].isalnum():
    left += 1
while left < right and not s[right].isalnum():
    right -= 1
```

The `left < right` check **inside** the skip loops is critical. Without it, your pointers can cross while skipping junk, and you'd compare garbage indices.

</details>

<details>
<summary>Hint 4 — comparing characters</summary>

Use `s[left].lower() != s[right].lower()` rather than lowercasing the whole string up front. Same correctness, less memory.

</details>

---

## Your solution

Write your solution in `04-valid-palindrome.py`.

Trace through `"A man, a plan, a canal: Panama"` step by step. Where do your pointers stop? What characters get compared?

---

## Follow-up questions

1. What if you were allowed to delete at most ONE character to make it a palindrome? (LeetCode 680. Surprisingly tricky — try to articulate the approach.)
2. What if you had to find the longest palindromic *substring* of s? (We'll handle this later — it uses "expand around center.")
3. How would you do this in a language without `.isalnum()` and `.lower()` built in? What's the underlying check?
