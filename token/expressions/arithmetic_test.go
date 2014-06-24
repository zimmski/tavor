package expressions

import (
	"testing"

	. "github.com/stretchr/testify/assert"

	"github.com/zimmski/tavor/test"
	"github.com/zimmski/tavor/token/lists"
	"github.com/zimmski/tavor/token/primitives"
)

func TestArithmeticExpressionTokensToBeTokens(t *testing.T) {
	var tok *lists.List

	Implements(t, tok, &AddArithmetic{})
	Implements(t, tok, &SubArithmetic{})
	Implements(t, tok, &MulArithmetic{})
	Implements(t, tok, &DivArithmetic{})
}

func TestAddArithmetic(t *testing.T) {
	a := primitives.NewRangeInt(1, 10)
	b := primitives.NewConstantInt(2)

	o := NewAddArithmetic(a, b)
	Equal(t, "3", o.String())

	i, err := o.Get(0)
	Nil(t, err)
	True(t, Exactly(t, a, i))
	i, err = o.Get(1)
	Nil(t, err)
	True(t, Exactly(t, b, i))
	i, err = o.Get(2)
	Equal(t, err.(*lists.ListError).Type, lists.ListErrorOutOfBound)
	Nil(t, i)

	r := test.NewRandTest(1)
	o.FuzzAll(r)
	Equal(t, "5", o.String())

	o2 := o.Clone()
	Equal(t, o.String(), o2.String())

	Equal(t, 1, o.Permutations())
	Equal(t, 10, o.PermutationsAll())
}

func TestSubArithmetic(t *testing.T) {
	a := primitives.NewRangeInt(1, 10)
	b := primitives.NewConstantInt(2)

	o := NewSubArithmetic(a, b)
	Equal(t, "-1", o.String())

	i, err := o.Get(0)
	Nil(t, err)
	True(t, Exactly(t, a, i))
	i, err = o.Get(1)
	Nil(t, err)
	True(t, Exactly(t, b, i))
	i, err = o.Get(2)
	Equal(t, err.(*lists.ListError).Type, lists.ListErrorOutOfBound)
	Nil(t, i)

	r := test.NewRandTest(1)
	o.FuzzAll(r)
	Equal(t, "1", o.String())

	o2 := o.Clone()
	Equal(t, o.String(), o2.String())

	Equal(t, 1, o.Permutations())
	Equal(t, 10, o.PermutationsAll())
}

func TestMulArithmetic(t *testing.T) {
	a := primitives.NewRangeInt(1, 10)
	b := primitives.NewConstantInt(2)

	o := NewMulArithmetic(a, b)
	Equal(t, "2", o.String())

	i, err := o.Get(0)
	Nil(t, err)
	True(t, Exactly(t, a, i))
	i, err = o.Get(1)
	Nil(t, err)
	True(t, Exactly(t, b, i))
	i, err = o.Get(2)
	Equal(t, err.(*lists.ListError).Type, lists.ListErrorOutOfBound)
	Nil(t, i)

	r := test.NewRandTest(1)
	o.FuzzAll(r)
	Equal(t, "6", o.String())

	o2 := o.Clone()
	Equal(t, o.String(), o2.String())

	Equal(t, 1, o.Permutations())
	Equal(t, 10, o.PermutationsAll())
}

func TestDivArithmetic(t *testing.T) {
	a := primitives.NewRangeInt(6, 10)
	b := primitives.NewConstantInt(2)

	o := NewDivArithmetic(a, b)
	Equal(t, "3", o.String())

	i, err := o.Get(0)
	Nil(t, err)
	True(t, Exactly(t, a, i))
	i, err = o.Get(1)
	Nil(t, err)
	True(t, Exactly(t, b, i))
	i, err = o.Get(2)
	Equal(t, err.(*lists.ListError).Type, lists.ListErrorOutOfBound)
	Nil(t, i)

	r := test.NewRandTest(2)
	o.FuzzAll(r)
	Equal(t, "4", o.String())

	o2 := o.Clone()
	Equal(t, o.String(), o2.String())

	Equal(t, 1, o.Permutations())
	Equal(t, 5, o.PermutationsAll())
}
