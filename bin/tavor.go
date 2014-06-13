package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/jessevdk/go-flags"

	"github.com/zimmski/tavor"
	"github.com/zimmski/tavor/fuzz/strategy"
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

	InputFile string `long:"input-file" description:"Input tavor file" required:"true" no-ini:"true"`
	Seed      int64  `long:"seed" description:"Seed for all the randomness"`
	Strategy  string `long:"strategy" description:"The fuzzing strategy" default:"random"`
	Verbose   bool   `long:"verbose" description:"Do verbose output."`
	Version   bool   `long:"version" description:"Print the version of this program." no-ini:"true"`

	Debug bool `long:"debug" description:"Temporary debugging argument"`

	configFile string
}

func V(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "[VERBOSE] "+msg+"\n", args...)
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

	if _, err := p.ParseArgs(os.Args); err != nil {
		if opts.Version {
			fmt.Printf("Tavor v%s\n", tavor.Version)

			os.Exit(returnOk)
		} else if opts.ListStrategies {
			for _, name := range strategy.List() {
				fmt.Println(name)
			}

			os.Exit(returnOk)
		}

		if e, ok := err.(*flags.Error); !ok || e.Type != flags.ErrHelp {
			panic(err)
		} else {
			p.WriteHelp(os.Stdout)

			os.Exit(returnHelp)
		}
	}

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
}

func main() {
	checkArguments()

	if opts.Verbose {
		V("Open file %s", opts.InputFile)
	}

	file, err := os.Open(opts.InputFile)
	if err != nil {
		panic(fmt.Errorf("cannot open tavor file %s: %v", opts.InputFile, err))
	}
	defer file.Close()

	if opts.Debug {
		tavor.DEBUG = true
	}

	doc, err := parser.ParseTavor(file)
	if err != nil {
		panic(fmt.Errorf("cannot parse tavor file: %v", err))
	}

	if opts.Verbose {
		V("Using seed %d", opts.Seed)
	}

	if opts.Verbose {
		V("Counted %d permutations", doc.Permutations())
	}

	r := rand.New(rand.NewSource(opts.Seed))

	strat, err := strategy.New(opts.Strategy, doc)
	if err != nil {
		panic(err)
	}

	if opts.Verbose {
		V("Using %s strategy", opts.Strategy)
	}

	strat.Fuzz(r)

	fmt.Print(doc.String())
}
