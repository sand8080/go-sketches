package search

import (
	"testing"
	"github.com/sand8080/go-sketches/utils"
)

func BenchmarkDisjointSetInt_Union(b *testing.B) {
	ds := NewDisjointSetInt(b.N)
	for i := 0; i < b.N; i++ {
		size := utils.Random(0, 4)
		ids := make([]int, size)
		for j := 0; j < size; j++ {
			ids[j] = utils.Random(1, b.N)
		}
		ds.Union(ids)
	}
}

func BenchmarkDisjointSetInt_EmitGroups(b *testing.B) {
	ds := NewDisjointSetInt(b.N)
	for i := 0; i < b.N; i++ {
		size := utils.Random(0, 4)
		ids := make([]int, size)
		for j := 0; j < size; j++ {
			ids[j] = utils.Random(1, b.N)
		}
		ds.Union(ids)
	}
	for range ds.EmitGroups() {
	}
}
