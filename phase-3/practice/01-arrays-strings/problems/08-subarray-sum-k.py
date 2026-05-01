# Problem 8 — Subarray Sum Equals K
# Time: O(?)
# Space: O(?)

from typing import List


def subarray_sum(nums: List[int], k: int) -> int:
    # TODO
    pass


if __name__ == "__main__":
    assert subarray_sum([1, 1, 1], 2) == 2
    assert subarray_sum([1, 2, 3], 3) == 2
    assert subarray_sum([1, -1, 1], 0) == 1
    assert subarray_sum([1, 2, 1, 2, 1], 3) == 4
    assert subarray_sum([1], 0) == 0
    assert subarray_sum([0, 0, 0, 0], 0) == 10  # 4 single + 3 pair + 2 triple + 1 quad = 10
    print("All tests passed.")
