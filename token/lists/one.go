package lists

import (
	"github.com/zimmski/tavor/rand"
	"github.com/zimmski/tavor/token"
)

type One struct {
	tokens []token.Token
	value  token.Token
}

func NewOne(toks ...token.Token) *One {
	if len(toks) == 0 {
		panic("at least one token needed")
	}

	return &One{
		tokens: toks,
		value:  toks[0],
	}
}

func (o *One) Clone() token.Token {
	c := One{
		tokens: make([]token.Token, len(o.tokens)),
		value:  o.value.Clone(),
	}

	for i, tok := range o.tokens {
		c.tokens[i] = tok.Clone()
	}

	return &c
}

func (o *One) Fuzz(r rand.Rand) {
	i := r.Intn(len(o.tokens))

	o.value = o.tokens[i]

	o.value.Fuzz(r)
}

func (o *One) String() string {
	return o.value.String()
}
