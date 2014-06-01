package lists

import (
	"testing"

	. "github.com/stretchr/testify/assert"

	"github.com/zimmski/tavor/test"
	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/primitives"
)

func TestMostTokensToBeTokens(t *testing.T) {
	var tok *token.Token

	a := primitives.NewRandomInt()

	Implements(t, tok, NewMost(a, 1))
}

func TestMost(t *testing.T) {
	a := primitives.NewConstantString("a")

	o := NewMost(a, 5)
	Equal(t, "aaaaa", o.String())

	r := test.NewRandTest(1)
	o.Fuzz(r)
	Equal(t, "aa", o.String())

	b := primitives.NewRangeInt(1, 3)
	o = NewMost(b, 4)
	Equal(t, "1111", o.String())

	r.Seed(2)
	o.Fuzz(r)
	Equal(t, "123", o.String())

	o2 := o.Clone()
	Equal(t, o.String(), o2.String())
}
