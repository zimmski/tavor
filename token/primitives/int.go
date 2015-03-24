package primitives

import (
	"fmt"
	"math"
	"strconv"

	"github.com/zimmski/tavor/log"
	"github.com/zimmski/tavor/token"
)

// ConstantInt implements an integer token which holds a constant integer
type ConstantInt struct {
	value int
}

// NewConstantInt returns a new instance of a ConstantInt token
func NewConstantInt(value int) *ConstantInt {
	return &ConstantInt{
		value: value,
	}
}

// SetValue sets the value of the token
func (p *ConstantInt) SetValue(v int) {
	p.value = v
}

// Value returns the value of the token
func (p *ConstantInt) Value() int {
	return p.value
}

// Token interface methods

// Clone returns a copy of the token and all its children
func (p *ConstantInt) Clone() token.Token {
	return &ConstantInt{
		value: p.value,
	}
}

// Parse tries to parse the token beginning from the current position in the parser data.
// If the parsing is successful the error argument is nil and the next current position after the token is returned.
func (p *ConstantInt) Parse(pars *token.InternalParser, cur int) (int, []error) {
	v := strconv.Itoa(p.value)
	vLen := len(v)

	nextIndex := vLen + cur

	if nextIndex > pars.DataLen {
		return cur, []error{&token.ParserError{
			Message: fmt.Sprintf("expected %q but got early EOF", v),
			Type:    token.ParseErrorUnexpectedEOF,

			Position: pars.GetPosition(cur),
		}}
	}

	if got := pars.Data[cur:nextIndex]; v != got {
		return cur, []error{&token.ParserError{
			Message: fmt.Sprintf("expected %q but got %q", v, got),
			Type:    token.ParseErrorUnexpectedData,

			Position: pars.GetPosition(cur),
		}}
	}

	log.Debugf("Parsed %q", p.value)

	return nextIndex, nil
}

// Permutation sets a specific permutation for this token
func (p *ConstantInt) Permutation(i uint) error {
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
func (p *ConstantInt) Permutations() uint {
	return 1
}

// PermutationsAll returns the number of all possible permutations for this token including its children
func (p *ConstantInt) PermutationsAll() uint {
	return p.Permutations()
}

func (p *ConstantInt) String() string {
	return strconv.Itoa(p.value)
}

// RangeInt implements an integer token holding a range of integers
// Every permutation generates a new value within the defined range and step. For example the range 1 to 10 with step 2 can hold the integers 1, 3, 5, 7 and 9.
type RangeInt struct {
	from int
	to   int
	step int

	value int
}

// NewRangeInt returns a new instance of a RangeInt token with the given range and step value of 1
func NewRangeInt(from, to int) *RangeInt {
	if from > to {
		panic("TODO implement that From can be bigger than To")
	}

	return &RangeInt{
		from: from,
		to:   to,
		step: 1,

		value: from,
	}
}

// NewRangeIntWithStep returns a new instance of a RangeInt token with the given range and step value
func NewRangeIntWithStep(from, to, step int) *RangeInt {
	if from > to {
		panic("TODO implement that From can be bigger than To")
	}
	if step < 1 {
		panic("TODO implement 0 and negative step")
	}

	return &RangeInt{
		from: from,
		to:   to,
		step: step,

		value: from,
	}
}

func init() {
	token.RegisterTyped("Int", func(argParser token.ArgumentsTypedParser) (token.Token, error) {
		from := argParser.GetInt("from", 0)
		to := argParser.GetInt("to", math.MaxInt32)
		step := argParser.GetInt("step", 1)

		if err := argParser.Err(); err != nil {
			return nil, err
		}

		return NewRangeIntWithStep(from, to, step), nil
	})
}

// From returns the from value of the range
func (p *RangeInt) From() int {
	return p.from
}

// To returns the to value of the range
func (p *RangeInt) To() int {
	return p.to
}

// Step returns the step value
func (p *RangeInt) Step() int {
	return p.step
}

// Token interface methods

// Clone returns a copy of the token and all its children
func (p *RangeInt) Clone() token.Token {
	return &RangeInt{
		from: p.from,
		to:   p.to,
		step: p.step,

		value: p.value,
	}
}

// Parse tries to parse the token beginning from the current position in the parser data.
// If the parsing is successful the error argument is nil and the next current position after the token is returned.
func (p *RangeInt) Parse(pars *token.InternalParser, cur int) (int, []error) {
	if cur == pars.DataLen {
		return cur, []error{&token.ParserError{
			Message: fmt.Sprintf("expected integer in range %d-%d with step %d but got early EOF", p.from, p.to, p.step),
			Type:    token.ParseErrorUnexpectedEOF,

			Position: pars.GetPosition(cur),
		}}
	}

	i := cur
	v := ""

	for {
		c := pars.Data[i]

		if c < '0' || c > '9' {
			break
		}

		v += string(c)

		if ci, _ := strconv.Atoi(v); ci > p.to {
			v = v[:len(v)-1] // remove last digit

			break
		}

		i++

		if i == pars.DataLen {
			break
		}
	}

	i--

	ci, _ := strconv.Atoi(v)

	if v == "" || (ci < p.from || ci > p.to) || ci%p.step != 0 {
		// is the first character already invalid
		if i < cur {
			i = cur
		}

		return cur, []error{&token.ParserError{
			Message: fmt.Sprintf("expected integer in range %d-%d with step %d but got %q", p.from, p.to, p.step, pars.Data[cur:i]),
			Type:    token.ParseErrorUnexpectedData,

			Position: pars.GetPosition(cur),
		}}
	}

	p.value = ci

	log.Debugf("Parsed %q", p.value)

	return i + 1, nil
}

func (p *RangeInt) permutation(i uint) {
	p.value = p.from + (int(i) * p.step)
}

// Permutation sets a specific permutation for this token
func (p *RangeInt) Permutation(i uint) error {
	permutations := p.Permutations()

	if i < 1 || i > permutations {
		return &token.PermutationError{
			Type: token.PermutationErrorIndexOutOfBound,
		}
	}

	p.permutation(i - 1)

	return nil
}

// Permutations returns the number of permutations for this token
func (p *RangeInt) Permutations() uint {
	// TODO FIXME this
	perms := (p.to-p.from)/p.step + 1

	if perms < 0 {
		return math.MaxUint32
	}

	return uint(perms)
}

// PermutationsAll returns the number of all possible permutations for this token including its children
func (p *RangeInt) PermutationsAll() uint {
	return p.Permutations()
}

func (p *RangeInt) String() string {
	return strconv.Itoa(p.value)
}
