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
	re := &Repeat{
		from:  from,
		to:    to,
		token: tok,
		value: make([]token.Token, from),
	}

	for i := range re.value {
		re.value[i] = tok.Clone()
	}

	return re
}

func (re *Repeat) Clone() token.Token {
	c := Repeat{
		from:  re.from,
		to:    re.to,
		token: re.token,
		value: make([]token.Token, len(re.value)),
	}

	for i, tok := range re.value {
		c.value[i] = tok.Clone()
	}

	return &c
}

func (re *Repeat) Fuzz(r rand.Rand) {
	n := r.Intn(int(re.to-re.from+1)) + int(re.from)
	toks := make([]token.Token, n)

	for i := range toks {
		toks[i] = re.token.Clone()

		toks[i].Fuzz(r)
	}

	re.value = toks
}

func (re *Repeat) String() string {
	var buffer bytes.Buffer

	for _, tok := range re.value {
		buffer.WriteString(tok.String())
	}

	return buffer.String()
}
