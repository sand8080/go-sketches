package utils

import (
	"testing"
	"reflect"
)

func TestIntsToStrings(t *testing.T) {
	cases := []struct{
		input []int
		expected []string
	}{
		{
			input: []int{},
			expected: []string{},
		},
		{
			input: []int{0, 1, 2},
			expected: []string{"0", "1", "2"},
		},
		{
			input: []int{-2, -1, 0, 1, 2},
			expected: []string{"-2", "-1", "0", "1", "2"},
		},
	}
	for _, c := range cases {
		actual := IntsToStrings(c.input)
		if !reflect.DeepEqual(c.expected, actual) {
			t.Errorf("Ints to strings conversion failed. %v != %v\n",
				c.expected, actual)
		}
	}
}

func BenchmarkIntsToStrings(b *testing.B) {
	ints_num := 100
	ints := make([]int, ints_num)
	for i := 0; i < ints_num; i++ {
		ints[i] = Random(-1000, 1000)
	}
	for i := 0; i < b.N; i++ {
		IntsToStrings(ints)
	}

}