package main

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/zimmski/osutil"

	"github.com/zimmski/tavor"
	tavorFuzzFilter "github.com/zimmski/tavor/fuzz/filter"
	tavorFuzzStrategy "github.com/zimmski/tavor/fuzz/strategy"
	"github.com/zimmski/tavor/graph"
	"github.com/zimmski/tavor/log"
	"github.com/zimmski/tavor/parser"
	tavorReduceStrategy "github.com/zimmski/tavor/reduce/strategy"
	"github.com/zimmski/tavor/token"
)

type exitCodeType int

const (
	exitCodeOk exitCodeType = iota
	exitCodeHelp
	exitCodeBashCompletion
	exitCodeInvalidInputFile
	exitCodeError
)

type options struct {
	General struct {
		Debug   bool `long:"debug" description:"Debug log output"`
		Help    bool `long:"help" description:"Show this help message"`
		Verbose bool `long:"verbose" description:"Verbose log output"`
		Version bool `long:"version" description:"Print the version of this program"`
	} `group:"General options"`

	Global struct {
		Seed      int64 `long:"seed" description:"Seed for all the randomness"`
		MaxRepeat int   `long:"max-repeat" description:"How many times loops and repetitions should be repeated" default:"2"`
	} `group:"Global options"`

	Format struct {
		Check         bool           `long:"check" description:"Just check the syntax of the format file and exit"`
		FormatFile    flags.Filename `long:"format-file" description:"Input tavor format file" required:"true"`
		Print         bool           `long:"print" description:"Prints the AST of the parsed format file"`
		PrintInternal bool           `long:"print-internal" description:"Prints the internal AST of the parsed format file"`
	} `group:"Format file options"`

	Fuzz struct {
		Exec struct {
			Exec                           string           `long:"exec" description:"Execute this binary with possible arguments to test a generation"`
			ExecExactExitCode              int              `long:"exec-exact-exit-code" description:"Same exit code has to be present" default:"-1"`
			ExecExactStderr                string           `long:"exec-exact-stderr" description:"Same stderr output has to be present"`
			ExecExactStdout                string           `long:"exec-exact-stdout" description:"Same stdout output has to be present"`
			ExecMatchStderr                string           `long:"exec-match-stderr" description:"Searches through stderr via the given regex. A match has to be present"`
			ExecMatchStdout                string           `long:"exec-match-stdout" description:"Searches through stdout via the given regex. A match has to be present"`
			ExecDoNotRemoveTmpFiles        bool             `long:"exec-do-not-remove-tmp-files" description:"If set, tmp files are not removed"`
			ExecDoNotRemoveTmpFilesOnError bool             `long:"exec-do-not-remove-tmp-files-on-error" description:"If set, tmp files are not removed on error"`
			ExecArgumentType               execArgumentType `long:"exec-argument-type" description:"How the generation is given to the binary" default:"stdin"`
			ListExecArgumentTypes          bool             `long:"list-exec-argument-types" description:"List all available exec argument types"`

			Script string `long:"script" description:"Execute this binary which gets fed with the generation and should return feedback"`

			ExitOnError bool `long:"exit-on-error" description:"Exit if an execution fails"`
		}

		Filter optsFuzzingFilters

		Strategy       fuzzStrategy `long:"strategy" description:"The fuzzing strategy" default:"random"`
		ListStrategies bool         `long:"list-strategies" description:"List all available fuzzing strategies"`

		ResultFolder     flags.Filename `long:"result-folder" description:"Save every fuzzing result with the MD5 checksum as filename in this folder"`
		ResultExtensions string         `long:"result-extension" description:"If result-folder is used this will be the extension of every filename"`
		ResultSeparator  string         `long:"result-separator" description:"Separates result outputs of each fuzzing step" default:"\n"`
	} `command:"fuzz" description:"Fuzz the given format file"`

	Graph struct {
		Filter optsFuzzingFilters
	} `command:"graph" description:"Generate a DOT file out of the internal AST"`

	Reduce struct {
		Exec struct {
			Exec                    string           `long:"exec" description:"Execute this binary with possible arguments to test a generation"`
			ExecExactExitCode       bool             `long:"exec-exact-exit-code" description:"Same exit code has to be present"`
			ExecExactStderr         bool             `long:"exec-exact-stderr" description:"Same stderr output has to be present"`
			ExecExactStdout         bool             `long:"exec-exact-stdout" description:"Same stdout output has to be present"`
			ExecMatchStderr         string           `long:"exec-match-stderr" description:"Searches through stderr via the given regex. A match has to be present"`
			ExecMatchStdout         string           `long:"exec-match-stdout" description:"Searches through stdout via the given regex. A match has to be present"`
			ExecDoNotRemoveTmpFiles bool             `long:"exec-do-not-remove-tmp-files" description:"If set, tmp files are not removed"`
			ExecArgumentType        execArgumentType `long:"exec-argument-type" description:"How the generation is given to the binary" default:"stdin"`
			ListExecArgumentTypes   bool             `long:"list-exec-argument-types" description:"List all available exec argument types"`

			Script string `long:"script" description:"Execute this binary which gets fed with the generation and should return feedback"`
		}

		InputFile flags.Filename `long:"input-file" description:"Input file which gets parsed, validated and delta-debugged via the format file" required:"true"`

		Strategy       reduceStrategy `long:"strategy" description:"The reducing strategy" default:"Linear"`
		ListStrategies bool           `long:"list-strategies" description:"List all available reducing strategies"`

		ResultSeparator string `long:"result-separator" description:"Separates result outputs of each reducing step" default:"\n"`
	} `command:"reduce" description:"Reduce the given input file"`

	Validate struct {
		InputFile flags.Filename `long:"input-file" description:"Input file which gets parsed and validated via the format file" required:"true"`
	} `command:"validate" description:"Validate the given input file"`
}

