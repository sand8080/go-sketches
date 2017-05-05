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

func TestRandomFailure(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Error("Panic should be initiated on wrong params " +
				"for random function")
		}
	}()
	random(5, 1)
}

func TestRandom(t *testing.T) {
	cases := []struct {
		min, max int
	}{
		{min: 1, max: 1},
		{min: 0, max: 5},
		{min: -5, max: 0},
		{min: 10, max: 20},
		{min: -10, max: 10},
	}
	for _, c := range cases {
		result := random(c.min, c.max)
		if result < c.min || (result >= c.max && c.min != c.max) {
			t.Errorf("Random value %d is out of bounds: '%d' - '%d'",
				result, c.min, c.max)
		}
	}
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

func BenchmarkDisjointSetInt_EmitGroups(b *testing.B) {
	ds := NewDisjointSetInt(b.N)
	for i := 0; i < b.N; i++ {
		size := random(0, 4)
		ids := make([]int, size)
		for j := 0; j < size; j++ {
			ids[j] = random(1, b.N)
		}
		ds.Union(ids)
	}
	for range ds.EmitGroups() {
	}
}
