package parser

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	. "github.com/zimmski/tavor/test/assert"

	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/constraints"
	"github.com/zimmski/tavor/token/lists"
	"github.com/zimmski/tavor/token/primitives"
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

	// range integer errors
	errs = ParseInternal(primitives.NewRangeInt(1, 10), strings.NewReader(""))
	Equal(t, len(errs), 1)
	Equal(t, token.ParseErrorUnexpectedEOF, errs[0].(*token.ParserError).Type)

	errs = ParseInternal(primitives.NewRangeInt(1, 10), strings.NewReader("0"))
	Equal(t, len(errs), 1)
	Equal(t, token.ParseErrorUnexpectedData, errs[0].(*token.ParserError).Type)

	errs = ParseInternal(primitives.NewRangeInt(1, 10), strings.NewReader("11"))
	Equal(t, len(errs), 1)
	Equal(t, token.ParseErrorExpectedEOF, errs[0].(*token.ParserError).Type)

	errs = ParseInternal(primitives.NewRangeIntWithStep(2, 10, 2), strings.NewReader("3"))
	Equal(t, len(errs), 1)
	Equal(t, token.ParseErrorUnexpectedData, errs[0].(*token.ParserError).Type)

	// constant string errors
	errs = ParseInternal(primitives.NewConstantString("a"), strings.NewReader(""))
	Equal(t, len(errs), 1)
	Equal(t, token.ParseErrorUnexpectedEOF, errs[0].(*token.ParserError).Type)

	errs = ParseInternal(primitives.NewConstantString("a"), strings.NewReader("b"))
	Equal(t, len(errs), 1)
	Equal(t, token.ParseErrorUnexpectedData, errs[0].(*token.ParserError).Type)

	// too much data left
	errs = ParseInternal(primitives.NewConstantString("123456"), strings.NewReader("123456abcdefgh"))
	Equal(t, len(errs), 1)
	Equal(t, token.ParseErrorExpectedEOF, errs[0].(*token.ParserError).Type)
}

func checkParse(t *testing.T, root token.Token, data string) {
	errs := ParseInternal(root, strings.NewReader(data))

	if len(errs) > 0 {
		fmt.Printf("ERRS: %+v\n", errs)

		panic(fmt.Sprintf("Expected nil, but got: %#v", errs))
	}
	if got := root.String(); !ObjectsAreEqual(data, got) {
		panic(fmt.Sprintf("Not equal: %#v != %#v", data, got))
	}
}

func TestInternalParse(t *testing.T) {
	var o token.Token
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

	// range integer
	checkParse(
		t,
		primitives.NewRangeInt(1, 1),
		"1",
	)
	for i := 1; i <= 10; i++ {
		checkParse(
			t,
			primitives.NewRangeInt(1, 10),
			strconv.Itoa(i),
		)
	}

	checkParse(
		t,
		primitives.NewRangeIntWithStep(2, 2, 2),
		"2",
	)
	for i := 2; i <= 10; i += 2 {
		checkParse(
			t,
			primitives.NewRangeIntWithStep(2, 10, 2),
			strconv.Itoa(i),
		)
	}

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

	// character class
	o = constraints.NewOptional(
		primitives.NewCharacterClass(`a`),
	)
	checkParse(
		t,
		o,
		"a",
	)

	errs = ParseInternal(o, strings.NewReader(""))
	Nil(t, errs)

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

	errs = ParseInternal(lists.NewAll(
		primitives.NewConstantInt(1),
		primitives.NewConstantString("a"),
	), strings.NewReader("1a2b"))

	Equal(t, token.ParseErrorExpectedEOF, errs[0].(*token.ParserError).Type)

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

	// optional
	o = lists.NewAll(
		constraints.NewOptional(primitives.NewConstantInt(1)),
		primitives.NewConstantString("a"),
	)

	checkParse(
		t,
		o,
		"1a",
	)

	checkParse(
		t,
		o,
		"a",
	)

	errs = ParseInternal(o, strings.NewReader("1c"))
	Equal(t, token.ParseErrorUnexpectedData, errs[0].(*token.ParserError).Type)

	errs = ParseInternal(o, strings.NewReader("21a"))
	Equal(t, token.ParseErrorUnexpectedData, errs[0].(*token.ParserError).Type)

	// repeat
	o = lists.NewAll(
		primitives.NewConstantInt(1),
		lists.NewRepeat(primitives.NewConstantInt(2), 2, 5),
	)

	checkParse(
		t,
		o,
		"122",
	)

	checkParse(
		t,
		o,
		"122222",
	)

	errs = ParseInternal(o, strings.NewReader("12"))
	Equal(t, token.ParseErrorUnexpectedEOF, errs[0].(*token.ParserError).Type)

	errs = ParseInternal(o, strings.NewReader("1222222"))
	Equal(t, token.ParseErrorExpectedEOF, errs[0].(*token.ParserError).Type)

	o = lists.NewAll(
		primitives.NewConstantInt(1),
		lists.NewRepeat(primitives.NewConstantInt(2), 0, 5),
	)

	checkParse(
		t,
		o,
		"1",
	)

	checkParse(
		t,
		o,
		"12",
	)

	checkParse(
		t,
		o,
		"122222",
	)

	// complex repeat
	o = lists.NewAll(
		primitives.NewConstantInt(1),
		lists.NewRepeat(lists.NewOne(
			primitives.NewConstantInt(2),
			primitives.NewConstantInt(3),
		), 1, 30),
		primitives.NewConstantInt(4),
	)

	checkParse(
		t,
		o,
		"13232323232323333323232224",
	)
}
