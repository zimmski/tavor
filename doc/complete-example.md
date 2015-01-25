# A complete example

This example provides a complete overview of Tavor. It does not utilize every single feature but every component which should give you an idea on how Tavor can be used in your own projects. The example tests an implementation of the following state machine.

![Basic states and actions](/examples/complete/fsm.png "Basic states and actions")

The example uses the *Tavor format* to define this state machine and the *Tavor binary* to generate key-driven test files. An *executor* translates keys of a file into actions for the implementation under test. After successfully testing the original implementation, some bugs will be introduced to show the failure of some tests as well as automatically delta-debugging the failed key-driven test files.

> **Note:** The implementation of this example has intentional concurrency and other problems. Future versions of Tavor will help to identify, test and resolve such flaws. As Tavor evolves this example will also evolve. This also concerns the given state machine and Tavor format definition. Both could be defined much more efficiently with the help of state variables which are common in model-based testing but not yet fully implemented in the Tavor format. However, this could be easily implemented via code using the Tavor framework.

The following components will be defined and described in the following sections:

- [Tavor format](/examples/complete/vending.tavor)
- [Executor](/examples/complete/executor.go)
- [Implementation](/examples/complete/implementation/vending.go)
- [Bash script to run all key-driven files](/examples/complete/run.sh)
- [Bash script to run all key-driven files and reduce failed ones](/examples/complete/run-and-reduce.sh)

## <a name="table-of-content"></a>Table of content

