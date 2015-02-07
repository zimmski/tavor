package filters

import (
	"github.com/zimmski/tavor/token"
)

// FuncFilter implements a filter token which takes a token and filters its output according to a fuzzing function
type FuncFilter struct {
	permutationFunc     func(state interface{}, tok token.Token, i uint) interface{}
	permutationsFunc    func(state interface{}, tok token.Token) uint
	permutationsAllFunc func(state interface{}, tok token.Token) uint
	stringFunc          func(state interface{}, tok token.Token) string
	state               interface{}
	token               token.Token
}

// NewFuncFilter returns a new instance of a FuncFilter token give the referenced token, a fuzzing and a stringer function
func NewFuncFilter(
	tok token.Token,
	state interface{},
	permutationFunc func(state interface{}, tok token.Token, i uint) interface{},
	permutationsFunc func(state interface{}, tok token.Token) uint,
	permutationsAllFunc func(state interface{}, tok token.Token) uint,
	stringFunc func(state interface{}, tok token.Token) string,
) *FuncFilter {
	return &FuncFilter{
		permutationFunc:     permutationFunc,
		permutationsFunc:    permutationsFunc,
		permutationsAllFunc: permutationsAllFunc,
		stringFunc:          stringFunc,
		state:               state,
		token:               tok,
	}
}

// Clone returns a copy of the token and all its children
func (f *FuncFilter) Clone() token.Token {
	return &FuncFilter{
		permutationFunc:     f.permutationFunc,
		permutationsFunc:    f.permutationsFunc,
		permutationsAllFunc: f.permutationsAllFunc,
		stringFunc:          f.stringFunc,
		state:               f.state,
		token:               f.token.Clone(),
	}
}

// Parse tries to parse the token beginning from the current position in the parser data.
// If the parsing is successful the error argument is nil and the next current position after the token is returned.
func (f *FuncFilter) Parse(pars *token.InternalParser, cur int) (int, []error) {
	panic("TODO implement")
}

// Permutation sets a specific permutation for this token
func (f *FuncFilter) Permutation(i uint) error {
	permutations := f.Permutations()

	if i < 1 || i > permutations {
		return &token.PermutationError{
			Type: token.PermutationErrorIndexOutOfBound,
		}
	}

	f.state = f.permutationFunc(f.state, f.token, i-1)

	return nil
}

// Permutations returns the number of permutations for this token
func (f *FuncFilter) Permutations() uint {
	return f.permutationsFunc(f.state, f.token)
}

// PermutationsAll returns the number of all possible permutations for this token including its children
func (f *FuncFilter) PermutationsAll() uint {
	return f.permutationsAllFunc(f.state, f.token)
}

func (f *FuncFilter) String() string {
	return f.stringFunc(f.state, f.token)
}

// ForwardToken interface methods

// Get returns the current referenced token
func (f *FuncFilter) Get() token.Token {
	return f.token
}

// InternalGet returns the current referenced internal token
func (f *FuncFilter) InternalGet() token.Token {
	return f.token
}

// InternalLogicalRemove removes the referenced internal token and returns the replacement for the current token or nil if the current token should be removed.
func (f *FuncFilter) InternalLogicalRemove(tok token.Token) token.Token {
	if f.token == tok {
		return nil
	}

	return f
}

// InternalReplace replaces an old with a new internal token if it is referenced by this token. The error return argument is not nil, if the replacement is not suitable.
func (f *FuncFilter) InternalReplace(oldToken, newToken token.Token) error {
	if f.token == oldToken {
		f.token = newToken
	}

	return nil
}
