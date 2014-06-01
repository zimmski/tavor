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
	m := &Most{
		n:     n,
		token: tok,
		value: make([]token.Token, n),
	}

	for i := range m.value {
		m.value[i] = tok.Clone()
	}

	return m
}

func (m *Most) Clone() token.Token {
	c := Most{
		n:     m.n,
		token: m.token,
		value: make([]token.Token, len(m.value)),
	}

	for i, tok := range m.value {
		c.value[i] = tok.Clone()
	}

	return &c
}

func (m *Most) Fuzz(r rand.Rand) {
	n := r.Intn(int(m.n) + 1)
	toks := make([]token.Token, n)

	for i := range toks {
		toks[i] = m.token.Clone()

		toks[i].Fuzz(r)
	}

	m.value = toks
}

func (m *Most) String() string {
	var buffer bytes.Buffer

	for _, tok := range m.value {
		buffer.WriteString(tok.String())
	}

	return buffer.String()
}
