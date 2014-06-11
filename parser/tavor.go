package parser

import (
	"fmt"
	"io"
	"reflect"
	"strconv"
	"text/scanner"

	"github.com/zimmski/container/list/linkedlist"

	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/aggregates"
	"github.com/zimmski/tavor/token/constraints"
	"github.com/zimmski/tavor/token/expressions"
	"github.com/zimmski/tavor/token/lists"
	"github.com/zimmski/tavor/token/primitives"
	"github.com/zimmski/tavor/token/sequences"
)

/*
	TODO

	Token names can only consist of letters, digits and "_"

	ShortAlternation = [123]

*/

//TODO remove this
var DEBUG = false

const zeroRune = 0

const (
	MaxRepeat = 2
)

type tavorParser struct {
	scan scanner.Scanner

	err string

	earlyUse             map[string]token.Token
	embeddedTokensInTerm map[string][]map[string]struct{}
	lookup               map[string]token.Token
	used                 map[string]struct{}
}

func (p *tavorParser) expectRune(expect rune, got rune) (rune, error) {
	if got != expect {
		return got, &ParserError{
			Message: fmt.Sprintf("Expected \"%s\" but got \"%s\"", scanner.TokenString(expect), scanner.TokenString(got)),
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
		case '$':
			c, err = p.parseSpecialTokenDefinition()
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

func (p *tavorParser) parseTerm(c rune) (rune, []token.Token, []map[string]struct{}, error) {
	var err error
	var embeddedTokens = make([]map[string]struct{}, 0)
	var embeddedToks = make(map[string]struct{}, 0)
	var tokens []token.Token

OUT:
	for {
		switch c {
		case scanner.Ident:
			n := p.scan.TokenText()

			if _, ok := p.lookup[n]; !ok {
				if DEBUG {
					fmt.Printf("parseTerm use empty pointer for %s\n", n)
				}

				var tokenInterface *token.Token

				p.lookup[n] = primitives.NewEmptyPointer(tokenInterface)
				p.earlyUse[n] = p.lookup[n]
			}

			embeddedToks[n] = struct{}{}
			p.used[n] = struct{}{}

			tokens = append(tokens, p.lookup[n])
		case scanner.Int:
			v, _ := strconv.Atoi(p.scan.TokenText())

			tokens = append(tokens, primitives.NewConstantInt(v))
		case scanner.String:
			s := p.scan.TokenText()

			if s[0] != '"' {
				panic("unknown " + s) // TODO remove this
			}

			if s[len(s)-1] != '"' {
				return zeroRune, nil, nil, &ParserError{
					Message: "String is not terminated",
					Type:    ParseErrorNonTerminatedString,
				}
			}

			s, _ = strconv.Unquote(s)

			tokens = append(tokens, primitives.NewConstantString(s))
		case '(':
			if DEBUG {
				fmt.Println("NEW group")
			}
			c = p.scan.Scan()
			if DEBUG {
				fmt.Printf("parseTerm Group %d:%v -> %v\n", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())
			}

			c, toks, embeddedToks, err := p.parseScope(c)
			if err != nil {
				return zeroRune, nil, nil, err
			}

			p.expectRune(')', c)

			switch len(toks) {
			case 0:
				// ignore
			case 1:
				tokens = append(tokens, toks[0])
			default:
				tokens = append(tokens, lists.NewAll(toks...))
			}

			if len(embeddedToks) != 0 {
				embeddedTokens = append(embeddedTokens, embeddedToks...)
			}

			if DEBUG {
				fmt.Println("END group")
			}
		case '?':
			if DEBUG {
				fmt.Println("NEW optional")
			}
			p.expectScanRune('(')

			c = p.scan.Scan()
			if DEBUG {
				fmt.Printf("parseTerm optional after ( %d:%v -> %v\n", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())
			}

			c, toks, _, err := p.parseScope(c)
			if err != nil {
				return zeroRune, nil, nil, err
			}

			p.expectRune(')', c)

			switch len(toks) {
			case 0:
				// ignore
			case 1:
				tokens = append(tokens, constraints.NewOptional(toks[0]))
			default:
				tokens = append(tokens, constraints.NewOptional(lists.NewAll(toks...)))
			}

			if DEBUG {
				fmt.Println("END optional")
			}
		case '+', '*':
			if DEBUG {
				fmt.Println("NEW repeat")
			}

			sym := c

			c = p.scan.Scan()
			if DEBUG {
				fmt.Printf("parseTerm repeat before ( %d:%v -> %v\n", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())
			}

			var from, to int

			if sym == '*' {
				from, to = 0, MaxRepeat
			} else {
				if c == scanner.Int {
					from, _ = strconv.Atoi(p.scan.TokenText())

					c = p.scan.Scan()
					if DEBUG {
						fmt.Printf("parseTerm repeat after from ( %d:%v -> %v\n", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())
					}

					// until there is an explicit "to" we can assume to==from
					to = from
				} else {
					from, to = 1, MaxRepeat
				}

				if c == ',' {
					c = p.scan.Scan()
					if DEBUG {
						fmt.Printf("parseTerm repeat after , ( %d:%v -> %v\n", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())
					}

					if c == scanner.Int {
						to, _ = strconv.Atoi(p.scan.TokenText())

						c = p.scan.Scan()
						if DEBUG {
							fmt.Printf("parseTerm repeat after to ( %d:%v -> %v\n", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())
						}
					} else {
						to = MaxRepeat
					}
				}
			}

			p.expectRune('(', c)

			if DEBUG {
				fmt.Printf("repeat from %v to %v\n", from, to)
			}

			c = p.scan.Scan()
			if DEBUG {
				fmt.Printf("parseTerm repeat after ( %d:%v -> %v\n", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())
			}

			c, toks, embeddedToks, err := p.parseScope(c)
			if err != nil {
				return zeroRune, nil, nil, err
			}

			p.expectRune(')', c)

			switch len(toks) {
			case 0:
				// ignore
			case 1:
				tokens = append(tokens, lists.NewRepeat(toks[0], int64(from), int64(to)))
			default:
				tokens = append(tokens, lists.NewRepeat(lists.NewAll(toks...), int64(from), int64(to)))
			}

			if from > 0 && len(embeddedToks) != 0 {
				embeddedTokens = append(embeddedTokens, embeddedToks...)
			}

			if DEBUG {
				fmt.Println("END repeat")
			}
		case '$':
			c = p.scan.Scan()
			if DEBUG {
				fmt.Printf("parseTerm after $ ( %d:%v -> %v\n", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())
			}

			var tok token.Token

			if c == '{' {
				tok, err = p.parseExpression(c)
			} else {
				tok, err = p.parseTokenAttribute(c)
			}

			if err != nil {
				return zeroRune, nil, nil, err
			}

			tokens = append(tokens, tok)
		case ',': // multi line token
			if _, err := p.expectScanRune('\n'); err != nil {
				return zeroRune, nil, nil, err
			}

			c = p.scan.Scan()
			if DEBUG {
				fmt.Printf("parseTerm multiline %d:%v -> %v\n", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())
			}

			if c == '\n' {
				return zeroRune, nil, nil, &ParserError{
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

	if len(embeddedToks) != 0 {
		embeddedTokens = append(embeddedTokens, embeddedToks)
	}

	return c, tokens, embeddedTokens, nil
}

func (p *tavorParser) parseExpression(c rune) (token.Token, error) {
	if DEBUG {
		fmt.Println("START expression")
	}

	_, err := p.expectRune('{', c)
	if err != nil {
		return nil, err
	}

	c = p.scan.Scan()
	if DEBUG {
		fmt.Printf("parseExpression after {} %d:%v -> %v\n", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())
	}

	c, tok, err := p.parseExpressionTerm(c)
	if err != nil {
		return nil, err
	} else if tok == nil {
		return nil, &ParserError{
			Message: "Empty expressions are not allowed",
			Type:    ParseErrorEmptyExpressionIsInvalid,
		}
	}

	_, err = p.expectRune('}', c)
	if err != nil {
		return nil, err
	}

	if DEBUG {
		fmt.Println("END expression")
	}

	return tok, nil
}

func (p *tavorParser) parseExpressionTerm(c rune) (rune, token.Token, error) {
	var tok token.Token
	var err error

	// single term
	switch c {
	case scanner.Ident:
		tok, err = p.parseTokenAttribute(c)
		if err != nil {
			return zeroRune, nil, err
		}
	case scanner.Int:
		v, _ := strconv.Atoi(p.scan.TokenText())

		tok = primitives.NewConstantInt(v)
	}

	if tok == nil {
		return zeroRune, nil, nil
	}

	c = p.scan.Scan()
	if DEBUG {
		fmt.Printf("parseExpressionTerm %d:%v -> %v\n", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())
	}

	// operators
	switch c {
	case '+', '-', '*', '/':
		sym := c

		c = p.scan.Scan()
		if DEBUG {
			fmt.Printf("parseExpressionTerm operator %d:%v -> %v\n", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())
		}

		var t token.Token
		c, t, err = p.parseExpressionTerm(c)
		if err != nil {
			return zeroRune, nil, err
		} else if t == nil {
			return zeroRune, nil, &ParserError{
				Message: "Expected another expression term after operator",
				Type:    ParseErrorExpectedExpressionTerm,
			}
		}

		switch sym {
		case '+':
			tok = expressions.NewAddArithmetic(tok, t)
		case '-':
			tok = expressions.NewSubArithmetic(tok, t)
		case '*':
			tok = expressions.NewMulArithmetic(tok, t)
		case '/':
			tok = expressions.NewDivArithmetic(tok, t)
		}
	}

	return c, tok, nil
}

func (p *tavorParser) parseTokenAttribute(c rune) (token.Token, error) {
	if DEBUG {
		fmt.Println("START token attribute")
	}

	_, err := p.expectRune(scanner.Ident, c)
	if err != nil {
		return nil, err
	}

	name := p.scan.TokenText()

	_, err = p.expectScanRune('.')
	if err != nil {
		return nil, err
	}

	_, err = p.expectScanRune(scanner.Ident)
	if err != nil {
		return nil, err
	}

	attribute := p.scan.TokenText()

	tok, ok := p.lookup[name]
	if !ok {
		return nil, &ParserError{
			Message: fmt.Sprintf("Token \"%s\" is not defined", name),
			Type:    ParseErrorTokenNotDefined,
		}
	}

	p.used[name] = struct{}{}

	if DEBUG {
		fmt.Println("END token attribute (or will be unknown token attribute)")
	}

	switch i := tok.(type) {
	case lists.List:
		switch attribute {
		case "Count":
			return aggregates.NewLen(i), nil
		}
	case *sequences.Sequence:
		switch attribute {
		case "Existing":
			return i.ExistingItem(), nil
		case "Next":
			return i.Item(), nil
		case "Reset":
			return i.ResetItem(), nil
		}
	}

	return nil, &ParserError{
		Message: fmt.Sprintf("Unknown token attribute \"%s\" for token type \"%s\"", attribute, reflect.TypeOf(tok)),
		Type:    ParseErrorUnknownTokenAttribute,
	}
}

func (p *tavorParser) parseScope(c rune) (rune, []token.Token, []map[string]struct{}, error) {
	var err error
	var embeddedTokens = make([]map[string]struct{}, 0)
	var tokens []token.Token

	var toks []token.Token
	var embeddedToks []map[string]struct{}

	c, toks, embeddedToks, err = p.parseTerm(c)
	if err != nil {
		return zeroRune, nil, nil, err
	} else if toks != nil {
		tokens = toks
	}

	if c == '|' {
		if DEBUG {
			fmt.Println("NEW or")
		}
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

			if embeddedTokens != nil {
				if len(embeddedToks) == 0 {
					// since there is a Or term without any token embedded we can say
					// that we can break out of a loop at this point
					embeddedTokens = nil
				} else {
					embeddedTokens = append(embeddedTokens, embeddedToks...)
				}
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

			c, toks, embeddedToks, err = p.parseTerm(c)
			if err != nil {
				return zeroRune, nil, nil, err
			}
		}

		or := lists.NewOne(orTerms...)

		if optional {
			tokens = []token.Token{constraints.NewOptional(or)}

			embeddedTokens = nil
		} else {
			tokens = []token.Token{or}
		}

		if DEBUG {
			fmt.Println("END or")
		}
	} else {
		if len(embeddedToks) != 0 {
			embeddedTokens = append(embeddedTokens, embeddedToks...)
		}
	}

	return c, tokens, embeddedTokens, nil
}

func (p *tavorParser) parseTokenDefinition() (rune, error) {
	var c rune
	var err error

	name := p.scan.TokenText()

	if tok, ok := p.lookup[name]; ok {
		// if there is a pointer in the lookup hash we can say that it was just used before
		if _, ok := tok.(*primitives.Pointer); !ok {
			return zeroRune, &ParserError{
				Message: "Token already defined",
				Type:    ParseErrorTokenAlreadyDefined,
			}
		}
	}

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

	c, tokens, embeddedToks, err := p.parseScope(c)
	if err != nil {
		return zeroRune, err
	}

	p.embeddedTokensInTerm[name] = embeddedToks

	if DEBUG {
		fmt.Printf("Token %s embeds %+v\n", name, embeddedToks)
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

	var tok token.Token

	switch len(tokens) {
	case 0:
		return zeroRune, &ParserError{
			Message: "Empty token definition",
			Type:    ParseErrorEmptyTokenDefinition,
		}
	case 1:
		tok = tokens[0]
	default:
		tok = lists.NewAll(tokens...)
	}

	// self loop?
	if pointer, ok := p.lookup[name]; ok {
		if DEBUG {
			fmt.Printf("parseTokenDefinition fill empty pointer for %s\n", name)
		}

		err = pointer.(*primitives.Pointer).Set(tok)
		if err != nil {
			return zeroRune, &ParserError{
				Message: fmt.Sprintf("Wrong token type for %s because of earlier usage: %s", name, err),
				Type:    ParseErrorInvalidTokenType,
			}
		}
	}

	// check for endless loop
	if len(embeddedToks) != 0 {
		foundExit := false

		if DEBUG {
			fmt.Println("Need to check for loops")
		}

	EMBEDDEDTOKS:
		for _, toks := range embeddedToks {
			checked := make(map[string]struct{})
			l := linkedlist.New()

			for n := range toks {
				l.Push(n)
			}

			for !l.Empty() {
				i, _ := l.Shift()
				n := i.(string)

				if name == n {
					if DEBUG {
						fmt.Println("Found one loop")
					}

					continue EMBEDDEDTOKS
				}

				checked[n] = struct{}{}

				for _, toks := range p.embeddedTokensInTerm[n] {
					for n := range toks {
						if _, ok := checked[n]; !ok {
							l.Push(n)
						}
					}
				}
			}
			if DEBUG {
				fmt.Println("Found exit, everything is fine.")
			}

			foundExit = true

			break
		}

		if !foundExit {
			if DEBUG {
				fmt.Println("There is no loop exit for this token, I'll through an error.")
			}

			return zeroRune, &ParserError{
				Message: fmt.Sprintf("Token \"%s\" has an endless loop without exit\n", name),
				Type:    ParseErrorEndlessLoopDetected,
			}
		}
	}

	p.lookup[name] = tok

	c = p.scan.Scan()
	if DEBUG {
		fmt.Printf("parseTokenDefinition after newline %d:%v -> %v\n", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())
	}

	return c, nil
}

func (p *tavorParser) parseSpecialTokenDefinition() (rune, error) {
	var c rune
	var err error

	if DEBUG {
		fmt.Println("START special token")
	}

	c = p.scan.Scan()
	if DEBUG {
		fmt.Printf("parseSpecialTokenDefinition after $ %d:%v -> %v\n", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())
	}

	name := p.scan.TokenText()
	if _, ok := p.lookup[name]; ok {
		return zeroRune, &ParserError{
			Message: "Token already defined",
			Type:    ParseErrorTokenAlreadyDefined,
		}
	}

	if c, err = p.expectScanRune('='); err != nil {
		return zeroRune, err
	}

	arguments := make(map[string]string)

	for {
		c, err = p.expectScanRune(scanner.Ident)
		if err != nil {
			return zeroRune, err
		}

		arg := p.scan.TokenText()

		_, err = p.expectScanRune(':')
		if err != nil {
			return zeroRune, err
		}

		c = p.scan.Scan()
		if DEBUG {
			fmt.Printf("parseSpecialTokenDefinition argument value %d:%v -> %v\n", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())
		}

		switch c {
		case scanner.Ident, scanner.String, scanner.Int:
			arguments[arg] = p.scan.TokenText()
		default:
			return zeroRune, &ParserError{
				Message: fmt.Sprintf("Invalid argument value %v", c),
				Type:    ParseErrorInvalidArgumentValue,
			}
		}

		c = p.scan.Scan()
		if DEBUG {
			fmt.Printf("parseSpecialTokenDefinition after argument value %d:%v -> %v\n", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())
		}

		if c != ',' {
			break
		}

		if c, err = p.expectScanRune('\n'); err != nil {
			return zeroRune, err
		}
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

	typ, ok := arguments["type"]
	if !ok {
		return zeroRune, &ParserError{
			Message: "Special token has no type argument",
			Type:    ParseErrorUnknownTypeForSpecialToken,
		}
	}

	var tok token.Token
	usedArguments := map[string]struct{}{
		"type": struct{}{},
	}

	switch typ {
	case "Int":
		rawFrom, okFrom := arguments["from"]
		rawTo, okTo := arguments["to"]

		if okFrom || okTo {
			if okFrom && !okTo {
				return zeroRune, &ParserError{
					Message: "Argument \"to\" is missing",
					Type:    ParseErrorMissingSpecialTokenArgument,
				}
			} else if !okFrom && okTo {
				return zeroRune, &ParserError{
					Message: "Argument \"from\" is missing",
					Type:    ParseErrorMissingSpecialTokenArgument,
				}
			}

			from, err := strconv.Atoi(rawFrom)
			if err != nil {
				return zeroRune, &ParserError{
					Message: "\"from\" needs an integer value",
					Type:    ParseErrorInvalidArgumentValue,
				}
			}

			to, err := strconv.Atoi(rawTo)
			if err != nil {
				return zeroRune, &ParserError{
					Message: "\"to\" needs an integer value",
					Type:    ParseErrorInvalidArgumentValue,
				}
			}

			usedArguments["from"] = struct{}{}
			usedArguments["to"] = struct{}{}

			tok = primitives.NewRangeInt(from, to)
		} else {
			tok = primitives.NewRandomInt()
		}
	case "Sequence":
		start := 1
		step := 1

		if raw, ok := arguments["start"]; ok {
			start, err = strconv.Atoi(raw)
			if err != nil {
				return zeroRune, &ParserError{
					Message: "\"start\" needs an integer value",
					Type:    ParseErrorInvalidArgumentValue,
				}
			}
		}

		if raw, ok := arguments["step"]; ok {
			step, err = strconv.Atoi(raw)
			if err != nil {
				return zeroRune, &ParserError{
					Message: "\"step\" needs an integer value",
					Type:    ParseErrorInvalidArgumentValue,
				}
			}
		}

		usedArguments["start"] = struct{}{}
		usedArguments["step"] = struct{}{}

		tok = sequences.NewSequence(start, step)
	default:
		return zeroRune, &ParserError{
			Message: fmt.Sprintf("Unknown special token type \"%s\"", typ),
			Type:    ParseErrorUnknownSpecialTokenType,
		}
	}

	for arg := range arguments {
		if _, ok := usedArguments[arg]; !ok {
			return zeroRune, &ParserError{
				Message: fmt.Sprintf("Unknown special token argument \"%s\"", arg),
				Type:    ParseErrorUnknownSpecialTokenArgument,
			}
		}
	}

	p.lookup[name] = tok

	c = p.scan.Scan()
	if DEBUG {
		fmt.Printf("parseSpecialTokenDefinition after newline %d:%v -> %v\n", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())
	}

	if DEBUG {
		fmt.Println("END special token")
	}

	return c, nil
}

func ParseTavor(src io.Reader) (token.Token, error) {
	p := &tavorParser{

		earlyUse:             make(map[string]token.Token),
		embeddedTokensInTerm: make(map[string][]map[string]struct{}),
		lookup:               make(map[string]token.Token),
		used:                 make(map[string]struct{}),
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

	for name, tok := range p.earlyUse {
		if tok.(*primitives.Pointer).Get() == nil {
			return nil, &ParserError{
				Message: fmt.Sprintf("Token \"%s\" is not defined", name),
				Type:    ParseErrorTokenNotDefined,
			}
		}
	}

	for name := range p.lookup {
		if _, ok := p.used[name]; !ok {
			return nil, &ParserError{
				Message: fmt.Sprintf("Token \"%s\" declared but not used", name),
				Type:    ParseErrorUnusedToken,
			}
		}
	}

	return p.lookup["START"], nil
}
