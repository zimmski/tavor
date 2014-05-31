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

func (s *ConstantString) Clone() token.Token {
	return &ConstantString{
		value: s.value,
	}
}

func (s *ConstantString) Fuzz(r rand.Rand) {
	// do nothing
}

func (s *ConstantString) String() string {
	return s.value
}
