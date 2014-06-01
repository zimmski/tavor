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

func (o *Optional) Clone() token.Token {
	return &Optional{
		token: o.token,
		value: o.value,
	}
}

func (o *Optional) Fuzz(r rand.Rand) {
	o.value = r.Int()%2 == 0

	if !o.value {
		o.token.Fuzz(r)
	}
}

func (o *Optional) String() string {
	if o.value {
		return ""
	}

	return o.token.String()
}
