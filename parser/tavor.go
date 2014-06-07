package parser

import (
	"fmt"
	"io"
	"strconv"
	"text/scanner"

	"github.com/zimmski/tavor/token"
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

func (p *tavorParser) parseTerm(c rune) (token.Token, error) {
	switch c {
	case scanner.Ident:
		n := p.scan.TokenText()

		if _, ok := p.lookup[n]; !ok {
			return nil, &ParserError{
				Message: fmt.Sprintf("Token %s does not exists", n),
				Type:    ParseErrorTokenDoesNotExists,
			}
		}

		p.used[n] = struct{}{}

		return p.lookup[n].Clone(), nil
	case scanner.Int:
		v, _ := strconv.Atoi(p.scan.TokenText())

		return primitives.NewConstantInt(v), nil
	case scanner.String:
		s := p.scan.TokenText()

		if s[0] != '"' {
			panic("unknown " + s) // TODO remove this
		}

		if s[len(s)-1] != '"' {
			return nil, &ParserError{
				Message: "String is not terminated",
				Type:    ParseErrorNonTerminatedString,
			}
		}

		return primitives.NewConstantString(s[1 : len(s)-1]), nil
	}

	return nil, nil
}

func (p *tavorParser) parseTokenDefinition() (rune, error) {
	name := p.scan.TokenText()

	if _, ok := p.lookup[name]; ok {
		return zeroRune, &ParserError{
			Message: "Token already exists",
			Type:    ParseErrorTokenExists,
		}
	}

	// do an empty definition to allow loops
	p.lookup[name] = nil

	if c, err := p.expectScanRune('='); err != nil {
		// unexpected new line?
		if c == '\n' {
			return zeroRune, &ParserError{
				Message: "New line inside single line token definitions is not allowed",
				Type:    ParseErrorEarlyNewLine,
			}
		}

		return zeroRune, err
	}

	tokens := make([]token.Token, 0)

	c := p.scan.Scan()
	if DEBUG {
		fmt.Printf("%d:%v -> %v\n", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())
	}

	for {
		tok, err := p.parseTerm(c)
		if err != nil {
			return zeroRune, err
		} else if tok != nil {
			tokens = append(tokens, tok)
		} else {
			switch c {
			case scanner.EOF:
				return zeroRune, &ParserError{
					Message: "New line at end of token definition needed",
					Type:    ParseErrorNewLineNeeded,
				}
			case ',': // multi line token
				if _, err := p.expectScanRune('\n'); err != nil {
					return zeroRune, err
				}

				c = p.scan.Scan()
				if DEBUG {
					fmt.Printf("%d:%v -> %v\n", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())
				}

				if c == '\n' {
					return zeroRune, &ParserError{
						Message: "Multi line token definition unexpectedly terminated",
						Type:    ParseErrorUnexpectedTokenDefinitionTermination,
					}
				}

				continue
			case '\n':
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
					fmt.Printf("%d:%v -> %v\n", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())
				}

				return c, nil
			default:
				panic("now what?") // TODO remove this
			}
		}

		c = p.scan.Scan()
		if DEBUG {
			fmt.Printf("aa%d:%v -> %v\n", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())
		}
	}

	return c, nil
}

func ParseTavor(src io.Reader) (token.Token, error) {
	var err error

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
				return nil, err
			}

			continue
		case scanner.Int:
			return nil, &ParserError{
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
