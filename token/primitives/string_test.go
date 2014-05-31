package primitives

import (
	"testing"

	. "github.com/stretchr/testify/assert"

	"github.com/zimmski/tavor/test"
)

func TestConstantString(t *testing.T) {
	o := NewConstantString("abc")
	Equal(t, "abc", o.String())

	r := test.NewRandTest(0)
	o.Fuzz(r)
	Equal(t, "abc", o.String())
}
