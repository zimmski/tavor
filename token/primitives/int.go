package primitives

import (
	"strconv"

	"github.com/zimmski/tavor/rand"
)

type ConstantInt struct {
	value int
}

func NewConstantInt(value int) *ConstantInt {
	return &ConstantInt{
		value: value,
	}
}

func (i *ConstantInt) Fuzz(r rand.Rand) {
	// do nothing
}

func (i *ConstantInt) String() string {
	return strconv.Itoa(i.value)
}

type RandomInt struct {
	value int
}

func NewRandomInt() *RandomInt {
	return &RandomInt{
		value: 0,
	}
}

func (i *RandomInt) Fuzz(r rand.Rand) {
	i.value = r.Int()
}

func (i *RandomInt) String() string {
	return strconv.Itoa(i.value)
}

type RangeInt struct {
	from int
	to   int

	value int
}

func NewRangeInt(from, to int) *RangeInt {
	return &RangeInt{
		from:  from,
		to:    to,
		value: from,
	}
}

func (i *RangeInt) Fuzz(r rand.Rand) {
	ri := r.Intn(i.to - i.from + 1)

	i.value = i.from + ri
}

func (i *RangeInt) String() string {
	return strconv.Itoa(i.value)
}
