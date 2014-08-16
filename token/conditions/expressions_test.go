package conditions

import (
	"testing"

	. "github.com/zimmski/tavor/test/assert"

	"github.com/zimmski/tavor/token/primitives"
)

func TestBooleanExpressionToBeBooleanExpression(t *testing.T) {
	var ex *BooleanExpression

	Implements(t, ex, &BooleanTrue{})
	Implements(t, ex, &BooleanEqual{})
}

func TestBooleanTrue(t *testing.T) {
	o := NewBooleanTrue()
	True(t, o.Evaluate())
}

func TestBooleanEqual(t *testing.T) {
	o := NewBooleanEqual(primitives.NewConstantInt(1), primitives.NewConstantInt(1))
	True(t, o.Evaluate())

	o = NewBooleanEqual(primitives.NewConstantInt(1), primitives.NewConstantInt(2))
	False(t, o.Evaluate())
}
