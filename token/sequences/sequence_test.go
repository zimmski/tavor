package sequences

import (
	"testing"

	. "github.com/zimmski/tavor/test/assert"

	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/lists"
	"github.com/zimmski/tavor/token/primitives"
)

func TestSequenceTokensToBeTokens(t *testing.T) {
	var tok *token.Token

	Implements(t, tok, &Sequence{})
	Implements(t, tok, &SequenceItem{})
	Implements(t, tok, &SequenceExistingItem{})
}

func TestSequence(t *testing.T) {
	o := NewSequence(10, 2)

	Equal(t, 10, o.Next())
	Equal(t, 12, o.Next())
	Equal(t, 14, o.Next())

	o.Reset()

	Equal(t, 10, o.Next())
	Equal(t, 12, o.Next())
	Equal(t, 14, o.Next())
}

func TestSequenceItem(t *testing.T) {
	s := NewSequence(10, 2)

	o := s.Item()
	Equal(t, "10", o.String())

	Nil(t, o.Permutation(1))
	Equal(t, "12", o.String())
	Equal(t, 14, s.Next())

	o2 := o.Clone()
	Equal(t, o.String(), o2.String())

	Equal(t, 1, o.Permutations())

	o.Reset()
	Equal(t, "16", o.String())
}

func TestExistingSequenceItem(t *testing.T) {
	s := NewSequence(10, 2)

	o := s.ExistingItem(nil)
	Equal(t, "-1", o.String())

	Equal(t, 10, s.Next())
	Equal(t, 12, s.Next())
	Equal(t, 14, s.Next())

	o = s.ExistingItem(nil)
	Equal(t, "10", o.String())
	Equal(t, 3, o.Permutations())

	Nil(t, o.Permutation(2))
	Equal(t, "12", o.String())

	Nil(t, o.Permutation(3))
	Equal(t, "14", o.String())

	Nil(t, o.Permutation(1))
	Equal(t, "10", o.String())

	// Except
	o = s.ExistingItem([]token.Token{primitives.NewConstantInt(10)})
	Equal(t, "10", o.String())
	Nil(t, o.Permutation(1))
	Equal(t, "12", o.String())

	o = s.ExistingItem([]token.Token{primitives.NewConstantInt(14)})
	Equal(t, "10", o.String())
	Nil(t, o.Permutation(1))
	Equal(t, "10", o.String())

	// Except with list token
	{
		a := lists.NewAll(primitives.NewConstantInt(10), primitives.NewConstantInt(12))
		o = s.ExistingItem([]token.Token{a})
		Equal(t, "10", o.String())
		Nil(t, o.Permutation(1))
		Equal(t, "14", o.String())
	}
}

func TestResetSequenceItem(t *testing.T) {
	s := NewSequence(10, 2)

	Equal(t, 10, s.Next())
	Equal(t, 12, s.Next())
	Equal(t, 14, s.Next())

	o := s.ResetItem()

	Nil(t, o.Permutation(1))

	Equal(t, 10, s.Next())
	Equal(t, 12, s.Next())
	Equal(t, 14, s.Next())

	Equal(t, 1, o.Permutations())
}
