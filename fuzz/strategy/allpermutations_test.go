package strategy

import (
	"strings"
	"testing"

	. "github.com/zimmski/tavor/test/assert"

	"github.com/zimmski/tavor/parser"
	"github.com/zimmski/tavor/test"
	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/constraints"
	"github.com/zimmski/tavor/token/lists"
	"github.com/zimmski/tavor/token/primitives"
	"github.com/zimmski/tavor/token/sequences"
)

func TestAllPermutationsStrategygetLevel(t *testing.T) {
	o := &allPermutations{}

	var nilChildren []allPermutationsLevel

	{
		a := primitives.NewConstantInt(1)
		b1 := primitives.NewConstantInt(2)
		b := constraints.NewOptional(b1)
		c1 := primitives.NewConstantInt(3)
		c := primitives.NewPointer(c1)
		d := lists.NewConcatenation(a, b, c)

		tree := o.getTree(d, false)

		Equal(t, tree, []allPermutationsLevel{
			allPermutationsLevel{
				token:       d,
				permutation: 0,

				children: []allPermutationsLevel{
					allPermutationsLevel{
						token:       a,
						permutation: 0,

						children: nilChildren,
					},
					allPermutationsLevel{
						token:       b,
						permutation: 0,

						children: nilChildren,
					},
					allPermutationsLevel{
						token:       c,
						permutation: 0,

						children: []allPermutationsLevel{
							allPermutationsLevel{
								token:       c1,
								permutation: 0,

								children: nilChildren,
							},
						},
					},
				},
			},
		})

		tree = o.getTree(d, true)

		Equal(t, tree, []allPermutationsLevel{
			allPermutationsLevel{
				token:       a,
				permutation: 0,

				children: nilChildren,
			},
			allPermutationsLevel{
				token:       b,
				permutation: 0,

				children: nilChildren,
			},
			allPermutationsLevel{
				token:       c,
				permutation: 0,

				children: []allPermutationsLevel{
					allPermutationsLevel{
						token:       c1,
						permutation: 0,

						children: nilChildren,
					},
				},
			},
		})
	}
}

