package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/jessevdk/go-flags"

	"github.com/zimmski/tavor"
	"github.com/zimmski/tavor/fuzz/strategy"
	"github.com/zimmski/tavor/log"
	"github.com/zimmski/tavor/parser"
)

const (
	returnOk = iota
	returnHelp
	returnBashCompletion
	returnInvalidInputFile
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
		Strategy       Strategy `long:"strategy" description:"The fuzzing strategy" default:"random"`
		ListStrategies bool     `long:"list-strategies" description:"List all available strategies"`
	} `command:"fuzz"`

	Reduce struct {
		InputFile flags.Filename `long:"input-file" description:"Input file which gets parsed, validated and delta-debugged via the format file" required:"true"`
	} `command:"reduce"`

	Validate struct {
		InputFile flags.Filename `long:"input-file" description:"Input file which gets parsed and validated via the format file" required:"true"`
	} `command:"validate"`
}

type Strategy string

func (s *Strategy) Complete(match string) []flags.Completion {
	var items []flags.Completion

	for _, name := range strategy.List() {
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
		panic(err)
	}

	_, err := p.Parse()
	if opts.General.Help || len(os.Args) == 1 {
		p.WriteHelp(os.Stdout)

		os.Exit(returnHelp)
	} else if opts.General.Version {
		fmt.Printf("Tavor v%s\n", tavor.Version)

		os.Exit(returnOk)
	} else if opts.Fuzz.ListStrategies {
		for _, name := range strategy.List() {
			fmt.Println(name)
		}

		os.Exit(returnOk)
	}

	if err != nil {
		panic(err)
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

func main() {
	command := checkArguments()

	log.Infof("Open file %s", opts.Format.FormatFile)

	file, err := os.Open(string(opts.Format.FormatFile))
	if err != nil {
		panic(fmt.Errorf("cannot open tavor file %s: %v", opts.Format.FormatFile, err))
	}
	defer file.Close()

	doc, err := parser.ParseTavor(file)
	if err != nil {
		panic(fmt.Errorf("cannot parse tavor file: %v", err))
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

		strat, err := strategy.New(string(opts.Fuzz.Strategy), doc)
		if err != nil {
			panic(err)
		}

		log.Infof("Using %s strategy", opts.Fuzz.Strategy)

		ch, err := strat.Fuzz(r)
		if err != nil {
			panic(err)
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
	case "reduce", "validate":
		inputFile := opts.Validate.InputFile

		if command == "reduce" {
			inputFile = opts.Reduce.InputFile
		}

		input, err := os.Open(string(inputFile))
		if err != nil {
			panic(fmt.Errorf("cannot open input file %s: %v", inputFile, err))
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
			panic("TODO not implemented yet")
		}
	default:
		log.Panicf("Unknown command %q", command)
	}

	os.Exit(returnOk)
}
