// +build example-main

package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/zimmski/tavor/examples/complete/implementation"
)

type action func(parameters []string) error

type command struct {
	key        string
	parameters []string
}

var actions = make(map[string]action)

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

	input, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Printf("Error: %v\n", err)

		os.Exit(exitError)
	}

	var cmds []*command

	for li, l := range strings.Split(string(input), "\n") {
		lc := strings.Split(l, "\t")

		for i := 0; i < len(lc); i++ {
			lc[i] = strings.Trim(lc[i], "\r ")
		}

		if len(lc[0]) != 0 {
			if _, ok := actions[lc[0]]; !ok {
				fmt.Printf("Error: Unknown key %q at line %d\n", lc[0], li+1)

				os.Exit(exitError)
			}

			cmds = append(cmds, &command{
				key:        lc[0],
				parameters: lc[1:],
			})
		}
	}

	for _, cmd := range cmds {
		fmt.Printf("%s %v\n", cmd.key, cmd.parameters)

		err := actions[cmd.key](cmd.parameters)
		if err != nil {
			fmt.Printf("Error: %v\n", err)

			os.Exit(exitFailed)
		}
	}

	os.Exit(exitPassed)
}

var (
	errInvalidParametersCount = errors.New("Invalid parmaters count")
)

func init() {
	machine := implementation.NewVendingMachine()

	actions["credit"] = func(parameters []string) error {
		if len(parameters) != 1 {
			return errInvalidParametersCount
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
	}

	actions["coin"] = func(parameters []string) error {
		if len(parameters) != 1 {
			return errInvalidParametersCount
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
	}

	actions["vend"] = func(parameters []string) error {
		if len(parameters) != 0 {
			return errInvalidParametersCount
		}

		vend := machine.Vend()
		if !vend {
			return fmt.Errorf("Could not vend")
		}

		return nil
	}
}
