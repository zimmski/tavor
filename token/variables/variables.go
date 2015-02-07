package variables

import (
	"strconv"

	"github.com/zimmski/tavor/log"
	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/primitives"
)

// Variable implements general variable token which references a token as its value and forwards all token functions to its token.
type Variable struct {
	name  string
	token token.Token
}

// NewVariable returns a new instance of a Variable token
func NewVariable(name string, token token.Token) *Variable {
	return &Variable{
		name:  name,
		token: token,
	}
}

// Name returns the name of the variable
func (v *Variable) Name() string {
	return v.name
}

// Len returns the number of the current referenced tokens
func (v *Variable) Len() int {
	if l, ok := v.token.(token.ListToken); ok {
		return l.Len()
	}

	return 1
}

// Token interface methods

// Clone returns a copy of the token and all its children
func (v *Variable) Clone() token.Token {
	return &Variable{
		name:  v.name,
		token: v.token.Clone(),
	}
}

// Parse tries to parse the token beginning from the current position in the parser data.
// If the parsing is successful the error argument is nil and the next current position after the token is returned.
func (v *Variable) Parse(pars *token.InternalParser, cur int) (int, []error) {
	panic("TODO implement")
}

// Permutation sets a specific permutation for this token
func (v *Variable) Permutation(i uint) error {
	permutations := v.Permutations()

	if i < 1 || i > permutations {
		return &token.PermutationError{
			Type: token.PermutationErrorIndexOutOfBound,
		}
	}
	// do nothing

	return nil
}

// Permutations returns the number of permutations for this token
func (v *Variable) Permutations() uint {
	return 1
}

// PermutationsAll returns the number of all possible permutations for this token including its children
func (v *Variable) PermutationsAll() uint {
	return v.token.PermutationsAll()
}

func (v *Variable) String() string {
	return v.token.String()
}

// ForwardToken interface methods

// Get returns the current referenced token
func (v *Variable) Get() token.Token {
	return v.token
}

// InternalGet returns the current referenced internal token
func (v *Variable) InternalGet() token.Token {
	return v.token
}

// InternalLogicalRemove removes the referenced internal token and returns the replacement for the current token or nil if the current token should be removed.
func (v *Variable) InternalLogicalRemove(tok token.Token) token.Token {
	if v.token == tok {
		return nil
	}

	return v
}

// InternalReplace replaces an old with a new internal token if it is referenced by this token. The error return argument is not nil, if the replacement is not suitable.
func (v *Variable) InternalReplace(oldToken, newToken token.Token) error {
	if v.token == oldToken {
		v.token = newToken
	}

	return nil
}

// IndexToken interface methods

// Index returns the index of this token in its parent token
func (v *Variable) Index() int {
	if p, ok := v.token.(token.IndexToken); ok {
		return p.Index()
	}

	return -1
}

// ScopeToken interface methods

// SetScope sets the scope of the token
func (v *Variable) SetScope(variableScope *token.VariableScope) {
	variableScope.Set(v.name, v)
}

// VariableItem implements a token which references a Variable token to output its referenced token
type VariableItem struct {
	index    token.Token
	variable token.VariableToken
}

// NewVariableItem returns a new instance of a VariableItem token
func NewVariableItem(index token.Token, variable token.VariableToken) *VariableItem {
	return &VariableItem{
		index:    index,
		variable: variable,
	}
}

// Token interface methods

// Clone returns a copy of the token and all its children
func (v *VariableItem) Clone() token.Token {
	return &VariableItem{
		index:    v.index.Clone(),
		variable: v.variable,
	}
}

// Parse tries to parse the token beginning from the current position in the parser data.
// If the parsing is successful the error argument is nil and the next current position after the token is returned.
func (v *VariableItem) Parse(pars *token.InternalParser, cur int) (int, []error) {
	panic("TODO implement")
}

// Permutation sets a specific permutation for this token
func (v *VariableItem) Permutation(i uint) error {
	// do nothing

	return nil
}

