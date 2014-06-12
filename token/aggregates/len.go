package aggregates

import (
	"strconv"

	"github.com/zimmski/tavor/rand"
	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/lists"
)

type Len struct {
	list lists.List
}

func NewLen(list lists.List) *Len {
	return &Len{
		list: list,
	}
}

func (a *Len) Clone() token.Token {
	return &Len{
		list: a.list,
	}
}

func (a *Len) FuzzAll(r rand.Rand) {
	// do nothing
}

func (a *Len) Permutations() int {
	return 1
}

func (a *Len) String() string {
	return strconv.Itoa(a.list.Len())
}
