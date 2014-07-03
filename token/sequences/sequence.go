package sequences

import (
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

func (s *Sequence) existing(r rand.Rand) int {
	n := s.value - s.start

	if n == 0 {
		return 0
	}

	n /= s.step

	return r.Intn(n)*s.step + s.start
}

func (s *Sequence) ExistingItem() *sequenceExistingItem {
	v := -1 // TODO there should be some kind of real nil value

	if s.value != s.start {
		v = s.start
	}

	return &sequenceExistingItem{
		sequence: s,
		value:    v,
	}
}

func (s *Sequence) Item() *sequenceItem {
	return &sequenceItem{
		sequence: s,
		value:    s.Next(),
	}
}

func (s *Sequence) Next() int {
	c := s.value

	s.value += s.step

	return c
}

func (s *Sequence) Reset() {
	s.value = s.start
}

func (s *Sequence) ResetItem() *sequenceResetItem {
	return &sequenceResetItem{
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
func (s *Sequence) Permutation(i int) error { panic("unusable token") }
func (s *Sequence) Permutations() int       { panic("unusable token") }
func (s *Sequence) PermutationsAll() int    { panic("unusable token") }
func (s *Sequence) String() string          { panic("unusable token") }

type sequenceItem struct {
	sequence *Sequence
	value    int
}

func (s *sequenceItem) Clone() token.Token {
	return &sequenceItem{
		sequence: s.sequence,
		value:    s.value,
	}
}

func (s *sequenceItem) Fuzz(r rand.Rand) {
	s.permutation(0)
}

func (s *sequenceItem) FuzzAll(r rand.Rand) {
	s.Fuzz(r)
}

func (s *sequenceItem) Parse(pars *token.InternalParser, cur int) (int, []error) {
	panic("TODO implement")
}

func (s *sequenceItem) permutation(i int) {
	s.value = s.sequence.Next()
}

func (s *sequenceItem) Permutation(i int) error {
	permutations := s.Permutations()

	if i < 1 || i > permutations {
		return &token.PermutationError{
			Type: token.PermutationErrorIndexOutOfBound,
		}
	}

	s.permutation(i - 1)

	return nil
}

func (s *sequenceItem) Permutations() int {
	return 1
}

func (s *sequenceItem) PermutationsAll() int {
	return s.Permutations()
}

func (s *sequenceItem) String() string {
	return strconv.Itoa(s.value)
}

// ResetToken interface methods

func (s *sequenceItem) Reset() {
	s.permutation(0)
}

type sequenceExistingItem struct {
	sequence *Sequence
	value    int
}

func (s *sequenceExistingItem) Clone() token.Token {
	return &sequenceExistingItem{
		sequence: s.sequence,
		value:    s.value,
	}
}

func (s *sequenceExistingItem) Fuzz(r rand.Rand) {
	s.permutation(r)
}

func (s *sequenceExistingItem) FuzzAll(r rand.Rand) {
	s.Fuzz(r)
}

func (s *sequenceExistingItem) Parse(pars *token.InternalParser, cur int) (int, []error) {
	panic("TODO implement")
}

func (s *sequenceExistingItem) permutation(r rand.Rand) {
	s.value = s.sequence.existing(r)
}

func (s *sequenceExistingItem) Permutation(i int) error {
	permutations := s.Permutations()

	if i < 1 || i > permutations {
		return &token.PermutationError{
			Type: token.PermutationErrorIndexOutOfBound,
		}
	}

	s.permutation(rand.NewConstantRand(0))

	return nil
}

func (s *sequenceExistingItem) Permutations() int {
	return 1
}

func (s *sequenceExistingItem) PermutationsAll() int {
	return s.Permutations()
}

func (s *sequenceExistingItem) String() string {
	return strconv.Itoa(s.value)
}

// ResetToken interface methods

func (s *sequenceExistingItem) Reset() {
	s.permutation(rand.NewConstantRand(0))
}

type sequenceResetItem struct {
	sequence *Sequence
}

func (s *sequenceResetItem) Clone() token.Token {
	return &sequenceResetItem{
		sequence: s.sequence,
	}
}

func (s *sequenceResetItem) Fuzz(r rand.Rand) {
	s.permutation(0)
}

func (s *sequenceResetItem) FuzzAll(r rand.Rand) {
	s.Fuzz(r)
}

func (s *sequenceResetItem) Parse(pars *token.InternalParser, cur int) (int, []error) {
	panic("TODO implement")
}

func (s *sequenceResetItem) permutation(i int) {
	s.sequence.Reset()
}

func (s *sequenceResetItem) Permutation(i int) error {
	permutations := s.Permutations()

	if i < 1 || i > permutations {
		return &token.PermutationError{
			Type: token.PermutationErrorIndexOutOfBound,
		}
	}

	s.permutation(i - 1)

	return nil
}

func (s *sequenceResetItem) Permutations() int {
	return 1
}

func (s *sequenceResetItem) PermutationsAll() int {
	return s.Permutations()
}

func (s *sequenceResetItem) String() string {
	return ""
}

// ResetToken interface methods

func (s *sequenceResetItem) Reset() {
	s.permutation(0)
}
