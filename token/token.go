package token

import (
	"fmt"

	"github.com/zimmski/tavor/rand"
)

type Token interface {
	fmt.Stringer

	Clone() Token

	Permutation(i int) error
	Permutations() int
	PermutationsAll() int

	Fuzz(r rand.Rand)
	FuzzAll(r rand.Rand)
}

type ForwardToken interface {
	Token

	Get() Token
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
