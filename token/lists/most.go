package lists

import (
	"bytes"

	"github.com/zimmski/tavor/rand"
	"github.com/zimmski/tavor/token"
)

type Most struct {
	n     int64
	token token.Token
	value []token.Token
}

func NewMost(tok token.Token, n int64) *Most {
	l := &Most{
		n:     n,
		token: tok,
		value: make([]token.Token, n),
	}

	for i := range l.value {
		l.value[i] = tok.Clone()
	}

	return l
}

func (l *Most) Clone() token.Token {
	c := Most{
		n:     l.n,
		token: l.token,
		value: make([]token.Token, len(l.value)),
	}

	for i, tok := range l.value {
		c.value[i] = tok.Clone()
	}

	return &c
}

func (l *Most) Fuzz(r rand.Rand) {
	n := r.Intn(int(l.n) + 1)
	toks := make([]token.Token, n)

	for i := range toks {
		toks[i] = l.token.Clone()

		toks[i].Fuzz(r)
	}

	l.value = toks
}

func (l *Most) Get(i int) (token.Token, error) {
	if i < 0 || i >= 1 {
		return nil, &ListError{ListErrorOutOfBound}
	}

	return l.token, nil
}

func (l *Most) Len() int {
	return len(l.value)
}

func (l *Most) String() string {
	var buffer bytes.Buffer

	for _, tok := range l.value {
		buffer.WriteString(tok.String())
	}

	return buffer.String()
}
