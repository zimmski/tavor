package lists

import (
	"github.com/zimmski/tavor/rand"
	"github.com/zimmski/tavor/token"
)

type One struct {
	tokens []token.Token
	value  int
}

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

func (l *One) Fuzz(r rand.Rand) {
	i := r.Intn(len(l.tokens))

	l.permutation(uint(i))
}

func (l *One) FuzzAll(r rand.Rand) {
	l.Fuzz(r)

	l.tokens[l.value].FuzzAll(r)
}

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

func (l *One) Permutations() uint {
	return uint(len(l.tokens))
}

func (l *One) PermutationsAll() uint {
	var sum uint = 0

	for _, tok := range l.tokens {
		sum += tok.PermutationsAll()
	}

	return sum
}

func (l *One) String() string {
	return l.tokens[l.value].String()
}

// List interface methods

func (l *One) Get(i int) (token.Token, error) {
	if i != 0 {
		return nil, &ListError{ListErrorOutOfBound}
	}

	return l.tokens[l.value], nil
}

func (l *One) Len() int {
	return 1
}

func (l *One) InternalGet(i int) (token.Token, error) {
	if i < 0 || i >= len(l.tokens) {
		return nil, &ListError{ListErrorOutOfBound}
	}

	return l.tokens[i], nil
}

func (l *One) InternalLen() int {
	return len(l.tokens)
}

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

	switch len(l.tokens) {
	case 0:
		return nil
	case 1:
		return l.tokens[0]
	}

	if l.value == -1 {
		l.value = 0
	}

	return l
}

func (l *One) InternalReplace(oldToken, newToken token.Token) {
	for i := 0; i < len(l.tokens); i++ {
		if l.tokens[i] == oldToken {
			l.tokens[i] = newToken
		}
	}
}
