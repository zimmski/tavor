package lists

import (
	"bytes"

	"github.com/zimmski/tavor/rand"
	"github.com/zimmski/tavor/token"
)

type Many struct {
	tokens []token.Token
	value  []token.Token
}

func NewMany(toks ...token.Token) *Many {
	if len(toks) == 0 {
		panic("at least one token needed")
	}

	return &Many{
		tokens: toks,
		value:  toks,
	}
}

func (l *Many) Clone() token.Token {
	c := Many{
		tokens: make([]token.Token, len(l.tokens)),
		value:  make([]token.Token, len(l.value)),
	}

	for i, tok := range l.tokens {
		c.tokens[i] = tok.Clone()
	}

	for i, tok := range l.value {
		c.value[i] = tok.Clone()
	}

	return &c
}

func (l *Many) Fuzz(r rand.Rand) {
	tl := len(l.tokens)

	n := r.Intn(tl) + 1
	toks := make([]token.Token, n)
	chosen := make(map[int]struct{})

	for i := range toks {
		for {
			ri := r.Intn(tl)

			if _, ok := chosen[ri]; !ok {
				toks[i] = l.value[ri]
				chosen[ri] = struct{}{}

				toks[i].Fuzz(r)

				break
			}
		}
	}

	l.value = toks
}

func (l *Many) Get(i int) (token.Token, error) {
	if i < 0 || i >= len(l.tokens) {
		return nil, &ListError{ListErrorOutOfBound}
	}

	return l.tokens[i], nil
}

func (l *Many) Len() int {
	return len(l.value)
}

func (l *Many) Permutations() int {
	panic("TODO make this precise")
}

func (l *Many) String() string {
	var buffer bytes.Buffer

	for _, tok := range l.value {
		buffer.WriteString(tok.String())
	}

	return buffer.String()
}
