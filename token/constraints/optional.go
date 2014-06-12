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
	c.value = r.Int()%2 == 0

	if !c.value {
		c.token.Fuzz(r)
	}
}

func (c *Optional) Permutations() int {
	return 2
}

func (c *Optional) String() string {
	if c.value {
		return ""
	}

	return c.token.String()
}
