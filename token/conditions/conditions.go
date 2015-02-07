package conditions

import (
	"fmt"

	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/lists"
)

// IfPair implements a condition token which holds an If condition with its head and body
type IfPair struct {
	Head BooleanExpression
	Body token.Token
}

// Token interface methods

// Clone returns a copy of the token and all its children
func (c *IfPair) Clone() token.Token {
	return &IfPair{
		Head: c.Head.Clone().(BooleanExpression),
		Body: c.Body.Clone(),
	}
}

// Parse tries to parse the token beginning from the current position in the parser data.
// If the parsing is successful the error argument is nil and the next current position after the token is returned.
func (c *IfPair) Parse(pars *token.InternalParser, cur int) (int, []error) {
	panic("This should never happen")
}

// Permutation sets a specific permutation for this token
func (c *IfPair) Permutation(i uint) error {
	// do nothing

	return nil
}

// Permutations returns the number of permutations for this token
func (c *IfPair) Permutations() uint {
	return 1
}

// PermutationsAll returns the number of all possible permutations for this token including its children
func (c *IfPair) PermutationsAll() uint {
	return 1
}

func (c *IfPair) String() string {
	return fmt.Sprintf("(%p)%#v -> (%p)%#v", c.Head, c.Head, c.Body, c.Body)
}

// List interface methods

// Get returns the current referenced token at the given index. The error return argument is not nil, if the index is out of bound.
func (c *IfPair) Get(i int) (token.Token, error) {
	return nil, &lists.ListError{
		Type: lists.ListErrorOutOfBound,
	}
}

// Len returns the number of the current referenced tokens
func (c *IfPair) Len() int {
	return 0
}

// InternalGet returns the current referenced internal token at the given index. The error return argument is not nil, if the index is out of bound.
func (c *IfPair) InternalGet(i int) (token.Token, error) {
	switch i {
	case 0:
		return c.Head, nil
	case 1:
		return c.Body, nil
	default:
		return nil, &lists.ListError{
			Type: lists.ListErrorOutOfBound,
		}
	}
}

// InternalLen returns the number of referenced internal tokens
func (c *IfPair) InternalLen() int {
	return 2
}

// InternalLogicalRemove removes the referenced internal token and returns the replacement for the current token or nil if the current token should be removed.
func (c *IfPair) InternalLogicalRemove(tok token.Token) token.Token {
	panic("This should never happen")
}

// InternalReplace replaces an old with a new internal token if it is referenced by this token. The error return argument is not nil, if the replacement is not suitable.
func (c *IfPair) InternalReplace(oldToken, newToken token.Token) error {
	panic("This should never happen")
}

// If implements a condition token which holds a list of IfPairs which belong together (e.g. If Elsif ... Else)
type If struct {
	Pairs []IfPair
}

// NewIf returns a new instance of a If token referencing a list of IfPairs
func NewIf(Pairs ...IfPair) *If {
	if len(Pairs) == 0 {
		panic("Must at least given one if pair")
	}

	return &If{
		Pairs: Pairs,
	}
}

// Token interface methods

// Clone returns a copy of the token and all its children
func (c *If) Clone() token.Token {
	nPairs := make([]IfPair, len(c.Pairs))
	for i := 0; i < len(c.Pairs); i++ {
		nPairs[i] = *c.Pairs[i].Clone().(*IfPair)
	}

	return &If{
		Pairs: nPairs,
	}
}

// Parse tries to parse the token beginning from the current position in the parser data.
// If the parsing is successful the error argument is nil and the next current position after the token is returned.
func (c *If) Parse(pars *token.InternalParser, cur int) (int, []error) {
	panic("TODO implement")
}

// Permutation sets a specific permutation for this token
func (c *If) Permutation(i uint) error {
	permutations := c.Permutations()

	if i < 1 || i > permutations {
		return &token.PermutationError{
			Type: token.PermutationErrorIndexOutOfBound,
		}
	}

	// do nothing

	return nil
}

// Permutations returns the number of permutations for this token
func (c *If) Permutations() uint {
	return 1
}

// PermutationsAll returns the number of all possible permutations for this token including its children
func (c *If) PermutationsAll() uint {
	return c.Permutations()
}

func (c *If) String() string {
	for _, pair := range c.Pairs {
		if pair.Head.Evaluate() {
			return pair.Body.String()
		}
	}

	return ""
}

/*
// List interface methods

// Get returns the current referenced token at the given index. The error return argument is not nil, if the index is out of bound.
func (c *If) Get(i int) (token.Token, error) {
	return nil, &lists.ListError{
		Type: lists.ListErrorOutOfBound,
	}
}

// Len returns the number of the current referenced tokens
func (c *If) Len() int {
	return 0
}

// InternalGet returns the current referenced internal token at the given index. The error return argument is not nil, if the index is out of bound.
func (c *If) InternalGet(i int) (token.Token, error) {
	if i < 0 || i >= len(c.Pairs) {
		return nil, &lists.ListError{
			Type: lists.ListErrorOutOfBound,
		}
	}

	return &c.Pairs[i], nil
}

// InternalLen returns the number of referenced internal tokens
func (c *If) InternalLen() int {
	return len(c.Pairs)
}

// InternalLogicalRemove removes the referenced internal token and returns the replacement for the current token or nil if the current token should be removed.
func (c *If) InternalLogicalRemove(tok token.Token) token.Token {
	panic("TODO")
}

// InternalReplace replaces an old with a new internal token if it is referenced by this token. The error return argument is not nil, if the replacement is not suitable.
func (c *If) InternalReplace(oldToken, newToken token.Token) error {
	panic("TODO")
}
*/

// ScopeToken interface methods

// SetScope sets the scope of the token
func (c *If) SetScope(variableScope *token.VariableScope) {
	for _, pair := range c.Pairs {
		token.SetScope(pair.Head, variableScope)
		token.SetScope(pair.Body, variableScope)
	}
}
