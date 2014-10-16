package variables

import (
	"github.com/zimmski/tavor/log"
	"github.com/zimmski/tavor/rand"
	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/primitives"
)

type Variable struct {
	name  string
	token token.Token
}

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

// Token interface methods

// Clone returns a copy of the token and all its children
func (v *Variable) Clone() token.Token {
	return &Variable{
		name:  v.name,
		token: v.token.Clone(),
	}
}

func (v *Variable) Fuzz(r rand.Rand) {
	// do nothing
}

func (v *Variable) FuzzAll(r rand.Rand) {
	v.Fuzz(r)

	v.token.Fuzz(r)
}

func (v *Variable) Parse(pars *token.InternalParser, cur int) (int, []error) {
	panic("TODO implement")
}

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

func (v *Variable) Permutations() uint {
	return 1
}

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

func NewVariableSave(name string, token token.Token) *VariableSave {
	return &VariableSave{
		Variable: Variable{
			name:  name,
			token: token,
		},
	}
}

type VariableValue struct {
	variable token.VariableToken
}

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

func (v *VariableValue) Fuzz(r rand.Rand) {
	// do nothing
}

func (v *VariableValue) FuzzAll(r rand.Rand) {
	v.Fuzz(r)
}

func (v *VariableValue) Parse(pars *token.InternalParser, cur int) (int, []error) {
	panic("TODO implement")
}

func (v *VariableValue) Permutation(i uint) error {
	// do nothing

	return nil
}

func (v *VariableValue) Permutations() uint {
	return 1
}

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
