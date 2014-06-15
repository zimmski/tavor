package lists

import (
	"github.com/zimmski/tavor/rand"
	"github.com/zimmski/tavor/token"
)

type One struct {
	tokens []token.Token
	value  token.Token
}

func NewOne(toks ...token.Token) *One {
	if len(toks) == 0 {
		panic("at least one token needed")
	}

	return &One{
		tokens: toks,
		value:  toks[0],
	}
}

func (l *One) Clone() token.Token {
	c := One{
		tokens: make([]token.Token, len(l.tokens)),
		value:  l.value.Clone(),
	}

	for i, tok := range l.tokens {
		c.tokens[i] = tok.Clone()
	}

	return &c
}

func (l *One) Fuzz(r rand.Rand) {
	i := r.Intn(len(l.tokens))

	l.permutation(i)
}

func (l *One) FuzzAll(r rand.Rand) {
	l.Fuzz(r)

	l.value.FuzzAll(r)
}

func (l *One) Get(i int) (token.Token, error) {
	if i != 0 {
		return nil, &ListError{ListErrorOutOfBound}
	}

	return l.value, nil
}

func (l *One) Len() int {
	return 1
}

func (l *One) permutation(i int) {
	l.value = l.tokens[i]
}

func (l *One) Permutation(i int) error {
	permutations := l.Permutations()

	if i < 1 || i > permutations {
		return &token.PermutationError{
			Type: token.PermutationErrorIndexOutOfBound,
		}
	}

	l.permutation(i - 1)

	return nil
}

func (l *One) Permutations() int {
	return len(l.tokens)
}

func (l *One) PermutationsAll() int {
	sum := 0

	for _, tok := range l.tokens {
		sum += tok.PermutationsAll()
	}

	return sum
}

func (l *One) String() string {
	return l.value.String()
}
