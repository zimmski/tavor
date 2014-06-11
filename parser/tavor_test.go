package parser

import (
	"strings"
	"testing"

	. "github.com/stretchr/testify/assert"

	"github.com/zimmski/tavor/test"
	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/aggregates"
	"github.com/zimmski/tavor/token/constraints"
	"github.com/zimmski/tavor/token/expressions"
	"github.com/zimmski/tavor/token/lists"
	"github.com/zimmski/tavor/token/primitives"
	"github.com/zimmski/tavor/token/sequences"
)

func TestTavorParseErrors(t *testing.T) {
	var tok token.Token
	var err error

	// empty file
	tok, err = ParseTavor(strings.NewReader(""))
	Equal(t, ParseErrorNoStart, err.(*ParserError).Type)
	Nil(t, tok)

	// empty file
	tok, err = ParseTavor(strings.NewReader("START = 123"))
	Equal(t, ParseErrorNewLineNeeded, err.(*ParserError).Type)
	Nil(t, tok)

	// new line before =
	tok, err = ParseTavor(strings.NewReader("START \n= 123\n"))
	Equal(t, ParseErrorEarlyNewLine, err.(*ParserError).Type)
	Nil(t, tok)

	// expect =
	tok, err = ParseTavor(strings.NewReader("START 123\n"))
	Equal(t, ParseErrorExpectRune, err.(*ParserError).Type)
	Nil(t, tok)

	// new line after =
	tok, err = ParseTavor(strings.NewReader("START = \n123\n"))
	Equal(t, ParseErrorEmptyTokenDefinition, err.(*ParserError).Type)
	Nil(t, tok)

	// invalid token name. does not start with letter
	tok, err = ParseTavor(strings.NewReader("3TART = 123\n"))
	Equal(t, ParseErrorInvalidTokenName, err.(*ParserError).Type)
	Nil(t, tok)

	// unused token
	tok, err = ParseTavor(strings.NewReader("START = 123\nNumber = 123\n"))
	Equal(t, ParseErrorUnusedToken, err.(*ParserError).Type)
	Nil(t, tok)

	// non-terminated string
	tok, err = ParseTavor(strings.NewReader("START = \"non-terminated string\n"))
	Equal(t, ParseErrorNonTerminatedString, err.(*ParserError).Type)
	Nil(t, tok)

	// token already defined
	tok, err = ParseTavor(strings.NewReader("START = 123\nSTART = 456\n"))
	Equal(t, ParseErrorTokenAlreadyDefined, err.(*ParserError).Type)
	Nil(t, tok)

	// token is not defined
	tok, err = ParseTavor(strings.NewReader("START = Token\n"))
	Equal(t, ParseErrorTokenNotDefined, err.(*ParserError).Type)
	Nil(t, tok)

	// unexpected multi line token termination
	tok, err = ParseTavor(strings.NewReader("Hello = 1,\n\n"))
	Equal(t, ParseErrorUnexpectedTokenDefinitionTermination, err.(*ParserError).Type)
	Nil(t, tok)

	// unexpected continue of multi line token
	tok, err = ParseTavor(strings.NewReader("Hello = 1,2\n"))
	Equal(t, ParseErrorExpectRune, err.(*ParserError).Type)
	Nil(t, tok)

	// unknown token attribute
	tok, err = ParseTavor(strings.NewReader("Token = 123\nSTART = $Token.yeah\n"))
	Equal(t, ParseErrorUnknownTokenAttribute, err.(*ParserError).Type)
	Nil(t, tok)

	// unknown token attribute
	tok, err = ParseTavor(strings.NewReader("Token = 123\nSTART = Token $Token.Count\n"))
	Equal(t, ParseErrorUnknownTokenAttribute, err.(*ParserError).Type)
	Nil(t, tok)

	// token not defined for token attribute
	tok, err = ParseTavor(strings.NewReader("START = $Token.Count\n"))
	Equal(t, ParseErrorTokenNotDefined, err.(*ParserError).Type)
	Nil(t, tok)

	// special token already defined
	tok, err = ParseTavor(strings.NewReader("START = 123\n$START = 456\n"))
	Equal(t, ParseErrorTokenAlreadyDefined, err.(*ParserError).Type)
	Nil(t, tok)

	// expect = in special token
	tok, err = ParseTavor(strings.NewReader("$START 123\n"))
	Equal(t, ParseErrorExpectRune, err.(*ParserError).Type)
	Nil(t, tok)

	// expect identifier in special token
	tok, err = ParseTavor(strings.NewReader("$START = 123\n"))
	Equal(t, ParseErrorExpectRune, err.(*ParserError).Type)
	Nil(t, tok)

	// expect : in special token
	tok, err = ParseTavor(strings.NewReader("$START = argument\n"))
	Equal(t, ParseErrorExpectRune, err.(*ParserError).Type)
	Nil(t, tok)

	// expect valid argument value in special token
	tok, err = ParseTavor(strings.NewReader("$START = argument:\n"))
	Equal(t, ParseErrorInvalidArgumentValue, err.(*ParserError).Type)
	Nil(t, tok)

	// expect no eof after argument value in special token
	tok, err = ParseTavor(strings.NewReader("$START = argument:value"))
	Equal(t, ParseErrorNewLineNeeded, err.(*ParserError).Type)
	Nil(t, tok)

	// expect new line in special token
	tok, err = ParseTavor(strings.NewReader("$START = argument:value$"))
	Equal(t, ParseErrorExpectRune, err.(*ParserError).Type)
	Nil(t, tok)

	// undefined type argument special token
	tok, err = ParseTavor(strings.NewReader("$START = hey: ok\n"))
	Equal(t, ParseErrorUnknownTypeForSpecialToken, err.(*ParserError).Type)
	Nil(t, tok)

	// unknown type argument special token
	tok, err = ParseTavor(strings.NewReader("$START = type: ok\n"))
	Equal(t, ParseErrorUnknownSpecialTokenType, err.(*ParserError).Type)
	Nil(t, tok)

	// unknown special token argument
	tok, err = ParseTavor(strings.NewReader("$START = type: Int,\nok: value\n"))
	Equal(t, ParseErrorUnknownSpecialTokenArgument, err.(*ParserError).Type)
	Nil(t, tok)

	// missing arguments for special token Int
	tok, err = ParseTavor(strings.NewReader("$START = type: Int,\nto:123\n"))
	Equal(t, ParseErrorMissingSpecialTokenArgument, err.(*ParserError).Type)
	Nil(t, tok)

	tok, err = ParseTavor(strings.NewReader("$START = type: Int,\nfrom:123\n"))
	Equal(t, ParseErrorMissingSpecialTokenArgument, err.(*ParserError).Type)
	Nil(t, tok)

	// invalid arguments for special token Int
	tok, err = ParseTavor(strings.NewReader("$START = type: Int,\nfrom:abc,\nto:123\n"))
	Equal(t, ParseErrorInvalidArgumentValue, err.(*ParserError).Type)
	Nil(t, tok)

	tok, err = ParseTavor(strings.NewReader("$START = type: Int,\nfrom:123,\nto:abc\n"))
	Equal(t, ParseErrorInvalidArgumentValue, err.(*ParserError).Type)
	Nil(t, tok)

	// invalid arguments for special token Sequence
	tok, err = ParseTavor(strings.NewReader("$START = type: Sequence,\nstart:abc\n"))
	Equal(t, ParseErrorInvalidArgumentValue, err.(*ParserError).Type)
	Nil(t, tok)

	tok, err = ParseTavor(strings.NewReader("$START = type: Sequence,\nstep:abc\n"))
	Equal(t, ParseErrorInvalidArgumentValue, err.(*ParserError).Type)
	Nil(t, tok)

	// empty expression
	tok, err = ParseTavor(strings.NewReader("START = ${}\n"))
	Equal(t, ParseErrorEmptyExpressionIsInvalid, err.(*ParserError).Type)
	Nil(t, tok)

	// open expression
	tok, err = ParseTavor(strings.NewReader("$Spec = type: Sequence\nSTART = ${Spec.Next\n"))
	Equal(t, ParseErrorExpectRune, err.(*ParserError).Type)
	Nil(t, tok)

	// missing operator expression term
	tok, err = ParseTavor(strings.NewReader("$Spec = type: Sequence\nSTART = ${Spec.Next +}\n"))
	Equal(t, ParseErrorExpectedExpressionTerm, err.(*ParserError).Type)
	Nil(t, tok)

	/*
		TODO this can maybe never happen as we do everything in one pass
		so we do not know that $List must implement lists.List

		// wrong token type because of earlier usage
		tok, err = ParseTavor(strings.NewReader("START = $List.Count\nList = 123"))
		panic(err)
		Equal(t, ParseErrorExpectedExpressionTerm, err.(*ParserError).Type)
		Nil(t, tok)
	*/
}

