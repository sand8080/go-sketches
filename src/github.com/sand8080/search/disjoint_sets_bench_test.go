package search

import (
	"math/rand"
	"testing"
	"time"
)

func random(min, max int) int {
	if min == max {
		return max
	}
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

func BenchmarkDisjointSetInt_Union(b *testing.B) {
	ds := NewDisjointSetInt(b.N)
	for i := 0; i < b.N; i++ {
		size := random(0, 4)
		ids := make([]int, size)
		for j := 0; j < size; j++ {
			ids[j] = random(1, b.N)
		}
		ds.Union(ids)
	}
}
