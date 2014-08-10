package primitives

import (
	"fmt"
	"reflect"

	"github.com/zimmski/tavor/rand"
	"github.com/zimmski/tavor/token"
)

type Pointer struct {
	token token.Token
	typ   reflect.Type

	cloned bool
}

func NewPointer(tok token.Token) *Pointer {
	return &Pointer{
		token: tok,
		typ:   reflect.TypeOf(tok).Elem(),
	}
}

func NewEmptyPointer(typ interface{}) *Pointer {
	return &Pointer{
		token: nil,
		typ:   reflect.TypeOf(typ).Elem(),
	}
}

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
		p.token = p.token.Clone()

		p.cloned = true
	}
}

func (p *Pointer) Fuzz(r rand.Rand) {
	p.cloneOnFirstUse()
}

func (p *Pointer) FuzzAll(r rand.Rand) {
	p.Fuzz(r)

	if p.token == nil {
		return
	}

	// fuzz with the clone not the original token
	p.token.FuzzAll(r)
}

func (p *Pointer) Parse(pars *token.InternalParser, cur int) (int, []error) {
	panic("Pointer token is not allowed during internal parsing")
}

func (p *Pointer) Permutation(i int) error {
	p.cloneOnFirstUse()

	permutations := p.Permutations()

	if i < 1 || i > permutations {
		return &token.PermutationError{
			Type: token.PermutationErrorIndexOutOfBound,
		}
	}

	// do nothing

	return nil
}

func (p *Pointer) Permutations() int {
	p.cloneOnFirstUse()

	return 1
}

func (p *Pointer) PermutationsAll() int {
	p.cloneOnFirstUse()

	if p.token == nil {
		return 1
	}

	return p.token.PermutationsAll()
}

func (p *Pointer) String() string {
	if p.token == nil {
		return ""
	}

	return p.token.String()
}

// ForwardToken interface methods

func (p *Pointer) Get() token.Token {
	return p.token
}

func (p *Pointer) InternalGet() token.Token {
	return p.token
}

func (p *Pointer) InternalLogicalRemove(tok token.Token) token.Token {
	if p.token == tok {
		return nil
	}

	return p
}

func (p *Pointer) InternalReplace(oldToken, newToken token.Token) {
	if p.token == oldToken {
		p.token = newToken
	}
}
