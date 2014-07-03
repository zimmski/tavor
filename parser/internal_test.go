package parser

import (
	"github.com/zimmski/tavor/token/primitives"
	"strings"
	"testing"

	. "github.com/stretchr/testify/assert"

	"github.com/zimmski/tavor/token"
)

func TestInternalParseErrors(t *testing.T) {
	var tok token.Token
	var err error

	// nil root token
	tok, err = ParseInternal(nil, strings.NewReader(""))
	Equal(t, token.ParseErrorRootIsNil, err.(*token.ParserError).Type)
	Nil(t, tok)

	// constant integer errors
	tok, err = ParseInternal(primitives.NewConstantInt(1), strings.NewReader(""))
	Equal(t, token.ParseErrorUnexpectedEOF, err.(*token.ParserError).Type)
	Nil(t, tok)

	tok, err = ParseInternal(primitives.NewConstantInt(1), strings.NewReader("2"))
	Equal(t, token.ParseErrorUnexpectedData, err.(*token.ParserError).Type)
	Nil(t, tok)
}

func TestInternalParse(t *testing.T) {
	var tok token.Token
	var err error

	// constant integer
	tok, err = ParseInternal(primitives.NewConstantInt(1), strings.NewReader("1"))
	Nil(t, err)
	Equal(t, tok, primitives.NewConstantInt(1))

	tok, err = ParseInternal(primitives.NewConstantInt(123), strings.NewReader("123"))
	Nil(t, err)
	Equal(t, tok, primitives.NewConstantInt(123))
}
