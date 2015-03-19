package parser

import (
	"fmt"
	"math"
	"strings"
	"testing"

	. "github.com/zimmski/tavor/test/assert"

	"github.com/zimmski/tavor"
	"github.com/zimmski/tavor/fuzz/strategy"
	"github.com/zimmski/tavor/test"
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

func TestTavorParseErrors(t *testing.T) {
	var tok token.Token
	var err error

	// empty file
	tok, err = ParseTavor(strings.NewReader(""))
	Equal(t, token.ParseErrorNoStart, err.(*token.ParserError).Type)
	Nil(t, tok)

	// empty file
	tok, err = ParseTavor(strings.NewReader("START = 123"))
	Equal(t, token.ParseErrorNewLineNeeded, err.(*token.ParserError).Type)
	Nil(t, tok)

	// new line before =
	tok, err = ParseTavor(strings.NewReader("START \n= 123\n"))
	Equal(t, token.ParseErrorEarlyNewLine, err.(*token.ParserError).Type)
	Nil(t, tok)

	// expect =
	tok, err = ParseTavor(strings.NewReader("START 123\n"))
	Equal(t, token.ParseErrorExpectRune, err.(*token.ParserError).Type)
	Nil(t, tok)

	// new line after =
	tok, err = ParseTavor(strings.NewReader("START = \n123\n"))
	Equal(t, token.ParseErrorEmptyTokenDefinition, err.(*token.ParserError).Type)
	Nil(t, tok)

	// invalid token name. does not start with letter
	tok, err = ParseTavor(strings.NewReader("3TART = 123\n"))
	Equal(t, token.ParseErrorInvalidTokenName, err.(*token.ParserError).Type)
	Nil(t, tok)

	tok, err = ParseTavor(strings.NewReader("+0,2(1)\n"))
	Equal(t, token.ParseErrorInvalidTokenName, err.(*token.ParserError).Type)
	Nil(t, tok)

	// unused token
	tok, err = ParseTavor(strings.NewReader("START = 123\nNumber = 123\n"))
	Equal(t, token.ParseErrorUnusedToken, err.(*token.ParserError).Type)
	Nil(t, tok)

	// non-terminated string
	tok, err = ParseTavor(strings.NewReader("START = \"non-terminated string\n"))
	Equal(t, token.ParseErrorNonTerminatedString, err.(*token.ParserError).Type)
	Nil(t, tok)

	// token already defined
	tok, err = ParseTavor(strings.NewReader("START = 123\nSTART = 456\n"))
	Equal(t, token.ParseErrorTokenAlreadyDefined, err.(*token.ParserError).Type)
	Nil(t, tok)

	// token is not defined
	tok, err = ParseTavor(strings.NewReader("START = Token\n"))
	Equal(t, token.ParseErrorTokenNotDefined, err.(*token.ParserError).Type)
	Nil(t, tok)

	// unexpected multi line token termination
	tok, err = ParseTavor(strings.NewReader("Hello = 1,\n\n"))
	Equal(t, token.ParseErrorUnexpectedTokenDefinitionTermination, err.(*token.ParserError).Type)
	Nil(t, tok)

	// unexpected continue of multi line token
	tok, err = ParseTavor(strings.NewReader("Hello = 1,2\n"))
	Equal(t, token.ParseErrorExpectRune, err.(*token.ParserError).Type)
	Nil(t, tok)

	// unknown token attribute
	tok, err = ParseTavor(strings.NewReader("Token = 123\nSTART = $Token.yeah\n"))
	Equal(t, token.ParseErrorUnknownTokenAttribute, err.(*token.ParserError).Type)
	Nil(t, tok)

	// unknown token attribute
	tok, err = ParseTavor(strings.NewReader("Token = 123\nSTART = Token $Token.Count\n"))
	Equal(t, token.ParseErrorUnknownTokenAttribute, err.(*token.ParserError).Type)
	Nil(t, tok)

	// token not defined for token attribute
	tok, err = ParseTavor(strings.NewReader("START = $Token.Count\n"))
	Equal(t, token.ParseErrorTokenNotDefined, err.(*token.ParserError).Type)
	Nil(t, tok)

	// typed token already defined
	tok, err = ParseTavor(strings.NewReader("START = 123\n$START = 456\n"))
	Equal(t, token.ParseErrorTokenAlreadyDefined, err.(*token.ParserError).Type)
	Nil(t, tok)

	// no type for typed token
	tok, err = ParseTavor(strings.NewReader("$START = To: 123\n"))
	Equal(t, token.ParseErrorTypeNotDefinedForTypedToken, err.(*token.ParserError).Type)
	Nil(t, tok)

	// expect identifier as type for typed token
	tok, err = ParseTavor(strings.NewReader("$START 123\n"))
	Equal(t, token.ParseErrorTypeNotDefinedForTypedToken, err.(*token.ParserError).Type)
	Nil(t, tok)

	// expect identifier in typed token
	tok, err = ParseTavor(strings.NewReader("$START Int = 123\n"))
	Equal(t, token.ParseErrorExpectRune, err.(*token.ParserError).Type)
	Nil(t, tok)

	// expect : in typed token
	tok, err = ParseTavor(strings.NewReader("$START Int = argument\n"))
	Equal(t, token.ParseErrorExpectRune, err.(*token.ParserError).Type)
	Nil(t, tok)

	// expect valid argument value in typed token
	tok, err = ParseTavor(strings.NewReader("$START Int = argument:\n"))
	Equal(t, token.ParseErrorInvalidArgumentValue, err.(*token.ParserError).Type)
	Nil(t, tok)

	// expect no eof after argument value in typed token
	tok, err = ParseTavor(strings.NewReader("$START Int = argument:value"))
	Equal(t, token.ParseErrorNewLineNeeded, err.(*token.ParserError).Type)
	Nil(t, tok)

	// expect new line in typed token
	tok, err = ParseTavor(strings.NewReader("$START Int = argument:value$"))
	Equal(t, token.ParseErrorExpectRune, err.(*token.ParserError).Type)
	Nil(t, tok)

	// unknown type argument typed token
	tok, err = ParseTavor(strings.NewReader("$START Unknown\n"))
	Equal(t, token.ParseErrorUnknownTypedTokenType, err.(*token.ParserError).Type)
	Nil(t, tok)

	// unknown typed token argument
	tok, err = ParseTavor(strings.NewReader("$START Int = ok: value\n"))
	Equal(t, token.ParseErrorUnknownTypedTokenArgument, err.(*token.ParserError).Type)
	Nil(t, tok)

	// invalid arguments for typed token Int
	tok, err = ParseTavor(strings.NewReader("$START Int = from:abc,\nto:123\n"))
	Equal(t, token.ParseErrorInvalidArgumentValue, err.(*token.ParserError).Type)
	Nil(t, tok)

	tok, err = ParseTavor(strings.NewReader("$START Int = from:123,\nto:abc\n"))
	Equal(t, token.ParseErrorInvalidArgumentValue, err.(*token.ParserError).Type)
	Nil(t, tok)

	tok, err = ParseTavor(strings.NewReader("$START Int = from:123,\nto:456,\nstep:abc\n"))
	Equal(t, token.ParseErrorInvalidArgumentValue, err.(*token.ParserError).Type)
	Nil(t, tok)

	// invalid arguments for typed token Sequence
	tok, err = ParseTavor(strings.NewReader("$START Sequence = start:abc\n"))
	Equal(t, token.ParseErrorInvalidArgumentValue, err.(*token.ParserError).Type)
	Nil(t, tok)

	tok, err = ParseTavor(strings.NewReader("$START Sequence = step:abc\n"))
	Equal(t, token.ParseErrorInvalidArgumentValue, err.(*token.ParserError).Type)
	Nil(t, tok)

	// empty expression
	tok, err = ParseTavor(strings.NewReader("START = ${}\n"))
	Equal(t, token.ParseErrorEmptyExpressionIsInvalid, err.(*token.ParserError).Type)
	Nil(t, tok)

	// open expression
	tok, err = ParseTavor(strings.NewReader("$Spec Sequence\nSTART = ${Spec.Next\n"))
	Equal(t, token.ParseErrorExpectRune, err.(*token.ParserError).Type)
	Nil(t, tok)

	// missing operator expression term
	tok, err = ParseTavor(strings.NewReader("$Spec Sequence\nSTART = ${Spec.Next +}\n"))
	Equal(t, token.ParseErrorExpectedExpressionTerm, err.(*token.ParserError).Type)
	Nil(t, tok)

	// wrong token type because of earlier usage
	tok, err = ParseTavor(strings.NewReader("START = $List.Count\nList = 123\n"))
	Equal(t, token.ParseErrorUnknownTokenAttribute, err.(*token.ParserError).Type)
	Nil(t, tok)

	// missing closing ] for character class
	tok, err = ParseTavor(strings.NewReader("START = [ab"))
	Equal(t, token.ParseErrorExpectRune, err.(*token.ParserError).Type)
	Nil(t, tok)
	tok, err = ParseTavor(strings.NewReader("START = [ab\n"))
	Equal(t, token.ParseErrorExpectRune, err.(*token.ParserError).Type)
	Nil(t, tok)

	// no token for variable
	tok, err = ParseTavor(strings.NewReader("START = <hey>\n"))
	Equal(t, token.ParseErrorNoTokenForVariable, err.(*token.ParserError).Type)
	Nil(t, tok)

	// variable not defined because of different scope
	tok, err = ParseTavor(strings.NewReader(`
		Save = "text"<var>
		Print = $var.Value
		START = Save Print
	`))
	Equal(t, token.ParseErrorTokenNotDefined, err.(*token.ParserError).Type)
	Nil(t, tok)

	tok, err = ParseTavor(strings.NewReader(`
		START = Token Print
		Token = "text"<var> Print
		Print = $var.Value
	`))
	Equal(t, token.ParseErrorTokenNotDefined, err.(*token.ParserError).Type)
	Nil(t, tok)

	// repeats with an optional term are not allowed
	tok, err = ParseTavor(strings.NewReader("START = +(?(1))\n"))
	Equal(t, token.ParseErrorRepeatWithOptionalTerm, err.(*token.ParserError).Type)
	Nil(t, tok)

	tok, err = ParseTavor(strings.NewReader("START = *(1 | )\n"))
	Equal(t, token.ParseErrorRepeatWithOptionalTerm, err.(*token.ParserError).Type)
	Nil(t, tok)

	tok, err = ParseTavor(strings.NewReader("START = +(1 | 2 | ?(3))\n"))
	Equal(t, token.ParseErrorRepeatWithOptionalTerm, err.(*token.ParserError).Type)
	Nil(t, tok)

	// empty strings are not allowed
	tok, err = ParseTavor(strings.NewReader("START = \"\"\n"))
	Equal(t, token.ParseErrorEmptyString, err.(*token.ParserError).Type)
	Nil(t, tok)

	// loops in list argument of path operator is not allowed
	tok, err = ParseTavor(strings.NewReader(`
			START = Pairs "->" Path

			Path = ${Pairs path from (2) over (e.Item(0)) connect by (e.Item(1)) without (0)}

			Pairs = (,
				(1 0 Path),
				(3 1),
				(2 3),
			)
		`))
	Equal(t, token.ParseErrorEndlessLoopDetected, err.(*token.ParserError).Type)
	Nil(t, tok)
}

func TestTavorParserSimple(t *testing.T) {
	var tok token.Token
	var err error

	// constant integer
	tok, err = ParseTavor(strings.NewReader("START = 123\n"))
	Nil(t, err)
	Equal(t, tok, primitives.NewScope(primitives.NewConstantInt(123)))

	// single line comment
	tok, err = ParseTavor(strings.NewReader("// hello\nSTART = 123\n"))
	Nil(t, err)
	Equal(t, tok, primitives.NewScope(primitives.NewConstantInt(123)))

	// single line multi line comment
	tok, err = ParseTavor(strings.NewReader("/* hello */\nSTART = 123\n"))
	Nil(t, err)
	Equal(t, tok, primitives.NewScope(primitives.NewConstantInt(123)))

	// multi line multi line comment
	tok, err = ParseTavor(strings.NewReader("/*\nh\ne\nl\nl\no\n*/\nSTART = 123\n"))
	Nil(t, err)
	Equal(t, tok, primitives.NewScope(primitives.NewConstantInt(123)))

	// inline comment
	tok, err = ParseTavor(strings.NewReader("START /* ok */= /* or so */ 123\n"))
	Nil(t, err)
	Equal(t, tok, primitives.NewScope(primitives.NewConstantInt(123)))

	// constant string
	tok, err = ParseTavor(strings.NewReader("START = \"abc\"\n"))
	Nil(t, err)
	Equal(t, tok, primitives.NewScope(primitives.NewConstantString("abc")))

	// constant string with whitespaces and epic chars
	tok, err = ParseTavor(strings.NewReader("START = \"a b c !\\n\\\"$%&/\"\n"))
	Nil(t, err)
	Equal(t, tok, primitives.NewScope(primitives.NewConstantString("a b c !\n\"$%&/")))

	// concatination
	tok, err = ParseTavor(strings.NewReader("START = \"I am a constant string\" 123\n"))
	Nil(t, err)
	Equal(t, tok, primitives.NewScope(lists.NewAll(
		primitives.NewConstantString("I am a constant string"),
		primitives.NewConstantInt(123),
	)))

	// embed token
	tok, err = ParseTavor(strings.NewReader("Token=123\nSTART = Token\n"))
	Nil(t, err)
	Equal(t, tok, primitives.NewScope(primitives.NewConstantInt(123)))

	// embed over token
	tok, err = ParseTavor(strings.NewReader("Token=123\nAnotherToken = Token\nSTART = AnotherToken\n"))
	Nil(t, err)
	Equal(t, tok, primitives.NewScope(primitives.NewConstantInt(123)))

	// multi line token
	tok, err = ParseTavor(strings.NewReader("START = 1,\n2,\n3\n"))
	Nil(t, err)
	Equal(t, tok, primitives.NewScope(lists.NewAll(
		primitives.NewConstantInt(1),
		primitives.NewConstantInt(2),
		primitives.NewConstantInt(3),
	)))

	// Umläüt
	tok, err = ParseTavor(strings.NewReader("Umläüt=123\nSTART = Umläüt\n"))
	Nil(t, err)
	Equal(t, tok, primitives.NewScope(primitives.NewConstantInt(123)))
}

func TestTavorParserAlternationsAndGroupings(t *testing.T) {
	var tok token.Token
	var err error

	// simple alternation
	tok, err = ParseTavor(strings.NewReader("START = 1 | 2 | 3\n"))
	Nil(t, err)
	Equal(t, tok, primitives.NewScope(lists.NewOne(
		primitives.NewConstantInt(1),
		primitives.NewConstantInt(2),
		primitives.NewConstantInt(3),
	)))

	// concatinated alternation
	tok, err = ParseTavor(strings.NewReader("START = 1 | 2 3 | 4\n"))
	Nil(t, err)
	Equal(t, tok, primitives.NewScope(lists.NewOne(
		primitives.NewConstantInt(1),
		lists.NewAll(
			primitives.NewConstantInt(2),
			primitives.NewConstantInt(3),
		),
		primitives.NewConstantInt(4),
	)))

	// optional alternation
	tok, err = ParseTavor(strings.NewReader("START = | 2 | 3\n"))
	Nil(t, err)
	Equal(t, tok, primitives.NewScope(constraints.NewOptional(lists.NewOne(
		primitives.NewConstantInt(2),
		primitives.NewConstantInt(3),
	))))

	tok, err = ParseTavor(strings.NewReader("START = 1 | | 3\n"))
	Nil(t, err)
	Equal(t, tok, primitives.NewScope(constraints.NewOptional(lists.NewOne(
		primitives.NewConstantInt(1),
		primitives.NewConstantInt(3),
	))))

	tok, err = ParseTavor(strings.NewReader("START = 1 | 2 |\n"))
	Nil(t, err)
	Equal(t, tok, primitives.NewScope(constraints.NewOptional(lists.NewOne(
		primitives.NewConstantInt(1),
		primitives.NewConstantInt(2),
	))))

	// alternation with embedded token
	tok, err = ParseTavor(strings.NewReader("Token = 2\nSTART = 1 | Token\n"))
	Nil(t, err)
	Equal(t, tok, primitives.NewScope(lists.NewOne(
		primitives.NewConstantInt(1),
		primitives.NewScope(primitives.NewConstantInt(2)),
	)))

	// simple group
	tok, err = ParseTavor(strings.NewReader("START = (1 2 3)\n"))
	Nil(t, err)
	Equal(t, tok, primitives.NewScope(lists.NewAll(
		primitives.NewConstantInt(1),
		primitives.NewConstantInt(2),
		primitives.NewConstantInt(3),
	)))

	// simple embedded group
	tok, err = ParseTavor(strings.NewReader("START = 0 (1 2 3) 4\n"))
	Nil(t, err)
	Equal(t, tok, primitives.NewScope(lists.NewAll(
		primitives.NewConstantInt(0),
		lists.NewAll(
			primitives.NewConstantInt(1),
			primitives.NewConstantInt(2),
			primitives.NewConstantInt(3),
		),
		primitives.NewConstantInt(4),
	)))

	// simple embedded or group
	tok, err = ParseTavor(strings.NewReader("START = 0 (1 | 2 | 3) 4\n"))
	Nil(t, err)
	Equal(t, tok, primitives.NewScope(lists.NewAll(
		primitives.NewConstantInt(0),
		lists.NewOne(
			primitives.NewConstantInt(1),
			primitives.NewConstantInt(2),
			primitives.NewConstantInt(3),
		),
		primitives.NewConstantInt(4),
	)))

	// Yo dog, I heard you like groups? so here is a group in a group
	tok, err = ParseTavor(strings.NewReader("START = (1 | (2 | 3)) | 4\n"))
	Nil(t, err)
	Equal(t, tok, primitives.NewScope(lists.NewOne(
		lists.NewOne(
			primitives.NewConstantInt(1),
			lists.NewOne(
				primitives.NewConstantInt(2),
				primitives.NewConstantInt(3),
			),
		),
		primitives.NewConstantInt(4),
	)))

	// simple optional
	tok, err = ParseTavor(strings.NewReader("START = 1 ?(2)\n"))
	Nil(t, err)
	Equal(t, tok, primitives.NewScope(lists.NewAll(
		primitives.NewConstantInt(1),
		constraints.NewOptional(primitives.NewConstantInt(2)),
	)))

	// or optional
	tok, err = ParseTavor(strings.NewReader("START = 1 ?(2 | 3) 4\n"))
	Nil(t, err)
	Equal(t, tok, primitives.NewScope(lists.NewAll(
		primitives.NewConstantInt(1),
		constraints.NewOptional(lists.NewOne(
			primitives.NewConstantInt(2),
			primitives.NewConstantInt(3),
		)),
		primitives.NewConstantInt(4),
	)))

	// simple repeat
	tok, err = ParseTavor(strings.NewReader("START = 1 +(2)\n"))
	Nil(t, err)
	Equal(t, tok, primitives.NewScope(lists.NewAll(
		primitives.NewConstantInt(1),
		lists.NewRepeat(primitives.NewConstantInt(2), 1, int64(tavor.MaxRepeat)),
	)))

	// or repeat
	tok, err = ParseTavor(strings.NewReader("START = 1 +(2 | 3) 4\n"))
	Nil(t, err)
	Equal(t, tok, primitives.NewScope(lists.NewAll(
		primitives.NewConstantInt(1),
		lists.NewRepeat(lists.NewOne(
			primitives.NewConstantInt(2),
			primitives.NewConstantInt(3),
		), 1, int64(tavor.MaxRepeat)),
		primitives.NewConstantInt(4),
	)))

	// simple optional repeat
	tok, err = ParseTavor(strings.NewReader("START = 1 *(2)\n"))
	Nil(t, err)
	Equal(t, tok, primitives.NewScope(lists.NewAll(
		primitives.NewConstantInt(1),
		lists.NewRepeat(primitives.NewConstantInt(2), 0, int64(tavor.MaxRepeat)),
	)))

	// or optional repeat
	tok, err = ParseTavor(strings.NewReader("START = 1 *(2 | 3) 4\n"))
	Nil(t, err)
	Equal(t, tok, primitives.NewScope(lists.NewAll(
		primitives.NewConstantInt(1),
		lists.NewRepeat(lists.NewOne(
			primitives.NewConstantInt(2),
			primitives.NewConstantInt(3),
		), 0, int64(tavor.MaxRepeat)),
		primitives.NewConstantInt(4),
	)))

	// simple optional repeat
	tok, err = ParseTavor(strings.NewReader("START = 1 *(2)\n"))
	Nil(t, err)
	Equal(t, tok, primitives.NewScope(lists.NewAll(
		primitives.NewConstantInt(1),
		lists.NewRepeat(primitives.NewConstantInt(2), 0, int64(tavor.MaxRepeat)),
	)))

	// exact repeat
	tok, err = ParseTavor(strings.NewReader("START = 1 +3(2)\n"))
	Nil(t, err)
	Equal(t, tok, primitives.NewScope(lists.NewAll(
		primitives.NewConstantInt(1),
		lists.NewRepeat(primitives.NewConstantInt(2), 3, 3),
	)))

	// at least repeat
	tok, err = ParseTavor(strings.NewReader("START = 1 +3,(2)\n"))
	Nil(t, err)
	Equal(t, tok, primitives.NewScope(lists.NewAll(
		primitives.NewConstantInt(1),
		lists.NewRepeat(primitives.NewConstantInt(2), 3, int64(tavor.MaxRepeat)),
	)))

	// at most repeat
	tok, err = ParseTavor(strings.NewReader("START = 1 +,3(2)\n"))
	Nil(t, err)
	Equal(t, tok, primitives.NewScope(lists.NewAll(
		primitives.NewConstantInt(1),
		lists.NewRepeat(primitives.NewConstantInt(2), 1, 3),
	)))

	// range repeat
	tok, err = ParseTavor(strings.NewReader("START = 1 +2,3(2)\n"))
	Nil(t, err)
	Equal(t, tok, primitives.NewScope(lists.NewAll(
		primitives.NewConstantInt(1),
		lists.NewRepeat(primitives.NewConstantInt(2), 2, 3),
	)))

	// once list
	tok, err = ParseTavor(strings.NewReader("START = @(1 | 2 | 3)\n"))
	Nil(t, err)
	Equal(t, tok, primitives.NewScope(lists.NewOnce(
		primitives.NewConstantInt(1),
		primitives.NewConstantInt(2),
		primitives.NewConstantInt(3),
	)))
}

func TestTavorParserTokenAttributes(t *testing.T) {
	// token attribute List.Count
	{
		tok, err := ParseTavor(strings.NewReader(`
			Digit = 1 | 2 | 3
			Digits = *(Digit)
			START = Digits "->" $Digits.Count
		`))
		Nil(t, err)

		v, _ := tok.(*primitives.Scope).InternalGet().(*lists.All).Get(0)
		list := v.(*primitives.Scope).InternalGet().(*lists.Repeat)

		Equal(t, tok, primitives.NewScope(lists.NewAll(
			primitives.NewScope(lists.NewRepeat(primitives.NewScope(lists.NewOne(
				primitives.NewConstantInt(1),
				primitives.NewConstantInt(2),
				primitives.NewConstantInt(3),
			)), 0, int64(tavor.MaxRepeat))),
			primitives.NewConstantString("->"),
			aggregates.NewLen(list),
		)))

		strat := strategy.NewRandomStrategy(tok)
		ch, err := strat.Fuzz(test.NewRandTest(1))
		Nil(t, err)

		for i := range ch {
			Equal(t, "2->1", tok.String())

			ch <- i
		}
	}
	// token attribute List.Item
	{
		tok, err := ParseTavor(strings.NewReader(`
			Digits = 1 2 3
			START = Digits "->" $Digits.Item(2) $Digits.Item(1) $Digits.Item(0)
		`))
		Nil(t, err)

		strat := strategy.NewRandomStrategy(tok)
		ch, err := strat.Fuzz(test.NewRandTest(1))
		Nil(t, err)

		for i := range ch {
			Equal(t, "123->321", tok.String())

			ch <- i
		}
	}
}

func TestTavorParserTypedTokens(t *testing.T) {
	var tok token.Token
	var err error

	// RangeInt
	tok, err = ParseTavor(strings.NewReader(
		"$Spec Int\nSTART = Spec\n",
	))
	Nil(t, err)
	Equal(t, tok, primitives.NewScope(primitives.NewRangeInt(0, math.MaxInt32)))

	tok, err = ParseTavor(strings.NewReader(
		"$Spec Int = from: 2,\nto: 10\nSTART = Spec\n",
	))
	Nil(t, err)
	Equal(t, tok, primitives.NewScope(primitives.NewRangeInt(2, 10)))

	tok, err = ParseTavor(strings.NewReader(
		"$Spec Int = from: 2,\nto: 10,\nstep: 2\nSTART = Spec\n",
	))
	Nil(t, err)
	Equal(t, tok, primitives.NewScope(primitives.NewRangeIntWithStep(2, 10, 2)))

	tok, err = ParseTavor(strings.NewReader(
		"$Spec Int = to: 10,\nstep: 2\nSTART = Spec\n",
	))
	Nil(t, err)
	Equal(t, tok, primitives.NewScope(primitives.NewRangeIntWithStep(0, 10, 2)))

	tok, err = ParseTavor(strings.NewReader(
		"$Spec Int = from: 2,\nstep: 2\nSTART = Spec\n",
	))
	Nil(t, err)
	Equal(t, tok, primitives.NewScope(primitives.NewRangeIntWithStep(2, math.MaxInt32, 2)))

	// Sequence
	{
		s := sequences.NewSequence(1, 1)
		tok, err = ParseTavor(strings.NewReader(
			"$Spec Sequence\nSTART = $Spec.Next\n",
		))
		Nil(t, err)
		Equal(t, tok, lists.NewAll(
			s.ResetItem(),
			primitives.NewScope(s.Item()),
		))

		s = sequences.NewSequence(2, 1)
		tok, err = ParseTavor(strings.NewReader(
			"$Spec Sequence = start: 2\nSTART = $Spec.Next\n",
		))
		Nil(t, err)
		Equal(t, tok, lists.NewAll(
			s.ResetItem(),
			primitives.NewScope(s.Item()),
		))

		s = sequences.NewSequence(1, 3)
		tok, err = ParseTavor(strings.NewReader(
			"$Spec Sequence = step: 3\nSTART = $Spec.Next\n",
		))
		Nil(t, err)
		Equal(t, tok, lists.NewAll(
			s.ResetItem(),
			primitives.NewScope(s.Item()),
		))

		s = sequences.NewSequence(1, 1)
		tok, err = ParseTavor(strings.NewReader(
			"$Spec Sequence\nSTART = $Spec.Existing\n",
		))
		Nil(t, err)
		Equal(t, tok, lists.NewAll(
			s.ResetItem(),
			primitives.NewScope(s.ExistingItem(nil)),
		))

		s = sequences.NewSequence(1, 1)
		tok, err = ParseTavor(strings.NewReader(
			"$Spec Sequence\nSTART = $Spec.Reset\n",
		))
		Nil(t, err)
		Equal(t, tok, lists.NewAll(
			s.ResetItem(),
			primitives.NewScope(s.ResetItem()),
		))
	}
}

func TestTavorParserExpressions(t *testing.T) {
	var tok token.Token
	var err error

	// token use in expression
	{
		tok, err = ParseTavor(strings.NewReader(`
			START = ${A}
			A = "a"
		`))
		Nil(t, err)
		Equal(t, tok, primitives.NewScope(primitives.NewConstantString("a")))
	}

	// variable use in expression
	{
		tok, err = ParseTavor(strings.NewReader(`
			START = "a"<A> ${A}
		`))
		Nil(t, err)
		v := variables.NewVariable("A", primitives.NewConstantString("a"))
		Equal(t, tok, primitives.NewScope(lists.NewAll(
			v,
			variables.NewVariableValue(v),
		)))
	}

	// token attribute use in expression
	{
		tok, err = ParseTavor(strings.NewReader(`
			START = "a"<A> ${A.Value}
		`))
		Nil(t, err)
		v := variables.NewVariable("A", primitives.NewConstantString("a"))
		Equal(t, tok, primitives.NewScope(lists.NewAll(
			v,
			variables.NewVariableValue(v),
		)))
	}

	// simple expression
	{
		s := sequences.NewSequence(1, 1)
		tok, err = ParseTavor(strings.NewReader(
			"$Spec Sequence\nSTART = ${Spec.Next}\n",
		))
		Nil(t, err)
		Equal(t, tok, lists.NewAll(
			s.ResetItem(),
			primitives.NewScope(s.Item()),
		))
	}

	// plus operator
	tok, err = ParseTavor(strings.NewReader(
		"START = ${1 + 2}\n",
	))
	Nil(t, err)
	Equal(t, tok, primitives.NewScope(expressions.NewAddArithmetic(
		primitives.NewConstantInt(1),
		primitives.NewConstantInt(2),
	)))

	tok, err = ParseTavor(strings.NewReader(`
		START = ${A + B}
		A = 1
		B = 2
	`))
	Nil(t, err)
	Equal(t, tok, primitives.NewScope(expressions.NewAddArithmetic(
		primitives.NewScope(primitives.NewConstantInt(1)),
		primitives.NewScope(primitives.NewConstantInt(2)),
	)))

	// sub operator
	tok, err = ParseTavor(strings.NewReader(
		"START = ${1 - 2}\n",
	))
	Nil(t, err)
	Equal(t, tok, primitives.NewScope(expressions.NewSubArithmetic(
		primitives.NewConstantInt(1),
		primitives.NewConstantInt(2),
	)))

	// mul operator
	tok, err = ParseTavor(strings.NewReader(
		"START = ${1 * 2}\n",
	))
	Nil(t, err)
	Equal(t, tok, primitives.NewScope(expressions.NewMulArithmetic(
		primitives.NewConstantInt(1),
		primitives.NewConstantInt(2),
	)))

	// div operator
	tok, err = ParseTavor(strings.NewReader(
		"START = ${1 / 2}\n",
	))
	Nil(t, err)
	Equal(t, tok, primitives.NewScope(expressions.NewDivArithmetic(
		primitives.NewConstantInt(1),
		primitives.NewConstantInt(2),
	)))

	// nested operator
	tok, err = ParseTavor(strings.NewReader(
		"START = ${1 + 2 + 3}\n",
	))
	Nil(t, err)
	Equal(t, tok, primitives.NewScope(expressions.NewAddArithmetic(
		primitives.NewConstantInt(1),
		expressions.NewAddArithmetic(
			primitives.NewConstantInt(2),
			primitives.NewConstantInt(3),
		),
	)))

	// mixed operator
	{
		s := sequences.NewSequence(1, 1)
		tok, err = ParseTavor(strings.NewReader(
			"$Spec Sequence\nSTART = ${Spec.Next + 1}\n",
		))
		Nil(t, err)
		Equal(t, tok, lists.NewAll(
			s.ResetItem(),
			primitives.NewScope(expressions.NewAddArithmetic(
				s.Item(),
				primitives.NewConstantInt(1),
			)),
		))
		Equal(t, "2", tok.String())
	}

	// path operator
	{
		tok, err = ParseTavor(strings.NewReader(`
				START = Pairs "->" Path

				Path = ${Pairs path from (2) over (e.Item(0)) connect by (e.Item(1)) without (0)}

				Pairs = (,
					(1 0),
					(3 1),
					(2 3),
				)
			`))
		Nil(t, err)
		Equal(t, "103123->231", tok.String())
	}
}

func TestTavorParserAndCuriousCaseOfFuzzing(t *testing.T) {
	var tok token.Token
	var err error

	// Additional forward declaration check
	tok, err = ParseTavor(strings.NewReader(
		"START = Token\nToken = 123\n",
	))
	Nil(t, err)
	Equal(t, tok, primitives.NewScope(primitives.NewConstantInt(123)))

	// double embedded forward token all the way
	tok, err = ParseTavor(strings.NewReader("A = B B\nB = 1\nSTART = A\n"))
	Nil(t, err)
	Equal(t, tok, primitives.NewScope(lists.NewAll(
		primitives.NewScope(primitives.NewConstantInt(1)),
		primitives.NewScope(primitives.NewConstantInt(1)),
	)))

	// Token attribute forward usage
	tok, err = ParseTavor(strings.NewReader(
		"START = $int.Value\n$int Int\n",
	))
	Nil(t, err)
	Equal(t, tok, primitives.NewScope(primitives.NewRangeInt(0, math.MaxInt32)))

	// Tokens should be cloned so they are different internally
	{
		tok, err = ParseTavor(strings.NewReader(
			"Token = 1 | 2\nSTART = Token Token\n",
		))
		Nil(t, err)
		Equal(t, tok, primitives.NewScope(lists.NewAll(
			primitives.NewScope(lists.NewOne(primitives.NewConstantInt(1), primitives.NewConstantInt(2))),
			primitives.NewScope(lists.NewOne(primitives.NewConstantInt(1), primitives.NewConstantInt(2))),
		)))

		va, _ := tok.(*primitives.Scope).InternalGet().(token.ListToken).Get(0)
		a := va.(*primitives.Scope).InternalGet().(*lists.One)
		vb, _ := tok.(*primitives.Scope).InternalGet().(token.ListToken).Get(1)
		b := vb.(*primitives.Scope).InternalGet().(*lists.One)

		True(t, Exactly(t, a, b))
		NotEqual(t, fmt.Sprintf("%p", a), fmt.Sprintf("%p", b))
	}

	// Correct sequence behaviour
	{
		tok, err = ParseTavor(strings.NewReader(`
			$Id Sequence = start: 2,
				step: 2

			START = $Id.Reset $Id.Next $Id.Existing
		`))
		Nil(t, err)

		strat := strategy.NewRandomStrategy(tok)
		ch, err := strat.Fuzz(test.NewRandTest(1))
		Nil(t, err)

		for i := range ch {
			Equal(t, "22", tok.String())

			ch <- i
		}
	}

	// Correct list behaviour
	{
		tok, err = ParseTavor(strings.NewReader(`
			A = +2(1)

			START = $A.Count A
		`))
		Nil(t, err)

		strat := strategy.NewRandomStrategy(tok)
		ch, err := strat.Fuzz(test.NewRandTest(1))
		Nil(t, err)

		for i := range ch {
			Equal(t, "211", tok.String())

			ch <- i
		}
	}

	// Attributes in repeats
	{
		tok, err = ParseTavor(strings.NewReader(`
			As = +3("a")
			Bs = +$As.Count("b")
			START = As Bs
		`))
		Nil(t, err)

		Equal(t, "aaabbb", tok.String())
	}
	{
		tok, err = ParseTavor(strings.NewReader(`
			As = +3("a")
			Bs = +2,$As.Count("b")
			START = As Bs
		`))
		Nil(t, err)

		Equal(t, "aaabb", tok.String())
	}

	// save variable
	{
		tok, err = ParseTavor(strings.NewReader(`
			START = "abc"<v> $v.Value
		`))
		Nil(t, err)

		Equal(t, "abcabc", tok.String())
	}
	{
		tok, err = ParseTavor(strings.NewReader(`
			START = "abc"<=v> $v.Value
		`))
		Nil(t, err)

		Equal(t, "abc", tok.String())
	}

	// Save variable scope and variable usage in expression
	{
		tok, err = ParseTavor(strings.NewReader(`
			START = Number<=a> Number<=b>,
			a " + " b " = " ${a.Value + b.Value} "\n",
			a " * " b " = " ${a.Value * b.Value} "\n"

			$Number Int = from: 1,
			              to:   2
		`))
		Nil(t, err)

		Equal(t, "1 + 1 = 2\n1 * 1 = 1\n", tok.String())
	}

	// Special path, with a loop which needs a variable to work
	{
		tok, err = ParseTavor(strings.NewReader(`
			$Literal Sequence = start: 2,
					step: 2

			ExistingLiteralAnd = 0,
				| 1,
				| ${Literal.Existing not in (AndCycle)}

			AndCycle = ${andList.Reference path from (andLiteral) over (e.Item(0)) connect by (e.Item(1), e.Item(2)) without (0, 1)}

			Ands = +4(And)
			And = $Literal.Next<andLiteral> " " ExistingLiteralAnd " " ExistingLiteralAnd "\n"

			START = Ands<andList>
		`))
		Nil(t, err)

		Equal(t, "2 0 0\n2 0 0\n2 0 0\n2 0 0\n", tok.String())
	}
}

func TestTavorParserLoops(t *testing.T) {
	var tok token.Token
	var err error

	// simplest valid loop
	tok, err = ParseTavor(strings.NewReader(`
		A = A | 1

		START = A
	`))
	Nil(t, err)
	{
		Equal(t, tok, primitives.NewScope(lists.NewOne(
			primitives.NewScope(lists.NewOne(
				primitives.NewScope(primitives.NewConstantInt(1)),
				primitives.NewConstantInt(1),
			)),
			primitives.NewConstantInt(1),
		)))

		Equal(t, "1", tok.String())
	}

	tok, err = ParseTavor(strings.NewReader(`
		A = A 1 | 2

		START = A
	`))
	Nil(t, err)
	{
		Equal(t, tok, primitives.NewScope(lists.NewOne(
			lists.NewAll(
				primitives.NewScope(lists.NewOne(
					lists.NewAll(
						primitives.NewScope(primitives.NewConstantInt(2)),
						primitives.NewConstantInt(1),
					),
					primitives.NewConstantInt(2),
				)),
				primitives.NewConstantInt(1),
			),
			primitives.NewConstantInt(2),
		)))

		Equal(t, "211", tok.String())
	}

	// optional loop
	tok, err = ParseTavor(strings.NewReader(`
		A = ?(A) 1

		START = A
	`))
	Nil(t, err)
	{
		Equal(t, tok, primitives.NewScope(lists.NewAll(
			constraints.NewOptional(
				primitives.NewScope(lists.NewAll(
					constraints.NewOptional(
						primitives.NewScope(primitives.NewConstantInt(1)),
					),
					primitives.NewConstantInt(1),
				)),
			),
			primitives.NewConstantInt(1),
		)))

		Equal(t, "111", tok.String())
	}

	// One loop
	tok, err = ParseTavor(strings.NewReader(`
		A = (A | 2) 1

		START = A
	`))
	Nil(t, err)
	{
		Equal(t, tok, primitives.NewScope(lists.NewAll(
			lists.NewOne(
				primitives.NewScope(lists.NewAll(
					lists.NewOne(
						primitives.NewScope(lists.NewAll(
							primitives.NewConstantInt(2),
							primitives.NewConstantInt(1),
						)),
						primitives.NewConstantInt(2),
					),
					primitives.NewConstantInt(1),
				)),
				primitives.NewConstantInt(2),
			),
			primitives.NewConstantInt(1),
		)))

		Equal(t, "2111", tok.String())
	}

	// two loops
	{
		tok, err = ParseTavor(strings.NewReader(`
			A = (A | 2) 1 ?(A) 3

			START = A
		`))
		Nil(t, err)
		Equal(t, tok, primitives.NewScope(lists.NewAll(
			lists.NewOne(
				primitives.NewScope(lists.NewAll(
					lists.NewOne(
						primitives.NewScope(lists.NewAll(
							primitives.NewConstantInt(2),
							primitives.NewConstantInt(1),
							primitives.NewConstantInt(3),
						)),
						primitives.NewConstantInt(2),
					),
					primitives.NewConstantInt(1),
					constraints.NewOptional(
						primitives.NewScope(lists.NewAll(
							primitives.NewConstantInt(2),
							primitives.NewConstantInt(1),
							primitives.NewConstantInt(3),
						)),
					),
					primitives.NewConstantInt(3),
				)),
				primitives.NewConstantInt(2),
			),
			primitives.NewConstantInt(1),
			constraints.NewOptional(
				primitives.NewScope(lists.NewAll(
					lists.NewOne(
						primitives.NewScope(lists.NewAll(
							primitives.NewConstantInt(2),
							primitives.NewConstantInt(1),
							primitives.NewConstantInt(3),
						)),
						primitives.NewConstantInt(2),
					),
					primitives.NewConstantInt(1),
					constraints.NewOptional(
						primitives.NewScope(lists.NewAll(
							primitives.NewConstantInt(2),
							primitives.NewConstantInt(1),
							primitives.NewConstantInt(3),
						)),
					),
					primitives.NewConstantInt(3),
				)),
			),
			primitives.NewConstantInt(3),
		)))

		Equal(t, "213121331213121333", tok.String())
	}

	// loop with token loop
	tok, err = ParseTavor(strings.NewReader(`
			B = A
			C = B

			A = ?(C) 1

			START = A
		`))
	Nil(t, err)
	{
		Equal(t, tok, primitives.NewScope(lists.NewAll(
			constraints.NewOptional(
				primitives.NewScope(lists.NewAll(
					constraints.NewOptional(
						primitives.NewScope(primitives.NewConstantInt(1)),
					),
					primitives.NewConstantInt(1),
				)),
			),
			primitives.NewConstantInt(1),
		)))

		Equal(t, "111", tok.String())
	}

	// repeated forward one
	tok, err = ParseTavor(strings.NewReader(`
			Action = SetParameter,
			       | GetParameter

			SetParameter = "setParam"
			GetParameter = "getParam" ("param 1" | "param 2")

			START = +(Action)
		`))
	Nil(t, err)
	{
		Equal(t, tok, primitives.NewScope(lists.NewRepeat(
			primitives.NewScope(lists.NewOne(
				primitives.NewScope(primitives.NewConstantString("setParam")),
				primitives.NewScope(lists.NewAll(
					primitives.NewConstantString("getParam"),
					lists.NewOne(
						primitives.NewConstantString("param 1"),
						primitives.NewConstantString("param 2"),
					),
				)),
			)),
			1,
			2,
		)))

		Equal(t, "setParam", tok.String())
	}

	// endless loop with exit through optional
	{
		tok, err = ParseTavor(strings.NewReader(`
			C = A
			B = C

			A = (B | 1)(B | 2) ?(B)

			START = A
		`))
		Nil(t, err)
		Equal(t, tok, primitives.NewScope(lists.NewAll(
			lists.NewOne(
				primitives.NewScope(lists.NewAll(
					lists.NewOne(
						primitives.NewScope(lists.NewAll(
							primitives.NewConstantInt(1),
							primitives.NewConstantInt(2),
						)),
						primitives.NewConstantInt(1),
					),
					lists.NewOne(
						primitives.NewScope(lists.NewAll(
							primitives.NewConstantInt(1),
							primitives.NewConstantInt(2),
						)),
						primitives.NewConstantInt(2),
					),
					constraints.NewOptional(
						primitives.NewScope(lists.NewAll(
							primitives.NewConstantInt(1),
							primitives.NewConstantInt(2),
						)),
					),
				)),
				primitives.NewConstantInt(1),
			),
			lists.NewOne(
				primitives.NewScope(lists.NewAll(
					lists.NewOne(
						primitives.NewScope(lists.NewAll(
							primitives.NewConstantInt(1),
							primitives.NewConstantInt(2),
						)),
						primitives.NewConstantInt(1),
					),
					lists.NewOne(
						primitives.NewScope(lists.NewAll(
							primitives.NewConstantInt(1),
							primitives.NewConstantInt(2),
						)),
						primitives.NewConstantInt(2),
					),
					constraints.NewOptional(
						primitives.NewScope(lists.NewAll(
							primitives.NewConstantInt(1),
							primitives.NewConstantInt(2),
						)),
					),
				)),
				primitives.NewConstantInt(2),
			),
			constraints.NewOptional(primitives.NewScope(lists.NewAll(
				lists.NewOne(
					primitives.NewScope(lists.NewAll(
						primitives.NewConstantInt(1),
						primitives.NewConstantInt(2),
					)),
					primitives.NewConstantInt(1),
				),
				lists.NewOne(
					primitives.NewScope(lists.NewAll(
						primitives.NewConstantInt(1),
						primitives.NewConstantInt(2),
					)),
					primitives.NewConstantInt(2),
				),
				constraints.NewOptional(
					primitives.NewScope(lists.NewAll(
						primitives.NewConstantInt(1),
						primitives.NewConstantInt(2),
					)),
				),
			))),
		)))

		Equal(t, "121212121212121212", tok.String())
	}
}

func TestTavorParserCornerCases(t *testing.T) {
	// early usage used twice deeper in token
	{
		tok, err := ParseTavor(strings.NewReader(`
			c = ?(d)

			a = c
			b = c

			d =  "TEXT"

			START = a | b
		`))
		Nil(t, err)
		Equal(t, tok, primitives.NewScope(lists.NewOne(
			primitives.NewScope(constraints.NewOptional(primitives.NewScope(primitives.NewConstantString("TEXT")))),
			primitives.NewScope(constraints.NewOptional(primitives.NewScope(primitives.NewConstantString("TEXT")))),
		)))

		Equal(t, "TEXT", tok.String())
	}
	// early usage used twice
	{
		tok, err := ParseTavor(strings.NewReader(`
			c = d

			a = c
			b = c

			d =  "TEXT"

			START = a | b
		`))
		Nil(t, err)
		Equal(t, tok, primitives.NewScope(lists.NewOne(
			primitives.NewScope(primitives.NewConstantString("TEXT")),
			primitives.NewScope(primitives.NewConstantString("TEXT")),
		)))

		Equal(t, "TEXT", tok.String())
	}
	// early usage used twice even deeper in token
	{
		tok, err := ParseTavor(strings.NewReader(`
			c = ?(?(d))

			a = c
			b = c

			d =  "TEXT"

			START = a | b
		`))
		Nil(t, err)
		Equal(t, tok, primitives.NewScope(lists.NewOne(
			primitives.NewScope(constraints.NewOptional(constraints.NewOptional(primitives.NewScope(primitives.NewConstantString("TEXT"))))),
			primitives.NewScope(constraints.NewOptional(constraints.NewOptional(primitives.NewScope(primitives.NewConstantString("TEXT"))))),
		)))

		Equal(t, "TEXT", tok.String())
	}
	// three times is the charme if you are doing early usage...
	{
		tok, err := ParseTavor(strings.NewReader(`
			START = B  B  B

			B = "B"
		`))
		Nil(t, err)
		Equal(t, tok, primitives.NewScope(lists.NewAll(
			primitives.NewScope(primitives.NewConstantString("B")),
			primitives.NewScope(primitives.NewConstantString("B")),
			primitives.NewScope(primitives.NewConstantString("B")),
		)))

		Equal(t, "BBB", tok.String())
	}
	{
		tok, err := ParseTavor(strings.NewReader(`
			START = B  B  B

			B = 1 2
		`))
		Nil(t, err)
		Equal(t, tok, primitives.NewScope(lists.NewAll(
			primitives.NewScope(lists.NewAll(primitives.NewConstantInt(1), primitives.NewConstantInt(2))),
			primitives.NewScope(lists.NewAll(primitives.NewConstantInt(1), primitives.NewConstantInt(2))),
			primitives.NewScope(lists.NewAll(primitives.NewConstantInt(1), primitives.NewConstantInt(2))),
		)))

		Equal(t, "121212", tok.String())
	}
	{
		tok, err := ParseTavor(strings.NewReader(`
			START = V

			$V Int = from: 1,
				to: 1
		`))
		Nil(t, err)
		Equal(t, tok, primitives.NewScope(primitives.NewRangeInt(1, 1)))

		Equal(t, "1", tok.String())
	}
}

func TestTavorParserCharacterClasses(t *testing.T) {
	{
		tok, err := ParseTavor(strings.NewReader(`
			START = [123]
		`))
		Nil(t, err)
		Equal(t, tok, primitives.NewScope(primitives.NewCharacterClass("123")))

		Equal(t, "1", tok.String())
	}
	{
		tok, err := ParseTavor(strings.NewReader(`
			START = [\w]
		`))
		Nil(t, err)
		Equal(t, tok, primitives.NewScope(primitives.NewCharacterClass(`\w`)))

		Equal(t, "0", tok.String())
	}
	{
		// do not parse spaces between character class brackets
		tok, err := ParseTavor(strings.NewReader(`
			START = [ ]
		`))
		Nil(t, err)
		Equal(t, tok, primitives.NewScope(primitives.NewCharacterClass(` `)))

		Equal(t, " ", tok.String())
	}
}

func TestTavorParserVariables(t *testing.T) {
	// simple save and value
	{
		tok, err := ParseTavor(strings.NewReader(`
			START = Save<var> Print

			Save = "text"

			Print = $var.Value
		`))
		Nil(t, err)
		variable := variables.NewVariable("var", primitives.NewScope(primitives.NewConstantString("text")))
		Equal(t, tok, primitives.NewScope(lists.NewAll(
			variable,
			primitives.NewScope(variables.NewVariableValue(variable)),
		)))

		Equal(t, "texttext", tok.String())
	}
	// correct scope of variables
	{
		tok, err := ParseTavor(strings.NewReader(`
			START = 1<var> Print 2<var> Print

			Print = $var.Value
		`))
		Nil(t, err)

		v1 := variables.NewVariable("var", primitives.NewConstantInt(1))
		v2 := variables.NewVariable("var", primitives.NewConstantInt(2))

		Equal(t, tok, primitives.NewScope(lists.NewAll(
			v1, primitives.NewScope(variables.NewVariableValue(v1)),
			v2, primitives.NewScope(variables.NewVariableValue(v2)),
		)))

		Equal(t, "1122", tok.String())
	}
	// forward variable declaration
	{
		tok, err := ParseTavor(strings.NewReader(`
			A = b

			START = "b"<b> A
		`))
		Nil(t, err)

		Equal(t, "bb", tok.String())
	}
	// forward variable declaration over path
	{
		tok, err := ParseTavor(strings.NewReader(`
			A = c
			B = A

			START = "c"<c> B
		`))
		Nil(t, err)

		Equal(t, "cc", tok.String())
	}
	// forward embedded variable declaration over path
	{
		tok, err := ParseTavor(strings.NewReader(`
			B = $var.Count
			A = B

			START = A<var>
		`))
		Nil(t, err)

		Equal(t, "1", tok.String())
	}
	// not in with variables
	{
		tok, err := ParseTavor(strings.NewReader(`
			$Literal Sequence

			And = $Literal.Next<x> " " ${Literal.Existing not in (x)} " " ${Literal.Existing not in (x)} "\n"

			START = $Literal.Reset $Literal.Next "\n" And And
		`))
		Nil(t, err)

		/* TODO if this example finally is correct.... do the token graph
		seq := sequences.NewSequence(1, 1)

		Equal(t, tok, primitives.NewScope(lists.NewAll(
			seq.ResetItem(),
			seq.Item(),

		))*/

		Equal(t, "2\n1 1 1\n1 1 1\n", tok.String())
	}
}

func TestTavorParserIfElseIfElsedd(t *testing.T) {
	// basic if, else if and else
	{
		tok, err := ParseTavor(strings.NewReader(`
			START = Choose<var> Print

			Choose = 1 | 2 | 3

			Print = {if var.Value == 1} "var is one" {else if var.Value == 2} "var is two" {else} "var is three" {endif}
		`))
		Nil(t, err)

		variable, _ := tok.(*primitives.Scope).InternalGet().(*lists.All).InternalGet(0)
		one := variable.(*variables.Variable).InternalGet().(*primitives.Scope).InternalGet()

		nOne := primitives.NewScope(lists.NewOne(
			primitives.NewConstantInt(1),
			primitives.NewConstantInt(2),
			primitives.NewConstantInt(3),
		))
		nVariable := variables.NewVariable("var", nOne)

		var ll token.Token = primitives.NewScope(lists.NewAll(
			nVariable,
			primitives.NewScope(conditions.NewIf(
				conditions.IfPair{ // TODO FIXME AND FIXME!!!!!! allow unrolling of IfPairs and BooleanEquals and pretty much all in token/conditions
					Head: conditions.NewBooleanEqual(primitives.NewPointer(primitives.NewTokenPointer(variables.NewVariableValue(nVariable))), primitives.NewConstantInt(1)),
					Body: primitives.NewConstantString("var is one"),
				},
				conditions.IfPair{
					Head: conditions.NewBooleanEqual(primitives.NewPointer(primitives.NewTokenPointer(variables.NewVariableValue(nVariable))), primitives.NewConstantInt(2)),
					Body: primitives.NewConstantString("var is two"),
				},
				conditions.IfPair{
					Head: conditions.NewBooleanTrue(),
					Body: primitives.NewConstantString("var is three"),
				},
			)),
		))

		Equal(t, tok, ll)

		Equal(t, "1var is one", tok.String())

		Equal(t, 3, one.Permutations())

		Nil(t, one.Permutation(2))
		Equal(t, "2var is two", tok.String())

		Nil(t, one.Permutation(3))
		Equal(t, "3var is three", tok.String())
	}
	// continued definition
	{
		tok, err := ParseTavor(strings.NewReader(`
			START = 1<var> {if var.Value == 1} 2 {endif} 3
		`))
		Nil(t, err)

		nVariable := variables.NewVariable("var", primitives.NewConstantInt(1))

		Equal(t, tok, primitives.NewScope(lists.NewAll(
			nVariable,
			conditions.NewIf(
				conditions.IfPair{
					Head: conditions.NewBooleanEqual(variables.NewVariableValue(nVariable), primitives.NewConstantInt(1)),
					Body: primitives.NewConstantInt(2),
				},
			),
			primitives.NewConstantInt(3),
		)))

		Equal(t, "123", tok.String())
		Equal(t, 1, tok.Permutations())
	}
	// if defined
	{
		tok, err := ParseTavor(strings.NewReader(`
			START = Token Print

			Token = "abc"<var> Print

			Print = {if defined var} "var is defined" {else} "var is not defined" {endif}
		`))
		Nil(t, err)

		nVariable := variables.NewVariable("var", primitives.NewConstantString("abc"))

		definedScope := token.NewVariableScope()
		definedScope = definedScope.Push().Push()
		definedScope.Set("var", nVariable)
		definedScope = definedScope.Push()

		notDefinedScope := token.NewVariableScope().Push().Push()

		Equal(t, tok, primitives.NewScope(lists.NewAll(
			primitives.NewScope(lists.NewAll(
				nVariable,
				primitives.NewScope(conditions.NewIf(
					conditions.IfPair{
						Head: conditions.NewVariableDefined("var", definedScope),
						Body: primitives.NewConstantString("var is defined"),
					},
					conditions.IfPair{
						Head: conditions.NewBooleanTrue(),
						Body: primitives.NewConstantString("var is not defined"),
					},
				)),
			)),
			primitives.NewScope(conditions.NewIf(
				conditions.IfPair{
					Head: conditions.NewVariableDefined("var", notDefinedScope),
					Body: primitives.NewConstantString("var is defined"),
				},
				conditions.IfPair{
					Head: conditions.NewBooleanTrue(),
					Body: primitives.NewConstantString("var is not defined"),
				},
			)),
		)))

		Equal(t, "abcvar is definedvar is not defined", tok.String())
		Equal(t, 1, tok.Permutations())
	}
}
