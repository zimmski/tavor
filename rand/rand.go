package rand

type Rand interface {
	Int() int
	Intn(n int) int
	Int63() int64
	Int63n(n int64) int64
	Seed(seed int64)
}
