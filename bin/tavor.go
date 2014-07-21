package main

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/jessevdk/go-flags"

	"github.com/zimmski/tavor"
	fuzzStrategy "github.com/zimmski/tavor/fuzz/strategy"
	"github.com/zimmski/tavor/graph"
	"github.com/zimmski/tavor/log"
	"github.com/zimmski/tavor/parser"
	reduceStrategy "github.com/zimmski/tavor/reduce/strategy"
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
		Strategy         FuzzStrategy   `long:"strategy" description:"The fuzzing strategy" default:"random"`
		ListStrategies   bool           `long:"list-strategies" description:"List all available strategies"`
		ResultFolder     flags.Filename `long:"result-folder" description:"Save every fuzzing result with the MD5 checksum as filename in this folder"`
		ResultExtensions string         `long:"result-extension" description:"If result-folder is used this will be the extension of every filename"`
		ResultSeparator  string         `long:"result-separator" description:"Separates result outputs of each fuzzing step" default:"\n"`
	} `command:"fuzz" description:"Fuzz the given format file"`

	Graph struct {
	} `command:"graph" description:"Generate a DOT file out of the internal AST"`

	Reduce struct {
		InputFile       flags.Filename `long:"input-file" description:"Input file which gets parsed, validated and delta-debugged via the format file" required:"true"`
		Strategy        ReduceStrategy `long:"strategy" description:"The reducing strategy" default:"BinarySearch"`
		ListStrategies  bool           `long:"list-strategies" description:"List all available strategies"`
		ResultSeparator string         `long:"result-separator" description:"Separates result outputs of each reducing step" default:"\n"`
	} `command:"reduce" description:"Reduce the given input file"`

	Validate struct {
		InputFile flags.Filename `long:"input-file" description:"Input file which gets parsed and validated via the format file" required:"true"`
	} `command:"validate" description:"Validate the given format file"`
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
	} else if opts.Fuzz.ListStrategies {
		for _, name := range fuzzStrategy.List() {
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

	log.Infof("Using seed %d", opts.Global.Seed)

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

func main() {
	command := checkArguments()

	log.Infof("Open file %s", opts.Format.FormatFile)

	file, err := os.Open(string(opts.Format.FormatFile))
	if err != nil {
		exitError("cannot open tavor file %s: %v", opts.Format.FormatFile, err)
	}
	defer file.Close()

	doc, err := parser.ParseTavor(file)
	if err != nil {
		exitError("cannot parse tavor file: %v", err)
	}

	log.Info("Format file is valid")

	if opts.Format.PrintInternal {
		tavor.PrettyPrintInternalTree(os.Stdout, doc)
	}

	if opts.Format.Validate {
		os.Exit(returnOk)
	}

	r := rand.New(rand.NewSource(opts.Global.Seed))

	switch command {
	case "fuzz":
		log.Infof("Counted %d overall permutations", doc.PermutationsAll())

		strat, err := fuzzStrategy.New(string(opts.Fuzz.Strategy), doc)
		if err != nil {
			exitError(err.Error())
		}

		log.Infof("Using %s strategy", opts.Fuzz.Strategy)

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

				log.Debug("Result:")
				fmt.Print(doc.String())
				fmt.Print(opts.Fuzz.ResultSeparator)
			} else {
				out := doc.String()
				sum := md5.Sum([]byte(out))

				file := fmt.Sprintf("%s%x%s", folder, sum, opts.Fuzz.ResultExtensions)

				log.Infof("Write result to %s", file)

				if err := ioutil.WriteFile(file, []byte(out), 0644); err != nil {
					exitError("error writing to %s: %v", file, err)
				}
			}

			ch <- i
		}
	case "graph":
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
			log.Info("Input file is valid")
		} else {
			log.Info("Input file is invalid")

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

			log.Infof("Using %s strategy", opts.Reduce.Strategy)

			contin, feedback, err := strat.Reduce()
			if err != nil {
				exitError(err.Error())
			}

			readCLI := bufio.NewReader(os.Stdin)

			for i := range contin {
				log.Debug("Result:")
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

			log.Info("Reduced to minimum")

			log.Debug("Result:")
			fmt.Print(doc.String())
			fmt.Print(opts.Reduce.ResultSeparator)
		}
	default:
		exitError("Unknown command %q", command)
	}

	os.Exit(returnOk)
}
