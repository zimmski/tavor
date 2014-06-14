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
	}

	l.value = toks
}

func (l *Repeat) FuzzAll(r rand.Rand) {
	l.Fuzz(r)

	for _, tok := range l.value {
		tok.FuzzAll(r)
	}
}

func (l *Repeat) Get(i int) (token.Token, error) {
	if i < 0 || i >= len(l.value) {
		return nil, &ListError{ListErrorOutOfBound}
	}

	return l.value[i], nil
}

func (l *Repeat) Len() int {
	return len(l.value)
}

func (l *Repeat) Permutations() int {
	if l.from == 0 {
		return int(l.to-l.from)*l.token.Permutations() + 1
	}

	return int(l.to-l.from+1) * l.token.Permutations()
}

func (l *Repeat) String() string {
	var buffer bytes.Buffer

	for _, tok := range l.value {
		buffer.WriteString(tok.String())
	}

	return buffer.String()
}

// OptionalToken interface methods
func (l *Repeat) IsOptional() bool { return l.from == 0 }
func (l *Repeat) Activate() {
	if l.from != 0 {
		return
	}

	l.value = []token.Token{
		l.token.Clone(),
	}
}
func (l *Repeat) Deactivate() {
	if l.from != 0 {
		return
	}

	l.value = []token.Token{}
}
