package main

import (
	"fmt"
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
		Strategy       FuzzStrategy `long:"strategy" description:"The fuzzing strategy" default:"random"`
		ListStrategies bool         `long:"list-strategies" description:"List all available strategies"`
	} `command:"fuzz" description:"Fuzz the given format file"`

	Graph struct {
	} `command:"graph" description:"Generate a DOT file out of the internal AST"`

	Reduce struct {
		InputFile      flags.Filename `long:"input-file" description:"Input file which gets parsed, validated and delta-debugged via the format file" required:"true"`
		Strategy       ReduceStrategy `long:"strategy" description:"The reducing strategy" default:"BinarySearch"`
		ListStrategies bool           `long:"list-strategies" description:"List all available strategies"`
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

	log.Infof("Using seed %d", opts.Global.Seed)

	return p.Active.Name
}

func exitError(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)

	os.Exit(returnError)
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

		another := false
		for i := range ch {
			if !opts.General.Debug {
				if another {
					fmt.Println()
				} else {
					another = true
				}
			}

			log.Debug("Result:")

			fmt.Print(doc.String())
			if opts.General.Debug {
				fmt.Println()
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

			for i := range contin {
				// TODO get user feedback for Good or Bad question, right now it is ok that everything is bad
				feedback <- reduceStrategy.Bad

				contin <- i
			}

			log.Info("Reduced to minimum")

			fmt.Print(doc.String())
			if opts.General.Debug {
				fmt.Println()
			}
		}
	default:
		exitError("Unknown command %q", command)
	}

	os.Exit(returnOk)
}
