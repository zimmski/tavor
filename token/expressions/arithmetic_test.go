package expressions

import (
	"testing"

	. "github.com/stretchr/testify/assert"

	"github.com/zimmski/tavor/test"
	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/primitives"
)

func TestArithmeticExpressionTokensToBeTokens(t *testing.T) {
	var tok *token.Token

	Implements(t, tok, &AddArithmetic{})
	Implements(t, tok, &SubArithmetic{})
	Implements(t, tok, &MulArithmetic{})
	Implements(t, tok, &DivArithmetic{})
}

func TestAddArithmetic(t *testing.T) {
	o := NewAddArithmetic(
		primitives.NewRangeInt(1, 10),
		primitives.NewConstantInt(2),
	)
	Equal(t, "3", o.String())

	r := test.NewRandTest(1)
	o.Fuzz(r)
	Equal(t, "5", o.String())

	o2 := o.Clone()
	Equal(t, o.String(), o2.String())
}

func TestSubArithmetic(t *testing.T) {
	o := NewSubArithmetic(
		primitives.NewRangeInt(1, 10),
		primitives.NewConstantInt(2),
	)
	Equal(t, "-1", o.String())

	r := test.NewRandTest(1)
	o.Fuzz(r)
	Equal(t, "1", o.String())

	o2 := o.Clone()
	Equal(t, o.String(), o2.String())
}

func TestMulArithmetic(t *testing.T) {
	o := NewMulArithmetic(
		primitives.NewRangeInt(1, 10),
		primitives.NewConstantInt(2),
	)
	Equal(t, "2", o.String())

	r := test.NewRandTest(1)
	o.Fuzz(r)
	Equal(t, "6", o.String())

	o2 := o.Clone()
	Equal(t, o.String(), o2.String())
}

func TestDivArithmetic(t *testing.T) {
	o := NewDivArithmetic(
		primitives.NewRangeInt(6, 10),
		primitives.NewConstantInt(2),
	)
	Equal(t, "3", o.String())

	r := test.NewRandTest(2)
	o.Fuzz(r)
	Equal(t, "4", o.String())

	o2 := o.Clone()
	Equal(t, o.String(), o2.String())
}
