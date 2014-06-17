package constraints

import (
	"github.com/zimmski/tavor/rand"
	"github.com/zimmski/tavor/token"
)

type Optional struct {
	token token.Token
	value bool
}

func NewOptional(tok token.Token) *Optional {
	return &Optional{
		token: tok,
		value: false,
	}
}

func (c *Optional) Clone() token.Token {
	return &Optional{
		token: c.token,
		value: c.value,
	}
}

func (c *Optional) Fuzz(r rand.Rand) {
	c.permutation(r.Int() % 2)
}

func (c *Optional) FuzzAll(r rand.Rand) {
	c.Fuzz(r)

	if !c.value {
		c.token.FuzzAll(r)
	}
}

func (c *Optional) Get() token.Token {
	return c.token
}

func (c *Optional) permutation(i int) {
	c.value = i == 0
}

func (c *Optional) Permutation(i int) error {
	permutations := c.Permutations()

	if i < 1 || i > permutations {
		return &token.PermutationError{
			Type: token.PermutationErrorIndexOutOfBound,
		}
	}

	c.permutation(i - 1)

	return nil
}

func (c *Optional) Permutations() int {
	return 2
}

func (c *Optional) PermutationsAll() int {
	return 1 + c.token.PermutationsAll()
}

func (c *Optional) String() string {
	if c.value {
		return ""
	}

	return c.token.String()
}

// OptionalToken interface methods

func (c *Optional) IsOptional() bool { return true }
func (c *Optional) Activate()        { c.value = false }
func (c *Optional) Deactivate()      { c.value = true }
