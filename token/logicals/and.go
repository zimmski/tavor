package logicals

import (
	"bytes"

	"github.com/zimmski/tavor/rand"
	"github.com/zimmski/tavor/token"
)

type And struct {
	tokens []token.Token
}

func NewAnd(tokens ...token.Token) *And {
	if len(tokens) == 0 {
		panic("at least one token needed")
	}

	return &And{
		tokens: tokens,
	}
}

func (a *And) Fuzz(r rand.Rand) {
	for _, tok := range a.tokens {
		tok.Fuzz(r)
	}
}

func (a *And) String() string {
	var buffer bytes.Buffer

	for _, tok := range a.tokens {
		buffer.WriteString(tok.String())
	}

	return buffer.String()
}
