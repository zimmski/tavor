package logicals

import (
	"testing"

	. "github.com/stretchr/testify/assert"

	"github.com/zimmski/tavor/test"
	"github.com/zimmski/tavor/token/primitives"
)

func TestOr(t *testing.T) {
	a := primitives.NewConstantString("a")
	b := primitives.NewConstantString("b")

	o := NewOr(a, b)
	Equal(t, "ab", o.String())

	r := test.NewRandTest(0)
	o.Fuzz(r)
	Equal(t, "ab", o.String())

	r.Seed(100)
	o.Fuzz(r)
	Equal(t, "b", o.String())

	c := primitives.NewRangeInt(5, 10)
	o = NewOr(c)
	Equal(t, "5", o.String())

	o.Fuzz(r)
	Equal(t, "6", o.String())
}
