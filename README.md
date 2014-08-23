# Tavor [![GoDoc](https://godoc.org/github.com/zimmski/tavor?status.png)](https://godoc.org/github.com/zimmski/tavor) [![Build Status](https://travis-ci.org/zimmski/tavor.svg?branch=master)](https://travis-ci.org/zimmski/tavor) [![Coverage Status](https://coveralls.io/repos/zimmski/tavor/badge.png)](https://coveralls.io/r/zimmski/tavor)

Tavor ([Sindarin](https://en.wikipedia.org/wiki/Sindarin) for "woodpecker") is a fuzzing and delta-debugging platform.

TODO SHORT description on what you can do with the platform and what is the purpose of it<br/>

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

You can download this file called [<code>basic.tavor</code> from here](examples/quick/basic.tavor).

Now we can use Tavor to [fuzz](#fuzzing) the format by issuing the following command:

```bash
tavor --format-file basic.tavor fuzz
```

On every call this command outputs random paths through the defined graph, since the default [fuzzing strategy](#fuzzing-strategy) of Tavor is the <code>random</code> strategy.

Here are some example outputs:

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

Generating data like this is just one example of the capabilities of Tavor. More interesting than generating data, is what you can do with it. Forwarding the data to a program to test the given vending machine is another possible use case of Tavor.

Please have a look [here](#bigexample) if you like to see a bigger example with a complete overview over the basic features or keep reading to find out more about the background and capabilities of Tavor.

## <a name="fuzzing"></a>What is fuzzing?

TODO in general, which types of fuzzing are there, what you can do with them, what are the pros and cons<br/>
TODO mention that it is pretty much just data generation and can be used for example for genetic programming<br/>

## <a name="delta-debugging"></a>What is delta-debugging?

TODO in general, which types of delta-debugging are there, what you can do with them, what are the pros and cons<br/>
TODO mention that delta-debugging and reducing are synonyms<br/>

## How does Tavor work and what does it provide?

TODO model-based concept with format files, doing almost everything with format files should be possible<br/>
TODO how fuzzing works in general with the model-based concept and the unrolling<br/>
TODO how delta-debugging works in general with the model-based concept, reading the input, delta-debug on it<br/>
TODO mention missing features -> link to it<br/>
TODO mention that it is a platform to extend on, so researchers and testers do not have to implement everything from scratch<br/>

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

### <a name="delta-debugging-strategy"></a>What are delta-debugging strategies?

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
- Format: Includes of external format files
- Fuzzing: Feedback-driven fuzzing -> transition into completely stateful fuzzing
- General: Parallel execution of fuzzing, delta-debugging, ...
- Binary: Online fuzzing
- Fuzzing: Mutation based fuzzing
- General: Encoding/Decoding of data e.g. to encrypt parts of data

## <a name="feature-request"></a>Can I make feature requests, report bugs and problems?

Sure, just submit an [issue via the project tracker](https://github.com/zimmski/tavor/issues/new) and I will see what I can do. Please note that I do not guarantee to implement anything soon and bugs and problems are more important to me than new features. If you need something implemented or fixed right away you can contact me via mail <mz@nethead.at> to do contract work for you.
