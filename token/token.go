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

type List interface {
	Token

	Get(i int) (Token, error)
	Len() int

	InternalGet(i int) (Token, error)
	InternalLen() int
	InternalLogicalRemove(tok Token) Token
	InternalReplace(oldToken, newToken Token)
}

type Forward interface {
	Get() Token

	InternalGet() Token
	InternalLogicalRemove(tok Token) Token
	InternalReplace(oldToken, newToken Token)
}

type ForwardToken interface {
	Token
	Forward
}

type Index interface {
	Index() int
}

type IndexToken interface {
	Token
	Index
}

type Optional interface {
	IsOptional() bool
	Activate()
	Deactivate()
}

type OptionalToken interface {
	Token
	Optional
}

type PermutationErrorType int

const (
	PermutationErrorIndexOutOfBound PermutationErrorType = iota
)

type PermutationError struct {
	Type PermutationErrorType
}

func (err *PermutationError) Error() string {
	switch err.Type {
	case PermutationErrorIndexOutOfBound:
		return "permutation index out of bound"
	default:
		return fmt.Sprintf("unknown permutation error type %#v", err.Type)
	}
}

type Reset interface {
	Reset()
}

type ResetToken interface {
	Token
	Reset
}

type ReduceErrorType int

const (
	ReduceErrorIndexOutOfBound ReduceErrorType = iota
)

type ReduceError struct {
	Type ReduceErrorType
}

func (err *ReduceError) Error() string {
	switch err.Type {
	case ReduceErrorIndexOutOfBound:
		return "reduce index out of bound"
	default:
		return fmt.Sprintf("unknown reduce error type %#v", err.Type)
	}
}

type Reduce interface {
	Reduce(i int) error
	Reduces() int
}

type ReduceToken interface {
	Token
	Reduce
}

type Scope interface {
	SetScope(variableScope map[string]Token)
}

type ScopeToken interface {
	Token
	Scope
}

type Variable interface {
	Forward
	Index
	Reset
	Scope

	Name() string
}

type VariableToken interface {
	Token
	Variable
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
	ParseErrorNoTokenForVariable
	ParseErrorTokenAlreadyDefined
	ParseErrorTokenNotDefined
	ParseErrorExpectRune
	ParseErrorExpectOperator
	ParseErrorUnknownBooleanOperator
	ParseErrorUnknownCondition
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
