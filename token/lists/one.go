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

func (l *One) Clone() token.Token {
	c := One{
		tokens: make([]token.Token, len(l.tokens)),
		value:  l.value.Clone(),
	}

	for i, tok := range l.tokens {
		c.tokens[i] = tok.Clone()
	}

	return &c
}

func (l *One) Fuzz(r rand.Rand) {
	i := r.Intn(len(l.tokens))

	l.value = l.tokens[i]

	l.value.Fuzz(r)
}

func (l *One) Len() int {
	return 1
}

func (l *One) String() string {
	return l.value.String()
}
