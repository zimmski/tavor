package parser

import (
	"fmt"
	"strings"
	"testing"

	. "github.com/stretchr/testify/assert"

	"github.com/zimmski/tavor"
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

	// special token already defined
	tok, err = ParseTavor(strings.NewReader("START = 123\n$START = 456\n"))
	Equal(t, token.ParseErrorTokenAlreadyDefined, err.(*token.ParserError).Type)
	Nil(t, tok)

	// expect = in special token
	tok, err = ParseTavor(strings.NewReader("$START 123\n"))
	Equal(t, token.ParseErrorExpectRune, err.(*token.ParserError).Type)
	Nil(t, tok)

	// expect identifier in special token
	tok, err = ParseTavor(strings.NewReader("$START = 123\n"))
	Equal(t, token.ParseErrorExpectRune, err.(*token.ParserError).Type)
	Nil(t, tok)

	// expect : in special token
	tok, err = ParseTavor(strings.NewReader("$START = argument\n"))
	Equal(t, token.ParseErrorExpectRune, err.(*token.ParserError).Type)
	Nil(t, tok)

	// expect valid argument value in special token
	tok, err = ParseTavor(strings.NewReader("$START = argument:\n"))
	Equal(t, token.ParseErrorInvalidArgumentValue, err.(*token.ParserError).Type)
	Nil(t, tok)

	// expect no eof after argument value in special token
	tok, err = ParseTavor(strings.NewReader("$START = argument:value"))
	Equal(t, token.ParseErrorNewLineNeeded, err.(*token.ParserError).Type)
	Nil(t, tok)

	// expect new line in special token
	tok, err = ParseTavor(strings.NewReader("$START = argument:value$"))
	Equal(t, token.ParseErrorExpectRune, err.(*token.ParserError).Type)
	Nil(t, tok)

	// undefined type argument special token
	tok, err = ParseTavor(strings.NewReader("$START = hey: ok\n"))
	Equal(t, token.ParseErrorUnknownTypeForSpecialToken, err.(*token.ParserError).Type)
	Nil(t, tok)

	// unknown type argument special token
	tok, err = ParseTavor(strings.NewReader("$START = type: ok\n"))
	Equal(t, token.ParseErrorUnknownSpecialTokenType, err.(*token.ParserError).Type)
	Nil(t, tok)

	// unknown special token argument
	tok, err = ParseTavor(strings.NewReader("$START = type: Int,\nok: value\n"))
	Equal(t, token.ParseErrorUnknownSpecialTokenArgument, err.(*token.ParserError).Type)
	Nil(t, tok)

	// missing arguments for special token Int
	tok, err = ParseTavor(strings.NewReader("$START = type: Int,\nto:123\n"))
	Equal(t, token.ParseErrorMissingSpecialTokenArgument, err.(*token.ParserError).Type)
	Nil(t, tok)

	tok, err = ParseTavor(strings.NewReader("$START = type: Int,\nfrom:123\n"))
	Equal(t, token.ParseErrorMissingSpecialTokenArgument, err.(*token.ParserError).Type)
	Nil(t, tok)

	// invalid arguments for special token Int
	tok, err = ParseTavor(strings.NewReader("$START = type: Int,\nfrom:abc,\nto:123\n"))
	Equal(t, token.ParseErrorInvalidArgumentValue, err.(*token.ParserError).Type)
	Nil(t, tok)

	tok, err = ParseTavor(strings.NewReader("$START = type: Int,\nfrom:123,\nto:abc\n"))
	Equal(t, token.ParseErrorInvalidArgumentValue, err.(*token.ParserError).Type)
	Nil(t, tok)

	// invalid arguments for special token Sequence
	tok, err = ParseTavor(strings.NewReader("$START = type: Sequence,\nstart:abc\n"))
	Equal(t, token.ParseErrorInvalidArgumentValue, err.(*token.ParserError).Type)
	Nil(t, tok)

	tok, err = ParseTavor(strings.NewReader("$START = type: Sequence,\nstep:abc\n"))
	Equal(t, token.ParseErrorInvalidArgumentValue, err.(*token.ParserError).Type)
	Nil(t, tok)

	// empty expression
	tok, err = ParseTavor(strings.NewReader("START = ${}\n"))
	Equal(t, token.ParseErrorEmptyExpressionIsInvalid, err.(*token.ParserError).Type)
	Nil(t, tok)

	// open expression
	tok, err = ParseTavor(strings.NewReader("$Spec = type: Sequence\nSTART = ${Spec.Next\n"))
	Equal(t, token.ParseErrorExpectRune, err.(*token.ParserError).Type)
	Nil(t, tok)

	// missing operator expression term
	tok, err = ParseTavor(strings.NewReader("$Spec = type: Sequence\nSTART = ${Spec.Next +}\n"))
	Equal(t, token.ParseErrorExpectedExpressionTerm, err.(*token.ParserError).Type)
	Nil(t, tok)

	// TODO this can maybe never happen as we do everything in one pass
	// so we do not know that $List must implement lists.List

	// // wrong token type because of earlier usage
	// tok, err = ParseTavor(strings.NewReader("START = $List.Count\nList = 123"))
	// panic(err)
	// Equal(t, token.ParseErrorExpectedExpressionTerm, err.(*token.ParserError).Type)
	// Nil(t, tok)
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
		lists.NewRepeat(primitives.NewConstantInt(2), 1, tavor.MaxRepeat),
	))

	// or repeat
	tok, err = ParseTavor(strings.NewReader("START = 1 +(2 | 3) 4\n"))
	Nil(t, err)
	Equal(t, tok, lists.NewAll(
		primitives.NewConstantInt(1),
		lists.NewRepeat(lists.NewOne(
			primitives.NewConstantInt(2),
			primitives.NewConstantInt(3),
		), 1, tavor.MaxRepeat),
		primitives.NewConstantInt(4),
	))

	// simple optional repeat
	tok, err = ParseTavor(strings.NewReader("START = 1 *(2)\n"))
	Nil(t, err)
	Equal(t, tok, lists.NewAll(
		primitives.NewConstantInt(1),
		lists.NewRepeat(primitives.NewConstantInt(2), 0, tavor.MaxRepeat),
	))

	// or optional repeat
	tok, err = ParseTavor(strings.NewReader("START = 1 *(2 | 3) 4\n"))
	Nil(t, err)
	Equal(t, tok, lists.NewAll(
		primitives.NewConstantInt(1),
		lists.NewRepeat(lists.NewOne(
			primitives.NewConstantInt(2),
			primitives.NewConstantInt(3),
		), 0, tavor.MaxRepeat),
		primitives.NewConstantInt(4),
	))

	// simple optional repeat
	tok, err = ParseTavor(strings.NewReader("START = 1 *(2)\n"))
	Nil(t, err)
	Equal(t, tok, lists.NewAll(
		primitives.NewConstantInt(1),
		lists.NewRepeat(primitives.NewConstantInt(2), 0, tavor.MaxRepeat),
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
		lists.NewRepeat(primitives.NewConstantInt(2), 3, tavor.MaxRepeat),
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
			), 0, tavor.MaxRepeat),
			primitives.NewConstantString("->"),
			aggregates.NewLen(list),
		))

		r := test.NewRandTest(1)
		tok.FuzzAll(r)
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

	// detect endless loops
	tok, err = ParseTavor(strings.NewReader(
		"B = 123\nA = A B\nSTART = A\n",
	))
	Equal(t, token.ParseErrorEndlessLoopDetected, err.(*token.ParserError).Type)
	Nil(t, tok)

	tok, err = ParseTavor(strings.NewReader(
		"A = B\nB = A\nSTART = A\n",
	))
	Equal(t, token.ParseErrorEndlessLoopDetected, err.(*token.ParserError).Type)
	Nil(t, tok)

	tok, err = ParseTavor(strings.NewReader(
		"B = A\nA = B A | B\nSTART = A\n",
	))
	Equal(t, token.ParseErrorEndlessLoopDetected, err.(*token.ParserError).Type)
	Nil(t, tok)

	tok, err = ParseTavor(strings.NewReader(
		"B = A\nA = (A | 1)(B | 2) A\nSTART = A\n",
	))
	Equal(t, token.ParseErrorEndlessLoopDetected, err.(*token.ParserError).Type)
	Nil(t, tok)

	// Additional forward declaration check
	tok, err = ParseTavor(strings.NewReader(
		"START = Token\nToken = 123\n",
	))
	Nil(t, err)
	Equal(t, tok, primitives.NewConstantInt(123))

	// double embedded forward token all the way
	tok, err = ParseTavor(strings.NewReader("A = B B\nB = 1\nSTART = A\n"))
	Nil(t, err)
	Equal(t, tok, lists.NewAll(
		primitives.NewConstantInt(1),
		primitives.NewConstantInt(1),
	))

	// Tokens should be cloned so they are different internally
	{
		tok, err = ParseTavor(strings.NewReader(
			"Token = 1 | 2\nSTART = Token Token\n",
		))
		Nil(t, err)
		Equal(t, tok, lists.NewAll(
			lists.NewOne(primitives.NewConstantInt(1), primitives.NewConstantInt(2)),
			lists.NewOne(primitives.NewConstantInt(1), primitives.NewConstantInt(2)),
		))

		va, _ := tok.(lists.List).Get(0)
		a := va.(*lists.One)
		vb, _ := tok.(lists.List).Get(1)
		b := vb.(*lists.One)

		True(t, Exactly(t, a, b))
		NotEqual(t, fmt.Sprintf("%p", a), fmt.Sprintf("%p", b))
	}

	// Correct sequence behaviour
	{
		tok, err = ParseTavor(strings.NewReader(`
			$Id = type: Sequence,
				start: 2,
				step: 2

			START = $Id.Reset $Id.Next $Id.Existing
		`))
		Nil(t, err)

		r := test.NewRandTest(1)

		tok.FuzzAll(r)

		Equal(t, "22", tok.String())
	}

	// Correct list behaviour
	{
		tok, err = ParseTavor(strings.NewReader(`
			A = *2(1)

			START = $A.Count A
		`))
		Nil(t, err)

		r := test.NewRandTest(1)

		tok.FuzzAll(r)

		Equal(t, "211", tok.String())
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
		Equal(t, tok, lists.NewOne(
			lists.NewOne(
				lists.NewOne(
					primitives.NewConstantInt(1),
				),
				primitives.NewConstantInt(1),
			),
			primitives.NewConstantInt(1),
		))

		Equal(t, "1", tok.String())
	}

	tok, err = ParseTavor(strings.NewReader(`
		A = A 1 | 2

		START = A
	`))
	Nil(t, err)
	{
		Equal(t, tok, lists.NewOne(
			lists.NewAll(
				lists.NewOne(
					lists.NewAll(
						lists.NewOne(
							primitives.NewConstantInt(2),
						),
						primitives.NewConstantInt(1),
					),
					primitives.NewConstantInt(2),
				),
				primitives.NewConstantInt(1),
			),
			primitives.NewConstantInt(2),
		))

		Equal(t, "211", tok.String())
	}

	// optional loop
	tok, err = ParseTavor(strings.NewReader(`
		A = ?(A) 1

		START = A
	`))
	Nil(t, err)
	{
		Equal(t, tok, lists.NewAll(
			constraints.NewOptional(
				lists.NewAll(
					constraints.NewOptional(
						lists.NewAll(
							primitives.NewConstantInt(1),
						),
					),
					primitives.NewConstantInt(1),
				),
			),
			primitives.NewConstantInt(1),
		))

		Equal(t, "111", tok.String())
	}

	// One loop
	tok, err = ParseTavor(strings.NewReader(`
		A = (A | 2) 1

		START = A
	`))
	Nil(t, err)
	{
		Equal(t, tok, lists.NewAll(
			lists.NewOne(
				lists.NewAll(
					lists.NewOne(
						lists.NewAll(
							lists.NewOne(
								primitives.NewConstantInt(2),
							),
							primitives.NewConstantInt(1),
						),
						primitives.NewConstantInt(2),
					),
					primitives.NewConstantInt(1),
				),
				primitives.NewConstantInt(2),
			),
			primitives.NewConstantInt(1),
		))

		Equal(t, "2111", tok.String())
	}

	// two loops
	/*
		TODO fix this. the problem is that the max loop is reached with 2 as max too soon
		since A is used twice in the whole token. so it is cloned and then gone through
		and the first entcounter exceeds the max repeats
		THE FIXME should be to remember which loop has which parent. so we now that
		we can repeat the something with the same parent but ignore something with a different
		parent


		log.LevelDebug()
		tok, err = ParseTavor(strings.NewReader(`
			A = (A | 2) 1 ?(A) 3

			START = A
		`))
		Nil(t, err)
		{
			tavor.PrettyPrintInternalTree(os.Stdout, tok)

			Equal(t, tok, lists.NewAll(
				lists.NewOne(
					lists.NewAll(
						lists.NewOne(
							lists.NewAll(
								primitives.NewConstantInt(2),
								primitives.NewConstantInt(1),
								primitives.NewConstantInt(3),
							),
							primitives.NewConstantInt(2),
						),
						primitives.NewConstantInt(1),
						constraints.NewOptional(
							lists.NewAll(
								primitives.NewConstantInt(2),
								primitives.NewConstantInt(1),
								primitives.NewConstantInt(3),
							),
						),
						primitives.NewConstantInt(3),
					),
					primitives.NewConstantInt(2),
				),
				primitives.NewConstantInt(1),
				constraints.NewOptional(
					lists.NewAll(
						lists.NewOne(
							lists.NewAll(
								primitives.NewConstantInt(2),
								primitives.NewConstantInt(1),
								primitives.NewConstantInt(3),
							),
							primitives.NewConstantInt(2),
						),
						primitives.NewConstantInt(1),
						constraints.NewOptional(
							lists.NewAll(
								primitives.NewConstantInt(2),
								primitives.NewConstantInt(1),
								primitives.NewConstantInt(3),
							),
						),
						primitives.NewConstantInt(3),
					),
				),
				primitives.NewConstantInt(3),
			))

			Equal(t, "213121331213121333", tok.String())
		}
		log.LevelWarn()
	*/

	// loop with token loop
	tok, err = ParseTavor(strings.NewReader(`
			B = A
			C = B

			A = ?(C) 1

			START = A
		`))
	Nil(t, err)
	{
		Equal(t, tok, lists.NewAll(
			constraints.NewOptional(
				lists.NewAll(
					constraints.NewOptional(
						lists.NewAll(
							primitives.NewConstantInt(1),
						),
					),
					primitives.NewConstantInt(1),
				),
			),
			primitives.NewConstantInt(1),
		))

		Equal(t, "111", tok.String())
	}
}
