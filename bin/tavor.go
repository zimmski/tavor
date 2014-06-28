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
)

var opts struct {
	Config         func(s flags.Filename) error `long:"config" description:"INI config file" no-ini:"true"`
	ConfigWrite    flags.Filename               `long:"config-write" description:"Write all arguments to an INI config file or to STDOUT with \"-\" as argument." no-ini:"true"`
	Debug          bool                         `long:"debug" description:"Temporary debugging argument"`
	FormatFile     flags.Filename               `long:"format-file" description:"Input tavor format file" required:"true" no-ini:"true"`
	Help           bool                         `long:"help" description:"Show this help message" no-ini:"true"`
	ListStrategies bool                         `long:"list-strategies" description:"List all available strategies." no-ini:"true"`
	PrintInternal  bool                         `long:"print-internal" description:"Prints the internal AST of the parsed file"`
	Seed           int64                        `long:"seed" description:"Seed for all the randomness"`
	Strategy       Strategy                     `long:"strategy" description:"The fuzzing strategy" default:"random"`
	Validate       bool                         `long:"validate" description:"Just validates the input file"`
	Verbose        bool                         `long:"verbose" description:"Do verbose output."`
	Version        bool                         `long:"version" description:"Print the version of this program." no-ini:"true"`

	configFile string
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

func checkArguments() {
	p := flags.NewNamedParser("tavor", flags.PassDoubleDash)

	p.ShortDescription = "A fuzzing and delta-debugging platform."

	opts.Config = func(s flags.Filename) error {
		ini := flags.NewIniParser(p)

		opts.configFile = string(s)

		return ini.ParseFile(string(s))
	}

	p.AddGroup("Tavor", "Tavor arguments", &opts)

	_, err := p.Parse()
	if opts.Help || len(os.Args) == 1 {
		p.WriteHelp(os.Stdout)

		os.Exit(returnHelp)
	} else if opts.Version {
		fmt.Printf("Tavor v%s\n", tavor.Version)

		os.Exit(returnOk)
	} else if opts.ListStrategies {
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

	if opts.ConfigWrite != "" {
		ini := flags.NewIniParser(p)

		var iniOptions flags.IniOptions = flags.IniIncludeComments | flags.IniIncludeDefaults | flags.IniCommentDefaults

		if opts.ConfigWrite == "-" {
			(ini.Write(os.Stdout, iniOptions))
		} else {
			ini.WriteFile(string(opts.ConfigWrite), iniOptions)
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

func main() {
	checkArguments()

	log.Infof("Open file %s", opts.FormatFile)

	file, err := os.Open(string(opts.FormatFile))
	if err != nil {
		panic(fmt.Errorf("cannot open tavor file %s: %v", opts.FormatFile, err))
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

	strat, err := strategy.New(string(opts.Strategy), doc)
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
