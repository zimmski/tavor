package rand

// ConstantRand implements a constant random generator.
// This generator returns the defined seed every time, unless a method has a bound then it returns at most the given bound. It is therefore not random at all and can be safely used when deterministic values are needed.
type ConstantRand struct {
	seed int64
}

// NewConstantRand returns a new instance of the constant random generator
func NewConstantRand(seed int64) *ConstantRand {
	return &ConstantRand{
		seed: seed,
	}
}

// Int returns a non-negative pseudo-random int
func (r *ConstantRand) Int() int {
	return int(r.seed)
}

// Intn returns, as an int, a non-negative pseudo-random number in [0,n). It panics if n <= 0.
func (r *ConstantRand) Intn(n int) int {
	if n <= 0 {
		panic("invalid argument to Intn")
	}

	if int(r.seed) > n-1 {
		return n
	}

	return int(r.seed)
}

// Int63 returns a non-negative pseudo-random 63-bit integer as an int64.
func (r *ConstantRand) Int63() int64 {
	return int64(r.Int())
}

// Int63n returns, as an int64, a non-negative pseudo-random number in [0,n). It panics if n <= 0.
func (r *ConstantRand) Int63n(n int64) int64 {
	if n <= 0 {
		panic("invalid argument to Int63n")
	}

	return int64(r.Intn(int(n)))
}

// Seed uses the provided seed value to initialize the generator to a deterministic state.
func (r *ConstantRand) Seed(seed int64) {
	r.seed = seed
}
