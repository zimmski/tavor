package lists

import (
	"bytes"

	"github.com/zimmski/tavor/token"
)

// Once implements a list token which holds a set of tokens that get shuffled on every permutation
type Once struct {
	tokens []token.Token
	values []int
}

// NewOnce returns a new instance of a Once token given the set of tokens
func NewOnce(toks ...token.Token) *Once {
	if len(toks) == 0 {
		panic("at least one token needed")
	}

	values := make([]int, len(toks))
	for i := 0; i < len(values); i++ {
		values[i] = i
	}

	return &Once{
		tokens: toks,
		values: values,
	}
}

// Token interface methods

// Clone returns a copy of the token and all its children
func (l *Once) Clone() token.Token {
	c := Once{
		tokens: make([]token.Token, len(l.tokens)),
		values: make([]int, len(l.values)),
	}

	for i, tok := range l.tokens {
		c.tokens[i] = tok.Clone()
	}

	for i, v := range l.values {
		c.values[i] = v
	}

	return &c
}

// Parse tries to parse the token beginning from the current position in the parser data.
// If the parsing is successful the error argument is nil and the next current position after the token is returned.
func (l *Once) Parse(pars *token.InternalParser, cur int) (int, []error) {
	panic("TODO implement")
}

func (l *Once) permutation(i uint) {
	le := len(l.tokens)
	rest := make([]int, le)
	for j := 0; j < len(rest); j++ {
		rest[j] = j
	}
	v := make([]int, 0, le)

	n := uint(le)
	for n > 0 {
		var pers uint = 1
		for j := uint(2); j <= n; j++ {
			pers *= j
		}
		split := pers / n

		ti := i / split
		i = i % split

		v = append(v, rest[ti])
		rest = append(rest[:ti], rest[ti+1:]...)

		n = uint(len(rest))
	}

	l.values = v
}

// Permutation sets a specific permutation for this token
func (l *Once) Permutation(i uint) error {
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
func (l *Once) Permutations() uint {
	var sum uint = 1

	for i := 2; i <= len(l.tokens); i++ {
		sum *= uint(i)
	}

	return sum
}

// PermutationsAll returns the number of all possible permutations for this token including its children
func (l *Once) PermutationsAll() uint {
	sum := l.Permutations()

	for _, tok := range l.tokens {
		sum *= tok.PermutationsAll()
	}

	return sum
}

func (l *Once) String() string {
	var buffer bytes.Buffer

	for i := range l.values {
		if _, err := buffer.WriteString(l.tokens[l.values[i]].String()); err != nil {
			panic(err)
		}
	}

	return buffer.String()
}

// List interface methods

// Get returns the current referenced token at the given index. The error return argument is not nil, if the index is out of bound.
func (l *Once) Get(i int) (token.Token, error) {
	if i < 0 || i >= len(l.values) {
		return nil, &ListError{ListErrorOutOfBound}
	}

	return l.tokens[l.values[i]], nil
}

// Len returns the number of the current referenced tokens
func (l *Once) Len() int {
	return len(l.values)
}

// InternalGet returns the current referenced internal token at the given index. The error return argument is not nil, if the index is out of bound.
func (l *Once) InternalGet(i int) (token.Token, error) {
	if i < 0 || i >= len(l.tokens) {
		return nil, &ListError{ListErrorOutOfBound}
	}

	return l.tokens[i], nil
}

// InternalLen returns the number of referenced internal tokens
func (l *Once) InternalLen() int {
	return len(l.tokens)
}

// InternalLogicalRemove removes the referenced internal token and returns the replacement for the current token or nil if the current token should be removed.
func (l *Once) InternalLogicalRemove(tok token.Token) token.Token {
	for i := 0; i < len(l.values); i++ {
		it := l.values[i]

		if l.tokens[it] == tok {
			l.tokens = append(l.tokens[:it], l.tokens[it+1:]...)
			l.values = append(l.values[:i], l.values[i+1:]...)

			for j := 0; j < len(l.values); j++ {
				if l.values[j] > it {
					l.values[j]--
				}
			}

			i--
		}
	}

	if len(l.values) == 0 {
		return nil
	}

	return l
}

// InternalReplace replaces an old with a new internal token if it is referenced by this token. The error return argument is not nil, if the replacement is not suitable.
func (l *Once) InternalReplace(oldToken, newToken token.Token) error {
	for i := 0; i < len(l.tokens); i++ {
		if l.tokens[i] == oldToken {
			l.tokens[i] = newToken
		}
	}

	return nil
}
