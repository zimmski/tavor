package lists

import (
	"github.com/zimmski/tavor/token"
)

type ListErrorType int

const (
	ListErrorOutOfBound ListErrorType = iota
)

type ListError struct {
	Type ListErrorType
}

func (err *ListError) Error() string {
	switch err.Type {
	default:
		return "Out of bound"
	}
}

type List interface {
	token.Token

	Get(i int) (token.Token, error)
	Len() int
}
