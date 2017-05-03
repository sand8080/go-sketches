package search

import (
	"testing"
	"fmt"
	"reflect"
)

func TestNewDisjointSetInt_getHeaviestFailed(t *testing.T) {
	ds := NewDisjointSetInt(0)
	cases := []struct{
		input []int
		msg string
	}{
		{input: []int{}, msg: "Can't find minimal value in the empty collection"},
		{input: []int{1}, msg: "'1' is not found in weights"},
	}
	for _, c := range cases {
		_, _, err := ds.getHeaviest(c.input, c.input)
		if err == nil || err.Error() != c.msg{
			t.Errorf("getHeaviest should be failed with: '%v', have: '%v'", c.msg, err)
		}
	}
}

func TestNewDisjointSetInt_getHeaviestSuccess(t *testing.T) {
	cases := []struct {
		ids []int
		expected int
	}{
		{ids: []int{1, 2}, expected: 2},
		{ids: []int{2, 1}, expected: 2},
		{ids: []int{1, 2, 3}, expected: 3},
		{ids: []int{2, 3, 1}, expected: 3},
		{ids: []int{1, 3, 2}, expected: 3},
	}
	for _, c := range cases {
		ds := NewDisjointSetInt(len(c.ids))
		// Assume all ids as disjoint sets
		for _, id := range c.ids {
			ds.Weights[id] = 1
			ds.Parents[id] = id
		}

		roots, err := ds.getRoots(c.ids)
		if err != nil {
			t.Errorf("Getting roots for %v failed with: %v", c.ids, err)
		}

		max, _, err := ds.getHeaviest(c.ids, roots)
		if err != nil {
			t.Errorf("Calculation of heaviest for '%v' failed " +
				"with error: %v", c.ids, err)
		}
		if max != c.expected {
			t.Errorf("Expected '%v' as heaviest. Got: '%v'",
				c.expected, max)
		}
	}
}

func TestDisjointSetInt_Union(t *testing.T) {
	ds := NewDisjointSetInt(0)
	ds.Union([]int{1, 2, 3})
	ds.Union([]int{4, 5, 6})
	ds.Union([]int{1, 7})
	ds.Union([]int{7, 4})
	fmt.Println("ds", ds.Parents, ds.Weights)
	expected_parents := map[int]int{1: 3, 2: 3, 3: 3, 4: 6, 5: 6, 6: 3, 7: 3}
	if !reflect.DeepEqual(expected_parents, ds.Parents) {
		t.Errorf("Union unexpected result: %v != %v", expected_parents, ds.Parents)
	}
	expected_weights := map[int]int{1: 1, 2: 1, 3: 7, 4: 1, 5: 1, 6: 3, 7: 1}
	if !reflect.DeepEqual(expected_weights, ds.Weights) {
		t.Errorf("Union unexpected result: %v != %v", expected_parents, ds.Parents)
	}
}
