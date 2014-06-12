package expressions

import (
	"testing"

	. "github.com/stretchr/testify/assert"

	"github.com/zimmski/tavor/test"
	"github.com/zimmski/tavor/token"
)

func TestFuncExpressionTokensToBeTokens(t *testing.T) {
	var tok *token.Token

	Implements(t, tok, &FuncExpression{})
}

func TestFuncExpression(t *testing.T) {
	s := "abc"

	o := NewFuncExpression(func() string {
		return s
	})
	Equal(t, "abc", o.String())

	r := test.NewRandTest(0)
	o.FuzzAll(r)
	Equal(t, "abc", o.String())

	o2 := o.Clone()
	Equal(t, o.String(), o2.String())

	Equal(t, 1, o.Permutations())
}
