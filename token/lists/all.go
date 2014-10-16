package lists

import (
	"bytes"

	"github.com/zimmski/tavor/rand"
	"github.com/zimmski/tavor/token"
)

type All struct {
	tokens []token.Token
}

func NewAll(toks ...token.Token) *All {
	if len(toks) == 0 {
		panic("at least one token needed")
	}

	return &All{
		tokens: toks,
	}
}

// Token interface methods

func (l *All) Clone() token.Token {
	c := All{
		tokens: make([]token.Token, len(l.tokens)),
	}

	for i, tok := range l.tokens {
		c.tokens[i] = tok.Clone()
	}

	return &c
}

func (l *All) Fuzz(r rand.Rand) {
	// do nothing
}

func (l *All) FuzzAll(r rand.Rand) {
	l.Fuzz(r)

	for _, tok := range l.tokens {
		tok.FuzzAll(r)
	}
}

func (l *All) Parse(pars *token.InternalParser, cur int) (int, []error) {
	for i := range l.tokens {
		nex, errs := l.tokens[i].Parse(pars, cur)

		if len(errs) > 0 {
			return nex, errs
		}

		cur = nex
	}

	return cur, nil
}

func (l *All) Permutation(i uint) error {
	permutations := l.Permutations()

	if i < 1 || i > permutations {
		return &token.PermutationError{
			Type: token.PermutationErrorIndexOutOfBound,
		}
	}

	// do nothing

	return nil
}

func (l *All) Permutations() uint {
	return 1
}

func (l *All) PermutationsAll() uint {
	sum := l.Permutations()

	for _, tok := range l.tokens {
		sum *= tok.PermutationsAll()
	}

	return sum
}

func (l *All) String() string {
	var buffer bytes.Buffer

	for _, tok := range l.tokens {
		if _, err := buffer.WriteString(tok.String()); err != nil {
			panic(err)
		}
	}

	return buffer.String()
}

// List interface methods

func (l *All) Get(i int) (token.Token, error) {
	if i < 0 || i >= len(l.tokens) {
		return nil, &ListError{ListErrorOutOfBound}
	}

	return l.tokens[i], nil
}

func (l *All) Len() int {
	return len(l.tokens)
}

func (l *All) InternalGet(i int) (token.Token, error) {
	return l.Get(i)
}

func (l *All) InternalLen() int {
	return l.Len()
}

func (l *All) InternalLogicalRemove(tok token.Token) token.Token {
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
				// if we remove one token from an All list we have to remove everything
				return nil
			}
		}
	}

	return l
}

func (l *All) InternalReplace(oldToken, newToken token.Token) {
	for i := 0; i < len(l.tokens); i++ {
		if l.tokens[i] == oldToken {
			l.tokens[i] = newToken
		}
	}
}
