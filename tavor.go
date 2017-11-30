// Package tavor provides all general properties, constants and functions for the Tavor framework and tools
package tavor

import (
	"fmt"
)

const (
	// Version of the framework and tools
	Version = "0.6"
)

// MaxRepeat determines the maximum copies in graph cycles.
var MaxRepeat = 2

// ErrNoSequenceValue there is no item left to choose an existing item.
var ErrNoSequenceValue = fmt.Errorf("There is no sequence value to choose from")
