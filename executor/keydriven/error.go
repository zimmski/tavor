package keydriven

import (
	"errors"
	"fmt"
)

var (
	// ErrKeyAlreadyDefined the given key is already defined in the executor
	ErrKeyAlreadyDefined = errors.New("Key is already defined")
	// ErrKeyNotDefined the given key is not defined in the executor
	ErrKeyNotDefined = errors.New("Key is not defined")
	// ErrInvalidParametersCount indicates that the given action has different parameter count requirements
	ErrInvalidParametersCount = errors.New("Invalid parmaters count")
)

// Error defines a key driven specific error
type Error struct {
	Message string
	Err     error
}

func (e Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Err, e.Message)
}
