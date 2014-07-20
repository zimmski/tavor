package strategy

import (
	"fmt"
	"sort"

	"github.com/zimmski/tavor/token"
)

type StrategyErrorType int

const (
	StrategyErrorEndlessLoopDetected StrategyErrorType = iota
)

type StrategyError struct {
	Message string
	Type    StrategyErrorType
}

func (err *StrategyError) Error() string {
	return err.Message
}

type ReduceFeedbackType int

const (
	Good ReduceFeedbackType = iota
	Bad
)

type Strategy interface {
	Reduce() (chan struct{}, chan<- ReduceFeedbackType, error)
}

var strategies = make(map[string]func(tok token.Token) Strategy)

func New(name string, tok token.Token) (Strategy, error) {
	strat, ok := strategies[name]
	if !ok {
		return nil, fmt.Errorf("unknown reduce strategy \"%s\"", name)
	}

	return strat(tok), nil
}

func List() []string {
	keyStrategies := make([]string, 0, len(strategies))

	for key := range strategies {
		keyStrategies = append(keyStrategies, key)
	}

	sort.Strings(keyStrategies)

	return keyStrategies
}

func Register(name string, strat func(tok token.Token) Strategy) {
	if strat == nil {
		panic("register reduce strategy is nil")
	}

	if _, ok := strategies[name]; ok {
		panic("reduce strategy " + name + " already registered")
	}

	strategies[name] = strat
}
