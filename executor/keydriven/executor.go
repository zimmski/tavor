package keydriven

import (
	"fmt"
	"sync"
)

// Action defines a key driven action
type Action func(key string, parameters ...string) error

// Command holds a key driven command instance
type Command struct {
	Key        string
	Parameters []string
}

// Executor holds a key driven executor instance
type Executor struct {
	sync.Mutex

	actions map[string]Action

	BeforeAction Action
	AfterAction  Action
}

// NewExecutor initialises and returns a new executor
func NewExecutor() *Executor {
	return &Executor{
		actions: make(map[string]Action),
	}
}

// Register adds a given key-action pair to the given executor
func (e *Executor) Register(key string, action Action) error {
	e.Lock()
	defer e.Unlock()

	if _, ok := e.actions[key]; ok {
		return &Error{
			Message: fmt.Sprintf("Key %q already defined", key),
			Err:     ErrKeyAlreadyDefined,
		}
	}

	e.actions[key] = action

	return nil
}

// MustRegister calls Register and panics if there is an error
func (e *Executor) MustRegister(key string, action Action) {
	if err := e.Register(key, action); err != nil {
		panic(err)
	}
}

// Execute executes a set of key driven commands.
func (e *Executor) Execute(cmds []Command) error {
	e.Lock()
	defer e.Unlock()

	for _, cmd := range cmds {
		a, ok := e.actions[cmd.Key]
		if !ok {
			return &Error{
				Message: fmt.Sprintf("Key %q is not defined", cmd.Key),
				Err:     ErrKeyNotDefined,
			}
		}

		if e.BeforeAction != nil {
			if err := e.BeforeAction(cmd.Key, cmd.Parameters...); err != nil {
				return err
			}
		}

		if err := a(cmd.Key, cmd.Parameters...); err != nil {
			return err
		}

		if e.AfterAction != nil {
			if err := e.AfterAction(cmd.Key, cmd.Parameters...); err != nil {
				return err
			}
		}
	}

	return nil
}
