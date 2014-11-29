package lists

import (
	"testing"

	. "github.com/zimmski/tavor/test/assert"

	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/primitives"
)

func TestOneTokensToBeTokens(t *testing.T) {
	var tok *token.ListToken

	Implements(t, tok, &One{})
}

func TestOne(t *testing.T) {
	a := primitives.NewConstantString("a")
	b := primitives.NewConstantString("b")

	o := NewOne(a, b)
	Equal(t, "a", o.String())
	Equal(t, 1, o.Len())
	Equal(t, 2, o.Permutations())
	Equal(t, 2, o.PermutationsAll())

	i, err := o.Get(0)
	Nil(t, err)
	Equal(t, a, i)
	i, err = o.Get(1)
	Equal(t, err.(*ListError).Type, ListErrorOutOfBound)
	Nil(t, i)

	Nil(t, o.Permutation(1))
	Equal(t, "a", o.String())
	Nil(t, o.Permutation(2))
	Equal(t, "b", o.String())

	Equal(t, o.Permutation(3).(*token.PermutationError).Type, token.PermutationErrorIndexOutOfBound)

	c := primitives.NewRangeInt(5, 10)
	o = NewOne(c)
	Equal(t, "5", o.String())
	Equal(t, 1, o.Len())
	Equal(t, 1, o.Permutations())
	Equal(t, 6, o.PermutationsAll())

	o2 := o.Clone()
	Equal(t, o.String(), o2.String())

	Nil(t, c.Permutation(2))
	Equal(t, "6", o.String())

	Equal(t, o.Permutation(2).(*token.PermutationError).Type, token.PermutationErrorIndexOutOfBound)
}
