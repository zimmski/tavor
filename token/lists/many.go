package lists

import (
	"bytes"

	"github.com/zimmski/tavor/rand"
	"github.com/zimmski/tavor/token"
)

type Many struct {
	tokens []token.Token
	value  []int
}

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

func (l *Many) FuzzAll(r rand.Rand) {
	l.Fuzz(r)

	for _, v := range l.value {
		l.tokens[v].FuzzAll(r)
	}
}

func (l *Many) Parse(pars *token.InternalParser, cur int) (int, []error) {
	panic("TODO implement")
}

func (l *Many) Permutation(i uint) error {
	panic("TODO not implemented")
}

func (l *Many) Permutations() uint {
	panic("TODO make this precise")
}

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

func (l *Many) Get(i int) (token.Token, error) {
	if i < 0 || i >= len(l.value) {
		return nil, &ListError{ListErrorOutOfBound}
	}

	return l.tokens[l.value[i]], nil
}

func (l *Many) Len() int {
	return len(l.value)
}

func (l *Many) InternalGet(i int) (token.Token, error) {
	if i < 0 || i >= len(l.tokens) {
		return nil, &ListError{ListErrorOutOfBound}
	}

	return l.tokens[i], nil
}

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
