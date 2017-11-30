package strategy

import (
	"github.com/zimmski/container/list/linkedlist"

	"github.com/zimmski/tavor/log"
	"github.com/zimmski/tavor/token"
)

func init() {
	Register("Linear", NewLinear)
}

type linearStrategyLevel struct {
	token         token.ReduceToken
	reduction     uint
	maxReductions uint

	children []linearStrategyLevel
}

type linearStrategy struct {
	root token.Token
}

// NewLinear implements a reduce strategy that reduces the data through a linear search algorithm.
// Every step of the strategy generates a new valid token graph state. The generation is deterministic. The algorithm starts by deactivating all optional tokens, this includes for example reducing lists to their minimum repetition. Each step uses the feedback to determine which tokens to reactivate next.
func NewLinear(root token.Token) (chan struct{}, chan<- ReduceFeedbackType, error) {
	if token.LoopExists(root) {
		return nil, nil, &Error{
			Message: "found endless loop in graph. Cannot proceed.",
			Type:    ErrEndlessLoopDetected,
		}
	}

	s := &linearStrategy{
		root: root,
	}

	continueReducing := make(chan struct{})
	feedbackReducing := make(chan ReduceFeedbackType)

	go func() {
		log.Debug("start linear routine")

		tree := s.getTree(s.root, false)

		if len(tree) > 0 {
			log.Debug("start reducing step")

			if contin := s.reduce(continueReducing, feedbackReducing, tree); !contin {
				return
			}
		} else {
			log.Debug("no reduceable tokens to begin with")
		}

		log.Debug("finished reducing")

		close(continueReducing)
		close(feedbackReducing)
	}()

	return continueReducing, feedbackReducing, nil
}

func (s *linearStrategy) reduce(continueReducing chan struct{}, feedbackReducing <-chan ReduceFeedbackType, tree []linearStrategyLevel) bool {
	log.Debugf("reducing level %d->%#v", len(tree), tree)

	// we always asume that the initial values are good so we ignore them right away

	for _, c := range tree {
		// reduce beginning from the first reduction
		c.reduction = 0

		for {
			if err := c.token.Reduce(c.reduction); err != nil {
				panic(err)
			}

			contin, feedback := s.nextStep(continueReducing, feedbackReducing)
			if !contin {
				return false
			} else if feedback == Good {
				break
			}

			if c.reduction == c.maxReductions-1 {
				log.Debugf("use initial value for (%p)%#v, nothing to reduce", c.token, c.token)

				c.reduction = c.maxReductions
				if err := c.token.Reduce(c.reduction); err != nil {
					panic(err)
				}

				break
			}

			c.reduction++
		}

		log.Debugf("reduced (%p)%#v to reduction %d/%d", c.token, c.token, c.reduction, c.maxReductions)

		c.children = s.getTree(c.token, true)

		if len(c.children) > 0 {
			log.Debugf("reduce the children of (%p)%#v %d/%d", c.token, c.token, c.reduction, c.maxReductions)

			s.reduce(continueReducing, feedbackReducing, c.children)
		}
	}

	return true
}

func (s *linearStrategy) nextStep(continueReducing chan struct{}, feedbackReducing <-chan ReduceFeedbackType) (bool, ReduceFeedbackType) {
	token.ResetScope(s.root)
	_ = token.ResetResetTokens(s.root)
	token.ResetScope(s.root)

	log.Debug("done with reducing step")

	// done with this reduce step
	continueReducing <- struct{}{}

	// wait until we got feedback to the current state
	feedback, ok := <-feedbackReducing
	if ok {
		log.Debugf("GOT FEEDBACK -> Looks %s", feedback)
	} else {
		log.Debug("reducing feedback channel closed from outside")

		return false, Unknown
	}

	// wait until we are allowed to continue
	if _, ok := <-continueReducing; !ok {
		log.Debug("reducing continue channel closed from outside")

		return false, Unknown
	}

	log.Debug("start reducing step")

	return true, feedback
}

func (s *linearStrategy) getTree(root token.Token, fromChildren bool) []linearStrategyLevel {
	var tree []linearStrategyLevel
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
		case token.ReduceToken:
			if t.Reduces() < 2 {
				continue
			}

			maxReductions := t.Reduces() - 1

			s.setReduction(t, maxReductions)

			tree = append(tree, linearStrategyLevel{
				token:         t,
				reduction:     maxReductions,
				maxReductions: maxReductions,
			})
		case token.ForwardToken:
			c := t.Get()

			queue.Unshift(c)
		case token.ListToken:
			for i := t.Len() - 1; i >= 0; i-- {
				c, _ := t.Get(i)

				queue.Unshift(c)
			}
		}
	}

	return tree
}

func (s *linearStrategy) setReduction(tok token.ReduceToken, reduction uint) {
	log.Debugf("set (%p)%#v to reduction %d", tok, tok, reduction)

	if err := tok.Reduce(reduction); err != nil {
		panic(err)
	}
}
