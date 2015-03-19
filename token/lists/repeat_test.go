package lists

import (
	"testing"

	"github.com/zimmski/go-leak"
	. "github.com/zimmski/tavor/test/assert"

	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/primitives"
)

func TestRepeatTokensToBeTokens(t *testing.T) {
	var tok *token.ListToken

	Implements(t, tok, &Repeat{})
}

func TestRepeat(t *testing.T) {
	a := primitives.NewConstantString("a")

	o := NewRepeat(a, 5, 10)
	Equal(t, "aaaaa", o.String())
	Equal(t, 5, o.Len())
	Equal(t, 6, o.Permutations())
	Equal(t, 6, o.PermutationsAll())

	i, err := o.Get(0)
	Nil(t, err)
	Equal(t, a, i)
	i, err = o.Get(6)
	Equal(t, err.(*ListError).Type, ListErrorOutOfBound)
	Nil(t, i)

	Nil(t, o.Permutation(1))
	Equal(t, "aaaaa", o.String())
	Nil(t, o.Permutation(2))
	Equal(t, "aaaaaa", o.String())
	Nil(t, o.Permutation(3))
	Equal(t, "aaaaaaa", o.String())

	Equal(t, o.Permutation(7).(*token.PermutationError).Type, token.PermutationErrorIndexOutOfBound)

	o = NewRepeat(primitives.NewRangeInt(1, 2), 0, 2)
	Equal(t, "", o.String())
	Equal(t, 0, o.Len())
	Equal(t, 3, o.Permutations())
	Equal(t, 7, o.PermutationsAll())

	o = NewRepeat(primitives.NewRangeInt(1, 2), 1, 2)
	Equal(t, "1", o.String())
	Equal(t, 1, o.Len())
	Equal(t, 2, o.Permutations())
	Equal(t, 6, o.PermutationsAll())

	o = NewRepeat(primitives.NewRangeInt(1, 2), 0, 3)
	Equal(t, "", o.String())
	Equal(t, 0, o.Len())
	Equal(t, 4, o.Permutations())
	Equal(t, 15, o.PermutationsAll())

	o = NewRepeat(primitives.NewRangeInt(1, 2), 1, 3)
	Equal(t, "1", o.String())
	Equal(t, 1, o.Len())
	Equal(t, 3, o.Permutations())
	Equal(t, 14, o.PermutationsAll())

	o = NewRepeat(primitives.NewRangeInt(1, 2), 3, 3)
	Equal(t, "111", o.String())
	Equal(t, 3, o.Len())
	Equal(t, 1, o.Permutations())
	Equal(t, 8, o.PermutationsAll())

	b := primitives.NewRangeInt(1, 3)
	o = NewRepeat(b, 2, 10)
	Equal(t, "11", o.String())
	Equal(t, 2, o.Len())
	Equal(t, 9, o.Permutations())
	Equal(t, 88569, o.PermutationsAll())

	Nil(t, o.Permutation(1))
	Equal(t, "11", o.String())
	Nil(t, o.Permutation(2))
	Equal(t, "111", o.String())
	Nil(t, o.Permutation(3))
	Equal(t, "1111", o.String())

	Equal(t, o.Permutation(10).(*token.PermutationError).Type, token.PermutationErrorIndexOutOfBound)

	o2 := o.Clone()
	Equal(t, o.String(), o2.String())
}

func TestRepeatReduces(t *testing.T) {
	a := primitives.NewConstantString("a")

	o := NewRepeat(a, 0, 1)
	Nil(t, o.Permutation(o.Permutations()))
	Equal(t, o.reduces(), []uint{1, 1})
	Equal(t, o.Reduces(), 2)

	// cannnot be reduces!
	o = NewRepeat(a, 1, 1)
	Nil(t, o.Permutation(o.Permutations()))
	Equal(t, o.reduces(), []uint{1})
	Equal(t, o.Reduces(), 0)

	o = NewRepeat(a, 0, 2)
	Nil(t, o.Permutation(o.Permutations()))
	Equal(t, o.reduces(), []uint{1, 2, 1})
	Equal(t, o.Reduces(), 4)

	o = NewRepeat(a, 1, 2)
	Nil(t, o.Permutation(o.Permutations()))
	Equal(t, o.reduces(), []uint{2, 1})
	Equal(t, o.Reduces(), 3)

	o = NewRepeat(a, 0, 3)
	Nil(t, o.Permutation(o.Permutations()))
	Equal(t, o.reduces(), []uint{1, 3, 3, 1})
	Equal(t, o.Reduces(), 8)

	o = NewRepeat(a, 1, 3)
	Nil(t, o.Permutation(o.Permutations()))
	Equal(t, o.reduces(), []uint{3, 3, 1})
	Equal(t, o.Reduces(), 7)

	o = NewRepeat(a, 2, 3)
	Nil(t, o.Permutation(o.Permutations()))
	Equal(t, o.reduces(), []uint{3, 1})
	Equal(t, o.Reduces(), 4)

	// cannnot be reduces!
	o = NewRepeat(a, 3, 3)
	Nil(t, o.Permutation(o.Permutations()))
	Equal(t, o.reduces(), []uint{1})
	Equal(t, o.Reduces(), 0)
}

