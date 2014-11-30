package lists

import (
	"testing"

	. "github.com/zimmski/tavor/test/assert"

	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/primitives"
)

func TestOnceTokensToBeTokens(t *testing.T) {
	var tok *token.ListToken

	Implements(t, tok, &Once{})
}

func TestOnce(t *testing.T) {
	a := primitives.NewConstantInt(10)
	b := primitives.NewConstantString("abc")
	c := primitives.NewConstantString("def")

	o := NewOnce(a, b, c)
	Equal(t, "10abcdef", o.String())
	Equal(t, 3, o.Len())
	Equal(t, 6, o.Permutations())
	Equal(t, 6, o.PermutationsAll())

	i, err := o.Get(0)
	Nil(t, err)
	Equal(t, a, i)
	i, err = o.Get(1)
	Nil(t, err)
	Equal(t, b, i)
	i, err = o.Get(2)
	Nil(t, err)
	Equal(t, c, i)
	i, err = o.Get(3)
	Equal(t, err.(*ListError).Type, ListErrorOutOfBound)
	Nil(t, i)

	for i, s := range []string{
		"10abcdef",
		"10defabc",
		"abc10def",
		"abcdef10",
		"def10abc",
		"defabc10",
	} {
		Nil(t, o.Permutation(uint(i+1)))
		Equal(t, s, o.String())
	}

	d := primitives.NewRangeInt(1, 2)
	o = NewOnce(a, b, d)
	Equal(t, "10abc1", o.String())
	Equal(t, 3, o.Len())
	Equal(t, 6, o.Permutations())
	Equal(t, 12, o.PermutationsAll())

	Nil(t, o.Permutation(2))
	Equal(t, "101abc", o.String())
	Equal(t, 3, o.Len())

	o2 := o.Clone()
	Equal(t, o.String(), o2.String())
}
