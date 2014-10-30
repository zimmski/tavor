package variables

import (
	"github.com/zimmski/tavor/log"
	"github.com/zimmski/tavor/rand"
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
	if l, ok := v.token.(token.List); ok {
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

// Fuzz fuzzes this token using the random generator by choosing one of the possible permutations for this token
func (v *Variable) Fuzz(r rand.Rand) {
	// do nothing
}

// FuzzAll calls Fuzz for this token and then FuzzAll for all children of this token
func (v *Variable) FuzzAll(r rand.Rand) {
	v.Fuzz(r)

	v.token.Fuzz(r)
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

// InternalReplace replaces an old with a new internal token if it is referenced by this token
func (v *Variable) InternalReplace(oldToken, newToken token.Token) {
	if v.token == oldToken {
		v.token = newToken
	}
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
func (v *Variable) SetScope(variableScope map[string]token.Token) {
	variableScope[v.name] = v
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

// Fuzz fuzzes this token using the random generator by choosing one of the possible permutations for this token
func (v *VariableValue) Fuzz(r rand.Rand) {
	// do nothing
}

// FuzzAll calls Fuzz for this token and then FuzzAll for all children of this token
func (v *VariableValue) FuzzAll(r rand.Rand) {
	v.Fuzz(r)
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

// InternalReplace replaces an old with a new internal token if it is referenced by this token
func (v *VariableValue) InternalReplace(oldToken, newToken token.Token) {
	if v.variable == oldToken {
		v.variable = newToken.(token.VariableToken)
	}
}

// IndexToken interface methods

// Index returns the index of this token in its parent token
func (v *VariableValue) Index() int {
	return v.variable.Index()
}

// ScopeToken interface methods

// SetScope sets the scope of the token
func (v *VariableValue) SetScope(variableScope map[string]token.Token) {
	tok := variableScope[v.variable.Name()]

	if p, ok := tok.(*primitives.Pointer); ok {
		for {
			tok = p.InternalGet()

			p, ok = tok.(*primitives.Pointer)
			if !ok {
				break
			}
		}
	}

	if tok == nil { // TODO
		log.Debugf("TODO this should not happen")

		return
	}

	if t, ok := tok.(*VariableValue); ok {
		v.variable = t.variable
	} else {
		v.variable = tok.(token.VariableToken)
	}
}
