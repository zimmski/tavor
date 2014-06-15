package lists

import (
	"testing"

	. "github.com/stretchr/testify/assert"

	"github.com/zimmski/tavor/test"
	"github.com/zimmski/tavor/token"
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
	Equal(t, 6, o.PermutationsAll())

	Nil(t, o.Permutation(1))
	Equal(t, "", o.String())
	Nil(t, o.Permutation(2))
	Equal(t, "a", o.String())
	Nil(t, o.Permutation(3))
	Equal(t, "aa", o.String())

	Equal(t, o.Permutation(7).(*token.PermutationError).Type, token.PermutationErrorIndexOutOfBound)

	i, err := o.Get(0)
	Nil(t, err)
	Equal(t, a, i)
	i, err = o.Get(5)
	Equal(t, err.(*ListError).Type, ListErrorOutOfBound)
	Nil(t, i)

	r := test.NewRandTest(1)
	o.FuzzAll(r)
	Equal(t, "aa", o.String())
	Equal(t, 2, o.Len())

	b := primitives.NewRangeInt(1, 3)
	o = NewMost(b, 4)
	Equal(t, "1111", o.String())
	Equal(t, 4, o.Len())
	Equal(t, 5, o.Permutations())
	Equal(t, 13, o.PermutationsAll())

	Nil(t, o.Permutation(1))
	Equal(t, "", o.String())
	Nil(t, o.Permutation(2))
	Equal(t, "1", o.String())
	Nil(t, o.Permutation(3))
	Equal(t, "11", o.String())

	Equal(t, o.Permutation(6).(*token.PermutationError).Type, token.PermutationErrorIndexOutOfBound)

	r.Seed(2)
	o.FuzzAll(r)
	Equal(t, "123", o.String())
	Equal(t, 3, o.Len())

	o2 := o.Clone()
	Equal(t, o.String(), o2.String())
}
