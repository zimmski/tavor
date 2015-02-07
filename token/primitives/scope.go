package primitives

import (
	"github.com/zimmski/tavor/token"
)

// Scope implements a general scope token which references a token
type Scope struct {
	token token.Token
}

// NewScope returns a new instance of a Scope token
func NewScope(tok token.Token) *Scope {
	return &Scope{
		token: tok,
	}
}

// Token interface methods

// Clone returns a copy of the token and all its children
func (p *Scope) Clone() token.Token {
	return &Scope{
		token: p.token.Clone(),
	}
}

// Parse tries to parse the token beginning from the current position in the parser data.
// If the parsing is successful the error argument is nil and the next current position after the token is returned.
func (p *Scope) Parse(pars *token.InternalParser, cur int) (int, []error) {
	return p.token.Parse(pars, cur)
}

// Permutation sets a specific permutation for this token
func (p *Scope) Permutation(i uint) error {
	permutations := p.Permutations()

	if i < 1 || i > permutations {
		return &token.PermutationError{
			Type: token.PermutationErrorIndexOutOfBound,
		}
	}

	// do nothing

	return nil
}

// Permutations returns the number of permutations for this token
func (p *Scope) Permutations() uint {
	return 1
}

// PermutationsAll returns the number of all possible permutations for this token including its children
func (p *Scope) PermutationsAll() uint {
	return p.token.PermutationsAll()
}

func (p *Scope) String() string {
	return p.token.String()
}

// ForwardToken interface methods

// Get returns the current referenced token
func (p *Scope) Get() token.Token {
	return p.token
}

// InternalGet returns the current referenced internal token
func (p *Scope) InternalGet() token.Token {
	return p.token
}

// InternalLogicalRemove removes the referenced internal token and returns the replacement for the current token or nil if the current token should be removed.
func (p *Scope) InternalLogicalRemove(tok token.Token) token.Token {
	if p.token == tok {
		return nil
	}

	return p
}

// InternalReplace replaces an old with a new internal token if it is referenced by this token. The error return argument is not nil, if the replacement is not suitable.
func (p *Scope) InternalReplace(oldToken, newToken token.Token) error {
	if p.token == oldToken {
		p.token = newToken
	}

	return nil
}

// Minimize interface methods

// Minimize tries to minimize itself and returns a token if it was successful, or nil if there was nothing to minimize
func (p *Scope) Minimize() token.Token {
	if _, ok := p.token.(*Scope); ok {
		return p.token
	}

	return nil
}

// Resolve interface methods

// Resolve returns the token which is referenced by the token, or a path of tokens
func (p *Scope) Resolve() token.Token {
	var ok bool

	po := p

	for {
		c := po.InternalGet()

		po, ok = c.(*Scope)
		if !ok {
			return c
		}
	}
}

// Scoping interface methods

// Scoping returns if the token holds a new scope
func (p *Scope) Scoping() bool {
	return true
}
