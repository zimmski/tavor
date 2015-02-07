package conditions

import (
	"fmt"

	"github.com/zimmski/tavor/log"
	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/lists"
	"github.com/zimmski/tavor/token/primitives"
)

// BooleanExpression defines a boolean expression
type BooleanExpression interface {
	token.Token

	// Evaluate evaluates the boolean expression and returns its result
	Evaluate() bool
}

// BooleanTrue implements a boolean expression which evaluates to always true
type BooleanTrue struct{}

// NewBooleanTrue returns a new instance of a BooleanTrue token
func NewBooleanTrue() *BooleanTrue {
	return &BooleanTrue{}
}

// Evaluate evaluates the boolean expression and returns its result
func (c *BooleanTrue) Evaluate() bool {
	return true
}

// Token interface methods

// Clone returns a copy of the token and all its children
func (c *BooleanTrue) Clone() token.Token {
	return &BooleanTrue{}
}

// Parse tries to parse the token beginning from the current position in the parser data.
// If the parsing is successful the error argument is nil and the next current position after the token is returned.
func (c *BooleanTrue) Parse(pars *token.InternalParser, cur int) (int, []error) {
	panic("This should never happen")
}

// Permutation sets a specific permutation for this token
func (c *BooleanTrue) Permutation(i uint) error {
	// do nothing

	return nil
}

// Permutations returns the number of permutations for this token
func (c *BooleanTrue) Permutations() uint {
	return 1
}

// PermutationsAll returns the number of all possible permutations for this token including its children
func (c *BooleanTrue) PermutationsAll() uint {
	return 1
}

func (c *BooleanTrue) String() string {
	return "true"
}

// List interface methods

// Get returns the current referenced token at the given index. The error return argument is not nil, if the index is out of bound.
func (c *BooleanTrue) Get(i int) (token.Token, error) {
	return nil, &lists.ListError{
		Type: lists.ListErrorOutOfBound,
	}
}

// Len returns the number of the current referenced tokens
func (c *BooleanTrue) Len() int {
	return 0
}

// InternalGet returns the current referenced internal token at the given index. The error return argument is not nil, if the index is out of bound.
func (c *BooleanTrue) InternalGet(i int) (token.Token, error) {
	return nil, &lists.ListError{
		Type: lists.ListErrorOutOfBound,
	}
}

// InternalLen returns the number of referenced internal tokens
func (c *BooleanTrue) InternalLen() int {
	return 0
}

// InternalLogicalRemove removes the referenced internal token and returns the replacement for the current token or nil if the current token should be removed.
func (c *BooleanTrue) InternalLogicalRemove(tok token.Token) token.Token {
	panic("This should never happen")
}

// InternalReplace replaces an old with a new internal token if it is referenced by this token. The error return argument is not nil, if the replacement is not suitable.
func (c *BooleanTrue) InternalReplace(oldToken, newToken token.Token) error {
	panic("This should never happen")
}

// BooleanEqual implements a boolean expression which compares the value of two tokens
type BooleanEqual struct {
	a, b token.Token
}

// NewBooleanEqual returns a new instance of a BooleanEqual token referencing two tokens
func NewBooleanEqual(a, b token.Token) *BooleanEqual {
	return &BooleanEqual{
		a: a,
		b: b,
	}
}

// Evaluate evaluates the boolean expression and returns its result
func (c *BooleanEqual) Evaluate() bool {
	return c.a.String() == c.b.String()
}

// Token interface methods

// Clone returns a copy of the token and all its children
func (c *BooleanEqual) Clone() token.Token {
	return &BooleanEqual{
		a: c.a,
		b: c.b,
	}
}

// Parse tries to parse the token beginning from the current position in the parser data.
// If the parsing is successful the error argument is nil and the next current position after the token is returned.
func (c *BooleanEqual) Parse(pars *token.InternalParser, cur int) (int, []error) {
	panic("This should never happen")
}

// Permutation sets a specific permutation for this token
func (c *BooleanEqual) Permutation(i uint) error {
	// do nothing

	return nil
}

// Permutations returns the number of permutations for this token
func (c *BooleanEqual) Permutations() uint {
	return 1
}

// PermutationsAll returns the number of all possible permutations for this token including its children
func (c *BooleanEqual) PermutationsAll() uint {
	return 1
}

func (c *BooleanEqual) String() string {
	return fmt.Sprintf("(%p)%#v == (%p)%#v", c.a, c.a, c.b, c.b)
}

// List interface methods

// Get returns the current referenced token at the given index. The error return argument is not nil, if the index is out of bound.
func (c *BooleanEqual) Get(i int) (token.Token, error) {
	return nil, &lists.ListError{
		Type: lists.ListErrorOutOfBound,
	}
}

// Len returns the number of the current referenced tokens
func (c *BooleanEqual) Len() int {
	return 0
}

// InternalGet returns the current referenced internal token at the given index. The error return argument is not nil, if the index is out of bound.
func (c *BooleanEqual) InternalGet(i int) (token.Token, error) {
	switch i {
	case 0:
		return c.a, nil
	case 1:
		return c.b, nil
	default:
		return nil, &lists.ListError{
			Type: lists.ListErrorOutOfBound,
		}
	}
}

// InternalLen returns the number of referenced internal tokens
func (c *BooleanEqual) InternalLen() int {
	return 2
}

// InternalLogicalRemove removes the referenced internal token and returns the replacement for the current token or nil if the current token should be removed.
func (c *BooleanEqual) InternalLogicalRemove(tok token.Token) token.Token {
	if tok == c.a || tok == c.b {
		return nil
	}

	return c
}

