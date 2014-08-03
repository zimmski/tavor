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

	"github.com/zimmski/tavor"
	fuzzFilter "github.com/zimmski/tavor/fuzz/filter"
	fuzzStrategy "github.com/zimmski/tavor/fuzz/strategy"
	"github.com/zimmski/tavor/graph"
	"github.com/zimmski/tavor/log"
	"github.com/zimmski/tavor/parser"
	reduceStrategy "github.com/zimmski/tavor/reduce/strategy"
	"github.com/zimmski/tavor/token"
)

const (
	returnOk = iota
	returnHelp
	returnBashCompletion
	returnInvalidInputFile
	returnError
)

var opts struct {
	General struct {
		Debug   bool `long:"debug" description:"Debug log output"`
		Help    bool `long:"help" description:"Show this help message"`
		Verbose bool `long:"verbose" description:"Verbose log output"`
		Version bool `long:"version" description:"Print the version of this program"`
	} `group:"General options"`

	Global struct {
		Seed int64 `long:"seed" description:"Seed for all the randomness"`
	} `group:"Global options"`

	Format struct {
		FormatFile    flags.Filename `long:"format-file" description:"Input tavor format file" required:"true"`
		PrintInternal bool           `long:"print-internal" description:"Prints the internal AST of the parsed format file"`
		Validate      bool           `long:"validate" description:"Just validate the format file and exit"`
	} `group:"Format file options"`

	Fuzz struct {
		Filters          FuzzFilters    `long:"filter" description:"Fuzzing filter to apply"`
		ListFilters      bool           `long:"list-filters" description:"List all available fuzzing filters"`
		Strategy         FuzzStrategy   `long:"strategy" description:"The fuzzing strategy" default:"random"`
		ListStrategies   bool           `long:"list-strategies" description:"List all available fuzzing strategies"`
		ResultFolder     flags.Filename `long:"result-folder" description:"Save every fuzzing result with the MD5 checksum as filename in this folder"`
		ResultExtensions string         `long:"result-extension" description:"If result-folder is used this will be the extension of every filename"`
		ResultSeparator  string         `long:"result-separator" description:"Separates result outputs of each fuzzing step" default:"\n"`
	} `command:"fuzz" description:"Fuzz the given format file"`

	Graph struct {
		Filters     FuzzFilters `long:"filter" description:"Fuzzing filter to apply"`
		ListFilters bool        `long:"list-filters" description:"List all available fuzzing filters"`
	} `command:"graph" description:"Generate a DOT file out of the internal AST"`

	Reduce struct {
		Exec                    string           `long:"exec" description:"Execute this binary with possible arguments to test a delta-debugging step"`
		ExecExactExitCode       bool             `long:"exec-exact-exit-code" description:"Same exit code has to be present to reduce further"`
		ExecExactStderr         bool             `long:"exec-exact-stderr" description:"Same stderr output has to be present to reduce further"`
		ExecExactStdout         bool             `long:"exec-exact-stdout" description:"Same stdout output has to be present to reduce further"`
		ExecMatchStderr         string           `long:"exec-match-stderr" description:"Searches through stderr via the given regex. A match has to be present to reduce further"`
		ExecMatchStdout         string           `long:"exec-match-stdout" description:"Searches through stdout via the given regex. A match has to be present to reduce further"`
		ExecDoNotRemoveTmpFiles bool             `long:"exec-do-not-remove-tmp-files" description:"If set tmp files for delta debugging are not removed"`
		ExecArgumentType        ExecArgumentType `long:"exec-argument-type" description:"How the delta-debugging step is given to the binary" default:"environment"`
		ListExecArgumentTypes   bool             `long:"list-exec-argument-types" description:"List all available exec argument types"`

		InputFile       flags.Filename `long:"input-file" description:"Input file which gets parsed, validated and delta-debugged via the format file" required:"true"`
		Strategy        ReduceStrategy `long:"strategy" description:"The reducing strategy" default:"BinarySearch"`
		ListStrategies  bool           `long:"list-strategies" description:"List all available reducing strategies"`
		ResultSeparator string         `long:"result-separator" description:"Separates result outputs of each reducing step" default:"\n"`
	} `command:"reduce" description:"Reduce the given input file"`

	Validate struct {
		InputFile flags.Filename `long:"input-file" description:"Input file which gets parsed and validated via the format file" required:"true"`
	} `command:"validate" description:"Validate the given input file"`
}

