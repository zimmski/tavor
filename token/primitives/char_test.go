package primitives

import (
	"testing"

	. "github.com/zimmski/tavor/test/assert"

	"github.com/zimmski/tavor/token"
)

func TestCharTokensToBeTokens(t *testing.T) {
	var tok *token.Token

	Implements(t, tok, &CharacterClass{})
}

func TestCharacterClass(t *testing.T) {
	o := NewCharacterClass("abc")
	Equal(t, "a", o.String())

	Equal(t, 3, o.Permutations())

	Nil(t, o.Permutation(1))
	Equal(t, "a", o.String())

	Nil(t, o.Permutation(2))
	Equal(t, "b", o.String())

	Nil(t, o.Permutation(3))
	Equal(t, "c", o.String())

	Equal(t, o.Permutation(4).(*token.PermutationError).Type, token.PermutationErrorIndexOutOfBound)

	o = NewCharacterClass(`\d`)
	Equal(t, "0", o.String())
	Equal(t, 10, o.Permutations())

	o = NewCharacterClass(`1-9`)
	Equal(t, "1", o.String())
	Equal(t, 9, o.Permutations())

	Nil(t, o.Permutation(2))
	Equal(t, "2", o.String())

	Nil(t, o.Permutation(9))
	Equal(t, "9", o.String())

	Equal(t, o.Permutation(10).(*token.PermutationError).Type, token.PermutationErrorIndexOutOfBound)

	o = NewCharacterClass(`a1-9`)
	Equal(t, "a", o.String())
	Equal(t, 10, o.Permutations())

	Nil(t, o.Permutation(2))
	Equal(t, "1", o.String())

	Nil(t, o.Permutation(10))
	Equal(t, "9", o.String())

	Equal(t, o.Permutation(11).(*token.PermutationError).Type, token.PermutationErrorIndexOutOfBound)

	o = NewCharacterClass(`a-fA-F0-9`)
	Equal(t, "a", o.String())
	Equal(t, 22, o.Permutations())

	Nil(t, o.Permutation(1))
	Equal(t, "a", o.String())

	Nil(t, o.Permutation(2))
	Equal(t, "b", o.String())

	Nil(t, o.Permutation(7))
	Equal(t, "A", o.String())

	Nil(t, o.Permutation(8))
	Equal(t, "B", o.String())

	Nil(t, o.Permutation(13))
	Equal(t, "0", o.String())

	Nil(t, o.Permutation(14))
	Equal(t, "1", o.String())

	Equal(t, o.Permutation(23).(*token.PermutationError).Type, token.PermutationErrorIndexOutOfBound)

	o2 := o.Clone()
	Equal(t, o.String(), o2.String())
}
