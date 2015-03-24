package parser

import (
	"fmt"
	"strconv"
)

type argumentsParser struct {
	arguments     map[string]string
	usedArguments map[string]struct{}
	err           error
}

func newArgumentsParser(arguments map[string]string) *argumentsParser {
	return &argumentsParser{
		arguments:     arguments,
		usedArguments: make(map[string]struct{}),
		err:           nil,
	}
}

// GetInt tries to parse the argument name and returns its integer value or defaultValue if the argument is not found.
// The return value is valid only if Err returns nil.
func (ap *argumentsParser) GetInt(name string, defaultValue int) int {
	if ap.err != nil {
		return -1
	}

	raw, found := ap.arguments[name]
	if !found {
		return defaultValue
	}

	val, err := strconv.Atoi(raw)
	if err != nil {
		ap.err = fmt.Errorf("%q needs an integer value", name)
		return -1
	}

	ap.usedArguments[name] = struct{}{}
	return val
}

// Err returns the first error encountered by the ArgumentsParser
func (ap *argumentsParser) Err() error {
	return ap.err
}

func (ap *argumentsParser) firstUnusedArgument() string {
	for arg := range ap.arguments {
		if _, ok := ap.usedArguments[arg]; !ok {
			return arg
		}
	}

	return ""
}