type optsFuzzingFilters struct {
	Filters     fuzzFilters `long:"filter" description:"Fuzzing filter to apply"`
	ListFilters bool        `long:"list-filters" description:"List all available fuzzing filters"`
}

var execArgumentTypes = []string{
	"argument",
	"environment",
	"stdin",
}

type execArgumentType string

func (e execArgumentType) Complete(match string) []flags.Completion {
	var items []flags.Completion

	for _, name := range execArgumentTypes {
		if strings.HasPrefix(name, match) {
			items = append(items, flags.Completion{
				Item: name,
			})
		}
	}

	return items
}

type fuzzFilter string
type fuzzFilters []fuzzFilter

func (s fuzzFilters) Complete(match string) []flags.Completion {
	var items []flags.Completion

	for _, name := range tavorFuzzFilter.List() {
		if strings.HasPrefix(name, match) {
			items = append(items, flags.Completion{
				Item: name,
			})
		}
	}

	return items
}

type fuzzStrategy string

func (s *fuzzStrategy) Complete(match string) []flags.Completion {
	var items []flags.Completion

	for _, name := range tavorFuzzStrategy.List() {
		if strings.HasPrefix(name, match) {
			items = append(items, flags.Completion{
				Item: name,
			})
		}
	}

	return items
}

type reduceStrategy string

func (s *reduceStrategy) Complete(match string) []flags.Completion {
	var items []flags.Completion

	for _, name := range tavorReduceStrategy.List() {
		if strings.HasPrefix(name, match) {
			items = append(items, flags.Completion{
				Item: name,
			})
		}
	}

	return items
}

