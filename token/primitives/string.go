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

func (p *ConstantString) Parse(pars *token.InternalParser, cur *token.ParserList) ([]token.ParserList, error) {
	vLen := len(p.value)

	nextIndex := vLen + cur.Index

	if nextIndex > pars.DataLen {
		return nil, &token.ParserError{
			Message: fmt.Sprintf("Expected \"%s\" but got early EOF", p.value),
			Type:    token.ParseErrorUnexpectedEOF,
		}
	}

	if got := pars.Data[cur.Index:nextIndex]; p.value != got {
		return nil, &token.ParserError{
			Message: fmt.Sprintf("Expected \"%s\" but got \"%s\"", p.value, got),
			Type:    token.ParseErrorUnexpectedData,
		}
	}

	return []token.ParserList{
		token.ParserList{
			Tokens: append(cur.Tokens, token.ParserToken{
				Token:    p.Clone(),
				Index:    1,
				MaxIndex: 1,
			}),
			Index: nextIndex,
		},
	}, nil
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
