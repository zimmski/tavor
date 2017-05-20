package lists

import (
	"bytes"

	"github.com/zimmski/tavor/token"
)

// Concatenation implements a list token which holds an ordered set of tokens
type Concatenation struct {
	tokens []token.Token
}

// NewConcatenation returns a new instance of a Concatenation token given the set of tokens
func NewConcatenation(toks ...token.Token) *Concatenation {
	if len(toks) == 0 {
		panic("at least one token needed")
	}

	return &Concatenation{
		tokens: toks,
	}
}

// Token interface methods

// Clone returns a copy of the token and all its children
func (l *Concatenation) Clone() token.Token {
	c := Concatenation{
		tokens: make([]token.Token, len(l.tokens)),
	}

	for i, tok := range l.tokens {
		c.tokens[i] = tok.Clone()
	}

	return &c
}

// Parse tries to parse the token beginning from the current position in the parser data.
// If the parsing is successful the error argument is nil and the next current position after the token is returned.
func (l *Concatenation) Parse(pars *token.InternalParser, cur int) (int, []error) {
	for i := range l.tokens {
		nex, errs := l.tokens[i].Parse(pars, cur)

		if len(errs) > 0 {
			return nex, errs
		}

		cur = nex
	}

	return cur, nil
}

// Permutation sets a specific permutation for this token
func (l *Concatenation) Permutation(i uint) error {
	permutations := l.Permutations()

	if i < 0 || i >= permutations {
		return &token.PermutationError{
			Type: token.PermutationErrorIndexOutOfBound,
		}
	}

	// do nothing

	return nil
}

// Permutations returns the number of permutations for this token
func (l *Concatenation) Permutations() uint {
	return 1
}

// PermutationsAll returns the number of all possible permutations for this token including its children
func (l *Concatenation) PermutationsAll() uint {
	sum := l.Permutations()

	for _, tok := range l.tokens {
		sum *= tok.PermutationsAll()
	}

	return sum
}

func (l *Concatenation) String() string {
	var buffer bytes.Buffer

	for _, tok := range l.tokens {
		if _, err := buffer.WriteString(tok.String()); err != nil {
			panic(err)
		}
	}

	return buffer.String()
}

// List interface methods

// Get returns the current referenced token at the given index. The error return argument is not nil, if the index is out of bound.
func (l *Concatenation) Get(i int) (token.Token, error) {
	if i < 0 || i >= len(l.tokens) {
		return nil, &ListError{ListErrorOutOfBound}
	}

	return l.tokens[i], nil
}

// Len returns the number of the current referenced tokens
func (l *Concatenation) Len() int {
	return len(l.tokens)
}

// InternalGet returns the current referenced internal token at the given index. The error return argument is not nil, if the index is out of bound.
func (l *Concatenation) InternalGet(i int) (token.Token, error) {
	return l.Get(i)
}

// InternalLen returns the number of referenced internal tokens
func (l *Concatenation) InternalLen() int {
	return l.Len()
}

// InternalLogicalRemove removes the referenced internal token and returns the replacement for the current token or nil if the current token should be removed.
func (l *Concatenation) InternalLogicalRemove(tok token.Token) token.Token {
	for i := 0; i < len(l.tokens); i++ {
		if l.tokens[i] == tok {
			switch tok.(type) {
			case token.OptionalToken:
				if i == len(l.tokens)-1 {
					l.tokens = l.tokens[:i]
				} else {
					l.tokens = append(l.tokens[:i], l.tokens[i+1:]...)
				}
			default:
				// if we remove one non-optional token from an Concatenation list we have to remove everything
				return nil
			}
		}
	}

	if len(l.tokens) == 1 {
		return l.tokens[0]
	}

	return l
}

// InternalReplace replaces an old with a new internal token if it is referenced by this token. The error return argument is not nil, if the replacement is not suitable.
func (l *Concatenation) InternalReplace(oldToken, newToken token.Token) error {
	for i := 0; i < len(l.tokens); i++ {
		if l.tokens[i] == oldToken {
			l.tokens[i] = newToken
		}
	}

	return nil
}

// Minimize interface methods

// Minimize tries to minimize itself and returns a token if it was successful, or nil if there was nothing to minimize
func (l *Concatenation) Minimize() token.Token {
	if len(l.tokens) == 1 {
		return l.tokens[0]
	}

	return nil
}
