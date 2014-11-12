package strategy

import (
	"github.com/zimmski/container/list/linkedlist"

	"github.com/zimmski/tavor/log"
	"github.com/zimmski/tavor/token"
)

type binarySearchLevel struct {
	token         token.ReduceToken
	reduction     uint
	maxReductions uint

	children []binarySearchLevel
}

// BinarySearchStrategy implements a reduce strategy that reduces the data through a binary search alike algorithm.
// Every step of the strategy generates a new valid token graph state. The generation is deterministic. The algorithm starts by deactivating all optional tokens, this includes for example reducing lists to their minimum repetition. Each step uses the feedback to determine which tokens to reactivate next. All steps use a binary search alike algorithm to focus the reactivation only on a portion of the available tokens.
type BinarySearchStrategy struct {
	root token.Token
}

// NewBinarySearch returns a new instance of the binary search reduce strategy
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
		case token.List:
			l := t.Len()

			for i := 0; i < l; i++ {
				c, _ := t.Get(i)

				queue.Push(c)
			}
		}
	}

	return tree
}

func (s *BinarySearchStrategy) setReduction(tok token.ReduceToken, reduction uint) {
	log.Debugf("set (%p)%#v to reduction %d", tok, tok, reduction)

	if err := tok.Reduce(reduction); err != nil {
		panic(err)
	}
}

// Reduce starts the first step of the reduce strategy returning a channel which controls the step flow and a channel for the feedback of the step.
// The channel returns a value if the step is complete and waits with calculating the next step until a value is put in and feedback is given. The channels are automatically closed when there are no more steps. The error return argument is not nil if an error occurs during the initialization of the reduce strategy.
func (s *BinarySearchStrategy) Reduce() (chan struct{}, chan<- ReduceFeedbackType, error) {
	if token.LoopExists(s.root) {
		return nil, nil, &Error{
			Message: "found endless loop in graph. Cannot proceed.",
			Type:    ErrorEndlessLoopDetected,
		}
	}

	continueReducing := make(chan struct{})
	feedbackReducing := make(chan ReduceFeedbackType)

	go func() {
		log.Debug("start binary search routine")

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

func (s *BinarySearchStrategy) reduce(continueReducing chan struct{}, feedbackReducing <-chan ReduceFeedbackType, tree []binarySearchLevel) bool {
	log.Debugf("reducing level %d->%#v", len(tree), tree)

	// we always asume that the initial values are bad so we ignore them right away

	// TODO do a binary search on the level entries
	for _, c := range tree {
		// reduce beginning from the first reduction
		c.reduction = 0

		for {
			// TODO do a binary search on the 1..maxReductions for this level entry
			c.reduction++
			if err := c.token.Reduce(c.reduction); err != nil {
				panic(err)
			}

			contin, feedback := s.nextStep(continueReducing, feedbackReducing)
			if !contin {
				return false
			} else if feedback == Bad {
				break
			}

			if c.reduction == c.maxReductions-1 {
				log.Debug("use initial value, nothing to reduce")

				c.reduction = c.maxReductions
				if err := c.token.Reduce(c.reduction); err != nil {
					panic(err)
				}

				break
			}
		}

		log.Debugf("reduced (%p)%#v to reduction %d/%d", c.token, c.token, c.reduction, c.maxReductions)

		c.children = s.getTree(c.token, true)

		if len(c.children) > 0 {
			log.Debugf("reduce the children of (%p)%#v", c.token, c.token, c.reduction, c.maxReductions)

			s.reduce(continueReducing, feedbackReducing, c.children)
		}
	}

	return true
}

func (s *BinarySearchStrategy) nextStep(continueReducing chan struct{}, feedbackReducing <-chan ReduceFeedbackType) (bool, ReduceFeedbackType) {
	token.ResetScope(s.root)
	token.ResetResetTokens(s.root)

	log.Debug("done with reducing step")

	// done with this reduce step
	continueReducing <- struct{}{}

	// wait until we got feedback to the current state
	feedback, ok := <-feedbackReducing
	if ok {
		// TODO implement the usage of the feedback
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
