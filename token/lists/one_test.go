package lists

import (
	"testing"

	. "github.com/stretchr/testify/assert"

	"github.com/zimmski/tavor/test"
	"github.com/zimmski/tavor/token/primitives"
)

func TestOneTokensToBeTokens(t *testing.T) {
	var tok *List

	Implements(t, tok, &One{})
}

func TestOne(t *testing.T) {
	a := primitives.NewConstantString("a")
	b := primitives.NewConstantString("b")

	o := NewOne(a, b)
	Equal(t, "a", o.String())
	Equal(t, 1, o.Len())
	Equal(t, 2, o.Permutations())

	i, err := o.Get(0)
	Nil(t, err)
	Equal(t, a, i)
	i, err = o.Get(1)
	Equal(t, err.(*ListError).Type, ListErrorOutOfBound)
	Nil(t, i)

	r := test.NewRandTest(0)
	o.FuzzAll(r)
	Equal(t, "b", o.String())
	Equal(t, 1, o.Len())

	c := primitives.NewRangeInt(5, 10)
	o = NewOne(c)
	Equal(t, "5", o.String())
	Equal(t, 1, o.Len())
	Equal(t, 6, o.Permutations())

	o.FuzzAll(r)
	Equal(t, "6", o.String())
	Equal(t, 1, o.Len())

	o2 := o.Clone()
	Equal(t, o.String(), o2.String())
}
