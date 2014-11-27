# A complete example

This example provides a complete overview of Tavor. It does not utilize every feature but every component which should give you an idea on how Tavor can be used in your own projects. The example tests an implementation of the following state machine.

![Basic states and actions](/examples/quick/basic.png "Basic states and actions")

It uses the *Tavor format* to define this state machine and the *Tavor binary* to generate key-driven test files. An *executor* translates keys of a file into actions for the implementation under test. After successfully testing the original implementation, some bugs will be introduced to show the failure of some tests as well as automatically delta-debugging the failed key-driven test files.

## Tavor format

The following Tavor format, saved as [vending.tavor](/examples/complete/vending.tavor), defines a model of all possible valid states of the given state machine. Since this model should generate key-driven files, a set of keys and a key-driven format have to be defined.

The following keys will be used:

- **credit** performs a validation of the current credit amount. It compares the current credit against the integer argument of the key.
- **coin** invokes the action to insert a coin. It uses the integer argument as the coin credit.
- **vend** invokes the action to vend.

The key-driven format is defined as the following.
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

Every execution will output in one generation of the format like the following.

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

This command outputs directly to STDOUT which is OK for one generation but since we are using the `AllPermutations` fuzzing strategy we have to deal with many generations of the format. Additionally we want to save every permutation in a file so we can build a regression test suit for our implementation under test. This can be done by the `--result-folder` fuzz command argument which saves each permutation in its own file in the given folder. Each file is named by the MD5 sum of the content and is given the extension `.test` by default.

```bash
mkdir testset
tavor --format-file vending.tavor --max-repeat 2 fuzz --strategy AllPermutations --result-folder testset
```

This command results into exactly **31** files created in the folder `testset` and concludes our fuzzing related work using the Tavor format.

## Executor

The executor connects the key-driven test files with the implementation under test. It reads, parses and validates one key-driven file, executes sequentially each key with its arguments by invoking actions of the implementation and validates these actions. We will first define the groundwork of the executor since it is not yet defined how the implementation can be contacted.
