# Tavor [![GoDoc](https://godoc.org/github.com/zimmski/tavor?status.png)](https://godoc.org/github.com/zimmski/tavor) [![Build Status](https://travis-ci.org/zimmski/tavor.svg?branch=master)](https://travis-ci.org/zimmski/tavor)

Tavor ([Sindarin](https://en.wikipedia.org/wiki/Sindarin) for "woodpecker") is a fuzzing and delta-debugging platform.

## How do I use Tavor?<a name="use"/>

### The Tavor binary

The Tavor binary provides fuzzing and delta-debugging functionality for Tavor format files as well as some other commands. Sane default arguments should provide a pleasant experience.

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

## How do I build Tavor?<a name="build"/>

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

## Where are the precompiled binaries?<a name="precompiled"/>

You can find all precompiled binaries on the [release page](https://github.com/zimmski/tavor/releases). The binaries are packed into archives that currently only hold the Tavor binary itself.

### Bash Completion

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

## Can I make feature requests, report bugs and problems?<a name="feature-request"/>

Sure, just submit an [issue via the project tracker](https://github.com/zimmski/tavor/issues/new) and I will see what I can do. Please note that I do not guarantee to implement anything soon and bugs and problems are more important to me than new features. If you need something implemented or fixed right away you can contact me via mail <mz@nethead.at> to do contract work for you.