func TestTavorParserSimple(t *testing.T) {
	var tok token.Token
	var err error

	// constant integer
	tok, err = ParseTavor(strings.NewReader("START = 123\n"))
	Nil(t, err)
	Equal(t, tok, primitives.NewConstantInt(123))

	// single line comment
	tok, err = ParseTavor(strings.NewReader("// hello\nSTART = 123\n"))
	Nil(t, err)
	Equal(t, tok, primitives.NewConstantInt(123))

	// single line multi line comment
	tok, err = ParseTavor(strings.NewReader("/* hello */\nSTART = 123\n"))
	Nil(t, err)
	Equal(t, tok, primitives.NewConstantInt(123))

	// multi line multi line comment
	tok, err = ParseTavor(strings.NewReader("/*\nh\ne\nl\nl\no\n*/\nSTART = 123\n"))
	Nil(t, err)
	Equal(t, tok, primitives.NewConstantInt(123))

	// inline comment
	tok, err = ParseTavor(strings.NewReader("START /* ok */= /* or so */ 123\n"))
	Nil(t, err)
	Equal(t, tok, primitives.NewConstantInt(123))

	// constant string
	tok, err = ParseTavor(strings.NewReader("START = \"abc\"\n"))
	Nil(t, err)
	Equal(t, tok, primitives.NewConstantString("abc"))

	// constant string with whitespaces and epic chars
	tok, err = ParseTavor(strings.NewReader("START = \"a b c !\\n\\\"$%&/\"\n"))
	Nil(t, err)
	Equal(t, tok, primitives.NewConstantString("a b c !\n\"$%&/"))

	// concatination
	tok, err = ParseTavor(strings.NewReader("START = \"I am a constant string\" 123\n"))
	Nil(t, err)
	Equal(t, tok, lists.NewAll(
		primitives.NewConstantString("I am a constant string"),
		primitives.NewConstantInt(123),
	))

	// embed token
	tok, err = ParseTavor(strings.NewReader("Token=123\nSTART = Token\n"))
	Nil(t, err)
	Equal(t, tok, primitives.NewConstantInt(123))

	// embed over token
	tok, err = ParseTavor(strings.NewReader("Token=123\nAnotherToken = Token\nSTART = AnotherToken\n"))
	Nil(t, err)
	Equal(t, tok, primitives.NewConstantInt(123))

	// multi line token
	tok, err = ParseTavor(strings.NewReader("START = 1,\n2,\n3\n"))
	Nil(t, err)
	Equal(t, tok, lists.NewAll(
		primitives.NewConstantInt(1),
		primitives.NewConstantInt(2),
		primitives.NewConstantInt(3),
	))

	// Umläüt
	tok, err = ParseTavor(strings.NewReader("Umläüt=123\nSTART = Umläüt\n"))
	Nil(t, err)
	Equal(t, tok, primitives.NewConstantInt(123))
}

