package filter

import (
	"testing"

	. "github.com/zimmski/tavor/test/assert"

	"github.com/zimmski/tavor/token/lists"
	"github.com/zimmski/tavor/token/primitives"
)

func TestNewNegativeBoundaryValueAnalysisFilterToBeFilter(t *testing.T) {
	var filt *Filter

	Implements(t, filt, &NegativeBoundaryValueAnalysisFilter{})
}

func TestNewNegativeBoundaryValueAnalysisFilter(t *testing.T) {
	f := NewNegativeBoundaryValueAnalysisFilter()

	// single value range
	{
		root := primitives.NewRangeInt(10, 10)
		replacements, err := f.Apply(root)
		Nil(t, err)
		Equal(t, replacements, lists.NewOne(
			primitives.NewConstantInt(9),
			primitives.NewConstantInt(11),
		))
	}
	// two value range
	{
		root := primitives.NewRangeInt(10, 11)
		replacements, err := f.Apply(root)
		Nil(t, err)
		Equal(t, replacements, lists.NewOne(
			primitives.NewConstantInt(9),
			primitives.NewConstantInt(12),
		))
	}
	// three value range
	{
		root := primitives.NewRangeInt(10, 12)
		replacements, err := f.Apply(root)
		Nil(t, err)
		Equal(t, replacements, lists.NewOne(
			primitives.NewConstantInt(9),
			primitives.NewConstantInt(13),
		))
	}
	// four value range
	{
		root := primitives.NewRangeInt(10, 13)
		replacements, err := f.Apply(root)
		Nil(t, err)
		Equal(t, replacements, lists.NewOne(
			primitives.NewConstantInt(9),
			primitives.NewConstantInt(14),
		))
	}
	// five value range
	{
		root := primitives.NewRangeInt(10, 14)
		replacements, err := f.Apply(root)
		Nil(t, err)
		Equal(t, replacements, lists.NewOne(
			primitives.NewConstantInt(9),
			primitives.NewConstantInt(15),
		))
	}
}
