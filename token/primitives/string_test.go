package primitives

import (
	"testing"

	. "github.com/zimmski/tavor/test/assert"

	"github.com/zimmski/tavor/token"
)

func TestStringTokensToBeTokens(t *testing.T) {
	var tok *token.Token

	Implements(t, tok, &ConstantString{})
}

func TestConstantString(t *testing.T) {
	o := NewConstantString("abc")
	Equal(t, "abc", o.String())

	Equal(t, 1, o.Permutations())

	Nil(t, o.Permutation(1))
	Equal(t, "abc", o.String())

	Equal(t, o.Permutation(2).(*token.PermutationError).Type, token.PermutationErrorIndexOutOfBound)

	o2 := o.Clone()
	Equal(t, o.String(), o2.String())
}
