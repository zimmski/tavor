package rand

type Rand interface {
	Int() int
	Intn(n int) int
	Seed(seed int64)
}
