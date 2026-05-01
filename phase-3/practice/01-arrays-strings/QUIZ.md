# Module 1 — Arrays & Strings: Quiz

> **Rules:** Answer out loud or write in a scratch file. No code. No peeking at INTUITION.md while answering. If you don't know, write "I don't know" — that's more useful than guessing.
>
> When you're done, check answers at the bottom. Score yourself. Re-attempt anything you missed in 48 hours.

---

## Section A — Mechanics (5 questions)

**A1.** What is the time complexity of `nums.append(x)` in Python? Why does the answer include the word "amortized"?

**A2.** Why is `nums.insert(0, x)` slow but `nums.append(x)` fast?

**A3.** You write `sub = nums[2:7]`. What is the time complexity of this line, and what is the space complexity?

**A4.** In Python, can you do `s[2] = 'x'` where `s` is a string? Why or why not? What's the workaround?

**A5.** What's the time complexity of `''.join(words)` if there are `n` words with average length `m`? (Be careful — the answer is not what most people guess.)

---

## Section B — Pattern recognition (5 questions)

For each problem statement below, name the pattern (hash map, two pointers, sliding window, prefix sum, sort-then-scan) and explain in one sentence why.

**B1.** "Given a sorted array, find two numbers that sum to a target."

**B2.** "Given an array of integers, return the count of subarrays whose sum equals k."

**B3.** "Given a string, find the length of the longest substring with no repeating characters."

**B4.** "Given an array, return whether any number appears more than once."

**B5.** "Given a list of meeting time intervals, determine if a person can attend all meetings."

---

## Section C — Trap-spotting (5 questions)

**C1.** A candidate writes this code to "remove the first element of a list 1 million times":

```python
for _ in range(1_000_000):
    nums.pop(0)
```

What's wrong with this? What's the real complexity? What should they use instead?

**C2.** A candidate writes:

```python
result = ""
for word in words:
    result += word
```

If `words` has 10,000 strings of average length 50, what's the complexity, and why?

**C3.** A candidate solves Two Sum like this:

```python
def two_sum(nums, target):
    for i in range(len(nums)):
        for j in range(i + 1, len(nums)):
            if nums[i] + nums[j] == target:
                return [i, j]
```

State the complexity, then describe the optimal solution in one sentence (no code).

**C4.** What is the difference between a **subarray** and a **subsequence**? Which one does sliding window help with?

**C5.** This loop is "supposed" to be O(n) but is actually slower. Why?

```python
for i in range(len(nums)):
    chunk = nums[i:i+5]
    process(chunk)
```

---

## Section D — The "explain it" questions (3 questions)

**D1.** Explain in plain English: why does the two-pointer technique work on a sorted array but not on an unsorted one?

**D2.** Explain in plain English: what does a prefix sum array let you do in O(1) that would otherwise be O(n)?

**D3.** An interviewer asks: "What's the time complexity of `x in nums`?" You correctly say "O(n)." They follow up: "What if `nums` were a `set` instead?" What do you say, and what is the catch you should mention?

---

## Section E — Open-ended (2 questions)

**E1.** Without coding it, describe step-by-step how you would find the **longest substring without repeating characters** in a string. Pretend you're walking through it on a whiteboard. Use the word "window" somewhere in your answer.

**E2.** Why is `' '.join(parts)` faster than building the string with `+=`? Describe what's happening in memory.

---

---

## Answers

> **Don't read until you've attempted every question.**

<details>
<summary>Section A answers</summary>

**A1.** O(1) amortized. "Amortized" means averaged over many operations. Most appends are truly O(1) — Python over-allocates the underlying array, so there's spare room. Occasionally the array fills up and Python has to allocate a bigger block and copy everything (O(n)). But that copy happens rarely enough that the *average* cost per append stays constant.

**A2.** `insert(0, x)` has to shift every existing element one slot to the right to make room at index 0 — that's O(n) work. `append(x)` just writes to the next empty slot at the end — O(1).

**A3.** Time: O(k) where k is the slice length (5 elements here). Space: O(k) — Python creates a new list. This is different from Go slices, which share memory with the original.

**A4.** No — strings are immutable in Python. To "modify" a character, convert to a list (`list(s)`), modify the list, then `''.join(...)` it back. Or build a new string by concatenation.

