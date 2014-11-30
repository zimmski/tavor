package token

import (
	"fmt"
	"text/scanner"
)

// Token defines a general token
type Token interface {
	fmt.Stringer

	// Clone returns a copy of the token and all its children
	Clone() Token

	// Permutation sets a specific permutation for this token
	Permutation(i uint) error
	// Permutations returns the number of permutations for this token
	Permutations() uint
	// PermutationsAll returns the number of all possible permutations for this token including its children
	PermutationsAll() uint

	// Parse tries to parse the token beginning from the current position in the parser data.
	// If the parsing is successful the error argument is nil and the next current position after the token is returned.
	Parse(pars *InternalParser, cur int) (int, []error)
}

// List defines a general list token
type List interface {
	// Get returns the current referenced token at the given index. The error return argument is not nil, if the index is out of bound.
	Get(i int) (Token, error)
	// Len returns the number of the current referenced tokens
	Len() int

	// InternalGet returns the current referenced internal token at the given index. The error return argument is not nil, if the index is out of bound.
	InternalGet(i int) (Token, error)
	// InternalLen returns the number of referenced internal tokens
	InternalLen() int
	// InternalLogicalRemove removes the referenced internal token and returns the replacement for the current token or nil if the current token should be removed.
	InternalLogicalRemove(tok Token) Token
	// InternalReplace replaces an old with a new internal token if it is referenced by this token
	InternalReplace(oldToken, newToken Token)
}

// ListToken combines the Token and List interface
type ListToken interface {
	Token
	List
}

// Forward defines a forward token which can reference another token
type Forward interface {
	// Get returns the current referenced token
	Get() Token

	// InternalGet returns the current referenced internal token
	InternalGet() Token
	// InternalLogicalRemove removes the referenced internal token and returns the replacement for the current token or nil if the current token should be removed.
	InternalLogicalRemove(tok Token) Token
	// InternalReplace replaces an old with a new internal token if it is referenced by this token
	InternalReplace(oldToken, newToken Token)
}

// ForwardToken combines the Token and Forward interface
type ForwardToken interface {
	Token
	Forward
}

// Index defines an index token which provides the index in its parent token
type Index interface {
	// Index returns the index of this token in its parent token
	Index() int
}

// IndexToken combines the Token and Index interface
type IndexToken interface {
	Token
	Index
}

// Len defines a general len token
type Len interface {
	// Len returns the number of the current referenced tokens
	Len() int
}

// LenToken combines the Token and Len interface
type LenToken interface {
	Token
	Len
}

// Minimize defines a minimize token which has methods to reduce itself to easier constructs
type Minimize interface {
	// Minimize tries to minimize itself and returns a token if it was successful, or nil if there was nothing to minimize
	Minimize() Token
}

// MinimizeToken combines the Token and Minimize interface
type MinimizeToken interface {
	Token
	Minimize
}

// Optional defines an optional token which can be (de)activated
type Optional interface {
	// IsOptional checks dynamically if this token is in the current state optional
	IsOptional() bool
	// Activate activates this token
	Activate()
	// Deactivate deactivates this token
	Deactivate()
}

// OptionalToken combines the Token and Optional interface
type OptionalToken interface {
	Token
	Optional
}

// Pointer defines a pointer token which can reference another token and can be reset to reference another token
type Pointer interface {
	Forward

	// Set sets the referenced token which must conform to the pointers token reference type
	Set(o Token) error
}

// PointerToken combines the Token and Pointer interface
type PointerToken interface {
	Token
	Pointer
}

// Reset defines a reset token which can reset its (internal) state
type Reset interface {
	// Reset resets the (internal) state of this token and its dependences
	Reset()
}

// ResetToken combines the Token and Index interface
type ResetToken interface {
	Token
	Reset
}

// Reduce defines a reduce token which provides methods to reduce itself and its children
type Reduce interface {
	// Reduce sets a specific reduction for this token
	Reduce(i uint) error
	// Reduces returns the number of reductions for this token
	Reduces() uint
}

// ReduceToken combines the Token and Reduce interface
type ReduceToken interface {
	Token
	Reduce
}

// Release defines a release token which provides methods to release resources on removal
type Release interface {
	// Release gives the token a chance to remove resources
	Release()
}

// ReleaseToken combines the Token and Release interface
type ReleaseToken interface {
	Token
	Release
}

// Scope defines a scope token which holds a scope
type Scope interface {
	// SetScope sets the scope of the token
	SetScope(variableScope map[string]Token)
}

// ScopeToken combines the Token and Scope interface
type ScopeToken interface {
	Token
	Scope
}

// Variable defines a variable token which holds a variable
type Variable interface {
	Forward
	Index
	Scope

	// Name returns the name of the variable
	Name() string

	// Len returns the number of the current referenced tokens
	Len() int
}

// VariableToken combines the Token and Variable interface
type VariableToken interface {
	Token
	Variable
}

