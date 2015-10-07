package strategy

import (
	"github.com/zimmski/tavor/log"
	"github.com/zimmski/tavor/rand"
	"github.com/zimmski/tavor/token"
)

type almostAllPermutationsLevel struct {
	parent      token.Token
	tokenIndex  int
	permutation uint
}

func (l almostAllPermutationsLevel) token() token.Token {
	if l.tokenIndex == -1 {
		return l.parent
	}

	switch t := l.parent.(type) {
	case token.ForwardToken:
		return t.Get()
	case token.ListToken:
		tt, _ := t.Get(l.tokenIndex)

		return tt
	}

	return nil
}

type almostAllPermutations struct {
	root token.Token

	resetedLookup map[token.Token]uint
	overextended  bool
}

func newAlmostAllPermutations(root token.Token) *almostAllPermutations {
	return &almostAllPermutations{
		root:          root,
		resetedLookup: make(map[token.Token]uint),
		overextended:  false,
	}
}

func init() {
	Register("AlmostAllPermutations", NewAlmostAllPermutations)
}

func (s *almostAllPermutations) getLevel(root token.Token, fromChildren bool) []almostAllPermutationsLevel {
	var level []almostAllPermutationsLevel

	if fromChildren {
		switch t := root.(type) {
		case token.ForwardToken:
			if v := t.Get(); v != nil {
				level = append(level, almostAllPermutationsLevel{
					parent:      root,
					tokenIndex:  0,
					permutation: 0,
				})
			}
		case token.ListToken:
			l := t.Len()

			for i := 0; i < l; i++ {
				level = append(level, almostAllPermutationsLevel{
					parent:      root,
					tokenIndex:  i,
					permutation: 0,
				})
			}
		}
	} else {
		level = append(level, almostAllPermutationsLevel{
			parent:      root,
			tokenIndex:  -1,
			permutation: 0,
		})
	}

	for _, l := range level {
		s.setTokenPermutation(l.token(), 0)
	}

	return level
}

// NewAlmostAllPermutations implements a fuzzing strategy that generates "almost" all possible permutations of a token graph.
// Every iteration of the strategy generates a new permutation. The generation is deterministic. This strategy does not cover all repititional permutations which can be helpful when less permutations are needed but a almost complete permutation coverage is still needed. For example the definition +2(?(1)?(2)) does not result in 16 permutations but instead it results in only 7.
func NewAlmostAllPermutations(root token.Token, r rand.Rand) (chan struct{}, error) {
	if token.LoopExists(root) {
		return nil, &Error{
			Message: "found endless loop in graph. Cannot proceed.",
			Type:    ErrEndlessLoopDetected,
		}
	}

	s := newAlmostAllPermutations(root)

	continueFuzzing := make(chan struct{})

	go func() {
		log.Debug("start almost all permutations routine")

		level := s.getLevel(s.root, false)

		if len(level) > 0 {
			log.Debug("start fuzzing step")

			if !s.fuzz(continueFuzzing, level) {
				return
			}
		}

		token.ResetCombinedScope(s.root)
		token.ResetResetTokens(s.root)
		token.ResetCombinedScope(s.root)

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

func (s *almostAllPermutations) setTokenPermutation(tok token.Token, permutation uint) {
	if per, ok := s.resetedLookup[tok]; ok && per == permutation {
		// Permutation already set in this step

		return
	}

	log.Debugf("set %#v(%p) to permutation %d of max permutations %d", tok, tok, permutation, tok.Permutations())

	if err := tok.Permutation(permutation); err != nil {
		panic(err)
	}

	s.resetedLookup[tok] = permutation
}

func (s *almostAllPermutations) fuzz(continueFuzzing chan struct{}, level []almostAllPermutationsLevel) bool {
	log.Debugf("fuzzing level %d->%#v", len(level), level)

	last := len(level) - 1

STEP:
	for {
		for i := range level {
			if level[i].permutation >= level[i].token().Permutations() {
				if i < last {
					log.Debugf("max reached redo everything <= %d and increment next", i)

					if level[i].token().Permutations() != 0 {
						log.Debug("Let's stay here")

						s.overextended = false
					}

					level[i+1].permutation++
					if level[i+1].permutation < level[i+1].token().Permutations() {
						s.setTokenPermutation(level[i+1].token(), level[i+1].permutation)
					}
					s.getLevel(level[i+1].token(), true) // set all children to permutation 0
				} else {
					log.Debug("Overextended our stay, let's get out of here!")

					s.overextended = true

					break STEP
				}

				for k := 0; k <= i; k++ {
					level[k].permutation = 0
					s.setTokenPermutation(level[k].token(), 0)
					s.getLevel(level[k].token(), true) // set all children to permutation 0
				}

				continue STEP
			}

			log.Debugf("permute %d->%#v", i, level[i])

			s.setTokenPermutation(level[i].token(), level[i].permutation)

			if t, ok := level[i].token().(token.OptionalToken); !ok || !t.IsOptional() || level[i].permutation > 0 {
				children := s.getLevel(level[i].token(), true) // set all children to permutation 0

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

		if level[0].permutation >= level[0].token().Permutations() {
			found := false
			for i := 1; i < len(level); i++ {
				if level[i].permutation < level[i].token().Permutations()-1 {
					found = true

					break
				}
			}
			if !found {
				break STEP
			}
		}

		if !s.overextended {
			token.ResetCombinedScope(s.root)
			token.ResetResetTokens(s.root)
			token.ResetCombinedScope(s.root)

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
	}

	log.Debug("done with fuzzing this level")

	return true
}
