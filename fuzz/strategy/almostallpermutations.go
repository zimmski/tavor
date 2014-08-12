package strategy

import (
	"github.com/zimmski/container/list/linkedlist"

	"github.com/zimmski/tavor"
	"github.com/zimmski/tavor/log"
	"github.com/zimmski/tavor/rand"
	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/lists"
)

/*

this strategy does not cover all repititional permutations
this can be helpful when less permutations are needed

e.g. +2(?(1)?(2))
does not result in 16 permutations but just 7

*/

type almostAllPermutationsLevel struct {
	token           token.Token
	permutation     int
	maxPermutations int
}

type AlmostAllPermutationsStrategy struct {
	root token.Token

	resetedLookup map[token.Token]int
}

func NewAlmostAllPermutationsStrategy(tok token.Token) *AlmostAllPermutationsStrategy {
	s := &AlmostAllPermutationsStrategy{
		root: tok,

		resetedLookup: make(map[token.Token]int),
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

func (s *AlmostAllPermutationsStrategy) Fuzz(r rand.Rand) (chan struct{}, error) {
	if tavor.LoopExists(s.root) {
		return nil, &StrategyError{
			Message: "found endless loop in graph. Cannot proceed.",
			Type:    StrategyErrorEndlessLoopDetected,
		}
	}

	continueFuzzing := make(chan struct{})

	s.resetedLookup = make(map[token.Token]int)

	go func() {
		log.Debug("start almost all permutations routine")

		level := s.getLevel(s.root, false)

		if len(level) != 0 {
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

func (s *AlmostAllPermutationsStrategy) setTokenPermutation(tok token.Token, permutation int) {
	if per, ok := s.resetedLookup[tok]; ok && per == permutation {
		// Permutation already set in this step
	} else {
		tok.Permutation(permutation)

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
					s.setTokenPermutation(level[i+1].token, level[i+1].permutation)
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

			if t, ok := level[i].token.(token.OptionalToken); !ok || !t.IsOptional() || level[i].permutation != 1 {
				children := s.getLevel(level[i].token, true) // set all children to permutation 1

				if len(children) != 0 {
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

		s.resetedLookup = make(map[token.Token]int)
	}

	return true
}