**A5.** O(n × m) — the total length of all the words combined. Common wrong answer is "O(n)." `join` allocates one big string of the right size, then copies each word into it. The copying is what makes it linear in total length, not just count.

</details>

<details>
<summary>Section B answers</summary>

**B1.** **Two pointers.** The array is sorted, so you can start with one pointer at each end and move them inward based on whether the sum is too big or too small.

**B2.** **Prefix sum** (with a hash map of seen prefix sums). The classic "subarray sum equals k" — count how often `currentPrefix - k` has appeared before.

**B3.** **Sliding window.** Contiguous, longest, condition on the chunk. Textbook variable-size window.

**B4.** **Hash set.** Just track what you've seen. (You could sort first for O(1) extra space, but hash set is the default answer.)

**B5.** **Sort first, then scan.** Sort by start time, then check if any meeting overlaps the previous one.

</details>

<details>
<summary>Section C answers</summary>

**C1.** `pop(0)` is O(n) because it shifts every remaining element left by one. Calling it a million times is O(n²) overall — for a million-element list that's 10¹² operations, completely unusable. Use `collections.deque`, which has O(1) `popleft()`.

**C2.** O(n × m) per concatenation, totaling O(n² × m). Strings are immutable, so `result += word` allocates a brand new string and copies the entire previous result every iteration. With 10,000 words this is roughly 50 million wasted character-copies. Fix: append to a list, then `''.join(list)` once at the end — O(n × m).

**C3.** Brute force is O(n²) time, O(1) space. The optimal is O(n) time, O(n) space — scan once, and as you visit each number, check if `target - num` is already in a hash map of values you've seen.

**C4.** A **subarray** is a contiguous chunk (e.g., `[3, 4, 5]` from `[1, 2, 3, 4, 5, 6]`). A **subsequence** preserves order but skips allowed (e.g., `[1, 3, 5]`). Sliding window only helps with subarrays. Subsequences usually need DP.

**C5.** `nums[i:i+5]` creates a new list of 5 elements every iteration — that's O(5) work, but more importantly, it allocates memory n times. The loop is technically O(n), but the constant factor is much worse than people expect, and it triggers garbage collection. Better: index directly with `nums[i], nums[i+1], ...` or pass indices into `process`.

</details>

<details>
<summary>Section D answers</summary>

**D1.** On a sorted array, the relationship between two values at positions `i` and `j` is monotonic: if `nums[i] + nums[j]` is too small, you *know* you need a bigger value, which can only come from increasing `i` (since `j` is already at the upper end). The sortedness gives you a direction to move. On an unsorted array there's no signal — moving a pointer doesn't tell you anything about whether you're getting closer to the answer.

**D2.** A prefix sum array lets you compute the sum of any contiguous range in O(1). Without it, summing `nums[i..j]` requires walking through every element in that range — O(j - i). After O(n) preprocessing, every range query is constant.

**D3.** "It becomes O(1) on average." The catch: the worst case is technically O(n) (if every key collides into one bucket), but in practice with Python's hash function this never happens. Also worth mentioning: sets only work with hashable values — you can't put a list in a set.

</details>

<details>
<summary>Section E answers</summary>

**E1.** Maintain a **window** with two indices, `left` and `right`, both starting at 0. Use a hash map to track each character's most recent index. Walk `right` forward through the string. For each new character: if it's already in the window (its last index is `>= left`), move `left` to one past that previous index — this shrinks the window so the duplicate is no longer inside. Update the character's recent index. Track the max window size as you go. Each character enters the window once and leaves at most once → O(n).

**E2.** Strings are immutable, so `result += word` can't append in place — it allocates a brand new string of length `len(result) + len(word)` and copies both. After 10,000 concatenations you've copied roughly the same characters 10,000 times. `' '.join(parts)` walks the list once to compute total length, allocates one final string of exactly that size, then copies each piece in — every character is copied exactly once.

</details>

---

## Scoring

- **18–20 correct:** Solid foundation. Move to the problems.
- **14–17 correct:** Reasonable. Re-read the relevant INTUITION.md sections for what you missed, then move on.
- **10–13 correct:** Re-read INTUITION.md fully. Re-attempt the quiz in 24 hours. Don't move to problems yet — they'll feel harder than they should.
- **Below 10:** Walk through INTUITION.md slowly with me. Talk it out. The problems will be frustrating without this foundation.
