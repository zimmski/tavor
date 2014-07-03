package primitives

import (
	"strconv"

	"github.com/zimmski/tavor/rand"
	"github.com/zimmski/tavor/token"
)

type ConstantInt struct {
	value int
}

func NewConstantInt(value int) *ConstantInt {
	return &ConstantInt{
		value: value,
	}
}

func (p *ConstantInt) Clone() token.Token {
	return &ConstantInt{
		value: p.value,
	}
}

func (p *ConstantInt) Fuzz(r rand.Rand) {
	// do nothing
}

func (p *ConstantInt) FuzzAll(r rand.Rand) {
	p.Fuzz(r)
}

func (p *ConstantInt) Parse(parser token.InternalParser, cur *token.ParserList) []token.ParserList {
	panic("TODO implement")
}

func (p *ConstantInt) Permutation(i int) error {
	permutations := p.Permutations()

	if i < 1 || i > permutations {
		return &token.PermutationError{
			Type: token.PermutationErrorIndexOutOfBound,
		}
	}

	// do nothing

	return nil
}

func (p *ConstantInt) Permutations() int {
	return 1
}

func (p *ConstantInt) PermutationsAll() int {
	return p.Permutations()
}

func (p *ConstantInt) String() string {
	return strconv.Itoa(p.value)
}

type RandomInt struct {
	value int
}

func NewRandomInt() *RandomInt {
	return &RandomInt{
		value: 0,
	}
}

func (p *RandomInt) Clone() token.Token {
	return &RandomInt{
		value: p.value,
	}
}

func (p *RandomInt) Fuzz(r rand.Rand) {
	p.value = r.Int()
}

func (p *RandomInt) FuzzAll(r rand.Rand) {
	p.Fuzz(r)
}

func (p *RandomInt) Parse(parser token.InternalParser, cur *token.ParserList) []token.ParserList {
	panic("TODO implement")
}

func (p *RandomInt) Permutation(i int) error {
	permutations := p.Permutations()

	if i < 1 || i > permutations {
		return &token.PermutationError{
			Type: token.PermutationErrorIndexOutOfBound,
		}
	}

	// TODO this could be done MUCH better
	p.value = 0

	return nil
}

func (p *RandomInt) Permutations() int {
	return 1 // TODO maybe this should be like RangeInt
}

func (p *RandomInt) PermutationsAll() int {
	return p.Permutations()
}

func (p *RandomInt) String() string {
	return strconv.Itoa(p.value)
}

type RangeInt struct {
	from int
	to   int

	value int
}

func NewRangeInt(from, to int) *RangeInt {
	return &RangeInt{
		from:  from,
		to:    to,
		value: from,
	}
}

func (p *RangeInt) Clone() token.Token {
	return &RangeInt{
		from:  p.from,
		to:    p.to,
		value: p.value,
	}
}

func (p *RangeInt) Fuzz(r rand.Rand) {
	i := r.Intn(p.Permutations())

	p.permutation(i)
}

func (p *RangeInt) FuzzAll(r rand.Rand) {
	p.Fuzz(r)
}

func (p *RangeInt) Parse(parser token.InternalParser, cur *token.ParserList) []token.ParserList {
	panic("TODO implement")
}

func (p *RangeInt) permutation(i int) {
	p.value = p.from + i
}

func (p *RangeInt) Permutation(i int) error {
	permutations := p.Permutations()

	if i < 1 || i > permutations {
		return &token.PermutationError{
			Type: token.PermutationErrorIndexOutOfBound,
		}
	}

	p.permutation(i - 1)

	return nil
}

func (p *RangeInt) Permutations() int {
	return p.to - p.from + 1
}

func (p *RangeInt) PermutationsAll() int {
	return p.Permutations()
}

func (p *RangeInt) String() string {
	return strconv.Itoa(p.value)
}
