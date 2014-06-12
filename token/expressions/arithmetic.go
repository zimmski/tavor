package expressions

import (
	"strconv"

	"github.com/zimmski/tavor/rand"
	"github.com/zimmski/tavor/token"
)

type AddArithmetic struct {
	a token.Token
	b token.Token
}

func NewAddArithmetic(a, b token.Token) *AddArithmetic {
	return &AddArithmetic{
		a: a,
		b: b,
	}
}

func (e *AddArithmetic) Clone() token.Token {
	return &AddArithmetic{
		a: e.a,
		b: e.b,
	}
}

func (e *AddArithmetic) Fuzz(r rand.Rand) {
	e.a.Fuzz(r)
	e.b.Fuzz(r)
}

func (e *AddArithmetic) Permutations() int {
	return 1
}

func (e *AddArithmetic) String() string {
	a, err := strconv.Atoi(e.a.String())
	if err != nil {
		panic(err)
	}
	b, err := strconv.Atoi(e.b.String())
	if err != nil {
		panic(err)
	}

	return strconv.Itoa(a + b)
}

type SubArithmetic struct {
	a token.Token
	b token.Token
}

func NewSubArithmetic(a, b token.Token) *SubArithmetic {
	return &SubArithmetic{
		a: a,
		b: b,
	}
}

func (e *SubArithmetic) Clone() token.Token {
	return &SubArithmetic{
		a: e.a,
		b: e.b,
	}
}

func (e *SubArithmetic) Fuzz(r rand.Rand) {
	e.a.Fuzz(r)
	e.b.Fuzz(r)
}

func (e *SubArithmetic) Permutations() int {
	return 1
}

func (e *SubArithmetic) String() string {
	a, err := strconv.Atoi(e.a.String())
	if err != nil {
		panic(err)
	}
	b, err := strconv.Atoi(e.b.String())
	if err != nil {
		panic(err)
	}

	return strconv.Itoa(a - b)
}

type MulArithmetic struct {
	a token.Token
	b token.Token
}

func NewMulArithmetic(a, b token.Token) *MulArithmetic {
	return &MulArithmetic{
		a: a,
		b: b,
	}
}

func (e *MulArithmetic) Clone() token.Token {
	return &MulArithmetic{
		a: e.a,
		b: e.b,
	}
}

func (e *MulArithmetic) Fuzz(r rand.Rand) {
	e.a.Fuzz(r)
	e.b.Fuzz(r)
}

func (e *MulArithmetic) Permutations() int {
	return 1
}

func (e *MulArithmetic) String() string {
	a, err := strconv.Atoi(e.a.String())
	if err != nil {
		panic(err)
	}
	b, err := strconv.Atoi(e.b.String())
	if err != nil {
		panic(err)
	}

	return strconv.Itoa(a * b)
}

type DivArithmetic struct {
	a token.Token
	b token.Token
}

func NewDivArithmetic(a, b token.Token) *DivArithmetic {
	return &DivArithmetic{
		a: a,
		b: b,
	}
}

func (e *DivArithmetic) Clone() token.Token {
	return &DivArithmetic{
		a: e.a,
		b: e.b,
	}
}

func (e *DivArithmetic) Fuzz(r rand.Rand) {
	e.a.Fuzz(r)
	e.b.Fuzz(r)
}

func (e *DivArithmetic) Permutations() int {
	return 1
}

func (e *DivArithmetic) String() string {
	a, err := strconv.Atoi(e.a.String())
	if err != nil {
		panic(err)
	}
	b, err := strconv.Atoi(e.b.String())
	if err != nil {
		panic(err)
	}

	return strconv.Itoa(a / b)
}