func TestTavorParserAlternationsAndGroupings(t *testing.T) {
	var tok token.Token
	var err error

	// simple alternation
	tok, err = ParseTavor(strings.NewReader("START = 1 | 2 | 3\n"))
	Nil(t, err)
	Equal(t, tok, lists.NewOne(
		primitives.NewConstantInt(1),
		primitives.NewConstantInt(2),
		primitives.NewConstantInt(3),
	))

	// concatinated alternation
	tok, err = ParseTavor(strings.NewReader("START = 1 | 2 3 | 4\n"))
	Nil(t, err)
	Equal(t, tok, lists.NewOne(
		primitives.NewConstantInt(1),
		lists.NewAll(
			primitives.NewConstantInt(2),
			primitives.NewConstantInt(3),
		),
		primitives.NewConstantInt(4),
	))

	// optional alternation
	tok, err = ParseTavor(strings.NewReader("START = | 2 | 3\n"))
	Nil(t, err)
	Equal(t, tok, constraints.NewOptional(lists.NewOne(
		primitives.NewConstantInt(2),
		primitives.NewConstantInt(3),
	)))

	tok, err = ParseTavor(strings.NewReader("START = 1 | | 3\n"))
	Nil(t, err)
	Equal(t, tok, constraints.NewOptional(lists.NewOne(
		primitives.NewConstantInt(1),
		primitives.NewConstantInt(3),
	)))

	tok, err = ParseTavor(strings.NewReader("START = 1 | 2 |\n"))
	Nil(t, err)
	Equal(t, tok, constraints.NewOptional(lists.NewOne(
		primitives.NewConstantInt(1),
		primitives.NewConstantInt(2),
	)))

	// alternation with embedded token
	tok, err = ParseTavor(strings.NewReader("Token = 2\nSTART = 1 | Token\n"))
	Nil(t, err)
	Equal(t, tok, lists.NewOne(
		primitives.NewConstantInt(1),
		primitives.NewConstantInt(2),
	))

	// simple group
	tok, err = ParseTavor(strings.NewReader("START = (1 2 3)\n"))
	Nil(t, err)
	Equal(t, tok, lists.NewAll(
		primitives.NewConstantInt(1),
		primitives.NewConstantInt(2),
		primitives.NewConstantInt(3),
	))

	// simple embedded group
	tok, err = ParseTavor(strings.NewReader("START = 0 (1 2 3) 4\n"))
	Nil(t, err)
	Equal(t, tok, lists.NewAll(
		primitives.NewConstantInt(0),
		lists.NewAll(
			primitives.NewConstantInt(1),
			primitives.NewConstantInt(2),
			primitives.NewConstantInt(3),
		),
		primitives.NewConstantInt(4),
	))

	// simple embedded or group
	tok, err = ParseTavor(strings.NewReader("START = 0 (1 | 2 | 3) 4\n"))
	Nil(t, err)
	Equal(t, tok, lists.NewAll(
		primitives.NewConstantInt(0),
		lists.NewOne(
			primitives.NewConstantInt(1),
			primitives.NewConstantInt(2),
			primitives.NewConstantInt(3),
		),
		primitives.NewConstantInt(4),
	))

	// Yo dog, I heard you like groups? so here is a group in a group
	tok, err = ParseTavor(strings.NewReader("START = (1 | (2 | 3)) | 4\n"))
	Nil(t, err)
	Equal(t, tok, lists.NewOne(
		lists.NewOne(
			primitives.NewConstantInt(1),
			lists.NewOne(
				primitives.NewConstantInt(2),
				primitives.NewConstantInt(3),
			),
		),
		primitives.NewConstantInt(4),
	))

	// simple optional
	tok, err = ParseTavor(strings.NewReader("START = 1 ?(2)\n"))
	Nil(t, err)
	Equal(t, tok, lists.NewAll(
		primitives.NewConstantInt(1),
		constraints.NewOptional(primitives.NewConstantInt(2)),
	))

	// or optional
	tok, err = ParseTavor(strings.NewReader("START = 1 ?(2 | 3) 4\n"))
	Nil(t, err)
	Equal(t, tok, lists.NewAll(
		primitives.NewConstantInt(1),
		constraints.NewOptional(lists.NewOne(
			primitives.NewConstantInt(2),
			primitives.NewConstantInt(3),
		)),
		primitives.NewConstantInt(4),
	))

	// simple repeat
	tok, err = ParseTavor(strings.NewReader("START = 1 +(2)\n"))
	Nil(t, err)
	Equal(t, tok, lists.NewAll(
		primitives.NewConstantInt(1),
		lists.NewRepeat(primitives.NewConstantInt(2), 1, MaxRepeat),
	))

	// or repeat
	tok, err = ParseTavor(strings.NewReader("START = 1 +(2 | 3) 4\n"))
	Nil(t, err)
	Equal(t, tok, lists.NewAll(
		primitives.NewConstantInt(1),
		lists.NewRepeat(lists.NewOne(
			primitives.NewConstantInt(2),
			primitives.NewConstantInt(3),
		), 1, MaxRepeat),
		primitives.NewConstantInt(4),
	))

	// simple optional repeat
	tok, err = ParseTavor(strings.NewReader("START = 1 *(2)\n"))
	Nil(t, err)
	Equal(t, tok, lists.NewAll(
		primitives.NewConstantInt(1),
		lists.NewRepeat(primitives.NewConstantInt(2), 0, MaxRepeat),
	))

	// or optional repeat
	tok, err = ParseTavor(strings.NewReader("START = 1 *(2 | 3) 4\n"))
	Nil(t, err)
	Equal(t, tok, lists.NewAll(
		primitives.NewConstantInt(1),
		lists.NewRepeat(lists.NewOne(
			primitives.NewConstantInt(2),
			primitives.NewConstantInt(3),
		), 0, MaxRepeat),
		primitives.NewConstantInt(4),
	))

	// simple optional repeat
	tok, err = ParseTavor(strings.NewReader("START = 1 *(2)\n"))
	Nil(t, err)
	Equal(t, tok, lists.NewAll(
		primitives.NewConstantInt(1),
		lists.NewRepeat(primitives.NewConstantInt(2), 0, MaxRepeat),
	))

	// exact repeat
	tok, err = ParseTavor(strings.NewReader("START = 1 +3(2)\n"))
	Nil(t, err)
	Equal(t, tok, lists.NewAll(
		primitives.NewConstantInt(1),
		lists.NewRepeat(primitives.NewConstantInt(2), 3, 3),
	))

	// at least repeat
	tok, err = ParseTavor(strings.NewReader("START = 1 +3,(2)\n"))
	Nil(t, err)
	Equal(t, tok, lists.NewAll(
		primitives.NewConstantInt(1),
		lists.NewRepeat(primitives.NewConstantInt(2), 3, MaxRepeat),
	))

	// at most repeat
	tok, err = ParseTavor(strings.NewReader("START = 1 +,3(2)\n"))
	Nil(t, err)
	Equal(t, tok, lists.NewAll(
		primitives.NewConstantInt(1),
		lists.NewRepeat(primitives.NewConstantInt(2), 1, 3),
	))

	// range repeat
	tok, err = ParseTavor(strings.NewReader("START = 1 +2,3(2)\n"))
	Nil(t, err)
	Equal(t, tok, lists.NewAll(
		primitives.NewConstantInt(1),
		lists.NewRepeat(primitives.NewConstantInt(2), 2, 3),
	))
}

