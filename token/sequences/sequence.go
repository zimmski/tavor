package sequences

import (
	"fmt"
	"strconv"

	"github.com/zimmski/tavor/rand"
	"github.com/zimmski/tavor/token"
)

type Sequence struct {
	start int
	step  int
	value int
}

func NewSequence(start, step int) *Sequence {
	return &Sequence{
		start: start,
		step:  step,
		value: start,
	}
}

func (s *Sequence) existing(r rand.Rand, except token.Token) int {
	n := s.value - s.start

	if n == 0 {
		return 0
	}

	n /= s.step

	if except == nil {
		return r.Intn(n)*s.step + s.start
	}

	ex, err := strconv.Atoi(except.String())
	if err != nil {
		panic(err) // TODO
	}

	if n == 1 && s.start == ex {
		panic(fmt.Sprintf("There is no sequence value to choose from")) // TODO
	}

	for {
		i := r.Intn(n)*s.step + s.start

		if i != ex {
			return i
		}
	}
}

func (s *Sequence) ExistingItem(except token.Token) *SequenceExistingItem {
	v := -1 // TODO there should be some kind of real nil value

	if s.value != s.start {
		v = s.start
	}

	return &SequenceExistingItem{
		sequence: s,
		value:    v,
		except:   except,
	}
}

func (s *Sequence) Item() *SequenceItem {
	return &SequenceItem{
		sequence: s,
		value:    s.Next(),
	}
}

func (s *Sequence) Next() int {
	c := s.value

	s.value += s.step

	return c
}

// ResetToken interface methods

// Reset resets the (internal) state of this token and its dependences
func (s *Sequence) Reset() {
	s.value = s.start
}

func (s *Sequence) ResetItem() *SequenceResetItem {
	return &SequenceResetItem{
		sequence: s,
	}
}

// Sequence is an unusable token

func (s *Sequence) Clone() token.Token  { panic("unusable token") }
func (s *Sequence) Fuzz(r rand.Rand)    { panic("unusable token") }
func (s *Sequence) FuzzAll(r rand.Rand) { panic("unusable token") }
func (s *Sequence) Parse(pars *token.InternalParser, cur int) (int, []error) {
	panic("unusable token")
}
func (s *Sequence) Permutation(i uint) error { panic("unusable token") }
func (s *Sequence) Permutations() uint       { panic("unusable token") }
func (s *Sequence) PermutationsAll() uint    { panic("unusable token") }
func (s *Sequence) String() string           { panic("unusable token") }

type SequenceItem struct {
	sequence *Sequence
	value    int
}

// Clone returns a copy of the token and all its children
func (s *SequenceItem) Clone() token.Token {
	return &SequenceItem{
		sequence: s.sequence,
		value:    s.value,
	}
}

func (s *SequenceItem) Fuzz(r rand.Rand) {
	s.permutation(0)
}

func (s *SequenceItem) FuzzAll(r rand.Rand) {
	s.Fuzz(r)
}

func (s *SequenceItem) Parse(pars *token.InternalParser, cur int) (int, []error) {
	panic("TODO implement")
}

func (s *SequenceItem) permutation(i uint) {
	s.value = s.sequence.Next()
}

func (s *SequenceItem) Permutation(i uint) error {
	permutations := s.Permutations()

	if i < 1 || i > permutations {
		return &token.PermutationError{
			Type: token.PermutationErrorIndexOutOfBound,
		}
	}

	s.permutation(i - 1)

	return nil
}

func (s *SequenceItem) Permutations() uint {
	return 1
}

func (s *SequenceItem) PermutationsAll() uint {
	return s.Permutations()
}

func (s *SequenceItem) String() string {
	return strconv.Itoa(s.value)
}

// ResetToken interface methods

// Reset resets the (internal) state of this token and its dependences
func (s *SequenceItem) Reset() {
	s.permutation(0)
}

type SequenceExistingItem struct {
	sequence *Sequence
	value    int
	except   token.Token
}

// Clone returns a copy of the token and all its children
func (s *SequenceExistingItem) Clone() token.Token {
	ex := s.except
	if ex != nil {
		ex = ex.Clone()
	}

	return &SequenceExistingItem{
		sequence: s.sequence,
		value:    s.value,
		except:   ex,
	}
}

func (s *SequenceExistingItem) Fuzz(r rand.Rand) {
	s.permutation(r)
}

func (s *SequenceExistingItem) FuzzAll(r rand.Rand) {
	s.Fuzz(r)
}

func (s *SequenceExistingItem) Parse(pars *token.InternalParser, cur int) (int, []error) {
	panic("TODO implement")
}

func (s *SequenceExistingItem) permutation(r rand.Rand) {
	s.value = s.sequence.existing(r, s.except)
}

func (s *SequenceExistingItem) Permutation(i uint) error {
	permutations := s.Permutations()

	if i < 1 || i > permutations {
		return &token.PermutationError{
			Type: token.PermutationErrorIndexOutOfBound,
		}
	}

	s.permutation(rand.NewIncrementRand(0))

	return nil
}

func (s *SequenceExistingItem) Permutations() uint {
	return 1
}

func (s *SequenceExistingItem) PermutationsAll() uint {
	return s.Permutations()
}

func (s *SequenceExistingItem) String() string {
	return strconv.Itoa(s.value)
}

// ForwardToken interface methods

func (s *SequenceExistingItem) Get() token.Token {
	return nil
}

func (s *SequenceExistingItem) InternalGet() token.Token {
	return s.except
}

func (s *SequenceExistingItem) InternalLogicalRemove(tok token.Token) token.Token {
	if s.except == tok {
		return nil
	}

	return s
}

func (s *SequenceExistingItem) InternalReplace(oldToken, newToken token.Token) {
	if s.except == oldToken {
		s.except = newToken
	}
}

// ResetToken interface methods

// Reset resets the (internal) state of this token and its dependences
func (s *SequenceExistingItem) Reset() {
	s.permutation(rand.NewIncrementRand(0))
}

// ScopeToken interface methods

func (s *SequenceExistingItem) SetScope(variableScope map[string]token.Token) {
	if s.except != nil {
		if tok, ok := s.except.(token.ScopeToken); ok {
			tok.SetScope(variableScope)
		}
	}
}

type SequenceResetItem struct {
	sequence *Sequence
}

// Clone returns a copy of the token and all its children
func (s *SequenceResetItem) Clone() token.Token {
	return &SequenceResetItem{
		sequence: s.sequence,
	}
}

func (s *SequenceResetItem) Fuzz(r rand.Rand) {
	s.permutation(0)
}

func (s *SequenceResetItem) FuzzAll(r rand.Rand) {
	s.Fuzz(r)
}

func (s *SequenceResetItem) Parse(pars *token.InternalParser, cur int) (int, []error) {
	panic("TODO implement")
}

func (s *SequenceResetItem) permutation(i uint) {
	s.sequence.Reset()
}

func (s *SequenceResetItem) Permutation(i uint) error {
	permutations := s.Permutations()

	if i < 1 || i > permutations {
		return &token.PermutationError{
			Type: token.PermutationErrorIndexOutOfBound,
		}
	}

	s.permutation(i - 1)

	return nil
}

func (s *SequenceResetItem) Permutations() uint {
	return 1
}

func (s *SequenceResetItem) PermutationsAll() uint {
	return s.Permutations()
}

func (s *SequenceResetItem) String() string {
	return ""
}

// ResetToken interface methods

// Reset resets the (internal) state of this token and its dependences
func (s *SequenceResetItem) Reset() {
	s.permutation(0)
}
