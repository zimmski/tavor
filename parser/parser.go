package parser

import (
	"fmt"
	"text/scanner"
)

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
	ParseErrorRootIsNil
)

type ParserError struct {
	Message string
	Type    ParserErrorType

	Position scanner.Position
}

func (err *ParserError) Error() string {
	return fmt.Sprintf("L:%d, C:%d - %s", err.Position.Line, err.Position.Column, err.Message)
}
