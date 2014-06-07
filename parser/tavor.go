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

func (p *tavorParser) expectRune(r rune) (rune, error) {
	c := p.scan.Scan()
	if DEBUG {
		fmt.Printf("%d:%v -> %v\n", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())
	}

	if c != r {
		return c, fmt.Errorf("Expected \"%c\" but got \"%c\"", r, c)
	}

	return c, nil
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

	if c, err := p.expectRune('='); err != nil {
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
		switch c {
		case scanner.Int:
			v, _ := strconv.Atoi(p.scan.TokenText())

			tokens = append(tokens, primitives.NewConstantInt(v))
		case scanner.String:
			s := p.scan.TokenText()

			if s[0] == '"' {
				if s[len(s)-1] != '"' {
					return zeroRune, &ParserError{
						Message: "String is not terminated",
						Type:    ParseErrorNonTerminatedString,
					}
				}

				tokens = append(tokens, primitives.NewConstantString(s[1:len(s)-1]))
			} else {
				panic("unknown " + s) // TODO remove this
			}
		case scanner.EOF:
			return zeroRune, &ParserError{
				Message: "New line at end of token definition needed",
				Type:    ParseErrorNewLineNeeded,
			}
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

		c = p.scan.Scan()
		if DEBUG {
			fmt.Printf("%d:%v -> %v\n", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())
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
