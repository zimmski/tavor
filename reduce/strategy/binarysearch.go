package strategy

import (
	"github.com/zimmski/tavor"
	"github.com/zimmski/tavor/log"
	"github.com/zimmski/tavor/token"
)

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

		log.Debug("Start reducing step")
		// TODO actually start the reducing process

		log.Debug("Finished reducing")

		close(continueReducing)
		close(feedbackReducing)
	}()

	return continueReducing, feedbackReducing, nil
}
