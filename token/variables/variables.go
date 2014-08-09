package variables

import (
	"github.com/zimmski/tavor/rand"
	"github.com/zimmski/tavor/token"
)

type Variable struct {
	value token.Token
}

func NewVariable(value token.Token) *Variable {
	return &Variable{
		value: value,
	}
}

func (v *Variable) Clone() token.Token {
	return &Variable{
		value: v.value,
	}
}

func (v *Variable) Fuzz(r rand.Rand) {
	// do nothing
}

func (v *Variable) FuzzAll(r rand.Rand) {
	v.Fuzz(r)
}

func (v *Variable) Parse(pars *token.InternalParser, cur int) (int, []error) {
	panic("TODO implement")
}

func (v *Variable) Permutation(i int) error {
	permutations := v.Permutations()

	if i < 1 || i > permutations {
		return &token.PermutationError{
			Type: token.PermutationErrorIndexOutOfBound,
		}
	}

	// do nothing

	return nil
}

func (v *Variable) Permutations() int {
	return 1
}

func (v *Variable) PermutationsAll() int {
	return v.Permutations()
}

func (v *Variable) String() string {
	return v.value.String()
}
