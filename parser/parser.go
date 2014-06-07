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
)

type ParserError struct {
	Message string
	Type    ParserErrorType
}

func (err *ParserError) Error() string {
	return err.Message
}
