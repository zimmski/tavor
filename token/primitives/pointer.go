package primitives

import (
	"fmt"
	"github.com/zimmski/tavor/token"
	"reflect"
)

// Pointer implements a general pointer token which references a token
type Pointer struct {
	token token.Token
	typ   reflect.Type

	cloned bool
}

// NewPointer returns a new instance of a Pointer token and sets the token reference type to the token's type
func NewPointer(tok token.Token) *Pointer {
	return &Pointer{
		token: tok,
		typ:   reflect.TypeOf(tok).Elem(),
	}
}

// NewEmptyPointer returns a new instance of a Pointer token with a token reference type but without a referenced token
func NewEmptyPointer(typ interface{}) *Pointer {
	return &Pointer{
		token: nil,
		typ:   reflect.TypeOf(typ).Elem(),
	}
}

// NewTokenPointer returns a new instance of a Pointer token with the token reference type Token but without a referenced token
func NewTokenPointer(tok token.Token) *Pointer {
	var tokenType *token.Token

	return &Pointer{
		token: tok,
		typ:   reflect.TypeOf(tokenType).Elem(),
	}
}

// Set sets the referenced token which must conform to the pointers token reference type
func (p *Pointer) Set(o token.Token) error {
	if o == nil {
		p.token = nil
		p.cloned = true
		return nil
	}

	oType := reflect.TypeOf(o)

	if !oType.AssignableTo(p.typ) && (p.typ.Kind() == reflect.Interface && !oType.Implements(p.typ)) {
		return fmt.Errorf("does not implement type %s", p.typ)
	}

	p.token = o

	return nil
}

// Token interface methods

// Clone returns a copy of the token and all its children
func (p *Pointer) Clone() token.Token {
	return &Pointer{
		token:  p.token, // do not clone further
		typ:    p.typ,
		cloned: false,
	}
}

func (p *Pointer) cloneOnFirstUse() {
	if !p.cloned && p.token != nil {
		// clone everything on first use until we hit pointers
		if _, ok := p.token.(*Pointer); !ok {
			p.token = p.token.Clone()

			p.cloned = true
		}
	}
}

// Parse tries to parse the token beginning from the current position in the parser data.
// If the parsing is successful the error argument is nil and the next current position after the token is returned.
func (p *Pointer) Parse(pars *token.InternalParser, cur int) (int, []error) {
	panic("Pointer token is not allowed during internal parsing")
}

// Permutation sets a specific permutation for this token
func (p *Pointer) Permutation(i uint) error {
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
func (p *Pointer) Permutations() uint {
	p.cloneOnFirstUse()

	return 1
}

// PermutationsAll returns the number of all possible permutations for this token including its children
func (p *Pointer) PermutationsAll() uint {
	p.cloneOnFirstUse()

	if p.token == nil {
		panic("Pointer token does not have a referencing token")
	}

	return p.token.PermutationsAll()
}

func (p *Pointer) String() string {
	if p.token == nil {
		panic("Pointer token does not have a referencing token")
	}

	return p.token.String()
}

// ForwardToken interface methods

// Get returns the current referenced token
func (p *Pointer) Get() token.Token {
	p.cloneOnFirstUse()

	return p.token
}

// InternalGet returns the current referenced internal token
func (p *Pointer) InternalGet() token.Token {
	return p.token
}

// InternalLogicalRemove removes the referenced internal token and returns the replacement for the current token or nil if the current token should be removed.
func (p *Pointer) InternalLogicalRemove(tok token.Token) token.Token {
	if p.token == tok {
		return nil
	}

	return p
}

// InternalReplace replaces an old with a new internal token if it is referenced by this token. The error return argument is not nil, if the replacement is not suitable.
func (p *Pointer) InternalReplace(oldToken, newToken token.Token) error {
	if p.token == oldToken {
		p.token = newToken
	}

	return nil
}

// BooleanExpression interface methods

/*func (p *Pointer) Evaluate() bool {
	if tok, ok := p.token.(conditions.BooleanExpression); ok {
		return tok.Evaluate()
	} else {
		panic(fmt.Errorf("TODO token %p(%#v) is not a BooleanExpression", p.token, p.token))
	}
}*/

// Minimize interface methods

// Minimize tries to minimize itself and returns a token if it was successful, or nil if there was nothing to minimize
func (p *Pointer) Minimize() token.Token {
	// Never ever _EVER_ minimize a pointer since it is normally there for a reason

	return nil
}

// Resolve interface methods

// Resolve returns the token which is referenced by the token, or a path of tokens
func (p *Pointer) Resolve() token.Token {
	var ok bool

	po := p

	for {
		c := po.InternalGet()

		po, ok = c.(*Pointer)
		if !ok {
			return c
		}
	}
}
