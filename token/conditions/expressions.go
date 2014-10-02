package conditions

import (
	"fmt"

	"github.com/zimmski/tavor/log"
	"github.com/zimmski/tavor/rand"
	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/lists"
	"github.com/zimmski/tavor/token/primitives"
)

type BooleanExpression interface {
	token.Token

	Evaluate() bool
}

type BooleanTrue struct{}

func NewBooleanTrue() *BooleanTrue {
	return &BooleanTrue{}
}

func (c *BooleanTrue) Evaluate() bool {
	return true
}

// Token interface methods

func (c *BooleanTrue) Clone() token.Token {
	return &BooleanTrue{}
}

func (c *BooleanTrue) Fuzz(r rand.Rand) {
	// do nothing
}

func (c *BooleanTrue) FuzzAll(r rand.Rand) {
	// do nothing
}

func (c *BooleanTrue) Parse(pars *token.InternalParser, cur int) (int, []error) {
	panic("This should never happen")
}

func (c *BooleanTrue) Permutation(i uint) error {
	// do nothing

	return nil
}

func (c *BooleanTrue) Permutations() uint {
	return 1
}

func (c *BooleanTrue) PermutationsAll() uint {
	return 1
}

func (c *BooleanTrue) String() string {
	return "true"
}

// List interface methods

func (c *BooleanTrue) Get(i int) (token.Token, error) {
	return nil, &lists.ListError{
		Type: lists.ListErrorOutOfBound,
	}
}

func (c *BooleanTrue) Len() int {
	return 0
}

func (c *BooleanTrue) InternalGet(i int) (token.Token, error) {
	return nil, &lists.ListError{
		Type: lists.ListErrorOutOfBound,
	}
}

func (c *BooleanTrue) InternalLen() int {
	return 0
}

func (c *BooleanTrue) InternalLogicalRemove(tok token.Token) token.Token {
	panic("This should never happen")
}

func (c *BooleanTrue) InternalReplace(oldToken, newToken token.Token) {
	panic("This should never happen")
}

type BooleanEqual struct {
	a, b token.Token
}

func NewBooleanEqual(a, b token.Token) *BooleanEqual {
	return &BooleanEqual{
		a: a,
		b: b,
	}
}

func (c *BooleanEqual) Evaluate() bool {
	return c.a.String() == c.b.String()
}

// Token interface methods

func (c *BooleanEqual) Clone() token.Token {
	return &BooleanEqual{
		a: c.a,
		b: c.b,
	}
}

func (c *BooleanEqual) Fuzz(r rand.Rand) {
	// do nothing
}

func (c *BooleanEqual) FuzzAll(r rand.Rand) {
	// do nothing
}

func (c *BooleanEqual) Parse(pars *token.InternalParser, cur int) (int, []error) {
	panic("This should never happen")
}

func (c *BooleanEqual) Permutation(i uint) error {
	// do nothing

	return nil
}

func (c *BooleanEqual) Permutations() uint {
	return 1
}

func (c *BooleanEqual) PermutationsAll() uint {
	return 1
}

func (c *BooleanEqual) String() string {
	return fmt.Sprintf("(%p)%#v == (%p)%#v", c.a, c.a, c.b, c.b)
}

// List interface methods

func (c *BooleanEqual) Get(i int) (token.Token, error) {
	return nil, &lists.ListError{
		Type: lists.ListErrorOutOfBound,
	}
}

func (c *BooleanEqual) Len() int {
	return 0
}

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

func (c *BooleanEqual) InternalLen() int {
	return 2
}

func (c *BooleanEqual) InternalLogicalRemove(tok token.Token) token.Token {
	if tok == c.a || tok == c.b {
		return nil
	}

	return c
}

func (c *BooleanEqual) InternalReplace(oldToken, newToken token.Token) {
	if oldToken == c.a {
		c.a = newToken
	}
	if oldToken == c.b {
		c.b = newToken
	}
}

type VariableDefined struct {
	name          string
	variableScope map[string]token.Token
}

func NewVariableDefined(name string, variableScope map[string]token.Token) *VariableDefined {
	return &VariableDefined{
		name:          name,
		variableScope: variableScope,
	}
}

func (c *VariableDefined) Evaluate() bool {
	_, ok := c.variableScope[c.name]

	return ok
}

// Token interface methods

func (c *VariableDefined) Clone() token.Token {
	return &VariableDefined{
		name:          c.name,
		variableScope: c.variableScope,
	}
}

func (c *VariableDefined) Fuzz(r rand.Rand) {
	// do nothing
}

func (c *VariableDefined) FuzzAll(r rand.Rand) {
	// do nothing
}

func (c *VariableDefined) Parse(pars *token.InternalParser, cur int) (int, []error) {
	panic("This should never happen")
}

func (c *VariableDefined) Permutation(i uint) error {
	// do nothing

	return nil
}

func (c *VariableDefined) Permutations() uint {
	return 1
}

func (c *VariableDefined) PermutationsAll() uint {
	return 1
}

func (c *VariableDefined) String() string {
	return fmt.Sprintf("defined(%q)", c.name)
}

// ScopeToken interface methods

func (c *VariableDefined) SetScope(variableScope map[string]token.Token) {
	nScope := make(map[string]token.Token, len(variableScope))
	for k, v := range variableScope {
		nScope[k] = v
	}

	c.variableScope = nScope
}

type ExpressionPointer struct {
	token token.Token
}

func NewExpressionPointer(token token.Token) *ExpressionPointer {
	return &ExpressionPointer{
		token: token,
	}
}

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

func (c *ExpressionPointer) Clone() token.Token {
	return &ExpressionPointer{
		token: c.token.Clone(),
	}
}

func (c *ExpressionPointer) Fuzz(r rand.Rand) {
	// do nothing
}

func (c *ExpressionPointer) FuzzAll(r rand.Rand) {
	// do nothing
}

func (c *ExpressionPointer) Parse(pars *token.InternalParser, cur int) (int, []error) {
	panic("This should never happen")
}

func (c *ExpressionPointer) Permutation(i uint) error {
	// do nothing

	return nil
}

func (c *ExpressionPointer) Permutations() uint {
	return 1
}

func (c *ExpressionPointer) PermutationsAll() uint {
	return 1
}

func (c *ExpressionPointer) String() string {
	return c.token.String()
}

// ForwardToken interface methods

func (c *ExpressionPointer) Get() token.Token {
	return nil
}

func (c *ExpressionPointer) InternalGet() token.Token {
	return c.token
}

func (c *ExpressionPointer) InternalLogicalRemove(tok token.Token) token.Token {
	if c.token == tok {
		return nil
	}

	return c
}

func (c *ExpressionPointer) InternalReplace(oldToken, newToken token.Token) {
	if c.token == oldToken {
		c.token = newToken
	}
}

// ScopeToken interface methods

func (c *ExpressionPointer) SetScope(variableScope map[string]token.Token) {
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