func TestRepeatCombinations(t *testing.T) {
	type tt struct {
		n        int
		k        int
		expected [][]int
	}
	tests := []tt{
		tt{
			n: 1, k: 0,
			expected: [][]int{
				[]int{},
			},
		},
		tt{
			n: 1, k: 1,
			expected: [][]int{
				[]int{0},
			},
		},
		tt{
			n: 2, k: 0,
			expected: [][]int{
				[]int{},
			},
		},
		tt{
			n: 2, k: 1,
			expected: [][]int{
				[]int{0},
				[]int{1},
			},
		},
		tt{
			n: 2, k: 2,
			expected: [][]int{
				[]int{0, 1},
			},
		},
		tt{
			n: 3, k: 0,
			expected: [][]int{
				[]int{},
			},
		},
		tt{
			n: 3, k: 1,
			expected: [][]int{
				[]int{0},
				[]int{1},
				[]int{2},
			},
		},
		tt{
			n: 3, k: 2,
			expected: [][]int{
				[]int{0, 1},
				[]int{0, 2},
				[]int{1, 2},
			},
		},
		tt{
			n: 3, k: 3,
			expected: [][]int{
				[]int{0, 1, 2},
			},
		},
		tt{
			n: 4, k: 0,
			expected: [][]int{
				[]int{},
			},
		},
		tt{
			n: 4, k: 1,
			expected: [][]int{
				[]int{0},
				[]int{1},
				[]int{2},
				[]int{3},
			},
		},
		tt{
			n: 4, k: 2,
			expected: [][]int{
				[]int{0, 1},
				[]int{0, 2},
				[]int{0, 3},
				[]int{1, 2},
				[]int{1, 3},
				[]int{2, 3},
			},
		},
		tt{
			n: 4, k: 3,
			expected: [][]int{
				[]int{0, 1, 2},
				[]int{0, 1, 3},
				[]int{0, 2, 3},
				[]int{1, 2, 3},
			},
		},
		tt{
			n: 4, k: 4,
			expected: [][]int{
				[]int{0, 1, 2, 3},
			},
		},
	}

	m := leak.MarkGoRoutines()

	for _, test := range tests {
		var actual [][]int

		ch, _ := combinations(test.n, test.k)
		for c := range ch {
			actual = append(actual, c)
		}

		Equal(t, test.expected, actual)
	}

	Equal(t, 0, m.Release(), "check for goroutine leaks")
}

func TestRepeatReduce(t *testing.T) {
	a := primitives.NewConstantString("a")

	type tt struct {
		from     int64
		to       int64
		expected []string
	}

	tests := []tt{
		tt{
			from: 0, to: 1,
			expected: []string{
				"",
				"0",
			},
		},
		tt{
			from: 1, to: 1,
			expected: nil,
		},
		tt{
			from: 0, to: 2,
			expected: []string{
				"",
				"0",
				"1",
				"01",
			},
		},
		tt{
			from: 1, to: 2,
			expected: []string{
				"0",
				"1",
				"01",
			},
		},
		tt{
			from: 2, to: 2,
			expected: nil,
		},
		tt{
			from: 0, to: 3,
			expected: []string{
				"",
				"0",
				"1",
				"2",
				"01",
				"02",
				"12",
				"012",
			},
		},
		tt{
			from: 1, to: 3,
			expected: []string{
				"0",
				"1",
				"2",
				"01",
				"02",
				"12",
				"012",
			},
		},
		tt{
			from: 2, to: 3,
			expected: []string{
				"01",
				"02",
				"12",
				"012",
			},
		},
		tt{
			from: 3, to: 3,
			expected: nil,
		},
		tt{
			from: 3, to: 5,
			expected: []string{
				"012",
				"013",
				"014",
				"023",
				"024",
				"034",
				"123",
				"124",
				"134",
				"234",
				"0123",
				"0124",
				"0134",
				"0234",
				"1234",
				"01234",
			},
		},
	}

	m := leak.MarkGoRoutines()

	for _, test := range tests {
		o := NewRepeat(a, test.from, test.to)
		o.value = make([]token.Token, test.to)
		for i := 0; i < len(o.value); i++ {
			o.value[i] = primitives.NewConstantInt(i)
		}

		reduces := uint(len(test.expected))

		Equal(t, reduces, o.Reduces())

		err := o.Reduce(0)
		NotNil(t, err)
		err = o.Reduce(reduces + 1)
		NotNil(t, err)

		if reduces > 0 {
			var actual []string

			for i := uint(1); i <= o.Reduces(); i++ {
				err := o.Reduce(i)
				Nil(t, err)

				actual = append(actual, o.String())
			}

			Equal(t, test.expected, actual)
		}
	}

	Equal(t, 0, m.Release(), "check for goroutine leaks")
}
