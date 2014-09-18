package strategy

import (
	"fmt"
	"sort"

	"github.com/zimmski/tavor/token"
)

type ErrorType int

const (
	ErrorEndlessLoopDetected ErrorType = iota
)

type Error struct {
	Message string
	Type    ErrorType
}

func (err *Error) Error() string {
	return err.Message
}

type ReduceFeedbackType int

const (
	Unknown ReduceFeedbackType = iota
	Good
	Bad
)

func (f ReduceFeedbackType) String() string {
	switch f {
	case Bad:
		return "bad"
	case Good:
		return "good"
	default:
		return "unknown feedback"
	}
}

type Strategy interface {
	Reduce() (chan struct{}, chan<- ReduceFeedbackType, error)
}

var strategyLookup = make(map[string]func(tok token.Token) Strategy)

func New(name string, tok token.Token) (Strategy, error) {
	strat, ok := strategyLookup[name]
	if !ok {
		return nil, fmt.Errorf("unknown reduce strategy %q", name)
	}

	return strat(tok), nil
}

func List() []string {
	keyStrategyLookup := make([]string, 0, len(strategyLookup))

	for key := range strategyLookup {
		keyStrategyLookup = append(keyStrategyLookup, key)
	}

	sort.Strings(keyStrategyLookup)

	return keyStrategyLookup
}

func Register(name string, strat func(tok token.Token) Strategy) {
	if strat == nil {
		panic("register reduce strategy is nil")
	}

	if _, ok := strategyLookup[name]; ok {
		panic("reduce strategy " + name + " already registered")
	}

	strategyLookup[name] = strat
}
