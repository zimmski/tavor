package rand

import (
	"testing"

	. "github.com/zimmski/tavor/test/assert"
)

func TestRandTest(t *testing.T) {
	o := NewIncrementRand(0)

	Equal(t, 0, o.Int())
	Equal(t, 0, o.Int())
	Equal(t, 0, o.Int())

	o.Seed(2)
	Equal(t, 0, o.Int())
	Equal(t, 1, o.Int())
	Equal(t, 0, o.Int())
	Equal(t, 1, o.Int())

	o.Seed(0)
	Equal(t, 0, o.Intn(2))
	Equal(t, 1, o.Intn(2))
	Equal(t, 0, o.Intn(2))
	Equal(t, 1, o.Intn(2))
}
