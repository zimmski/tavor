package strategy

import (
	"github.com/zimmski/tavor/log"
	"github.com/zimmski/tavor/rand"
	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/sequences"
)

// RandomStrategy implements a fuzzing strategy that generates a random permutation of a token graph.
// The strategy does exactly one iteration which permutates at random all reachable tokens in the graph. The determinism is dependent on the random generator and is therefore for example deterministic if a seed for the random generator produces always the same outputs.
type RandomStrategy struct {
	root token.Token
}

// NewRandomStrategy returns a new instance of the random fuzzing strategy
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

// Fuzz starts the first iteration of the fuzzing strategy returning a channel which controls the iteration flow.
// The channel returns a value if the iteration is complete and waits with calculating the next iteration until a value is put in. The channel is automatically closed when there are no more iterations. The error return argument is not nil if an error occurs during the setup of the fuzzing strategy.
func (s *RandomStrategy) Fuzz(r rand.Rand) (chan struct{}, error) {
	if token.LoopExists(s.root) {
		return nil, &Error{
			Message: "found endless loop in graph. Cannot proceed.",
			Type:    ErrorEndlessLoopDetected,
		}
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

func (s *RandomStrategy) fuzz(tok token.Token, r rand.Rand, variableScope *token.VariableScope) {
	log.Debugf("Fuzz (%p)%#v with maxPermutations %d", tok, tok, tok.Permutations())

	if t, ok := tok.(token.Scoping); ok && t.Scoping() {
		variableScope = variableScope.Push()
	}

	err := tok.Permutation(uint(r.Int63n(int64(tok.Permutations())) + 1))
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

func (s *RandomStrategy) fuzzYADDA(root token.Token, r rand.Rand) {
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

			err := tok.Permutation(uint(r.Int63n(int64(tok.Permutations())) + 1))
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