func TestTavorParserTokenAttributes(t *testing.T) {
	var tok token.Token
	var err error

	// token attribute List.Count
	tok, err = ParseTavor(strings.NewReader(
		"Digit = 1 | 2 | 3\n" +
			"Digits = *(Digit)\n" +
			"START = Digits \"->\" $Digits.Count\n",
	))
	Nil(t, err)
	{
		v, _ := tok.(*lists.All).Get(0)
		list := v.(*lists.Repeat)

		Equal(t, tok, lists.NewAll(
			lists.NewRepeat(lists.NewOne(
				primitives.NewConstantInt(1),
				primitives.NewConstantInt(2),
				primitives.NewConstantInt(3),
			), 0, MaxRepeat),
			primitives.NewConstantString("->"),
			aggregates.NewLen(list),
		))

		r := test.NewRandTest(1)
		tok.Fuzz(r)
		Equal(t, "12->2", tok.String())
	}
}

func TestTavorParserSpecialTokens(t *testing.T) {
	var tok token.Token
	var err error

	// RandomInt
	tok, err = ParseTavor(strings.NewReader(
		"$Spec = type: Int\nSTART = Spec\n",
	))
	Nil(t, err)
	Equal(t, tok, primitives.NewRandomInt())

	// RangeInt
	tok, err = ParseTavor(strings.NewReader(
		"$Spec = type: Int,\nfrom: 2,\nto: 10\nSTART = Spec\n",
	))
	Nil(t, err)
	Equal(t, tok, primitives.NewRangeInt(2, 10))

	// Sequence
	tok, err = ParseTavor(strings.NewReader(
		"$Spec = type: Sequence\nSTART = $Spec.Next\n",
	))
	Nil(t, err)
	Equal(t, tok, sequences.NewSequence(1, 1).Item())

	tok, err = ParseTavor(strings.NewReader(
		"$Spec = type: Sequence,\nstart: 2\nSTART = $Spec.Next\n",
	))
	Nil(t, err)
	Equal(t, tok, sequences.NewSequence(2, 1).Item())

	tok, err = ParseTavor(strings.NewReader(
		"$Spec = type: Sequence,\nstep: 3\nSTART = $Spec.Next\n",
	))
	Nil(t, err)
	Equal(t, tok, sequences.NewSequence(1, 3).Item())

	tok, err = ParseTavor(strings.NewReader(
		"$Spec = type: Sequence\nSTART = $Spec.Existing\n",
	))
	Nil(t, err)
	Equal(t, tok, sequences.NewSequence(1, 1).ExistingItem())

	tok, err = ParseTavor(strings.NewReader(
		"$Spec = type: Sequence\nSTART = $Spec.Reset\n",
	))
	Nil(t, err)
	Equal(t, tok, sequences.NewSequence(1, 1).ResetItem())
}

