package lists

import (
	"bytes"
	"math"

	"github.com/zimmski/tavor/rand"
	"github.com/zimmski/tavor/token"
)

// Least implements a list token which holds a token that gets repeated minimum repetition value to infinite times
// Every permutation chooses a repetition value out of the range minimum repetition value to infinite and then generates the repetition of the token.
type Least struct {
	n     int64
	token token.Token
	value []token.Token
}

// NewLeast returns a new instance of a Least token given the set of tokens and the minimum repetition value
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

// Clone returns a copy of the token and all its children
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

// Fuzz fuzzes this token using the random generator by choosing one of the possible permutations for this token
func (l *Least) Fuzz(r rand.Rand) {
	n := int64(r.Intn(int(math.MaxInt64-l.n))) + l.n
	toks := make([]token.Token, int(n))

	for i := range toks {
		toks[i] = l.token.Clone()
	}

	l.value = toks
}

// FuzzAll calls Fuzz for this token and then FuzzAll for all children of this token
func (l *Least) FuzzAll(r rand.Rand) {
	l.Fuzz(r)

	for _, tok := range l.value {
		tok.FuzzAll(r)
	}
}

// Parse tries to parse the token beginning from the current position in the parser data.
// If the parsing is successful the error argument is nil and the next current position after the token is returned.
func (l *Least) Parse(pars *token.InternalParser, cur int) (int, []error) {
	panic("TODO implement")
}

// Permutation sets a specific permutation for this token
func (l *Least) Permutation(i uint) error {
	panic("TODO not implemented")
}

// Permutations returns the number of permutations for this token
func (l *Least) Permutations() uint {
	panic("TODO this might be hard to fit in 64bit")
}

// PermutationsAll returns the number of all possible permutations for this token including its children
func (l *Least) PermutationsAll() uint {
	panic("TODO this might be hard to fit in 64bit")
}

func (l *Least) String() string {
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
func (l *Least) Get(i int) (token.Token, error) {
	if i < 0 || i >= len(l.value) {
		return nil, &ListError{ListErrorOutOfBound}
	}

	return l.value[i], nil
}

// Len returns the number of the current referenced tokens
func (l *Least) Len() int {
	return len(l.value)
}

// InternalGet returns the current referenced internal token at the given index. The error return argument is not nil, if the index is out of bound.
func (l *Least) InternalGet(i int) (token.Token, error) {
	if i != 0 {
		return nil, &ListError{ListErrorOutOfBound}
	}

	return l.token, nil
}

// InternalLen returns the number of referenced internal tokens
func (l *Least) InternalLen() int {
	return 1
}

// InternalLogicalRemove removes the referenced internal token and returns the replacement for the current token or nil if the current token should be removed.
func (l *Least) InternalLogicalRemove(tok token.Token) token.Token {
	if l.token == tok {
		return nil
	}

	return l
}

// InternalReplace replaces an old with a new internal token if it is referenced by this token
func (l *Least) InternalReplace(oldToken, newToken token.Token) {
	if l.token == oldToken {
		l.token = newToken

		for i := range l.value {
			l.value[i] = l.token.Clone()
		}
	}
}

// OptionalToken interface methods

// IsOptional checks dynamically if this token is in the current state optional
func (l *Least) IsOptional() bool { return l.n == 0 }

// Activate activates this token
func (l *Least) Activate() {
	if l.n > 0 {
		return
	}

	l.value = []token.Token{
		l.token.Clone(),
	}
}

// Deactivate deactivates this token
func (l *Least) Deactivate() {
	if l.n > 0 {
		return
	}

	l.value = []token.Token{}
}
