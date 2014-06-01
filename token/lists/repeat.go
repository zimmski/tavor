package lists

import (
	"bytes"

	"github.com/zimmski/tavor/rand"
	"github.com/zimmski/tavor/token"
)

type Repeat struct {
	from  int64
	to    int64
	token token.Token
	value []token.Token
}

func NewRepeat(tok token.Token, from, to int64) *Repeat {
	l := &Repeat{
		from:  from,
		to:    to,
		token: tok,
		value: make([]token.Token, from),
	}

	for i := range l.value {
		l.value[i] = tok.Clone()
	}

	return l
}

func (l *Repeat) Clone() token.Token {
	c := Repeat{
		from:  l.from,
		to:    l.to,
		token: l.token,
		value: make([]token.Token, len(l.value)),
	}

	for i, tok := range l.value {
		c.value[i] = tok.Clone()
	}

	return &c
}

func (l *Repeat) Fuzz(r rand.Rand) {
	n := r.Intn(int(l.to-l.from+1)) + int(l.from)
	toks := make([]token.Token, n)

	for i := range toks {
		toks[i] = l.token.Clone()

		toks[i].Fuzz(r)
	}

	l.value = toks
}

func (l *Repeat) Len() int {
	return len(l.value)
}

func (l *Repeat) String() string {
	var buffer bytes.Buffer

	for _, tok := range l.value {
		buffer.WriteString(tok.String())
	}

	return buffer.String()
}
