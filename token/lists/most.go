package lists

import (
	"bytes"
	"math"

	"github.com/zimmski/tavor/rand"
	"github.com/zimmski/tavor/token"
)

// Most implements a list token which holds a token that gets repeated 0 to maximum repetition value times
// Every permutation chooses a repetition value out of the range 0 to maximum repetition value and then generates the repetition of the token.
type Most struct {
	n     uint
	token token.Token
	value []token.Token
}

// NewMost returns a new instance of a Most token given the set of tokens and the maximum repetition value
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

// Clone returns a copy of the token and all its children
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

// Fuzz fuzzes this token using the random generator by choosing one of the possible permutations for this token
func (l *Most) Fuzz(r rand.Rand) {
	i := r.Int63n(int64(l.n + 1))

	l.permutation(uint(i))
}

// FuzzAll calls Fuzz for this token and then FuzzAll for all children of this token
func (l *Most) FuzzAll(r rand.Rand) {
	l.Fuzz(r)

	for _, tok := range l.value {
		tok.FuzzAll(r)
	}
}

// Parse tries to parse the token beginning from the current position in the parser data.
// If the parsing is successful the error argument is nil and the next current position after the token is returned.
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

// Permutation sets a specific permutation for this token
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

// Permutations returns the number of permutations for this token
func (l *Most) Permutations() uint {
	return l.n + 1
}

// PermutationsAll returns the number of all possible permutations for this token including its children
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
		if _, err := buffer.WriteString(tok.String()); err != nil {
			panic(err)
		}
	}

	return buffer.String()
}

// List interface methods

// Get returns the current referenced token at the given index. The error return argument is not nil, if the index is out of bound.
func (l *Most) Get(i int) (token.Token, error) {
	if i < 0 || i >= len(l.value) {
		return nil, &ListError{ListErrorOutOfBound}
	}

	return l.value[i], nil
}

// Len returns the number of the current referenced tokens
func (l *Most) Len() int {
	return len(l.value)
}

// InternalGet returns the current referenced internal token at the given index. The error return argument is not nil, if the index is out of bound.
func (l *Most) InternalGet(i int) (token.Token, error) {
	if i != 0 {
		return nil, &ListError{ListErrorOutOfBound}
	}

	return l.token, nil
}

// InternalLen returns the number of referenced internal tokens
func (l *Most) InternalLen() int {
	return 1
}

// InternalLogicalRemove removes the referenced internal token and returns the replacement for the current token or nil if the current token should be removed.
func (l *Most) InternalLogicalRemove(tok token.Token) token.Token {
	if l.token == tok {
		return nil
	}

	return l
}

// InternalReplace replaces an old with a new internal token if it is referenced by this token
func (l *Most) InternalReplace(oldToken, newToken token.Token) {
	if l.token == oldToken {
		l.token = newToken

		for i := range l.value {
			l.value[i] = l.token.Clone()
		}
	}
}

// OptionalToken interface methods

// IsOptional checks dynamically if this token is in the current state optional
func (l *Most) IsOptional() bool { return true }

// Activate activates this token
func (l *Most) Activate() {
	l.value = []token.Token{
		l.token.Clone(),
	}
}

// Deactivate deactivates this token
func (l *Most) Deactivate() {
	l.value = []token.Token{}
}