func TestAllPermutationsStrategy(t *testing.T) {
	r := test.NewRandTest(1)

	{
		a := primitives.NewConstantInt(1)

		ch, err := NewAllPermutations(a, r)
		Nil(t, err)

		_, ok := <-ch
		True(t, ok)
		Equal(t, "1", a.String())
		ch <- struct{}{}

		_, ok = <-ch
		False(t, ok)
	}
	{
		a := constraints.NewOptional(primitives.NewConstantInt(1))

		validateTokenAllPermutations(
			t,
			a,
			[]string{
				"",
				"1",
			},
		)
	}
	{
		a := lists.NewOne(
			primitives.NewConstantInt(1),
			primitives.NewConstantInt(2),
			primitives.NewConstantInt(3),
		)

		validateTokenAllPermutations(
			t,
			a,
			[]string{
				"1",
				"2",
				"3",
			},
		)
	}
	{
		a := constraints.NewOptional(primitives.NewConstantInt(1))
		b := constraints.NewOptional(primitives.NewConstantInt(2))
		c := constraints.NewOptional(primitives.NewConstantInt(3))
		abc := lists.NewConcatenation(a, b, c)

		validateTokenAllPermutations(
			t,
			abc,
			[]string{
				"",
				"1",
				"2",
				"12",
				"3",
				"13",
				"23",
				"123",
			},
		)
	}
	{
		abc := lists.NewConcatenation(
			constraints.NewOptional(lists.NewConcatenation(
				lists.NewOne(
					primitives.NewConstantInt(1),
					primitives.NewConstantInt(2),
				),
				primitives.NewConstantInt(3),
			)),
			primitives.NewConstantInt(4),
		)

		validateTokenAllPermutations(
			t,
			abc,
			[]string{
				"4",
				"134",
				"234",
			},
		)
	}
	{
		a := constraints.NewOptional(primitives.NewConstantInt(1))
		b := primitives.NewConstantInt(2)
		c := constraints.NewOptional(primitives.NewConstantInt(3))
		d := lists.NewConcatenation(a, b, c)

		ch, err := NewAllPermutations(d, r)
		Nil(t, err)

		_, ok := <-ch
		True(t, ok)
		Equal(t, "2", d.String())
		ch <- struct{}{}

		_, ok = <-ch
		True(t, ok)
		Equal(t, "12", d.String())
		ch <- struct{}{}

		_, ok = <-ch
		True(t, ok)
		Equal(t, "23", d.String())
		ch <- struct{}{}

		_, ok = <-ch
		True(t, ok)
		Equal(t, "123", d.String())
		ch <- struct{}{}

		_, ok = <-ch
		False(t, ok)

		// rerun
		ch, err = NewAllPermutations(d, r)
		Nil(t, err)

		_, ok = <-ch
		True(t, ok)
		Equal(t, "2", d.String())

		close(ch)

		// run with range
		var got []string

		ch, err = NewAllPermutations(d, r)
		Nil(t, err)
		for i := range ch {
			got = append(got, d.String())

			ch <- i
		}

		Equal(t, got, []string{
			"2",
			"12",
			"23",
			"123",
		})
	}
	{
		a1 := constraints.NewOptional(primitives.NewConstantInt(1))
		a2 := constraints.NewOptional(primitives.NewConstantInt(11))
		a := constraints.NewOptional(lists.NewConcatenation(a1, a2, primitives.NewConstantString("a")))
		b := constraints.NewOptional(primitives.NewConstantString("b"))
		c := lists.NewConcatenation(a, b, primitives.NewConstantString("c"))
		d := constraints.NewOptional(c)

		validateTokenAllPermutations(
			t,
			d,
			[]string{
				"",
				"c",
				"ac",
				"1ac",
				"11ac",
				"111ac",
				"bc",
				"abc",
				"1abc",
				"11abc",
				"111abc",
			},
		)
	}
	{
		a := lists.NewConcatenation(
			constraints.NewOptional(primitives.NewConstantInt(1)),
			constraints.NewOptional(primitives.NewConstantInt(2)),
		)
		b := lists.NewRepeat(a, 0, 2)

		validateTokenAllPermutations(
			t,
			b,
			[]string{
				"", // 0x
				"", // 1x
				"1",
				"2",
				"12",
				"", // 2x
				"1",
				"2",
				"12",
				"1",
				"11",
				"21",
				"121",
				"2",
				"12",
				"22",
				"122",
				"12",
				"112",
				"212",
				"1212",
			},
		)
	}
	{
		s := sequences.NewSequence(10, 2)

		Equal(t, 10, s.Next())
		Equal(t, 12, s.Next())

		a := lists.NewConcatenation(
			constraints.NewOptional(primitives.NewConstantString("a")),
			constraints.NewOptional(primitives.NewConstantString("b")),
			s.ResetItem(),
			s.Item(),
			s.ExistingItem(nil),
		)
		b := lists.NewRepeat(a, 0, 1)

		validateTokenAllPermutations(
			t,
			b,
			[]string{
				"",
				"1010",
				"a1010",
				"b1010",
				"ab1010",
			},
		)
	}
	{
		// correct sequence and multi-OR token behaviour
		validateTavorAllPermutations(
			t,
			`
				$Id Sequence = start: 2,
					step: 2

				ExistingLiteral = 1,
					| $Id.Existing,
					| ${Id.Existing + 1}

				And = $Id.Next " " ExistingLiteral " " ExistingLiteral

				START = $Id.Reset And
			`,
			[]string{
				"2 1 1",
				"2 2 1",
				"2 3 1",
				"2 1 2",
				"2 2 2",
				"2 3 2",
				"2 1 3",
				"2 2 3",
				"2 3 3",
			},
		)
	}
	{
		// Correct list pointer behaviour
		validateTavorAllPermutations(
			t,
			`
				$Id Sequence = start: 2,
					step: 2

				Inputs = *(Input)
				Input = $Id.Next

				START = $Id.Reset Inputs
			`,
			[]string{
				"",
				"2",
				"24",
			},
		)
	}
	{
		// Correct sequence deep or behaviour
		validateTavorAllPermutations(
			t,
			`
				$Id Sequence = start: 2,
					step: 2

				A = $Id.Next
				B = $Id.Next (1 | 2 | 3)

				START = $Id.Reset A B
			`,
			[]string{
				"241",
				"242",
				"243",
			},
		)
	}
	{
		// bug
		s := sequences.NewSequence(10, 2)

		a := lists.NewConcatenation(
			s.ResetItem(),
			lists.NewRepeat(lists.NewOne(
				primitives.NewConstantInt(1),
				primitives.NewConstantInt(2),
			), 1, 1),
		)

		validateTokenAllPermutations(
			t,
			a,
			[]string{
				"1",
				"2",
			},
		)
	}
	{
		// dynamic repeat
		validateTavorAllPermutations(
			t,
			`
				As = +0,3(A)
				A = "a"

				Bs = +$As.Count(B)
				B = "b"

				START = As Bs
			`,
			[]string{
				"",
				"ab",
				"aabb",
				"aaabbb",
			},
		)
	}
	{
		// unrolling always ending at end state
		validateTavorAllPermutations(
			t,
			`
				A = "a" (B | C | )
				B = "b" C
				C = "c" A

				START = A
			`,
			[]string{
				"a",
				"abca",
				"abcabca",
				"abcaca",
				"aca",
				"acabca",
				"acaca",
			},
		)
	}
	{
		// corner case of character classes
		validateTavorAllPermutations(
			t,
			`
			START = 1+1([0-1])
			`,
			[]string{
				"10",
				"11",
			},
		)
	}
	{
		// language is not regular let's test that this works
		validateTavorAllPermutations(
			t,
			`
			START = +("a")<A> +$A.Count("b") +$A.Count("c")
			`,
			[]string{
				"abc",
				"aabbcc",
			},
		)
	}
	{
		// correct unqiue behavior
		validateTavorAllPermutations(
			t,
			`
				Items = "a" "b" "c"
				START = Items " -> " $Items.Unique
			`,
			[]string{
				"abc -> a",
				"abc -> b",
				"abc -> c",
			},
		)
	}
	{
		// check if the strategy really works as expected
		validateTavorAllPermutations(
			t,
			`
				START = +2(?(1)?(2))
			`,
			[]string{
				"",
				"1",
				"2",
				"12",
				"1",
				"11",
				"21",
				"121",
				"2",
				"12",
				"22",
				"122",
				"12",
				"112",
				"212",
				"1212",
			},
		)
	}
	{
		// sequences should always start reseted
		validateTavorAllPermutations(
			t,
			`
				$Id Sequence = start: 0,
					step:  2

				START = +1,5($Id.Next " ")
			`,
			[]string{
				"0 ",
				"0 2 ",
				"0 2 4 ",
				"0 2 4 6 ",
				"0 2 4 6 8 ",
			},
		)
	}
	{
		// If a sequence item does not "exist" do not fail on the execution
		validateTavorAllPermutations(
			t,
			`
				$Literal Sequence

				START = "test" $Literal.Existing
			`,
			[]string{
				"test0", // TODO this test should not output any generation, since there is no existing item for $Literal. https://github.com/zimmski/tavor/issues/103
			},
		)
	}
	{
		// Correct next and existing behavior of sequences
		validateTavorAllPermutations(
			t,
			`
				$Literal Sequence

				START = +($Literal.Next " " $Literal.Existing "\n")
			`,
			[]string{
				"1 1\n",
				"1 1\n2 1\n",
				"1 1\n2 1\n", // TODO the number "2" is not used for the existing part https://github.com/zimmski/tavor/issues/12
				"1 1\n2 1\n",
				"1 1\n2 1\n",
			},
		)
	}
}

func validateTavorAllPermutations(t *testing.T, format string, expect []string) {
	r := test.NewRandTest(1)

	o, err := parser.ParseTavor(strings.NewReader(format))
	Nil(t, err)

	var got []string

	ch, err := NewAllPermutations(o, r)
	Nil(t, err)
	for i := range ch {
		got = append(got, o.String())

		ch <- i
	}

	Equal(t, expect, got)
}

func validateTokenAllPermutations(t *testing.T, tok token.Token, expect []string) {
	r := test.NewRandTest(1)

	ch, err := NewAllPermutations(tok, r)
	Nil(t, err)

	var got []string

	for i := range ch {
		got = append(got, tok.String())

		ch <- i
	}

	Equal(t, expect, got)
}

func TestAllPermutationsStrategyLoopDetection(t *testing.T) {
	testStrategyLoopDetection(t, NewAllPermutations)
}
