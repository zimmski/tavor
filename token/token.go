package token

import (
	"fmt"

	"github.com/zimmski/tavor/rand"
)

type Token interface {
	fmt.Stringer

	Clone() Token

	Fuzz(r rand.Rand)
	FuzzAll(r rand.Rand)

	Permutation(i int) error
	Permutations() int
	PermutationsAll() int

	Parse(parser InternalParser, cur *ParserList) []ParserList
}

type ForwardToken interface {
	Token

	Get() Token

	InternalGet() Token
	InternalLogicalRemove(tok Token) Token
	InternalReplace(oldToken, newToken Token)
}

type OptionalToken interface {
	Token

	IsOptional() bool
	Activate()
	Deactivate()
}

type PermutationErrorType int

const (
	PermutationErrorIndexOutOfBound = iota
)

type PermutationError struct {
	Type PermutationErrorType
}

func (err *PermutationError) Error() string {
	switch err.Type {
	case PermutationErrorIndexOutOfBound:
		return "Permutation index out of bound"
	default:
		return fmt.Sprintf("Unknown permutation error type %#v", err.Type)
	}
}

type ResetToken interface {
	Token

	Reset()
}

type InternalParser struct { // TODO move this some place else
	Data string
}

type ParserList struct { // TODO move this some place else
	Tokens []ParserToken
	Index  int
}

type ParserToken struct { // TODO move this some place else
	Token Token
	Index int
}
