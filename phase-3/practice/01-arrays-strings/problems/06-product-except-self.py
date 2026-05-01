# Problem 6 — Product of Array Except Self
# Time: O(?)
# Space: O(?) extra (output array doesn't count)

from typing import List


def product_except_self(nums: List[int]) -> List[int]:
    # TODO
    pass


if __name__ == "__main__":
    assert product_except_self([1, 2, 3, 4]) == [24, 12, 8, 6]
    assert product_except_self([-1, 1, 0, -3, 3]) == [0, 0, 9, 0, 0]
    assert product_except_self([2, 3]) == [3, 2]
    assert product_except_self([0, 0]) == [0, 0]
    print("All tests passed.")
