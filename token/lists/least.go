package lists

import (
	"bytes"
	"math"

	"github.com/zimmski/tavor/rand"
	"github.com/zimmski/tavor/token"
)

type Least struct {
	n     int64
	token token.Token
	value []token.Token
}

func NewLeast(tok token.Token, n int64) *Least {
	l := &Least{
		n:     n,
		token: tok,
		value: make([]token.Token, n),
	}

	for i, _ := range l.value {
		l.value[i] = tok.Clone()
	}

	return l
}

func (l *Least) Clone() token.Token {
	c := Least{
		n:     l.n,
		token: l.token,
		value: make([]token.Token, len(l.value)),
	}

	for i, tok := range l.value {
		c.value[i] = tok.Clone()
	}

	return &c
}

func (l *Least) Fuzz(r rand.Rand) {
	n := int64(r.Intn(int(math.MaxInt64-l.n))) + l.n
	toks := make([]token.Token, int(n))

	for i, _ := range toks {
		toks[i] = l.token.Clone()

		toks[i].Fuzz(r)
	}

	l.value = toks
}

func (l *Least) String() string {
	var buffer bytes.Buffer

	for _, tok := range l.value {
		buffer.WriteString(tok.String())
	}

	return buffer.String()
}
