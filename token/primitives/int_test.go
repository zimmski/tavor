package primitives

import (
	"testing"

	. "github.com/stretchr/testify/assert"

	"github.com/zimmski/tavor/test"
	"github.com/zimmski/tavor/token"
)

func TestIntTokensToBeTokens(t *testing.T) {
	var tok *token.Token

	Implements(t, tok, &ConstantInt{})
	Implements(t, tok, &RandomInt{})
	Implements(t, tok, &RangeInt{})
}

func TestConstantInt(t *testing.T) {
	o := NewConstantInt(10)
	Equal(t, "10", o.String())

	r := test.NewRandTest(0)
	o.Fuzz(r)
	Equal(t, "10", o.String())

	o2 := o.Clone()
	Equal(t, o.String(), o2.String())
}

func TestRandomInt(t *testing.T) {
	o := NewRandomInt()
	Equal(t, "0", o.String())

	r := test.NewRandTest(0)
	o.Fuzz(r)
	Equal(t, "1", o.String())

	o2 := o.Clone()
	Equal(t, o.String(), o2.String())
}

func TestRangeInt(t *testing.T) {
	o := NewRangeInt(2, 4)
	Equal(t, "2", o.String())

	r := test.NewRandTest(0)
	o.Fuzz(r)
	Equal(t, "3", o.String())
	o.Fuzz(r)
	Equal(t, "4", o.String())
	o.Fuzz(r)
	Equal(t, "2", o.String())

	o2 := o.Clone()
	Equal(t, o.String(), o2.String())
}
