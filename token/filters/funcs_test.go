package filters

import (
	"testing"

	. "github.com/zimmski/tavor/test/assert"

	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/primitives"
)

func TestFuncFilterTokensToBeTokens(t *testing.T) {
	var tok *token.Token

	Implements(t, tok, &FuncFilter{})
}

func TestFuncExpression(t *testing.T) {
	o := NewFuncFilter(
		primitives.NewRangeInt(1, 10),
		false,
		func(state interface{}, tok token.Token, i uint) interface{} {
			return i == 1
		},
		func(state interface{}, tok token.Token) uint {
			return 2
		},
		func(state interface{}, tok token.Token) uint {
			return 2 * tok.PermutationsAll()
		},
		func(state interface{}, tok token.Token) string {
			i, ok := state.(bool)
			if !ok {
				panic("unknown type")
			}

			if i {
				return tok.String()
			}

			return ""
		},
	)
	Equal(t, 2, o.Permutations())
	Equal(t, 20, o.PermutationsAll())
	Equal(t, "", o.String())

	Nil(t, o.Permutation(1))
	Equal(t, "", o.String())

	Nil(t, o.Permutation(2))
	Equal(t, "1", o.String())

	o2 := o.Clone()
	Equal(t, o.String(), o2.String())
}
