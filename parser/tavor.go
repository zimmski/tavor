package parser

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"text/scanner"

	"github.com/zimmski/container/list/linkedlist"

	"github.com/zimmski/tavor"
	"github.com/zimmski/tavor/log"
	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/aggregates"
	"github.com/zimmski/tavor/token/conditions"
	"github.com/zimmski/tavor/token/constraints"
	"github.com/zimmski/tavor/token/expressions"
	"github.com/zimmski/tavor/token/lists"
	"github.com/zimmski/tavor/token/primitives"
	"github.com/zimmski/tavor/token/sequences"
	"github.com/zimmski/tavor/token/variables"
)

/*
	TODO

	Token names can only consist of letters, digits and "_"
	Allow forward usage of token attributes
	Allow correct forward usage of ExistingSequenceItems
*/

const zeroRune = 0

type tokenUse struct {
	token    token.Token
	position scanner.Position
}

type tavorParser struct {
	scan scanner.Scanner

	err string

	earlyUse             map[string][]tokenUse
	embeddedTokensInTerm map[string][]map[string]struct{}
	lookup               map[string]tokenUse
	lookupUsage          map[token.Token]struct{}
	used                 map[string]scanner.Position
	variableUsages       []*variables.Variable
}

func (p *tavorParser) expectRune(expect rune, got rune) (rune, error) {
	if got != expect {
		return got, &token.ParserError{
			Message:  fmt.Sprintf("expected %q but got %q", scanner.TokenString(expect), scanner.TokenString(got)),
			Type:     token.ParseErrorExpectRune,
			Position: p.scan.Pos(),
		}
	}

	return got, nil
}

func (p *tavorParser) expectScanRune(expect rune) (rune, error) {
	got := p.scan.Scan()

	log.Debugf("%d:%v -> %v", p.scan.Line, scanner.TokenString(got), p.scan.TokenText())

	return p.expectRune(expect, got)
}

func (p *tavorParser) parseGlobalScope(variableScope map[string]*variables.Variable) error {
	var err error

	c := p.scan.Scan()
	log.Debugf("%d:%v -> %v", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())

	for c != scanner.EOF {
		switch c {
		case '\n':
			// ignore new lines in the global scope
		case scanner.Ident:
			c, err = p.parseTokenDefinition(variableScope)
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
			return &token.ParserError{
				Message:  fmt.Sprintf("token names have to start with a letter and not with %s", scanner.TokenString(c)),
				Type:     token.ParseErrorInvalidTokenName,
				Position: p.scan.Pos(),
			}
		}

		c = p.scan.Scan()
		log.Debugf("%d:%v -> %v", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())
	}

	return nil
}

