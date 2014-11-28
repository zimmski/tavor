# A complete example

This example provides a complete overview of Tavor. It does not utilize every feature but every component which should give you an idea on how Tavor can be used in your own projects. The example tests an implementation of the following state machine.

![Basic states and actions](/examples/quick/basic.png "Basic states and actions")

It uses the *Tavor format* to define this state machine and the *Tavor binary* to generate key-driven test files. An *executor* translates keys of a file into actions for the implementation under test. After successfully testing the original implementation, some bugs will be introduced to show the failure of some tests as well as automatically delta-debugging the failed key-driven test files.

**Please note:** The implementation of this example has intentional concurrency and other problems. Future versions of Tavor will help to identify, test and resolve such flaws. As Tavor evolves this example will evolve. This also concerns the given state machine and Tavor format. Both could be defined much more efficient with the help of state variables which are common in model-based testing tools but not yet implemented in the Tavor format. However, state variables can be easily implemented via code.

The following components will be defined:

- [Tavor format](/examples/complete/vending.tavor)
- [Executor](/examples/complete/executor.go)
- [Implementation](/examples/complete/implementation/vending.go)
- [Bash script to run all key-driven files](/examples/complete/run.sh)

## Tavor format

The following Tavor format, saved as [vending.tavor](/examples/complete/vending.tavor), defines a model of all possible valid states of the given state machine. Since this model should generate key-driven files, a set of keys and a key-driven format have to be defined.

The following keys will be used:

- **credit** performs a validation of the current credit amount. It compares the current credit against the integer argument of the key.
- **coin** invokes the action to insert a coin. It uses the integer argument as the coin credit.
- **vend** invokes the action to vend.

The key-driven format is defined as followed:

- Every line holds exactly one key
- A line begins with a key word
- After the key word zero, one or more arguments can be defined
- Each argument prepends a tab character
- A line ends with the new line character

Putting the given state machine, the defined keys and the rules for the key-driven format together, results in the following Tavor format:

```tavor
START = Credit0

Credit0   = "credit" "\t"   0 "\n" ( Coin25 Credit25 | Coin50 Credit50 | )
Credit25  = "credit" "\t"  25 "\n" ( Coin25 Credit50 | Coin50 Credit75 )
Credit50  = "credit" "\t"  50 "\n" ( Coin25 Credit75 | Coin50 Credit100 )
Credit75  = "credit" "\t"  75 "\n" Coin25 Credit100
Credit100 = "credit" "\t" 100 "\n" Vend Credit0

Coin25 = "coin" "\t" 25 "\n"
Coin50 = "coin" "\t" 50 "\n"

Vend = "vend" "\n"
```

This format file can now be easily fuzzed using the Tavor binary.

```bash
tavor --format-file vending.tavor fuzz
```

Every execution will output one generation of the format.

```
credit  0
coin    25
credit  25
coin    50
credit  75
coin    25
credit  100
vend
credit  0
```

Since there is a loop in the state machine, the graph can be traversed more than once.

```
credit  0
coin    25
credit  25
coin    50
credit  75
coin    25
credit  100
vend
credit  0
coin    25
credit  25
coin    50
credit  75
coin    25
credit  100
vend
credit  0
```

