package filters

import (
	"testing"

	. "github.com/stretchr/testify/assert"

	"github.com/zimmski/tavor/rand"
	"github.com/zimmski/tavor/test"
	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/primitives"
)

func TestFuncFilterTokensToBeTokens(t *testing.T) {
	var tok *token.Token

	Implements(t, tok, &FuncFilter{})
}

func TestFuncExpression(t *testing.T) {
	o := NewFuncFilter(
		primitives.NewConstantInt(1),
		func(r rand.Rand, tok token.Token) interface{} {
			c := r.Int()%2 == 0

			if c {
				tok.Fuzz(r)
			}

			return c
		},
		func(state interface{}, tok token.Token) string {
			switch i := state.(type) {
			case bool:
				if i {
					return tok.String()
				}

				return ""
			case nil:
				return tok.String()
			}

			panic("unknown type")
		},
	)
	Equal(t, "1", o.String())

	r := test.NewRandTest(0)
	o.Fuzz(r)
	Equal(t, "", o.String())

	o.Fuzz(r)
	Equal(t, "1", o.String())

	o2 := o.Clone()
	Equal(t, o.String(), o2.String())
}