func checkArguments(args []string, opts *options) (string, exitCodeType) {
	p := flags.NewNamedParser("tavor", flags.None)

	p.ShortDescription = "A generic fuzzing and delta-debugging framework."

	if _, err := p.AddGroup("Tavor", "Tavor arguments", opts); err != nil {
		return "", exitError(err.Error())
	}

	completion := len(os.Getenv("GO_FLAGS_COMPLETION")) > 0

	_, err := p.ParseArgs(args)
	if (opts.General.Help || len(args) == 0) && !completion {
		p.WriteHelp(os.Stdout)

		return "", exitCodeHelp
	} else if opts.General.Version {
		fmt.Printf("Tavor v%s\n", tavor.Version)

		return "", exitCodeOk
	} else if opts.Fuzz.Filter.ListFilters || opts.Graph.Filter.ListFilters {
		for _, name := range tavorFuzzFilter.List() {
			fmt.Println(name)
		}

		return "", exitCodeOk
	} else if opts.Fuzz.ListStrategies {
		for _, name := range tavorFuzzStrategy.List() {
			fmt.Println(name)
		}

		return "", exitCodeOk
	} else if opts.Fuzz.Exec.ListExecArgumentTypes || opts.Reduce.Exec.ListExecArgumentTypes {
		for _, name := range execArgumentTypes {
			fmt.Println(name)
		}

		return "", exitCodeOk
	} else if opts.Reduce.ListStrategies {
		for _, name := range tavorReduceStrategy.List() {
			fmt.Println(name)
		}

		return "", exitCodeOk
	}

	if err != nil {
		return "", exitError(err.Error())
	}

	if completion {
		return "", exitCodeBashCompletion
	}

	if opts.General.Debug {
		log.LevelDebug()
	} else if opts.General.Verbose {
		log.LevelInfo()
	} else {
		log.LevelWarn()
	}

	if opts.Global.Seed == 0 {
		opts.Global.Seed = time.Now().UTC().UnixNano()
	}

	if opts.Global.MaxRepeat < 1 {
		return "", exitError("max repeats has to be at least 1")
	}

	if opts.Fuzz.ResultFolder != "" {
		if err := osutil.DirExists(string(opts.Fuzz.ResultFolder)); err != nil {
			return "", exitError("result-folder invalid: %v", err)
		}
	}
	if opts.Fuzz.ResultSeparator != "" {
		if t, err := strconv.Unquote(`"` + opts.Fuzz.ResultSeparator + `"`); err == nil {
			opts.Fuzz.ResultSeparator = t
		}
	}

	if opts.Reduce.Exec.ExecArgumentType != "" {
		found := false

		for _, v := range execArgumentTypes {
			if string(opts.Reduce.Exec.ExecArgumentType) == v {
				found = true

				break
			}
		}

		if !found {
			return "", exitError(fmt.Sprintf("%q is an unknown exec argument type", opts.Reduce.Exec.ExecArgumentType))
		}
	}
	if opts.Reduce.Exec.Exec != "" {
		if !opts.Reduce.Exec.ExecExactExitCode && !opts.Reduce.Exec.ExecExactStderr && !opts.Reduce.Exec.ExecExactStdout && opts.Reduce.Exec.ExecMatchStderr == "" && opts.Reduce.Exec.ExecMatchStdout == "" {
			return "", exitError("At least one exec-exact or exec-match argument has to be given")
		}
	}
	if opts.Reduce.ResultSeparator != "" {
		if t, err := strconv.Unquote(`"` + opts.Reduce.ResultSeparator + `"`); err == nil {
			opts.Reduce.ResultSeparator = t
		}
	}

	log.Infof("using seed %d", opts.Global.Seed)
	log.Infof("using max repeat %d", opts.Global.MaxRepeat)

	return p.Active.Name, exitCodeOk
}

func exitError(format string, args ...interface{}) exitCodeType {
	fmt.Fprintf(os.Stderr, format+"\n", args...)

	return exitCodeError
}

func applyFilters(opts *options, filterNames []fuzzFilter, doc token.Token) (token.Token, error) {
	if len(filterNames) > 0 {
		var err error
		var filters []tavorFuzzFilter.Filter

		for _, name := range filterNames {
			filt, err := tavorFuzzFilter.New(string(name))
			if err != nil {
				return nil, err
			}

			filters = append(filters, filt)

			log.Infof("using %s fuzzing filter", name)
		}

		doc, err = tavorFuzzFilter.ApplyFilters(filters, doc)
		if err != nil {
			return nil, err
		}

		if opts.Format.PrintInternal {
			log.Info("Internal AST:")

			token.PrettyPrintInternalTree(os.Stdout, doc)
		}

		if opts.Format.Print {
			log.Info("AST:")

			token.PrettyPrintTree(os.Stdout, doc)
		}
	}

	return doc, nil
}

