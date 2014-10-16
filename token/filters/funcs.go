package filters

import (
	"github.com/zimmski/tavor/rand"
	"github.com/zimmski/tavor/token"
)

// FuncFilter implements a filter token which takes a token and filters its output according to a fuzzing function
type FuncFilter struct {
	fuzzFunc   func(r rand.Rand, tok token.Token) interface{}
	stringFunc func(state interface{}, tok token.Token) string
	state      interface{}
	token      token.Token
}

// NewFuncFilter returns a new instance of a FuncFilter token give the referenced token, a fuzzing and a stringer function
func NewFuncFilter(
	tok token.Token,
	fuzzFunc func(r rand.Rand, tok token.Token) interface{},
	stringFunc func(state interface{}, tok token.Token) string,
) *FuncFilter {
	return &FuncFilter{
		fuzzFunc:   fuzzFunc,
		stringFunc: stringFunc,
		state:      nil,
		token:      tok,
	}
}

// Clone returns a copy of the token and all its children
func (f *FuncFilter) Clone() token.Token {
	return &FuncFilter{
		fuzzFunc:   f.fuzzFunc,
		stringFunc: f.stringFunc,
		state:      f.state,
		token:      f.token.Clone(),
	}
}

// Fuzz fuzzes this token using the random generator by choosing one of the possible permutations for this token
func (f *FuncFilter) Fuzz(r rand.Rand) {
	f.state = f.fuzzFunc(r, f.token)
}

// FuzzAll calls Fuzz for this token and then FuzzAll for all children of this token
func (f *FuncFilter) FuzzAll(r rand.Rand) {
	f.Fuzz(r)
}

// Parse tries to parse the token beginning from the current position in the parser data.
// If the parsing is successful the error argument is nil and the next current position after the token is returned.
func (f *FuncFilter) Parse(pars *token.InternalParser, cur int) (int, []error) {
	panic("TODO implement")
}

// Permutation sets a specific permutation for this token
func (f *FuncFilter) Permutation(i uint) error {
	panic("TODO implemented")
}

// Permutations returns the number of permutations for this token
func (f *FuncFilter) Permutations() uint {
	return 1 // TODO this depends on the function
}

// PermutationsAll returns the number of all possible permutations for this token including its children
func (f *FuncFilter) PermutationsAll() uint {
	return f.Permutations()
}

func (f *FuncFilter) String() string {
	return f.stringFunc(f.state, f.token)
}
