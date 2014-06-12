package sequences

import (
	"testing"

	. "github.com/stretchr/testify/assert"

	"github.com/zimmski/tavor/test"
	"github.com/zimmski/tavor/token"
)

func TestSequenceTokensToBeTokens(t *testing.T) {
	var tok *token.Token

	Implements(t, tok, &Sequence{})
	Implements(t, tok, &sequenceItem{})
	Implements(t, tok, &sequenceExistingItem{})
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

	r := test.NewRandTest(0)
	o.FuzzAll(r)
	Equal(t, "12", o.String())
	Equal(t, 14, s.Next())

	o2 := o.Clone()
	Equal(t, o.String(), o2.String())

	Equal(t, 1, o.Permutations())
}

func TestExistingSequenceItem(t *testing.T) {
	s := NewSequence(10, 2)

	o := s.ExistingItem()
	Equal(t, "-1", o.String())

	Equal(t, 10, s.Next())
	Equal(t, 12, s.Next())
	Equal(t, 14, s.Next())

	o = s.ExistingItem()
	Equal(t, "10", o.String())

	r := test.NewRandTest(0)
	o.FuzzAll(r)
	Equal(t, "12", o.String())

	o.FuzzAll(r)
	Equal(t, "14", o.String())

	o.FuzzAll(r)
	Equal(t, "10", o.String())

	Equal(t, 1, o.Permutations())
}

func TestResetSequenceItem(t *testing.T) {
	s := NewSequence(10, 2)

	Equal(t, 10, s.Next())
	Equal(t, 12, s.Next())
	Equal(t, 14, s.Next())

	r := test.NewRandTest(0)
	o := s.ResetItem()
	o.FuzzAll(r)

	Equal(t, 10, s.Next())
	Equal(t, 12, s.Next())
	Equal(t, 14, s.Next())

	Equal(t, 1, o.Permutations())
}
