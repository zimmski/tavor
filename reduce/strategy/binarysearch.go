package strategy

import (
	"github.com/zimmski/container/list/linkedlist"

	"github.com/zimmski/tavor"
	"github.com/zimmski/tavor/log"
	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/lists"
)

type binarySearchLevel struct {
	token         token.ReduceToken
	reduction     int
	maxReductions int

	children []binarySearchLevel
}

type BinarySearchStrategy struct {
	root token.Token
}

func NewBinarySearch(tok token.Token) *BinarySearchStrategy {
	s := &BinarySearchStrategy{
		root: tok,
	}

	return s
}

func init() {
	Register("BinarySearch", func(tok token.Token) Strategy {
		return NewBinarySearch(tok)
	})
}

func (s *BinarySearchStrategy) getTree(root token.Token, fromChildren bool) []binarySearchLevel {
	var tree []binarySearchLevel
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
		case token.ReduceToken:
			if t.Reduces() < 2 {
				continue
			}

			s.setReduction(t, 1)

			tree = append(tree, binarySearchLevel{
				token:         t,
				reduction:     1,
				maxReductions: t.Reduces(),

				children: s.getTree(t, true),
			})
		case token.ForwardToken:
			c := t.Get()

			queue.Push(c)
		case lists.List:
			l := t.Len()

			for i := 0; i < l; i++ {
				c, _ := t.Get(i)

				queue.Push(c)
			}
		}
	}

	return tree
}

func (s *BinarySearchStrategy) resetResetTokens() {
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

func (s *BinarySearchStrategy) setReduction(tok token.ReduceToken, reduction int) {
	log.Debugf("Set %#v(%p) to reduction %d", tok, tok, reduction)

	if err := tok.Reduce(reduction); err != nil {
		panic(err)
	}
}

func (s *BinarySearchStrategy) Reduce() (chan struct{}, chan<- ReduceFeedbackType, error) {
	if tavor.LoopExists(s.root) {
		return nil, nil, &StrategyError{
			Message: "Found endless loop in graph. Cannot proceed.",
			Type:    StrategyErrorEndlessLoopDetected,
		}
	}

	continueReducing := make(chan struct{})
	feedbackReducing := make(chan ReduceFeedbackType)

	go func() {
		log.Debug("Start binary search routine")

		tree := s.getTree(s.root, false)

		if len(tree) != 0 {
			log.Debug("Start reducing step")

			if contin, _ := s.reduce(continueReducing, feedbackReducing, tree, false); !contin {
				return
			}
		} else {
			log.Debug("No reduceable tokens to begin with")
		}

		s.resetResetTokens()

		log.Debug("Finished reducing")

		close(continueReducing)
		close(feedbackReducing)
	}()

	return continueReducing, feedbackReducing, nil
}

func (s *BinarySearchStrategy) reduce(continueReducing chan struct{}, feedbackReducing <-chan ReduceFeedbackType, tree []binarySearchLevel, justastep bool) (bool, bool) {
	log.Debugf("Reducing level %d->%#v", len(tree), tree)

STEP:
	for {
		if justastep && len(tree[0].children) != 0 {
			log.Debugf("STEP FURTHER INTO")

			if contin, step := s.reduce(continueReducing, feedbackReducing, tree[0].children, justastep); !contin {
				return false, false
			} else if step {
				log.Debugf("CONTINUE after child step")

				return true, true
			}

			log.Debugf("REDUCE after child step")
		} else {
			log.Debugf("Reduce %d->%#v", 0, tree[0])

			if tree[0].reduction != 1 {
				s.setReduction(tree[0].token, tree[0].reduction)
				tree[0].children = s.getTree(tree[0].token, true)

				if justastep {
					log.Debugf("CONTINUE after reduction")

					return true, true
				}
			}

			if len(tree[0].children) != 0 {
				if contin, step := s.reduce(continueReducing, feedbackReducing, tree[0].children, justastep); !contin {
					return false, false
				} else if step {
					log.Debugf("CONTINUE after child step")

					return true, true
				}
			} else {
				if !justastep && (tree[0].token != s.root || tree[0].reduction <= tree[0].maxReductions) && !s.nextStep(continueReducing, feedbackReducing) {
					return false, false
				}
			}
		}

		tree[0].reduction++

		if tree[0].reduction > tree[0].maxReductions {
			for i := 0; i < len(tree); i++ {
				log.Debugf("Check %d vs %d for %#v", tree[i].reduction, tree[i].maxReductions, tree[i])
			}

			i := 0

			for {
				if i == len(tree)-1 {
					log.Debugf("Done with reducing this level because %#v", tree)

					break STEP
				}

				i++

				if len(tree[i].children) != 0 {
					log.Debugf("CHECK children %#v", tree[i])

					if contin, step := s.reduce(continueReducing, feedbackReducing, tree[i].children, true); !contin {
						return false, false
					} else if step {
						for j := 0; j < i; j++ {
							tree[j].reduction = 1
							s.setReduction(tree[j].token, tree[j].reduction)
							tree[j].children = s.getTree(tree[j].token, true)
						}

						if justastep {
							return true, true
						}

						log.Debugf("STEP continue")

						continue STEP
					}

					log.Debugf("REDUCE continue")
				}

				tree[i].reduction++

				if tree[i].reduction <= tree[i].maxReductions {
					for j := 0; j < i; j++ {
						tree[j].reduction = 1
						s.setReduction(tree[j].token, tree[j].reduction)
						tree[j].children = s.getTree(tree[j].token, true)
					}

					log.Debugf("Reduce %d->%#v", i, tree[i])

					s.setReduction(tree[i].token, tree[i].reduction)
					tree[i].children = s.getTree(tree[i].token, true)

					if justastep {
						return true, true
					}

					continue STEP
				}
			}
		} else if justastep {
			s.setReduction(tree[0].token, tree[0].reduction)
			tree[0].children = s.getTree(tree[0].token, true)

			log.Debugf("CONTINUE after reduction")

			return true, true
		}
	}

	return true, false
}

func (s *BinarySearchStrategy) nextStep(continueReducing chan struct{}, feedbackReducing <-chan ReduceFeedbackType) bool {
	s.resetResetTokens()

	log.Debug("Done with reducing step")

	// done with this reduce step
	continueReducing <- struct{}{}

	// wait until we got feedback to the current state
	if feedback, ok := <-feedbackReducing; ok {
		// TODO implement the usage of the feedback
		log.Debugf("GOT FEEDBACK -> Looks %s", feedback)
	} else {
		log.Debug("Reducing feedback channel closed from outside")

		return false
	}

	// wait until we are allowed to continue
	if _, ok := <-continueReducing; !ok {
		log.Debug("Reducing continue channel closed from outside")

		return false
	}

	log.Debug("Start reducing step")

	return true
}
