package strategy

import (
	"testing"

	. "github.com/zimmski/tavor/test/assert"

	"github.com/zimmski/tavor/test"
	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/constraints"
	"github.com/zimmski/tavor/token/lists"
	"github.com/zimmski/tavor/token/primitives"
	"github.com/zimmski/tavor/token/sequences"
)

func TestPermuteOptionalsStrategyToBeStrategy(t *testing.T) {
	var strat *Strategy

	Implements(t, strat, &PermuteOptionalsStrategy{})
}

func TestPermuteOptionalsfindOptionals(t *testing.T) {
	r := test.NewRandTest(1)

	o := NewPermuteOptionalsStrategy(nil)

	{
		a := primitives.NewConstantInt(1)
		b := constraints.NewOptional(primitives.NewConstantInt(2))
		c := primitives.NewPointer(primitives.NewConstantInt(3))
		d := lists.NewAll(a, b, c)

		optionals := o.findOptionals(r, d, false)

		Equal(t, optionals, []token.OptionalToken{
			b,
		})
	}
	{
		a := constraints.NewOptional(primitives.NewConstantInt(1))
		b := constraints.NewOptional(primitives.NewConstantInt(2))
		c := lists.NewAll(a, b)
		d := constraints.NewOptional(c)

		optionals := o.findOptionals(r, d, false)

		Equal(t, optionals, []token.OptionalToken{
			d,
		})

		optionals[0].(token.OptionalToken).Activate()
		children := o.findOptionals(r, optionals[0], true)

		Equal(t, children, []token.OptionalToken{
			a,
			b,
		})
	}
	{
		a := lists.NewRepeat(primitives.NewConstantInt(1), 0, 10)

		optionals := o.findOptionals(r, a, false)

		Equal(t, optionals, []token.OptionalToken{
			a,
		})

		b := lists.NewRepeat(primitives.NewConstantInt(1), 1, 10)

		optionals = o.findOptionals(r, b, false)

		var nilOpts []token.OptionalToken
		Equal(t, optionals, nilOpts)
	}
}

func TestPermuteOptionalsStrategy(t *testing.T) {
	r := test.NewRandTest(1)

	{
		a := constraints.NewOptional(primitives.NewConstantInt(1))
		b := primitives.NewConstantInt(2)
		c := constraints.NewOptional(primitives.NewConstantInt(3))
		d := lists.NewAll(a, b, c)

		o := NewPermuteOptionalsStrategy(d)

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
		a := constraints.NewOptional(primitives.NewConstantInt(1))
		b := constraints.NewOptional(primitives.NewConstantInt(2))
		c := lists.NewAll(a, b)
		d := constraints.NewOptional(c)

		o := NewPermuteOptionalsStrategy(d)

		var got []string

		ch, err := o.Fuzz(r)
		Nil(t, err)
		for i := range ch {
			got = append(got, d.String())

			ch <- i
		}

		Equal(t, got, []string{
			"",
			"",
			"1",
			"2",
			"12",
		})
	}
	{
		a1 := constraints.NewOptional(primitives.NewConstantInt(1))
		a2 := constraints.NewOptional(primitives.NewConstantInt(11))
		a := constraints.NewOptional(lists.NewAll(a1, a2, primitives.NewConstantString("a")))
		b := constraints.NewOptional(primitives.NewConstantString("b"))
		c := lists.NewAll(a, b, primitives.NewConstantString("c"))
		d := constraints.NewOptional(c)

		o := NewPermuteOptionalsStrategy(d)

		var got []string

		ch, err := o.Fuzz(r)
		Nil(t, err)
		for i := range ch {
			got = append(got, d.String())

			ch <- i
		}

		Equal(t, got, []string{
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
		})
	}
	{
		a := lists.NewAll(
			constraints.NewOptional(primitives.NewConstantInt(1)),
			constraints.NewOptional(primitives.NewConstantInt(2)),
		)
		b := lists.NewRepeat(a, 0, 10)

		o := NewPermuteOptionalsStrategy(b)

		var got []string

		ch, err := o.Fuzz(r)
		Nil(t, err)
		for i := range ch {
			got = append(got, b.String())

			ch <- i
		}

		Equal(t, got, []string{
			"",
			"",
			"1",
			"2",
			"12",
		})
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
		b := lists.NewRepeat(a, 0, 10)

		o := NewPermuteOptionalsStrategy(b)

		var got []string

		ch, err := o.Fuzz(r)
		Nil(t, err)
		for i := range ch {
			got = append(got, b.String())

			ch <- i
		}

		Equal(t, got, []string{
			"",
			"1010",
			"a1010",
			"b1010",
			"ab1010",
		})
	}
}

func TestPermuteOptionalsStrategyLoopDetection(t *testing.T) {
	testStrategyLoopDetection(t, func(root token.Token) Strategy {
		return NewPermuteOptionalsStrategy(root)
	})
}