// Permutations returns the number of permutations for this token
func (v *VariableItem) Permutations() uint {
	return 1
}

// PermutationsAll returns the number of all possible permutations for this token including its children
func (v *VariableItem) PermutationsAll() uint {
	return 1
}

func (v *VariableItem) String() string {
	i := v.Index()

	l, ok := v.variable.Get().(token.ListToken)
	if !ok {
		// TODO

		return ""
	}

	tok, err := l.Get(i)
	if err != nil {
		panic(err) // TODO
	}

	return tok.String()
}

// ForwardToken interface methods

// Get returns the current referenced token
func (v *VariableItem) Get() token.Token {
	return nil
}

// InternalGet returns the current referenced internal token
func (v *VariableItem) InternalGet() token.Token {
	return v.variable
}

// InternalLogicalRemove removes the referenced internal token and returns the replacement for the current token or nil if the current token should be removed.
func (v *VariableItem) InternalLogicalRemove(tok token.Token) token.Token {
	if v.variable == tok {
		return nil
	}

	return v
}

// InternalReplace replaces an old with a new internal token if it is referenced by this token. The error return argument is not nil, if the replacement is not suitable.
func (v *VariableItem) InternalReplace(oldToken, newToken token.Token) error {
	if v.variable == oldToken {
		v.variable = newToken.(token.VariableToken)
	}

	return nil
}

// Follow returns if the children of the token should be traversed
func (v *VariableItem) Follow() bool {
	return false
}

// IndexToken interface methods

// Index returns the index of this token in its parent token
func (v *VariableItem) Index() int {
	i, err := strconv.Atoi(v.index.String())
	if err != nil {
		panic(err) // TODO
	}

	return i
}

// ScopeToken interface methods

// SetScope sets the scope of the token
func (v *VariableItem) SetScope(variableScope *token.VariableScope) {
	if tok, ok := v.index.(token.ScopeToken); ok {
		tok.SetScope(variableScope)
	}

	tok := variableScope.Get(v.variable.Name())

	if p, ok := tok.(*primitives.Pointer); ok {
		tok = p.Resolve()
	}

	if tok == nil { // TODO
		log.Debugf("TODO VariableItem: this should not happen")

		return
	}

	if t, ok := tok.(*VariableItem); ok {
		v.variable = t.variable
	} else {
		v.variable = tok.(token.VariableToken)
	}
}

// VariableSave is based on the general Variable token but does prevent the output of the referenced token
type VariableSave struct {
	Variable
}

// Overwrite Variable methods

// Clone returns a copy of the token and all its children
func (v *VariableSave) Clone() token.Token {
	return &VariableSave{
		Variable: Variable{
			name:  v.name,
			token: v.token.Clone(),
		},
	}
}

func (v *VariableSave) String() string {
	return ""
}

// NewVariableSave returns a new instance of a VariableSave token
func NewVariableSave(name string, token token.Token) *VariableSave {
	return &VariableSave{
		Variable: Variable{
			name:  name,
			token: token,
		},
	}
}

// VariableReference implements a token which references a Variable token to output its referenced token
type VariableReference struct {
	variable token.VariableToken
}

// NewVariableReference returns a new instance of a VariableReference token
func NewVariableReference(variable token.VariableToken) *VariableReference {
	return &VariableReference{
		variable: variable,
	}
}

// Reference returns the referenced token
func (v *VariableReference) Reference() token.Token {
	return v.variable.Get()
}

// Token interface methods

// Clone returns a copy of the token and all its children
func (v *VariableReference) Clone() token.Token {
	return &VariableReference{
		variable: v.variable,
	}
}

// Parse tries to parse the token beginning from the current position in the parser data.
// If the parsing is successful the error argument is nil and the next current position after the token is returned.
func (v *VariableReference) Parse(pars *token.InternalParser, cur int) (int, []error) {
	panic("TODO implement")
}

