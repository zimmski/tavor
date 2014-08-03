package strategy

import (
	"bytes"
	"testing"

	. "github.com/stretchr/testify/assert"

	"github.com/zimmski/tavor/parser"
	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/constraints"
	"github.com/zimmski/tavor/token/lists"
	"github.com/zimmski/tavor/token/primitives"
)

func TestBinarySearchStrategyToBeStrategy(t *testing.T) {
	var strat *Strategy

	Implements(t, strat, &BinarySearchStrategy{})
}

func TestBinarySearchStrategy(t *testing.T) {
	{
		root := primitives.NewConstantInt(1)

		o := NewBinarySearch(root)

		contin, _, err := o.Reduce()
		Nil(t, err)

		_, ok := <-contin
		False(t, ok)

		Equal(t, "1", root.String())
	}
	{
		c := constraints.NewOptional(
			primitives.NewConstantInt(2),
		)
		c.Activate()
		root := lists.NewAll(
			primitives.NewConstantInt(1),
			c,
		)

		o := NewBinarySearch(root)

		contin, feedback, err := o.Reduce()
		Nil(t, err)

		_, ok := <-contin
		True(t, ok)

		Equal(t, "1", root.String())

		feedback <- Bad
		contin <- struct{}{}

		_, ok = <-contin
		False(t, ok)

		Equal(t, "1", root.String())
	}
	{
		c := constraints.NewOptional(
			primitives.NewConstantInt(2),
		)
		c.Activate()
		root := lists.NewAll(
			primitives.NewConstantInt(1),
			c,
		)

		o := NewBinarySearch(root)

		contin, feedback, err := o.Reduce()
		Nil(t, err)

		_, ok := <-contin
		True(t, ok)

		Equal(t, "1", root.String())

		feedback <- Good
		contin <- struct{}{}

		_, ok = <-contin
		False(t, ok)

		Equal(t, "12", root.String())
	}
	// Test that inputs are never changed if they cannot be reduced
	{
		root := lists.NewRepeat(primitives.NewCharacterClass(`\w`), 10, 10)
		input := "KrOxDOj4fU"

		errs := parser.ParseInternal(root, bytes.NewBufferString(input))
		Nil(t, errs)

		Equal(t, input, root.String())

		o := NewBinarySearch(root)

		contin, _, err := o.Reduce()
		Nil(t, err)

		_, ok := <-contin
		False(t, ok)

		Equal(t, input, root.String())
	}
}

func TestBinarySearchStrategyLoopDetection(t *testing.T) {
	testStrategyLoopDetection(t, func(root token.Token) Strategy {
		return NewBinarySearch(root)
	})
}
