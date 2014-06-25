package lists

import (
	"bytes"
	"math"

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

// Token interface methods

func (l *Repeat) Clone() token.Token {
	c := Repeat{
		from:  l.from,
		to:    l.to,
		token: l.token.Clone(),
		value: make([]token.Token, len(l.value)),
	}

	for i, tok := range l.value {
		c.value[i] = tok.Clone()
	}

	return &c
}

func (l *Repeat) Fuzz(r rand.Rand) {
	i := r.Intn(int(l.to - l.from + 1))

	l.permutation(i)
}

func (l *Repeat) FuzzAll(r rand.Rand) {
	l.Fuzz(r)

	for _, tok := range l.value {
		tok.FuzzAll(r)
	}
}

func (l *Repeat) permutation(i int) {
	toks := make([]token.Token, i+int(l.from))

	for i := range toks {
		toks[i] = l.token.Clone()
	}

	l.value = toks
}

func (l *Repeat) Permutation(i int) error {
	permutations := l.Permutations()

	if i < 1 || i > permutations {
		return &token.PermutationError{
			Type: token.PermutationErrorIndexOutOfBound,
		}
	}

	l.permutation(i - 1)

	return nil
}

func (l *Repeat) Permutations() int {
	return int(l.to - l.from + 1)
}

func (l *Repeat) PermutationsAll() int {
	sum := 0
	from := l.from

	if l.from == 0 {
		sum++
		from++
	}

	tokenPermutations := l.token.PermutationsAll()

	for i := from; i <= l.to; i++ {
		sum += int(math.Pow(float64(tokenPermutations), float64(i)))
	}

	return sum
}

func (l *Repeat) String() string {
	var buffer bytes.Buffer

	for _, tok := range l.value {
		buffer.WriteString(tok.String())
	}

	return buffer.String()
}

// List interface methods

func (l *Repeat) Get(i int) (token.Token, error) {
	if i < 0 || i >= len(l.value) {
		return nil, &ListError{ListErrorOutOfBound}
	}

	return l.value[i], nil
}

func (l *Repeat) Len() int {
	return len(l.value)
}

func (l *Repeat) InternalGet(i int) (token.Token, error) {
	if i != 1 {
		return nil, &ListError{ListErrorOutOfBound}
	}

	return l.token, nil
}

func (l *Repeat) InternalLen() int {
	return 1
}

func (l *Repeat) InternalLogicalRemove(tok token.Token) token.Token {
	if l.token == tok {
		return nil
	}

	return l
}

func (l *Repeat) InternalReplace(oldToken, newToken token.Token) {
	if l.token == oldToken {
		l.token = newToken

		for i := range l.value {
			l.value[i] = l.token.Clone()
		}
	}
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
