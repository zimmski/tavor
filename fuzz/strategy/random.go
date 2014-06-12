package strategy

import (
	"github.com/zimmski/tavor/rand"
	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/lists"
)

type RandomStrategy struct {
	root token.Token
}

func NewRandomStrategy(tok token.Token) *RandomStrategy {
	return &RandomStrategy{
		root: tok,
	}
}

func init() {
	Register("random", func(tok token.Token) Strategy {
		return NewRandomStrategy(tok)
	})
}

func (s *RandomStrategy) Fuzz(r rand.Rand) {
	s.fuzz(s.root, r)
}

func (s *RandomStrategy) fuzz(tok token.Token, r rand.Rand) {
	tok.Fuzz(r)

	switch t := tok.(type) {
	case token.ForwardToken:
		s.fuzz(t.Get(), r)
	case lists.List:
		l := t.Len()

		for i := 0; i < l; i++ {
			c, _ := t.Get(i)
			s.fuzz(c, r)
		}
	}
}
