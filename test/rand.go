package test

import (
	"github.com/zimmski/tavor/rand"
)

// NewRandTest returns a new instance of an increment random generator which is perfect for testing
func NewRandTest(seed int64) rand.Rand {
	return rand.NewIncrementRand(seed)
}
