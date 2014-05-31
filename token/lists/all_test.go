package logicals

import (
	"testing"

	. "github.com/stretchr/testify/assert"

	"github.com/zimmski/tavor/test"
	"github.com/zimmski/tavor/token/primitives"
)

func TestAll(t *testing.T) {
	a := primitives.NewConstantInt(10)
	b := primitives.NewConstantString("abc")

	o := NewAll(a, b)
	Equal(t, "10abc", o.String())

	r := test.NewRandTest(0)
	o.Fuzz(r)
	Equal(t, "10abc", o.String())

	c := primitives.NewRangeInt(1, 2)
	o = NewAll(a, b, c)
	Equal(t, "10abc1", o.String())

	o.Fuzz(r)
	Equal(t, "10abc2", o.String())
}
