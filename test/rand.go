package test

import (
	"math"
)

type RandTest struct {
	seed int64
}

func NewRandTest(seed int64) *RandTest {
	return &RandTest{
		seed: seed,
	}
}

func (r *RandTest) rand(n int64) int {
	if r.seed == math.MaxInt64 {
		r.seed = 0
	}

	if r.seed < n {
		r.seed++
	} else {
		r.seed = 0
	}

	return int(r.seed)
}

func (r *RandTest) Int() int {
	return r.rand(r.seed + 1)
}

func (r *RandTest) Intn(n int) int {
	return r.rand(int64(n - 1))
}

func (r *RandTest) Int63() int64 {
	return int64(r.Int())
}

func (r *RandTest) Int63n(n int64) int64 {
	return int64(r.Intn(int(n)))
}

func (r *RandTest) Seed(seed int64) {
	r.seed = seed
}
