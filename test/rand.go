package test

import (
	"github.com/zimmski/tavor/rand"
)

func NewRandTest(seed int64) rand.Rand {
	return rand.NewIncrementRand(seed)
}