- [Tavor format and fuzzing](#tavor-format-and-fuzzing)
- [Implementing an executor](#executor)
- [Test execution](#test-execution)
- [Introducing intentional bugs](#bugs)
- [Delta-debugging of inputs](#delta-debugging)

## <a name="tavor-format-and-fuzzing"></a>Tavor format and fuzzing

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
START = Credit0 *( Coin25 Credit25 | Coin50 Credit50 )

Credit0   = "credit" "\t"   0 "\n"
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

The default fuzzing strategy `random` can create all possible permutations of a format but since it is random, it will need enough time to do so. Lots of duplicated results will be generated, since even random events often lead to the same results. To work around this problem the `AllPermutations` strategy can be used which, as its name suggests, generates all possible permutations of a graph. This strategy should be used wisely since even small graphs can have [an enormous amount of permutations](https://en.wikipedia.org/wiki/Combinatorial_explosion). Since the example graph has a loop, we can state that there are an infinite amount of permutations. To work around this additional problem, the `--max-repeat` argument will be used with a suitable value. It enforces a maximum of loop traversals and repetitions.

> **Note:** There is no easy way to ensure that the `--max-repeat` value is correct. A high value can lead to many repetitive permutations which will often not improve the testing process. A small value can lead to a bad coverage which means that some code branches would not be taken. Future versions of Tavor will aid this process by implementing better fuzzing strategies as well as additional metrics like token and edge coverage.

Since the given state machine is pretty easy to understand one may imply that a loop should be generated at least twice to include the repetition part of each loop. In conclusion, this leads to the following arguments for the Tavor binary.

```bash
tavor --format-file vending.tavor --max-repeat 2 fuzz --strategy AllPermutations
```

This command outputs directly to STDOUT which is OK for one generation but since we are using the `AllPermutations` fuzzing strategy we have to deal with many generations of the format. Additionally we want to save every permutation in a file so we can build a regression test suite for our implementation under test. This can be done using the `--result-folder` fuzz command argument which saves each permutation in its own file in the given folder. Each file is named by the MD5 sum of the content and is given the extension `.test` by the `--result-extension` fuzz command argument.

```bash
mkdir testset
tavor --format-file vending.tavor --max-repeat 2 fuzz --strategy AllPermutations --result-folder testset --result-extension ".test"
```

This command results into exactly **31** files created in the folder `testset` and concludes our fuzzing related work using the Tavor format.

## <a name="executor"></a>Implementing an executor

The executor connects the key-driven test files with the implementation under test. It reads, parses and validates one key-driven file, executes sequentially each key with its arguments by invoking actions of the implementation and validates these actions. A test passes if each key executes without any problem. We will first define the groundwork of the executor since it is not yet defined how the implementation can be contacted. The executor will be written in Go, since all Tavor examples are written in Go. However, the generated key-driven test files are independent of the programming language which means that the executor could be implemented in any language too.

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

This program reads a file given as CLI argument, parses it according to our key-driven format rules and then executes each key with its parameters using the map `actions`. Since we did not fill the map yet, every key-driven file will fail. However, the groundwork of the executor is hereby done and we can move on to define the actions for our keys. Since this example should be kept simple, we will use an additional package as our implementation. It should be noted that the same mechanisms could be used to test implementations of external processes, web APIs or any other implementation as long as an interface can be used.

The following code introduces the implementation of the given state machine which we will declare in its own package.

```go
package implementation

import (
	"errors"
)

const (
	coin25 = 25
	coin50 = 50
)

var (
	// ErrUnknownCoin states that a given coin is unknown to the vending machine
	ErrUnknownCoin = errors.New("Unknown coin")
)

// VendingMachine holds the state of a vending machine
type VendingMachine struct {
	credit int
}

// NewVendingMachine returns a instantiated state of a vending machine
func NewVendingMachine() *VendingMachine {
	return &VendingMachine{
		credit: 0,
	}
}

// Credit returns the current credit of the vending machine
func (v VendingMachine) Credit() int {
	return v.credit
}

// Coin inserts a coin into the vending machine. On success the credit of the machine will be increased by the coin.
func (v *VendingMachine) Coin(credit int) error {
	switch credit {
	case coin25:
		v.credit += credit
	case coin50:
		v.credit += credit
	default:
		return ErrUnknownCoin
	}

	return nil
}

// Vend executes a vend of the machine if enough credit (100) has been put in and returns true.
func (v *VendingMachine) Vend() bool {
	if v.credit < 100 {
		return false
	}

	v.credit -= 100

	return true
}
```

This implementation can be now used in the executor to define actions for the given keys.

```go
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
```

As you can see, each action does exactly what its key suggests. `credit` checks if the current credit of the vending machine is equal to the given argument, `coin` inserts a new coin and validates that the invoked action is successful and `vend` invokes the vending action and also validates that it was successful. Thereby each action to a key is isolated from one another.

Since we have now defined all components for testing the given state machine, we can proceed to execute the actual tests.

## <a name="test-execution"></a>Test execution

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

`go run` does compile the whole executor for every execution. This is very slow and it is therefore advisable to compile the executor with the following command to the current folder.

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

Executing this script reveals no errors meaning all tests passed. Since this is not very exciting we will integrate in the next section some bugs into the implementation.

## <a name="bugs"></a>Introducing intentional bugs

The following subsections will introduce bugs into the implementation. Each bug will fail at least one test of our generated test set and will be studied isolated from other bugs. Meaning each section starts with a fresh original version of the implementation.

### <a name="bugs-coin"></a>The `Coin` method does not increase the credit

This bug can be introduced easily with one of the following code replacements for the `Coin` method of our implementation:

- Remove the addition statements.

	```go
	func (v *VendingMachine) Coin(credit int) error {
		switch credit {
		case coin25:
		case coin50:
		default:
			return ErrUnknownCoin
		}

		return nil
	}
	```
- Using a non-pointer type as receiver for the `Coin` method which will leave the state of the machine untouched.

	```go
	func (v VendingMachine) Coin(credit int) error {
		switch credit {
		case coin25:
			v.credit += credit
		case coin50:
			v.credit += credit
		default:
			return ErrUnknownCoin
		}

		return nil
	}
	```

Running our script to execute all key-driven files will immediately result in an failed test. For example with the file `testset/1f6b08c8273b8e46128e4d84e4e7e621.test`:

```
Test testset/1f6b08c8273b8e46128e4d84e4e7e621.test
credit [0]
coin [50]
credit [50]
Error: Credit should be 50 but was 0
Error detected, will exit loop
```

### <a name="bugs-vend"></a>The `Vend` method does not decrease the credit

Similar to the previous example we can modify the code to leave the `credit` member variable untouched with one of the following replacements of the `Vend` method:

- Remove the subtraction statements.

	```go
	func (v *VendingMachine) Vend() bool {
		if v.credit < 100 {
			return false
		}

		return true
	}
	```
- Using a non-pointer type as receiver for the `Vend` method which will leave the state of the machine untouched.

	```go
	func (v VendingMachine) Vend() bool {
		if v.credit < 100 {
			return false
		}

		v.credit -= 100

		return true
	}
	```

Running our script to execute all key-driven files will immediately result in an failed test. For example with the file `testset/1f6b08c8273b8e46128e4d84e4e7e621.test`:

```
Test testset/1f6b08c8273b8e46128e4d84e4e7e621.test
credit [0]
coin [50]
credit [50]
coin [25]
credit [75]
coin [25]
credit [100]
vend []
credit [0]
Error: Credit should be 0 but was 100
Error detected, will exit loop
```

### <a name="bugs-second-25-coin"></a>Every second 25 coin does not increase the credit

The bug type "works the first time but not the second" is very common in most programs. Since our vending machine implementation is too easy we have to introduce an additional state member variable to trigger such a bug. The following code snippets has to replace the original implementation:

```go
type VendingMachine struct {
	credit    int
	coinsOf25 int
}

func NewVendingMachine() *VendingMachine {
	return &VendingMachine{
		credit:    0,
		coinsOf25: 0,
	}
}

func (v *VendingMachine) Coin(credit int) error {
	switch credit {
	case coin25:
		if v.coinsOf25%2 == 0 {
			v.credit += credit
		}

		v.coinsOf25++
	case coin50:
		v.credit += credit
	default:
		return ErrUnknownCoin
	}

	return nil
}
```

Running our script to execute all key-driven files will immediately result in an failed test. For example with the file `testset/1f6b08c8273b8e46128e4d84e4e7e621.test`:

```
Test testset/1f6b08c8273b8e46128e4d84e4e7e621.test
credit [0]
coin [50]
credit [50]
coin [25]
credit [75]
coin [25]
credit [100]
Error: Credit should be 100 but was 75
Error detected, will exit loop
```

This especially interesting with a long running key-driven file which inserts two 25 coins in a second or third vending loop. For example the file `testset/fba58bb35d28010b61c8004fadcb88a3.test` triggers the bug in the second loop.

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
Error: Credit should be 100 but was 75
```

This is an interesting test case since the first iteration of the vending loop is not relevant to the bug. It shows that actions which trigger a flaw most not be a minimal set of actions but they can be reduced to such a set. This is one of the major operations of Tavor which is called **delta-debugging** or in general **reducing**. The next main section will cover how Tavor can be used to reduce an input to its minimum.

## <a name="delta-debugging"></a>Delta-debugging of inputs

**Delta-debugging** or in general **reducing** is a method to reduce data to ideally its minimum while still complying to defined constraints. In our example the data is a key-driven test file which fails and the constraint is that the reduced test case should still fail. Therefore the final result of the delta-debugging process should be a minimal test case which still triggers the same bug as the original test case. This can be automatically or semi-automatically done by the `reduce` command of the Tavor binary. The binary uses our Tavor format file to parse and validate the given key-driven file and tries to reduce its data according to rules defined by the format file. For instance optional content like repetitions can be reduced to a minimal repetition. In our example we can reduce the iterations of the vending loop.

We will use the bug and the key-driven test file `testset/fba58bb35d28010b61c8004fadcb88a3.test` which were introduced in [one of the subsections of "Introducing bugs"](#bugs-second-25-coin). The file has the following content.

```
credit	0
coin	50
credit	50
coin	50
credit	100
vend
credit	0
coin	50
credit	50
coin	25
credit	75
coin	25
credit	100
vend
credit	0
```

The introduced bug will be triggered in the second vending iteration. Every second 25 coin does not increase the machine's credit counter. This can be easily tested with our generated test set but the given file shows that there are key-driven files for this bug that could be reduced because of unnecessary loops.

We will first use the semi-automatic method of the Tavor `reduce` command. The given format file will be used to reduce the given input. Every reduction step displays the question "Do the constraints of the original input still hold for this generation?" to the user. The user's task is to inspect and validate the reduced output of the original data and decide by giving feedback if the bug is triggered (**yes**) or not (**no**). The following command starts this process.

```bash
tavor --format-file vending.tavor reduce --input-file testset/fba58bb35d28010b61c8004fadcb88a3.test
```

This should result in the following output and feedbacks.

```
credit  0


Do the constraints of the original input still hold for this generation? [yes|no]: no
credit  0
coin    50
credit  50
coin    50
credit  100
vend
credit  0


Do the constraints of the original input still hold for this generation? [yes|no]: no
credit  0
coin    50
credit  50
coin    25
credit  75
coin    25
credit  100
vend
credit  0


Do the constraints of the original input still hold for this generation? [yes|no]: yes
credit  0
coin    50
credit  50
coin    25
credit  75
coin    25
credit  100
vend
credit  0
```

The last reduction output is the minimum which still triggers the same bug as the original test case. Additionally it is shown that the default reduce strategy of the Tavor `reduce` command tries to output the smallest generation first which is simply the `credit  0` command.

This semi-automatic process can be tedious for big data especially because of the manual validation. The Tavor binary does therefore provide several methods to reduce completely automatically. Since we already have a executor which tests key-driven files we can use it in this process. This is additionally aided by the executor which exits with different exit status codes on success or failure. We can therefore conclude that a reduced generation of our original failing key-driven file has to have the same exit status code as the original one. This can be automatically done by the following command. Which uses the executor to validate reduced data which is temporary written to a file. Each exit status code of the executor is compared to the original exit status code. If it is not equal, the reduction process will try an alternative reduction step until a reduction path is found which leads to the minimum.

```bash
tavor --format-file vending.tavor reduce --input-file testset/fba58bb35d28010b61c8004fadcb88a3.test --exec "./executor TAVOR_DD_FILE" --exec-argument-type argument --exec-exact-exit-code
```

Which results into the following output.

```
credit  0
coin    50
credit  50
coin    25
credit  75
coin    25
credit  100
vend
credit  0
```

As you can see this is the minimum which still triggers the same bug as the original key-driven file.

Reducing key-driven test files allows developers to always debug with the minimum set of actions to trigger a bug which can save a lot of debugging time. It is therefore a handy addition to the execution of a test suite. We can modify our bash script to automatically reduce failed files.

```bash
#!/bin/bash

shopt -s nullglob

for file in testset/*.test
do
	echo "Test $file"

	./executor $file

	if [ $? -ne 0 ]; then
		echo "Error detected."

		echo "Reduce original file to its minimum."

		tavor --format-file vending.tavor reduce --input-file $file --exec "./executor TAVOR_DD_FILE" --exec-argument-type argument --exec-exact-exit-code > $file.reduced

		echo "Saved to $file.reduced"

		break
	fi
done
```

This script will run the executor with every key-driven test file of the `testset` folder and stop at the first failed file. The failed file will be then reduced to its minimum which will be saved next to the original file with the extension `.reduced`.

Executing this script with the introduced bug will for example result in the following output.

```
Test testset/1f6b08c8273b8e46128e4d84e4e7e621.test
credit [0]
coin [50]
credit [50]
coin [25]
credit [75]
coin [25]
credit [100]
Error: Credit should be 100 but was 75
Error detected.
Reduce original file to its minimum.
Saved to testset/1f6b08c8273b8e46128e4d84e4e7e621.test.reduced
```
