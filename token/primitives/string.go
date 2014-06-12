package primitives

import (
	"github.com/zimmski/tavor/rand"
	"github.com/zimmski/tavor/token"
)

type ConstantString struct {
	value string
}

func NewConstantString(value string) *ConstantString {
	return &ConstantString{
		value: value,
	}
}

func (p *ConstantString) Clone() token.Token {
	return &ConstantString{
		value: p.value,
	}
}

func (p *ConstantString) Fuzz(r rand.Rand) {
	// do nothing
}

func (p *ConstantString) Permutations() int {
	return 1
}

func (p *ConstantString) String() string {
	return p.value
}
