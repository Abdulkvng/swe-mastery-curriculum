# Problem 5 — Three Sum
# Time: O(?)
# Space: O(?)

from typing import List


def three_sum(nums: List[int]) -> List[List[int]]:
    # TODO
    pass


def _normalize(triplets: List[List[int]]) -> List[tuple]:
    """Helper: sort within and across triplets so order doesn't break tests."""
    return sorted(tuple(sorted(t)) for t in triplets)


if __name__ == "__main__":
    assert _normalize(three_sum([-1, 0, 1, 2, -1, -4])) == _normalize([[-1, -1, 2], [-1, 0, 1]])
    assert _normalize(three_sum([0, 1, 1])) == []
    assert _normalize(three_sum([0, 0, 0])) == _normalize([[0, 0, 0]])
    assert _normalize(three_sum([-2, 0, 0, 2, 2])) == _normalize([[-2, 0, 2]])
    assert _normalize(three_sum([0, 0, 0, 0])) == _normalize([[0, 0, 0]])
    print("All tests passed.")
