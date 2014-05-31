package logicals

import (
	"testing"

	. "github.com/stretchr/testify/assert"

	"github.com/zimmski/tavor/test"
	"github.com/zimmski/tavor/token/primitives"
)

func TestXOr(t *testing.T) {
	a := primitives.NewConstantString("a")
	b := primitives.NewConstantString("b")

	o := NewXOr(a, b)
	Equal(t, "a", o.String())

	r := test.NewRandTest(0)
	o.Fuzz(r)
	Equal(t, "b", o.String())

	c := primitives.NewRangeInt(5, 10)
	o = NewXOr(c)
	Equal(t, "5", o.String())

	o.Fuzz(r)
	Equal(t, "6", o.String())
}
