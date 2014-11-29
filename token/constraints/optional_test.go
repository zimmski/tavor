package constraints

import (
	"testing"

	. "github.com/zimmski/tavor/test/assert"

	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/primitives"
)

func TestOptionalTokensToBeTokens(t *testing.T) {
	var tok *token.Token

	Implements(t, tok, &Optional{})

	var forward *token.ForwardToken

	Implements(t, forward, &Optional{})
}

func TestOptional(t *testing.T) {
	a := primitives.NewConstantInt(1)

	o := NewOptional(a)
	Equal(t, "1", o.String())
	True(t, Exactly(t, a, o.Get()))
	Equal(t, 2, o.Permutations())
	Equal(t, 2, o.PermutationsAll())

	Nil(t, o.Permutation(1))
	Equal(t, "", o.String())
	Nil(t, o.Permutation(2))
	Equal(t, "1", o.String())

	Equal(t, o.Permutation(3).(*token.PermutationError).Type, token.PermutationErrorIndexOutOfBound)

	o2 := o.Clone()
	Equal(t, o.String(), o2.String())
}

func TestOptionalOptionalTokenInterface(t *testing.T) {
	a := primitives.NewConstantInt(1)

	o := NewOptional(a)

	var optionalTok *token.OptionalToken

	Implements(t, optionalTok, o)

	Equal(t, "1", o.String())

	o.Deactivate()
	Nil(t, o.Get())
	Equal(t, "", o.String())

	o.Activate()
	Equal(t, o.Get(), a)
	Equal(t, "1", o.String())
}
