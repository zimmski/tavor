package primitives

import (
	"github.com/zimmski/tavor/rand"
)

type ConstantString struct {
	value string
}

func NewConstantString(value string) *ConstantString {
	return &ConstantString{
		value: value,
	}
}

func (i *ConstantString) Fuzz(r rand.Rand) {
	// do nothing
}

func (i *ConstantString) String() string {
	return i.value
}
