package sketches

import (
	"testing"
)

func TestFindMaxCrossingArray(t *testing.T) {
	var tests = []struct {
		arr []int
		low, mid, high int
		expected ArrayRangeSum
	}{
		{
			arr: []int {1, 2, 3, 4, 5},
			low: 0, mid: 2, high: 4,
			expected: ArrayRangeSum{low: 0, high: 4, sum: 15},
		},
		{
			arr: []int {1, 2, 0, 4, 5},
			low: 0, mid: 2, high: 4,
			expected: ArrayRangeSum{low: 0, high: 4, sum: 12},
		},
		{
			arr: []int {1, 2, 0, 4, 5},
			low: 0, mid: 2, high: 4,
			expected: ArrayRangeSum{low: 0, high: 4, sum: 12},
		},
		{
			arr: []int {1, 2, 3, 4},
			low: 0, mid: 2, high: 3,
			expected: ArrayRangeSum{low: 0, high: 3, sum: 10},
		},
		{
			arr: []int {-1, 2, -3, 4},
			low: 0, mid: 2, high: 3,
			expected: ArrayRangeSum{low: 1, high: 3, sum: 3},
		},
		{
			arr: []int {1, 2, -3, 4},
			low: 0, mid: 2, high: 3,
			expected: ArrayRangeSum{low: 0, high: 3, sum: 4},
		},
		{
			arr: []int {1, 2, -3, 4, -3},
			low: 0, mid: 2, high: 4,
			expected: ArrayRangeSum{low: 0, high: 3, sum: 4},
		},
	}

	for _, c := range tests {
		actual := FindMaxCrossingArray(c.arr, c.low, c.mid, c.high)
		if c.expected != actual {
			t.Errorf("Expected: %v, actual: %v", c.expected, actual)
		}
	}
}

