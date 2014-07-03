package token

import (
	"fmt"
	"text/scanner"

	"github.com/zimmski/tavor/rand"
)

type Token interface {
	fmt.Stringer

	Clone() Token

	Fuzz(r rand.Rand)
	FuzzAll(r rand.Rand)

	Permutation(i int) error
	Permutations() int
	PermutationsAll() int

	Parse(pars *InternalParser, cur int) (int, []error)
}

type ForwardToken interface {
	Token

	Get() Token

	InternalGet() Token
	InternalLogicalRemove(tok Token) Token
	InternalReplace(oldToken, newToken Token)
}

type OptionalToken interface {
	Token

	IsOptional() bool
	Activate()
	Deactivate()
}

type PermutationErrorType int

const (
	PermutationErrorIndexOutOfBound = iota
)

type PermutationError struct {
	Type PermutationErrorType
}

func (err *PermutationError) Error() string {
	switch err.Type {
	case PermutationErrorIndexOutOfBound:
		return "Permutation index out of bound"
	default:
		return fmt.Sprintf("Unknown permutation error type %#v", err.Type)
	}
}

type ResetToken interface {
	Token

	Reset()
}

type InternalParser struct { // TODO move this some place else
	Data    string
	DataLen int
}

////////////////////////
// was in parser.go but "import cycle not allowed" forced me to do this

type ParserErrorType int

const (
	ParseErrorNoStart ParserErrorType = iota
	ParseErrorNewLineNeeded
	ParseErrorEarlyNewLine
	ParseErrorEmptyExpressionIsInvalid
	ParseErrorEmptyTokenDefinition
	ParseErrorInvalidArgumentValue
	ParseErrorInvalidTokenName
	ParseErrorInvalidTokenType
	ParseErrorUnusedToken
	ParseErrorMissingSpecialTokenArgument
	ParseErrorNonTerminatedString
	ParseErrorTokenAlreadyDefined
	ParseErrorTokenNotDefined
	ParseErrorExpectRune
	ParseErrorUnknownSpecialTokenArgument
	ParseErrorUnknownSpecialTokenType
	ParseErrorUnknownTokenAttribute
	ParseErrorUnknownTypeForSpecialToken
	ParseErrorUnexpectedTokenDefinitionTermination
	ParseErrorExpectedExpressionTerm
	ParseErrorEndlessLoopDetected

	ParseErrorExpectedEOF
	ParseErrorRootIsNil
	ParseErrorUnexpectedEOF
	ParseErrorUnexpectedData
)

type ParserError struct {
	Message string
	Type    ParserErrorType

	Position scanner.Position
}

func (err *ParserError) Error() string {
	return fmt.Sprintf("L:%d, C:%d - %s", err.Position.Line, err.Position.Column, err.Message)
}

////////////////////////
