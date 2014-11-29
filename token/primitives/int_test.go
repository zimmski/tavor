package primitives

import (
	"testing"

	. "github.com/zimmski/tavor/test/assert"

	"github.com/zimmski/tavor/token"
)

func TestIntTokensToBeTokens(t *testing.T) {
	var tok *token.Token

	Implements(t, tok, &ConstantInt{})
	Implements(t, tok, &RangeInt{})
}

func TestConstantInt(t *testing.T) {
	o := NewConstantInt(10)
	Equal(t, "10", o.String())

	Equal(t, 1, o.Permutations())

	Nil(t, o.Permutation(1))
	Equal(t, "10", o.String())

	Equal(t, o.Permutation(2).(*token.PermutationError).Type, token.PermutationErrorIndexOutOfBound)

	o2 := o.Clone()
	Equal(t, o.String(), o2.String())
}

func TestRangeInt(t *testing.T) {
	o := NewRangeInt(2, 4)
	Equal(t, "2", o.String())

	Equal(t, 3, o.Permutations())

	Nil(t, o.Permutation(1))
	Equal(t, "2", o.String())
	Nil(t, o.Permutation(2))
	Equal(t, "3", o.String())

	Equal(t, o.Permutation(4).(*token.PermutationError).Type, token.PermutationErrorIndexOutOfBound)

	o2 := o.Clone()
	Equal(t, o.String(), o2.String())

	// range with step 2
	o = NewRangeIntWithStep(2, 6, 2)
	Equal(t, "2", o.String())
}
