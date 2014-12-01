package primitives

import (
	"fmt"

	"github.com/zimmski/tavor/log"
	"github.com/zimmski/tavor/token"
)

// ConstantString implements a string token which holds a constant string
type ConstantString struct {
	value string
}

// NewConstantString returns a new instance of a ConstantString token
func NewConstantString(value string) *ConstantString {
	return &ConstantString{
		value: value,
	}
}

// Clone returns a copy of the token and all its children
func (p *ConstantString) Clone() token.Token {
	return &ConstantString{
		value: p.value,
	}
}

// Parse tries to parse the token beginning from the current position in the parser data.
// If the parsing is successful the error argument is nil and the next current position after the token is returned.
func (p *ConstantString) Parse(pars *token.InternalParser, cur int) (int, []error) {
	vLen := len(p.value)

	nextIndex := vLen + cur

	if nextIndex > pars.DataLen {
		return cur, []error{&token.ParserError{
			Message: fmt.Sprintf("expected %q but got early EOF", p.value),
			Type:    token.ParseErrorUnexpectedEOF,

			Position: pars.GetPosition(cur),
		}}
	}

	if got := pars.Data[cur:nextIndex]; p.value != got {
		return cur, []error{&token.ParserError{
			Message: fmt.Sprintf("expected %q but got %q", p.value, got),
			Type:    token.ParseErrorUnexpectedData,

			Position: pars.GetPosition(cur),
		}}
	}

	log.Debugf("Parsed %q", p.value)

	return nextIndex, nil
}

// Permutation sets a specific permutation for this token
func (p *ConstantString) Permutation(i uint) error {
	permutations := p.Permutations()

	if i < 1 || i > permutations {
		return &token.PermutationError{
			Type: token.PermutationErrorIndexOutOfBound,
		}
	}

	// do nothing

	return nil
}

// Permutations returns the number of permutations for this token
func (p *ConstantString) Permutations() uint {
	return 1
}

// PermutationsAll returns the number of all possible permutations for this token including its children
func (p *ConstantString) PermutationsAll() uint {
	return p.Permutations()
}

func (p *ConstantString) String() string {
	return p.value
}
