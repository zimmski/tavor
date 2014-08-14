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

func (v *Variable) Name() string {
	return v.name
}

// Token interface methods

func (v *Variable) Clone() token.Token {
	return &Variable{
		name:  v.name,
		token: v.token.Clone(),
	}
}

func (v *Variable) Fuzz(r rand.Rand) {
	v.token.Fuzz(r)
}

func (v *Variable) FuzzAll(r rand.Rand) {
	v.Fuzz(r)
}

func (v *Variable) Parse(pars *token.InternalParser, cur int) (int, []error) {
	panic("TODO implement")
}

func (v *Variable) Permutation(i int) error {
	return v.token.Permutation(i)
}

func (v *Variable) Permutations() int {
	return v.token.Permutations()
}

func (v *Variable) PermutationsAll() int {
	return v.token.PermutationsAll()
}

func (v *Variable) String() string {
	return v.token.String()
}

// ForwardToken interface methods

func (v *Variable) Get() token.Token {
	return v.token
}

func (v *Variable) InternalGet() token.Token {
	return v.token
}

func (v *Variable) InternalLogicalRemove(tok token.Token) token.Token {
	if v.token == tok {
		return nil
	}

	return v
}

func (v *Variable) InternalReplace(oldToken, newToken token.Token) {
	if v.token == oldToken {
		v.token = newToken
	}
}

// IndexToken interface methods

func (v *Variable) Index() int {
	if p, ok := v.token.(token.IndexToken); ok {
		return p.Index()
	}

	return -1
}

// ResetToken interface methods

func (v *Variable) Reset() {
	if p, ok := v.token.(token.ResetToken); ok {
		p.Reset()
	}
}

// ScopeToken interface methods

func (v *Variable) SetScope(variableScope map[string]token.Token) {
	variableScope[v.name] = v
}

type VariableSave struct {
	Variable
}

// Overwrite Variable methods

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

func (v *VariableValue) Permutation(i int) error {
	// do nothing

	return nil
}

func (v *VariableValue) Permutations() int {
	return 1
}

func (v *VariableValue) PermutationsAll() int {
	return 1
}

func (v *VariableValue) String() string {
	return v.variable.InternalGet().String()
}

// ForwardToken interface methods

func (v *VariableValue) Get() token.Token {
	return nil
}

func (v *VariableValue) InternalGet() token.Token {
	return v.variable
}

func (v *VariableValue) InternalLogicalRemove(tok token.Token) token.Token {
	if v.variable == tok {
		return nil
	}

	return v
}

func (v *VariableValue) InternalReplace(oldToken, newToken token.Token) {
	if v.variable == oldToken {
		v.variable = newToken.(token.VariableToken)
	}
}

// IndexToken interface methods

func (v *VariableValue) Index() int {
	return v.variable.Index()
}

// ScopeToken interface methods

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
