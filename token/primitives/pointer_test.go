package primitives

import (
	"github.com/zimmski/tavor/token"
	"testing"

	. "github.com/stretchr/testify/assert"

	"github.com/zimmski/tavor/test"
)

func TestPointerTokensToBeTokens(t *testing.T) {
	var tok *token.Token

	Implements(t, tok, &Pointer{})
}

func TestPointer(t *testing.T) {
	a := NewRangeInt(4, 10)

	o := NewPointer(a)
	Equal(t, "4", o.String())

	r := test.NewRandTest(0)
	o.Fuzz(r)
	// this uses a clone
	Equal(t, "5", o.String())
	// this is the original one which must be untouched
	Equal(t, "4", a.String())

	o2 := o.Clone()
	Equal(t, o.String(), o2.String())

	o = NewEmptyPointer()
	Nil(t, o.Tok)
}
