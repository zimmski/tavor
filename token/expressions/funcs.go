package expressions

import (
	"github.com/zimmski/tavor/rand"
	"github.com/zimmski/tavor/token"
)

type FuncExpression struct {
	function func() string
}

func NewFuncExpression(f func() string) *FuncExpression {
	return &FuncExpression{
		function: f,
	}
}

func (e *FuncExpression) Clone() token.Token {
	return &FuncExpression{
		function: e.function,
	}
}

func (e *FuncExpression) Fuzz(r rand.Rand) {
	// do nothing
}

func (e *FuncExpression) FuzzAll(r rand.Rand) {
	e.Fuzz(r)
}

func (e *FuncExpression) Permutations() int {
	return 1 // TODO this depends on the function
}

func (e *FuncExpression) String() string {
	return e.function()
}
