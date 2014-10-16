package aggregates

import (
	"strconv"

	"github.com/zimmski/tavor/rand"
	"github.com/zimmski/tavor/token"
)

// Len implements a aggregation token that returns the length of a List token
type Len struct {
	list token.List
}

// NewLen returns a new instance of a Len token referencing the given List token
func NewLen(list token.List) *Len {
	return &Len{
		list: list,
	}
}

// Clone returns a copy of the token and all its children
func (a *Len) Clone() token.Token {
	return &Len{
		list: a.list,
	}
}

// Fuzz fuzzes this token using the random generator by choosing one of the possible permutations for this token
func (a *Len) Fuzz(r rand.Rand) {
	// do nothing
}

// FuzzAll calls Fuzz for this token and then FuzzAll for all children of this token
func (a *Len) FuzzAll(r rand.Rand) {
	a.Fuzz(r)
}

// Parse tries to parse the token beginning from the current position in the parser data.
// If the parsing is successful the error argument is nil and the next current position after the token is returned.
func (a *Len) Parse(pars *token.InternalParser, cur int) (int, []error) {
	panic("TODO implement")
}

// Permutation sets a specific permutation for this token
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

// Permutations returns the number of permutations for this token
func (a *Len) Permutations() uint {
	return 1
}

// PermutationsAll returns the number of all possible permutations for this token including its children
func (a *Len) PermutationsAll() uint {
	return a.Permutations()
}

func (a *Len) String() string {
	return strconv.Itoa(a.list.Len())
}
