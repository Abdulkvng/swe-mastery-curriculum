# Problem 3 — Contains Duplicate
# Time: O(?)
# Space: O(?)
#
# Tradeoffs to consider:
#   Approach A (brute force):     O(n²) time, O(1) space
#   Approach B (sort + scan):     O(n log n) time, O(1) extra space
#   Approach C (hash set):        O(n) time, O(n) space      <-- default
#
# Implement approach C below.

from typing import List


def contains_duplicate(nums: List[int]) -> bool:
    # TODO
    pass


if __name__ == "__main__":
    assert contains_duplicate([1, 2, 3, 1]) == True
    assert contains_duplicate([1, 2, 3, 4]) == False
    assert contains_duplicate([1, 1, 1, 3, 3, 4, 3, 2, 4, 2]) == True
    assert contains_duplicate([1]) == False
    print("All tests passed.")
