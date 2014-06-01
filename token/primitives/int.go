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

func (p *ConstantInt) Clone() token.Token {
	return &ConstantInt{
		value: p.value,
	}
}

func (p *ConstantInt) Fuzz(r rand.Rand) {
	// do nothing
}

func (p *ConstantInt) String() string {
	return strconv.Itoa(p.value)
}

type RandomInt struct {
	value int
}

func NewRandomInt() *RandomInt {
	return &RandomInt{
		value: 0,
	}
}

func (p *RandomInt) Clone() token.Token {
	return &RandomInt{
		value: p.value,
	}
}

func (p *RandomInt) Fuzz(r rand.Rand) {
	p.value = r.Int()
}

func (p *RandomInt) String() string {
	return strconv.Itoa(p.value)
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

func (p *RangeInt) Clone() token.Token {
	return &RangeInt{
		from:  p.from,
		to:    p.to,
		value: p.value,
	}
}

func (p *RangeInt) Fuzz(r rand.Rand) {
	ri := r.Intn(p.to - p.from + 1)

	p.value = p.from + ri
}

func (p *RangeInt) String() string {
	return strconv.Itoa(p.value)
}
