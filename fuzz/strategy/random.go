package strategy

import (
	"github.com/zimmski/container/list/linkedlist"

	"github.com/zimmski/tavor"
	"github.com/zimmski/tavor/log"
	"github.com/zimmski/tavor/rand"
	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/lists"
	"github.com/zimmski/tavor/token/sequences"
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
		tavor.ResetScope(s.root)
		s.fuzzYADDA(s.root, r)

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
	case token.List:
		l := t.Len()

		for i := 0; i < l; i++ {
			c, _ := t.Get(i)
			s.fuzz(c, r)
		}
	}
}

func (s *RandomStrategy) fuzzYADDA(root token.Token, r rand.Rand) {

	// TODO FIXME AND FIXME FIXME FIXME this should be done automatically somehow

	queue := linkedlist.New()

	queue.Push(root)

	for !queue.Empty() {
		t, _ := queue.Shift()
		tok := t.(token.Token)

		switch tok.(type) {
		case *sequences.SequenceExistingItem, *lists.UniqueItem:
			log.Debugf("fuzz again %#v(%p)", tok, tok)

			tok.Fuzz(r)
		}

		switch t := tok.(type) {
		case token.ForwardToken:
			if v := t.Get(); v != nil {
				queue.Push(v)
			}
		case token.List:
			for i := 0; i < t.Len(); i++ {
				c, _ := t.Get(i)

				queue.Push(c)
			}
		}
	}
}
