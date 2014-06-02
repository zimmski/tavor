package sequences

import (
	"testing"

	. "github.com/stretchr/testify/assert"

	"github.com/zimmski/tavor/test"
	"github.com/zimmski/tavor/token"
)

func TestSequenceTokensToBeTokens(t *testing.T) {
	var tok *token.Token

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
	o.Fuzz(r)
	Equal(t, "12", o.String())
	Equal(t, 14, s.Next())

	o2 := o.Clone()
	Equal(t, o.String(), o2.String())
}

func TestExistingSequenceItem(t *testing.T) {
	s := NewSequence(10, 2)

	Equal(t, 10, s.Next())
	Equal(t, 12, s.Next())
	Equal(t, 14, s.Next())

	o := s.ExistingItem()
	Equal(t, "0", o.String())

	r := test.NewRandTest(0)
	o.Fuzz(r)
	Equal(t, "12", o.String())

	o.Fuzz(r)
	Equal(t, "14", o.String())

	o.Fuzz(r)
	Equal(t, "10", o.String())
}
