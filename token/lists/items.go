package lists

import (
	"strconv"

	"github.com/zimmski/tavor/rand"
	"github.com/zimmski/tavor/token"
)

// ListItem implements a list item token which references a List token and holds one index of the list to reference a list item
type ListItem struct {
	index int
	list  token.List
}

// NewListItem returns a new instance of a ListItem token referencing the given list and the given index
func NewListItem(index int, list token.List) *ListItem {
	return &ListItem{
		index: index,
		list:  list,
	}
}

// Token interface methods

// Clone returns a copy of the token and all its children
func (l *ListItem) Clone() token.Token {
	return &ListItem{
		index: l.index,
		list:  l.list,
	}
}

// Fuzz fuzzes this token using the random generator by choosing one of the possible permutations for this token
func (l *ListItem) Fuzz(r rand.Rand) {
	// do nothing
}

// FuzzAll calls Fuzz for this token and then FuzzAll for all children of this token
func (l *ListItem) FuzzAll(r rand.Rand) {
	l.Fuzz(r)
}

// Parse tries to parse the token beginning from the current position in the parser data.
// If the parsing is successful the error argument is nil and the next current position after the token is returned.
func (l *ListItem) Parse(pars *token.InternalParser, cur int) (int, []error) {
	panic("TODO implement")
}

// Permutation sets a specific permutation for this token
func (l *ListItem) Permutation(i uint) error {
	permutations := l.Permutations()

	if i < 1 || i > permutations {
		return &token.PermutationError{
			Type: token.PermutationErrorIndexOutOfBound,
		}
	}

	// do nothing

	return nil
}

// Permutations returns the number of permutations for this token
func (l *ListItem) Permutations() uint {
	return 1
}

// PermutationsAll returns the number of all possible permutations for this token including its children
func (l *ListItem) PermutationsAll() uint {
	return l.Permutations()
}

func (l *ListItem) String() string {
	tok, err := l.list.Get(l.index)
	if err != nil {
		panic(err) // TODO
	}

	return tok.String()
}

// IndexToken interface methods

// Index returns the index of this token in its parent token
func (l *ListItem) Index() int {
	return l.index
}

// IndexItem implements a list item which references an Index token to represent the index itself of this token
type IndexItem struct {
	token token.IndexToken
}

// NewIndexItem returns a new instance of a IndexItem token referencing the given Index token
func NewIndexItem(tok token.IndexToken) *IndexItem {
	return &IndexItem{
		token: tok.Clone().(token.IndexToken),
	}
}

// Token interface methods

// Clone returns a copy of the token and all its children
func (l *IndexItem) Clone() token.Token {
	return &IndexItem{
		token: l.token.Clone().(token.IndexToken),
	}
}

// Fuzz fuzzes this token using the random generator by choosing one of the possible permutations for this token
func (l *IndexItem) Fuzz(r rand.Rand) {
	// do nothing
}

// FuzzAll calls Fuzz for this token and then FuzzAll for all children of this token
func (l *IndexItem) FuzzAll(r rand.Rand) {
	l.Fuzz(r)
}

// Parse tries to parse the token beginning from the current position in the parser data.
// If the parsing is successful the error argument is nil and the next current position after the token is returned.
func (l *IndexItem) Parse(pars *token.InternalParser, cur int) (int, []error) {
	panic("TODO implement")
}

// Permutation sets a specific permutation for this token
func (l *IndexItem) Permutation(i uint) error {
	permutations := l.Permutations()

	if i < 1 || i > permutations {
		return &token.PermutationError{
			Type: token.PermutationErrorIndexOutOfBound,
		}
	}

	// do nothing

	return nil
}

// Permutations returns the number of permutations for this token
func (l *IndexItem) Permutations() uint {
	return 1
}

// PermutationsAll returns the number of all possible permutations for this token including its children
func (l *IndexItem) PermutationsAll() uint {
	return l.Permutations()
}

func (l *IndexItem) String() string {
	return strconv.Itoa(l.token.Index())
}

// ScopeToken interface methods

// SetScope sets the scope of the token
func (l *IndexItem) SetScope(variableScope map[string]token.Token) {
	if tok, ok := l.token.(token.ScopeToken); ok {
		tok.SetScope(variableScope)
	}
}

// UniqueItem implements a list item token which holds an distinct list item of a referenced List token
type UniqueItem struct {
	original *UniqueItem
	list     token.List
	picked   map[int]struct{}

	index int
}

// NewUniqueItem returns a new instance of a UniqueItem token referencing the given List token
func NewUniqueItem(list token.List) *UniqueItem {
	l := &UniqueItem{
		list:   list,
		picked: make(map[int]struct{}),

		index: -1,
	}

	l.original = l

	return l
}

func (l *UniqueItem) pick(r rand.Rand) {
	nList := l.original.list.Len()
	nPicked := len(l.original.picked)

	if nPicked >= nList {
		panic("already picked everything!") // TODO
	}

	// TODO make this WAYYYYYYYYY more effiecent
	for {
		c := r.Intn(nList)

		if _, ok := l.original.picked[c]; !ok {
			l.index = c
			l.original.picked[c] = struct{}{}

			break
		}
	}
}

// Token interface methods

// Clone returns a copy of the token and all its children
func (l *UniqueItem) Clone() token.Token {
	n := &UniqueItem{
		original: l.original,
		list:     nil,
		picked:   nil,

		index: -1,
	}

	return n
}

// Fuzz fuzzes this token using the random generator by choosing one of the possible permutations for this token
func (l *UniqueItem) Fuzz(r rand.Rand) {
	if l.index == -1 {
		l.pick(r)
	}
}

// FuzzAll calls Fuzz for this token and then FuzzAll for all children of this token
func (l *UniqueItem) FuzzAll(r rand.Rand) {
	l.Fuzz(r)
}

// Parse tries to parse the token beginning from the current position in the parser data.
// If the parsing is successful the error argument is nil and the next current position after the token is returned.
func (l *UniqueItem) Parse(pars *token.InternalParser, cur int) (int, []error) {
	panic("TODO implement")
}

// Permutation sets a specific permutation for this token
func (l *UniqueItem) Permutation(i uint) error {
	permutations := l.Permutations()

	if i < 1 || i > permutations {
		return &token.PermutationError{
			Type: token.PermutationErrorIndexOutOfBound,
		}
	}

	// do nothing

	return nil
}

// Permutations returns the number of permutations for this token
func (l *UniqueItem) Permutations() uint {
	return 1
}

// PermutationsAll returns the number of all possible permutations for this token including its children
func (l *UniqueItem) PermutationsAll() uint {
	return l.Permutations()
}

func (l *UniqueItem) String() string {
	i := l.Index()

	tok, err := l.original.list.Get(i)
	if err != nil {
		panic(err) // TODO
	}

	return tok.String()
}

// IndexToken interface methods

// Index returns the index of this token in its parent token
func (l *UniqueItem) Index() int {
	if l.index == -1 {
		l.pick(rand.NewIncrementRand(0))
	}

	return l.index
}

// ResetToken interface methods

// Reset resets the (internal) state of this token and its dependences
func (l *UniqueItem) Reset() {
	if l.index != -1 {
		delete(l.original.picked, l.index)

		l.index = -1
	}
}
