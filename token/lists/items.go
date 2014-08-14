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

func (l *ListItem) Permutation(i int) error {
	permutations := l.Permutations()

	if i < 1 || i > permutations {
		return &token.PermutationError{
			Type: token.PermutationErrorIndexOutOfBound,
		}
	}

	// do nothing

	return nil
}

func (l *ListItem) Permutations() int {
	return 1
}

func (l *ListItem) PermutationsAll() int {
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

func (l *ListItem) Index() int {
	return l.index
}

type IndexItem struct {
	token token.IndexToken
}

func NewIndexItem(token token.IndexToken) *IndexItem {
	return &IndexItem{
		token: token,
	}
}

// Token interface methods

func (l *IndexItem) Clone() token.Token {
	return &IndexItem{
		token: l.token,
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

func (l *IndexItem) Permutation(i int) error {
	permutations := l.Permutations()

	if i < 1 || i > permutations {
		return &token.PermutationError{
			Type: token.PermutationErrorIndexOutOfBound,
		}
	}

	// do nothing

	return nil
}

func (l *IndexItem) Permutations() int {
	return 1
}

func (l *IndexItem) PermutationsAll() int {
	return l.Permutations()
}

func (l *IndexItem) String() string {
	return strconv.Itoa(l.token.Index())
}

// ScopeToken interface methods

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

	return l
}

func (l *UniqueItem) pick(r rand.Rand) {
	nList := l.list.Len()
	nPicked := len(l.picked)

	if nPicked >= nList {
		panic("already picked everything!") // TODO
	}

	// TODO make this WAYYYYYYYYY more effiecent
	for {
		c := r.Intn(nList)

		if _, ok := l.picked[c]; !ok {
			l.index = c
			l.picked[c] = struct{}{}

			break
		}
	}
}

// Token interface methods

func (l *UniqueItem) Clone() token.Token {
	n := &UniqueItem{
		list:   l.list,
		picked: l.picked,

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

func (l *UniqueItem) Permutation(i int) error {
	permutations := l.Permutations()

	if i < 1 || i > permutations {
		return &token.PermutationError{
			Type: token.PermutationErrorIndexOutOfBound,
		}
	}

	// do nothing

	return nil
}

func (l *UniqueItem) Permutations() int {
	return 1
}

func (l *UniqueItem) PermutationsAll() int {
	return l.Permutations()
}

func (l *UniqueItem) String() string {
	i := l.Index()

	tok, err := l.list.Get(i)
	if err != nil {
		panic(err) // TODO
	}

	return tok.String()
}

// IndexToken interface methods

func (l *UniqueItem) Index() int {
	if l.index == -1 {
		l.pick(rand.NewConstantRand(0))
	}

	return l.index
}
