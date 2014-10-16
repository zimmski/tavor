package test

import (
	"testing"

	. "github.com/zimmski/tavor/test/assert"
)

func TestRandTest(t *testing.T) {
	o := NewRandTest(0)

	Equal(t, 0, o.Int())
	Equal(t, 1, o.Int())
	Equal(t, 2, o.Int())

	o.Seed(0)
	Equal(t, 0, o.Int())
	Equal(t, 1, o.Int())

	Equal(t, 2, o.Intn(10))
	Equal(t, 0, o.Intn(1))
	Equal(t, 1, o.Intn(2))
	Equal(t, 0, o.Intn(2))
	Equal(t, 1, o.Intn(2))
}
