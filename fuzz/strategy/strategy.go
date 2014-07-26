package strategy

import (
	"fmt"
	"sort"

	"github.com/zimmski/tavor/rand"
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

type Strategy interface {
	Fuzz(r rand.Rand) (chan struct{}, error)
}

var strategyLookup = make(map[string]func(tok token.Token) Strategy)

func New(name string, tok token.Token) (Strategy, error) {
	strat, ok := strategyLookup[name]
	if !ok {
		return nil, fmt.Errorf("unknown fuzzing strategy %q", name)
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
		panic("register fuzzing strategy is nil")
	}

	if _, ok := strategyLookup[name]; ok {
		panic("fuzzing strategy " + name + " already registered")
	}

	strategyLookup[name] = strat
}
