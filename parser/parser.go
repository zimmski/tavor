package parser

type ParserErrorType int

const (
	ParseErrorNoStart ParserErrorType = iota
	ParseErrorNewLineNeeded
	ParseErrorEarlyNewLine
	ParseErrorEmptyExpressionIsInvalid
	ParseErrorEmptyTokenDefinition
	ParseErrorInvalidArgumentValue
	ParseErrorInvalidTokenName
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
)

type ParserError struct {
	Message string
	Type    ParserErrorType
}

func (err *ParserError) Error() string {
	return err.Message
}
