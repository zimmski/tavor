package rand

import (
	"testing"

	. "github.com/zimmski/tavor/test/assert"
)

func TestIncrementRand(t *testing.T) {
	o := NewIncrementRand(1)

	Equal(t, 1, o.Int())
	Equal(t, 2, o.Int())
	Equal(t, 3, o.Int())

	o.Seed(2)
	Equal(t, 2, o.Int())
	Equal(t, 3, o.Int())
	Equal(t, 4, o.Int())
	Equal(t, 5, o.Int())

	o.Seed(0)
	Equal(t, 0, o.Intn(2))
	Equal(t, 1, o.Intn(2))
	Equal(t, 0, o.Intn(2))
	Equal(t, 1, o.Intn(2))
}
