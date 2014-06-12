package expressions

import (
	"strconv"

	"github.com/zimmski/tavor/rand"
	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/lists"
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
	// do nothing
}

func (e *AddArithmetic) FuzzAll(r rand.Rand) {
	e.Fuzz(r)

	e.a.FuzzAll(r)
	e.b.FuzzAll(r)
}

func (e *AddArithmetic) Get(i int) (token.Token, error) {
	switch i {
	case 0:
		return e.a, nil
	case 1:
		return e.b, nil
	default:
		return nil, &lists.ListError{lists.ListErrorOutOfBound}
	}
}

func (e *AddArithmetic) Len() int {
	return 2
}

func (e *AddArithmetic) Permutations() int {
	return e.a.Permutations() * e.b.Permutations()
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
	// do nothing
}

func (e *SubArithmetic) FuzzAll(r rand.Rand) {
	e.Fuzz(r)

	e.a.FuzzAll(r)
	e.b.FuzzAll(r)
}

func (e *SubArithmetic) Get(i int) (token.Token, error) {
	switch i {
	case 0:
		return e.a, nil
	case 1:
		return e.b, nil
	default:
		return nil, &lists.ListError{lists.ListErrorOutOfBound}
	}
}

func (e *SubArithmetic) Len() int {
	return 2
}

func (e *SubArithmetic) Permutations() int {
	return e.a.Permutations() * e.b.Permutations()
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
	// do nothing
}

func (e *MulArithmetic) FuzzAll(r rand.Rand) {
	e.Fuzz(r)

	e.a.FuzzAll(r)
	e.b.FuzzAll(r)
}

func (e *MulArithmetic) Get(i int) (token.Token, error) {
	switch i {
	case 0:
		return e.a, nil
	case 1:
		return e.b, nil
	default:
		return nil, &lists.ListError{lists.ListErrorOutOfBound}
	}
}

func (e *MulArithmetic) Len() int {
	return 2
}

func (e *MulArithmetic) Permutations() int {
	return e.a.Permutations() * e.b.Permutations()
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
	// do nothing
}

func (e *DivArithmetic) FuzzAll(r rand.Rand) {
	e.Fuzz(r)

	e.a.FuzzAll(r)
	e.b.FuzzAll(r)
}

func (e *DivArithmetic) Get(i int) (token.Token, error) {
	switch i {
	case 0:
		return e.a, nil
	case 1:
		return e.b, nil
	default:
		return nil, &lists.ListError{lists.ListErrorOutOfBound}
	}
}

func (e *DivArithmetic) Len() int {
	return 2
}

func (e *DivArithmetic) Permutations() int {
	return e.a.Permutations() * e.b.Permutations()
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
