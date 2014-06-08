package parser

import (
	"fmt"
	"io"
	"strconv"
	"text/scanner"

	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/constraints"
	"github.com/zimmski/tavor/token/lists"
	"github.com/zimmski/tavor/token/primitives"
)

//TODO remove this
var DEBUG = false

const zeroRune = 0

type tavorParser struct {
	scan scanner.Scanner

	err string

	lookup map[string]token.Token
	used   map[string]struct{}
}

func (p *tavorParser) expectRune(expect rune, got rune) (rune, error) {
	if got != expect {
		return got, &ParserError{
			Message: fmt.Sprintf("Expected \"%c\" but got \"%c\"", expect, got),
			Type:    ParseErrorExpectRune,
		}
	}

	return got, nil
}

func (p *tavorParser) expectScanRune(expect rune) (rune, error) {
	got := p.scan.Scan()
	if DEBUG {
		fmt.Printf("%d:%v -> %v\n", p.scan.Line, scanner.TokenString(got), p.scan.TokenText())
	}

	return p.expectRune(expect, got)
}

func (p *tavorParser) parseGlobalScope() error {
	var err error

	c := p.scan.Scan()
	if DEBUG {
		fmt.Printf("%d:%v -> %v\n", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())
	}

	for c != scanner.EOF {
		switch c {
		case '\n':
			// ignore new lines in the global scope
		case scanner.Ident:
			c, err = p.parseTokenDefinition()
			if err != nil {
				return err
			}

			continue
		case scanner.Int:
			return &ParserError{
				Message: "Token names have to start with a letter",
				Type:    ParseErrorInvalidTokenName,
			}
		default:
			panic("what am i to do now") // TODO remove this
		}

		c = p.scan.Scan()
		if DEBUG {
			fmt.Printf("%d:%v -> %v\n", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())
		}
	}

	return nil
}

func (p *tavorParser) parseTerm(c rune) (rune, []token.Token, error) {
	tokens := make([]token.Token, 0)

OUT:
	for {
		switch c {
		case scanner.Ident:
			n := p.scan.TokenText()

			if _, ok := p.lookup[n]; !ok {
				return zeroRune, nil, &ParserError{
					Message: fmt.Sprintf("Token %s does not exists", n),
					Type:    ParseErrorTokenDoesNotExists,
				}
			}

			p.used[n] = struct{}{}

			tokens = append(tokens, p.lookup[n].Clone())
		case scanner.Int:
			v, _ := strconv.Atoi(p.scan.TokenText())

			tokens = append(tokens, primitives.NewConstantInt(v))
		case scanner.String:
			s := p.scan.TokenText()

			if s[0] != '"' {
				panic("unknown " + s) // TODO remove this
			}

			if s[len(s)-1] != '"' {
				return zeroRune, nil, &ParserError{
					Message: "String is not terminated",
					Type:    ParseErrorNonTerminatedString,
				}
			}

			tokens = append(tokens, primitives.NewConstantString(s[1:len(s)-1]))
		case ',': // multi line token
			if _, err := p.expectScanRune('\n'); err != nil {
				return zeroRune, nil, err
			}

			c = p.scan.Scan()
			if DEBUG {
				fmt.Printf("parseTerm multiline %d:%v -> %v\n", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())
			}

			if c == '\n' {
				return zeroRune, nil, &ParserError{
					Message: "Multi line token definition unexpectedly terminated",
					Type:    ParseErrorUnexpectedTokenDefinitionTermination,
				}
			}

			continue
		default:
			if DEBUG {
				fmt.Println("break out parseTerm")
			}
			break OUT
		}

		c = p.scan.Scan()
		if DEBUG {
			fmt.Printf("parseTerm %d:%v -> %v\n", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())
		}
	}

	return c, tokens, nil
}

