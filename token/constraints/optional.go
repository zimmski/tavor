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
}

func (c *Optional) FuzzAll(r rand.Rand) {
	c.Fuzz(r)

	if !c.value {
		c.token.FuzzAll(r)
	}
}

func (c *Optional) Permutations() int {
	return 2 * c.token.Permutations()
}

func (c *Optional) String() string {
	if c.value {
		return ""
	}

	return c.token.String()
}
