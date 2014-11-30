package expressions

import (
	"github.com/zimmski/tavor/token"
)

// FuncExpression implements a expression token which executes a given list on output
type FuncExpression struct {
	permutationFunc     func(state interface{}, i uint) interface{}
	permutationsFunc    func(state interface{}) uint
	permutationsAllFunc func(state interface{}) uint
	stringFunc          func(state interface{}) string
	state               interface{}
}

// NewFuncExpression returns a new instance of a FuncExpression token given the output function
func NewFuncExpression(
	state interface{},
	permutationFunc func(state interface{}, i uint) interface{},
	permutationsFunc func(state interface{}) uint,
	permutationsAllFunc func(state interface{}) uint,
	stringFunc func(state interface{}) string,
) *FuncExpression {
	return &FuncExpression{
		permutationFunc:     permutationFunc,
		permutationsFunc:    permutationsFunc,
		permutationsAllFunc: permutationsAllFunc,
		stringFunc:          stringFunc,
		state:               state,
	}
}

// Clone returns a copy of the token and all its children
func (e *FuncExpression) Clone() token.Token {
	return &FuncExpression{
		permutationFunc:     e.permutationFunc,
		permutationsFunc:    e.permutationsFunc,
		permutationsAllFunc: e.permutationsAllFunc,
		stringFunc:          e.stringFunc,
		state:               e.state,
	}
}

// Parse tries to parse the token beginning from the current position in the parser data.
// If the parsing is successful the error argument is nil and the next current position after the token is returned.
func (e *FuncExpression) Parse(pars *token.InternalParser, cur int) (int, []error) {
	panic("TODO implement")
}

// Permutation sets a specific permutation for this token
func (e *FuncExpression) Permutation(i uint) error {
	permutations := e.Permutations()

	if i < 1 || i > permutations {
		return &token.PermutationError{
			Type: token.PermutationErrorIndexOutOfBound,
		}
	}

	e.state = e.permutationFunc(e.state, i-1)

	return nil
}

// Permutations returns the number of permutations for this token
func (e *FuncExpression) Permutations() uint {
	return e.permutationsFunc(e.state)
}

// PermutationsAll returns the number of all possible permutations for this token including its children
func (e *FuncExpression) PermutationsAll() uint {
	return e.permutationsAllFunc(e.state)
}

func (e *FuncExpression) String() string {
	return e.stringFunc(e.state)
}
