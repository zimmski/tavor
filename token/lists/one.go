package lists

import (
	"github.com/zimmski/tavor/token"
)

// One implements a list token which chooses of a set of referenced token exactly one token
// Every permutation chooses one token out of the token set.
type One struct {
	tokens []token.Token
	value  int
}

// NewOne returns a new instance of a One token given the set of tokens
func NewOne(toks ...token.Token) *One {
	if len(toks) == 0 {
		panic("at least one token needed")
	}

	return &One{
		tokens: toks,
		value:  0,
	}
}

// Token interface methods

// Clone returns a copy of the token and all its children
func (l *One) Clone() token.Token {
	c := One{
		tokens: make([]token.Token, len(l.tokens)),
		value:  l.value,
	}

	for i, tok := range l.tokens {
		c.tokens[i] = tok.Clone()
	}

	return &c
}

// Parse tries to parse the token beginning from the current position in the parser data.
// If the parsing is successful the error argument is nil and the next current position after the token is returned.
func (l *One) Parse(pars *token.InternalParser, cur int) (int, []error) {
	var nex int
	var es, errs []error

	for i := range l.tokens {
		nex, es = l.tokens[i].Parse(pars, cur)

		if len(es) == 0 {
			l.value = i

			return nex, nil
		}

		errs = append(errs, es...)
	}

	return cur, errs
}

func (l *One) permutation(i uint) {
	l.value = int(i)
}

// Permutation sets a specific permutation for this token
func (l *One) Permutation(i uint) error {
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
func (l *One) Permutations() uint {
	return uint(len(l.tokens))
}

// PermutationsAll returns the number of all possible permutations for this token including its children
func (l *One) PermutationsAll() uint {
	var sum uint

	for _, tok := range l.tokens {
		sum += tok.PermutationsAll()
	}

	return sum
}

func (l *One) String() string {
	return l.tokens[l.value].String()
}

// List interface methods

// Get returns the current referenced token at the given index. The error return argument is not nil, if the index is out of bound.
func (l *One) Get(i int) (token.Token, error) {
	if i != 0 {
		return nil, &ListError{ListErrorOutOfBound}
	}

	return l.tokens[l.value], nil
}

// Len returns the number of the current referenced tokens
func (l *One) Len() int {
	return 1
}

// InternalGet returns the current referenced internal token at the given index. The error return argument is not nil, if the index is out of bound.
func (l *One) InternalGet(i int) (token.Token, error) {
	if i < 0 || i >= len(l.tokens) {
		return nil, &ListError{ListErrorOutOfBound}
	}

	return l.tokens[i], nil
}

// InternalLen returns the number of referenced internal tokens
func (l *One) InternalLen() int {
	return len(l.tokens)
}

// InternalLogicalRemove removes the referenced internal token and returns the replacement for the current token or nil if the current token should be removed.
func (l *One) InternalLogicalRemove(tok token.Token) token.Token {
	for i := 0; i < len(l.tokens); i++ {
		if l.tokens[i] == tok {
			if l.value == i {
				l.value--
			}

			if i == len(l.tokens)-1 {
				l.tokens = l.tokens[:i]
			} else {
				l.tokens = append(l.tokens[:i], l.tokens[i+1:]...)
			}

			i--
		}
	}

	if l.value == -1 {
		l.value = 0
	}

	switch len(l.tokens) {
	case 0:
		return nil
	case 1:
		return l.tokens[0]
	}

	return l
}

// InternalReplace replaces an old with a new internal token if it is referenced by this token. The error return argument is not nil, if the replacement is not suitable.
func (l *One) InternalReplace(oldToken, newToken token.Token) error {
	for i := 0; i < len(l.tokens); i++ {
		if l.tokens[i] == oldToken {
			l.tokens[i] = newToken
		}
	}

	return nil
}

// Minimize interface methods

// Minimize tries to minimize itself and returns a token if it was successful, or nil if there was nothing to minimize
func (l *One) Minimize() token.Token {
	if len(l.tokens) == 1 {
		return l.tokens[0]
	}

	return nil
}
