package lists

import (
	"testing"

	. "github.com/stretchr/testify/assert"

	"github.com/zimmski/tavor/test"
	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/primitives"
)

func TestRepeatTokensToBeTokens(t *testing.T) {
	var tok *token.Token

	a := primitives.NewRandomInt()

	Implements(t, tok, NewRepeat(a, 1, 5))
}

func TestRepeat(t *testing.T) {
	a := primitives.NewConstantString("a")

	o := NewRepeat(a, 5, 10)
	Equal(t, "aaaaa", o.String())

	r := test.NewRandTest(1)
	o.Fuzz(r)
	Equal(t, "aaaaaaa", o.String())

	b := primitives.NewRangeInt(1, 3)
	o = NewRepeat(b, 2, 10)
	Equal(t, "11", o.String())

	r.Seed(2)
	o.Fuzz(r)
	Equal(t, "12312", o.String())

	o2 := o.Clone()
	Equal(t, o.String(), o2.String())
}