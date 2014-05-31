package primitives

import (
	"strconv"

	"github.com/zimmski/tavor/rand"
	"github.com/zimmski/tavor/token"
)

type ConstantInt struct {
	value int
}

func NewConstantInt(value int) *ConstantInt {
	return &ConstantInt{
		value: value,
	}
}

func (i *ConstantInt) Clone() token.Token {
	return &ConstantInt{
		value: i.value,
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

func (i *RandomInt) Clone() token.Token {
	return &RandomInt{
		value: i.value,
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

func (i *RangeInt) Clone() token.Token {
	return &RangeInt{
		from:  i.from,
		to:    i.to,
		value: i.value,
	}
}

func (i *RangeInt) Fuzz(r rand.Rand) {
	ri := r.Intn(i.to - i.from + 1)

	i.value = i.from + ri
}

func (i *RangeInt) String() string {
	return strconv.Itoa(i.value)
}
