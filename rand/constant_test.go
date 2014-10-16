package rand

import (
	"testing"

	. "github.com/zimmski/tavor/test/assert"
)

func TestConstantRand(t *testing.T) {
	o := NewConstantRand(3)

	Equal(t, 3, o.Int())
	Equal(t, 3, o.Int())
	Equal(t, 3, o.Int())

	Equal(t, 1, o.Intn(1))
	Equal(t, 2, o.Intn(2))
	Equal(t, 3, o.Intn(3))

	o.Seed(2)
	Equal(t, 2, o.Int())
	Equal(t, 2, o.Int())
	Equal(t, 2, o.Int())

	Equal(t, 1, o.Intn(1))
	Equal(t, 2, o.Intn(2))
	Equal(t, 2, o.Intn(3))

	o.Seed(0)
	Equal(t, 0, o.Int())
	Equal(t, 0, o.Int())
	Equal(t, 0, o.Int())

	Equal(t, 0, o.Intn(1))
	Equal(t, 0, o.Intn(2))
	Equal(t, 0, o.Intn(3))

	o.Seed(4)
	Equal(t, 4, o.Int63())
	Equal(t, 4, o.Int63n(5))
}
