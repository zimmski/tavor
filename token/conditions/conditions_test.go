package conditions

import (
	"testing"

	. "github.com/zimmski/tavor/test/assert"

	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/lists"
	"github.com/zimmski/tavor/token/primitives"
)

func TestConditionsTokensToBeTokens(t *testing.T) {
	var tok *token.Token

	Implements(t, tok, &If{})
}

func TestVariableIf(t *testing.T) {
	var o token.Token = NewIf(IfPair{
		Head: NewBooleanEqual(primitives.NewConstantInt(1), primitives.NewConstantInt(1)),
		Body: primitives.NewConstantString("a"),
	})
	Equal(t, "a", o.String())
	Equal(t, 1, o.Permutations())
	Equal(t, 1, o.PermutationsAll())

	o2 := o.Clone()
	Equal(t, o.String(), o2.String())

	Equal(t, 1, o.Permutations())

	Nil(t, o.Permutation(1))
	Equal(t, "a", o.String())

	Equal(t, o.Permutation(2).(*token.PermutationError).Type, token.PermutationErrorIndexOutOfBound)

	o = lists.NewAll(
		NewIf(IfPair{
			Head: NewBooleanEqual(primitives.NewConstantInt(1), primitives.NewConstantInt(2)),
			Body: primitives.NewConstantString("a"),
		}),
		primitives.NewConstantString("b"),
	)
	Equal(t, "b", o.String())
}

func TestVariableElse(t *testing.T) {
	o := NewIf(
		IfPair{
			Head: NewBooleanEqual(primitives.NewConstantInt(1), primitives.NewConstantInt(2)),
			Body: primitives.NewConstantString("a"),
		},
		IfPair{
			Head: NewBooleanTrue(),
			Body: primitives.NewConstantString("b"),
		},
	)
	Equal(t, "b", o.String())
	Equal(t, 1, o.Permutations())
	Equal(t, 1, o.PermutationsAll())

	o2 := o.Clone()
	Equal(t, o.String(), o2.String())

	Equal(t, 1, o.Permutations())

	Nil(t, o.Permutation(1))
	Equal(t, "b", o.String())

	Equal(t, o.Permutation(2).(*token.PermutationError).Type, token.PermutationErrorIndexOutOfBound)
}
