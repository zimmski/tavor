package expressions

import (
	"testing"

	. "github.com/zimmski/tavor/test/assert"

	"github.com/zimmski/tavor/token"
)

func TestFuncExpressionTokensToBeTokens(t *testing.T) {
	var tok *token.Token

	Implements(t, tok, &FuncExpression{})
}

func TestFuncExpression(t *testing.T) {
	o := NewFuncExpression(
		false,
		func(state interface{}, i uint) interface{} {
			return i == 1
		},
		func(state interface{}) uint {
			return 2
		},
		func(state interface{}) uint {
			return 2
		},
		func(state interface{}) string {
			i, ok := state.(bool)
			if !ok {
				panic("unknown type")
			}

			if i {
				return "abc"
			}

			return ""
		},
	)
	Equal(t, 2, o.Permutations())
	Equal(t, 2, o.PermutationsAll())
	Equal(t, "", o.String())

	Nil(t, o.Permutation(1))
	Equal(t, "", o.String())

	Nil(t, o.Permutation(2))
	Equal(t, "abc", o.String())

	o2 := o.Clone()
	Equal(t, o.String(), o2.String())
}