func (p *tavorParser) parseTerm(c rune, variableScope map[string]*variables.Variable) (rune, []token.Token, []map[string]struct{}, error) {
	var err error
	var embeddedTokens = make([]map[string]struct{}, 0)
	var embeddedToks = make(map[string]struct{}, 0)
	var tokens []token.Token

	addToken := func(tok token.Token) {
		tokens = append(tokens, tok)
	}

OUT:
	for {
		switch c {
		case scanner.Ident:
			name := p.scan.TokenText()

			_, ok := p.lookup[name]
			if !ok {
				log.Debugf("parseTerm use empty pointer for %s", name)

				var tokenInterface *token.Token
				b := primitives.NewEmptyPointer(tokenInterface)
				n := primitives.NewPointer(b)

				p.lookup[name] = tokenUse{
					token:    n,
					position: p.scan.Position,
				}
				p.earlyUse[name] = append(p.earlyUse[name], tokenUse{
					token:    b,
					position: p.scan.Position,
				})
			}

			embeddedToks[name] = struct{}{}
			p.used[name] = p.scan.Position

			/*

				THIS IS A TERRIBLE HACK THIS SHOULD BE -REMOVED- fixed ASAP

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
				if t, ok := tok.(*primitives.Pointer); ok && t.Get() == nil {
					// FIXME if tok is directly given to NewPointer we get a panic: reflect: non-interface type passed to Type.Implements
					var tokInterface *token.Token
					ntok := primitives.NewEmptyPointer(tokInterface)
					ntok.Set(tok)

					log.Debugf("token %s (%p)%#v is an empty pointer, better just forward to it (%p)%#v", name, tok, tok, ntok, ntok)

					tok = ntok
				} else {
					ntok := tok.Clone()

					log.Debugf("token %s (%p)%#v was already used once. Cloned as (%p)%#v", name, tok, tok, ntok, ntok)

					p.lookup[name] = tokenUse{
						token:    ntok,
						position: p.scan.Position,
					}
					if t, ok := tok.(*primitives.Pointer); ok && t.Get() == nil {
						p.earlyUse[name] = append(p.earlyUse[name], tokenUse{
							token:    ntok,
							position: p.scan.Position,
						})
					}

					tok = ntok
				}
			} else {
				log.Debugf("use token (%p)%#v", tok, tok)
			}

			p.lookupUsage[tok] = struct{}{}

			addToken(tok)
		case scanner.Int:
			v, _ := strconv.Atoi(p.scan.TokenText())

			addToken(primitives.NewConstantInt(v))
		case scanner.String:
			s := p.scan.TokenText()

			if s[len(s)-1] != '"' {
				return zeroRune, nil, nil, &token.ParserError{
					Message:  "string is not terminated",
					Type:     token.ParseErrorNonTerminatedString,
					Position: p.scan.Pos(),
				}
			}

			s, _ = strconv.Unquote(s)

			addToken(primitives.NewConstantString(s))
		case '(':
			log.Debug("NEW group")

			c = p.scan.Scan()
			log.Debugf("parseTerm Group %d:%v -> %v", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())

			c, toks, embeddedToks, err := p.parseScope(c, variableScope)
			if err != nil {
				return zeroRune, nil, nil, err
			}

			_, err = p.expectRune(')', c)
			if err != nil {
				return zeroRune, nil, nil, err
			}

			switch len(toks) {
			case 0:
				// ignore
			case 1:
				addToken(toks[0])
			default:
				addToken(lists.NewAll(toks...))
			}

			if len(embeddedToks) != 0 {
				embeddedTokens = append(embeddedTokens, embeddedToks...)
			}

			log.Debug("END group")
		case '?':
			log.Debug("NEW optional")

			_, err = p.expectScanRune('(')
			if err != nil {
				return zeroRune, nil, nil, err
			}

			c = p.scan.Scan()
			log.Debugf("parseTerm optional after ( %d:%v -> %v", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())

			c, toks, _, err := p.parseScope(c, variableScope)
			if err != nil {
				return zeroRune, nil, nil, err
			}

			_, err = p.expectRune(')', c)
			if err != nil {
				return zeroRune, nil, nil, err
			}

			switch len(toks) {
			case 0:
				// ignore
			case 1:
				addToken(constraints.NewOptional(toks[0]))
			default:
				addToken(constraints.NewOptional(lists.NewAll(toks...)))
			}

			log.Debug("END optional")
		case '+', '*':
			log.Debug("NEW repeat")

			sym := c

			c = p.scan.Scan()
			log.Debugf("parseTerm repeat before ( %d:%v -> %v", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())

			var from, to int

			if sym == '*' {
				from, to = 0, tavor.MaxRepeat
			} else {
				if c == scanner.Int {
					from, _ = strconv.Atoi(p.scan.TokenText())

					c = p.scan.Scan()
					log.Debugf("parseTerm repeat after from ( %d:%v -> %v", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())

					// until there is an explicit "to" we can assume to==from
					to = from
				} else {
					from, to = 1, tavor.MaxRepeat
				}

				if c == ',' {
					c = p.scan.Scan()
					log.Debugf("parseTerm repeat after , ( %d:%v -> %v", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())

					if c == scanner.Int {
						to, _ = strconv.Atoi(p.scan.TokenText())

						c = p.scan.Scan()
						log.Debugf("parseTerm repeat after to ( %d:%v -> %v", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())
					} else {
						to = tavor.MaxRepeat
					}
				}
			}

			_, err = p.expectRune('(', c)
			if err != nil {
				return zeroRune, nil, nil, err
			}

			log.Debugf("repeat from %v to %v", from, to)

			c = p.scan.Scan()
			log.Debugf("parseTerm repeat after ( %d:%v -> %v", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())

			c, toks, embeddedToks, err := p.parseScope(c, variableScope)
			if err != nil {
				return zeroRune, nil, nil, err
			}

			_, err = p.expectRune(')', c)
			if err != nil {
				return zeroRune, nil, nil, err
			}

			switch len(toks) {
			case 0:
				// ignore
			case 1:
				addToken(lists.NewRepeat(toks[0], int64(from), int64(to)))
			default:
				addToken(lists.NewRepeat(lists.NewAll(toks...), int64(from), int64(to)))
			}

			if from > 0 && len(embeddedToks) != 0 {
				embeddedTokens = append(embeddedTokens, embeddedToks...)
			}

			log.Debug("END repeat")
		case '$':
			c = p.scan.Scan()
			log.Debugf("parseTerm after $ ( %d:%v -> %v", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())

			var tok token.Token

			if c == '{' {
				c, tok, err = p.parseExpression(variableScope)
				if err != nil {
					return zeroRune, nil, nil, err
				}

				_, err = p.expectRune('}', c)
				if err != nil {
					return zeroRune, nil, nil, err
				}
			} else {
				tok, err = p.parseTokenAttribute(c, variableScope)
			}

			if err != nil {
				return zeroRune, nil, nil, err
			}

			addToken(tok)
		case '[':
			log.Debug("NEW character class")

			var pattern bytes.Buffer

			c = p.scan.Scan()

			for c != ']' && c != '\n' && c != scanner.EOF {
				pattern.WriteString(p.scan.TokenText())

				c = p.scan.Scan()
			}

			if c != ']' {
				_, err := p.expectRune(']', c)

				return zeroRune, nil, nil, err
			}

			addToken(primitives.NewCharacterClass(pattern.String()))

			log.Debug("END character class")
		case '<':
			log.Debug("NEW variable")

			_, err = p.expectScanRune(scanner.Ident)
			if err != nil {
				return zeroRune, nil, nil, err
			}

			if len(tokens) == 0 {
				return zeroRune, nil, nil, &token.ParserError{
					Message:  "Variable has to be assigned to a token",
					Type:     token.ParseErrorNoTokenForVariable,
					Position: p.scan.Pos(),
				}
			}

			variableName := p.scan.TokenText()

			_, err = p.expectScanRune('>')
			if err != nil {
				return zeroRune, nil, nil, err
			}

			variable := variables.NewVariable(tokens[len(tokens)-1])

			tokens[len(tokens)-1] = variable
			variableScope[variableName] = variable

			p.variableUsages = append(p.variableUsages, variable)

			log.Debug("END variable")
		case ',': // multi line token
			if _, err := p.expectScanRune('\n'); err != nil {
				return zeroRune, nil, nil, err
			}

			c = p.scan.Scan()
			log.Debugf("parseTerm multiline %d:%v -> %v", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())

			if c == '\n' {
				return zeroRune, nil, nil, &token.ParserError{
					Message:  "multi line token definition unexpectedly terminated",
					Type:     token.ParseErrorUnexpectedTokenDefinitionTermination,
					Position: p.scan.Pos(),
				}
			}

			continue
		default:
			log.Debug("break out parseTerm")
			break OUT
		}

		c = p.scan.Scan()
		log.Debugf("parseTerm %d:%v -> %v", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())
	}

	if len(embeddedToks) != 0 {
		embeddedTokens = append(embeddedTokens, embeddedToks)
	}

	return c, tokens, embeddedTokens, nil
}

func (p *tavorParser) parseExpression(variableScope map[string]*variables.Variable) (rune, token.Token, error) {
	log.Debug("START expression")

	c := p.scan.Scan()
	log.Debugf("parseExpression %d:%v -> %v", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())

	c, tok, err := p.parseExpressionTerm(c, variableScope)
	if err != nil {
		return zeroRune, nil, err
	} else if tok == nil {
		return zeroRune, nil, &token.ParserError{
			Message:  "empty expressions are not allowed",
			Type:     token.ParseErrorEmptyExpressionIsInvalid,
			Position: p.scan.Pos(), // TODO correct position
		}
	}

	log.Debug("END expression")

	return c, tok, nil
}

func (p *tavorParser) parseExpressionTerm(c rune, variableScope map[string]*variables.Variable) (rune, token.Token, error) {
	var tok token.Token
	var err error

	// single term
	switch c {
	case scanner.Ident:
		tok, err = p.parseTokenAttribute(c, variableScope)
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
	log.Debugf("parseExpressionTerm %d:%v -> %v", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())

	// operators
	switch c {
	case '+', '-', '*', '/':
		sym := c

		c = p.scan.Scan()
		log.Debugf("parseExpressionTerm operator %d:%v -> %v", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())

		var t token.Token
		c, t, err = p.parseExpressionTerm(c, variableScope)
		if err != nil {
			return zeroRune, nil, err
		} else if t == nil {
			return zeroRune, nil, &token.ParserError{
				Message:  "expected another expression term after operator",
				Type:     token.ParseErrorExpectedExpressionTerm,
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

func (p *tavorParser) parseTokenAttribute(c rune, variableScope map[string]*variables.Variable) (token.Token, error) {
	log.Debug("new token attribute")

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

	var tok token.Token

	use, ok := p.lookup[name]
	if ok {
		tok = use.token
	} else {
		tok, ok = variableScope[name]

		if !ok {
			return nil, &token.ParserError{
				Message:  fmt.Sprintf("token or variable %q is not defined", name),
				Type:     token.ParseErrorTokenNotDefined,
				Position: p.scan.Pos(),
			}
		}
	}

	log.Debugf("use (%p)%#v as token", tok, tok)

	p.used[name] = tokenPosition

	log.Debug("finished token attribute (or will be unknown token attribute)")

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
	case *primitives.RandomInt, *primitives.RangeInt:
		switch attribute {
		case "Value":
			return i.Clone(), nil
		}
	case *variables.Variable:
		switch attribute {
		case "Value":
			v := variables.NewVariableValue(i)

			p.variableUsages = append(p.variableUsages, i)

			return v, nil
		}
	}

	return nil, &token.ParserError{
		Message:  fmt.Sprintf("unknown token attribute %q for token type %q", attribute, reflect.TypeOf(tok)),
		Type:     token.ParseErrorUnknownTokenAttribute,
		Position: p.scan.Pos(),
	}
}

func (p *tavorParser) parseScope(c rune, variableScope map[string]*variables.Variable) (rune, []token.Token, []map[string]struct{}, error) {
	var err error
	var embeddedTokens = make([]map[string]struct{}, 0)
	var tokens []token.Token

	var toks []token.Token
	var embeddedToks []map[string]struct{}

	c, toks, embeddedToks, err = p.parseTerm(c, variableScope)
	if err != nil {
		return zeroRune, nil, nil, err
	} else if toks != nil {
		tokens = toks
	}

	var ifPairs []conditions.IfPair

SCOPE:
	for {
		switch c {
		case '|':
			log.Debug("NEW or")

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
					log.Debugf("parseScope Or %d:%v -> %v", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())
				} else {
					log.Debug("parseScope break out or")

					break OR
				}

				c, toks, embeddedToks, err = p.parseTerm(c, variableScope)
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

			log.Debug("END or")
		case '{': // TODO make conditions work with ORs...
			log.Debug("NEW condition")

			c = p.scan.Scan()
			condition := p.scan.TokenText()

			var conditionExpression conditions.BooleanExpression

			switch condition {
			case "if":
				log.Debug("found IF")

				c, conditionExpression, err = p.parseConditionExpression(variableScope)
				if err != nil {
					return zeroRune, nil, nil, err
				}
			case "else":
				if len(ifPairs) == 0 {
					panic("TODO else/elseif without if")
				}

				c = p.scan.Scan()

				if p.scan.TokenText() == "if" {
					log.Debug("found ELSEIF")

					c, conditionExpression, err = p.parseConditionExpression(variableScope)
					if err != nil {
						return zeroRune, nil, nil, err
					}
				} else {
					log.Debug("found ELSE")

					conditionExpression = conditions.NewBooleanTrue()
				}
			case "endif":
				log.Debug("found ENDIF")

				if len(ifPairs) == 0 {
					panic("TODO endif without if")
				}

				c = p.scan.Scan()
			default:
				return zeroRune, nil, nil, &token.ParserError{
					Message:  fmt.Sprintf("unknown condition %q", condition),
					Type:     token.ParseErrorUnknownCondition,
					Position: p.scan.Pos(),
				}
			}

			_, err = p.expectRune('}', c)
			if err != nil {
				return zeroRune, nil, nil, err
			}

			c = p.scan.Scan()

			if condition != "endif" {
				// TODO do something with embeddedTOks here
				c, toks, embeddedToks, err = p.parseTerm(c, variableScope) // TODO this should be a a global scope or so ... we can do nesting
				if err != nil {
					return zeroRune, nil, nil, err
				}

				embeddedTokens = append(embeddedTokens, embeddedToks...)

				var tok token.Token

				switch len(toks) {
				case 0:
					panic("empty body is not allowed") // TODO catch this
				case 1:
					tok = toks[0]
				default:
					tok = lists.NewAll(toks...)
				}

				ifPairs = append(ifPairs, conditions.IfPair{
					Head: conditionExpression,
					Body: tok,
				})
			} else {
				tokens = []token.Token{conditions.NewIf(ifPairs...)}

				ifPairs = nil

				break SCOPE
			}

			log.Debug("END condition")

			continue SCOPE
		default:
			if len(embeddedToks) != 0 {
				embeddedTokens = append(embeddedTokens, embeddedToks...)
			}
		}

		break SCOPE
	}

	if len(ifPairs) != 0 {
		panic("TODO if without endif")
	}

	return c, tokens, embeddedTokens, nil
}

func (p *tavorParser) parseConditionExpression(variableScope map[string]*variables.Variable) (rune, conditions.BooleanExpression, error) {
	c, a, err := p.parseExpression(variableScope)
	if err != nil {
		return zeroRune, nil, err
	}

	switch c {
	case '=':
		_, err = p.expectScanRune('=')
		if err != nil {
			return zeroRune, nil, err
		}
	default:
		return zeroRune, nil, &token.ParserError{
			Message:  fmt.Sprintf("unknown boolean operator %q", c),
			Type:     token.ParseErrorUnknownBooleanOperator,
			Position: p.scan.Pos(),
		}
	}

	c, b, err := p.parseExpression(variableScope)
	if err != nil {
		return zeroRune, nil, err
	}

	return c, conditions.NewBooleanEqual(a, b), nil
}

func (p *tavorParser) parseTokenDefinition(variableScope map[string]*variables.Variable) (rune, error) {
	var c rune
	var err error

	name := p.scan.TokenText()

	if use, ok := p.lookup[name]; ok {
		// if there is a pointer in the lookup hash we can say that it was just used before
		if _, ok := use.token.(*primitives.Pointer); !ok {
			return zeroRune, &token.ParserError{
				Message:  "token already defined",
				Type:     token.ParseErrorTokenAlreadyDefined,
				Position: p.scan.Pos(),
			}
		}
	}

	tokenPosition := p.scan.Position

	if c, err = p.expectScanRune('='); err != nil {
		// unexpected new line?
		if c == '\n' {
			return zeroRune, &token.ParserError{
				Message:  "new line inside single line token definitions is not allowed",
				Type:     token.ParseErrorEarlyNewLine,
				Position: p.scan.Pos(),
			}
		}

		return zeroRune, err
	}

	c = p.scan.Scan()
	log.Debugf("parseTokenDefinition after = %d:%v -> %v", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())

	c, tokens, embeddedToks, err := p.parseScope(c, variableScope)
	if err != nil {
		return zeroRune, err
	}

	p.embeddedTokensInTerm[name] = embeddedToks

	log.Debugf("token %s embeds %+v", name, embeddedToks)

	log.Debugf("back to token definition with c=%c", c)

	// we always want a new line at the end of the file
	if c == scanner.EOF {
		return zeroRune, &token.ParserError{
			Message:  "new line at end of token definition needed",
			Type:     token.ParseErrorNewLineNeeded,
			Position: p.scan.Pos(),
		}
	}

	if c, err = p.expectRune('\n', c); err != nil {
		return zeroRune, err
	}

	var tok token.Token

	switch len(tokens) {
	case 0:
		return zeroRune, &token.ParserError{
			Message:  "empty token definition",
			Type:     token.ParseErrorEmptyTokenDefinition,
			Position: p.scan.Pos(),
		}
	case 1:
		tok = tokens[0]
	default:
		tok = lists.NewAll(tokens...)
	}

	// self loop?
	if uses, ok := p.earlyUse[name]; ok {
		log.Debugf("parseTokenDefinition fill empty pointer for %s", name)

		for _, use := range uses {
			log.Debugf("use (%p)%#v for pointer (%p)%#v", tok, tok, use.token, use.token)

			err = use.token.(*primitives.Pointer).Set(tok)
			if err != nil {
				return zeroRune, &token.ParserError{
					Message:  fmt.Sprintf("wrong token type for %s because of earlier usage: %s", name, err),
					Type:     token.ParseErrorInvalidTokenType,
					Position: p.scan.Pos(),
				}
			}
		}

		delete(p.earlyUse, name)
	}

	// check for endless loop
	if len(embeddedToks) != 0 {
		foundExit := false

		log.Debugf("need to check for loops in %s with embedding %+v", name, embeddedToks)
		log.Debugf("use embedding lookup map %+v", p.embeddedTokensInTerm)

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
					log.Debug("found one loop")

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

			log.Debug("found exit, everything is fine.")

			foundExit = true

			break
		}

		if !foundExit {
			log.Debug("there is no loop exit for this token, I'll throw an error.")

			return zeroRune, &token.ParserError{
				Message:  fmt.Sprintf("token %q has an endless loop without exit\n", name),
				Type:     token.ParseErrorEndlessLoopDetected,
				Position: p.scan.Pos(), // TODO correct position
			}
		}
	}

	p.lookup[name] = tokenUse{
		token:    tok,
		position: tokenPosition,
	}

	log.Debugf("added (%p)%#v as token %s", tok, tok, name)

	c = p.scan.Scan()
	log.Debugf("parseTokenDefinition after newline %d:%v -> %v", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())

	return c, nil
}

func (p *tavorParser) parseSpecialTokenDefinition() (rune, error) {
	var c rune
	var err error

	log.Debug("START special token")

	c = p.scan.Scan()
	log.Debugf("parseSpecialTokenDefinition after $ %d:%v -> %v", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())

	name := p.scan.TokenText()
	if _, ok := p.lookup[name]; ok {
		return zeroRune, &token.ParserError{
			Message:  "token already defined",
			Type:     token.ParseErrorTokenAlreadyDefined,
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
		log.Debugf("parseSpecialTokenDefinition argument value %d:%v -> %v", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())

		switch c {
		case scanner.Ident, scanner.String, scanner.Int:
			arguments[arg] = p.scan.TokenText()
		default:
			return zeroRune, &token.ParserError{
				Message:  fmt.Sprintf("invalid argument value %v", c),
				Type:     token.ParseErrorInvalidArgumentValue,
				Position: p.scan.Pos(),
			}
		}

		c = p.scan.Scan()
		log.Debugf("parseSpecialTokenDefinition after argument value %d:%v -> %v", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())

		if c != ',' {
			break
		}

		if c, err = p.expectScanRune('\n'); err != nil {
			return zeroRune, err
		}
	}

	// we always want a new line at the end of the file
	if c == scanner.EOF {
		return zeroRune, &token.ParserError{
			Message:  "new line at end of token definition needed",
			Type:     token.ParseErrorNewLineNeeded,
			Position: p.scan.Pos(),
		}
	}

	if c, err = p.expectRune('\n', c); err != nil {
		return zeroRune, err
	}

	typ, ok := arguments["type"]
	if !ok {
		return zeroRune, &token.ParserError{
			Message:  "special token has no type argument",
			Type:     token.ParseErrorUnknownTypeForSpecialToken,
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
				return zeroRune, &token.ParserError{
					Message:  "argument \"to\" is missing",
					Type:     token.ParseErrorMissingSpecialTokenArgument,
					Position: p.scan.Pos(),
				}
			} else if !okFrom && okTo {
				return zeroRune, &token.ParserError{
					Message:  "argument \"from\" is missing",
					Type:     token.ParseErrorMissingSpecialTokenArgument,
					Position: p.scan.Pos(),
				}
			}

			from, err := strconv.Atoi(rawFrom)
			if err != nil {
				return zeroRune, &token.ParserError{
					Message:  "\"from\" needs an integer value",
					Type:     token.ParseErrorInvalidArgumentValue,
					Position: p.scan.Pos(),
				}
			}

			to, err := strconv.Atoi(rawTo)
			if err != nil {
				return zeroRune, &token.ParserError{
					Message:  "\"to\" needs an integer value",
					Type:     token.ParseErrorInvalidArgumentValue,
					Position: p.scan.Pos(),
				}
			}

			step := 1

			if raw, ok := arguments["step"]; ok {
				step, err = strconv.Atoi(raw)
				if err != nil {
					return zeroRune, &token.ParserError{
						Message:  "\"step\" needs an integer value",
						Type:     token.ParseErrorInvalidArgumentValue,
						Position: p.scan.Pos(),
					}
				}

				usedArguments["step"] = struct{}{}
			}

			usedArguments["from"] = struct{}{}
			usedArguments["to"] = struct{}{}

			tok = primitives.NewRangeIntWithStep(from, to, step)
		} else {
			tok = primitives.NewRandomInt()
		}
	case "Sequence":
		start := 1
		step := 1

		if raw, ok := arguments["start"]; ok {
			start, err = strconv.Atoi(raw)
			if err != nil {
				return zeroRune, &token.ParserError{
					Message:  "\"start\" needs an integer value",
					Type:     token.ParseErrorInvalidArgumentValue,
					Position: p.scan.Pos(),
				}
			}
		}

		if raw, ok := arguments["step"]; ok {
			step, err = strconv.Atoi(raw)
			if err != nil {
				return zeroRune, &token.ParserError{
					Message:  "\"step\" needs an integer value",
					Type:     token.ParseErrorInvalidArgumentValue,
					Position: p.scan.Pos(),
				}
			}
		}

		usedArguments["start"] = struct{}{}
		usedArguments["step"] = struct{}{}

		tok = sequences.NewSequence(start, step)
	default:
		return zeroRune, &token.ParserError{
			Message:  fmt.Sprintf("unknown special token type %q", typ),
			Type:     token.ParseErrorUnknownSpecialTokenType,
			Position: p.scan.Pos(),
		}
	}

	for arg := range arguments {
		if _, ok := usedArguments[arg]; !ok {
			return zeroRune, &token.ParserError{
				Message:  fmt.Sprintf("unknown special token argument %q", arg),
				Type:     token.ParseErrorUnknownSpecialTokenArgument,
				Position: p.scan.Pos(),
			}
		}
	}

	p.lookup[name] = tokenUse{
		token:    tok,
		position: tokenPosition,
	}

	log.Debugf("added (%p)%#v as token %s", tok, tok, name)

	c = p.scan.Scan()
	log.Debugf("parseSpecialTokenDefinition after newline %d:%v -> %v", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())

	log.Debug("END special token")

	return c, nil
}

func ParseTavor(src io.Reader) (token.Token, error) {
	p := &tavorParser{
		earlyUse:             make(map[string][]tokenUse),
		embeddedTokensInTerm: make(map[string][]map[string]struct{}),
		lookup:               make(map[string]tokenUse),
		lookupUsage:          make(map[token.Token]struct{}),
		used:                 make(map[string]scanner.Position),
	}

	log.Debug("start parsing tavor file")

	p.scan.Init(src)

	p.scan.Error = func(s *scanner.Scanner, msg string) {
		p.err = msg
	}
	p.scan.Whitespace = 1<<'\t' | 1<<' ' | 1<<'\r'

	variableScope := make(map[string]*variables.Variable)

	if err := p.parseGlobalScope(variableScope); err != nil {
		return nil, err
	}

	if _, ok := p.lookup["START"]; !ok {
		return nil, &token.ParserError{
			Message:  "no START token defined",
			Type:     token.ParseErrorNoStart,
			Position: p.scan.Pos(), // TODO correct position
		}
	}

	p.used["START"] = p.scan.Position

	for name, uses := range p.earlyUse {
		for _, use := range uses {
			if use.token.(*primitives.Pointer).Get() == nil {
				return nil, &token.ParserError{
					Message:  fmt.Sprintf("token %q is not defined", name),
					Type:     token.ParseErrorTokenNotDefined,
					Position: use.position,
				}
			}
		}
	}

	for name, use := range p.lookup {
		if _, ok := p.used[name]; !ok {
			return nil, &token.ParserError{
				Message:  fmt.Sprintf("token %q declared but not used", name),
				Type:     token.ParseErrorUnusedToken,
				Position: use.position,
			}
		}
	}

	for _, variable := range p.variableUsages {
		tok := variable.InternalGet()

		if po, ok := tok.(*primitives.Pointer); ok {
			log.Debugf("Found pointer in variable %p(%#v)", variable, variable)

			for {
				c := po.InternalGet()

				po, ok = c.(*primitives.Pointer)
				if !ok {
					log.Debugf("Replaced pointer %p(%#v) with %p(%#v)", tok, tok, c, c)

					variable.InternalReplace(tok, c)

					break
				}
			}
		}
	}

	start := p.lookup["START"].token

	start = tavor.UnrollPointers(start)

	log.Debug("finished parsing")

	return start, nil
}
