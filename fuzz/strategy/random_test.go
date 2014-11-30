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
)

func TestRandomStrategyToBeStrategy(t *testing.T) {
	var strat *Strategy

	Implements(t, strat, &RandomStrategy{})
}

func TestRandomStrategy(t *testing.T) {
	a := primitives.NewRangeInt(5, 10)
	b := primitives.NewRangeInt(5, 10)

	c := constraints.NewOptional(
		lists.NewAll(a, b),
	)

	o, err := New("random", c)
	NotNil(t, o)
	Nil(t, err)

	r := test.NewRandTest(1)

	ch, err := o.Fuzz(r)
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

	ch, err = o.Fuzz(r)
	Nil(t, err)

	_, ok = <-ch
	True(t, ok)

	Equal(t, "", c.String())
	Equal(t, "6", a.String())
	Equal(t, "7", b.String())

	close(ch)

	// run with range
	r.Seed(1)

	ch, err = o.Fuzz(r)
	Nil(t, err)
	for i := range ch {
		Equal(t, "67", c.String())
		Equal(t, "6", a.String())
		Equal(t, "7", b.String())

		ch <- i
	}
}

func TestRandomStrategyCases(t *testing.T) {
	r := test.NewRandTest(1)

	{
		root, err := parser.ParseTavor(strings.NewReader(`
			Items = "a" "b" "c"
			Choice = $Items.Unique<=v> $v.Index " " $v.Value
			START = Items +$Items.Count(Choice)
		`))
		Nil(t, err)

		o, err := New("random", root)
		NotNil(t, o)
		Nil(t, err)

		// run
		{
			r.Seed(0)

			ch, err := o.Fuzz(r)
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

			ch, err := o.Fuzz(r)
			Nil(t, err)

			_, ok := <-ch
			True(t, ok)

			Equal(t, "abc0 a1 b2 c", root.String())

			ch <- struct{}{}

			_, ok = <-ch
			False(t, ok)
		}
	}
}

func TestRandomStrategyLoopDetection(t *testing.T) {
	testStrategyLoopDetection(t, func(root token.Token) Strategy {
		return NewRandomStrategy(root)
	})
}
