package variables

import (
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

// Token interface methods

func (v *Variable) Clone() token.Token {
	/*return &Variable{
		token: v.token,
	}*/
	return v
}

func (v *Variable) Fuzz(r rand.Rand) {
	// do nothing
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
	return nil
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

// ScopeToken interface methods

func (v *Variable) SetScope(variableScope map[string]token.Token) {
	variableScope[v.name] = v
}

type VariableValue struct {
	variable *Variable
}

func NewVariableValue(variable *Variable) *VariableValue {
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

// ScopeToken interface methods

func (v *VariableValue) SetScope(variableScope map[string]token.Token) {
	tok := variableScope[v.variable.name]

	if p, ok := tok.(*primitives.Pointer); ok {
		for {
			tok = p.InternalGet()

			p, ok = tok.(*primitives.Pointer)
			if !ok {
				break
			}
		}
	}

	if t, ok := tok.(*VariableValue); ok {
		v.variable = t.variable
	} else {
		v.variable = tok.(*Variable)
	}
}
