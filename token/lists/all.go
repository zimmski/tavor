package lists

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

func (l *All) Clone() token.Token {
	c := All{
		tokens: make([]token.Token, len(l.tokens)),
	}

	for i, tok := range l.tokens {
		c.tokens[i] = tok.Clone()
	}

	return &c
}

func (l *All) Fuzz(r rand.Rand) {
	for _, tok := range l.tokens {
		tok.Fuzz(r)
	}
}

func (l *All) Get(i int) (token.Token, error) {
	if i < 0 || i >= len(l.tokens) {
		return nil, &ListError{ListErrorOutOfBound}
	}

	return l.tokens[i], nil
}

func (l *All) Len() int {
	return len(l.tokens)
}

func (l *All) String() string {
	var buffer bytes.Buffer

	for _, tok := range l.tokens {
		buffer.WriteString(tok.String())
	}

	return buffer.String()
}