// Permutation sets a specific permutation for this token
func (v *VariableReference) Permutation(i uint) error {
	// do nothing

	return nil
}

// Permutations returns the number of permutations for this token
func (v *VariableReference) Permutations() uint {
	return 1
}

// PermutationsAll returns the number of all possible permutations for this token including its children
func (v *VariableReference) PermutationsAll() uint {
	return 1
}

func (v *VariableReference) String() string {
	return ""
}

// Follow returns if the children of the token should be traversed
func (v *VariableReference) Follow() bool {
	return false
}

// IndexToken interface methods

// Index returns the index of this token in its parent token
func (v *VariableReference) Index() int {
	return v.variable.Index()
}

// ScopeToken interface methods

// SetScope sets the scope of the token
func (v *VariableReference) SetScope(variableScope *token.VariableScope) {
	tok := variableScope.Get(v.variable.Name())

	if p, ok := tok.(*primitives.Pointer); ok {
		tok = p.Resolve()
	}

	if tok == nil { // TODO
		log.Debugf("TODO VariableReference: this should not happen")

		return
	}

	if t, ok := tok.(*VariableReference); ok {
		v.variable = t.variable
	} else {
		v.variable = tok.(token.VariableToken)
	}
}

// VariableValue implements a token which references a Variable token to output its referenced token
type VariableValue struct {
	variable token.VariableToken
}

// NewVariableValue returns a new instance of a VariableValue token
func NewVariableValue(variable token.VariableToken) *VariableValue {
	return &VariableValue{
		variable: variable,
	}
}

// Token interface methods

// Clone returns a copy of the token and all its children
func (v *VariableValue) Clone() token.Token {
	return &VariableValue{
		variable: v.variable,
	}
}

// Parse tries to parse the token beginning from the current position in the parser data.
// If the parsing is successful the error argument is nil and the next current position after the token is returned.
func (v *VariableValue) Parse(pars *token.InternalParser, cur int) (int, []error) {
	panic("TODO implement")
}

// Permutation sets a specific permutation for this token
func (v *VariableValue) Permutation(i uint) error {
	// do nothing

	return nil
}

// Permutations returns the number of permutations for this token
func (v *VariableValue) Permutations() uint {
	return 1
}

// PermutationsAll returns the number of all possible permutations for this token including its children
func (v *VariableValue) PermutationsAll() uint {
	return 1
}

func (v *VariableValue) String() string {
	return v.variable.InternalGet().String()
}

// ForwardToken interface methods

// Get returns the current referenced token
func (v *VariableValue) Get() token.Token {
	return nil
}

// InternalGet returns the current referenced internal token
func (v *VariableValue) InternalGet() token.Token {
	return v.variable
}

// InternalLogicalRemove removes the referenced internal token and returns the replacement for the current token or nil if the current token should be removed.
func (v *VariableValue) InternalLogicalRemove(tok token.Token) token.Token {
	if v.variable == tok {
		return nil
	}

	return v
}

// InternalReplace replaces an old with a new internal token if it is referenced by this token. The error return argument is not nil, if the replacement is not suitable.
func (v *VariableValue) InternalReplace(oldToken, newToken token.Token) error {
	if v.variable == oldToken {
		v.variable = newToken.(token.VariableToken)
	}

	return nil
}

// IndexToken interface methods

// Index returns the index of this token in its parent token
func (v *VariableValue) Index() int {
	return v.variable.Index()
}

// ScopeToken interface methods

// SetScope sets the scope of the token
func (v *VariableValue) SetScope(variableScope *token.VariableScope) {
	tok := variableScope.Get(v.variable.Name())

	if p, ok := tok.(*primitives.Pointer); ok {
		tok = p.Resolve()
	}

	if tok == nil { // TODO
		log.Debugf("TODO VariableValue: this should not happen")

		return
	}

	if t, ok := tok.(*VariableValue); ok {
		v.variable = t.variable
	} else {
		v.variable = tok.(token.VariableToken)
	}
}
