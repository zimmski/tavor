package strategy

import (
	"fmt"
	"sort"

	"github.com/zimmski/tavor/rand"
	"github.com/zimmski/tavor/token"
)

type Strategy interface {
	Fuzz(r rand.Rand)
}

var strategies = make(map[string]func(tok token.Token) Strategy)

func New(name string, tok token.Token) (Strategy, error) {
	if strat, ok := strategies[name]; ok {
		return strat(tok), nil
	} else {
		return nil, fmt.Errorf("Unknown fuzzing strategy \"%s\"", name)
	}
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
		panic("Register fuzzing strategy is nil")
	}

	if _, ok := strategies[name]; ok {
		panic("Fuzzing strategy " + name + " already registered")
	}

	strategies[name] = strat
}
