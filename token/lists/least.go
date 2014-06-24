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

	for i := range l.value {
		l.value[i] = tok.Clone()
	}

	return l
}

// Token interface methods

func (l *Least) Clone() token.Token {
	c := Least{
		n:     l.n,
		token: l.token.Clone(),
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

	for i := range toks {
		toks[i] = l.token.Clone()
	}

	l.value = toks
}

func (l *Least) FuzzAll(r rand.Rand) {
	l.Fuzz(r)

	for _, tok := range l.value {
		tok.FuzzAll(r)
	}
}

func (l *Least) Permutation(i int) error {
	panic("TODO not implemented")
}

func (l *Least) Permutations() int {
	panic("TODO this might be hard to fit in 64bit")
}

func (l *Least) PermutationsAll() int {
	panic("TODO this might be hard to fit in 64bit")
}

func (l *Least) String() string {
	var buffer bytes.Buffer

	for _, tok := range l.value {
		buffer.WriteString(tok.String())
	}

	return buffer.String()
}

// List interface methods

func (l *Least) Get(i int) (token.Token, error) {
	if i < 0 || i >= len(l.value) {
		return nil, &ListError{ListErrorOutOfBound}
	}

	return l.value[i], nil
}

func (l *Least) Len() int {
	return len(l.value)
}

func (l *Least) LogicalRemove(tok token.Token) token.Token {
	if l.token == tok {
		return nil
	}

	return l
}

func (l *Least) Replace(oldToken, newToken token.Token) {
	if l.token == oldToken {
		l.token = newToken

		for i := range l.value {
			l.value[i] = l.token.Clone()
		}
	}
}

// OptionalToken interface methods

func (l *Least) IsOptional() bool { return l.n == 0 }
func (l *Least) Activate() {
	if l.n != 0 {
		return
	}

	l.value = []token.Token{
		l.token.Clone(),
	}
}
func (l *Least) Deactivate() {
	if l.n != 0 {
		return
	}

	l.value = []token.Token{}
}