// InternalReplace replaces an old with a new internal token if it is referenced by this token. The error return argument is not nil, if the replacement is not suitable.
func (c *BooleanEqual) InternalReplace(oldToken, newToken token.Token) error {
	if oldToken == c.a {
		c.a = newToken
	}
	if oldToken == c.b {
		c.b = newToken
	}

	return nil
}

// VariableDefined implements a boolean expression which evaluates if a variable is defined in a given scope
type VariableDefined struct {
	name          string
	variableScope *token.VariableScope
}

// NewVariableDefined returns a new instance of a VariableDefined token initialzed with the given name and scope
func NewVariableDefined(name string, variableScope *token.VariableScope) *VariableDefined {
	return &VariableDefined{
		name:          name,
		variableScope: variableScope,
	}
}

// Evaluate evaluates the boolean expression and returns its result
func (c *VariableDefined) Evaluate() bool {
	return c.variableScope.Get(c.name) != nil
}

// Token interface methods

// Clone returns a copy of the token and all its children
func (c *VariableDefined) Clone() token.Token {
	return &VariableDefined{
		name:          c.name,
		variableScope: c.variableScope,
	}
}

// Parse tries to parse the token beginning from the current position in the parser data.
// If the parsing is successful the error argument is nil and the next current position after the token is returned.
func (c *VariableDefined) Parse(pars *token.InternalParser, cur int) (int, []error) {
	panic("This should never happen")
}

// Permutation sets a specific permutation for this token
func (c *VariableDefined) Permutation(i uint) error {
	// do nothing

	return nil
}

// Permutations returns the number of permutations for this token
func (c *VariableDefined) Permutations() uint {
	return 1
}

// PermutationsAll returns the number of all possible permutations for this token including its children
func (c *VariableDefined) PermutationsAll() uint {
	return 1
}

func (c *VariableDefined) String() string {
	return fmt.Sprintf("defined(%q)", c.name)
}

// ScopeToken interface methods

// SetScope sets the scope of the token
func (c *VariableDefined) SetScope(variableScope *token.VariableScope) {
	c.variableScope = variableScope
}

// ExpressionPointer implements a token pointer to an expression token
type ExpressionPointer struct {
	token token.Token
}

// NewExpressionPointer returns a new instance of a ExpressionPointer token referencing the given token
func NewExpressionPointer(token token.Token) *ExpressionPointer {
	return &ExpressionPointer{
		token: token,
	}
}

// Evaluate evaluates the boolean expression and returns its result
func (c *ExpressionPointer) Evaluate() bool {
	tok := c.token

	if po, ok := tok.(*primitives.Pointer); ok {
		log.Debugf("Found pointer in ExpressionPointer %p(%#v)", c, c)

		for {
			c := po.InternalGet()
			c = c.Clone()
			_ = po.Set(c)

			po, ok = c.(*primitives.Pointer)
			if !ok {
				log.Debugf("Replaced pointer %p(%#v) with %p(%#v)", tok, tok, c, c)

				tok = c

				break
			}
		}
	}

	if t, ok := tok.(BooleanExpression); ok {
		return t.Evaluate()
	}

	panic(fmt.Sprintf("token %p(%#v) does not implement BooleanExpression interface", c.token, c.token))
}

// Token interface methods

// Clone returns a copy of the token and all its children
func (c *ExpressionPointer) Clone() token.Token {
	return &ExpressionPointer{
		token: c.token.Clone(),
	}
}

// Parse tries to parse the token beginning from the current position in the parser data.
// If the parsing is successful the error argument is nil and the next current position after the token is returned.
func (c *ExpressionPointer) Parse(pars *token.InternalParser, cur int) (int, []error) {
	panic("This should never happen")
}

// Permutation sets a specific permutation for this token
func (c *ExpressionPointer) Permutation(i uint) error {
	// do nothing

	return nil
}

// Permutations returns the number of permutations for this token
func (c *ExpressionPointer) Permutations() uint {
	return 1
}

// PermutationsAll returns the number of all possible permutations for this token including its children
func (c *ExpressionPointer) PermutationsAll() uint {
	return 1
}

func (c *ExpressionPointer) String() string {
	return c.token.String()
}

// ForwardToken interface methods

// Get returns the current referenced token
func (c *ExpressionPointer) Get() token.Token {
	return nil
}

// InternalGet returns the current referenced internal token
func (c *ExpressionPointer) InternalGet() token.Token {
	return c.token
}

// InternalLogicalRemove removes the referenced internal token and returns the replacement for the current token or nil if the current token should be removed.
func (c *ExpressionPointer) InternalLogicalRemove(tok token.Token) token.Token {
	if c.token == tok {
		return nil
	}

	return c
}

// InternalReplace replaces an old with a new internal token if it is referenced by this token. The error return argument is not nil, if the replacement is not suitable.
func (c *ExpressionPointer) InternalReplace(oldToken, newToken token.Token) error {
	if c.token == oldToken {
		c.token = newToken
	}

	return nil
}

// ScopeToken interface methods

// SetScope sets the scope of the token
func (c *ExpressionPointer) SetScope(variableScope *token.VariableScope) {
	tok := c.token

	if po, ok := tok.(*primitives.Pointer); ok {
		log.Debugf("Found pointer in ExpressionPointer %p(%#v)", c, c)

		for {
			c := po.InternalGet()
			c = c.Clone()
			_ = po.Set(c)

			po, ok = c.(*primitives.Pointer)
			if !ok {
				log.Debugf("Replaced pointer %p(%#v) with %p(%#v)", tok, tok, c, c)

				tok = c

				break
			}
		}
	}

	if t, ok := tok.(token.ScopeToken); ok {
		t.SetScope(variableScope)
	}
}
