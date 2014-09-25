package aggregates

import (
	"strconv"

	"github.com/zimmski/tavor/rand"
	"github.com/zimmski/tavor/token"
)

type Len struct {
	list token.List
}

func NewLen(list token.List) *Len {
	return &Len{
		list: list,
	}
}

func (a *Len) Clone() token.Token {
	return &Len{
		list: a.list,
	}
}

func (a *Len) Fuzz(r rand.Rand) {
	// do nothing
}

func (a *Len) FuzzAll(r rand.Rand) {
	a.Fuzz(r)
}

func (a *Len) Parse(pars *token.InternalParser, cur int) (int, []error) {
	panic("TODO implement")
}

func (a *Len) Permutation(i uint) error {
	permutations := a.Permutations()

	if i < 1 || i > permutations {
		return &token.PermutationError{
			Type: token.PermutationErrorIndexOutOfBound,
		}
	}

	// do nothing

	return nil
}

func (a *Len) Permutations() uint {
	return 1
}

func (a *Len) PermutationsAll() uint {
	return a.Permutations()
}

func (a *Len) String() string {
	return strconv.Itoa(a.list.Len())
}
