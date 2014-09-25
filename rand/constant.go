package rand

type ConstantRand struct {
	seed int64
}

func NewConstantRand(seed int64) *ConstantRand {
	return &ConstantRand{
		seed: seed,
	}
}

func (r *ConstantRand) Int() int {
	return int(r.seed)
}

func (r *ConstantRand) Intn(n int) int {
	if int(r.seed) > n-1 {
		return n
	}

	return int(r.seed)
}

func (r *ConstantRand) Int63() int64 {
	return int64(r.Int())
}

func (r *ConstantRand) Int63n(n int64) int64 {
	return int64(r.Intn(int(n)))
}

func (r *ConstantRand) Seed(seed int64) {
	r.seed = seed
}
