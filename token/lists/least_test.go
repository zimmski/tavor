package logicals

import (
	"testing"

	. "github.com/stretchr/testify/assert"

	"github.com/zimmski/tavor/test"
	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/primitives"
)

func TestLeastTokensToBeTokens(t *testing.T) {
	var tok *token.Token

	a := primitives.NewRandomInt()

	Implements(t, tok, NewLeast(a, 1))
}

func TestLeast(t *testing.T) {
	a := primitives.NewConstantString("a")

	o := NewLeast(a, 1)
	Equal(t, "a", o.String())

	r := test.NewRandTest(1)
	o.Fuzz(r)
	Equal(t, "aaa", o.String())

	b := primitives.NewRangeInt(1, 3)
	o = NewLeast(b, 2)
	Equal(t, "11", o.String())

	r.Seed(2)
	o.Fuzz(r)
	Equal(t, "12312", o.String())

	o2 := o.Clone()
	Equal(t, o.String(), o2.String())
}