func TestTavorParserExpressions(t *testing.T) {
	var tok token.Token
	var err error

	// simple expression
	tok, err = ParseTavor(strings.NewReader(
		"$Spec = type: Sequence\nSTART = ${Spec.Next}\n",
	))
	Nil(t, err)
	Equal(t, tok, sequences.NewSequence(1, 1).Item())

	// plus operator
	tok, err = ParseTavor(strings.NewReader(
		"START = ${1 + 2}\n",
	))
	Nil(t, err)
	Equal(t, tok, expressions.NewAddArithmetic(
		primitives.NewConstantInt(1),
		primitives.NewConstantInt(2),
	))

	// sub operator
	tok, err = ParseTavor(strings.NewReader(
		"START = ${1 - 2}\n",
	))
	Nil(t, err)
	Equal(t, tok, expressions.NewSubArithmetic(
		primitives.NewConstantInt(1),
		primitives.NewConstantInt(2),
	))

	// mul operator
	tok, err = ParseTavor(strings.NewReader(
		"START = ${1 * 2}\n",
	))
	Nil(t, err)
	Equal(t, tok, expressions.NewMulArithmetic(
		primitives.NewConstantInt(1),
		primitives.NewConstantInt(2),
	))

	// div operator
	tok, err = ParseTavor(strings.NewReader(
		"START = ${1 / 2}\n",
	))
	Nil(t, err)
	Equal(t, tok, expressions.NewDivArithmetic(
		primitives.NewConstantInt(1),
		primitives.NewConstantInt(2),
	))

	// nested operator
	tok, err = ParseTavor(strings.NewReader(
		"START = ${1 + 2 + 3}\n",
	))
	Nil(t, err)
	Equal(t, tok, expressions.NewAddArithmetic(
		primitives.NewConstantInt(1),
		expressions.NewAddArithmetic(
			primitives.NewConstantInt(2),
			primitives.NewConstantInt(3),
		),
	))

	// mixed operator
	tok, err = ParseTavor(strings.NewReader(
		"$Spec = type: Sequence\nSTART = ${Spec.Next + 1}\n",
	))
	Nil(t, err)
	Equal(t, tok, expressions.NewAddArithmetic(
		sequences.NewSequence(1, 1).Item(),
		primitives.NewConstantInt(1),
	))
	Equal(t, "2", tok.String())
}

