package expressions

import (
	"github.com/zimmski/tavor/rand"
	"github.com/zimmski/tavor/token"
)

// FuncExpression implements a expression token which executes a given list on output
type FuncExpression struct {
	function func() string
}

// NewFuncExpression returns a new instance of a FuncExpression token given the output function
func NewFuncExpression(f func() string) *FuncExpression {
	return &FuncExpression{
		function: f,
	}
}

// Clone returns a copy of the token and all its children
func (e *FuncExpression) Clone() token.Token {
	return &FuncExpression{
		function: e.function,
	}
}

// Fuzz fuzzes this token using the random generator by choosing one of the possible permutations for this token
func (e *FuncExpression) Fuzz(r rand.Rand) {
	// do nothing
}

// FuzzAll calls Fuzz for this token and then FuzzAll for all children of this token
func (e *FuncExpression) FuzzAll(r rand.Rand) {
	e.Fuzz(r)
}

// Parse tries to parse the token beginning from the current position in the parser data.
// If the parsing is successful the error argument is nil and the next current position after the token is returned.
func (e *FuncExpression) Parse(pars *token.InternalParser, cur int) (int, []error) {
	panic("TODO implement")
}

// Permutation sets a specific permutation for this token
func (e *FuncExpression) Permutation(i uint) error {
	panic("TODO Not implemented")
}

// Permutations returns the number of permutations for this token
func (e *FuncExpression) Permutations() uint {
	return 1 // TODO this depends on the function
}

// PermutationsAll returns the number of all possible permutations for this token including its children
func (e *FuncExpression) PermutationsAll() uint {
	return e.Permutations()
}

func (e *FuncExpression) String() string {
	return e.function()
}
