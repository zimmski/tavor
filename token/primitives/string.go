package primitives

import (
	"fmt"

	"github.com/zimmski/tavor/rand"
	"github.com/zimmski/tavor/token"
)

type ConstantString struct {
	value string
}

func NewConstantString(value string) *ConstantString {
	return &ConstantString{
		value: value,
	}
}

func (p *ConstantString) Clone() token.Token {
	return &ConstantString{
		value: p.value,
	}
}

func (p *ConstantString) Fuzz(r rand.Rand) {
	// do nothing
}

func (p *ConstantString) FuzzAll(r rand.Rand) {
	p.Fuzz(r)
}

func (p *ConstantString) Parse(pars *token.InternalParser, cur int) (int, []error) {
	vLen := len(p.value)

	nextIndex := vLen + cur

	if nextIndex > pars.DataLen {
		return cur, []error{&token.ParserError{
			Message: fmt.Sprintf("expected %q but got early EOF", p.value),
			Type:    token.ParseErrorUnexpectedEOF,
		}}
	}

	if got := pars.Data[cur:nextIndex]; p.value != got {
		return cur, []error{&token.ParserError{
			Message: fmt.Sprintf("expected %q but got %q", p.value, got),
			Type:    token.ParseErrorUnexpectedData,
		}}
	}

	return nextIndex, nil
}

func (p *ConstantString) Permutation(i int) error {
	permutations := p.Permutations()

	if i < 1 || i > permutations {
		return &token.PermutationError{
			Type: token.PermutationErrorIndexOutOfBound,
		}
	}

	// do nothing

	return nil
}

func (p *ConstantString) Permutations() int {
	return 1
}

func (p *ConstantString) PermutationsAll() int {
	return p.Permutations()
}

func (p *ConstantString) String() string {
	return p.value
}
