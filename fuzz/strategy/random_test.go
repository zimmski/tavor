package strategy

import (
	"testing"

	. "github.com/stretchr/testify/assert"

	"github.com/zimmski/tavor/test"
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

	r := test.NewRandTest(0)
	o.Fuzz(r)
	Equal(t, "78", c.String())
	Equal(t, "7", a.String())
	Equal(t, "8", b.String())

	// optional should be off but since we fuzz everything no matter
	// what, a and b should still get fuzzed
	r.Seed(1)
	o.Fuzz(r)
	Equal(t, "", c.String())
	Equal(t, "8", a.String())
	Equal(t, "9", b.String())
}