var ExecArgumentTypes = []string{
	"argument",
	"environment",
	"stdin",
}

type ExecArgumentType string

func (e ExecArgumentType) Complete(match string) []flags.Completion {
	var items []flags.Completion

	for _, name := range ExecArgumentTypes {
		if strings.HasPrefix(name, match) {
			items = append(items, flags.Completion{
				Item: name,
			})
		}
	}

	return items
}

type FuzzFilter string
type FuzzFilters []FuzzFilter

func (s FuzzFilters) Complete(match string) []flags.Completion {
	var items []flags.Completion

	for _, name := range fuzzFilter.List() {
		if strings.HasPrefix(name, match) {
			items = append(items, flags.Completion{
				Item: name,
			})
		}
	}

	return items
}

type FuzzStrategy string

func (s *FuzzStrategy) Complete(match string) []flags.Completion {
	var items []flags.Completion

	for _, name := range fuzzStrategy.List() {
		if strings.HasPrefix(name, match) {
			items = append(items, flags.Completion{
				Item: name,
			})
		}
	}

	return items
}

type ReduceStrategy string

func (s *ReduceStrategy) Complete(match string) []flags.Completion {
	var items []flags.Completion

	for _, name := range reduceStrategy.List() {
		if strings.HasPrefix(name, match) {
			items = append(items, flags.Completion{
				Item: name,
			})
		}
	}

	return items
}

func checkArguments() string {
	p := flags.NewNamedParser("tavor", flags.PassDoubleDash)

	p.ShortDescription = "A fuzzing and delta-debugging platform."

	if _, err := p.AddGroup("Tavor", "Tavor arguments", &opts); err != nil {
		exitError(err.Error())
	}

	_, err := p.Parse()
	if opts.General.Help || len(os.Args) == 1 {
		p.WriteHelp(os.Stdout)

		os.Exit(returnHelp)
	} else if opts.General.Version {
		fmt.Printf("Tavor v%s\n", tavor.Version)

		os.Exit(returnOk)
	} else if opts.Fuzz.ListFilters || opts.Graph.ListFilters {
		for _, name := range fuzzFilter.List() {
			fmt.Println(name)
		}

		os.Exit(returnOk)
	} else if opts.Fuzz.ListStrategies {
		for _, name := range fuzzStrategy.List() {
			fmt.Println(name)
		}

		os.Exit(returnOk)
	} else if opts.Reduce.ListExecArgumentTypes {
		for _, name := range ExecArgumentTypes {
			fmt.Println(name)
		}

		os.Exit(returnOk)
	} else if opts.Reduce.ListStrategies {
		for _, name := range reduceStrategy.List() {
			fmt.Println(name)
		}

		os.Exit(returnOk)
	}

	if err != nil {
		exitError(err.Error())
	}

	if len(os.Getenv("GO_FLAGS_COMPLETION")) != 0 {
		os.Exit(returnBashCompletion)
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

	if opts.Fuzz.ResultFolder != "" {
		if err := folderExists(string(opts.Fuzz.ResultFolder)); err != nil {
			exitError("result-folder invalid: %v", err)
		}
	}
	if opts.Fuzz.ResultSeparator != "" {
		if t, err := strconv.Unquote(`"` + opts.Fuzz.ResultSeparator + `"`); err == nil {
			opts.Fuzz.ResultSeparator = t
		}
	}

	if opts.Reduce.ExecArgumentType != "" {
		found := false

		for _, v := range ExecArgumentTypes {
			if string(opts.Reduce.ExecArgumentType) == v {
				found = true

				break
			}
		}

		if !found {
			exitError(fmt.Sprintf("%q is an unknown exec argument type", opts.Reduce.ExecArgumentType))
		}
	}
	if opts.Reduce.Exec != "" {
		if !opts.Reduce.ExecExactExitCode && !opts.Reduce.ExecExactStderr && !opts.Reduce.ExecExactStdout && opts.Reduce.ExecMatchStderr == "" && opts.Reduce.ExecMatchStdout == "" {
			exitError("At least one exec-exact or exec-match argument has to be given")
		}
	}
	if opts.Reduce.ResultSeparator != "" {
		if t, err := strconv.Unquote(`"` + opts.Reduce.ResultSeparator + `"`); err == nil {
			opts.Reduce.ResultSeparator = t
		}
	}

	log.Infof("using seed %d", opts.Global.Seed)

	return p.Active.Name
}

func exitError(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)

	os.Exit(returnError)
}