////////////////////////

// TODO put this somewhere else?

// PermutationErrorType the permutation error type
type PermutationErrorType int

const (
	// PermutationErrorIndexOutOfBound an index not in the bound of available permutations was used.
	PermutationErrorIndexOutOfBound PermutationErrorType = iota
)

// PermutationError holds a permutation error
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

// ReduceErrorType the reduce error type
type ReduceErrorType int

const (
	// ReduceErrorIndexOutOfBound an index not in the bound of available reductions was used.
	ReduceErrorIndexOutOfBound ReduceErrorType = iota
)

// ReduceError holds a reduce error
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

////////////////////////

// InternalParser holds the data information for an internal parser
type InternalParser struct { // TODO move this some place else
	Data    string
	DataLen int
}

// GetPosition returns a text position in the data given an index of the data
func (p InternalParser) GetPosition(i int) scanner.Position {
	// TODO this could be done MUCH better e.g. memorize or keep count while parsing is happening
	l := 1
	c := 1

	for j := 0; j < i; j++ {
		if p.Data[j] == '\n' {
			l++
			c = 1
		} else {
			c++
		}
	}

	return scanner.Position{
		Line:   l,
		Column: c,
	}
}

////////////////////////
// TODO was in parser.go but "import cycle not allowed" forced me to do this

// ParserErrorType the parser error type
type ParserErrorType int

const (
	// ParseErrorNoStart no Start token was defined
	ParseErrorNoStart ParserErrorType = iota
	// ParseErrorNewLineNeeded a new line is needed
	ParseErrorNewLineNeeded
	// ParseErrorEarlyNewLine an unexpected new line was encountered
	ParseErrorEarlyNewLine
	// ParseErrorEmptyExpressionIsInvalid empty expressions are not allowed
	ParseErrorEmptyExpressionIsInvalid
	// ParseErrorEmptyTokenDefinition empty token definitions are not allowed
	ParseErrorEmptyTokenDefinition
	// ParseErrorInvalidArgumentValue invalid argument value
	ParseErrorInvalidArgumentValue
	// ParseErrorInvalidTokenName invalid token name
	ParseErrorInvalidTokenName
	// ParseErrorInvalidTokenType invalid token type
	ParseErrorInvalidTokenType
	// ParseErrorUnusedToken token is unused
	ParseErrorUnusedToken
	// ParseErrorMissingTypedTokenArgument a typed token argument is missing
	ParseErrorMissingTypedTokenArgument
	// ParseErrorNonTerminatedString string is not properly terminated
	ParseErrorNonTerminatedString
	// ParseErrorNoTokenForVariable variable is not assigned to a token
	ParseErrorNoTokenForVariable
	// ParseErrorTokenAlreadyDefined token name is already in use
	ParseErrorTokenAlreadyDefined
	// ParseErrorTokenNotDefined there is no token with this name
	ParseErrorTokenNotDefined
	// ParseErrorTypeNotDefinedForTypedToken type is not defined for this typed token
	ParseErrorTypeNotDefinedForTypedToken
	// ParseErrorExpectRune the given rune would be expected
	ParseErrorExpectRune
	// ParseErrorExpectOperator the given operator would be expected
	ParseErrorExpectOperator
	// ParseErrorUnknownBooleanOperator the boolean operator is unknown
	ParseErrorUnknownBooleanOperator
	// ParseErrorUnknownCondition the condition is unknown
	ParseErrorUnknownCondition
	// ParseErrorUnknownTypedTokenArgument the typed token argument is unknown
	ParseErrorUnknownTypedTokenArgument
	// ParseErrorUnknownTypedTokenType the typed token type is unknown
	ParseErrorUnknownTypedTokenType
	// ParseErrorUnknownTokenAttribute the token attribute is unknown
	ParseErrorUnknownTokenAttribute
	// ParseErrorUnexpectedTokenDefinitionTermination token definition was unexpectedly terminated
	ParseErrorUnexpectedTokenDefinitionTermination
	// ParseErrorExpectedExpressionTerm expression term is expected
	ParseErrorExpectedExpressionTerm
	// ParseErrorEndlessLoopDetected an invalid loop was detected
	ParseErrorEndlessLoopDetected

	// ParseErrorExpectedEOF expected EOF
	ParseErrorExpectedEOF
	// ParseErrorRootIsNil root token is nil
	ParseErrorRootIsNil
	// ParseErrorUnexpectedEOF EOF was not expected
	ParseErrorUnexpectedEOF
	// ParseErrorUnexpectedData additional data was not expected
	ParseErrorUnexpectedData
)

// ParserError holds a parser error
type ParserError struct {
	Message string
	Type    ParserErrorType

	Position scanner.Position
}

func (err *ParserError) Error() string {
	return fmt.Sprintf("L:%d, C:%d - %s", err.Position.Line, err.Position.Column, err.Message)
}

////////////////////////
