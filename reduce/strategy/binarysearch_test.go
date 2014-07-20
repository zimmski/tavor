package strategy

import (
	"testing"

	. "github.com/stretchr/testify/assert"

	"github.com/zimmski/tavor/token"
)

func TestBinarySearchStrategyToBeStrategy(t *testing.T) {
	var strat *Strategy

	Implements(t, strat, &BinarySearchStrategy{})
}

func TestBinarySearchStrategy(t *testing.T) {

}

func TestBinarySearchStrategyLoopDetection(t *testing.T) {
	testStrategyLoopDetection(t, func(root token.Token) Strategy {
		return NewBinarySearch(root)
	})
}
