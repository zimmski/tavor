package variables

import (
	"testing"

	. "github.com/stretchr/testify/assert"

	"github.com/zimmski/tavor/test"
	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/primitives"
)

func TestVariablesTokensToBeTokens(t *testing.T) {
	var tok *token.Token

	Implements(t, tok, &Variable{})
}

func TestVariable(t *testing.T) {
	o := NewVariable(primitives.NewConstantInt(10))
	Equal(t, "10", o.String())

	r := test.NewRandTest(0)
	o.FuzzAll(r)
	Equal(t, "10", o.String())

	o2 := o.Clone()
	Equal(t, o.String(), o2.String())

	Equal(t, 1, o.Permutations())

	Nil(t, o.Permutation(1))
	Equal(t, "10", o.String())

	Equal(t, o.Permutation(2).(*token.PermutationError).Type, token.PermutationErrorIndexOutOfBound)
}
