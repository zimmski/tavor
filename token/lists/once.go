package lists

import (
	"bytes"
	"math"

	"github.com/zimmski/tavor/rand"
	"github.com/zimmski/tavor/token"
)

type Once struct {
	tokens []token.Token
	values []int
}

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

func (l *Once) Fuzz(r rand.Rand) {
	ii := int64(l.Permutations())
	if ii < 0 { // TODO FIXME
		ii = math.MaxInt64
	}
	l.permutation(uint(r.Int63n(ii)))
}

func (l *Once) FuzzAll(r rand.Rand) {
	l.Fuzz(r)

	for i := range l.values {
		l.tokens[l.values[i]].FuzzAll(r)
	}
}

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

func (l *Once) Permutations() uint {
	var sum uint = 1

	for i := 2; i <= len(l.tokens); i++ {
		sum *= uint(i)
	}

	return sum
}

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

func (l *Once) Get(i int) (token.Token, error) {
	if i < 0 || i >= len(l.values) {
		return nil, &ListError{ListErrorOutOfBound}
	}

	return l.tokens[l.values[i]], nil
}

func (l *Once) Len() int {
	return len(l.values)
}

func (l *Once) InternalGet(i int) (token.Token, error) {
	if i < 0 || i >= len(l.tokens) {
		return nil, &ListError{ListErrorOutOfBound}
	}

	return l.tokens[i], nil
}

func (l *Once) InternalLen() int {
	return len(l.tokens)
}

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

func (l *Once) InternalReplace(oldToken, newToken token.Token) {
	for i := 0; i < len(l.tokens); i++ {
		if l.tokens[i] == oldToken {
			l.tokens[i] = newToken
		}
	}
}
