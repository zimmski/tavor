package lists

import (
	"bytes"

	"github.com/zimmski/tavor/rand"
	"github.com/zimmski/tavor/token"
)

// Many implements a list token which chooses some tokens out of set of tokens
// Every permutation chooses a count of 1 to length(set of tokens) and chooses this count of tokens out of the set of tokens.
type Many struct {
	tokens []token.Token
	value  []int
}

// NewMany returns a new instance of a Many token given the set of tokens
func NewMany(toks ...token.Token) *Many {
	if len(toks) == 0 {
		panic("at least one token needed")
	}

	return &Many{
		tokens: toks,
		value:  []int{0},
	}
}

// Token interface methods

// Clone returns a copy of the token and all its children
func (l *Many) Clone() token.Token {
	c := Many{
		tokens: make([]token.Token, len(l.tokens)),
		value:  make([]int, len(l.value)),
	}

	for i, tok := range l.tokens {
		c.tokens[i] = tok.Clone()
	}

	for i, v := range l.value {
		c.value[i] = v
	}

	return &c
}

// Fuzz fuzzes this token using the random generator by choosing one of the possible permutations for this token
func (l *Many) Fuzz(r rand.Rand) {
	tl := len(l.tokens)

	n := r.Intn(tl) + 1
	toks := make([]int, n)
	chosen := make(map[int]struct{})

	for i := range toks {
		for {
			ri := r.Intn(tl)

			if _, ok := chosen[ri]; !ok {
				toks[i] = ri
				chosen[ri] = struct{}{}

				break
			}
		}
	}

	l.value = toks
}

// FuzzAll calls Fuzz for this token and then FuzzAll for all children of this token
func (l *Many) FuzzAll(r rand.Rand) {
	l.Fuzz(r)

	for _, v := range l.value {
		l.tokens[v].FuzzAll(r)
	}
}

// Parse tries to parse the token beginning from the current position in the parser data.
// If the parsing is successful the error argument is nil and the next current position after the token is returned.
func (l *Many) Parse(pars *token.InternalParser, cur int) (int, []error) {
	panic("TODO implement")
}

// Permutation sets a specific permutation for this token
func (l *Many) Permutation(i uint) error {
	panic("TODO not implemented")
}

// Permutations returns the number of permutations for this token
func (l *Many) Permutations() uint {
	panic("TODO make this precise")
}

// PermutationsAll returns the number of all possible permutations for this token including its children
func (l *Many) PermutationsAll() uint {
	panic("TODO make this precise")
}

func (l *Many) String() string {
	var buffer bytes.Buffer

	for _, v := range l.value {
		if _, err := buffer.WriteString(l.tokens[v].String()); err != nil {
			panic(err)
		}
	}

	return buffer.String()
}

// List interface methods

// Get returns the current referenced token at the given index. The error return argument is not nil, if the index is out of bound.
func (l *Many) Get(i int) (token.Token, error) {
	if i < 0 || i >= len(l.value) {
		return nil, &ListError{ListErrorOutOfBound}
	}

	return l.tokens[l.value[i]], nil
}

// Len returns the number of the current referenced tokens
func (l *Many) Len() int {
	return len(l.value)
}

// InternalGet returns the current referenced internal token at the given index. The error return argument is not nil, if the index is out of bound.
func (l *Many) InternalGet(i int) (token.Token, error) {
	if i < 0 || i >= len(l.tokens) {
		return nil, &ListError{ListErrorOutOfBound}
	}

	return l.tokens[i], nil
}

// InternalLen returns the number of referenced internal tokens
func (l *Many) InternalLen() int {
	return len(l.tokens)
}

// InternalLogicalRemove removes the referenced internal token and returns the replacement for the current token or nil if the current token should be removed.
func (l *Many) InternalLogicalRemove(tok token.Token) token.Token {
	for i := 0; i < len(l.tokens); i++ {
		if l.tokens[i] == tok {
			for i, v := range l.value {
				if v == -1 {
					l.value[i]--
				}
			}

			if i == len(l.tokens)-1 {
				l.tokens = l.tokens[:i]
			} else {
				l.tokens = append(l.tokens[:i], l.tokens[i+1:]...)
			}

			i--
		}
	}

	if len(l.tokens) == 0 {
		return nil
	}

	for i, v := range l.value {
		if v == -1 {
			l.value[i] = 0
		}
	}

	return l
}

// InternalReplace replaces an old with a new internal token if it is referenced by this token
func (l *Many) InternalReplace(oldToken, newToken token.Token) {
	for i := 0; i < len(l.tokens); i++ {
		if l.tokens[i] == oldToken {
			l.tokens[i] = newToken
		}
	}
}
