package strategy

import (
	"github.com/zimmski/container/list/linkedlist"

	"github.com/zimmski/tavor"
	"github.com/zimmski/tavor/log"
	"github.com/zimmski/tavor/rand"
	"github.com/zimmski/tavor/token"
)

type almostAllPermutationsLevel struct {
	token           token.Token
	permutation     uint
	maxPermutations uint
}

// AlmostAllPermutationsStrategy implements a fuzzing strategy that generates "almost" all possible permutations of a token graph.
// Every iteration of the strategy generates a new permutation. The generation is deterministically. This strategy does not cover all repititional permutations which can be helpful when less permutations are needed but a almost complete permutation coverage is still needed. For example the definition +2(?(1)?(2)) does not result in 16 permutations but instead it results in only 7.
type AlmostAllPermutationsStrategy struct {
	root token.Token

	resetedLookup map[token.Token]uint
}

// NewAlmostAllPermutationsStrategy returns a new instance of the Almost All Permutations fuzzing strategy
func NewAlmostAllPermutationsStrategy(tok token.Token) *AlmostAllPermutationsStrategy {
	s := &AlmostAllPermutationsStrategy{
		root: tok,

		resetedLookup: make(map[token.Token]uint),
	}

	return s
}

func init() {
	Register("AlmostAllPermutations", func(tok token.Token) Strategy {
		return NewAlmostAllPermutationsStrategy(tok)
	})
}

func (s *AlmostAllPermutationsStrategy) getLevel(root token.Token, fromChildren bool) []almostAllPermutationsLevel {
	var level []almostAllPermutationsLevel
	var queue = linkedlist.New()

	if fromChildren {
		switch t := root.(type) {
		case token.ForwardToken:
			if v := t.Get(); v != nil {
				queue.Push(v)
			}
		case token.List:
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
		v, _ := queue.Shift()
		tok, _ := v.(token.Token)

		s.setTokenPermutation(tok, 1)

		level = append(level, almostAllPermutationsLevel{
			token:           tok,
			permutation:     1,
			maxPermutations: tok.Permutations(),
		})
	}

	return level
}

// Fuzz starts the first iteration of the fuzzing strategy returning a channel which controls the iteration flow.
// The channel returns a value if the iteration is complete and waits with calculating the next iteration until a value is put in. The channel is automatically closed when there are no more iterations. The error return argument is not nil if an error occurs during the setup of the fuzzing strategy.
func (s *AlmostAllPermutationsStrategy) Fuzz(r rand.Rand) (chan struct{}, error) {
	if tavor.LoopExists(s.root) {
		return nil, &Error{
			Message: "found endless loop in graph. Cannot proceed.",
			Type:    ErrorEndlessLoopDetected,
		}
	}

	continueFuzzing := make(chan struct{})

	s.resetedLookup = make(map[token.Token]uint)

	go func() {
		log.Debug("start almost all permutations routine")

		level := s.getLevel(s.root, false)

		if len(level) > 0 {
			log.Debug("start fuzzing step")

			if !s.fuzz(continueFuzzing, level) {
				return
			}
		}

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

func (s *AlmostAllPermutationsStrategy) setTokenPermutation(tok token.Token, permutation uint) {
	if per, ok := s.resetedLookup[tok]; ok && per == permutation {
		// Permutation already set in this step
	} else {
		if err := tok.Permutation(permutation); err != nil {
			panic(err)
		}

		s.resetedLookup[tok] = permutation
	}
}

func (s *AlmostAllPermutationsStrategy) fuzz(continueFuzzing chan struct{}, level []almostAllPermutationsLevel) bool {
	log.Debugf("fuzzing level %d->%#v", len(level), level)

	last := len(level) - 1

STEP:
	for {
		for i := range level {
			if level[i].permutation > level[i].maxPermutations {
				if i <= last {
					log.Debugf("max reached redo everything <= %d and increment next", i)

					level[i+1].permutation++
					if level[i+1].permutation <= level[i+1].maxPermutations {
						s.setTokenPermutation(level[i+1].token, level[i+1].permutation)
					}
					s.getLevel(level[i+1].token, true) // set all children to permutation 1
				}

				for k := 0; k <= i; k++ {
					level[k].permutation = 1
					s.setTokenPermutation(level[k].token, 1)
					s.getLevel(level[k].token, true) // set all children to permutation 1
				}

				continue STEP
			}

			log.Debugf("permute %d->%#v", i, level[i])

			s.setTokenPermutation(level[i].token, level[i].permutation)

			if t, ok := level[i].token.(token.OptionalToken); !ok || !t.IsOptional() || level[i].permutation > 1 {
				children := s.getLevel(level[i].token, true) // set all children to permutation 1

				if len(children) > 0 {
					if !s.fuzz(continueFuzzing, children) {
						return false
					}
				}
			}

			if i == 0 {
				level[i].permutation++
			}
		}

		if level[0].permutation > level[0].maxPermutations {
			found := false
			for i := 1; i < len(level); i++ {
				if level[i].permutation < level[i].maxPermutations {
					found = true

					break
				}
			}
			if !found {
				log.Debug("done with fuzzing this level")

				break STEP
			}
		}

		tavor.ResetScope(s.root)
		tavor.ResetResetTokens(s.root)

		log.Debug("done with fuzzing step")

		// done with this fuzzing step
		continueFuzzing <- struct{}{}

		// wait until we are allowed to continue
		if _, ok := <-continueFuzzing; !ok {
			log.Debug("fuzzing channel closed from outside")

			return false
		}

		log.Debug("start fuzzing step")

		s.resetedLookup = make(map[token.Token]uint)
	}

	return true
}
