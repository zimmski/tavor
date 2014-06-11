package primitives

import (
	"github.com/zimmski/tavor/rand"
	"github.com/zimmski/tavor/token"
)

type Pointer struct {
	Tok token.Token
}

func NewPointer(tok token.Token) *Pointer {
	return &Pointer{
		Tok: tok,
	}
}

func NewEmptyPointer() *Pointer {
	return &Pointer{
		Tok: nil,
	}
}

func (p *Pointer) Clone() token.Token {
	return &Pointer{
		Tok: p.Tok, // do not clone further
	}
}

func (p *Pointer) Fuzz(r rand.Rand) {
	// clone everything until we hit pointers
	p.Tok = p.Tok.Clone()

	// fuzz with the clone not the original token
	p.Tok.Fuzz(r)
}

func (p *Pointer) String() string {
	return p.Tok.String()
}
