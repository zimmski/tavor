package lists

import (
	"testing"

	. "github.com/stretchr/testify/assert"

	"github.com/zimmski/tavor/test"
	"github.com/zimmski/tavor/token/primitives"
)

func TestMostTokensToBeTokens(t *testing.T) {
	var tok *List

	Implements(t, tok, &Most{})
}

func TestMost(t *testing.T) {
	a := primitives.NewConstantString("a")

	o := NewMost(a, 5)
	Equal(t, "aaaaa", o.String())
	Equal(t, 5, o.Len())
	Equal(t, 6, o.Permutations())

	i, err := o.Get(0)
	Nil(t, err)
	Equal(t, a, i)
	i, err = o.Get(1)
	Equal(t, err.(*ListError).Type, ListErrorOutOfBound)
	Nil(t, i)

	r := test.NewRandTest(1)
	o.Fuzz(r)
	Equal(t, "aa", o.String())
	Equal(t, 2, o.Len())

	b := primitives.NewRangeInt(1, 3)
	o = NewMost(b, 4)
	Equal(t, "1111", o.String())
	Equal(t, 4, o.Len())
	Equal(t, 13, o.Permutations())

	r.Seed(2)
	o.Fuzz(r)
	Equal(t, "123", o.String())
	Equal(t, 3, o.Len())

	o2 := o.Clone()
	Equal(t, o.String(), o2.String())
}
