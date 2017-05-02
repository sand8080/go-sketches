package search

import (
	"testing"
	"fmt"
)

func TestDisjointSetInt_Union(t *testing.T) {
	ds := NewDisjointSetInt(0)
	ds.Union([]int{1, 2, 3})
	//fmt.Println("ds", ds)
}

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
		_, _, err := ds.getHeaviest(c.input)
		if err.Error() != c.msg{
			t.Error(fmt.Sprintf("getHeaviest should be failed with: %q, have: %q",
				c.msg, err))
		}
	}
}

func TestNewDisjointSetInt_getHeaviestSuccess(t *testing.T) {
	ds := NewDisjointSetInt(0)
	ds.Weights[1] = 1
	ds.Weights[2] = 1
	ds.Parents[1] = 1
	ds.Parents[2] = 2
	max, max_idx, err := ds.getHeaviest(ds.getRoots([]int{2, 1}))
	fmt.Println("#####", max, max_idx, err)
	//cases := []struct{
	//	input []int
	//	msg string
	//}{
	//	{input: []int{}, msg: "Can't find minimal value in the empty collection"},
	//	{input: []int{1}, msg: "'1' is not found in weights"},
	//}
	//for _, c := range cases {
	//	_, _, err := ds.getHeaviest(c.input)
	//	if err.Error() != c.msg{
	//		t.Error(fmt.Sprintf("getHeaviest should be failed with: %q, have: %q",
	//			c.msg, err))
	//	}
	//}
}