func TestTavorParserAndCuriousCaseOfFuzzing(t *testing.T) {
	var tok token.Token
	var err error

	// loop term (use a term in its definition)
	tok, err = ParseTavor(strings.NewReader(
		"B = 123\nA = B A | 456\nSTART = A\n",
	))
	Nil(t, err)
	{
		a, _ := tok.(*lists.One).Get(0)
		p, _ := a.(*lists.All).Get(1)

		Equal(t, tok, p.(*primitives.Pointer).Get())

		Equal(t, tok, lists.NewOne(
			lists.NewAll(
				primitives.NewConstantInt(123),
				p,
			),
			primitives.NewConstantInt(456),
		))
		r := test.NewRandTest(1)
		tok.Fuzz(r)
		Equal(t, "123456", tok.String())
	}

	// detect endless loops
	tok, err = ParseTavor(strings.NewReader(
		"B = 123\nA = A B\nSTART = A\n",
	))
	Equal(t, ParseErrorEndlessLoopDetected, err.(*ParserError).Type)
	Nil(t, tok)

	tok, err = ParseTavor(strings.NewReader(
		"A = B\nB = A\nSTART = A\n",
	))
	Equal(t, ParseErrorEndlessLoopDetected, err.(*ParserError).Type)
	Nil(t, tok)

	tok, err = ParseTavor(strings.NewReader(
		"B = 123\nA = B A | B\nSTART = A\n",
	))
	Nil(t, err)
	{
		a, _ := tok.(*lists.One).Get(0)
		p, _ := a.(*lists.All).Get(1)

		Equal(t, tok, p.(*primitives.Pointer).Get())

		Equal(t, tok, lists.NewOne(
			lists.NewAll(
				primitives.NewConstantInt(123),
				p,
			),
			primitives.NewConstantInt(123),
		))
	}

	tok, err = ParseTavor(strings.NewReader(
		"B = A\nA = B A | B\nSTART = A\n",
	))
	Equal(t, ParseErrorEndlessLoopDetected, err.(*ParserError).Type)
	Nil(t, tok)

	tok, err = ParseTavor(strings.NewReader(
		"B = A\nA = (A | 1)(B | 2) A\nSTART = A\n",
	))
	Equal(t, ParseErrorEndlessLoopDetected, err.(*ParserError).Type)
	Nil(t, tok)

	tok, err = ParseTavor(strings.NewReader(
		"C = A\nB = C\nA = (B | 1)(B | 2) ?(B)\nSTART = A\n",
	))
	Nil(t, err)
	{
		o, _ := tok.(*lists.All).Get(0)
		p, _ := o.(*lists.One).Get(0)

		Equal(t, tok, p.(*primitives.Pointer).Get())

		Equal(t, tok, lists.NewAll(
			lists.NewOne(
				p,
				primitives.NewConstantInt(1),
			),
			lists.NewOne(
				p,
				primitives.NewConstantInt(2),
			),
			constraints.NewOptional(
				p,
			),
		))
	}

	// Additional forward declaration check
	tok, err = ParseTavor(strings.NewReader(
		"START = Token\nToken = 123\n",
	))
	Nil(t, err)
	Equal(t, tok.(*primitives.Pointer).Get(), primitives.NewConstantInt(123))
}
