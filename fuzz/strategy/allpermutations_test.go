package strategy

import (
	"strings"
	"testing"

	"github.com/zimmski/go-leak"
	. "github.com/zimmski/tavor/test/assert"

	"github.com/zimmski/tavor/parser"
	"github.com/zimmski/tavor/test"
	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/constraints"
	"github.com/zimmski/tavor/token/lists"
	"github.com/zimmski/tavor/token/primitives"
	"github.com/zimmski/tavor/token/sequences"
)

func TestAllPermutationsStrategyToBeStrategy(t *testing.T) {
	var strat *Strategy

	Implements(t, strat, &AllPermutationsStrategy{})
}

func TestAllPermutationsStrategygetLevel(t *testing.T) {
	o := NewAllPermutationsStrategy(nil)

	var nilChildren []allPermutationsLevel

	{
		a := primitives.NewConstantInt(1)
		b1 := primitives.NewConstantInt(2)
		b := constraints.NewOptional(b1)
		c1 := primitives.NewConstantInt(3)
		c := primitives.NewPointer(c1)
		d := lists.NewAll(a, b, c)

		tree := o.getTree(d, false)

		Equal(t, tree, []allPermutationsLevel{
			allPermutationsLevel{
				token:       d,
				permutation: 1,

				children: []allPermutationsLevel{
					allPermutationsLevel{
						token:       a,
						permutation: 1,

						children: nilChildren,
					},
					allPermutationsLevel{
						token:       b,
						permutation: 1,

						children: nilChildren,
					},
					allPermutationsLevel{
						token:       c,
						permutation: 1,

						children: []allPermutationsLevel{
							allPermutationsLevel{
								token:       c1,
								permutation: 1,

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
				permutation: 1,

				children: nilChildren,
			},
			allPermutationsLevel{
				token:       b,
				permutation: 1,

				children: nilChildren,
			},
			allPermutationsLevel{
				token:       c,
				permutation: 1,

				children: []allPermutationsLevel{
					allPermutationsLevel{
						token:       c1,
						permutation: 1,

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

		o := NewAllPermutationsStrategy(a)

		ch, err := o.Fuzz(r)
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
		abc := lists.NewAll(a, b, c)

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
		abc := lists.NewAll(
			constraints.NewOptional(lists.NewAll(
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
		d := lists.NewAll(a, b, c)

		o := NewAllPermutationsStrategy(d)

		ch, err := o.Fuzz(r)
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
		ch, err = o.Fuzz(r)
		Nil(t, err)

		_, ok = <-ch
		True(t, ok)
		Equal(t, "2", d.String())

		close(ch)

		// run with range
		var got []string

		ch, err = o.Fuzz(r)
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
		a := constraints.NewOptional(lists.NewAll(a1, a2, primitives.NewConstantString("a")))
		b := constraints.NewOptional(primitives.NewConstantString("b"))
		c := lists.NewAll(a, b, primitives.NewConstantString("c"))
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
		a := lists.NewAll(
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

		a := lists.NewAll(
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

		a := lists.NewAll(
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
}

func validateTavorAllPermutations(t *testing.T, format string, expect []string) {
	m := leak.MarkGoRoutines()

	r := test.NewRandTest(1)

	o, err := parser.ParseTavor(strings.NewReader(format))
	Nil(t, err)

	s := NewAllPermutationsStrategy(o)

	var got []string

	ch, err := s.Fuzz(r)
	Nil(t, err)
	for i := range ch {
		got = append(got, o.String())

		ch <- i
	}

	Equal(t, got, expect)

	Equal(t, 0, m.Release(), "check for goroutine leaks")
}

func validateTokenAllPermutations(t *testing.T, tok token.Token, expect []string) {
	r := test.NewRandTest(1)

	o := NewAllPermutationsStrategy(tok)

	ch, err := o.Fuzz(r)
	Nil(t, err)

	var got []string

	for i := range ch {
		got = append(got, tok.String())

		ch <- i
	}

	Equal(t, got, expect)
}

func TestAllPermutationsStrategyLoopDetection(t *testing.T) {
	testStrategyLoopDetection(t, func(root token.Token) Strategy {
		return NewAllPermutationsStrategy(root)
	})
}
