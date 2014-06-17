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
	i := r.Intn(int(l.n) + 1)

	l.permutation(i)
}

func (l *Most) FuzzAll(r rand.Rand) {
	l.Fuzz(r)

	for _, tok := range l.value {
		tok.FuzzAll(r)
	}
}

func (l *Most) Get(i int) (token.Token, error) {
	if i < 0 || i >= len(l.value) {
		return nil, &ListError{ListErrorOutOfBound}
	}

	return l.value[i], nil
}

func (l *Most) Len() int {
	return len(l.value)
}

func (l *Most) permutation(i int) {
	toks := make([]token.Token, i)

	for i := range toks {
		toks[i] = l.token.Clone()
	}

	l.value = toks
}

func (l *Most) Permutation(i int) error {
	permutations := l.Permutations()

	if i < 1 || i > permutations {
		return &token.PermutationError{
			Type: token.PermutationErrorIndexOutOfBound,
		}
	}

	l.permutation(i - 1)

	return nil
}

func (l *Most) Permutations() int {
	return int(l.n) + 1
}

func (l *Most) PermutationsAll() int {
	return int(l.n)*l.token.PermutationsAll() + 1
}

func (l *Most) String() string {
	var buffer bytes.Buffer

	for _, tok := range l.value {
		buffer.WriteString(tok.String())
	}

	return buffer.String()
}

// OptionalToken interface methods

func (l *Most) IsOptional() bool { return true }
func (l *Most) Activate() {
	l.value = []token.Token{
		l.token.Clone(),
	}
}
func (l *Most) Deactivate() {
	l.value = []token.Token{}
}
