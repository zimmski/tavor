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

func TestConstantInt(t *testing.T) {
	o := NewOptional(primitives.NewConstantInt(1))
	Equal(t, "1", o.String())

	r := test.NewRandTest(1)
	o.Fuzz(r)
	Equal(t, "", o.String())

	o.Fuzz(r)
	Equal(t, "1", o.String())

	o2 := o.Clone()
	Equal(t, o.String(), o2.String())

	Equal(t, 2, o.Permutations())
}
