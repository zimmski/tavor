package logicals

import (
	"testing"

	. "github.com/stretchr/testify/assert"

	"github.com/zimmski/tavor/test"
	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/primitives"
)

func TestManyTokensToBeTokens(t *testing.T) {
	var tok *token.Token

	a := primitives.NewRandomInt()

	Implements(t, tok, NewMany(a))
}

func TestMany(t *testing.T) {
	a := primitives.NewConstantString("a")
	b := primitives.NewConstantString("b")

	o := NewMany(a, b)
	Equal(t, "ab", o.String())

	r := test.NewRandTest(0)
	o.Fuzz(r)
	Equal(t, "ab", o.String())

	r.Seed(100)
	o.Fuzz(r)
	Equal(t, "b", o.String())

	c := primitives.NewRangeInt(5, 10)
	o = NewMany(c)
	Equal(t, "5", o.String())

	o.Fuzz(r)
	Equal(t, "6", o.String())
}
