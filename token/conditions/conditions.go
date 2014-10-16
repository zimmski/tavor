package conditions

import (
	"fmt"

	"github.com/zimmski/tavor"
	"github.com/zimmski/tavor/rand"
	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/lists"
)

type IfPair struct {
	Head BooleanExpression
	Body token.Token
}

// Token interface methods

func (c *IfPair) Clone() IfPair {
	return IfPair{
		Head: c.Head.Clone().(BooleanExpression),
		Body: c.Body.Clone(),
	}
}

func (c *IfPair) Fuzz(r rand.Rand) {
	// do nothing
}

func (c *IfPair) FuzzAll(r rand.Rand) {
	// do nothing
}

func (c *IfPair) Parse(pars *token.InternalParser, cur int) (int, []error) {
	panic("This should never happen")
}

func (c *IfPair) Permutation(i uint) error {
	// do nothing

	return nil
}

func (c *IfPair) Permutations() uint {
	return 1
}

func (c *IfPair) PermutationsAll() uint {
	return 1
}

func (c *IfPair) String() string {
	return fmt.Sprintf("(%p)%#v -> (%p)%#v", c.Head, c.Head, c.Body, c.Body)
}

// List interface methods

func (c *IfPair) Get(i int) (token.Token, error) {
	return nil, &lists.ListError{
		Type: lists.ListErrorOutOfBound,
	}
}

func (c *IfPair) Len() int {
	return 0
}

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

func (c *IfPair) InternalLen() int {
	return 2
}

func (c *IfPair) InternalLogicalRemove(tok token.Token) token.Token {
	panic("This should never happen")
}

func (c *IfPair) InternalReplace(oldToken, newToken token.Token) {
	panic("This should never happen")
}

type If struct {
	Pairs []IfPair
}

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
		nPairs[i] = c.Pairs[i].Clone()
	}

	return &If{
		Pairs: nPairs,
	}
}

func (c *If) Fuzz(r rand.Rand) {
	// do nothing
}

func (c *If) FuzzAll(r rand.Rand) {
	c.Fuzz(r)
}

func (c *If) Parse(pars *token.InternalParser, cur int) (int, []error) {
	panic("TODO implement")
}

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

func (c *If) Permutations() uint {
	return 1
}

func (c *If) PermutationsAll() uint {
	return c.Permutations()
}

func (c *If) String() string {
	for _, pair := range c.Pairs {
		if pair.Head.Evaluate() {
			return pair.Body.String()
		}
	}

	panic("This should not happen")
}

/*
// List interface methods

func (c *If) Get(i int) (token.Token, error) {
	return nil, &lists.ListError{
		Type: lists.ListErrorOutOfBound,
	}
}

func (c *If) Len() int {
	return 0
}

func (c *If) InternalGet(i int) (token.Token, error) {
	if i < 0 || i >= len(c.Pairs) {
		return nil, &lists.ListError{
			Type: lists.ListErrorOutOfBound,
		}
	}

	return &c.Pairs[i], nil
}

func (c *If) InternalLen() int {
	return len(c.Pairs)
}

func (c *If) InternalLogicalRemove(tok token.Token) token.Token {
	panic("TODO")
}

func (c *If) InternalReplace(oldToken, newToken token.Token) {
	panic("TODO")
}
*/

// ScopeToken interface methods

func (c *If) SetScope(variableScope map[string]token.Token) {
	for _, pair := range c.Pairs {
		tavor.SetInternalScope(pair.Head, variableScope)
		tavor.SetInternalScope(pair.Body, variableScope)
	}
}
