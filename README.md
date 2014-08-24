# Tavor [![GoDoc](https://godoc.org/github.com/zimmski/tavor?status.png)](https://godoc.org/github.com/zimmski/tavor) [![Build Status](https://travis-ci.org/zimmski/tavor.svg?branch=master)](https://travis-ci.org/zimmski/tavor) [![Coverage Status](https://coveralls.io/repos/zimmski/tavor/badge.png)](https://coveralls.io/r/zimmski/tavor)

Tavor ([Sindarin](https://en.wikipedia.org/wiki/Sindarin) for "woodpecker") is a fuzzing and delta-debugging platform, which provides a framework and binary to not only implement and do everyday fuzzing and delta-debugging but to also do research on new methods without implementing basics again. A EBNF-like notations allows the definition of data (e.g. file formats and protocols) without the need of programming. Tavor also relaxes on the definitions of fuzzing and delta-debugging allowing the user to utilize implemented techniques universally e.g. for key-driven testing, model-based testing, simulating user-behavior and genetic programming.

### <a name="quick-example"></a>A quick example

Imagine a vending machine which ejects a product after receiving 100 worth of credits. It is possible to input 25 and 50 credit coins into the machine. After receiving enough credits the machine ejects a product and resets the credit counter to zero. To keep it simple, we specify that the machine does not handle credit overflows. A representation of the states and actions of the machine could look like this:

![Basic states and actions](examples/quick/basic.png "Basic states and actions")

This state machine can be defined using the following [Tavor format file](#format):

```
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

You can download this file called [<code>basic.tavor</code> from here](examples/quick/basic.tavor).

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

Generating data like this is just one example of the capabilities of Tavor. Please have a look [here](#bigexample) if you like to see a bigger example with a complete overview over the basic features or keep reading to find out more about the background and capabilities of Tavor.

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

	There are no common techniques for generation-based fuzzing but most algorithms choose a graph as underlying representation of the data model. The graph is then traversed and each node outputs a part of the data. The traversal algorithms and the complexity and abilities of the data model like constraints between nodes or adding nodes during the traversal distinguish generation-based fuzzers and contribute in general to their mightiness.

## <a name="delta-debugging"></a>What is delta-debugging?

> The Delta Debugging algorithm isolates failure causes automatically - by systematically narrowing down failure-inducing circumstances until a minimal set remains.
> -- <cite>[https://en.wikipedia.org/wiki/Delta_Debugging](https://en.wikipedia.org/wiki/Delta_Debugging)</cite>

E.g. we feed a given data to a program which fails on executing. By delta-debugging this data we can reduce it to its minimum while still failing the execution. The reduction of the data is handled by software heuristics (semi-)automatically. The obvious advantage of this method, besides being done (semi-)automatically, is that we do not need to handle uninteresting parts of the data while debugging the problem, we can focus on the important parts which actually lead to the failure.

**Note**: Since delta-debugging reduces data it is also called <code>reducing</code>.

Delta-debugging consists of three areas:
- A heuristic has to decide which parts of the data will be reduced next
- The reduction itself e.g.
	- Remove repetitions
	- Remove optional data
	- Replace data with something else e.g. replace an uninteresting complex function with a constant value
- Testing the new resulting data concerning the failure

Although delta-debugging is described as method to isolate failure causes, it can be also used to isolate anything given isolating constraints. For example we could reduce an input for a program which leads to a positive outcome to its minimum.

## How does Tavor work and what does it provide?

Tavor combines both fuzzing and delta-debugging into one platform by allowing all implemented methods to operate on one internal model-based structure represented by a graph. This structure can be defined and generated programmatically or by using a format file. Out of the box Tavor comes with its own [format](#format) which covers all functionality of the Tavor framework itself.

Tavor's fuzzing implementation is generically and not fixed to one technique nor format. Instead different fuzzing techniques and heuristics can be implemented and executed independently as [Tavor fuzzing strategies](fuzzing-strategy). The same principle is used for delta-debugging where so called [Tavor reduce strategies](#reduce-strategy) can be implemented and used. Both types of strategies operate on the same internal structure independent of the format.

Even tough Tavor provides loads of functionality out of the box a lot is still missing. A list of missing but planed features can be found in the [missing features section](#missing-features). For feature request please have a look at the [feature request section](#feature-request).

### Unrolling loops

Although the internal structure allows loops in its graph, Tavor currently unrolles loops for easier usage. This is a trade-off that is currently in place but will be eliminated in future versions of Tavor.

E.g. This graph loops between the states <code>Idle</code> and <code>Action</code>:

![Looping](/doc/images/README/unroll-loop.png "Looping")

Will result in the following internal graph given a maximum of two repetitions:

![Unrolled](/doc/images/README/unroll-unrolled.png "Unrolled")

## <a name="format"></a>The Tavor format file

TODO -> put this in its own .md and do not skimp on examples<br/>
TODO explain every aspect. basics first<br/>

## <a name="use"></a>How do I use Tavor?

TODO explain the currently three ways to use Tavor: binary, programmatically using Tavor and extending Tavor itself<br/>

TODO explain that all three ways work in the same way filters fuzzing and delta debugging<br/>
TODO mention bigger example scenario -> link to it<br/>

### <a name="fuzzing-filter"></a>What are fuzzing filters?

TODO<br/>

TODO available filters -> link to godoc and explain the filters in the code<br/>

### <a name="fuzzing-strategy"></a>What are fuzzing strategies?

TODO<br/>

TODO available strategies -> link to godoc and explain the strategies in the code<br/>

### <a name="reduce-strategy"></a>What are reduce strategies?

TODO<br/>

TODO available strategies -> link to godoc and explain the strategies in the code<br/>

## <a name="binary"></a>The Tavor binary

The [Tavor binary](#precompiled) provides fuzzing and delta-debugging functionality for Tavor format files as well as some other commands. Sane default arguments should provide a pleasant experience.

Since the binary acts on Tavor format files, the <code>--format-file</code> argument has to be used for every non-informational action. E.g. the following commands fuzzes the given format file with the default fuzzing strategy:

```bash
tavor --format-file file.tavor fuzz
```

In contrast listing all available fuzzing strategies does not require the <code>--format-file</code> argument:

```bash
tavor fuzz --list-strategies
```

To learn more about available arguments and commands, you can invoke the binary's help by executing the binary without any arguments or with the <code>--help</code> argument. Here is a complete overview of all arguments, commands and their options.

```
Usage:
  tavor <options> <command> <command options>

General options:
  --debug             Debug log output
  --help              Show this help message
  --verbose           Verbose log output
  --version           Print the version of this program

Global options:
  --seed=             Seed for all the randomness
  --max-repeat=       How many times loops and repetitions should be repeated (2)

Format file options:
  --format-file=      Input tavor format file
  --print             Prints the AST of the parsed format file
  --print-internal    Prints the internal AST of the parsed format file
  --validate          Just validate the format file and exit

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
      --exec-argument-type=                      How the generation is given to the binary (environment)
      --list-exec-argument-types                 List all available exec argument types
      --script=                                  Execute this binary which gets fed with the generation and should return
                                                 feedback
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
      --exec-argument-type=             How the generation is given to the binary (environment)
      --list-exec-argument-types        List all available exec argument types
      --script=                         Execute this binary which gets fed with the generation and should return feedback
      --input-file=                     Input file which gets parsed, validated and delta-debugged via the format file
      --strategy=                       The reducing strategy (BinarySearch)
      --list-strategies                 List all available reducing strategies
      --result-separator=               Separates result outputs of each reducing step ("\n")

[validate command options]
      --input-file=   Input file which gets parsed and validated via the format file
```

### Graphing

TODO with examples<br/>

 | dot -Tsvg -o outfile.svg

### Fuzzing

TODO bigger example with example commands and files<br/>

### Delta-debugging

TODO bigger example with example commands and files<br/>

## <a name="bigexample"></a>A complete example for fuzzing, executing and delta-debugging

TODO this example should give a complete overview of how Tavor can be used.<br/>
TODO do a key-word driven format-file<br/>
TODO executor for the key-words<br/>
TODO delta-debug keywords because of an intentional error<br/>

## <a name="programmatically"></a>How do I use the Tavor platform programmatically?

TODO<br/>
TODO explain creating internal structures (instead of using a format file) with examples<br/>
TODO explain how to use filters, fuzzers and delta debugging<br/>

## <a name="extend"></a>How do I extend Tavor?

TODO<br/>
TODO mention feature request section, but if someone is interested in really extending Tavor by her/himself... read on<br/>

### Filters

TODO<br/>

### Fuzzing strategies

TODO<br/>

### Delta-debugging strategies

TODO<br/>

### Tokens

TODO<br/>
TODO explain the different interfaces for tokens<br/>

### Attributes for tokens

TODO<br/>

### Special tokens

TODO<br/>

### Still looking for something else?

TODO explain if the reader has not find what she/he looks for -> link to the feature request section<br/>

## <a name="build"></a>How do I build Tavor?

Tavor provides [precompiled 64 bit Linux binaries](#precompiled). Other platforms are currently not supported, but might work. Please have a look at the [feature request section](#feature-request) if you need them to work or you want more binaries.

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

You now have a binary "tavor" in your GOPATH/bin (or if set GOBIN folder) folder which can be used without any further actions.

## <a name="precompiled"></a>Where are the precompiled binaries?

You can find all precompiled binaries on the [release page](https://github.com/zimmski/tavor/releases). The binaries are packed into archives that currently only hold the Tavor binary itself.

### <a name="bash-completion"></a>Bash Completion

If you like Bash Completion for Tavor make sure that you have Bash Completion installed and then copy the [bash completion Tavor script](https://raw.githubusercontent.com/zimmski/tavor/master/bin/tavor-bash_completion.sh) into your Bash Completion folder.

```bash
mkdir -p $HOME/.bash_completion
wget -P $HOME/.bash_completion https://raw.githubusercontent.com/zimmski/tavor/master/bin/tavor-bash_completion.sh
. ~/.bashrc
```

Bash Completion for Tavor should now be working. If not, one reason could be that your distribution does not include user defined Bash Completion scripts in .bashrc so just add it to your .bashrc:

```bash
echo ". ~/.bash_completion/tavor-bash_completion.sh" >> ~/.bashrc
. ~/.bashrc
```

## <a name="missing-features"></a>Missing features

- Format: Format files for binary data and different character sets (currently only UTF-8 is supported)
- General: Direct support for protocols (can be currently only done with fuzzing an output and putting this input into an executor)
- Format: Functions with parameters to reduce clutter
- General: Remove the need for unrolling and allow real loops
- Format: Includes of external format files
- Fuzzing: Feedback-driven fuzzing -> transition into completely stateful fuzzing
- General: Parallel execution of fuzzing, delta-debugging, ...
- Binary: Online fuzzing
- Fuzzing: Mutation based fuzzing
- General: Encoding/Decoding of data e.g. to encrypt parts of data

## <a name="feature-request"></a>Can I make feature requests, report bugs and problems?

Sure, just submit an [issue via the project tracker](https://github.com/zimmski/tavor/issues/new) and I will see what I can do. Please note that I do not guarantee to implement anything soon and bugs and problems are more important to me than new features. If you need something implemented or fixed right away you can contact me via mail <mz@nethead.at> to do contract work for you.
