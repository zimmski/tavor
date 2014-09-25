package lists

import (
	"bytes"
	"math"

	"github.com/zimmski/tavor/rand"
	"github.com/zimmski/tavor/token"
)

type Most struct {
	n     uint
	token token.Token
	value []token.Token
}

func NewMost(tok token.Token, n uint) *Most {
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

// Token interface methods

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
	i := r.Int63n(int64(l.n + 1))

	l.permutation(uint(i))
}

func (l *Most) FuzzAll(r rand.Rand) {
	l.Fuzz(r)

	for _, tok := range l.value {
		tok.FuzzAll(r)
	}
}

func (l *Most) Parse(pars *token.InternalParser, cur int) (int, []error) {
	panic("TODO implement")
}

func (l *Most) permutation(i uint) {
	toks := make([]token.Token, i)

	for i := range toks {
		toks[i] = l.token.Clone()
	}

	l.value = toks
}

func (l *Most) Permutation(i uint) error {
	permutations := l.Permutations()

	if i < 1 || i > permutations {
		return &token.PermutationError{
			Type: token.PermutationErrorIndexOutOfBound,
		}
	}

	l.permutation(i - 1)

	return nil
}

func (l *Most) Permutations() uint {
	return l.n + 1
}

func (l *Most) PermutationsAll() uint {
	var sum uint = 1

	tokenPermutations := l.token.PermutationsAll()

	for i := 1; i <= int(l.n); i++ {
		sum += uint(math.Pow(float64(tokenPermutations), float64(i)))
	}

	return sum
}

func (l *Most) String() string {
	var buffer bytes.Buffer

	for _, tok := range l.value {
		buffer.WriteString(tok.String())
	}

	return buffer.String()
}

// List interface methods

func (l *Most) Get(i int) (token.Token, error) {
	if i < 0 || i >= len(l.value) {
		return nil, &ListError{ListErrorOutOfBound}
	}

	return l.value[i], nil
}

func (l *Most) Len() int {
	return len(l.value)
}

func (l *Most) InternalGet(i int) (token.Token, error) {
	if i != 0 {
		return nil, &ListError{ListErrorOutOfBound}
	}

	return l.token, nil
}

func (l *Most) InternalLen() int {
	return 1
}

func (l *Most) InternalLogicalRemove(tok token.Token) token.Token {
	if l.token == tok {
		return nil
	}

	return l
}

func (l *Most) InternalReplace(oldToken, newToken token.Token) {
	if l.token == oldToken {
		l.token = newToken

		for i := range l.value {
			l.value[i] = l.token.Clone()
		}
	}
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
