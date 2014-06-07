package parser

type ParserErrorType int

const (
	ParseErrorNoStart ParserErrorType = iota
	ParseErrorNewLineNeeded
	ParseErrorEarlyNewLine
	ParseErrorEmptyTokenDefinition
	ParseErrorInvalidTokenName
	ParseErrorUnusedToken
	ParseErrorNonTerminatedString
	ParseErrorTokenExists
	ParseErrorTokenDoesNotExists
	ParseErrorExpectRune
	ParseErrorUnexpectedTokenDefinitionTermination
)

type ParserError struct {
	Message string
	Type    ParserErrorType
}

func (err *ParserError) Error() string {
	return err.Message
}
