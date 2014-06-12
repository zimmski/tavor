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
func (s *Sequence) Permutations() int   { panic("unusable token") }
func (s *Sequence) String() string      { panic("unusable token") }

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
	s.value = s.sequence.Next()
}

func (s *sequenceItem) FuzzAll(r rand.Rand) {
	s.Fuzz(r)
}

func (s *sequenceItem) Permutations() int {
	return 1
}

func (s *sequenceItem) String() string {
	return strconv.Itoa(s.value)
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
	s.value = s.sequence.existing(r)
}

func (s *sequenceExistingItem) FuzzAll(r rand.Rand) {
	s.Fuzz(r)
}

func (s *sequenceExistingItem) Permutations() int {
	return 1
}

func (s *sequenceExistingItem) String() string {
	return strconv.Itoa(s.value)
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
	s.sequence.Reset()
}

func (s *sequenceResetItem) FuzzAll(r rand.Rand) {
	s.Fuzz(r)
}

func (s *sequenceResetItem) Permutations() int {
	return 1
}

func (s *sequenceResetItem) String() string {
	return ""
}
