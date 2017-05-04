package search

import (
	"reflect"
	"testing"
)

func TestFindMaxCrossingArray(t *testing.T) {
	var tests = []struct {
		arr            []int
		low, mid, high int
		expected       ArrayRangeSum
	}{
		{
			arr: []int{1, 2, 3, 4, 5},
			low: 0, mid: 2, high: 4,
			expected: ArrayRangeSum{low: 0, high: 4, sum: 15},
		},
		{
			arr: []int{1, 2, 0, 4, 5},
			low: 0, mid: 2, high: 4,
			expected: ArrayRangeSum{low: 0, high: 4, sum: 12},
		},
		{
			arr: []int{1, 2, 0, 4, 5},
			low: 0, mid: 2, high: 4,
			expected: ArrayRangeSum{low: 0, high: 4, sum: 12},
		},
		{
			arr: []int{1, 2, 3, 4},
			low: 0, mid: 2, high: 3,
			expected: ArrayRangeSum{low: 0, high: 3, sum: 10},
		},
		{
			arr: []int{-1, 2, -3, 4},
			low: 0, mid: 2, high: 3,
			expected: ArrayRangeSum{low: 1, high: 3, sum: 3},
		},
		{
			arr: []int{1, 2, -3, 4},
			low: 0, mid: 2, high: 3,
			expected: ArrayRangeSum{low: 0, high: 3, sum: 4},
		},
		{
			arr: []int{1, 2, -3, 4, -3},
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

func TestFindMaxSubarray(t *testing.T) {
	var tests = []struct {
		input    []int
		expected ArrayRangeSum
	}{
		// empty array
		{
			input:    []int{},
			expected: ArrayRangeSum{low: 0, high: 0, sum: 0},
		},
		// one element array
		{
			input:    []int{1},
			expected: ArrayRangeSum{low: 0, high: 0, sum: 1},
		},
		{
			input:    []int{-1},
			expected: ArrayRangeSum{low: 0, high: 0, sum: -1},
		},
		{
			input:    []int{10},
			expected: ArrayRangeSum{low: 0, high: 0, sum: 10},
		},
		// positive elements
		{
			input:    []int{1, 10, 2},
			expected: ArrayRangeSum{low: 0, high: 2, sum: 13},
		},
		{
			input:    []int{1, 10, 2, 0},
			expected: ArrayRangeSum{low: 0, high: 2, sum: 13},
		},
		{
			input:    []int{1, 10, 2, 0, 0, 4},
			expected: ArrayRangeSum{low: 0, high: 5, sum: 17},
		},
		{
			input:    []int{0, 1, 10, 2},
			expected: ArrayRangeSum{low: 1, high: 3, sum: 13},
		},
		{
			input:    []int{5, 0, 1, 10, 2},
			expected: ArrayRangeSum{low: 0, high: 4, sum: 18},
		},
		// negative elements
		{
			input:    []int{-1, -10, -2},
			expected: ArrayRangeSum{low: 0, high: 0, sum: -1},
		},
		{
			input:    []int{-1, -10, 0, -2},
			expected: ArrayRangeSum{low: 2, high: 2, sum: 0},
		},
		{
			input:    []int{0, -1, -10, -2},
			expected: ArrayRangeSum{low: 0, high: 0, sum: 0},
		},
		{
			input:    []int{-1, -10, -2, 0},
			expected: ArrayRangeSum{low: 3, high: 3, sum: 0},
		},
	}

	for _, c := range tests {
		actual := FindMaxSubArray(c.input)
		if c.expected != actual {
			t.Errorf("Expected: %v, actual: %v",
				c.expected, actual)
		}
	}

}

func TestFindMaxSubArrayReferenceValue(t *testing.T) {
	var arr = []int{
		13, -3, -25, 20, -3, -16, -23, 18, 20, -7, 12, -5,
		-22, 15, -4, 7,
	}
	expectedRange := ArrayRangeSum{low: 7, high: 10, sum: 43}
	actualRange := FindMaxSubArray(arr)
	if expectedRange != actualRange {
		t.Errorf("Expected range: %v, actual: %v",
			expectedRange, actualRange)
	}
	expectedArray := []int{18, 20, -7, 12}
	actualArray := arr[actualRange.low : actualRange.high+1]
	if !reflect.DeepEqual(expectedArray, actualArray) {
		t.Errorf("Expected array: %v, actual: %v",
			expectedArray, actualArray)
	}
}
