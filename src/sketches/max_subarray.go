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
		if left > leftSum {
			break
		}
		leftIdx = idx
		left = leftSum
	}

	var right int = math.MinInt64
	var rightSum int = 0
	var rightIdx int

	for idx := mid + 1; idx <= high; idx++ {
		rightSum += arr[idx]
		if right > rightSum {
			break
		}
		rightIdx = idx
		right = rightSum
	}

	return ArrayRangeSum{low: leftIdx, high: rightIdx, sum: left + right}
}
