package aggregates

import (
	"strconv"

	"github.com/zimmski/tavor/token"
)

// Len implements a aggregation token that returns the length of a LenToken token
type Len struct {
	token token.LenToken
}

// NewLen returns a new instance of a Len token referencing the given LenToken token
func NewLen(token token.LenToken) *Len {
	return &Len{
		token: token,
	}
}

// Clone returns a copy of the token and all its children
func (a *Len) Clone() token.Token {
	return &Len{
		token: a.token,
	}
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
	return strconv.Itoa(a.token.Len())
}