func mainCmd(args []string) exitCodeType {
	var opts = new(options)

	command, exitCode := checkArguments(args, opts)
	if command == "" {
		return exitCode
	}

	tavor.MaxRepeat = opts.Global.MaxRepeat

	log.Infof("open file %s", opts.Format.FormatFile)

	file, err := os.Open(string(opts.Format.FormatFile))
	if err != nil {
		return exitError("cannot open tavor file %s: %v", opts.Format.FormatFile, err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	doc, err := parser.ParseTavor(file)
	if err != nil {
		return exitError("cannot parse tavor file: %v", err)
	}

	log.Info("format file is valid")

	if opts.Format.PrintInternal {
		log.Info("Internal AST:")

		token.PrettyPrintInternalTree(os.Stdout, doc)
	}

	if opts.Format.Print {
		log.Info("AST:")

		token.PrettyPrintTree(os.Stdout, doc)
	}

	if opts.Format.Check {
		return exitCodeOk
	}

	r := rand.New(rand.NewSource(opts.Global.Seed))

	switch command {
	case "fuzz":
		doc, err = applyFilters(opts, opts.Fuzz.Filter.Filters, doc)
		if err != nil {
			return exitError("cannot apply filters: %v", err)
		}

		log.Infof("counted %d overall permutations", doc.PermutationsAll())

		strat, err := tavorFuzzStrategy.New(string(opts.Fuzz.Strategy), doc)
		if err != nil {
			return exitError(err.Error())
		}

		log.Infof("using %s fuzzing strategy", opts.Fuzz.Strategy)

		folder := opts.Fuzz.ResultFolder
		if len(folder) > 0 && folder[len(folder)-1] != '/' {
			folder += "/"
		}

		ch, err := strat.Fuzz(r)
		if err != nil {
			return exitError(err.Error())
		}

		if opts.Fuzz.Exec.Exec != "" {
			execs := strings.Split(opts.Fuzz.Exec.Exec, " ")
			var execFileArguments []int
			for i, v := range execs {
				if v == "TAVOR_FUZZ_FILE" {
					execFileArguments = append(execFileArguments, i)
				}
			}

			var matchStderr *regexp.Regexp
			var matchStdout *regexp.Regexp

			if opts.Fuzz.Exec.ExecMatchStderr != "" {
				matchStderr = regexp.MustCompile(opts.Fuzz.Exec.ExecMatchStderr)
			}
			if opts.Fuzz.Exec.ExecMatchStdout != "" {
				matchStdout = regexp.MustCompile(opts.Fuzz.Exec.ExecMatchStdout)
			}

			stepID := 1

			writeTmpFile := func(docOut string) (*os.File, error) {
				tmp, err := ioutil.TempFile(string(folder), fmt.Sprintf("fuzz-%d-", stepID))
				if err != nil {
					return nil, fmt.Errorf("Cannot create tmp file: %v", err)
				}
				_, err = tmp.WriteString(docOut)
				if err != nil {
					return nil, fmt.Errorf("Cannot write to tmp file: %v", err)
				}

				return tmp, nil
			}

		GENERATION:
			for i := range ch {
				docOut := doc.String()

				log.Infof("Test %d", stepID)

				var tmp *os.File

				var cmdExitCode int
				var cmdStderr bytes.Buffer
				var cmdStdout bytes.Buffer

				if string(opts.Fuzz.Exec.ExecArgumentType) == "argument" {
					tmp, err = writeTmpFile(docOut)
					if err != nil {
						return exitError(err.Error())
					}

					for _, v := range execFileArguments {
						execs[v] = tmp.Name()
					}
				}

				execCommand := exec.Command(execs[0], execs[1:]...)

				if string(opts.Fuzz.Exec.ExecArgumentType) == "environment" {
					tmp, err = writeTmpFile(docOut)
					if err != nil {
						return exitError(err.Error())
					}

					execCommand.Env = []string{fmt.Sprintf("TAVOR_FUZZ_FILE=%s", tmp.Name())}
				}

				if opts.General.Verbose || opts.General.Debug {
					execCommand.Stderr = io.MultiWriter(&cmdStderr, os.Stderr)
					execCommand.Stdout = io.MultiWriter(&cmdStdout, os.Stdout)
				} else {
					execCommand.Stderr = &cmdStderr
					execCommand.Stdout = &cmdStdout
				}

				stdin, err := execCommand.StdinPipe()
				if err != nil {
					return exitError("Could not get stdin pipe: %s", err)
				}

				if tmp != nil {
					log.Infof("Test %q", tmp.Name())
				}

				err = execCommand.Start()
				if err != nil {
					return exitError("Could not start exce: %s", err)
				}

				if string(opts.Fuzz.Exec.ExecArgumentType) == "stdin" {
					_, err := stdin.Write([]byte(docOut))
					if err != nil {
						return exitError("Could not write stdin to exec: %s", err)
					}

					if err := stdin.Close(); err != nil {
						panic(err)
					}
				}

				err = execCommand.Wait()

				if err == nil {
					cmdExitCode = 0
				} else if e, ok := err.(*exec.ExitError); ok {
					cmdExitCode = e.Sys().(syscall.WaitStatus).ExitStatus()
				} else {
					return exitError("Could not execute exec successfully: %s", err)
				}

				log.Infof("Exit status was %d", cmdExitCode)

				oks := 0
				oksNeeded := 0

				gotError := false

				if opts.Fuzz.Exec.ExecExactExitCode != -1 {
					oksNeeded++

					if opts.Fuzz.Exec.ExecExactExitCode == cmdExitCode {
						log.Infof("Same exit code")

						oks++
					} else {
						log.Infof("Not the same exit code")
					}
				}
				if opts.Fuzz.Exec.ExecExactStderr != "" {
					oksNeeded++

					if opts.Fuzz.Exec.ExecExactStderr == cmdStderr.String() {
						log.Infof("Same stderr")

						oks++
					} else {
						log.Infof("Not the same stderr")
					}
				}
				if opts.Fuzz.Exec.ExecExactStdout != "" {
					oksNeeded++

					if opts.Fuzz.Exec.ExecExactStdout == cmdStdout.String() {
						log.Infof("Same stdout")

						oks++
					} else {
						log.Infof("Not the same stdout")
					}
				}
				if matchStderr != nil {
					oksNeeded++

					if matchStderr.Match(cmdStderr.Bytes()) {
						log.Infof("Same stderr matching")

						oks++
					} else {
						log.Infof("Not the same stderr matching")
					}
				}
				if matchStdout != nil {
					oksNeeded++

					if matchStdout.Match(cmdStdout.Bytes()) {
						log.Infof("Same stdout matching")

						oks++
					} else {
						log.Infof("Not the same stdout matching")
					}
				}

				if oksNeeded == 0 {
					log.Warnf("Not defined what to compare")
				} else {
					if oks == oksNeeded {
						log.Infof("Same output")
					} else {
						log.Infof("Not the same output")

						gotError = true

						if opts.Fuzz.Exec.ExecDoNotRemoveTmpFilesOnError || string(folder) != "" {
							if tmp == nil {
								tmp, err = writeTmpFile(docOut)
								if err != nil {
									return exitError(err.Error())
								}
							}

							log.Infof("Written to %q", tmp.Name())
						}

						if opts.Fuzz.Exec.ExitOnError {
							break GENERATION
						}
					}
				}

				if !opts.Fuzz.Exec.ExecDoNotRemoveTmpFiles && (!gotError || (!opts.Fuzz.Exec.ExecDoNotRemoveTmpFilesOnError && string(folder) == "")) {
					if tmp != nil {
						err = os.Remove(tmp.Name())
						if err != nil {
							log.Errorf("Could not remove tmp file %q: %s", tmp.Name(), err)
						}
					}
				}

				ch <- i

				stepID++
			}
		} else if opts.Fuzz.Exec.Script != "" {
			execs := strings.Split(opts.Fuzz.Exec.Script, " ")

			execCommand := exec.Command(execs[0], execs[1:]...)

			stdin, err := execCommand.StdinPipe()
			if err != nil {
				return exitError("Could not get stdin pipe: %s", err)
			}
			defer func() {
				if err := stdin.Close(); err != nil {
					panic(err)
				}
			}()

			stdout, err := execCommand.StdoutPipe()
			if err != nil {
				return exitError("Could not get stdout pipe: %s", err)
			}
			defer func() {
				if err := stdout.Close(); err != nil {
					panic(err)
				}
			}()

			execCommand.Stderr = os.Stderr

			stdoutReader := bufio.NewReader(stdout)

			log.Infof("Execute script %q", opts.Fuzz.Exec.Script)

			err = execCommand.Start()
			if err != nil {
				return exitError("Could not start script: %s", err)
			}

		GENERATIONSC:
			for i := range ch {
				_, err = stdin.Write([]byte("Generation\n"))
				if err != nil {
					return exitError("Could not write stdin to script: %s", err)
				}
				_, err = stdin.Write([]byte(doc.String()))
				if err != nil {
					return exitError("Could not write stdin to script: %s", err)
				}
				_, err = stdin.Write([]byte(opts.Fuzz.ResultSeparator))
				if err != nil {
					return exitError("Could not write stdin to script: %s", err)
				}

				feed, err := stdoutReader.ReadString('\n')
				if err != nil {
					return exitError("Could not read stdout from script: %s", err)
				}

				switch feed {
				case "YES\n":
					log.Infof("Same output")
				case "NO\n":
					log.Infof("Not the same output")

					if opts.Fuzz.Exec.ExitOnError {
						break GENERATIONSC
					}
				default:
					return exitError("Feedback from script was not YES nor NO: %s", feed)
				}

				ch <- i
			}

			_, err = stdin.Write([]byte("Exit\n"))
			if err != nil {
				return exitError("Could not write stdin to script: %s", err)
			}

			if err := stdin.Close(); err != nil {
				panic(err)
			}
			if err := stdout.Close(); err != nil {
				panic(err)
			}

			err = execCommand.Wait()

			var execExitCode int

			if err == nil {
				execExitCode = 0
			} else if e, ok := err.(*exec.ExitError); ok {
				execExitCode = e.Sys().(syscall.WaitStatus).ExitStatus()
			} else {
				return exitError("Could not execute exec successfully: %s", err)
			}

			log.Infof("Exit status was %d", execExitCode)
		} else {
			another := false

			for i := range ch {
				if folder == "" {
					if opts.General.Debug {
						if another {
							fmt.Println()
						} else {
							another = true
						}
					}

					log.Debug("result:")
					fmt.Print(doc.String())
					fmt.Print(opts.Fuzz.ResultSeparator)
				} else {
					out := doc.String()
					sum := md5.Sum([]byte(out))

					file := fmt.Sprintf("%s%x%s", folder, sum, opts.Fuzz.ResultExtensions)

					log.Infof("write result to %s", file)

					if err := ioutil.WriteFile(file, []byte(out), 0644); err != nil {
						return exitError("error writing to %s: %v", file, err)
					}
				}

				ch <- i
			}
		}
	case "graph":
		doc, err = applyFilters(opts, opts.Graph.Filter.Filters, doc)
		if err != nil {
			return exitError("cannot apply filters: %v", err)
		}

		graph.WriteDot(doc, os.Stdout)
	case "reduce", "validate":
		inputFile := opts.Validate.InputFile

		if command == "reduce" {
			inputFile = opts.Reduce.InputFile
		}

		input, err := os.Open(string(inputFile))
		if err != nil {
			return exitError("cannot open input file %s: %v", inputFile, err)
		}
		defer func() {
			if err := input.Close(); err != nil {
				panic(err)
			}
		}()

		errs := parser.ParseInternal(doc, input)

		if len(errs) == 0 {
			log.Info("input file is valid")
		} else {
			log.Info("input file is invalid")

			for _, err := range errs {
				log.Error(err)
			}

			return exitCodeInvalidInputFile
		}

		if command == "reduce" {
			strat, err := tavorReduceStrategy.New(string(opts.Reduce.Strategy), doc)
			if err != nil {
				return exitError(err.Error())
			}

			log.Infof("using %s reducing strategy", opts.Reduce.Strategy)

			if opts.Reduce.Exec.Exec != "" {
				execs := strings.Split(opts.Reduce.Exec.Exec, " ")
				var execFileArguments []int
				for i, v := range execs {
					if v == "TAVOR_DD_FILE" {
						execFileArguments = append(execFileArguments, i)
					}
				}

				stepID := 1

				docOut := doc.String()

				tmp, err := ioutil.TempFile("", fmt.Sprintf("dd-%d-", stepID))
				if err != nil {
					return exitError("Cannot create tmp file: %s", err)
				}
				_, err = tmp.WriteString(docOut)
				if err != nil {
					return exitError("Cannot write to tmp file: %s", err)
				}

				log.Infof("Execute %q to get original outputs with %q", opts.Reduce.Exec.Exec, tmp.Name())

				var execExitCode int
				var execStderr bytes.Buffer
				var execStdout bytes.Buffer

				var matchStderr *regexp.Regexp
				var matchStdout *regexp.Regexp

				if opts.Reduce.Exec.ExecMatchStderr != "" {
					matchStderr = regexp.MustCompile(opts.Reduce.Exec.ExecMatchStderr)
				}
				if opts.Reduce.Exec.ExecMatchStdout != "" {
					matchStdout = regexp.MustCompile(opts.Reduce.Exec.ExecMatchStdout)
				}

				if string(opts.Reduce.Exec.ExecArgumentType) == "argument" {
					for _, v := range execFileArguments {
						execs[v] = tmp.Name()
					}
				}

				execCommand := exec.Command(execs[0], execs[1:]...)

				if string(opts.Reduce.Exec.ExecArgumentType) == "environment" {
					execCommand.Env = []string{fmt.Sprintf("TAVOR_DD_FILE=%s", tmp.Name())}
				}

				if opts.General.Verbose || opts.General.Debug {
					execCommand.Stderr = io.MultiWriter(&execStderr, os.Stderr)
					execCommand.Stdout = io.MultiWriter(&execStdout, os.Stdout)
				} else {
					execCommand.Stderr = &execStderr
					execCommand.Stdout = &execStdout
				}

				stdin, err := execCommand.StdinPipe()
				if err != nil {
					return exitError("Could not get stdin pipe: %s", err)
				}

				err = execCommand.Start()
				if err != nil {
					return exitError("Could not start exce: %s", err)
				}

				if string(opts.Reduce.Exec.ExecArgumentType) == "stdin" {
					_, err := stdin.Write([]byte(docOut))
					if err != nil {
						return exitError("Could not write stdin to exec: %s", err)
					}

					if err := stdin.Close(); err != nil {
						panic(err)
					}
				}

				err = execCommand.Wait()

				if err == nil {
					execExitCode = 0
				} else if e, ok := err.(*exec.ExitError); ok {
					execExitCode = e.Sys().(syscall.WaitStatus).ExitStatus()
				} else {
					return exitError("Could not execute exec successfully: %s", err)
				}

				log.Infof("Exit status was %d", execExitCode)

				if matchStderr != nil && !matchStderr.Match(execStderr.Bytes()) {
					return exitError("Original output does not match stderr match pattern")
				}
				if matchStdout != nil && !matchStdout.Match(execStdout.Bytes()) {
					return exitError("Original output does not match stdout match pattern")
				}

				if !opts.Reduce.Exec.ExecDoNotRemoveTmpFiles {
					err = os.Remove(tmp.Name())
					if err != nil {
						log.Errorf("Could not remove tmp file %q: %s", tmp.Name(), err)
					}
				}

				contin, feedback, err := strat.Reduce()
				if err != nil {
					return exitError(err.Error())
				}

				for i := range contin {
					stepID++

					docOut := doc.String()

					tmp, err := ioutil.TempFile("", fmt.Sprintf("dd-%d-", stepID))
					if err != nil {
						return exitError("Cannot create tmp file: %s", err)
					}
					_, err = tmp.WriteString(docOut)
					if err != nil {
						return exitError("Cannot write to tmp file: %s", err)
					}

					log.Infof("Test %q", tmp.Name())

					var cmdExitCode int
					var cmdStderr bytes.Buffer
					var cmdStdout bytes.Buffer

					if string(opts.Reduce.Exec.ExecArgumentType) == "argument" {
						for _, v := range execFileArguments {
							execs[v] = tmp.Name()
						}
					}

					execCommand := exec.Command(execs[0], execs[1:]...)

					if string(opts.Reduce.Exec.ExecArgumentType) == "environment" {
						execCommand.Env = []string{fmt.Sprintf("TAVOR_DD_FILE=%s", tmp.Name())}
					}

					if opts.General.Verbose || opts.General.Debug {
						execCommand.Stderr = io.MultiWriter(&cmdStderr, os.Stderr)
						execCommand.Stdout = io.MultiWriter(&cmdStdout, os.Stdout)
					} else {
						execCommand.Stderr = &cmdStderr
						execCommand.Stdout = &cmdStdout
					}

					stdin, err := execCommand.StdinPipe()
					if err != nil {
						return exitError("Could not get stdin pipe: %s", err)
					}

					err = execCommand.Start()
					if err != nil {
						return exitError("Could not start exce: %s", err)
					}

					if string(opts.Reduce.Exec.ExecArgumentType) == "stdin" {
						_, err := stdin.Write([]byte(docOut))
						if err != nil {
							return exitError("Could not write stdin to exec: %s", err)
						}

						if err := stdin.Close(); err != nil {
							panic(err)
						}
					}

					err = execCommand.Wait()

					if err == nil {
						cmdExitCode = 0
					} else if e, ok := err.(*exec.ExitError); ok {
						cmdExitCode = e.Sys().(syscall.WaitStatus).ExitStatus()
					} else {
						return exitError("Could not execute exec successfully: %s", err)
					}

					log.Infof("Exit status was %d", cmdExitCode)

					oks := 0
					oksNeeded := 0

					if opts.Reduce.Exec.ExecExactExitCode {
						oksNeeded++

						if execExitCode == cmdExitCode {
							log.Infof("Same exit code")

							oks++
						} else {
							log.Infof("Not the same exit code")
						}
					}
					if opts.Reduce.Exec.ExecExactStderr {
						oksNeeded++

						if execStderr.String() == cmdStderr.String() {
							log.Infof("Same stderr")

							oks++
						} else {
							log.Infof("Not the same stderr")
						}
					}
					if opts.Reduce.Exec.ExecExactStdout {
						oksNeeded++

						if execStdout.String() == cmdStdout.String() {
							log.Infof("Same stdout")

							oks++
						} else {
							log.Infof("Not the same stdout")
						}
					}
					if matchStderr != nil {
						oksNeeded++

						if matchStderr.Match(cmdStderr.Bytes()) {
							log.Infof("Same stderr matching")

							oks++
						} else {
							log.Infof("Not the same stderr matching")
						}
					}
					if matchStdout != nil {
						oksNeeded++

						if matchStdout.Match(cmdStdout.Bytes()) {
							log.Infof("Same stdout matching")

							oks++
						} else {
							log.Infof("Not the same stdout matching")
						}
					}

					if oksNeeded == 0 {
						log.Warnf("Not defined what to compare")
					} else {
						if oks == oksNeeded {
							log.Infof("Same output, continue delta")

							feedback <- tavorReduceStrategy.Good
						} else {
							log.Infof("Not the same output, do another step")

							feedback <- tavorReduceStrategy.Bad
						}
					}

					if !opts.Reduce.Exec.ExecDoNotRemoveTmpFiles {
						err = os.Remove(tmp.Name())
						if err != nil {
							log.Errorf("Could not remove tmp file %q: %s", tmp.Name(), err)
						}
					}

					contin <- i
				}
			} else if opts.Reduce.Exec.Script != "" {
				execs := strings.Split(opts.Reduce.Exec.Script, " ")

				execCommand := exec.Command(execs[0], execs[1:]...)

				stdin, err := execCommand.StdinPipe()
				if err != nil {
					return exitError("Could not get stdin pipe: %s", err)
				}
				defer func() {
					if err := stdin.Close(); err != nil {
						panic(err)
					}
				}()

				stdout, err := execCommand.StdoutPipe()
				if err != nil {
					return exitError("Could not get stdout pipe: %s", err)
				}
				defer func() {
					if err := stdout.Close(); err != nil {
						panic(err)
					}
				}()

				execCommand.Stderr = os.Stderr

				stdoutReader := bufio.NewReader(stdout)

				log.Infof("Execute script %q", opts.Reduce.Exec.Script)

				err = execCommand.Start()
				if err != nil {
					return exitError("Could not start script: %s", err)
				}

				log.Infof("Send original output to script")

				_, err = stdin.Write([]byte(doc.String()))
				if err != nil {
					return exitError("Could not write stdin to script: %s", err)
				}
				_, err = stdin.Write([]byte(opts.Reduce.ResultSeparator))
				if err != nil {
					return exitError("Could not write stdin to script: %s", err)
				}

				feed, err := stdoutReader.ReadString('\n')
				if err != nil {
					return exitError("Could not read stdout from script: %s", err)
				}

				if feed != "OK\n" {
					return exitError("Feedback from script to orignal was not OK: %s", feed)
				}

				contin, feedback, err := strat.Reduce()
				if err != nil {
					return exitError(err.Error())
				}

				for i := range contin {
					_, err = stdin.Write([]byte(doc.String()))
					if err != nil {
						return exitError("Could not write stdin to script: %s", err)
					}
					_, err = stdin.Write([]byte(opts.Reduce.ResultSeparator))
					if err != nil {
						return exitError("Could not write stdin to script: %s", err)
					}

					feed, err := stdoutReader.ReadString('\n')
					if err != nil {
						return exitError("Could not read stdout from script: %s", err)
					}

					switch feed {
					case "YES\n":
						log.Infof("Same output, continue delta")

						feedback <- tavorReduceStrategy.Good
					case "NO\n":
						log.Infof("Not the same output, do another step")

						feedback <- tavorReduceStrategy.Bad
					default:
						return exitError("Feedback from script to orignal was not YES nor NO: %s", feed)
					}

					contin <- i
				}

				if err := stdin.Close(); err != nil {
					panic(err)
				}
				if err := stdout.Close(); err != nil {
					panic(err)
				}

				err = execCommand.Wait()

				var execExitCode int

				if err == nil {
					execExitCode = 0
				} else if e, ok := err.(*exec.ExitError); ok {
					execExitCode = e.Sys().(syscall.WaitStatus).ExitStatus()
				} else {
					return exitError("Could not execute exec successfully: %s", err)
				}

				log.Infof("Exit status was %d", execExitCode)
			} else {
				readCLI := bufio.NewReader(os.Stdin)

				contin, feedback, err := strat.Reduce()
				if err != nil {
					return exitError(err.Error())
				}

				for i := range contin {
					log.Debug("result:")
					fmt.Print(doc.String())
					fmt.Print(opts.Reduce.ResultSeparator)

					for {
						fmt.Printf("\nDo the constraints of the original input still hold for this generation? [yes|no]: ")

						line, _, err := readCLI.ReadLine()
						if err != nil {
							return exitError("reading from CLI failed: %v", err)
						}

						if s := strings.ToUpper(string(line)); s == "YES" {
							feedback <- tavorReduceStrategy.Good

							break
						} else if s == "NO" {
							feedback <- tavorReduceStrategy.Bad

							break
						}
					}

					contin <- i
				}
			}

			log.Info("reduced to minimum")

			log.Debug("result:")
			fmt.Print(doc.String())
			fmt.Print(opts.Reduce.ResultSeparator)
		}
	default:
		return exitError("unknown command %q", command)
	}

	return exitCodeOk
}

func main() {
	exitCode := int(mainCmd(os.Args[1:]))

	os.Exit(exitCode)
}
