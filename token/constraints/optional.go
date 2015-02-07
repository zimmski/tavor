package constraints

import (
	"github.com/zimmski/tavor/token"
)

// Optional implements a constraint and optional token which references another token which can be de(activated)
type Optional struct {
	token token.Token
	value bool

	reducing              bool
	reducingOriginalValue bool
}

// NewOptional returns a new instance of a Optional token referencing the given token and setting the initial state to deactivated
func NewOptional(tok token.Token) *Optional {
	return &Optional{
		token: tok,
		value: false,
	}
}

// Token interface methods

// Clone returns a copy of the token and all its children
func (c *Optional) Clone() token.Token {
	return &Optional{
		token: c.token.Clone(),
		value: c.value,
	}
}

// Parse tries to parse the token beginning from the current position in the parser data.
// If the parsing is successful the error argument is nil and the next current position after the token is returned.
func (c *Optional) Parse(pars *token.InternalParser, cur int) (int, []error) {
	nex, errs := c.token.Parse(pars, cur)

	if len(errs) == 0 {
		c.value = false

		return nex, nil
	}

	c.value = true

	return cur, nil
}

func (c *Optional) permutation(i uint) {
	c.value = i == 0
}

// Permutation sets a specific permutation for this token
func (c *Optional) Permutation(i uint) error {
	permutations := c.Permutations()

	if i < 1 || i > permutations {
		return &token.PermutationError{
			Type: token.PermutationErrorIndexOutOfBound,
		}
	}

	c.permutation(i - 1)

	return nil
}

// Permutations returns the number of permutations for this token
func (c *Optional) Permutations() uint {
	return 2
}

// PermutationsAll returns the number of all possible permutations for this token including its children
func (c *Optional) PermutationsAll() uint {
	return 1 + c.token.PermutationsAll()
}

func (c *Optional) String() string {
	if c.value {
		return ""
	}

	return c.token.String()
}

// ForwardToken interface methods

// Get returns the current referenced token
func (c *Optional) Get() token.Token {
	if c.value {
		return nil
	}

	return c.token
}

// InternalGet returns the current referenced internal token
func (c *Optional) InternalGet() token.Token {
	return c.token
}

// InternalLogicalRemove removes the referenced internal token and returns the replacement for the current token or nil if the current token should be removed.
func (c *Optional) InternalLogicalRemove(tok token.Token) token.Token {
	if c.token == tok {
		return nil
	}

	return c
}

// InternalReplace replaces an old with a new internal token if it is referenced by this token. The error return argument is not nil, if the replacement is not suitable.
func (c *Optional) InternalReplace(oldToken, newToken token.Token) error {
	if c.token == oldToken {
		c.token = newToken
	}

	return nil
}

// OptionalToken interface methods

// IsOptional checks dynamically if this token is in the current state optional
func (c *Optional) IsOptional() bool { return true }

// Activate activates this token
func (c *Optional) Activate() { c.value = false }

// Deactivate deactivates this token
func (c *Optional) Deactivate() { c.value = true }

// ReduceToken interface methods

// Reduce sets a specific reduction for this token
func (c *Optional) Reduce(i uint) error {
	reduces := c.Permutations()

	if reduces == 0 || i < 1 || i > reduces {
		return &token.ReduceError{
			Type: token.ReduceErrorIndexOutOfBound,
		}
	}

	if !c.reducing {
		c.reducing = true
		c.reducingOriginalValue = c.value
	}

	c.permutation(i - 1)

	return nil
}

// Reduces returns the number of reductions for this token
func (c *Optional) Reduces() uint {
	if c.reducing || !c.value {
		return 2
	}

	return 0
}
