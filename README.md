# Tavor [![GoDoc](https://godoc.org/github.com/zimmski/tavor?status.png)](https://godoc.org/github.com/zimmski/tavor) [![Build Status](https://travis-ci.org/zimmski/tavor.svg?branch=master)](https://travis-ci.org/zimmski/tavor) [![Coverage Status](https://coveralls.io/repos/zimmski/tavor/badge.png)](https://coveralls.io/r/zimmski/tavor)

Tavor ([Sindarin](https://en.wikipedia.org/wiki/Sindarin) for "woodpecker") is a framework for implementing and doing everyday [fuzzing](#fuzzing) and [delta-debugging](#delta-debugging) as well as doing research on new methods without reimplementing basic algorithms. An [EBNF-like notation](#format) allows the definition of data (e.g. file formats and protocols) without the need of programming. Tavor also relaxes on the definitions of fuzzing and delta-debugging allowing the user to utilize implemented algorithms universally e.g. for key-driven testing, model-based testing, simulating user-behavior and genetic programming.

### <a name="quick-example"></a>A quick example

Imagine a vending machine which ejects a product after receiving 100 worth of credits. It is possible to input 25 and 50 credit coins into the machine. After receiving enough credits the machine ejects a product and resets the credit counter to zero. To keep it simple, we specify that the machine does not handle credit overflows. A representation of the states and actions of the machine could look like this:

![Basic states and actions](examples/quick/basic.png "Basic states and actions")

This state machine can be defined using the following [Tavor format](#format):

```tavor
START = Credit0

Credit0   = "Credit0"   "\n" ( Coin25 Credit25 | Coin50 Credit50 | )
Credit25  = "Credit25"  "\n" ( Coin25 Credit50 | Coin50 Credit75 )
Credit50  = "Credit50"  "\n" ( Coin25 Credit75 | Coin50 Credit100 )
Credit75  = "Credit75"  "\n" Coin25 Credit100
Credit100 = "Credit100" "\n" Vend Credit0

Coin25 = "Coin25" "\n"
Coin50 = "Coin50" "\n"

Vend = "Vend" "\n"
```
(Please note that the new line escape sequences "\n" are just defined to make the output prettier)

You can download this file called [`basic.tavor` from here](examples/quick/basic.tavor).

Now we can use Tavor to [fuzz](#fuzzing) the format by issuing the following command:

```bash
tavor --format-file basic.tavor fuzz
```

On every call this command outputs random paths through the defined graph. Here are some example outputs:

```
Credit0
```

```
Credit0
Coin50
Credit50
Coin50
Credit100
Vend
Credit0
```

```
Credit0
Coin25
Credit25
Coin50
Credit75
Coin25
Credit100
Vend
Credit0
```

Generating data like this is just one example of the capabilities of Tavor. Please have a look at [the bigger example](#bigexample) with a complete overview over the basic features or keep reading to find out more about the background and capabilities of Tavor.

Additionally you can find functional Tavor format files and fuzzer applications at [https://github.com/zimmski/fuzzer](https://github.com/zimmski/fuzzer).

## <a name="table-of-content"></a>Table of content

- [What is fuzzing?](#fuzzing)
- [What is delta-debugging?](#delta-debugging)
- [What does Tavor provide and how does it work?](#tavor-provides)
- [The Tavor format](#format)
- [How do I use Tavor?](#use)
- [The Tavor binary](#binary)
- [A complete example for fuzzing, executing and delta-debugging](#bigexample)
- [Where are the precompiled binaries?](#precompiled)
- [How do I build Tavor?](#build)
- [How do I develop applications with the Tavor framework?](#develop)
- [How do I extend the Tavor framework?](#extend)
- [How stable is Tavor?](#stability)
- [Missing features](#missing-features)
- [Can I make feature requests and report bugs and problems?](#feature-request)

TODO update table of content with header3 sections when done

## <a name="fuzzing"></a>What is fuzzing?

> Fuzz testing or fuzzing is a software testing technique, often automated or semi-automated, that involves providing invalid, unexpected, or random data to the inputs of a computer program. The program is then monitored for exceptions such as crashes, or failing built-in code assertions or for finding potential memory leaks. Fuzzing is commonly used to test for security problems in software or computer systems.
> -- <cite>[https://en.wikipedia.org/wiki/Fuzz_testing](https://en.wikipedia.org/wiki/Fuzz_testing)</cite>

Although this is the common definition of fuzzing, it is nowadays often just one view on capabilities of fuzzing tools. Fuzzing is in general just the generation of data and it does not matter if it is invalid or valid and what the type of data (e.g. files, protocol data) is. The use case of the data itself is also often broadly defined as it can be used to test algorithms, programs or hardware but it can be practically used everywhere where data is needed.

Fuzzing algorithms can be categorized into two areas:

- **mutation-based**

	Mutation-based fuzzing takes existing data and simply changes it. This often leads to invalid data, as most techniques are not obeying constraints nor rules of the underlying data.

	Some common technique for mutation-based fuzzing are:
	- Bit flipping: Random chosen bits of the data are flipped.
	- Prepending/Appending: New data is prepended/appended to the given data.
	- Repeating: Random chosen parts of the given data are repeated.
	- Removal: Random chosen parts of the given data are removed.

- **generation-based**

	Generation-based algorithms have one big advantage over mutation-based in that they have to understand and obey the underlying constraints and rules of the data itself. This property can be used to generate valid as well as invalid data. Another property is that generation-based algorithms generate data from scratch which eliminates the need for gathering data and keeping it up to date.

	There are no common techniques for generation-based fuzzing but most algorithms choose a graph as underlying representation of the data model. The graph is then traversed and each node outputs a part of the data. The traversal algorithms and the complexity and abilities of the data model like constraints between nodes or adding nodes during the traversal distinguish generation-based fuzzers and contributes in general to their mightiness.

## <a name="delta-debugging"></a>What is delta-debugging?

> The Delta Debugging algorithm isolates failure causes automatically - by systematically narrowing down failure-inducing circumstances until a minimal set remains.
> -- <cite>[https://en.wikipedia.org/wiki/Delta_Debugging](https://en.wikipedia.org/wiki/Delta_Debugging)</cite>

E.g. we feed a given data to a program which fails on executing. By delta-debugging this data we can reduce it to its minimum while still failing the execution. The reduction of the data is handled by software heuristics (semi-)automatically. The obvious advantage of this method, besides being done (semi-)automatically, is that we do not need to then handle uninteresting parts of the data while debugging the problem, we can focus on the important parts which actually lead to the failure.

**Note**: Since delta-debugging reduces data it is also called `reducing`.

Delta-debugging consists of three areas:

- A heuristic has to decide which parts of the data will be reduced next
- The reduction itself e.g.
	+ Remove repetitions
	+ Remove optional data
	+ Replace data with something else e.g. replace an uninteresting complex function with a constant value
- Testing the new resulting data concerning the failure

Although delta-debugging is described as method to isolate failure causes, it can be also used to isolate anything given isolating constraints. For example we could reduce an input for a program which leads to a positive outcome to its minimum.

## <a name="tavor-provides"></a>What does Tavor provide and how does it work?

Tavor combines both fuzzing and delta-debugging by allowing all implemented methods to operate on one internal model-based structure represented by a graph. This structure can be defined and generated via code or by using a format file. Out of the box Tavor comes with its own [format](#format) which covers all functionality of the framework.

Tavor's generic fuzzing implementation is not fixed to one technique. Instead different fuzzing techniques and heuristics can be implemented and executed independently as [Tavor fuzzing strategies](#fuzzing-strategy). The same principle is used for delta-debugging where so called [Tavor reduce strategies](#reduce-strategy) can be implemented and used. Both types of strategies operate on the same internal structure independent of the format. This structure is basically a graph of nodes which are called [tokens](#token) throughout the Tavor framework. The structure itself is not fixed to a static definition but can be changed by so called [fuzzing filters](#fuzzing-filter).

Even tough Tavor provides loads of functionality out of the box, a lot is still missing. A list of missing but planed features can be found in the [missing features section](#missing-features). For feature requests please have a look at the [feature request section](#feature-request).

### <a name="token"></a>What are tokens?

Tavor's tokens differ from *lexical analysis tokens* in that they represent not just a group of characters but different kind of data with additional properties and abilities. Tokens can be constant integers and strings of all kind but also dynamic data like integer ranges, sequences and character classes. Furthermore tokens can encapsulate other tokens to not only group them together but to create building blocks that can be reused to, for example, repeat a group of tokens. Tokens can have states, conditions and logic. They can create new tokens dynamically and can depend on other tokens to generate data. Tavor's tokens are basically the foundation of the whole framework and every algorithm for fuzzing, parsing and delta-debugging uses them.

If you want to know more about Tavor's tokens you can read through [Tavor's format definition](#format) or you can read about them in depth in the [developing](#develop) and [extending](#extend) sections.

### <a name="fuzzing-strategy"></a>What are fuzzing strategies?

Each fuzzing strategy represents one fuzzing technique. This can be a heuristic for walking through the internal structure, how tokens of the structure are fuzzed or even both. Tavor currently distinguishes between two techniques of token fuzzing. One is to deterministically choose one possible permutation of the token the other is choosing randomly out of all permutations of a token.

An example for a fuzzing strategy is the [random fuzzing strategy](https://godoc.org/github.com/zimmski/tavor/fuzz/strategy#RandomStrategy) which is Tavor's default. It traverses through the whole internal structure and randomly permutates each token.

Please have a look at [the documentation](https://godoc.org/github.com/zimmski/tavor/fuzz/strategy) for an overview of all officially available fuzzing strategies of Tavor.

### <a name="fuzzing-filter"></a>What are fuzzing filters?

Fuzzing filters mutate the internal structure and can be applied after it is ready for fuzzing thus after creating it e.g. after parsing and unrolling. This can be associated to [mutation-based fuzzing](#fuzzing) where not the generating structure but the data itself is mutated.

An example use-case for fuzzing filters is the [boundary-value analysis](https://en.wikipedia.org/wiki/Boundary-value_analysis) software testing technique. Imagine a function which should be tested having one integer which has valid values from 1 to 100. This would lead to 100 possible values which have to be tested just for this one integer and thus to at least 100 permutations of the internal structure. Boundary-value analysis reduces these permutations to e.g. 1, 50 and 100 so just three instead of 100 cases. This is exactly what the [PositiveBoundaryValueAnalysis fuzzing filter](https://godoc.org/github.com/zimmski/tavor/fuzz/filter#PositiveBoundaryValueAnalysisFilter) does. It traverses the whole internal structure and replaces every range token with its boundary values.

Please have a look at [the documentation](https://godoc.org/github.com/zimmski/tavor/fuzz/filter) for an overview of all officially available fuzzing filters of Tavor.

### <a name="reduce-strategy"></a>What are reduce strategies?

Reduce strategies are strongly comparable to fuzzing strategies. Each reduce strategy represents one reduce/delta-debugging technique. This can be a heuristic for walking through the internal structure, how tokens of the structure are reduced or even both. The reduction method is depending on the token type. For example a constant integer cannot be reduced any further but a repetition of optional strings can be minimized or even left out.

Please have a look at [the documentation](https://godoc.org/github.com/zimmski/tavor/reduce/strategy) for an overview of all officially available reduce strategies of Tavor.

### <a name="unrolling"></a>Unrolling loops

Although the internal structure allows loops in its graph, Tavor currently unrolls loops for easier algorithm implementations and usage. A future version will supplement this by allowing loops by default.

This graph for example loops between the states `Idle` and `Action`:

![Looping](/doc/images/README/unroll-loop.png "Looping")

Unrolling the graph results in the following internal graph given a maximum of two repetitions:

![Unrolled](/doc/images/README/unroll-unrolled.png "Unrolled")

## <a name="format"></a>The Tavor format

The Tavor format documentation has its own [page which can be found here](/doc/format.md).

## <a name="use"></a>How do I use Tavor?

Tavor can be used in three different ways:
- [Using the binary](#binary) which makes everything officially provided by the Tavor framework available via the command line.
- [Developing applications with the Tavor framework](#develop) by implementing the internal structure via code and doing everything else like fuzzing and delta-debugging too via code.
- [Extending the Tavor framework](#extend) because of research or missing features.

## <a name="binary"></a>The Tavor binary

The [Tavor binary](#precompiled) provides fuzzing and delta-debugging functionality for Tavor format files as well as some other commands. Sane default arguments should provide a pleasant experience.

Since the binary acts on Tavor format files, the `--format-file` argument has to be used for every non-informational action. E.g. the following command fuzzes the given format file with the default fuzzing strategy:

```bash
tavor --format-file file.tavor fuzz
```

In contrast listing all available fuzzing strategies does not require the `--format-file` argument:

```bash
tavor fuzz --list-strategies
```

To learn more about available arguments and commands, you can invoke the binary's help by executing the binary without any arguments or with the `--help` argument.

Here is a complete overview of all arguments, commands and their options:

```
Usage:
  tavor [options] <command> [command options]

General options:
  --debug             Debug log output
  --help              Show this help message
  --verbose           Verbose log output
  --version           Print the version of this program

Global options:
  --seed=             Seed for all the randomness
  --max-repeat=       How many times loops and repetitions should be repeated (2)

Format file options:
  --check             Just check the syntax of the format file and exit
  --format-file=      Input tavor format file
  --print             Prints the AST of the parsed format file
  --print-internal    Prints the internal AST of the parsed format file

Available commands:
  fuzz      Fuzz the given format file
  graph     Generate a DOT file out of the internal AST
  reduce    Reduce the given input file
  validate  Validate the given input file

[fuzz command options]
      --exec=                                    Execute this binary with possible arguments to test a generation
      --exec-exact-exit-code=                    Same exit code has to be present (-1)
      --exec-exact-stderr=                       Same stderr output has to be present
      --exec-exact-stdout=                       Same stdout output has to be present
      --exec-match-stderr=                       Searches through stderr via the given regex. A match has to be present
      --exec-match-stdout=                       Searches through stdout via the given regex. A match has to be present
      --exec-do-not-remove-tmp-files             If set, tmp files are not removed
      --exec-do-not-remove-tmp-files-on-error    If set, tmp files are not removed on error
      --exec-argument-type=                      How the generation is given to the binary (stdin)
      --list-exec-argument-types                 List all available exec argument types
      --script=                                  Execute this binary which gets fed with the generation and should return feedback
      --exit-on-error                            Exit if an execution fails
      --filter=                                  Fuzzing filter to apply
      --list-filters                             List all available fuzzing filters
      --strategy=                                The fuzzing strategy (random)
      --list-strategies                          List all available fuzzing strategies
      --result-folder=                           Save every fuzzing result with the MD5 checksum as filename in this folder
      --result-extension=                        If result-folder is used this will be the extension of every filename
      --result-separator=                        Separates result outputs of each fuzzing step ("\n")

[graph command options]
      --filter=         Fuzzing filter to apply
      --list-filters    List all available fuzzing filters

[reduce command options]
      --exec=                           Execute this binary with possible arguments to test a generation
      --exec-exact-exit-code            Same exit code has to be present
      --exec-exact-stderr               Same stderr output has to be present
      --exec-exact-stdout               Same stdout output has to be present
      --exec-match-stderr=              Searches through stderr via the given regex. A match has to be present
      --exec-match-stdout=              Searches through stdout via the given regex. A match has to be present
      --exec-do-not-remove-tmp-files    If set, tmp files are not removed
      --exec-argument-type=             How the generation is given to the binary (stdin)
      --list-exec-argument-types        List all available exec argument types
      --script=                         Execute this binary which gets fed with the generation and should return feedback
      --input-file=                     Input file which gets parsed, validated and delta-debugged via the format file
      --strategy=                       The reducing strategy (BinarySearch)
      --list-strategies                 List all available reducing strategies
      --result-separator=               Separates result outputs of each reducing step ("\n")

[validate command options]
      --input-file=   Input file which gets parsed and validated via the format file
```

### General options

The Tavor binary provides different kinds of general options which are informative or apply to all commands. Besides the `--format-file` general format option the following options are noteworthy:

- **--max-repeat** sets the maximum repetition of loops and repeating tokens. If not set, the default value (currently 2) is used. 0, meaning no maximum repetition, is currently now allowed because of the limitation mentioned [here](#unrolling).
- **--seed** defines the seed for all random generators. If not set, a random value will be chosen. This argument makes the execution of every Tavor command deterministic. Meaning that a result or failure can be replayed with the same `--seed` argument, other arguments and Tavor version.
- **--verbose** switches Tavor into verbose mode which prints additional information, like the used seed, to STDERR.

Please have a look at the help for more options and descriptions:

```bash
tavor --help
```

### Command: `fuzz`

The `fuzz` command generates data using the given format file and prints it directly to STDOUT.

```bash
tavor --format-file file.tavor fuzz
```

By default the `random` fuzzing strategy is used which can be altered using the `--strategy` fuzz command option.

```bash
tavor --format-file file.tavor fuzz --strategy AllPermutations
```

Fuzzing filters can be applied before the fuzzing generation using the `--filter` fuzz command option. Filters are applied in the same order as they are defined, meaning from left to right.

The following command will apply the `PositiveBoundaryValueAnalysis` fuzzing filter and then the `NegativeBoundaryValueAnalysis`:

```bash
tavor --format-file file.tavor fuzz --filter PositiveBoundaryValueAnalysis --filter NegativeBoundaryValueAnalysis
```

Alternatively to printing to STDOUT an executable (or script) can be fed with the generated data. You can find examples for executables and scripts [here](/examples/fuzzing). There are two types of arguments to execute commands:

- #### exec

	Executes a given command for every data generation. The validation of the data can be done via the executable and by using additional `--exec-*` fuzz command options.

	For example the following command will execute a binary called `validate` with the default exec settings which feed the generation via STDIN to the started process and apply no validation at all.

	```bash
	tavor --format-file file.tavor fuzz --exec validate
	```

	The method of feeding the generation to the process can be changed using the `--exec-argument-type` fuzz command option. The following command puts the generation into a temporary file which is defined using the environment variable `TAVOR_FUZZ_FILE`:

	```bash
	tavor --format-file file.tavor fuzz --exec validate --exec-argument-type environment
	```

	The fuzz command allows to validate the execution using additional `--exec-*` fuzz command options. For example the exit code of the process can be validated to be `0` using the `--exec-exact-exit-code` fuzz command option:

	```bash
	tavor --format-file file.tavor fuzz --exec validate --exec-exact-exit-code 0
	```

- #### script

	Executes a given command and feeds every data generation to the running process using STDIN. Feedback is read using STDOUT. The running process can therefore control the fuzzing generation while it has to do all validation on its own.

	The following command will execute a binary called `validate`:

	```bash
	tavor --format-file file.tavor fuzz --script validate
	```

	Feedback commands control the fuzzing generation and are read by Tavor using STDOUT of the running process. Each command has to end with a new line delimiter and exactly one command has to be given for every generation.

	- **YES** reports a positive outcome for the generation.
	- **NO** reports a negative outcome for the generation. This is an error and will terminate the fuzzing generation if the `--exit-on-error` fuzz command option is used.

`--result-*` is an additional fuzz command option kind which can be used to influence the fuzzing generation itself. For example the `--result-separator` fuzz command option changes the separator of the generations if they are printed to STDOUT. The following command will use `@@@@` instead of the default `\n` separator to feed the fuzzing generations to the running process:

```bash
tavor --format-file file.tavor fuzz --script validate --result-separator "@@@@"
```

Please have a look at the fuzz command help for more options and descriptions:

```bash
tavor --help fuzz
```

### Command: `graph`

The `graph` command prints out a graph of the internal structure. This is needed since textual formats like the [Tavor format](#format) can be often difficult to mentally visualize. Currently only the DOT format is supported therefore third-party tools like [Graphviz](http://graphviz.org/) have to be used to convert the DOT data to other formats like JPEG, PNG or SVG.

The following command prints the DOT graph of a format file to STDOUT:

```bash
tavor --format-file file.tavor graph
```

To save the graph to a file the output of the command can be simply redirected:

```bash
tavor --format-file file.tavor graph > file.dot
```

The output can also be piped directly into the `dot` command of [Graphviz](http://graphviz.org/):

```
tavor --format-file file.tavor graph | dot -Tsvg -o outfile.svg
```

Please have a look at the graph command help for more options and descriptions:

```bash
tavor --help graph
```

To define the graph notation, the following image will be explained:

![Graph notation](doc/images/README/graph-notaton.png "Graph notation")

- Each circle represents a token of the internal structure (a to f)
- Arrows represents connections between tokens, meaning the token the arrow points to can come next (e.g a to c, c to d, d to f)
- Dotted arrows represent optional connections (a to b)
- The small dot is the start of the whole graph (arrow to a)
- Double bordered circles represent end-state tokens (f)

### Command: `reduce`

The `reduce` command applies delta-debugging to a given input according to the given format file. The reduction generates reduced generations of the original input which have to be tested either by the user or a program. Every generation has to correspond to the given format file which implies that the original input has to be valid too. This is checked using the same mechanisms as used by the `validate` command.

By default the reduction generation is printed to STDOUT and feedback is given through STDIN.

```bash
tavor --format-file file.tavor reduce --input-file file.input
```

By default the `BinarySearch` reduce strategy is used which can be altered using the `--strategy` reduce command option.

```bash
tavor --format-file file.tavor reduce --input-file file.input --strategy random
```

Alternatively to printing to STDOUT an executable (or script) can be fed with the generated data. You can find examples for executables and scripts [here](/examples/deltadebugging). There are two types of arguments to execute commands:

- #### exec

	Executes a given command for every generation. The validation of the data can be done via the executable and by using additional `--exec-\*` reduce command options. At least one `--exec-\*` matcher must be used to validate the reduced generations.

	For example the following command will execute a binary called `validate` with the default exec settings which feed the generation via STDIN to the started process. The `--exec-exact-exit-code` reduce command option is used to validate that the exit code of the original data matches the exit codes of reduce generations.

	```bash
	tavor --format-file file.tavor reduce --input-file file.input --exec validate --exec-exact-exit-code
	```

	The method of feeding the generation to the process can be changed using the `--exec-argument-type` reduce command option. The following command puts the generation into a temporary file which is defined using the environment variable `TAVOR_DD_FILE`:

	```bash
	tavor --format-file file.tavor reduce --input-file file.input --exec validate --exec-exact-exit-code --exec-argument-type environment
	```

	The reduce command allows to validate the execution using additional `--exec-*` reduce command options. For example the exit code and STDERR can be validated to match the original input using the `--exec-exact-exit-code` and `--exec-exact-stderr` reduce command options:

	```bash
	tavor --format-file file.tavor reduce --input-file file.input --exec validate --exec-exact-exit-code --exec-exact-stderr
	```

- #### script

	Executes a given command and feeds every data generation to the running process using STDIN. Feedback is read using STDOUT. The running process can therefore control the reduce generation while it has to do all validation on its own.

	The following command will execute a binary called `validate`:

	```bash
	tavor --format-file file.tavor reduce --input-file file.input --script validate
	```

	Feedback commands control the reducing generation and are read by Tavor using STDOUT of the running process. Each command has to end with a new line delimiter and exactly one command has to be given for every generation.

	- **YES** reports a positive outcome for the generation. The reduce strategy will therefore continue reducing this generation.
	- **NO** reports a negative outcome for the generation. This is an error and the reduce strategy will therefore produce a different generation.

`--result-*` is an additional reduce command option kind which can be used to influence the reduce generation itself. For example the `--result-separator` reduce command option changes the separator of the generations if they are printed to STDOUT. The following command will use `@@@@` instead of the default `\n` separator to feed the reduce generations to the running process:

```bash
tavor --format-file file.tavor reduce --input-file file.input --script validate --result-separator "@@@@"
```

Please have a look at the reduce command help for more options and descriptions:

```bash
tavor --help reduce
```

### Command: `validate`

The `validate` command validates a given input file according to the given format file. This can be helpful since this is for instance needed for the `reduce` command which does apply delta-debugging only on valid inputs or in the general case to check an input which was not generated through the given format file.

```bash
tavor --format-file file.tavor validate --input-file file.input
```

Please have a look at the validate command help for more options and descriptions:

```bash
tavor --help validate
```

## <a name="bigexample"></a>A complete example for fuzzing, executing and delta-debugging

TODO this example should give a complete overview of how Tavor can be used.<br/>
TODO do a key-word driven format-file<br/>
TODO executor for the key-words<br/>
TODO delta-debug keywords because of an intentional error<br/>

## <a name="precompiled"></a>Where are the precompiled binaries?

You can find all precompiled binaries on the [release page](https://github.com/zimmski/tavor/releases). The binaries are packed into archives that currently only hold the Tavor binary itself. Have a look at the [How do I build Tavor?](#build) section if you like to compile Tavor yourself.

### <a name="bash-completion"></a>Bash Completion

If you like Bash Completion for Tavor make sure that you have Bash Completion installed and then copy the [bash completion Tavor script](https://raw.githubusercontent.com/zimmski/tavor/master/cmd/tavor-bash_completion.sh) into your Bash Completion folder.

```bash
mkdir -p $HOME/.bash_completion
wget -P $HOME/.bash_completion https://raw.githubusercontent.com/zimmski/tavor/master/cmd/tavor-bash_completion.sh
. ~/.bashrc
```

Bash Completion for Tavor should now be working. If not, one reason could be that your distribution does not include user defined Bash Completion scripts in .bashrc so just add it to your .bashrc:

```bash
echo ". ~/.bash_completion/tavor-bash_completion.sh" >> ~/.bashrc
. ~/.bashrc
```

## <a name="build"></a>How do I build Tavor?

Tavor provides [precompiled 64 bit Linux binaries](#precompiled). Other architectures are currently not supported, but might work. Please have a look at the [feature request section](#feature-request) if you need them to work or you want more binaries.

If you do not want to use the [precompiled binaries](#precompiled) but instead want to compile Tavor from scratch, just follow the these steps (NOTE: All steps must execute without any errors):

1. Install and configure Go.

	At least version 1.3 must be used. Your distribution will most definitely have some packages or you can be brave and just install it yourself. Have a look at [the official documentation](http://golang.org/doc/install). Good luck!

2. Go-get Tavor

	```bash
	go get -t -v github.com/zimmski/tavor
	```

3. Install dependencies

	```bash
	cd $GOPATH/src/github.com/zimmski/tavor
	make dependencies
	```

3. Compile

	```bash
	cd $GOPATH/src/github.com/zimmski/tavor
	make install
	```

4. Run tests

	```bash
	cd $GOPATH/src/github.com/zimmski/tavor
	make test
	```

You now have a binary "tavor" in your `$GOPATH/bin` (or if set `$GOBIN` folder) folder which can be used without any further actions.

## <a name="develop"></a>How do I develop applications with the Tavor framework?

TODO<br/>
TODO create a header outline<br/>
TODO explain creating internal structures (instead of using a format file) with examples<br/>
TODO explain how to use filters, fuzzers and delta debugging<br/>

## <a name="extend"></a>How do I extend the Tavor framework?

If the [Tavor format](#format) and the [implemented functionality of the framework](#develop) is not enough to implement your applications needs, you can easily extend and change the Tavor framework. The following sections will provide starting points, hints and conventions to help you write your own Tavor extensions like fuzzing and reduce strategies or even your own tokens.

Since implementing new extensions and doing changes is trickier than using the existing framework, it is advisable to read the code documentation, which can be found in a nice representation on [https://godoc.org/github.com/zimmski/tavor/](https://godoc.org/github.com/zimmski/tavor/), and of course the actual code.

Code for the Tavor framework has to be deterministic. This means that no functionality is allowed to have its own source or seed of randomness. Methods of interfaces that define a random generator have to be implemented deterministically so that the same random seed will always result in the same result. This also applies to hand written tests and code who is concurrent.

If you are aiming to get your extensions and changes offically incorporated into the Tavor framework, please **first** [create an issue](https://github.com/zimmski/tavor/issues/new) in the issue tracker and discuss your implementation goals and plans with an outline. Please note that every feature and change has to be tested with handwritten tests, so please include a test plan in your outline too.

If extending Tavor yourself is not for you, but you still need new features, you can take a look at the [feature request section](#feature-request).

### Fuzzing filters [![GoDoc](https://godoc.org/github.com/zimmski/tavor?status.png)](https://godoc.org/github.com/zimmski/tavor/fuzz/filter)

The fuzzing filter code and all officially implemented fuzzing filters can be found in the [github.com/zimmski/tavor/fuzz/filter](/fuzz/filter) package and its sub-packages.

A fuzzing filter has to implement the `Filter` interface which is exported by the [github.com/zimmski/tavor/fuzz/filter](/fuzz/filter) package. The interface defines the `Apply` method that applies the filter onto a token which is passed to the method. The method's concern is therefore only one token at a time. The error return argument is not nil, if an error is encountered during the filter execution. On success a replacement for the token is returned. This can be either `nil`, meaning the token should not be replaced, or a slice of tokens which will replace the old token using an alternation group.

Applying a filter can be done manually or using the `ApplyFilters` function exported by the [github.com/zimmski/tavor/fuzz/filter](/fuzz/filter) package. It can apply more than one filter, does correctly traverse the graph, handle errors of filters and does not apply filters onto filter generated tokens. The last property is needed to avoid filter loops.

The `Register` function of the [github.com/zimmski/tavor/fuzz/filter](/fuzz/filter) package allows to register filters based on an identifier which can be then used within the framework. The function `New` of the [github.com/zimmski/tavor/fuzz/filter](/fuzz/filter) package allows to generate a new instance of the registered filter given the identifier. This is for example needed for the Tavor binary, which applies filters defined by CLI arguments.

**Examples**

The following fuzzing filter searches the token graph for constant string tokens which hold the string "old" and replaced them with a constant string token holding the string "new".

```go
import (
	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/primitives"
)

type SampleFilter struct{}

func NewSampleFilter() *SampleFilter {
	return &SampleFilter{}
}

func (f *SampleFilter) Apply(tok token.Token) ([]token.Token, error) {
	c, ok := tok.(*primitives.ConstantString)
	if !ok || c.String() != "old" {
		return nil, nil
	}

	return []token.Token{
		primitives.NewConstantString("new"),
	}, nil
}
```
One way to apply this filter is by using the following code. Which will change the generation from "old string" to "new string" after the filter is applied.

```go
import (
	"fmt"

	"github.com/zimmski/tavor/fuzz/filter"
	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/lists"
	"github.com/zimmski/tavor/token/primitives"
)

func main() {
	var doc token.Token = lists.NewAll(
		primitives.NewConstantString("old"),
		primitives.NewConstantString(" "),
		primitives.NewConstantString("string"),
	)

	var filters = []filter.Filter{
		NewSampleFilter(),
	}

	doc, err := filter.ApplyFilters(filters, doc)
	if err != nil {
		panic(err)
	}

	fmt.Println(doc.String())
}
```

The filter can be registered as a framework-wide usable filter using the following code. Please note that this should be usually done in an `init` function inside the package of a filter.

```go
import (
	"github.com/zimmski/tavor/fuzz/filter"
)

func init() {
	filter.Register("SampleFilter", func() filter.Filter {
		return NewSampleFilter()
	})
}
```

### Fuzzing strategies [![GoDoc](https://godoc.org/github.com/zimmski/tavor?status.png)](https://godoc.org/github.com/zimmski/tavor/fuzz/strategy)

The fuzzing strategy code and all officially implemented fuzzing strategies can be found in the [github.com/zimmski/tavor/fuzz/strategy](/fuzz/strategy) package and its sub-packages.

Each fuzzing strategy instance has to be associated on construction with exactly one token. This allows an instance to hold a dedicated state of the given token's graph, which makes optimizations for multiple fuzzing operations possible.

A fuzzing strategy has to implement the `Strategy` interface which is exported by the [github.com/zimmski/tavor/fuzz/strategy](/fuzz/strategy) package. The interface defines the `Fuzz` method which starts the first iteration of the fuzzing strategy in a new goroutine and returns a channel which controls the fuzzing process. The error return argument is not nil, if an error is encountered during the initialization. On success a value is returned by the channel which marks the completion of the iteration. A value has to be put back in, to initiate the calculation of the next fuzzing iteration. This passing of values is needed to avoid data races within the token graph. The channel must be closed when there are no more iterations or the strategy caller wants to end the fuzzing process. Please note that this can also occur right after receiving the channel. Hence when there are no iterations. Since the `Fuzz` method is running in its own goroutine, it can be implemented statefully without using savepoints.

The `Register` function of the [github.com/zimmski/tavor/fuzz/strategy](/fuzz/strategy) package allows to register strategies based on an identifier which can be then used within the framework. The function `New` of the [github.com/zimmski/tavor/fuzz/strategy](/fuzz/strategy) package allows to generate a new instance of the registered strategy given the identifier. This is for example needed for the Tavor binary, which can execute a specific strategy defined by a CLI argument.

**Examples**

The following fuzzing strategy searches the token graph for constant integer tokens which have a value within 1 and 10 and increments their content by replacing the original value. This strategy falls therefore in the category of mutation-based fuzzing, since it does change the original data. It is also stateless since there is no need to keep track of current events between iterations. The graph is simply searched and changed once per iteration. An additional property is that the defined operation allows the strategy to end, which is strictly not necessary.

Please note that the search of tokens could be cached, which would give the fuzzing strategy instance a kind of state. It would strictly speaking still not be a stateful fuzzing strategy, since there is no connection between iterations. That is no iteration depends on a previous iteration. Caching tokens would also imply that changes in the graph, which could be made outside the fuzzing go routing, must be handled or a contract has to be made between caller and callee of the fuzzing strategy.

```go
import (
	"github.com/zimmski/tavor/rand"
	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/primitives"
)

type SampleStrategy struct {
	root token.Token
}

func NewSampleStrategy(tok token.Token) *SampleStrategy {
	return &SampleStrategy{
		root: tok,
	}
}

func (s *SampleStrategy) Fuzz(r rand.Rand) (chan struct{}, error) {
	continueFuzzing := make(chan struct{})

	go func() {
		for {
			found := false

			err := token.Walk(s.root, func(tok token.Token) error {
				intTok, ok := tok.(*primitives.ConstantInt)
				if !ok {
					return nil
				}

				v := intTok.Value()
				if v >= 1 && v <= 10 {
					found = true

					v++

					intTok.SetValue(v)
				}

				return nil
			})
			if err != nil {
				panic(err)
			}

			if !found {
				break
			}

			continueFuzzing <- struct{}{}

			if _, ok := <-continueFuzzing; !ok {
				// the fuzzing strategy caller can close the channel too

				return
			}
		}

		close(continueFuzzing)
	}()

	return continueFuzzing, nil
}
```

One way to execute this strategy is by using the following code. Please note that the implemented strategy does not need a random generator. The argument for the `Fuzz` method can be therefore just `nil`.

```go
import (
	"fmt"

	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/lists"
	"github.com/zimmski/tavor/token/primitives"
)

func main() {
	var doc token.Token = lists.NewAll(
		primitives.NewConstantInt(7),
		primitives.NewConstantString(" "),
		primitives.NewConstantInt(9),
	)

	strat := NewSampleStrategy(doc)

	continueFuzzing, err := strat.Fuzz(nil)
	if err != nil {
		panic(err)
	}

	for i := range continueFuzzing {
		fmt.Println(doc.String())

		continueFuzzing <- i
	}
}
```

This program results in the following output.

```
8 10
9 11
10 11
11 11
```

The strategy can be registered as a framework-wide usable strategy using the following code. Please note that this should be usually done in an `init` function inside the package of a strategy.

```go
import (
	"github.com/zimmski/tavor/fuzz/strategy"
	"github.com/zimmski/tavor/token"
)

func init() {
	strategy.Register("SampleStrategy", func(tok token.Token) strategy.Strategy {
		return NewSampleStrategy(tok)
	})
}
```

### Reduce strategies [![GoDoc](https://godoc.org/github.com/zimmski/tavor?status.png)](https://godoc.org/github.com/zimmski/tavor/reduce/strategy)

The reduce strategy code and all officially implemented reduce strategies can be found in the [github.com/zimmski/tavor/reduce/strategy](/reduce/strategy) package and its sub-packages.

Each reduce strategy instance has to be associated on construction with exactly one token. This allows an instance to hold a dedicated state of the given token's graph, which makes optimizations for multiple reduce operations possible.

A reduce strategy has to implement the `Strategy` interface which is exported by the [github.com/zimmski/tavor/reduce/strategy](/reduce/strategy) package. The interface defines the `Reduce` method which starts the first step of the reduce strategy in a new goroutine and returns two channels to control the reduce process. The error return argument is not nil, if an error is encountered during the initialization. On success a value is returned by the control channel which marks the completion of the iteration. A feedback has to be given through the feedback channel as well as a value to the control channel to initiate the calculation of the next reduce step. This passing of values is needed to avoid data races within the token graph. The channels must be closed when there are no more steps or the strategy caller wants to end the reduce process. Please note that this can also occur right after receiving the channels. Hence when there are no steps. Since the `Reduce` method is running in its own goroutine, it can be implemented statefully without using savepoints.

Currently only two different feedback answers can be given. They are defined by the `ReduceFeedbackType` type which is exported by the [github.com/zimmski/tavor/reduce/strategy](/reduce/strategy) package. One feedback answer is `Good` which communicates to the reduce strategy that the current step produced a successful result. This can mean for example that the result has the right syntax or is better than the last good result. The meaning of the feedback and the response of the strategy to the feedback are purely dependent on the application. Responses could be for example to proceed with a given optimization path or to simply end the whole reducing process, since it is often enough to find one solution. The second feedback answer is `Bad` which communicates exactly the opposite of `Good` to the strategy. This answer is often more complicated to handle since it means that in some scenarios a revert of the current step to the last good step has to occur before the reduce process can continue.

Please note that all reduce strategies should currently implement algorithms which produce valid generations according to the internal token graph. Hence for example a constant integer should not be replaced by a constant string. This is a convention which is not enforced but highly recommended to avoid problems until it is safely supported by a future version of Tavor.

The `Register` function of the [github.com/zimmski/tavor/reduce/strategy](/reduce/strategy) package allows to register strategies based on an identifier which can be then used within the framework. The function `New` of the [github.com/zimmski/tavor/reduce/strategy](/reduce/strategy) package allows to generate a new instance of the registered strategy given the identifier. This is for example needed for the Tavor binary, which can execute a specific strategy defined by a CLI argument.

**Examples**

The following reduce strategy searches the token graph for repeat tokens holding one constant string token as its internal token and reduces the repetition by one token for every reduce step. This could be done more efficiently but it is very simple and should demonstrate that reduce strategies can be implemented very easily to start with. It is a stateful strategy since every step depends on the previous one. Additionally it is guaranteed to end, since only a finite amount of tokens is targeted without generating new ones.

```go
import (
	"errors"

	"github.com/zimmski/tavor/reduce/strategy"
	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/lists"
	"github.com/zimmski/tavor/token/primitives"
)

type SampleStrategy struct {
	root token.Token
}

func NewSampleStrategy(tok token.Token) *SampleStrategy {
	return &SampleStrategy{
		root: tok,
	}
}

func (s *SampleStrategy) Reduce() (chan struct{}, chan<- strategy.ReduceFeedbackType, error) {
	continueReducing := make(chan struct{})
	feedbackReducing := make(chan strategy.ReduceFeedbackType)

	go func() {
		done := errors.New("done")

		err := token.Walk(s.root, func(tok token.Token) error {
			repeat, ok := tok.(*lists.Repeat)
			if !ok || repeat.InternalLen() != 1 {
				return nil
			}

			char, err := repeat.InternalGet(0)
			if err != nil {
				return err
			}
			if _, ok := char.(*primitives.ConstantString); !ok {
				return nil
			}

			for i := repeat.Reduces(); i >= 1; {
				found := false
				l := len(repeat.String())
				for ; i >= 1; i-- {
					if err := repeat.Reduce(i); err != nil {
						return err
					}

					if l-1 == len(repeat.String()) {
						found = true

						break
					}
				}

				if !found {
					break
				}

				continueReducing <- struct{}{}

				feedback, ok := <-feedbackReducing
				if !ok {
					return nil
				}

				if _, ok := <-continueReducing; !ok {
					return nil
				}

				if feedback == strategy.Good {
					return done
				}
			}

			return nil
		})
		if err != nil && err != done {
			panic(err)
		}

		close(continueReducing)
		close(feedbackReducing)
	}()

	return continueReducing, feedbackReducing, nil
}
```

One way to execute this strategy is by using the following code.

```go
import (
	"fmt"

	"github.com/zimmski/tavor/reduce/strategy"
	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/lists"
	"github.com/zimmski/tavor/token/primitives"
)

func main() {
	aRepeat := lists.NewRepeat(primitives.NewConstantString("a"), 0, 100)
	aRepeat.Permutation(7)
	bRepeat := lists.NewRepeat(primitives.NewConstantString("b"), 1, 100)
	bRepeat.Permutation(5)
	cRepeat := lists.NewRepeat(primitives.NewConstantString("c"), 7, 100)
	cRepeat.Permutation(9)
	dRepeat := lists.NewRepeat(primitives.NewConstantString("d"), 1, 100)
	dRepeat.Permutation(2)

	var doc token.Token = lists.NewAll(
		aRepeat,
		bRepeat,
		cRepeat,
		dRepeat,
	)

	fmt.Println(doc.String())

	strat := NewSampleStrategy(doc)

	continueFuzzing, feedbackReducing, err := strat.Reduce()
	if err != nil {
		panic(err)
	}

	for i := range continueFuzzing {
		out := doc.String()

		fmt.Println(out)

		if len(out) <= 10 {
			feedbackReducing <- strategy.Good
		} else {
			feedbackReducing <- strategy.Bad
		}

		continueFuzzing <- i
	}
}
```

This program results in the following output.

```
aaaaaabbbbbcccccccccccccccdd
aaaaabbbbbcccccccccccccccdd
aaaabbbbbcccccccccccccccdd
aaabbbbbcccccccccccccccdd
aabbbbbcccccccccccccccdd
abbbbbcccccccccccccccdd
bbbbbcccccccccccccccdd
bbbbcccccccccccccccdd
bbbcccccccccccccccdd
bbcccccccccccccccdd
bcccccccccccccccdd
bccccccccccccccdd
bcccccccccccccdd
bccccccccccccdd
bcccccccccccdd
bccccccccccdd
bcccccccccdd
bccccccccdd
bcccccccdd
```

The strategy can be registered as a framework-wide usable strategy using the following code. Please note that this should be usually done in an `init` function inside the package of a strategy.

```go
import (
	"github.com/zimmski/tavor/reduce/strategy"
	"github.com/zimmski/tavor/token"
)

func init() {
	strategy.Register("SampleStrategy", func(tok token.Token) strategy.Strategy {
		return NewSampleStrategy(tok)
	})
}
```

### Tokens [![GoDoc](https://godoc.org/github.com/zimmski/tavor?status.png)](https://godoc.org/github.com/zimmski/tavor/tokens)

TODO<br/>
TODO explain the different interfaces for tokens<br/>

### Attributes for tokens

TODO<br/>

### Special tokens

TODO<br/>

## <a name="stability"></a>How stable is Tavor?

Tavor is still in development and fare from a 1.0 release. There are [some bugs](https://github.com/zimmski/tavor/issues?q=is%3Aopen+is%3Aissue+label%3Abug) and a lot of [functionality is still missing](https://github.com/zimmski/tavor/issues?q=is%3Aopen+is%3Aissue+label%3Aenhancement) but basic features are stable enough and are successufully used in production by many projects.

[Individual package code coverage](https://coveralls.io/r/zimmski/tavor) is currently low but since most tests do cover a lot of Tavor's components this is not a big issue. However 100% coverage using hand written tests is a required feature of the 1.0 release as well as fully fuzzing the Tavor format and the Tavor binary. This means that Tavor will be equipped to test every feature of the binary and the Tavor format itself.

Since Tavor is still a moving target, backwards-incompatible changes will happen but are documented for every release. This is necessary to make working with Tavor as easy as possible while still providing loads of functionality.

## <a name="missing-features"></a>Missing features

- Format: Format files for binary data and different character sets (currently only UTF-8 is supported)
- General: Direct support for protocols (can be currently only done with fuzzing an output and putting this input into an executor)
- General: Direct support for source code generation and execution (needs an execution layer as-well)
- Format: Functions with parameters to reduce clutter
- General: Remove the need for unrolling and allow real loops
- Format: Includes of external format files
- Fuzzing: Feedback-driven fuzzing -> transition into completely stateful fuzzing
- General: Parallel execution of fuzzing, delta-debugging, ...
- Binary: Online fuzzing
- Fuzzing: Mutation based fuzzing
- General: Encoding/Decoding of data e.g. to encrypt parts of data

There are also a lot of smaller features and enhancements waiting in the [issue tracker](https://github.com/zimmski/tavor/issues).

## <a name="feature-request"></a>Can I make feature requests and report bugs and problems?

Sure, just submit an [issue via the project tracker](https://github.com/zimmski/tavor/issues/new) and I will see what I can do. Please note that I do not guarantee to implement anything soon and bugs and problems are more important to me than new features. If you need something implemented or fixed right away you can contact me via mail <mz@nethead.at> to do contract work for you.
