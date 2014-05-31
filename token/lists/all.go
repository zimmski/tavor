package logicals

import (
	"bytes"

	"github.com/zimmski/tavor/rand"
	"github.com/zimmski/tavor/token"
)

type All struct {
	tokens []token.Token
}

func NewAll(toks ...token.Token) *All {
	if len(toks) == 0 {
		panic("at least one token needed")
	}

	return &All{
		tokens: toks,
	}
}

func (a *All) Fuzz(r rand.Rand) {
	for _, tok := range a.tokens {
		tok.Fuzz(r)
	}
}

func (a *All) String() string {
	var buffer bytes.Buffer

	for _, tok := range a.tokens {
		buffer.WriteString(tok.String())
	}

	return buffer.String()
}
