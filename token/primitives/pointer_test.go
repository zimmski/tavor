package primitives

import (
	"testing"

	. "github.com/zimmski/tavor/test/assert"

	"github.com/zimmski/tavor/test"
	"github.com/zimmski/tavor/token"
)

func TestPointerTokensToBeTokens(t *testing.T) {
	var tok *token.Token

	Implements(t, tok, &Pointer{})

	var forward *token.ForwardToken

	Implements(t, forward, &Pointer{})
}

func TestPointer(t *testing.T) {
	a := NewRangeInt(4, 10)

	o := NewPointer(a)
	Equal(t, "4", o.String())
	Equal(t, 1, o.Permutations())
	Equal(t, 7, o.PermutationsAll())

	r := test.NewRandTest(1)
	o.FuzzAll(r)
	// this uses a clone
	Equal(t, "5", o.String())
	// this is the original one which must be untouched
	Equal(t, "4", a.String())

	o2 := o.Clone()

	// cloned pointers are always different to their original one

	o.FuzzAll(r)
	o2.FuzzAll(r)

	// original token still untouched
	Equal(t, "4", a.String())
	// first cloned token
	Equal(t, "6", o.String())
	// second cloned token
	Equal(t, "7", o2.String())

	Nil(t, o.Permutation(1))
	Equal(t, "6", o.String())

	Equal(t, o.Permutation(8).(*token.PermutationError).Type, token.PermutationErrorIndexOutOfBound)
}
