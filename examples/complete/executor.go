// +build example-main

package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/zimmski/tavor/executor/keydriven"

	"github.com/zimmski/tavor/examples/complete/implementation"
)

const (
	exitPassed = iota
	exitFailed
	exitError
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <inputfile>\n", os.Args[0])

		os.Exit(exitError)
	}

	cmds, err := keydriven.ReadKeyDrivenFile(os.Args[1])
	if err != nil {
		fmt.Printf("Error: %v\n", err)

		os.Exit(exitError)
	}

	executor := initExecutor()

	if err := executor.Execute(cmds); err != nil {
		fmt.Printf("Failed: %v\n", err)

		os.Exit(exitFailed)
	}

	os.Exit(exitPassed)
}

func initExecutor() *keydriven.Executor {
	executor := keydriven.NewExecutor()

	executor.BeforeAction = func(key string, parameters ...string) error {
		fmt.Printf("%s %v\n", key, parameters)

		return nil
	}

	machine := implementation.NewVendingMachine()

	executor.MustRegister("credit", func(key string, parameters ...string) error {
		if err := checkParameterCount(key, len(parameters), 1); err != nil {
			return err
		}

		expected, err := strconv.Atoi(parameters[0])
		if err != nil {
			return err
		}

		got := machine.Credit()

		if expected != got {
			return fmt.Errorf("Credit should be %d but was %d", expected, got)
		}

		return nil
	})

	executor.MustRegister("coin", func(key string, parameters ...string) error {
		if err := checkParameterCount(key, len(parameters), 1); err != nil {
			return err
		}

		coin, err := strconv.Atoi(parameters[0])
		if err != nil {
			return err
		}

		err = machine.Coin(coin)
		if err != nil {
			return err
		}

		return nil
	})

	executor.MustRegister("vend", func(key string, parameters ...string) error {
		if err := checkParameterCount(key, len(parameters), 0); err != nil {
			return err
		}

		vend := machine.Vend()
		if !vend {
			return fmt.Errorf("Could not vend")
		}

		return nil
	})

	return executor
}

func checkParameterCount(key string, got int, expected int) error {
	if got != expected {
		return &keydriven.Error{
			Message: fmt.Sprintf("Key %q requires %d parameters not %d", key, expected, got),
			Err:     keydriven.ErrInvalidParametersCount,
		}
	}

	return nil
}
