package strategy

import (
	"fmt"
	"sort"

	"github.com/zimmski/tavor/rand"
	"github.com/zimmski/tavor/token"
)

// ErrorType the fuzzing strategy error type
type ErrorType int

const (
	// ErrorEndlessLoopDetected the token graph has a cycle which is not allowed.
	ErrorEndlessLoopDetected ErrorType = iota
)

// Error holds a fuzzing strategy error
type Error struct {
	Message string
	Type    ErrorType
}

func (err *Error) Error() string {
	return err.Message
}

// Strategy defines a fuzzing strategy
type Strategy interface {
	// Fuzz starts the first iteration of the fuzzing strategy returning a channel which controls the iteration flow.
	// The channel returns a value if the iteration is complete and waits with calculating the next iteration until a value is put in. The channel is automatically closed when there are no more iterations. The error return argument is not nil if an error occurs during the setup of the fuzzing strategy.
	Fuzz(r rand.Rand) (chan struct{}, error)
}

var strategyLookup = make(map[string]func(tok token.Token) Strategy)

// New returns a new fuzzing strategy instance given the registered name of the strategy.
// The error return argument is not nil, if the name does not exist in the registered fuzzing strategy list.
func New(name string, tok token.Token) (Strategy, error) {
	strat, ok := strategyLookup[name]
	if !ok {
		return nil, fmt.Errorf("unknown fuzzing strategy %q", name)
	}

	return strat(tok), nil
}

// List returns a list of all registered fuzzing strategy names.
func List() []string {
	keyStrategyLookup := make([]string, 0, len(strategyLookup))

	for key := range strategyLookup {
		keyStrategyLookup = append(keyStrategyLookup, key)
	}

	sort.Strings(keyStrategyLookup)

	return keyStrategyLookup
}

// Register registers a fuzzing strategy instance function with the given name.
func Register(name string, strat func(tok token.Token) Strategy) {
	if strat == nil {
		panic("register fuzzing strategy is nil")
	}

	if _, ok := strategyLookup[name]; ok {
		panic("fuzzing strategy " + name + " already registered")
	}

	strategyLookup[name] = strat
}
