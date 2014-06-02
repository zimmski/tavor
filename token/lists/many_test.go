package lists

import (
	"testing"

	. "github.com/stretchr/testify/assert"

	"github.com/zimmski/tavor/test"
	"github.com/zimmski/tavor/token/primitives"
)

func TestManyTokensToBeTokens(t *testing.T) {
	var tok *List

	Implements(t, tok, &Many{})
}

func TestMany(t *testing.T) {
	a := primitives.NewConstantString("a")
	b := primitives.NewConstantString("b")

	o := NewMany(a, b)
	Equal(t, "ab", o.String())
	Equal(t, 2, o.Len())

	r := test.NewRandTest(0)
	o.Fuzz(r)
	Equal(t, "ab", o.String())
	Equal(t, 2, o.Len())

	r.Seed(100)
	o.Fuzz(r)
	Equal(t, "b", o.String())
	Equal(t, 1, o.Len())

	c := primitives.NewRangeInt(5, 10)
	o = NewMany(c)
	Equal(t, "5", o.String())
	Equal(t, 1, o.Len())

	o.Fuzz(r)
	Equal(t, "6", o.String())
	Equal(t, 1, o.Len())

	o2 := o.Clone()
	Equal(t, o.String(), o2.String())
}
