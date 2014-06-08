package lists

import (
	"testing"

	. "github.com/stretchr/testify/assert"

	"github.com/zimmski/tavor/test"
	"github.com/zimmski/tavor/token/primitives"
)

func TestLeastTokensToBeTokens(t *testing.T) {
	var tok *List

	Implements(t, tok, &Least{})
}

func TestLeast(t *testing.T) {
	a := primitives.NewConstantString("a")

	o := NewLeast(a, 1)
	Equal(t, "a", o.String())
	Equal(t, 1, o.Len())

	i, err := o.Get(0)
	Nil(t, err)
	Equal(t, a, i)
	i, err = o.Get(1)
	Equal(t, err.(*ListError).Type, ListErrorOutOfBound)
	Nil(t, i)

	r := test.NewRandTest(1)
	o.Fuzz(r)
	Equal(t, "aaa", o.String())
	Equal(t, 3, o.Len())

	b := primitives.NewRangeInt(1, 3)
	o = NewLeast(b, 2)
	Equal(t, "11", o.String())
	Equal(t, 2, o.Len())

	r.Seed(2)
	o.Fuzz(r)
	Equal(t, "12312", o.String())
	Equal(t, 5, o.Len())

	o2 := o.Clone()
	Equal(t, o.String(), o2.String())
}
