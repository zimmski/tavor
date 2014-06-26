package main

import (
	"fmt"
	"math/rand"
	"os"
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
)

var opts struct {
	Config      func(s string) error `long:"config" description:"INI config file" no-ini:"true"`
	ConfigWrite string               `long:"config-write" description:"Write all arguments to an INI config file or to STDOUT with \"-\" as argument." no-ini:"true"`

	ListStrategies bool `long:"list-strategies" description:"List all available strategies." no-ini:"true"`

	InputFile     string `long:"input-file" description:"Input tavor file" required:"true" no-ini:"true"`
	PrintInternal bool   `long:"print-internal" description:"Prints the interal AST of the parsed file"`
	Seed          int64  `long:"seed" description:"Seed for all the randomness"`
	Strategy      string `long:"strategy" description:"The fuzzing strategy" default:"random"`
	Validate      bool   `long:"validate" description:"Just validates the input file"`
	Verbose       bool   `long:"verbose" description:"Do verbose output."`
	Version       bool   `long:"version" description:"Print the version of this program." no-ini:"true"`

	Debug bool `long:"debug" description:"Temporary debugging argument"`

	configFile string
}

func checkArguments() {
	p := flags.NewNamedParser("tavor", flags.HelpFlag)
	p.ShortDescription = "A fuzzing and delta-debugging platform."

	opts.Config = func(s string) error {
		ini := flags.NewIniParser(p)

		opts.configFile = s

		return ini.ParseFile(s)
	}

	p.AddGroup("Tavor", "Tavor arguments", &opts)

	if len(os.Args) == 1 {
		p.WriteHelp(os.Stdout)

		os.Exit(returnHelp)
	}

	if _, err := p.ParseArgs(os.Args); err != nil {
		doListArguments()

		if e, ok := err.(*flags.Error); !ok || e.Type != flags.ErrHelp {
			panic(err)
		} else {
			p.WriteHelp(os.Stdout)

			os.Exit(returnHelp)
		}
	}

	doListArguments()

	if opts.ConfigWrite != "" {
		ini := flags.NewIniParser(p)

		var iniOptions flags.IniOptions = flags.IniIncludeComments | flags.IniIncludeDefaults | flags.IniCommentDefaults

		if opts.ConfigWrite == "-" {
			(ini.Write(os.Stdout, iniOptions))
		} else {
			ini.WriteFile(opts.ConfigWrite, iniOptions)
		}

		os.Exit(returnOk)
	}

	if opts.Seed == 0 {
		opts.Seed = time.Now().UTC().UnixNano()
	}

	if opts.Debug {
		log.LevelDebug()
	} else if opts.Verbose {
		log.LevelInfo()
	} else {
		log.LevelWarn()
	}
}

func doListArguments() {
	if opts.Version {
		fmt.Printf("Tavor v%s\n", tavor.Version)

		os.Exit(returnOk)
	} else if opts.ListStrategies {
		for _, name := range strategy.List() {
			fmt.Println(name)
		}

		os.Exit(returnOk)
	}
}

func main() {
	checkArguments()

	log.Infof("Open file %s", opts.InputFile)

	file, err := os.Open(opts.InputFile)
	if err != nil {
		panic(fmt.Errorf("cannot open tavor file %s: %v", opts.InputFile, err))
	}
	defer file.Close()

	doc, err := parser.ParseTavor(file)
	if err != nil {
		panic(fmt.Errorf("cannot parse tavor file: %v", err))
	}

	log.Info("File is valid")

	if opts.PrintInternal {
		tavor.PrettyPrintInternalTree(os.Stdout, doc)
	}

	if opts.Validate {
		os.Exit(returnOk)
	}

	log.Infof("Using seed %d", opts.Seed)
	log.Infof("Counted %d overall permutations", doc.PermutationsAll())

	r := rand.New(rand.NewSource(opts.Seed))

	strat, err := strategy.New(opts.Strategy, doc)
	if err != nil {
		panic(err)
	}

	log.Infof("Using %s strategy", opts.Strategy)

	ch, err := strat.Fuzz(r)
	if err != nil {
		panic(err)
	}

	another := false
	for i := range ch {
		if !opts.Debug {
			if another {
				fmt.Println()
			} else {
				another = true
			}
		}

		log.Debug("Result:")

		fmt.Print(doc.String())
		if opts.Debug {
			fmt.Println()
		}

		ch <- i
	}
}
