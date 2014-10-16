package rand

// IncrementRand implements an incrementing random generator.
// This generator returns beginning from 0 an incremented value on every call. If the defined seed is equal or greater to the current generator value, the generator resets its value to 0. The generator will reset itself also to zero, if a method has a bound and the current generator value is equal or greater to the bound.
type IncrementRand struct {
	seed  int64
	value int64
}

// NewIncrementRand returns a new instance of the increment random generator
func NewIncrementRand(seed int64) *IncrementRand {
	return &IncrementRand{
		seed:  seed,
		value: seed,
	}
}

// Int returns a non-negative pseudo-random int
func (r *IncrementRand) Int() int {
	return r.Intn(int(r.value + 1))
}

// Intn returns, as an int, a non-negative pseudo-random number in [0,n). It panics if n <= 0.
func (r *IncrementRand) Intn(n int) int {
	if n <= 0 {
		panic("invalid argument to Intn")
	}

	if n <= int(r.value) {
		r.value = 0
	}

	v := r.value

	r.value++

	return int(v)
}

// Int63 returns a non-negative pseudo-random 63-bit integer as an int64.
func (r *IncrementRand) Int63() int64 {
	return int64(r.Int())
}

// Int63n returns, as an int64, a non-negative pseudo-random number in [0,n). It panics if n <= 0.
func (r *IncrementRand) Int63n(n int64) int64 {
	if n <= 0 {
		panic("invalid argument to Int63n")
	}

	return int64(r.Intn(int(n)))
}

// Seed uses the provided seed value to initialize the generator to a deterministic state.
func (r *IncrementRand) Seed(seed int64) {
	r.seed = seed
	r.value = seed
}
