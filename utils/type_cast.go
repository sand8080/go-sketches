package utils

import "strconv"

func IntsToStrings(vals []int) []string {
	result := make([]string, len(vals))
	for idx, val := range vals {
		result[idx] = strconv.Itoa(val)
	}
	return result
}
