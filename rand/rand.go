package rand

// Rand defines random data generator for the Tavor framework
type Rand interface {
	// Int returns a non-negative pseudo-random int
	Int() int
	// Intn returns, as an int, a non-negative pseudo-random number in [0,n). It panics if n <= 0.
	Intn(n int) int
	// Int63 returns a non-negative pseudo-random 63-bit integer as an int64.
	Int63() int64
	// Int63n returns, as an int64, a non-negative pseudo-random number in [0,n). It panics if n <= 0.
	Int63n(n int64) int64
	// Seed uses the provided seed value to initialize the generator to a deterministic state.
	Seed(seed int64)
}
