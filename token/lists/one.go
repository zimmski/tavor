package logicals

import (
	"github.com/zimmski/tavor/rand"
	"github.com/zimmski/tavor/token"
)

type XOr struct {
	tokens []token.Token
	value  token.Token
}

func NewXOr(tokens ...token.Token) *XOr {
	if len(tokens) == 0 {
		panic("at least one token needed")
	}

	return &XOr{
		tokens: tokens,
		value:  tokens[0],
	}
}

func (o *XOr) Fuzz(r rand.Rand) {
	i := r.Intn(len(o.tokens))

	o.value = o.tokens[i]

	o.value.Fuzz(r)
}

func (o *XOr) String() string {
	return o.value.String()
}
