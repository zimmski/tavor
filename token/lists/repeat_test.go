package lists

import (
	"testing"

	. "github.com/stretchr/testify/assert"

	"github.com/zimmski/tavor/test"
	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/primitives"
)

func TestRepeatTokensToBeTokens(t *testing.T) {
	var tok *List

	Implements(t, tok, &Repeat{})
}

func TestRepeat(t *testing.T) {
	a := primitives.NewConstantString("a")

	o := NewRepeat(a, 5, 10)
	Equal(t, "aaaaa", o.String())
	Equal(t, 5, o.Len())
	Equal(t, 6, o.Permutations())
	Equal(t, 6, o.PermutationsAll())

	i, err := o.Get(0)
	Nil(t, err)
	Equal(t, a, i)
	i, err = o.Get(6)
	Equal(t, err.(*ListError).Type, ListErrorOutOfBound)
	Nil(t, i)

	Nil(t, o.Permutation(1))
	Equal(t, "aaaaa", o.String())
	Nil(t, o.Permutation(2))
	Equal(t, "aaaaaa", o.String())
	Nil(t, o.Permutation(3))
	Equal(t, "aaaaaaa", o.String())

	Equal(t, o.Permutation(7).(*token.PermutationError).Type, token.PermutationErrorIndexOutOfBound)

	r := test.NewRandTest(1)
	o.FuzzAll(r)
	Equal(t, "aaaaaaa", o.String())
	Equal(t, 7, o.Len())

	o = NewRepeat(primitives.NewRangeInt(1, 2), 0, 2)
	Equal(t, "", o.String())
	Equal(t, 0, o.Len())
	Equal(t, 3, o.Permutations())
	Equal(t, 7, o.PermutationsAll())

	o = NewRepeat(primitives.NewRangeInt(1, 2), 1, 2)
	Equal(t, "1", o.String())
	Equal(t, 1, o.Len())
	Equal(t, 2, o.Permutations())
	Equal(t, 6, o.PermutationsAll())

	o = NewRepeat(primitives.NewRangeInt(1, 2), 0, 3)
	Equal(t, "", o.String())
	Equal(t, 0, o.Len())
	Equal(t, 4, o.Permutations())
	Equal(t, 15, o.PermutationsAll())

	o = NewRepeat(primitives.NewRangeInt(1, 2), 1, 3)
	Equal(t, "1", o.String())
	Equal(t, 1, o.Len())
	Equal(t, 3, o.Permutations())
	Equal(t, 14, o.PermutationsAll())

	b := primitives.NewRangeInt(1, 3)
	o = NewRepeat(b, 2, 10)
	Equal(t, "11", o.String())
	Equal(t, 2, o.Len())
	Equal(t, 9, o.Permutations())
	Equal(t, 29523, o.PermutationsAll())

	Nil(t, o.Permutation(1))
	Equal(t, "11", o.String())
	Nil(t, o.Permutation(2))
	Equal(t, "111", o.String())
	Nil(t, o.Permutation(3))
	Equal(t, "1111", o.String())

	Equal(t, o.Permutation(10).(*token.PermutationError).Type, token.PermutationErrorIndexOutOfBound)

	r.Seed(2)
	o.FuzzAll(r)
	Equal(t, "12312", o.String())
	Equal(t, 5, o.Len())

	o2 := o.Clone()
	Equal(t, o.String(), o2.String())
}
