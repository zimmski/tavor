package filter

import (
	"testing"

	. "github.com/zimmski/tavor/test/assert"

	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/lists"
	"github.com/zimmski/tavor/token/primitives"
)

type mockSeeFilter struct {
	Seen map[token.Token]struct{}
}

func (f *mockSeeFilter) Apply(tok token.Token) (token.Token, error) {
	f.Seen[tok] = struct{}{}

	return nil, nil
}

func TestStrategySeen(t *testing.T) {
	c1 := primitives.NewConstantInt(1)
	s1 := primitives.NewConstantString("a")
	p := primitives.NewPointer(s1)
	one := lists.NewOne(
		c1,
		p,
	)
	c2 := primitives.NewConstantInt(2)
	rep := lists.NewRepeat(c2, 2, 10)
	root := lists.NewAll(
		one,
		rep,
	)

	see := &mockSeeFilter{
		Seen: make(map[token.Token]struct{}),
	}
	filters := []Filter{
		see,
	}

	rootNew, err := ApplyFilters(filters, root)
	True(t, Exactly(t, root, rootNew))
	Nil(t, err)

	keyExists := func(key token.Token) bool {
		_, ok := see.Seen[key]

		return ok
	}

	Equal(t, 7, len(see.Seen))
	True(t, keyExists(root))
	True(t, keyExists(one))
	True(t, keyExists(c1))
	True(t, keyExists(p))
	True(t, keyExists(s1))
	True(t, keyExists(rep))
	True(t, keyExists(c2))
}

type mockReplaceFilter struct {
	stringSuffix string
}

func (f *mockReplaceFilter) Apply(tok token.Token) (token.Token, error) {
	if t, ok := tok.(*primitives.ConstantString); ok {
		return primitives.NewConstantString(t.String() + f.stringSuffix), nil
	}

	return nil, nil
}

func TestStrategyReplaces(t *testing.T) {
	filters := []Filter{
		&mockReplaceFilter{"b"},
	}

	// root replace
	{
		root := primitives.NewConstantString("a")

		rootNew, err := ApplyFilters(filters, root)
		Nil(t, err)
		Equal(t, "ab", rootNew.String())
	}
	// replace all over
	{
		c1 := primitives.NewConstantString("1")
		s1 := primitives.NewConstantString("a")
		one := lists.NewOne(
			c1,
		)
		p := primitives.NewPointer(s1)
		c2 := primitives.NewConstantString("2")
		rep := lists.NewRepeat(c2, 2, 10)
		root := lists.NewAll(
			one,
			p,
			rep,
		)

		rootNew, err := ApplyFilters(filters, root)
		Nil(t, err)
		Equal(t, "1bab2b2b", rootNew.String())
	}
	// double the replace: only the first filter should be applied
	{
		filters2 := []Filter{
			&mockReplaceFilter{"b"},
			&mockReplaceFilter{"c"},
		}
		root := primitives.NewConstantString("a")

		rootNew, err := ApplyFilters(filters2, root)
		Nil(t, err)
		Equal(t, 1, rootNew.Permutations())
		Equal(t, "ab", rootNew.String())
	}
}
