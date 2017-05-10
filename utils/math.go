package utils

import (
	"math/rand"
	"time"
)

func InitRandom() {
	rand.Seed(time.Now().UnixNano())
}

func Random(min, max int) int {
	if min == max {
		return max
	}
	return rand.Intn(max-min) + min
}


