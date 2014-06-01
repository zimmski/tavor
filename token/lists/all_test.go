package lists

import (
	"testing"

	. "github.com/stretchr/testify/assert"

	"github.com/zimmski/tavor/test"
	"github.com/zimmski/tavor/token/primitives"
)

func TestAllTokensToBeTokens(t *testing.T) {
	var tok *List

	a := primitives.NewRandomInt()

	Implements(t, tok, NewAll(a))
}

func TestAll(t *testing.T) {
	a := primitives.NewConstantInt(10)
	b := primitives.NewConstantString("abc")

	o := NewAll(a, b)
	Equal(t, "10abc", o.String())
	Equal(t, 2, o.Len())

	r := test.NewRandTest(0)
	o.Fuzz(r)
	Equal(t, "10abc", o.String())
	Equal(t, 2, o.Len())

	c := primitives.NewRangeInt(1, 2)
	o = NewAll(a, b, c)
	Equal(t, "10abc1", o.String())
	Equal(t, 3, o.Len())

	o.Fuzz(r)
	Equal(t, "10abc2", o.String())
	Equal(t, 3, o.Len())

	o2 := o.Clone()
	Equal(t, o.String(), o2.String())
}
