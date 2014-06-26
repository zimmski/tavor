package strategy

import (
	"math"

	"github.com/zimmski/container/list/linkedlist"

	"github.com/zimmski/tavor"
	"github.com/zimmski/tavor/log"
	"github.com/zimmski/tavor/rand"
	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/lists"
)

type PermuteOptionalsStrategy struct {
	root token.Token
}

func NewPermuteOptionalsStrategy(tok token.Token) *PermuteOptionalsStrategy {
	s := &PermuteOptionalsStrategy{
		root: tok,
	}

	return s
}

func init() {
	Register("PermuteOptionals", func(tok token.Token) Strategy {
		return NewPermuteOptionalsStrategy(tok)
	})
}

func (s *PermuteOptionalsStrategy) findOptionals(r rand.Rand, root token.Token, fromChildren bool) []token.OptionalToken {
	var optionals []token.OptionalToken
	var queue = linkedlist.New()

	if fromChildren {
		switch t := root.(type) {
		case token.ForwardToken:
			queue.Push(t.Get())
		case lists.List:
			l := t.Len()

			for i := 0; i < l; i++ {
				c, _ := t.Get(i)
				queue.Push(c)
			}
		}
	} else {
		queue.Push(root)
	}

	for !queue.Empty() {
		tok, _ := queue.Shift()

		switch t := tok.(type) {
		case token.OptionalToken:
			if !t.IsOptional() {
				opts := s.findOptionals(r, t, true)

				if len(opts) != 0 {
					optionals = append(optionals, opts...)
				}

				continue
			}

			log.Debugf("Found optional %#v", t)

			t.Deactivate()

			optionals = append(optionals, t)
		case token.ForwardToken:
			c := t.Get()

			c.Fuzz(r)

			queue.Push(c)
		case lists.List:
			l := t.Len()

			for i := 0; i < l; i++ {
				c, _ := t.Get(i)

				c.Fuzz(r)

				queue.Push(c)
			}
		}
	}

	return optionals
}

func (s *PermuteOptionalsStrategy) resetResetTokens() {
	var queue = linkedlist.New()

	queue.Push(s.root)

	for !queue.Empty() {
		v, _ := queue.Shift()

		switch tok := v.(type) {
		case token.ResetToken:
			log.Debugf("Reset %#v(%p)", tok, tok)

			tok.Reset()
		}

		switch tok := v.(type) {
		case token.ForwardToken:
			if v := tok.Get(); v != nil {
				queue.Push(v)
			}
		case lists.List:
			for i := 0; i < tok.Len(); i++ {
				c, _ := tok.Get(i)
				queue.Push(c)
			}
		}
	}
}

func (s *PermuteOptionalsStrategy) Fuzz(r rand.Rand) (chan struct{}, error) {
	if tavor.LoopExists(s.root) {
		return nil, &StrategyError{
			Message: "Found endless loop in graph. Cannot proceed.",
			Type:    StrategyErrorEndlessLoopDetected,
		}
	}

	continueFuzzing := make(chan struct{})

	go func() {
		log.Debug("Start permute optionals routine")

		optionals := s.findOptionals(r, s.root, false)

		if len(optionals) != 0 {
			if !s.fuzz(r, continueFuzzing, optionals) {
				return
			}
		}

		s.resetResetTokens()

		log.Debug("Done with fuzzing step")

		// done with the last fuzzing step
		continueFuzzing <- struct{}{}

		log.Debug("Finished fuzzing. Wait till the outside is ready to close.")

		if _, ok := <-continueFuzzing; ok {
			log.Debug("Close fuzzing channel")

			close(continueFuzzing)
		}
	}()

	return continueFuzzing, nil
}

func (s *PermuteOptionalsStrategy) fuzz(r rand.Rand, continueFuzzing chan struct{}, optionals []token.OptionalToken) bool {
	log.Debugf("Fuzzing optionals %#v", optionals)

	// TODO make this WAYYYYYYYYYYYYYYYYYYYYYYYYYYYYYY smarter
	// since we can only fuzz 64 optionals at max
	// https://en.wikipedia.org/wiki/Steinhaus%E2%80%93Johnson%E2%80%93Trotter_algorithm
	p := 0
	maxSteps := int(math.Pow(2, float64(len(optionals))))

	for {
		log.Debugf("Fuzzing step %b", p)

		ith := 1

		for i := range optionals {
			if p&ith == 0 {
				optionals[i].Deactivate()
			} else {
				optionals[i].Activate()

				children := s.findOptionals(r, optionals[i], true)

				if len(children) != 0 {
					if !s.fuzz(r, continueFuzzing, children) {
						return false
					}
				}
			}

			ith = ith << 1
		}

		p++

		if p == maxSteps {
			log.Debug("Done with fuzzing these optionals")

			break
		}

		s.resetResetTokens()

		log.Debug("Done with fuzzing step")

		// done with this fuzzing step
		continueFuzzing <- struct{}{}

		// wait until we are allowed to continue
		if _, ok := <-continueFuzzing; !ok {
			log.Debug("Fuzzing channel closed from outside")

			return false
		}
	}

	return true
}
