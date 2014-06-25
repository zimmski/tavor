package parser

import (
	"fmt"
	"io"
	"reflect"
	"strconv"
	"text/scanner"

	"github.com/zimmski/container/list/linkedlist"

	"github.com/zimmski/tavor"
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
	Allow forward usage of token attributes
	Allow correct forward usage of ExistingSequenceItems
*/

const zeroRune = 0

const (
	MaxRepeat = 2
)

type tokenUse struct {
	token    token.Token
	position scanner.Position
}

type tavorParser struct {
	scan scanner.Scanner

	err string

	earlyUse             map[string]tokenUse
	embeddedTokensInTerm map[string][]map[string]struct{}
	lookup               map[string]tokenUse
	lookupUsage          map[token.Token]struct{}
	used                 map[string]scanner.Position
}

func (p *tavorParser) expectRune(expect rune, got rune) (rune, error) {
	if got != expect {
		return got, &ParserError{
			Message:  fmt.Sprintf("Expected \"%s\" but got \"%s\"", scanner.TokenString(expect), scanner.TokenString(got)),
			Type:     ParseErrorExpectRune,
			Position: p.scan.Pos(),
		}
	}

	return got, nil
}

func (p *tavorParser) expectScanRune(expect rune) (rune, error) {
	got := p.scan.Scan()
	if tavor.DEBUG {
		fmt.Printf("%d:%v -> %v\n", p.scan.Line, scanner.TokenString(got), p.scan.TokenText())
	}

	return p.expectRune(expect, got)
}

func (p *tavorParser) parseGlobalScope() error {
	var err error

	c := p.scan.Scan()
	if tavor.DEBUG {
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
		default:
			return &ParserError{
				Message:  fmt.Sprintf("Token names have to start with a letter and not with %s", scanner.TokenString(c)),
				Type:     ParseErrorInvalidTokenName,
				Position: p.scan.Pos(),
			}
		}

		c = p.scan.Scan()
		if tavor.DEBUG {
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
			name := p.scan.TokenText()

			_, ok := p.lookup[name]
			if !ok {
				if tavor.DEBUG {
					fmt.Printf("parseTerm use empty pointer for %s\n", name)
				}

				var tokenInterface *token.Token
				n := primitives.NewEmptyPointer(tokenInterface)

				p.lookup[name] = tokenUse{
					token:    n,
					position: p.scan.Position,
				}
				p.earlyUse[name] = tokenUse{
					token:    n,
					position: p.scan.Position,
				}
			}

			embeddedToks[name] = struct{}{}
			p.used[name] = p.scan.Position

			/*

				THIS IS A TERRIBLE HACK THIS SHOULD BE REMOVED ASAP

				but let me explain

				B = 1 | 2
				A = B B

				This should result in 4 permutations 11 21 12 and 22
				but without this condition this would result in only
				11s and 22s. The clone alone would be fine but this
				leads to more problems if the clone is not saved back
				into the lookup.

				For example

				Bs = +(1)
				A = $Bs.Count Bs

				would not work without saving back.

				But what if somebody writes

				Cs = +(1)
				B = $Cs.Count Cs
				A = +(B)

				or even

				Cs = +(1)
				B = $Cs.Count Cs
				A = $Cs.Count +(B)

				So TODO and FIXME all over this

			*/
			tok := p.lookup[name].token

			if _, ok := p.lookupUsage[tok]; ok {
				ntok := tok.Clone()

				if tavor.DEBUG {
					fmt.Printf("token %s %#v(%p) was already used once. Cloned as %#v(%p)\n", name, tok, tok, ntok, ntok)
				}

				p.lookup[name] = tokenUse{
					token:    ntok,
					position: p.scan.Position,
				}
				tok = ntok
			} else {
				if tavor.DEBUG {
					fmt.Printf("Use token %#v(%p)\n", tok, tok)
				}
			}

			p.lookupUsage[tok] = struct{}{}

			tokens = append(tokens, tok)
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
					Message:  "String is not terminated",
					Type:     ParseErrorNonTerminatedString,
					Position: p.scan.Pos(),
				}
			}

			s, _ = strconv.Unquote(s)

			tokens = append(tokens, primitives.NewConstantString(s))
		case '(':
			if tavor.DEBUG {
				fmt.Println("NEW group")
			}
			c = p.scan.Scan()
			if tavor.DEBUG {
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

			if tavor.DEBUG {
				fmt.Println("END group")
			}
		case '?':
			if tavor.DEBUG {
				fmt.Println("NEW optional")
			}
			p.expectScanRune('(')

			c = p.scan.Scan()
			if tavor.DEBUG {
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

			if tavor.DEBUG {
				fmt.Println("END optional")
			}
		case '+', '*':
			if tavor.DEBUG {
				fmt.Println("NEW repeat")
			}

			sym := c

			c = p.scan.Scan()
			if tavor.DEBUG {
				fmt.Printf("parseTerm repeat before ( %d:%v -> %v\n", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())
			}

			var from, to int

			if sym == '*' {
				from, to = 0, MaxRepeat
			} else {
				if c == scanner.Int {
					from, _ = strconv.Atoi(p.scan.TokenText())

					c = p.scan.Scan()
					if tavor.DEBUG {
						fmt.Printf("parseTerm repeat after from ( %d:%v -> %v\n", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())
					}

					// until there is an explicit "to" we can assume to==from
					to = from
				} else {
					from, to = 1, MaxRepeat
				}

				if c == ',' {
					c = p.scan.Scan()
					if tavor.DEBUG {
						fmt.Printf("parseTerm repeat after , ( %d:%v -> %v\n", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())
					}

					if c == scanner.Int {
						to, _ = strconv.Atoi(p.scan.TokenText())

						c = p.scan.Scan()
						if tavor.DEBUG {
							fmt.Printf("parseTerm repeat after to ( %d:%v -> %v\n", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())
						}
					} else {
						to = MaxRepeat
					}
				}
			}

			p.expectRune('(', c)

			if tavor.DEBUG {
				fmt.Printf("repeat from %v to %v\n", from, to)
			}

			c = p.scan.Scan()
			if tavor.DEBUG {
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

			if tavor.DEBUG {
				fmt.Println("END repeat")
			}
		case '$':
			c = p.scan.Scan()
			if tavor.DEBUG {
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
			if tavor.DEBUG {
				fmt.Printf("parseTerm multiline %d:%v -> %v\n", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())
			}

			if c == '\n' {
				return zeroRune, nil, nil, &ParserError{
					Message:  "Multi line token definition unexpectedly terminated",
					Type:     ParseErrorUnexpectedTokenDefinitionTermination,
					Position: p.scan.Pos(),
				}
			}

			continue
		default:
			if tavor.DEBUG {
				fmt.Println("break out parseTerm")
			}
			break OUT
		}

		c = p.scan.Scan()
		if tavor.DEBUG {
			fmt.Printf("parseTerm %d:%v -> %v\n", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())
		}
	}

	if len(embeddedToks) != 0 {
		embeddedTokens = append(embeddedTokens, embeddedToks)
	}

	return c, tokens, embeddedTokens, nil
}

func (p *tavorParser) parseExpression(c rune) (token.Token, error) {
	if tavor.DEBUG {
		fmt.Println("START expression")
	}

	_, err := p.expectRune('{', c)
	if err != nil {
		return nil, err
	}

	c = p.scan.Scan()
	if tavor.DEBUG {
		fmt.Printf("parseExpression after {} %d:%v -> %v\n", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())
	}

	c, tok, err := p.parseExpressionTerm(c)
	if err != nil {
		return nil, err
	} else if tok == nil {
		return nil, &ParserError{
			Message:  "Empty expressions are not allowed",
			Type:     ParseErrorEmptyExpressionIsInvalid,
			Position: p.scan.Pos(), // TODO correct position
		}
	}

	_, err = p.expectRune('}', c)
	if err != nil {
		return nil, err
	}

	if tavor.DEBUG {
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
	if tavor.DEBUG {
		fmt.Printf("parseExpressionTerm %d:%v -> %v\n", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())
	}

	// operators
	switch c {
	case '+', '-', '*', '/':
		sym := c

		c = p.scan.Scan()
		if tavor.DEBUG {
			fmt.Printf("parseExpressionTerm operator %d:%v -> %v\n", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())
		}

		var t token.Token
		c, t, err = p.parseExpressionTerm(c)
		if err != nil {
			return zeroRune, nil, err
		} else if t == nil {
			return zeroRune, nil, &ParserError{
				Message:  "Expected another expression term after operator",
				Type:     ParseErrorExpectedExpressionTerm,
				Position: p.scan.Pos(),
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
	if tavor.DEBUG {
		fmt.Println("START token attribute")
	}

	_, err := p.expectRune(scanner.Ident, c)
	if err != nil {
		return nil, err
	}

	name := p.scan.TokenText()

	tokenPosition := p.scan.Position

	_, err = p.expectScanRune('.')
	if err != nil {
		return nil, err
	}

	_, err = p.expectScanRune(scanner.Ident)
	if err != nil {
		return nil, err
	}

	attribute := p.scan.TokenText()

	use, ok := p.lookup[name]
	if !ok {
		return nil, &ParserError{
			Message:  fmt.Sprintf("Token \"%s\" is not defined", name),
			Type:     ParseErrorTokenNotDefined,
			Position: p.scan.Pos(),
		}
	}

	tok := use.token

	if tavor.DEBUG {
		fmt.Printf("Use %#v(%p) as token\n", tok, tok)
	}

	p.used[name] = tokenPosition

	if tavor.DEBUG {
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
		Message:  fmt.Sprintf("Unknown token attribute \"%s\" for token type \"%s\"", attribute, reflect.TypeOf(tok)),
		Type:     ParseErrorUnknownTokenAttribute,
		Position: p.scan.Pos(),
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
		if tavor.DEBUG {
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
				if tavor.DEBUG {
					fmt.Printf("parseScope Or %d:%v -> %v\n", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())
				}
			} else {
				if tavor.DEBUG {
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

		if tavor.DEBUG {
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

	if use, ok := p.lookup[name]; ok {
		// if there is a pointer in the lookup hash we can say that it was just used before
		if _, ok := use.token.(*primitives.Pointer); !ok {
			return zeroRune, &ParserError{
				Message:  "Token already defined",
				Type:     ParseErrorTokenAlreadyDefined,
				Position: p.scan.Pos(),
			}
		}
	}

	tokenPosition := p.scan.Position

	if c, err = p.expectScanRune('='); err != nil {
		// unexpected new line?
		if c == '\n' {
			return zeroRune, &ParserError{
				Message:  "New line inside single line token definitions is not allowed",
				Type:     ParseErrorEarlyNewLine,
				Position: p.scan.Pos(),
			}
		}

		return zeroRune, err
	}

	c = p.scan.Scan()
	if tavor.DEBUG {
		fmt.Printf("parseTokenDefinition after = %d:%v -> %v\n", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())
	}

	c, tokens, embeddedToks, err := p.parseScope(c)
	if err != nil {
		return zeroRune, err
	}

	p.embeddedTokensInTerm[name] = embeddedToks

	if tavor.DEBUG {
		fmt.Printf("Token %s embeds %+v\n", name, embeddedToks)
	}

	if tavor.DEBUG {
		fmt.Printf("back to token definition with c=%c\n", c)
	}

	// we always want a new line at the end of the file
	if c == scanner.EOF {
		return zeroRune, &ParserError{
			Message:  "New line at end of token definition needed",
			Type:     ParseErrorNewLineNeeded,
			Position: p.scan.Pos(),
		}
	}

	if c, err = p.expectRune('\n', c); err != nil {
		return zeroRune, err
	}

	var tok token.Token

	switch len(tokens) {
	case 0:
		return zeroRune, &ParserError{
			Message:  "Empty token definition",
			Type:     ParseErrorEmptyTokenDefinition,
			Position: p.scan.Pos(),
		}
	case 1:
		tok = tokens[0]
	default:
		tok = lists.NewAll(tokens...)
	}

	// self loop?
	if use, ok := p.lookup[name]; ok {
		if tavor.DEBUG {
			fmt.Printf("parseTokenDefinition fill empty pointer for %s\n", name)
		}

		err = use.token.(*primitives.Pointer).Set(tok)
		if err != nil {
			return zeroRune, &ParserError{
				Message:  fmt.Sprintf("Wrong token type for %s because of earlier usage: %s", name, err),
				Type:     ParseErrorInvalidTokenType,
				Position: p.scan.Pos(),
			}
		}
	}

	// check for endless loop
	if len(embeddedToks) != 0 {
		foundExit := false

		if tavor.DEBUG {
			fmt.Printf("Need to check for loops in %s with embedding %+v\n", name, embeddedToks)
			fmt.Printf("Use embedding lookup map %+v\n", p.embeddedTokensInTerm)
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
					if tavor.DEBUG {
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
			if tavor.DEBUG {
				fmt.Println("Found exit, everything is fine.")
			}

			foundExit = true

			break
		}

		if !foundExit {
			if tavor.DEBUG {
				fmt.Println("There is no loop exit for this token, I'll throw an error.")
			}

			return zeroRune, &ParserError{
				Message:  fmt.Sprintf("Token \"%s\" has an endless loop without exit\n", name),
				Type:     ParseErrorEndlessLoopDetected,
				Position: p.scan.Pos(), // TODO correct position
			}
		}
	}

	p.lookup[name] = tokenUse{
		token:    tok,
		position: tokenPosition,
	}

	if tavor.DEBUG {
		fmt.Printf("Added %#v(%p) as token %s\n", tok, tok, name)
	}

	c = p.scan.Scan()
	if tavor.DEBUG {
		fmt.Printf("parseTokenDefinition after newline %d:%v -> %v\n", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())
	}

	return c, nil
}

func (p *tavorParser) parseSpecialTokenDefinition() (rune, error) {
	var c rune
	var err error

	if tavor.DEBUG {
		fmt.Println("START special token")
	}

	c = p.scan.Scan()
	if tavor.DEBUG {
		fmt.Printf("parseSpecialTokenDefinition after $ %d:%v -> %v\n", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())
	}

	name := p.scan.TokenText()
	if _, ok := p.lookup[name]; ok {
		return zeroRune, &ParserError{
			Message:  "Token already defined",
			Type:     ParseErrorTokenAlreadyDefined,
			Position: p.scan.Pos(),
		}
	}

	tokenPosition := p.scan.Position

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
		if tavor.DEBUG {
			fmt.Printf("parseSpecialTokenDefinition argument value %d:%v -> %v\n", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())
		}

		switch c {
		case scanner.Ident, scanner.String, scanner.Int:
			arguments[arg] = p.scan.TokenText()
		default:
			return zeroRune, &ParserError{
				Message:  fmt.Sprintf("Invalid argument value %v", c),
				Type:     ParseErrorInvalidArgumentValue,
				Position: p.scan.Pos(),
			}
		}

		c = p.scan.Scan()
		if tavor.DEBUG {
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
			Message:  "New line at end of token definition needed",
			Type:     ParseErrorNewLineNeeded,
			Position: p.scan.Pos(),
		}
	}

	if c, err = p.expectRune('\n', c); err != nil {
		return zeroRune, err
	}

	typ, ok := arguments["type"]
	if !ok {
		return zeroRune, &ParserError{
			Message:  "Special token has no type argument",
			Type:     ParseErrorUnknownTypeForSpecialToken,
			Position: p.scan.Pos(),
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
					Message:  "Argument \"to\" is missing",
					Type:     ParseErrorMissingSpecialTokenArgument,
					Position: p.scan.Pos(),
				}
			} else if !okFrom && okTo {
				return zeroRune, &ParserError{
					Message:  "Argument \"from\" is missing",
					Type:     ParseErrorMissingSpecialTokenArgument,
					Position: p.scan.Pos(),
				}
			}

			from, err := strconv.Atoi(rawFrom)
			if err != nil {
				return zeroRune, &ParserError{
					Message:  "\"from\" needs an integer value",
					Type:     ParseErrorInvalidArgumentValue,
					Position: p.scan.Pos(),
				}
			}

			to, err := strconv.Atoi(rawTo)
			if err != nil {
				return zeroRune, &ParserError{
					Message:  "\"to\" needs an integer value",
					Type:     ParseErrorInvalidArgumentValue,
					Position: p.scan.Pos(),
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
					Message:  "\"start\" needs an integer value",
					Type:     ParseErrorInvalidArgumentValue,
					Position: p.scan.Pos(),
				}
			}
		}

		if raw, ok := arguments["step"]; ok {
			step, err = strconv.Atoi(raw)
			if err != nil {
				return zeroRune, &ParserError{
					Message:  "\"step\" needs an integer value",
					Type:     ParseErrorInvalidArgumentValue,
					Position: p.scan.Pos(),
				}
			}
		}

		usedArguments["start"] = struct{}{}
		usedArguments["step"] = struct{}{}

		tok = sequences.NewSequence(start, step)
	default:
		return zeroRune, &ParserError{
			Message:  fmt.Sprintf("Unknown special token type \"%s\"", typ),
			Type:     ParseErrorUnknownSpecialTokenType,
			Position: p.scan.Pos(),
		}
	}

	for arg := range arguments {
		if _, ok := usedArguments[arg]; !ok {
			return zeroRune, &ParserError{
				Message:  fmt.Sprintf("Unknown special token argument \"%s\"", arg),
				Type:     ParseErrorUnknownSpecialTokenArgument,
				Position: p.scan.Pos(),
			}
		}
	}

	p.lookup[name] = tokenUse{
		token:    tok,
		position: tokenPosition,
	}

	if tavor.DEBUG {
		fmt.Printf("Added %#v(%p) as token %s\n", tok, tok, name)
	}

	c = p.scan.Scan()
	if tavor.DEBUG {
		fmt.Printf("parseSpecialTokenDefinition after newline %d:%v -> %v\n", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())
	}

	if tavor.DEBUG {
		fmt.Println("END special token")
	}

	return c, nil
}

func (p *tavorParser) unrollLoops(root token.Token) token.Token {
	type unrollToken struct {
		tok    token.Token
		parent *unrollToken
	}

	if tavor.DEBUG {
		fmt.Println("Unroll loops by cloning pointers")
	}

	checked := make(map[token.Token]token.Token)
	counters := make(map[token.Token]int)

	queue := linkedlist.New()

	queue.Push(&unrollToken{
		tok:    root,
		parent: nil,
	})

	for !queue.Empty() {
		v, _ := queue.Shift()
		iTok, _ := v.(*unrollToken)

		switch t := iTok.tok.(type) {
		case *primitives.Pointer:
			o := t.InternalGet()

			parent, ok := checked[o]
			times := 0

			if ok {
				times = counters[parent]
			} else {
				parent = o.Clone()
				checked[o] = parent
			}

			if times != MaxRepeat {
				if tavor.DEBUG {
					fmt.Printf("Clone (%p)%#v with parent (%p)%#v\n", t, t, parent, parent)
				}

				c := parent.Clone()

				t.Set(c)

				counters[parent] = times + 1
				checked[c] = parent

				if iTok.parent != nil {
					switch tt := iTok.parent.tok.(type) {
					case token.ForwardToken:
						tt.InternalReplace(t, c)
					case lists.List:
						tt.InternalReplace(t, c)
					}
				} else {
					root = c
				}

				queue.Unshift(&unrollToken{
					tok:    c,
					parent: iTok.parent,
				})
			} else {
				if tavor.DEBUG {
					fmt.Printf("Reached max repeat of %d for (%p)%#v with parent (%p)%#v\n", MaxRepeat, t, t, parent, parent)
				}

				t.Set(nil)

				ta := iTok.tok
				tt := iTok.parent

			REMOVE:
				for tt != nil {
					switch l := tt.tok.(type) {
					case token.ForwardToken:
						if tavor.DEBUG {
							fmt.Printf("Remove (%p)%#v from (%p)%#v\n", ta, ta, l, l)
						}

						c := l.InternalLogicalRemove(ta)

						if c != nil {
							break REMOVE
						}

						ta = l
						tt = tt.parent
					case lists.List:
						if tavor.DEBUG {
							fmt.Printf("Remove (%p)%#v from (%p)%#v\n", ta, ta, l, l)
						}

						c := l.InternalLogicalRemove(ta)

						if c != nil {
							break REMOVE
						}

						ta = l
						tt = tt.parent
					}
				}
			}
		case token.ForwardToken:
			if v := t.InternalGet(); v != nil {
				queue.Push(&unrollToken{
					tok:    v,
					parent: iTok,
				})
			}
		case lists.List:
			for i := 0; i < t.InternalLen(); i++ {
				c, _ := t.InternalGet(i)

				queue.Push(&unrollToken{
					tok:    c,
					parent: iTok,
				})
			}
		}
	}

	if tavor.DEBUG {
		fmt.Println("Done unrolling")
	}

	return root
}

func ParseTavor(src io.Reader) (token.Token, error) {
	p := &tavorParser{
		earlyUse:             make(map[string]tokenUse),
		embeddedTokensInTerm: make(map[string][]map[string]struct{}),
		lookup:               make(map[string]tokenUse),
		lookupUsage:          make(map[token.Token]struct{}),
		used:                 make(map[string]scanner.Position),
	}

	if tavor.DEBUG {
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
			Message:  "No START token defined",
			Type:     ParseErrorNoStart,
			Position: p.scan.Pos(), // TODO correct position
		}
	}

	p.used["START"] = p.scan.Position

	for name, use := range p.earlyUse {
		if use.token.(*primitives.Pointer).Get() == nil {
			return nil, &ParserError{
				Message:  fmt.Sprintf("Token \"%s\" is not defined", name),
				Type:     ParseErrorTokenNotDefined,
				Position: use.position,
			}
		}
	}

	for name, use := range p.lookup {
		if _, ok := p.used[name]; !ok {
			return nil, &ParserError{
				Message:  fmt.Sprintf("Token \"%s\" declared but not used", name),
				Type:     ParseErrorUnusedToken,
				Position: use.position,
			}
		}
	}

	start := p.lookup["START"].token

	start = p.unrollLoops(start)

	return start, nil
}
