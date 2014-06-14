package constraints

import (
	"testing"

	. "github.com/stretchr/testify/assert"

	"github.com/zimmski/tavor/test"
	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/primitives"
)

func TestOptionalTokensToBeTokens(t *testing.T) {
	var tok *token.Token

	Implements(t, tok, &Optional{})
}

func TestOptional(t *testing.T) {
	a := primitives.NewConstantInt(1)

	o := NewOptional(a)
	Equal(t, "1", o.String())
	Exactly(t, a, o.Get())

	r := test.NewRandTest(1)
	o.FuzzAll(r)
	Equal(t, "", o.String())

	o.FuzzAll(r)
	Equal(t, "1", o.String())

	o.Fuzz(r)
	Equal(t, "", o.String())

	o2 := o.Clone()
	Equal(t, o.String(), o2.String())

	Equal(t, 2, o.Permutations())
}

func TestOptionalOptionalTokenInterface(t *testing.T) {
	o := NewOptional(primitives.NewConstantInt(1))

	var optionalTok *token.OptionalToken

	Implements(t, optionalTok, o)

	Equal(t, "1", o.String())

	o.Deactivate()
	Equal(t, "", o.String())

	o.Activate()
	Equal(t, "1", o.String())
}
