package strategy

import (
	"strings"
	"testing"

	. "github.com/zimmski/tavor/test/assert"

	"github.com/zimmski/tavor/parser"
	"github.com/zimmski/tavor/test"
	"github.com/zimmski/tavor/token/constraints"
	"github.com/zimmski/tavor/token/lists"
	"github.com/zimmski/tavor/token/primitives"
)

func TestRandomStrategyNilRandomGenerator(t *testing.T) {
	ch, err := NewRandom(nil, nil)
	Nil(t, ch)
	Equal(t, ErrNilRandomGenerator, err.(*Error).Type)
}

func TestRandomStrategy(t *testing.T) {
	a := primitives.NewRangeInt(5, 10)
	b := primitives.NewRangeInt(5, 10)

	c := constraints.NewOptional(
		lists.NewConcatenation(a, b),
	)

	r := test.NewRandTest(1)

	ch, err := NewRandom(c, r)
	Nil(t, err)

	_, ok := <-ch
	True(t, ok)

	Equal(t, "67", c.String())
	Equal(t, "6", a.String())
	Equal(t, "7", b.String())

	ch <- struct{}{}

	_, ok = <-ch
	False(t, ok)

	// rerun
	r.Seed(0)

	ch, err = NewRandom(c, r)
	Nil(t, err)

	_, ok = <-ch
	True(t, ok)

	Equal(t, "", c.String())
	Equal(t, "6", a.String())
	Equal(t, "7", b.String())

	close(ch)

	// run with range
	r.Seed(1)

	ch, err = NewRandom(c, r)
	Nil(t, err)
	for i := range ch {
		Equal(t, "67", c.String())
		Equal(t, "6", a.String())
		Equal(t, "7", b.String())

		ch <- i
	}
}

func TestRandomStrategyCases(t *testing.T) {
	{
		r := test.NewRandTest(1)

		root, err := parser.ParseTavor(strings.NewReader(`
			Items = "a" "b" "c"
			Choice = $Items.Unique<=v> $v.Index " " $v.Value
			START = Items +$Items.Count(Choice)
		`))
		Nil(t, err)

		// run
		{
			r.Seed(0)

			ch, err := NewRandom(root, r)
			Nil(t, err)

			_, ok := <-ch
			True(t, ok)

			Equal(t, "abc0 a1 b2 c", root.String())

			ch <- struct{}{}

			_, ok = <-ch
			False(t, ok)
		}

		// rerun
		{
			r.Seed(1)

			ch, err := NewRandom(root, r)
			Nil(t, err)

			_, ok := <-ch
			True(t, ok)

			Equal(t, "abc0 a1 b2 c", root.String())

			ch <- struct{}{}

			_, ok = <-ch
			False(t, ok)
		}
	}
	{
		// sequences should always start reseted
		validateTavorRandom(
			t,
			1,
			`
				$Id Sequence = start: 0,
					step:  2

				START = +1,5($Id.Next " ")
			`,
			[]string{
				"0 2 ",
			},
		)
	}
	{
		// If a sequence item does not "exist" do not fail on the execution
		validateTavorRandom(
			t,
			1,
			`
				$Literal Sequence

				START = "test" $Literal.Existing
			`,
			[]string{
				"test0", // TODO this test should not output any generation, since there is no existing item for $Literal. https://github.com/zimmski/tavor/issues/103
			},
		)
	}
}

func validateTavorRandom(t *testing.T, seed int, format string, expect []string) {
	root, err := parser.ParseTavor(strings.NewReader(format))
	Nil(t, err)

	r := test.NewRandTest(int64(seed))

	ch, err := NewRandom(root, r)
	Nil(t, err)

	_, ok := <-ch
	True(t, ok)

	var got []string
	got = append(got, root.String())

	ch <- struct{}{}

	_, ok = <-ch
	False(t, ok)

	Equal(t, got, expect)
}

func TestRandomStrategyLoopDetection(t *testing.T) {
	testStrategyLoopDetection(t, NewRandom)
}
