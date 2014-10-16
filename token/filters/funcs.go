package filters

import (
	"github.com/zimmski/tavor/rand"
	"github.com/zimmski/tavor/token"
)

type FuncFilter struct {
	fuzzFunc   func(r rand.Rand, tok token.Token) interface{}
	stringFunc func(state interface{}, tok token.Token) string
	state      interface{}
	token      token.Token
}

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

func (f *FuncFilter) Fuzz(r rand.Rand) {
	f.state = f.fuzzFunc(r, f.token)
}

func (f *FuncFilter) FuzzAll(r rand.Rand) {
	f.Fuzz(r)
}

func (f *FuncFilter) Parse(pars *token.InternalParser, cur int) (int, []error) {
	panic("TODO implement")
}

func (f *FuncFilter) Permutation(i uint) error {
	panic("TODO implemented")
}

func (f *FuncFilter) Permutations() uint {
	return 1 // TODO this depends on the function
}

func (f *FuncFilter) PermutationsAll() uint {
	return f.Permutations()
}

func (f *FuncFilter) String() string {
	return f.stringFunc(f.state, f.token)
}
