package strategy

import (
	"fmt"
	"sort"

	"github.com/zimmski/tavor/token"
)

// ErrorType the reduce strategy error type
type ErrorType int

const (
	// ErrEndlessLoopDetected the token graph has a cycle which is not allowed.
	ErrEndlessLoopDetected ErrorType = iota
)

// Error holds a reduce strategy error
type Error struct {
	Message string
	Type    ErrorType
}

func (err *Error) Error() string {
	return err.Message
}

// ReduceFeedbackType the reduce strategy feedback type
type ReduceFeedbackType int

//go:generate stringer -type=ReduceFeedbackType
const (
	// Unknown the feedback is of unknown type, this is always a fatal error
	Unknown ReduceFeedbackType = iota
	// Good the reduce step produced a successful result
	Good
	// Bad the reduce step produced an unsuccessful result
	Bad
)

// Strategy defines a reduce strategy
// The function starts the first step of the reduce strategy returning a channel which controls the step flow and a channel for the feedback of the step. The channel returns a value if the step is complete and waits with calculating the next step until a value is put in and feedback is given. The channels are automatically closed when there are no more steps. The error return argument is not nil if an error occurs during the initialization of the reduce strategy.
type Strategy func(root token.Token) (chan struct{}, chan<- ReduceFeedbackType, error)

var strategyLookup = make(map[string]Strategy)

// New returns a new reduce strategy instance given the registered name of the strategy.
// The error return argument is not nil, if the name does not exist in the registered reduce strategy list.
func New(name string) (Strategy, error) {
	strat, ok := strategyLookup[name]
	if !ok {
		return nil, fmt.Errorf("unknown reduce strategy %q", name)
	}

	return strat, nil
}

// List returns a list of all registered reduce strategy names.
func List() []string {
	keyStrategyLookup := make([]string, 0, len(strategyLookup))

	for key := range strategyLookup {
		keyStrategyLookup = append(keyStrategyLookup, key)
	}

	sort.Strings(keyStrategyLookup)

	return keyStrategyLookup
}

// Register registers a reduce strategy instance function with the given name.
func Register(name string, strat Strategy) {
	if strat == nil {
		panic("register reduce strategy is nil")
	}

	if _, ok := strategyLookup[name]; ok {
		panic("reduce strategy " + name + " already registered")
	}

	strategyLookup[name] = strat
}
