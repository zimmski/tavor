package strategy

import (
	"github.com/zimmski/tavor"
	"github.com/zimmski/tavor/log"
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

func (s *RandomStrategy) Fuzz(r rand.Rand) (chan struct{}, error) {
	if tavor.LoopExists(s.root) {
		return nil, &StrategyError{
			Message: "found endless loop in graph. Cannot proceed.",
			Type:    StrategyErrorEndlessLoopDetected,
		}
	}

	continueFuzzing := make(chan struct{})

	go func() {
		log.Debug("start random fuzzing routine")

		s.fuzz(s.root, r)

		tavor.ResetScope(s.root)
		tavor.ResetResetTokens(s.root)

		log.Debug("done with fuzzing step")

		// done with the last fuzzing step
		continueFuzzing <- struct{}{}

		log.Debug("finished fuzzing. Wait till the outside is ready to close.")

		if _, ok := <-continueFuzzing; ok {
			log.Debug("close fuzzing channel")

			close(continueFuzzing)
		}
	}()

	return continueFuzzing, nil
}

func (s *RandomStrategy) fuzz(tok token.Token, r rand.Rand) {
	tok.Fuzz(r)

	switch t := tok.(type) {
	case token.ForwardToken:
		if v := t.Get(); v != nil {
			s.fuzz(v, r)
		}
	case lists.List:
		l := t.Len()

		for i := 0; i < l; i++ {
			c, _ := t.Get(i)
			s.fuzz(c, r)
		}
	}
}
