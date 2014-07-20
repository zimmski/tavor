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

			maxReductions := t.Reduces()

			s.setReduction(t, maxReductions)

			tree = append(tree, binarySearchLevel{
				token:         t,
				reduction:     maxReductions,
				maxReductions: maxReductions,
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
	log.Debugf("Set (%p)%#v to reduction %d", tok, tok, reduction)

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

			if contin := s.reduce(continueReducing, feedbackReducing, tree); !contin {
				return
			}
		} else {
			log.Debug("No reduceable tokens to begin with")
		}

		log.Debug("Finished reducing")

		close(continueReducing)
		close(feedbackReducing)
	}()

	return continueReducing, feedbackReducing, nil
}

func (s *BinarySearchStrategy) reduce(continueReducing chan struct{}, feedbackReducing <-chan ReduceFeedbackType, tree []binarySearchLevel) bool {
	log.Debugf("Reducing level %d->%#v", len(tree), tree)

	// we always asume that the initial values are bad so we ignore them right away

	// TODO do a binary search on the level entries
	for _, c := range tree {
		for {
			// TODO do a binary search on the 1..maxReductions for this level entry
			c.reduction--
			c.token.Reduce(c.reduction)

			contin, feedback := s.nextStep(continueReducing, feedbackReducing)
			if !contin {
				return false
			} else if feedback == Good {
				log.Debugf("Go back a reduction")

				c.reduction++
				c.token.Reduce(c.reduction)

				break
			}

			if c.reduction == 1 {
				break
			}
		}

		log.Debugf("Reduced (%p)%#v to reduction %d/%d", c.token, c.token, c.reduction, c.maxReductions)

		c.children = s.getTree(c.token, true)

		if len(c.children) != 0 {
			log.Debugf("Reduce the children of (%p)%#v", c.token, c.token, c.reduction, c.maxReductions)

			s.reduce(continueReducing, feedbackReducing, c.children)
		}
	}

	return true
}

func (s *BinarySearchStrategy) nextStep(continueReducing chan struct{}, feedbackReducing <-chan ReduceFeedbackType) (bool, ReduceFeedbackType) {
	s.resetResetTokens()

	log.Debug("Done with reducing step")

	// done with this reduce step
	continueReducing <- struct{}{}

	// wait until we got feedback to the current state
	feedback, ok := <-feedbackReducing
	if ok {
		// TODO implement the usage of the feedback
		log.Debugf("GOT FEEDBACK -> Looks %s", feedback)
	} else {
		log.Debug("Reducing feedback channel closed from outside")

		return false, Unknown
	}

	// wait until we are allowed to continue
	if _, ok := <-continueReducing; !ok {
		log.Debug("Reducing continue channel closed from outside")

		return false, Unknown
	}

	log.Debug("Start reducing step")

	return true, feedback
}
