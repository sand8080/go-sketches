package utils

import "testing"

func TestRandomFailure(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Error("Panic should be initiated on wrong params " +
				"for random function")
		}
	}()
	Random(5, 1)
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
		result := Random(c.min, c.max)
		if result < c.min || (result >= c.max && c.min != c.max) {
			t.Errorf("Random value %d is out of bounds: '%d' - '%d'",
				result, c.min, c.max)
		}
	}
}
