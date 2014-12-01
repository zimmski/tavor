package strategy

import (
	"math"

	"github.com/zimmski/container/list/linkedlist"

	"github.com/zimmski/tavor/log"
	"github.com/zimmski/tavor/rand"
	"github.com/zimmski/tavor/token"
)

// PermuteOptionalsStrategy implements a fuzzing strategy that generates permutations of only optional tokens of a token graph.
// Every iteration of the strategy generates a new permutation. The generation is deterministic. This strategy searches the graph for tokens who implement the OptionalToken interface and permutates over them by deactivating or activating them. The permutations always start from the deactivated states so that minimum data is generated first.
type PermuteOptionalsStrategy struct {
	root token.Token
}

// NewPermuteOptionalsStrategy returns a new instance of the Permute Optionals fuzzing strategy
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
			queue.Unshift(t.Get())
		case token.ListToken:
			for i := t.Len() - 1; i >= 0; i-- {
				c, _ := t.Get(i)

				queue.Unshift(c)
			}
		}
	} else {
		queue.Unshift(root)
	}

	for !queue.Empty() {
		tok, _ := queue.Shift()

		switch t := tok.(type) {
		case token.OptionalToken:
			if !t.IsOptional() {
				opts := s.findOptionals(r, t, true)

				if len(opts) > 0 {
					optionals = append(optionals, opts...)
				}

				continue
			}

			log.Debugf("found optional %#v", t)

			t.Deactivate()

			optionals = append(optionals, t)
		case token.ForwardToken:
			c := t.Get()

			if c != nil {
				err := c.Permutation(uint(r.Intn(int(c.Permutations())) + 1))
				if err != nil {
					log.Panic(err)
				}

				queue.Unshift(c)
			}
		case token.ListToken:
			for i := t.Len() - 1; i >= 0; i-- {
				c, _ := t.Get(i)

				err := c.Permutation(uint(r.Intn(int(c.Permutations())) + 1))
				if err != nil {
					log.Panic(err)
				}

				queue.Unshift(c)
			}
		}
	}

	return optionals
}

// Fuzz starts the first iteration of the fuzzing strategy returning a channel which controls the iteration flow.
// The channel returns a value if the iteration is complete and waits with calculating the next iteration until a value is put in. The channel is automatically closed when there are no more iterations. The error return argument is not nil if an error occurs during the setup of the fuzzing strategy.
func (s *PermuteOptionalsStrategy) Fuzz(r rand.Rand) (chan struct{}, error) {
	if token.LoopExists(s.root) {
		return nil, &Error{
			Message: "found endless loop in graph. Cannot proceed.",
			Type:    ErrorEndlessLoopDetected,
		}
	}

	continueFuzzing := make(chan struct{})

	go func() {
		log.Debug("start permute optionals routine")

		optionals := s.findOptionals(r, s.root, false)

		if len(optionals) > 0 {
			if !s.fuzz(r, continueFuzzing, optionals) {
				return
			}
		}

		token.ResetScope(s.root)
		token.ResetResetTokens(s.root)
		token.ResetScope(s.root)

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

func (s *PermuteOptionalsStrategy) fuzz(r rand.Rand, continueFuzzing chan struct{}, optionals []token.OptionalToken) bool {
	log.Debugf("fuzzing optionals %#v", optionals)

	// TODO make this WAYYYYYYYYYYYYYYYYYYYYYYYYYYYYYY smarter
	// since we can only fuzz 64 optionals at max
	// https://en.wikipedia.org/wiki/Steinhaus%E2%80%93Johnson%E2%80%93Trotter_algorithm
	p := 0
	maxSteps := int(math.Pow(2, float64(len(optionals))))

	for {
		log.Debugf("fuzzing step %b", p)

		ith := 1

		for i := range optionals {
			if p&ith == 0 {
				optionals[i].Deactivate()
			} else {
				optionals[i].Activate()

				children := s.findOptionals(r, optionals[i], true)

				if len(children) > 0 {
					if !s.fuzz(r, continueFuzzing, children) {
						return false
					}
				}
			}

			ith = ith << 1
		}

		p++

		if p == maxSteps {
			log.Debug("done with fuzzing these optionals")

			break
		}

		token.ResetScope(s.root)
		token.ResetResetTokens(s.root)
		token.ResetScope(s.root)

		log.Debug("done with fuzzing step")

		// done with this fuzzing step
		continueFuzzing <- struct{}{}

		// wait until we are allowed to continue
		if _, ok := <-continueFuzzing; !ok {
			log.Debug("fuzzing channel closed from outside")

			return false
		}
	}

	return true
}
