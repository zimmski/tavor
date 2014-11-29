package lists

import (
	"testing"

	. "github.com/zimmski/tavor/test/assert"

	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/primitives"
)

func TestAllTokensToBeTokens(t *testing.T) {
	var tok *token.ListToken

	Implements(t, tok, &All{})
}

func TestAll(t *testing.T) {
	a := primitives.NewConstantInt(10)
	b := primitives.NewConstantString("abc")

	o := NewAll(a, b)
	Equal(t, "10abc", o.String())
	Equal(t, 2, o.Len())
	Equal(t, 1, o.Permutations())
	Equal(t, 1, o.PermutationsAll())

	Nil(t, o.Permutation(1))
	Equal(t, "10abc", o.String())

	i, err := o.Get(0)
	Nil(t, err)
	Equal(t, a, i)
	i, err = o.Get(1)
	Nil(t, err)
	Equal(t, b, i)
	i, err = o.Get(2)
	Equal(t, err.(*ListError).Type, ListErrorOutOfBound)
	Nil(t, i)

	c := primitives.NewRangeInt(1, 2)
	o = NewAll(a, b, c)
	Equal(t, "10abc1", o.String())
	Equal(t, 3, o.Len())
	Equal(t, 1, o.Permutations())
	Equal(t, 2, o.PermutationsAll())

	Nil(t, o.Permutation(1))
	Equal(t, "10abc1", o.String())

	Nil(t, c.Permutation(2))
	Equal(t, "10abc2", o.String())
	Equal(t, 3, o.Len())

	o2 := o.Clone()
	Equal(t, o.String(), o2.String())
}
