package strategy

import (
	"github.com/zimmski/container/list/linkedlist"

	"github.com/zimmski/tavor"
	"github.com/zimmski/tavor/log"
	"github.com/zimmski/tavor/rand"
	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/lists"
)

type allPermutationsLevel struct {
	token           token.Token
	permutation     int
	maxPermutations int

	children []allPermutationsLevel
}

type AllPermutationsStrategy struct {
	root token.Token
}

func NewAllPermutationsStrategy(tok token.Token) *AllPermutationsStrategy {
	s := &AllPermutationsStrategy{
		root: tok,
	}

	return s
}

func init() {
	Register("AllPermutations", func(tok token.Token) Strategy {
		return NewAllPermutationsStrategy(tok)
	})
}

func (s *AllPermutationsStrategy) getTree(root token.Token, fromChildren bool) []allPermutationsLevel {
	var tree []allPermutationsLevel

	add := func(tok token.Token) {
		s.setPermutation(tok, 1)

		tree = append(tree, allPermutationsLevel{
			token:           tok,
			permutation:     1,
			maxPermutations: tok.Permutations(),

			children: s.getTree(tok, true),
		})
	}

	if fromChildren {
		switch t := root.(type) {
		case token.ForwardToken:
			if v := t.Get(); v != nil {
				add(v)
			}
		case lists.List:
			for i := 0; i < t.Len(); i++ {
				c, _ := t.Get(i)

				add(c)
			}
		}
	} else {
		add(root)
	}

	return tree
}

func (s *AllPermutationsStrategy) setPermutation(tok token.Token, permutation int) {
	log.Debugf("set %#v(%p) to permutation %d", tok, tok, permutation)

	if err := tok.Permutation(permutation); err != nil {
		panic(err)
	}
}

func (s *AllPermutationsStrategy) Fuzz(r rand.Rand) (chan struct{}, error) {
	if tavor.LoopExists(s.root) {
		return nil, &StrategyError{
			Message: "found endless loop in graph. Cannot proceed.",
			Type:    StrategyErrorEndlessLoopDetected,
		}
	}

	continueFuzzing := make(chan struct{})

	go func() {
		log.Debug("start all permutations routine")

		tree := s.getTree(s.root, false)

		log.Debug("start fuzzing step")

		if contin, _ := s.fuzz(continueFuzzing, tree, false); !contin {
			return
		}

		log.Debug("finished fuzzing.")

		close(continueFuzzing)
	}()

	return continueFuzzing, nil
}

func (s *AllPermutationsStrategy) fuzz(continueFuzzing chan struct{}, tree []allPermutationsLevel, justastep bool) (bool, bool) {
	log.Debugf("fuzzing level %d->%#v", len(tree), tree)

STEP:
	for {
		if justastep && len(tree[0].children) != 0 {
			log.Debugf("STEP FURTHER INTO")

			if contin, step := s.fuzz(continueFuzzing, tree[0].children, justastep); !contin {
				return false, false
			} else if step {
				log.Debugf("CONTINUE after child step")

				return true, true
			}

			log.Debugf("PERMUTATE after child step")
		} else {
			log.Debugf("permute %d->%#v", 0, tree[0])

			if tree[0].permutation != 1 {
				s.setPermutation(tree[0].token, tree[0].permutation)
				tree[0].children = s.getTree(tree[0].token, true)

				if justastep {
					log.Debugf("CONTINUE after permutate")

					return true, true
				}
			}

			if len(tree[0].children) != 0 {
				if contin, step := s.fuzz(continueFuzzing, tree[0].children, justastep); !contin {
					return false, false
				} else if step {
					log.Debugf("CONTINUE after child step")

					return true, true
				}
			} else {
				if !justastep && (tree[0].token != s.root || tree[0].permutation <= tree[0].maxPermutations) && !s.nextStep(continueFuzzing) {
					return false, false
				}
			}
		}

		tree[0].permutation++

		if tree[0].permutation > tree[0].maxPermutations {
			for i := 0; i < len(tree); i++ {
				log.Debugf("check %d vs %d for %#v", tree[i].permutation, tree[i].maxPermutations, tree[i])
			}

			i := 0

			for {
				if i == len(tree)-1 {
					log.Debugf("done with fuzzing this level because %#v", tree)

					break STEP
				}

				i++

				if len(tree[i].children) != 0 {
					log.Debugf("CHECK children %#v", tree[i])

					if contin, step := s.fuzz(continueFuzzing, tree[i].children, true); !contin {
						return false, false
					} else if step {
						for j := 0; j < i; j++ {
							tree[j].permutation = 1
							s.setPermutation(tree[j].token, tree[j].permutation)
							tree[j].children = s.getTree(tree[j].token, true)
						}

						if justastep {
							return true, true
						}

						log.Debugf("STEP continue")

						continue STEP
					}

					log.Debugf("PERMUTATE continue")
				}

				tree[i].permutation++

				if tree[i].permutation <= tree[i].maxPermutations {
					for j := 0; j < i; j++ {
						tree[j].permutation = 1
						s.setPermutation(tree[j].token, tree[j].permutation)
						tree[j].children = s.getTree(tree[j].token, true)
					}

					log.Debugf("permute %d->%#v", i, tree[i])

					s.setPermutation(tree[i].token, tree[i].permutation)
					tree[i].children = s.getTree(tree[i].token, true)

					if justastep {
						return true, true
					}

					continue STEP
				}
			}
		} else if justastep {
			s.setPermutation(tree[0].token, tree[0].permutation)
			tree[0].children = s.getTree(tree[0].token, true)

			log.Debugf("CONTINUE after permutate")

			return true, true
		}
	}

	return true, false
}

func (s *AllPermutationsStrategy) nextStep(continueFuzzing chan struct{}) bool {
	s.resetResetTokens()

	log.Debug("done with fuzzing step")

	// done with this fuzzing step
	continueFuzzing <- struct{}{}

	// wait until we are allowed to continue
	if _, ok := <-continueFuzzing; !ok {
		log.Debug("fuzzing channel closed from outside")

		return false
	}

	log.Debug("start fuzzing step")

	return true
}

func (s *AllPermutationsStrategy) resetResetTokens() {
	var queue = linkedlist.New()

	queue.Push(s.root)

	for !queue.Empty() {
		v, _ := queue.Shift()

		switch tok := v.(type) {
		case token.ResetToken:
			log.Debugf("reset %#v(%p)", tok, tok)

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
