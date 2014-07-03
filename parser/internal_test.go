package parser

import (
	"fmt"
	"github.com/zimmski/tavor/token/lists"
	"github.com/zimmski/tavor/token/primitives"
	"strings"
	"testing"

	. "github.com/stretchr/testify/assert"

	"github.com/zimmski/tavor/token"
)

func TestInternalParseErrors(t *testing.T) {
	var errs []error

	// nil root token
	errs = ParseInternal(nil, strings.NewReader(""))
	Equal(t, len(errs), 1)
	Equal(t, token.ParseErrorRootIsNil, errs[0].(*token.ParserError).Type)

	// constant integer errors
	errs = ParseInternal(primitives.NewConstantInt(1), strings.NewReader(""))
	Equal(t, len(errs), 1)
	Equal(t, token.ParseErrorUnexpectedEOF, errs[0].(*token.ParserError).Type)

	errs = ParseInternal(primitives.NewConstantInt(1), strings.NewReader("2"))
	Equal(t, len(errs), 1)
	Equal(t, token.ParseErrorUnexpectedData, errs[0].(*token.ParserError).Type)

	errs = ParseInternal(primitives.NewConstantInt(1), strings.NewReader("123"))
	Equal(t, len(errs), 1)
	Equal(t, token.ParseErrorExpectedEOF, errs[0].(*token.ParserError).Type)

	errs = ParseInternal(primitives.NewConstantInt(1), strings.NewReader("1234567890"))
	Equal(t, len(errs), 1)
	Equal(t, token.ParseErrorExpectedEOF, errs[0].(*token.ParserError).Type)

	// constant string errors
	errs = ParseInternal(primitives.NewConstantString("a"), strings.NewReader(""))
	Equal(t, len(errs), 1)
	Equal(t, token.ParseErrorUnexpectedEOF, errs[0].(*token.ParserError).Type)

	errs = ParseInternal(primitives.NewConstantString("a"), strings.NewReader("b"))
	Equal(t, len(errs), 1)
	Equal(t, token.ParseErrorUnexpectedData, errs[0].(*token.ParserError).Type)
}

func checkParse(t *testing.T, root token.Token, data string) {
	errs := ParseInternal(root, strings.NewReader(data))
	if len(errs) != 0 {
		fmt.Printf("ERRS: %+v\n", errs)
	}
	Nil(t, errs)
	Equal(t, data, root.String())
}

func TestInternalParse(t *testing.T) {
	var tok, o token.Token
	var errs []error

	// constant integer
	checkParse(
		t,
		primitives.NewConstantInt(1),
		"1",
	)

	checkParse(
		t,
		primitives.NewConstantInt(123),
		"123",
	)

	// constant string
	checkParse(
		t,
		primitives.NewConstantString("a"),
		"a",
	)

	checkParse(
		t,
		primitives.NewConstantString("abc"),
		"abc",
	)

	// All
	checkParse(
		t,
		lists.NewAll(
			primitives.NewConstantInt(1),
			primitives.NewConstantString("a"),
		),
		"1a",
	)

	errs = ParseInternal(lists.NewAll(
		primitives.NewConstantInt(1),
		primitives.NewConstantString("a"),
		primitives.NewConstantInt(2),
	), strings.NewReader("1a"))

	Equal(t, token.ParseErrorUnexpectedEOF, errs[0].(*token.ParserError).Type)
	Nil(t, tok)

	errs = ParseInternal(lists.NewAll(
		primitives.NewConstantInt(1),
		primitives.NewConstantString("a"),
	), strings.NewReader("1a2b"))

	Equal(t, token.ParseErrorExpectedEOF, errs[0].(*token.ParserError).Type)
	Nil(t, tok)

	// One
	o = lists.NewOne(
		primitives.NewConstantInt(1),
		primitives.NewConstantString("a"),
	)

	checkParse(
		t,
		o,
		"1",
	)

	checkParse(
		t,
		o,
		"a",
	)

	errs = ParseInternal(o, strings.NewReader("2"))
	Equal(t, token.ParseErrorUnexpectedData, errs[0].(*token.ParserError).Type)
	Nil(t, tok)

	// combine
	o = lists.NewOne(
		lists.NewAll(
			primitives.NewConstantInt(1),
			primitives.NewConstantString("a"),
		),
		lists.NewAll(
			primitives.NewConstantInt(1),
			primitives.NewConstantString("b"),
		),
	)

	checkParse(
		t,
		o,
		"1a",
	)

	checkParse(
		t,
		o,
		"1b",
	)

	errs = ParseInternal(o, strings.NewReader("1c"))
	Equal(t, token.ParseErrorUnexpectedData, errs[0].(*token.ParserError).Type)
	Nil(t, tok)
}
