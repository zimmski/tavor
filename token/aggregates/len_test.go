package aggregates

import (
	"testing"

	. "github.com/zimmski/tavor/test/assert"

	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/lists"
	"github.com/zimmski/tavor/token/primitives"
)

func TestLenTokensToBeTokens(t *testing.T) {
	var tok *token.Token

	Implements(t, tok, &Len{})
}

func TestConstantInt(t *testing.T) {
	list := lists.NewRepeat(primitives.NewConstantInt(1), 1, 10)
	Equal(t, "1", list.String())

	o := NewLen(list)
	Equal(t, "1", o.String())
	Equal(t, 1, o.Permutations())
	Equal(t, 1, o.Permutations())

	Nil(t, list.Permutation(1))
	Equal(t, "1", o.String())

	Nil(t, list.Permutation(2))
	Equal(t, "11", list.String())
	Equal(t, "2", o.String())

	Nil(t, list.Permutation(3))
	Equal(t, "111", list.String())
	Equal(t, "3", o.String())

	Nil(t, list.Permutation(4))
	Equal(t, "1111", list.String())
	Equal(t, "4", o.String())

	Nil(t, list.Permutation(5))
	Equal(t, "11111", list.String())
	Equal(t, "5", o.String())

	o2 := o.Clone()
	Equal(t, o.String(), o2.String())

	Equal(t, 1, o.Permutations())

	Nil(t, o.Permutation(1))
	Equal(t, "5", o.String())

	Equal(t, o.Permutation(2).(*token.PermutationError).Type, token.PermutationErrorIndexOutOfBound)
}
