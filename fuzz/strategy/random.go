package strategy

import (
	"github.com/zimmski/tavor/log"
	"github.com/zimmski/tavor/rand"
	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/sequences"
)

func init() {
	Register("random", NewRandom)
}

type random struct {
	root token.Token
}

// NewRandom implements a fuzzing strategy that generates a random permutation of a token graph.
// The strategy does exactly one iteration which permutates at random all reachable tokens in the graph. The determinism is dependent on the random generator and is therefore for example deterministic if a seed for the random generator produces always the same outputs.
func NewRandom(root token.Token, r rand.Rand) (chan struct{}, error) {
	if r == nil {
		return nil, &Error{
			Message: "random generator is nil",
			Type:    ErrNilRandomGenerator,
		}
	}

	if token.LoopExists(root) {
		return nil, &Error{
			Message: "found endless loop in graph. Cannot proceed.",
			Type:    ErrEndlessLoopDetected,
		}
	}

	s := &random{
		root: root,
	}

	continueFuzzing := make(chan struct{})

	go func() {
		log.Debug("start random fuzzing routine")

		s.fuzz(s.root, r, token.NewVariableScope())

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

func (s *random) fuzz(tok token.Token, r rand.Rand, variableScope *token.VariableScope) {
	log.Debugf("Fuzz (%p)%#v with maxPermutations %d", tok, tok, tok.Permutations())

	if t, ok := tok.(token.Scoping); ok && t.Scoping() {
		variableScope = variableScope.Push()
	}

	err := tok.Permutation(uint(r.Int63n(int64(tok.Permutations()))))
	if err != nil {
		log.Panic(err)
	}

	if t, ok := tok.(token.Follow); !ok || t.Follow() {
		switch t := tok.(type) {
		case token.ForwardToken:
			if v := t.Get(); v != nil {
				s.fuzz(v, r, variableScope)
			}
		case token.ListToken:
			l := t.Len()

			for i := 0; i < l; i++ {
				c, _ := t.Get(i)
				s.fuzz(c, r, variableScope)
			}
		}
	}

	if t, ok := tok.(token.Scoping); ok && t.Scoping() {
		variableScope = variableScope.Pop()
	}
}

func (s *random) fuzzYADDA(root token.Token, r rand.Rand) {
	// TODO FIXME AND FIXME FIXME FIXME this should be done automatically somehow
	// since this doesn't work in other heuristics...
	// especially the fuzz again part is tricky. the whole reason is because of dynamic repeats that clone during a reset. so the "reset" or regenerating of new child tokens has to be done better

	token.ResetCombinedScope(root)
	token.ResetResetTokens(root)
	token.ResetCombinedScope(root)

	err := token.Walk(root, func(tok token.Token) error {
		switch tok.(type) {
		case *sequences.SequenceExistingItem:
			log.Debugf("Fuzz again %p(%#v)", tok, tok)

			err := tok.Permutation(uint(r.Int63n(int64(tok.Permutations()))))
			if err != nil {
				log.Panic(err)
			}
		}

		return nil
	})
	if err != nil {
		panic(err)
	}
}