The default fuzzing strategy `random` can create all possible permutations of a format but since it is random, it will need enough time to do so. Lots of duplicated results will be generated, since even random events often lead to the same results. To work around this problem the `AllPermutations` strategy can be used which, as its name states, generates all possible permutations of a graph. This strategy should be used wisely since even small graphs can have [an enormous amount of permutations](https://en.wikipedia.org/wiki/Combinatorial_explosion). We can say that there are theoretically an infinite amount of permutations, since the example graph has a loop. To work around this additional problem, the `--max-repeat` argument will be used with a suitable value. It enforces a maximum of loop traversals and repetitions.

**Please note:** There is no easy way to ensure that the `--max-repeat` value is correct. A high value can lead to many repetitive permutations which will often not improve the testing process. A small value can lead to a bad coverage which means that some code branches would not be taken. Future versions of Tavor will aid this process by implementing better fuzzing strategies as well as additional metrics like token and edge coverage.

Since the given state machine is pretty easy to understand we can imply that a loop should be generated at least twice to include the repetition part of each loop. Putting this together we arrive at the following Tavor binary arguments.

```bash
tavor --format-file vending.tavor --max-repeat 2 fuzz --strategy AllPermutations
```

This command outputs directly to STDOUT which is OK for one generation but since we are using the `AllPermutations` fuzzing strategy we have to deal with many generations of the format. Additionally we want to save every permutation in a file so we can build a regression test suit for our implementation under test. This can be done using the `--result-folder` fuzz command argument which saves each permutation in its own file in the given folder. Each file is named by the MD5 sum of the content and is given the extension `.test` with the `--result-extension` fuzz command argument.

```bash
mkdir testset
tavor --format-file vending.tavor --max-repeat 2 fuzz --strategy AllPermutations --result-folder testset --result-extension ".test"
```

This command results into exactly **31** files created in the folder `testset` and concludes our fuzzing related work using the Tavor format.

## Executor

The executor connects the key-driven test files with the implementation under test. It reads, parses and validates one key-driven file, executes sequentially each key with its arguments by invoking actions of the implementation and validates these actions. A test passes if each key executes without any problem. We will first define the groundwork of the executor since it is not yet defined how the implementation can be contacted.

```go
import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
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
```

This program reads a file given as CLI argument, parses it according to our key-driven format rules and then executes each key with its parameters using the map `actions`. Since we did not fill the map yet, every key-driven file will fail. However, the groundwork of the executor is hereby done and we can move on to define the actions for our keys. Since this example should be kept simple, we will use a package directly as our implementation. It should be noted that the same mechanisms could be used to test implementations of external processes, web APIs or any other implementation.

The following code introduces the implementation of the given state machine which we will declare in its own package.

```go
package implementation

import (
	"errors"
)

const (
	Coin25 = 25
	Coin50 = 50
)

var (
	ErrUnknownCoin = errors.New("Unknown coin")
)

type VendingMachine struct {
	credit int
}

func NewVendingMachine() *VendingMachine {
	return &VendingMachine{
		credit: 0,
	}
}

func (v VendingMachine) Credit() int {
	return v.credit
}

func (v *VendingMachine) Coin(credit int) error {
	switch credit {
	case Coin25:
		v.credit += credit
	case Coin50:
		v.credit += credit
	default:
		return ErrUnknownCoin
	}

	return nil
}

func (v *VendingMachine) Vend() bool {
	if v.credit < 100 {
		return false
	}

	v.credit -= 100

	return true
}
```

This implementation can be now used in the executor to define actions to the defined keys.

```go
var (
	ErrInvalidParametersCount = errors.New("Invalid parmaters count")
)

func init() {
	machine := implementation.NewVendingMachine()

	actions["credit"] = func(parameters []string) error {
		if len(parameters) != 1 {
			return ErrInvalidParametersCount
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
			return ErrInvalidParametersCount
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
			return ErrInvalidParametersCount
		}

		vend := machine.Vend()
		if !vend {
			return fmt.Errorf("Could not vend")
		}

		return nil
	}
}
```

As you can see, each action does exactly what the key suggests. `credit` checks if the current credit of the vending machine is equal to the given argument, `coin` inserts a new coin and validates that the invoked action is successful and `vend` invokes the vending action and also validates that it was successful. Thereby each action to a key is isolated from one another.

Since we have now defined all components for testing the given state machine we can proceed to execute the actual tests.

## Test execution

Since all components have been defined and the test set has been generated we can execute single key-driven files via the executor.

```bash
go run executor.go testset/fba58bb35d28010b61c8004fadcb88a3
```

Results in the following output.

```
credit [0]
coin [50]
credit [50]
coin [50]
credit [100]
vend []
credit [0]
coin [50]
credit [50]
coin [25]
credit [75]
coin [25]
credit [100]
vend []
credit [0]
```

We can also look at the exit code of the program with the following command.

```bash
echo $?
```

Which results in the output `0`, meaning that the test has passed.

`go run` does compile the whole executor with each execution. This is very slow. It is therefore advisable to compile the executor with the following command to the current folder.

```bash
go build executor.go
```

This will create a binary called `executor` which we will use in the following examples.

Executing each key-driven file is tedious. A solution would be to extend the executor but this would also mean more restrictions and more flaw possibilities in the executor code. Alternatively a simple bash script which executes each key-driven file of the `testset` folder and immediately exits if a file fails can be used.

```bash
#!/bin/bash

shopt -s nullglob

for file in testset/*.test
do
	echo "Test $file"

	./executor $file

	if [ $? -ne 0 ]; then
		echo "Error detected, will exit loop"

		break
	fi
done
```

Executing this script reveals no error meaning all tests passed.
