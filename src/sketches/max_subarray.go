package sketches

import (
	"math"
)

type ArrayRangeSum struct {
	low, high, sum int
}

func FindMaxCrossingArray(arr []int, low, mid, high int) ArrayRangeSum {

	var left int = math.MinInt64
	var leftSum int = 0
	var leftIdx int
	for idx := mid; idx >= low; idx-- {
		leftSum += arr[idx]
		if leftSum > left {
			leftIdx = idx
			left = leftSum
		}
	}

	var right int = math.MinInt64
	var rightSum int = 0
	var rightIdx int

	for idx := mid + 1; idx <= high; idx++ {
		rightSum += arr[idx]
		if rightSum > right {
			rightIdx = idx
			right = rightSum
		}
	}

	return ArrayRangeSum{low: leftIdx, high: rightIdx, sum: left + right}
}

func FindMaxSubArray(arr []int) ArrayRangeSum {
	if len(arr) > 0 {
		return FindMaxSubArrayRange(arr, 0, len(arr) - 1)
	} else {
		return ArrayRangeSum{}
	}
}

func FindMaxSubArrayRange(arr []int, low, high int) ArrayRangeSum {
	if low == high {
		return ArrayRangeSum{low, high, arr[low]}
	}

	mid := (high + low) / 2
	leftRange := FindMaxSubArrayRange(arr, low, mid)
	rightRange := FindMaxSubArrayRange(arr, mid + 1, high)
	crossRange := FindMaxCrossingArray(arr, low, mid, high)
	if leftRange.sum >= rightRange.sum && leftRange.sum >= crossRange.sum {
		return leftRange
	} else if rightRange.sum >= leftRange.sum && rightRange.sum >= crossRange.sum {
		return rightRange
	} else {
		return crossRange
	}
}