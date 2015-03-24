package token

import (
	"fmt"
	"sort"
	"text/scanner"
)

// ArgumentsTypedParser defines a parser for the arguments of a typed token.
// Parsing stops unrecoverably at the first error. The return value of Err must be checked before using the values returned by precedings calls to GetInt.
type ArgumentsTypedParser interface {
	// GetInt tries to parse the argument name and returns its integer value or defaultValue if the argument is not found.
	// The return value is valid only if Err returns nil.
	GetInt(name string, defaultValue int) int
	// Err returns the first error encountered by the ArgumentsTypedParser.
	Err() error
}

// CreateTypedFunc defines a function to create a typed token
type CreateTypedFunc func(argParser ArgumentsTypedParser) (Token, error)

// typedLookup is a mapping from typed token names to their creation functions
var typedLookup = make(map[string]CreateTypedFunc)

// NewTyped creates a typed token by calling the function registered with the given name.
// The error return argument is not nil if the name does not exist in the registered typed token list or if the token creation failed.
func NewTyped(name string, argParser ArgumentsTypedParser, pos scanner.Position) (Token, error) {
	createTok, ok := typedLookup[name]
	if !ok {
		return nil, &ParserError{
			Message:  fmt.Sprintf("unknown typed token %q", name),
			Type:     ParseErrorUnknownTypedTokenType,
			Position: pos,
		}
	}

	tok, err := createTok(argParser)
	if err != nil {
		return nil, &ParserError{
			Message:  err.Error(),
			Type:     ParseErrorInvalidArgumentValue,
			Position: pos,
		}
	}

	return tok, nil
}

// ListTyped returns a list of all registered typed token names.
func ListTyped() []string {
	typedNames := make([]string, 0, len(typedLookup))

	for key := range typedLookup {
		typedNames = append(typedNames, key)
	}

	sort.Strings(typedNames)

	return typedNames
}

// RegisterTyped registers a typed token creation function with the given name.
func RegisterTyped(name string, ct CreateTypedFunc) {
	if ct == nil {
		panic("register typed token is nil")
	}

	if _, ok := typedLookup[name]; ok {
		panic("typed token " + name + " already registered")
	}

	typedLookup[name] = ct
}
