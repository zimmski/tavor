package conditions

import (
	"fmt"

	"github.com/zimmski/tavor/rand"
	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/lists"
)

type BooleanExpression interface {
	lists.List

	Evaluate() bool
}

type BooleanTrue struct{}

func NewBooleanTrue() *BooleanTrue {
	return &BooleanTrue{}
}

func (c *BooleanTrue) Evaluate() bool {
	return true
}

// Token interface methods

func (c *BooleanTrue) Clone() token.Token {
	return &BooleanTrue{}
}

func (c *BooleanTrue) Fuzz(r rand.Rand) {
	// do nothing
}

func (c *BooleanTrue) FuzzAll(r rand.Rand) {
	// do nothing
}

func (c *BooleanTrue) Parse(pars *token.InternalParser, cur int) (int, []error) {
	panic("This should never happen")
}

func (c *BooleanTrue) Permutation(i int) error {
	// do nothing

	return nil
}

func (c *BooleanTrue) Permutations() int {
	return 1
}

func (c *BooleanTrue) PermutationsAll() int {
	return 1
}

func (c *BooleanTrue) String() string {
	return "true"
}

// List interface methods

func (c *BooleanTrue) Get(i int) (token.Token, error) {
	return nil, &lists.ListError{
		Type: lists.ListErrorOutOfBound,
	}
}

func (c *BooleanTrue) Len() int {
	return 0
}

func (c *BooleanTrue) InternalGet(i int) (token.Token, error) {
	return nil, &lists.ListError{
		Type: lists.ListErrorOutOfBound,
	}
}

func (c *BooleanTrue) InternalLen() int {
	return 0
}

func (c *BooleanTrue) InternalLogicalRemove(tok token.Token) token.Token {
	panic("This should never happen")
}

func (c *BooleanTrue) InternalReplace(oldToken, newToken token.Token) {
	panic("This should never happen")
}

type BooleanEqual struct {
	a, b token.Token
}

func NewBooleanEqual(a, b token.Token) *BooleanEqual {
	return &BooleanEqual{
		a: a,
		b: b,
	}
}

func (c *BooleanEqual) Evaluate() bool {
	return c.a.String() == c.b.String()
}

// Token interface methods

func (c *BooleanEqual) Clone() token.Token {
	return &BooleanEqual{
		a: c.a,
		b: c.b,
	}
}

func (c *BooleanEqual) Fuzz(r rand.Rand) {
	// do nothing
}

func (c *BooleanEqual) FuzzAll(r rand.Rand) {
	// do nothing
}

func (c *BooleanEqual) Parse(pars *token.InternalParser, cur int) (int, []error) {
	panic("This should never happen")
}

func (c *BooleanEqual) Permutation(i int) error {
	// do nothing

	return nil
}

func (c *BooleanEqual) Permutations() int {
	return 1
}

func (c *BooleanEqual) PermutationsAll() int {
	return 1
}

func (c *BooleanEqual) String() string {
	return fmt.Sprintf("(%p)%#v == (%p)%#v", c.a, c.a, c.b, c.b)
}

// List interface methods

func (c *BooleanEqual) Get(i int) (token.Token, error) {
	return nil, &lists.ListError{
		Type: lists.ListErrorOutOfBound,
	}
}

func (c *BooleanEqual) Len() int {
	return 0
}

func (c *BooleanEqual) InternalGet(i int) (token.Token, error) {
	switch i {
	case 0:
		return c.a, nil
	case 1:
		return c.b, nil
	default:
		return nil, &lists.ListError{
			Type: lists.ListErrorOutOfBound,
		}
	}
}

func (c *BooleanEqual) InternalLen() int {
	return 2
}

func (c *BooleanEqual) InternalLogicalRemove(tok token.Token) token.Token {
	if tok == c.a || tok == c.b {
		return nil
	}

	return c
}

func (c *BooleanEqual) InternalReplace(oldToken, newToken token.Token) {
	if oldToken == c.a {
		c.a = newToken
	}
	if oldToken == c.b {
		c.b = newToken
	}
}
