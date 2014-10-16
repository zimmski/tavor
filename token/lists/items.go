package lists

import (
	"strconv"

	"github.com/zimmski/tavor/rand"
	"github.com/zimmski/tavor/token"
)

type ListItem struct {
	index int
	list  token.List
}

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

func (l *ListItem) Fuzz(r rand.Rand) {
	// do nothing
}

func (l *ListItem) FuzzAll(r rand.Rand) {
	l.Fuzz(r)
}

func (l *ListItem) Parse(pars *token.InternalParser, cur int) (int, []error) {
	panic("TODO implement")
}

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

func (l *ListItem) Permutations() uint {
	return 1
}

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

type IndexItem struct {
	token token.IndexToken
}

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

func (l *IndexItem) Fuzz(r rand.Rand) {
	// do nothing
}

func (l *IndexItem) FuzzAll(r rand.Rand) {
	l.Fuzz(r)
}

func (l *IndexItem) Parse(pars *token.InternalParser, cur int) (int, []error) {
	panic("TODO implement")
}

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

func (l *IndexItem) Permutations() uint {
	return 1
}

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

type UniqueItem struct {
	original *UniqueItem
	list     token.List
	picked   map[int]struct{}

	index int
}

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

func (l *UniqueItem) Fuzz(r rand.Rand) {
	if l.index == -1 {
		l.pick(r)
	}
}

func (l *UniqueItem) FuzzAll(r rand.Rand) {
	l.Fuzz(r)
}

func (l *UniqueItem) Parse(pars *token.InternalParser, cur int) (int, []error) {
	panic("TODO implement")
}

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

func (l *UniqueItem) Permutations() uint {
	return 1
}

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