func (p *tavorParser) parseScope(c rune) (rune, []token.Token, error) {
	var err error

	tokens := make([]token.Token, 0)

OUT:
	for {
		// identifier and literals
		var toks []token.Token
		c, toks, err = p.parseTerm(c)
		if err != nil {
			return zeroRune, nil, err
		} else if toks != nil {
			tokens = append(tokens, toks...)
			if DEBUG {
				fmt.Println("add these tokens in parseScope")
			}
		}

		// alternations and groupings
		switch c {
		case '|':
			var orTerms []token.Token
			optional := false

			toks = tokens

		OR:
			for {
				switch len(toks) {
				case 0:
					optional = true
				case 1:
					orTerms = append(orTerms, toks[0])
				default:
					orTerms = append(orTerms, lists.NewAll(toks...))
				}

				if c == '|' {
					c = p.scan.Scan()
					if DEBUG {
						fmt.Printf("parseScope Or %d:%v -> %v\n", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())
					}
				} else {
					if DEBUG {
						fmt.Println("parseScope break out or")
					}
					break OR
				}

				c, toks, err = p.parseTerm(c)
				if err != nil {
					return zeroRune, nil, err
				}
			}

			or := lists.NewOne(orTerms...)

			if optional {
				tokens = []token.Token{constraints.NewOptional(or)}
			} else {
				tokens = []token.Token{or}
			}

			continue
		default:
			break OUT
		}

		c = p.scan.Scan()
		if DEBUG {
			fmt.Printf("parseScope %d:%v -> %v\n", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())
		}
	}

	return c, tokens, nil
}

func (p *tavorParser) parseTokenDefinition() (rune, error) {
	var c rune
	var err error

	name := p.scan.TokenText()

	if _, ok := p.lookup[name]; ok {
		return zeroRune, &ParserError{
			Message: "Token already exists",
			Type:    ParseErrorTokenExists,
		}
	}

	// do an empty definition to allow loops
	p.lookup[name] = nil

	if c, err = p.expectScanRune('='); err != nil {
		// unexpected new line?
		if c == '\n' {
			return zeroRune, &ParserError{
				Message: "New line inside single line token definitions is not allowed",
				Type:    ParseErrorEarlyNewLine,
			}
		}

		return zeroRune, err
	}

	c = p.scan.Scan()
	if DEBUG {
		fmt.Printf("parseTokenDefinition after = %d:%v -> %v\n", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())
	}

	c, tokens, err := p.parseScope(c)
	if err != nil {
		return zeroRune, err
	}

	if DEBUG {
		fmt.Printf("back to token definition with c=%c\n", c)
	}

	// we always want a new line at the end of the file
	if c == scanner.EOF {
		return zeroRune, &ParserError{
			Message: "New line at end of token definition needed",
			Type:    ParseErrorNewLineNeeded,
		}
	}

	if c, err = p.expectRune('\n', c); err != nil {
		return zeroRune, err
	}

	switch len(tokens) {
	case 0:
		return zeroRune, &ParserError{
			Message: "Empty token definition",
			Type:    ParseErrorEmptyTokenDefinition,
		}
	case 1:
		p.lookup[name] = tokens[0]
	default:
		p.lookup[name] = lists.NewAll(tokens...)
	}

	c = p.scan.Scan()
	if DEBUG {
		fmt.Printf("parseTokenDefinition after newline %d:%v -> %v\n", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())
	}

	return c, nil
}

func ParseTavor(src io.Reader) (token.Token, error) {
	p := &tavorParser{
		lookup: make(map[string]token.Token),
		used:   make(map[string]struct{}),
	}

	if DEBUG {
		fmt.Println("INIT")
	}

	p.scan.Init(src)

	p.scan.Error = func(s *scanner.Scanner, msg string) {
		p.err = msg
	}
	p.scan.Whitespace = 1<<'\t' | 1<<' ' | 1<<'\r'

	if err := p.parseGlobalScope(); err != nil {
		return nil, err
	}

	if _, ok := p.lookup["START"]; !ok {
		return nil, &ParserError{
			Message: "No START token defined",
			Type:    ParseErrorNoStart,
		}
	}

	p.used["START"] = struct{}{}

	for key := range p.lookup {
		if _, ok := p.used[key]; !ok {
			return nil, &ParserError{
				Message: fmt.Sprintf("Token %s declared but not used", key),
				Type:    ParseErrorUnusedToken,
			}
		}
	}

	return p.lookup["START"], nil
}
