package filter

import (
	"testing"

	. "github.com/zimmski/tavor/test/assert"

	"github.com/zimmski/tavor/token/lists"
	"github.com/zimmski/tavor/token/primitives"
)

func TestNewPositiveBoundaryValueAnalysisFilter(t *testing.T) {
	// single value range
	{
		root := primitives.NewRangeInt(10, 10)
		replacements, err := NewPositiveBoundaryValueAnalysis(root)
		Nil(t, err)
		Equal(t, replacements, primitives.NewConstantInt(10))
	}
	// two value range
	{
		root := primitives.NewRangeInt(10, 11)
		replacements, err := NewPositiveBoundaryValueAnalysis(root)
		Nil(t, err)
		Equal(t, replacements, lists.NewOne(
			primitives.NewConstantInt(10),
			primitives.NewConstantInt(11),
		))
	}
	// three value range
	{
		root := primitives.NewRangeInt(10, 12)
		replacements, err := NewPositiveBoundaryValueAnalysis(root)
		Nil(t, err)
		Equal(t, replacements, lists.NewOne(
			primitives.NewConstantInt(10),
			primitives.NewConstantInt(11),
			primitives.NewConstantInt(12),
		))
	}
	// four value range
	{
		root := primitives.NewRangeInt(10, 13)
		replacements, err := NewPositiveBoundaryValueAnalysis(root)
		Nil(t, err)
		Equal(t, replacements, lists.NewOne(
			primitives.NewConstantInt(10),
			primitives.NewConstantInt(12),
			primitives.NewConstantInt(13),
		))
	}
	// five value range
	{
		root := primitives.NewRangeInt(10, 14)
		replacements, err := NewPositiveBoundaryValueAnalysis(root)
		Nil(t, err)
		Equal(t, replacements, lists.NewOne(
			primitives.NewConstantInt(10),
			primitives.NewConstantInt(12),
			primitives.NewConstantInt(14),
		))
	}
	// negative range
	{
		root := primitives.NewRangeInt(-14, -10)
		replacements, err := NewPositiveBoundaryValueAnalysis(root)
		Nil(t, err)
		Equal(t, replacements, lists.NewOne(
			primitives.NewConstantInt(-14),
			primitives.NewConstantInt(-12),
			primitives.NewConstantInt(-10),
		))
	}
	// negative to positive range
	{
		root := primitives.NewRangeInt(-5, 10)
		replacements, err := NewPositiveBoundaryValueAnalysis(root)
		Nil(t, err)
		Equal(t, replacements, lists.NewOne(
			primitives.NewConstantInt(-5),
			primitives.NewConstantInt(-1),
			primitives.NewConstantInt(0),
			primitives.NewConstantInt(1),
			primitives.NewConstantInt(10),
		))
	}
	// three value CharacterClass
	{
		root := primitives.NewCharacterClass("a-z")
		replacements, err := NewPositiveBoundaryValueAnalysis(root)
		Nil(t, err)
		Equal(t, replacements, lists.NewOne(
			primitives.NewConstantString("a"),
			primitives.NewConstantString("n"),
			primitives.NewConstantString("z"),
		))
	}
}
