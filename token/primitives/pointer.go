package primitives

import (
	"fmt"
	"github.com/zimmski/tavor/rand"
	"github.com/zimmski/tavor/token"
	"reflect"
)

type Pointer struct {
	tok token.Token
	typ reflect.Type

	cloned bool
}

func NewPointer(tok token.Token) *Pointer {
	return &Pointer{
		tok: tok,
		typ: reflect.TypeOf(tok).Elem(),
	}
}

func NewEmptyPointer(typ interface{}) *Pointer {
	return &Pointer{
		tok: nil,
		typ: reflect.TypeOf(typ).Elem(),
	}
}

func (p *Pointer) Clone() token.Token {
	return &Pointer{
		tok:    p.tok, // do not clone further
		typ:    p.typ,
		cloned: false,
	}
}

func (p *Pointer) cloneOnFirstUse() {
	if !p.cloned {
		// clone everything on first use until we hit pointers
		p.tok = p.tok.Clone()

		p.cloned = true
	}
}

func (p *Pointer) Fuzz(r rand.Rand) {
	p.cloneOnFirstUse()
}

func (p *Pointer) FuzzAll(r rand.Rand) {
	p.Fuzz(r)

	// fuzz with the clone not the original token
	p.tok.FuzzAll(r)
}

func (p *Pointer) Get() token.Token {
	return p.tok
}

func (p *Pointer) Permutations() int {
	return 1 // TODO this could run forever if there is a loop so just return 1 for now
}

func (p *Pointer) Set(o token.Token) error {
	if !reflect.TypeOf(o).Implements(p.typ) {
		return fmt.Errorf("does not implement type %s", p.typ)
	}

	p.tok = o

	return nil
}

func (p *Pointer) String() string {
	return p.tok.String()
}