func folderExists(folder string) error {
	f, err := os.Open(folder)
	if err != nil {
		return fmt.Errorf("%q does not exist", folder)
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return fmt.Errorf("could not stat %q", folder)
	}

	if !fi.Mode().IsDir() {
		return fmt.Errorf("%q is not a folder", folder)
	}

	return nil
}

func applyFilters(filterNames []FuzzFilter, doc token.Token) token.Token {
	if len(filterNames) != 0 {
		var err error
		var filters []fuzzFilter.Filter

		for _, name := range filterNames {
			filt, err := fuzzFilter.New(string(name))
			if err != nil {
				exitError(err.Error())
			}

			filters = append(filters, filt)

			log.Infof("using %s fuzzing filter", name)
		}

		doc, err = fuzzFilter.ApplyFilters(filters, doc)
		if err != nil {
			exitError(err.Error())
		}
	}

	return doc
}

func main() {
	command := checkArguments()

	log.Infof("open file %s", opts.Format.FormatFile)

	file, err := os.Open(string(opts.Format.FormatFile))
	if err != nil {
		exitError("cannot open tavor file %s: %v", opts.Format.FormatFile, err)
	}
	defer file.Close()

	doc, err := parser.ParseTavor(file)
	if err != nil {
		exitError("cannot parse tavor file: %v", err)
	}

	log.Info("format file is valid")

	if opts.Format.PrintInternal {
		tavor.PrettyPrintInternalTree(os.Stdout, doc)
	}

	if opts.Format.Validate {
		os.Exit(returnOk)
	}

	r := rand.New(rand.NewSource(opts.Global.Seed))

	switch command {
	case "fuzz":
		doc = applyFilters(opts.Fuzz.Filters, doc)

		log.Infof("counted %d overall permutations", doc.PermutationsAll())

		strat, err := fuzzStrategy.New(string(opts.Fuzz.Strategy), doc)
		if err != nil {
			exitError(err.Error())
		}

		log.Infof("using %s fuzzing strategy", opts.Fuzz.Strategy)

		ch, err := strat.Fuzz(r)
		if err != nil {
			exitError(err.Error())
		}

		folder := opts.Fuzz.ResultFolder
		if len(folder) != 0 && folder[len(folder)-1] != '/' {
			folder += "/"
		}

		another := false
		for i := range ch {
			if folder == "" {
				if !opts.General.Debug {
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
					exitError("error writing to %s: %v", file, err)
				}
			}

			ch <- i
		}
	case "graph":
		doc = applyFilters(opts.Graph.Filters, doc)

		graph.WriteDot(doc, os.Stdout)
	case "reduce", "validate":
		inputFile := opts.Validate.InputFile

		if command == "reduce" {
			inputFile = opts.Reduce.InputFile
		}

		input, err := os.Open(string(inputFile))
		if err != nil {
			exitError("cannot open input file %s: %v", inputFile, err)
		}
		defer input.Close()

		errs := parser.ParseInternal(doc, input)

		if len(errs) == 0 {
			log.Info("input file is valid")
		} else {
			log.Info("input file is invalid")

			for _, err := range errs {
				log.Error(err)
			}

			os.Exit(returnInvalidInputFile)
		}

		if command == "reduce" {
			strat, err := reduceStrategy.New(string(opts.Reduce.Strategy), doc)
			if err != nil {
				exitError(err.Error())
			}

			log.Infof("using %s reducing strategy", opts.Reduce.Strategy)

			if opts.Reduce.Exec != "" {
				execs := strings.Split(opts.Reduce.Exec, " ")
				var execDDFileArguments []int
				for i, v := range execs {
					if v == "TAVOR_DD_FILE" {
						execDDFileArguments = append(execDDFileArguments, i)
					}
				}

				stepId := 0

				docOut := doc.String()

				tmp, err := ioutil.TempFile("", fmt.Sprintf("dd-%d-", stepId))
				if err != nil {
					exitError("Cannot create tmp file: %s", err)
				}
				_, err = tmp.WriteString(docOut)
				if err != nil {
					exitError("Cannot write to tmp file: %s", err)
				}

				log.Infof("Execute %q to get original outputs with %q", opts.Reduce.Exec, tmp.Name())

				var execExitCode int
				var execStderr bytes.Buffer
				var execStdout bytes.Buffer

				var matchStderr *regexp.Regexp
				var matchStdout *regexp.Regexp

				if opts.Reduce.ExecMatchStderr != "" {
					matchStderr = regexp.MustCompile(opts.Reduce.ExecMatchStderr)
				}
				if opts.Reduce.ExecMatchStdout != "" {
					matchStdout = regexp.MustCompile(opts.Reduce.ExecMatchStdout)
				}

				if string(opts.Reduce.ExecArgumentType) == "argument" {
					for _, v := range execDDFileArguments {
						execs[v] = tmp.Name()
					}
				}

				execCommand := exec.Command(execs[0], execs[1:]...)

				if string(opts.Reduce.ExecArgumentType) == "environment" {
					execCommand.Env = []string{fmt.Sprintf("TAVOR_DD_FILE=%s", tmp.Name())}
				}

				execCommand.Stderr = io.MultiWriter(&execStderr, os.Stderr)
				execCommand.Stdout = io.MultiWriter(&execStdout, os.Stdout)

				stdin, err := execCommand.StdinPipe()
				if err != nil {
					exitError("Could not get stdin pipe: %s", err)
				}

				err = execCommand.Start()
				if err != nil {
					exitError("Could not start exce: %s", err)
				}

				if string(opts.Reduce.ExecArgumentType) == "stdin" {
					_, err := stdin.Write([]byte(docOut))
					if err != nil {
						exitError("Could not write stdin to exec: %s", err)
					}

					stdin.Close()
				}

				err = execCommand.Wait()

				if err == nil {
					execExitCode = 0
				} else if e, ok := err.(*exec.ExitError); ok {
					execExitCode = e.Sys().(syscall.WaitStatus).ExitStatus()
				} else {
					exitError("Could not execute exec successfully: %s", err)
				}

				log.Infof("Exit status was %d", execExitCode)

				if matchStderr != nil && !matchStderr.Match(execStderr.Bytes()) {
					exitError("Original output does not match stderr match pattern")
				}
				if matchStdout != nil && !matchStdout.Match(execStdout.Bytes()) {
					exitError("Original output does not match stdout match pattern")
				}

				if !opts.Reduce.ExecDoNotRemoveTmpFiles {
					err = os.Remove(tmp.Name())
					if err != nil {
						log.Errorf("Could not remove tmp file %q: %s", tmp.Name(), err)
					}
				}

				contin, feedback, err := strat.Reduce()
				if err != nil {
					exitError(err.Error())
				}

				for i := range contin {
					stepId++

					docOut := doc.String()

					tmp, err := ioutil.TempFile("", fmt.Sprintf("dd-%d-", stepId))
					if err != nil {
						exitError("Cannot create tmp file: %s", err)
					}
					_, err = tmp.WriteString(docOut)
					if err != nil {
						exitError("Cannot write to tmp file: %s", err)
					}

					log.Infof("Test %q", tmp.Name())

					var ddExitCode int
					var ddStderr bytes.Buffer
					var ddStdout bytes.Buffer

					if string(opts.Reduce.ExecArgumentType) == "argument" {
						for _, v := range execDDFileArguments {
							execs[v] = tmp.Name()
						}
					}

					execCommand := exec.Command(execs[0], execs[1:]...)

					if string(opts.Reduce.ExecArgumentType) == "environment" {
						execCommand.Env = []string{fmt.Sprintf("TAVOR_DD_FILE=%s", tmp.Name())}
					}

					execCommand.Stderr = io.MultiWriter(&ddStderr, os.Stderr)
					execCommand.Stdout = io.MultiWriter(&ddStdout, os.Stdout)

					stdin, err := execCommand.StdinPipe()
					if err != nil {
						exitError("Could not get stdin pipe: %s", err)
					}

					err = execCommand.Start()
					if err != nil {
						exitError("Could not start exce: %s", err)
					}

					if string(opts.Reduce.ExecArgumentType) == "stdin" {
						_, err := stdin.Write([]byte(docOut))
						if err != nil {
							exitError("Could not write stdin to exec: %s", err)
						}

						stdin.Close()
					}

					err = execCommand.Wait()

					if err == nil {
						ddExitCode = 0
					} else if e, ok := err.(*exec.ExitError); ok {
						ddExitCode = e.Sys().(syscall.WaitStatus).ExitStatus()
					} else {
						exitError("Could not execute exec successfully: %s", err)
					}

					log.Infof("Exit status was %d", ddExitCode)

					if !opts.Reduce.ExecDoNotRemoveTmpFiles {
						err = os.Remove(tmp.Name())
						if err != nil {
							log.Errorf("Could not remove tmp file %q: %s", tmp.Name(), err)
						}
					}

					oks := 0
					oksNeeded := 0

					if opts.Reduce.ExecExactExitCode {
						oksNeeded++

						if execExitCode == ddExitCode {
							log.Infof("Same exit code")

							oks++
						} else {
							log.Infof("Not the same exit code")
						}
					}
					if opts.Reduce.ExecExactStderr {
						oksNeeded++

						if execStderr.String() == ddStderr.String() {
							log.Infof("Same stderr")

							oks++
						} else {
							log.Infof("Not the same stderr")
						}
					}
					if opts.Reduce.ExecExactStdout {
						oksNeeded++

						if execStdout.String() == ddStdout.String() {
							log.Infof("Same stdout")

							oks++
						} else {
							log.Infof("Not the same stdout")
						}
					}
					if matchStderr != nil {
						oksNeeded++

						if matchStderr.Match(ddStderr.Bytes()) {
							log.Infof("Same stderr matching")

							oks++
						} else {
							log.Infof("Not the same stderr matching")
						}
					}
					if matchStdout != nil {
						oksNeeded++

						if matchStdout.Match(ddStdout.Bytes()) {
							log.Infof("Same stdout matching")

							oks++
						} else {
							log.Infof("Not the same stdout matching")
						}
					}

					if oks == oksNeeded {
						log.Infof("Same output, continue delta")

						feedback <- reduceStrategy.Bad
					} else {
						log.Infof("Not the same output, do another step")

						feedback <- reduceStrategy.Good
					}

					contin <- i
				}
			} else {
				readCLI := bufio.NewReader(os.Stdin)

				contin, feedback, err := strat.Reduce()
				if err != nil {
					exitError(err.Error())
				}

				for i := range contin {
					log.Debug("result:")
					fmt.Print(doc.String())
					fmt.Print(opts.Reduce.ResultSeparator)

					for {
						fmt.Printf("\nDoes the error still exist? [yes|no]: ")

						line, _, err := readCLI.ReadLine()
						if err != nil {
							exitError("reading from CLI failed: %v", err)
						}

						if s := strings.ToUpper(string(line)); s == "YES" {
							feedback <- reduceStrategy.Bad

							break
						} else if s == "NO" {
							feedback <- reduceStrategy.Good

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
		exitError("unknown command %q", command)
	}

	os.Exit(returnOk)
}
