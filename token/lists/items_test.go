package lists

import (
	"testing"

	. "github.com/zimmski/tavor/test/assert"

	"github.com/zimmski/tavor/token/primitives"
)

func TestUniqueItem(t *testing.T) {
	list := NewAll(
		primitives.NewConstantInt(1),
		primitives.NewConstantInt(2),
		primitives.NewConstantInt(3),
	)

	a := NewUniqueItem(list)
	Equal(t, 3, a.Permutations())
	Equal(t, 3, a.PermutationsAll())
	Equal(t, "1", a.String())
	Equal(t, 3, a.Permutations())
	Equal(t, 3, a.PermutationsAll())

	Nil(t, a.Permutation(3))
	Equal(t, "3", a.String())
	Equal(t, 3, a.Permutations())
	Equal(t, 3, a.PermutationsAll())

	Nil(t, a.Permutation(2))
	Equal(t, "2", a.String())
	Equal(t, 3, a.Permutations())
	Equal(t, 3, a.PermutationsAll())

	Nil(t, a.Permutation(1))
	Equal(t, "1", a.String())
	Equal(t, 3, a.Permutations())
	Equal(t, 3, a.PermutationsAll())

	b := a.Clone().(*UniqueItem)
	Equal(t, 2, b.Permutations())
	Equal(t, 2, b.PermutationsAll())
	Equal(t, "2", b.String())
	Equal(t, 2, b.Permutations())
	Equal(t, 2, b.PermutationsAll())

	Nil(t, b.Permutation(1))
	Equal(t, "2", b.String())
	Equal(t, 2, b.Permutations())
	Equal(t, 2, b.PermutationsAll())

	Nil(t, b.Permutation(2))
	Equal(t, "3", b.String())
	Equal(t, 2, b.Permutations())
	Equal(t, 2, b.PermutationsAll())
}
