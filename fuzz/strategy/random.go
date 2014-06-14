package strategy

import (
	"fmt"

	"github.com/zimmski/tavor"
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

func (s *RandomStrategy) Fuzz(r rand.Rand) chan struct{} {
	continueFuzzing := make(chan struct{})

	go func() {
		if tavor.DEBUG {
			fmt.Println("Start random fuzzing routine")
		}

		s.fuzz(s.root, r)

		if tavor.DEBUG {
			fmt.Println("Done with fuzzing step")
		}

		// done with the last fuzzing step
		continueFuzzing <- struct{}{}

		if tavor.DEBUG {
			fmt.Println("Finished fuzzing. Wait till the outside is ready to close.")
		}

		if _, ok := <-continueFuzzing; ok {
			if tavor.DEBUG {
				fmt.Println("Close fuzzing channel")
			}

			close(continueFuzzing)
		}
	}()

	return continueFuzzing
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
