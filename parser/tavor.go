package parser

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"reflect"
	"strconv"
	"text/scanner"

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
	Allow correct forward usage of ExistingSequenceItems
*/

const zeroRune = 0

type tokenUsage struct {
	token         token.Token
	position      scanner.Position
	variableScope map[string]token.Token
}

type attributeForwardUsage struct {
	definitionName    string
	tokenName         string
	tokenPosition     scanner.Position
	attribute         string
	attributePosition scanner.Position
	operator          string
	operatorToken     token.Token
	pointer           *primitives.Pointer
	variableScope     map[string]token.Token
}

type tavorParser struct {
	scan scanner.Scanner

	err string

	earlyUse       map[string][]tokenUsage
	lookup         map[string]tokenUsage
	lookupUsage    map[token.Token]struct{}
	used           map[string][]tokenUsage
	variableUsages []token.Token

	forwardAttributeUsage []attributeForwardUsage
}

func (p *tavorParser) expectRune(expect rune, got rune) (rune, error) {
	if got != expect {
		return got, &token.ParserError{
			Message:  fmt.Sprintf("expected %s but got %s", scanner.TokenString(expect), scanner.TokenString(got)),
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

func (p *tavorParser) parseGlobalScope(variableScope map[string]token.Token) error {
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

func (p *tavorParser) parseTerm(definitionName string, c rune, variableScope map[string]token.Token) (rune, []token.Token, error) {
	var err error
	var tokens []token.Token

	addToken := func(tok token.Token) {
		tokens = append(tokens, tok)
	}

OUT:
	for {
		switch c {
		case scanner.Ident:
			name := p.scan.TokenText()

			nVariableScope := make(map[string]token.Token, len(variableScope))
			for k, v := range variableScope {
				nVariableScope[k] = v
			}

			_, ok := p.lookup[name]
			if !ok {
				log.Debugf("parseTerm use empty pointer for %s", name)

				var tokenInterface *token.Token
				b := primitives.NewEmptyPointer(tokenInterface)
				n := primitives.NewPointer(b)

				p.lookup[name] = tokenUsage{
					token:    n,
					position: p.scan.Position,
				}
				p.earlyUse[name] = append(p.earlyUse[name], tokenUsage{
					token:         b,
					position:      p.scan.Position,
					variableScope: nVariableScope,
				})
			}

			p.used[name] = append(p.used[name], tokenUsage{
				token:         nil,
				position:      p.scan.Position,
				variableScope: nVariableScope,
			})

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

					p.lookup[name] = tokenUsage{
						token:    ntok,
						position: p.scan.Position,
					}
					if t, ok := tok.(*primitives.Pointer); ok && t.Get() == nil {
						p.earlyUse[name] = append(p.earlyUse[name], tokenUsage{
							token:         ntok,
							position:      p.scan.Position,
							variableScope: nVariableScope,
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
				return zeroRune, nil, &token.ParserError{
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

			c, toks, err := p.parseScope(definitionName, c, variableScope)
			if err != nil {
				return zeroRune, nil, err
			}

			_, err = p.expectRune(')', c)
			if err != nil {
				return zeroRune, nil, err
			}

			switch len(toks) {
			case 0:
				// ignore
			case 1:
				addToken(toks[0])
			default:
				addToken(lists.NewAll(toks...))
			}

			log.Debug("END group")
		case '?':
			log.Debug("NEW optional")

			_, err = p.expectScanRune('(')
			if err != nil {
				return zeroRune, nil, err
			}

			c = p.scan.Scan()
			log.Debugf("parseTerm optional after ( %d:%v -> %v", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())

			c, toks, err := p.parseScope(definitionName, c, variableScope)
			if err != nil {
				return zeroRune, nil, err
			}

			_, err = p.expectRune(')', c)
			if err != nil {
				return zeroRune, nil, err
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

			var from, to token.Token

			if sym == '*' {
				from, to = primitives.NewConstantInt(0), primitives.NewConstantInt(tavor.MaxRepeat)
			} else {
				if c == scanner.Int {
					iFrom, _ := strconv.Atoi(p.scan.TokenText())
					from = primitives.NewConstantInt(iFrom)

					c = p.scan.Scan()
					log.Debugf("parseTerm repeat after from ( %d:%v -> %v", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())

					// until there is an explicit "to" we can assume to==from
					to = from // do not clone here! since really to==from
				} else if c == '$' {
					c = p.scan.Scan()

					c, from, err = p.parseTokenAttribute(definitionName, c, variableScope)

					if err != nil {
						return zeroRune, nil, err
					}

					// until there is an explicit "to" we can assume to==from
					to = from // do not clone here! since really to==from
				} else {
					from, to = primitives.NewConstantInt(1), primitives.NewConstantInt(tavor.MaxRepeat)
				}

				if c == ',' {
					c = p.scan.Scan()
					log.Debugf("parseTerm repeat after , ( %d:%v -> %v", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())

					if c == scanner.Int {
						iTo, _ := strconv.Atoi(p.scan.TokenText())
						to = primitives.NewConstantInt(iTo)

						c = p.scan.Scan()
						log.Debugf("parseTerm repeat after to ( %d:%v -> %v", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())
					} else if c == '$' {
						c = p.scan.Scan()

						c, to, err = p.parseTokenAttribute(definitionName, c, variableScope)

						if err != nil {
							return zeroRune, nil, err
						}
					} else {
						to = primitives.NewConstantInt(tavor.MaxRepeat)
					}
				}
			}

			_, err = p.expectRune('(', c)
			if err != nil {
				return zeroRune, nil, err
			}

			log.Debugf("repeat from %v to %v", from, to)

			c = p.scan.Scan()
			log.Debugf("parseTerm repeat after ( %d:%v -> %v", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())

			c, toks, err := p.parseScope(definitionName, c, variableScope)
			if err != nil {
				return zeroRune, nil, err
			}

			_, err = p.expectRune(')', c)
			if err != nil {
				return zeroRune, nil, err
			}

			switch len(toks) {
			case 0:
				// ignore
			case 1:
				addToken(lists.NewRepeatWithTokens(toks[0], from, to))
			default:
				addToken(lists.NewRepeatWithTokens(lists.NewAll(toks...), from, to))
			}

			log.Debug("END repeat")
		case '$':
			c = p.scan.Scan()
			log.Debugf("parseTerm after $ ( %d:%v -> %v", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())

			var tok token.Token

			if c == '{' {
				c, tok, err = p.parseExpression(definitionName, variableScope)
				if err != nil {
					return zeroRune, nil, err
				}

				_, err = p.expectRune('}', c)
				if err != nil {
					return zeroRune, nil, err
				}

				c = p.scan.Scan()
			} else {
				c, tok, err = p.parseTokenAttribute(definitionName, c, variableScope)
			}

			if err != nil {
				return zeroRune, nil, err
			}

			addToken(tok)

			continue
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

				return zeroRune, nil, err
			}

			addToken(primitives.NewCharacterClass(pattern.String()))

			log.Debug("END character class")
		case '<':
			log.Debug("NEW variable")

			c = p.scan.Scan()

			justSave := false

			if c == '=' {
				justSave = true

				c = p.scan.Scan()
			}

			_, err = p.expectRune(scanner.Ident, c)
			if err != nil {
				return zeroRune, nil, err
			}

			if len(tokens) == 0 {
				return zeroRune, nil, &token.ParserError{
					Message:  "Variable has to be assigned to a token",
					Type:     token.ParseErrorNoTokenForVariable,
					Position: p.scan.Pos(),
				}
			}

			variableName := p.scan.TokenText()

			_, err = p.expectScanRune('>')
			if err != nil {
				return zeroRune, nil, err
			}

			// TODO do not overwrite Token names... this sould lead to an already defined error, only variables can overwrite each other

			var variable token.Token

			if justSave {
				variable = variables.NewVariableSave(variableName, tokens[len(tokens)-1])
			} else {
				variable = variables.NewVariable(variableName, tokens[len(tokens)-1])
			}

			tokens[len(tokens)-1] = variable
			variableScope[variableName] = variable

			p.variableUsages = append(p.variableUsages, variable)

			log.Debug("END variable")
		case ',': // multi line token
			if _, err := p.expectScanRune('\n'); err != nil {
				return zeroRune, nil, err
			}

			c = p.scan.Scan()
			log.Debugf("parseTerm multiline %d:%v -> %v", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())

			if c == '\n' {
				return zeroRune, nil, &token.ParserError{
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

	return c, tokens, nil
}

func (p *tavorParser) parseExpression(definitionName string, variableScope map[string]token.Token) (rune, token.Token, error) {
	log.Debug("START expression")

	c := p.scan.Scan()
	log.Debugf("parseExpression %d:%v -> %v", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())

	c, tok, err := p.parseExpressionTerm(definitionName, c, variableScope)
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

func (p *tavorParser) parseExpressionTerm(definitionName string, c rune, variableScope map[string]token.Token) (rune, token.Token, error) {
	var tok token.Token
	var err error

	// single term
	switch c {
	case scanner.Ident:
		attribute := p.scan.TokenText()
		attributePosition := p.scan.Position

		switch attribute {
		case "defined":
			_, err := p.expectScanRune(scanner.Ident)
			if err != nil {
				return zeroRune, nil, err
			}

			name := p.scan.TokenText()
			namePosition := p.scan.Position

			tok, ok := variableScope[name]

			isPointer := false

			if ok {
				if _, pp := tok.(*primitives.Pointer); pp {
					isPointer = true
				}
			}

			nVariableScope := make(map[string]token.Token, len(variableScope))
			for k, v := range variableScope {
				nVariableScope[k] = v
			}

			if !ok || isPointer {
				log.Debugf("parseTokenAttribute use empty pointer for %s.%s", name, attribute)

				var tokenInterface *token.Token

				pointer := primitives.NewEmptyPointer(tokenInterface)
				nPointer := primitives.NewPointer(pointer)

				p.forwardAttributeUsage = append(p.forwardAttributeUsage, attributeForwardUsage{
					definitionName:    definitionName,
					tokenName:         name,
					tokenPosition:     namePosition,
					attribute:         attribute,
					attributePosition: attributePosition,
					pointer:           pointer,
					variableScope:     nVariableScope,
				})

				tok = nPointer
			} else {
				tok, err = p.selectTokenAttribute(tok, name, attribute, attributePosition, "", nil, nVariableScope)
				if err != nil {
					return zeroRune, nil, err
				}
			}

			variableScope[name] = tok
			nVariableScope[name] = tok

			c = p.scan.Scan()

			return c, conditions.NewExpressionPointer(tok), nil
		default:
			c, tok, err = p.parseTokenAttribute(definitionName, c, variableScope)
			if err != nil {
				return zeroRune, nil, err
			}
		}
	case scanner.Int:
		v, _ := strconv.Atoi(p.scan.TokenText())

		tok = primitives.NewConstantInt(v)

		c = p.scan.Scan()
	}

	if tok == nil {
		return zeroRune, nil, nil
	}

	log.Debugf("parseExpressionTerm %d:%v -> %v", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())

	// operators
	switch c {
	case '+', '-', '*', '/':
		sym := c

		c = p.scan.Scan()
		log.Debugf("parseExpressionTerm operator %d:%v -> %v", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())

		var t token.Token
		c, t, err = p.parseExpressionTerm(definitionName, c, variableScope)
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

func (p *tavorParser) parseTokenAttribute(definitionName string, c rune, variableScope map[string]token.Token) (rune, token.Token, error) {
	log.Debug("new token attribute")

	_, err := p.expectRune(scanner.Ident, c)
	if err != nil {
		return zeroRune, nil, err
	}

	name := p.scan.TokenText()

	tokenPosition := p.scan.Position

	_, err = p.expectScanRune('.')
	if err != nil {
		return zeroRune, nil, err
	}

	_, err = p.expectScanRune(scanner.Ident)
	if err != nil {
		return zeroRune, nil, err
	}

	attribute := p.scan.TokenText()
	attributePosition := p.scan.Position

	var op string
	var opToken token.Token

	c = p.scan.Scan()

	if c == scanner.Ident && p.scan.TokenText() == "not" {
		_, err := p.expectScanRune(scanner.Ident)
		if err != nil {
			return zeroRune, nil, err
		}

		if p.scan.TokenText() != "in" {
			return zeroRune, nil, &token.ParserError{
				Message:  fmt.Sprintf("expected operator %q not %q", "in", p.scan.TokenText()),
				Type:     token.ParseErrorExpectOperator,
				Position: p.scan.Pos(),
			}
		}

		log.Debug("operator: not in")

		_, err = p.expectScanRune('(')
		if err != nil {
			return zeroRune, nil, err
		}

		_, err = p.expectScanRune(scanner.Ident)
		if err != nil {
			return zeroRune, nil, err
		}

		expectName := p.scan.TokenText()
		expectNamePosition := p.scan.Position

		expectTok, ok := variableScope[expectName]

		isPointer := false

		if ok {
			if _, pp := expectTok.(*primitives.Pointer); pp {
				isPointer = true
			}
		}

		attribute := "Value"

		if !ok || isPointer {
			log.Debugf("parseTokenAttribute use empty pointer for %s.%s", expectName, attribute)

			var tokenInterface *token.Token

			pointer := primitives.NewEmptyPointer(tokenInterface)
			nPointer := primitives.NewPointer(pointer)

			variableScope[expectName] = nPointer

			nVariableScope := make(map[string]token.Token, len(variableScope))
			for k, v := range variableScope {
				nVariableScope[k] = v
			}

			p.forwardAttributeUsage = append(p.forwardAttributeUsage, attributeForwardUsage{
				definitionName:    definitionName,
				tokenName:         expectName,
				tokenPosition:     expectNamePosition,
				attribute:         attribute,
				attributePosition: expectNamePosition,
				pointer:           pointer,
				variableScope:     nVariableScope,
			})

			expectTok = nPointer
		} else {
			expectTok, err = p.selectTokenAttribute(expectTok, name, attribute, expectNamePosition, "", nil, variableScope)
			if err != nil {
				return zeroRune, nil, err
			}
		}

		_, err = p.expectScanRune(')')
		if err != nil {
			return zeroRune, nil, err
		}

		op = "not in"
		opToken = expectTok

		c = p.scan.Scan()
	}

	var tok token.Token

	use, ok := p.lookup[name]
	if ok {
		tok = use.token
	} else {
		tok, ok = variableScope[name]

		isPointer := false

		if ok {
			if _, pp := tok.(*primitives.Pointer); pp {
				isPointer = true
			}
		}

		if !ok || isPointer {
			log.Debugf("parseTokenAttribute use empty pointer for %s.%s", name, attribute)

			var tokenInterface *token.Token

			pointer := primitives.NewEmptyPointer(tokenInterface)
			nPointer := primitives.NewPointer(pointer)

			variableScope[name] = nPointer

			nVariableScope := make(map[string]token.Token, len(variableScope))
			for k, v := range variableScope {
				nVariableScope[k] = v
			}

			p.forwardAttributeUsage = append(p.forwardAttributeUsage, attributeForwardUsage{
				definitionName:    definitionName,
				tokenName:         name,
				tokenPosition:     tokenPosition,
				attribute:         attribute,
				attributePosition: attributePosition,
				operator:          op,
				operatorToken:     opToken,
				pointer:           pointer,
				variableScope:     nVariableScope,
			})

			return c, nPointer, nil
		}
	}

	p.used[name] = append(p.earlyUse[name], tokenUsage{
		token:         nil,
		position:      tokenPosition,
		variableScope: variableScope,
	})

	rtok, err := p.selectTokenAttribute(tok, name, attribute, attributePosition, op, opToken, variableScope)

	if err == nil {
		log.Debugf("Insert token attribute %p(%#v)", rtok, rtok)
	}

	return c, rtok, err
}

func (p *tavorParser) selectTokenAttribute(tok token.Token, tokenName string, attribute string, attributePosition scanner.Position, operator string, operatorToken token.Token, variableScope map[string]token.Token) (token.Token, error) {
	log.Debugf("use (%p)%#v as token", tok, tok)

	log.Debug("finished token attribute (or will be unknown token attribute)")

	switch i := tok.(type) {
	case token.List:
		switch attribute {
		case "Count":
			return aggregates.NewLen(i), nil
		case "Unique":
			return lists.NewUniqueItem(i), nil
		}
	case *sequences.Sequence:
		switch attribute {
		case "Existing":
			if operator == "not in" {
				return i.ExistingItem(operatorToken), nil
			}

			return i.ExistingItem(nil), nil
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
	case token.VariableToken:
		switch attribute {
		case "defined":
			return conditions.NewVariableDefined(tokenName, variableScope), nil
		case "Index":
			return lists.NewIndexItem(variables.NewVariableValue(i)), nil
		case "Value":
			v := variables.NewVariableValue(i)

			p.variableUsages = append(p.variableUsages, i)

			return v, nil
		}
	}

	return nil, &token.ParserError{
		Message:  fmt.Sprintf("unknown token attribute %q for token type %q", attribute, reflect.TypeOf(tok)),
		Type:     token.ParseErrorUnknownTokenAttribute,
		Position: attributePosition,
	}
}

func (p *tavorParser) parseScope(definitionName string, c rune, variableScope map[string]token.Token) (rune, []token.Token, error) {
	var err error
	var tokens []token.Token

	var toks []token.Token

	c, toks, err = p.parseTerm(definitionName, c, variableScope)
	if err != nil {
		return zeroRune, nil, err
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

				if c == '|' {
					c = p.scan.Scan()
					log.Debugf("parseScope Or %d:%v -> %v", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())
				} else {
					log.Debug("parseScope break out or")

					break OR
				}

				c, toks, err = p.parseTerm(definitionName, c, variableScope)
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

			log.Debug("END or")
		case '{': // TODO make conditions work with ORs...
			log.Debug("NEW condition")

			c = p.scan.Scan()
			condition := p.scan.TokenText()

			var conditionExpression conditions.BooleanExpression

			switch condition {
			case "if":
				log.Debug("found IF")

				c, conditionExpression, err = p.parseConditionExpression(definitionName, variableScope)
				if err != nil {
					return zeroRune, nil, err
				}
			case "else":
				if len(ifPairs) == 0 {
					panic("TODO else/elseif without if")
				}

				c = p.scan.Scan()

				if p.scan.TokenText() == "if" {
					log.Debug("found ELSEIF")

					c, conditionExpression, err = p.parseConditionExpression(definitionName, variableScope)
					if err != nil {
						return zeroRune, nil, err
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
				return zeroRune, nil, &token.ParserError{
					Message:  fmt.Sprintf("unknown condition %q", condition),
					Type:     token.ParseErrorUnknownCondition,
					Position: p.scan.Pos(),
				}
			}

			_, err = p.expectRune('}', c)
			if err != nil {
				return zeroRune, nil, err
			}

			c = p.scan.Scan()

			if condition != "endif" {
				c, toks, err = p.parseTerm(definitionName, c, variableScope) // TODO this should be a a global scope or so ... we can do nesting
				if err != nil {
					return zeroRune, nil, err
				}

				var tok token.Token

				switch len(toks) {
				case 0:
					panic("empty body is not allowed") // TODO catch this
				case 1:
					tok = toks[0]
				default:
					tok = lists.NewAll(toks...)
				}

				ifPart := conditions.IfPair{
					Head: conditionExpression,
					Body: tok,
				}

				log.Debugf("Add if part %#v", ifPart)

				ifPairs = append(ifPairs, ifPart)
			} else {
				tokens = append(tokens, conditions.NewIf(ifPairs...))

				ifPairs = nil

				c, toks, err = p.parseTerm(definitionName, c, variableScope) // TODO this should be a a global scope or so ... we can do nesting
				if err != nil {
					return zeroRune, nil, err
				}

				tokens = append(tokens, toks...)
			}

			log.Debug("END condition")

			continue SCOPE
		}

		break SCOPE
	}

	if len(ifPairs) != 0 {
		panic("TODO if without endif")
	}

	return c, tokens, nil
}

func (p *tavorParser) parseConditionExpression(definitionName string, variableScope map[string]token.Token) (rune, conditions.BooleanExpression, error) {
	c, a, err := p.parseExpression(definitionName, variableScope)
	if err != nil {
		return zeroRune, nil, err
	}

	if ex, ok := a.(conditions.BooleanExpression); ok {
		return c, ex, nil
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

	c, b, err := p.parseExpression(definitionName, variableScope)
	if err != nil {
		return zeroRune, nil, err
	}

	return c, conditions.NewBooleanEqual(a, b), nil
}

func (p *tavorParser) parseTokenDefinition(variableScope map[string]token.Token) (rune, error) {
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

	// each definition start its own scope
	nVariableScope := make(map[string]token.Token, len(variableScope))
	for k, v := range variableScope {
		nVariableScope[k] = v
	}
	variableScope = nVariableScope

	// start reading definition
	c = p.scan.Scan()
	log.Debugf("parseTokenDefinition after = %d:%v -> %v", p.scan.Line, scanner.TokenString(c), p.scan.TokenText())

	c, tokens, err := p.parseScope(name, c, variableScope)
	if err != nil {
		return zeroRune, err
	}

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

	p.lookup[name] = tokenUsage{
		token:         tok,
		position:      tokenPosition,
		variableScope: variableScope,
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
			from := 0
			if okFrom {
				from, err = strconv.Atoi(rawFrom)
				if err != nil {
					return zeroRune, &token.ParserError{
						Message:  "\"from\" needs an integer value",
						Type:     token.ParseErrorInvalidArgumentValue,
						Position: p.scan.Pos(),
					}
				}
			}

			to := math.MaxInt64
			if okTo {
				to, err = strconv.Atoi(rawTo)
				if err != nil {
					return zeroRune, &token.ParserError{
						Message:  "\"to\" needs an integer value",
						Type:     token.ParseErrorInvalidArgumentValue,
						Position: p.scan.Pos(),
					}
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

	p.lookup[name] = tokenUsage{
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
		earlyUse:    make(map[string][]tokenUsage),
		lookup:      make(map[string]tokenUsage),
		lookupUsage: make(map[token.Token]struct{}),
		used:        make(map[string][]tokenUsage),
	}

	log.Debug("start parsing tavor file")

	p.scan.Init(src)

	p.scan.Error = func(s *scanner.Scanner, msg string) {
		p.err = msg
	}
	p.scan.Whitespace = 1<<'\t' | 1<<' ' | 1<<'\r'

	variableScope := make(map[string]token.Token)

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

	p.used["START"] = append(p.used["START"], tokenUsage{
		token:         nil,
		position:      p.scan.Position,
		variableScope: variableScope,
	})

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

	for _, forwardUse := range p.forwardAttributeUsage {
		var tok token.Token

		use, ok := p.lookup[forwardUse.tokenName]
		if ok {
			tok = use.token
		} else {
			found := false
			var typ reflect.Type

			for _, usage := range p.used[forwardUse.definitionName] {
				if t, ok := usage.variableScope[forwardUse.tokenName]; ok {
					found = true

					if typ == nil {
						tok = t
						typ = reflect.TypeOf(t)
					} else {
						tType := reflect.TypeOf(t)

						if typ != tType {
							panic(fmt.Errorf("TODO handle this better: different types %v vs %v", typ, tType))
						}
					}
				} else if forwardUse.attribute != "defined" {
					return nil, &token.ParserError{
						Message:  fmt.Sprintf("token or variable %q is not defined", forwardUse.tokenName),
						Type:     token.ParseErrorTokenNotDefined,
						Position: usage.position,
					}
				}
			}

			if !found {
				return nil, &token.ParserError{
					Message:  fmt.Sprintf("token or variable %q is not defined", forwardUse.tokenName),
					Type:     token.ParseErrorTokenNotDefined,
					Position: forwardUse.tokenPosition,
				}
			}
		}

		if po, ok := tok.(*primitives.Pointer); ok {
			log.Debugf("Found pointer in forwardAttributeUsage %#v", forwardUse)

			for {
				c := po.InternalGet()

				po, ok = c.(*primitives.Pointer)
				if !ok {
					log.Debugf("Replaced pointer %p(%#v) with %p(%#v)", tok, tok, c, c)

					tok = c

					break
				}
			}
		}

		rtok, err := p.selectTokenAttribute(tok, forwardUse.tokenName, forwardUse.attribute, forwardUse.attributePosition, forwardUse.operator, forwardUse.operatorToken, variableScope)
		if err != nil {
			return nil, err
		}

		err = forwardUse.pointer.Set(rtok)
		if err != nil {
			return nil, err
		}

		p.used[forwardUse.tokenName] = append(p.used[forwardUse.tokenName], tokenUsage{
			token:         nil,
			position:      forwardUse.tokenPosition,
			variableScope: forwardUse.variableScope,
		})
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
		tok := variable.(token.ForwardToken).InternalGet()

		if po, ok := tok.(*primitives.Pointer); ok {
			log.Debugf("Found pointer in variable %p(%#v)", variable, variable)

			for {
				c := po.InternalGet()

				po, ok = c.(*primitives.Pointer)
				if !ok {
					log.Debugf("Replaced pointer %p(%#v) with %p(%#v)", tok, tok, c, c)

					variable.(token.ForwardToken).InternalReplace(tok, c)

					break
				}
			}
		}
	}

	start := p.lookup["START"].token

	start = tavor.UnrollPointers(start)

	variableScope = make(map[string]token.Token)
	tavor.SetInternalScope(start, variableScope)

	log.Debug("finished parsing")

	return start, nil
}
