package filter

import (
	"math"
	"strconv"

	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/primitives"
)

type PositiveBoundaryValueAnalysisFilter struct{}

func NewPositiveBoundaryValueAnalysisFilter() *PositiveBoundaryValueAnalysisFilter {
	return &PositiveBoundaryValueAnalysisFilter{}
}

func init() {
	Register("PositiveBoundaryValueAnalysis", func() Filter {
		return NewPositiveBoundaryValueAnalysisFilter()
	})
}

func (f *PositiveBoundaryValueAnalysisFilter) Apply(tok token.Token) ([]token.Token, error) {
	if t, ok := tok.(*primitives.RangeInt); ok {
		l := t.Permutations()

		var replacements []token.Token

		if err := t.Permutation(1); err != nil {
			panic(err)
		}

		i, _ := strconv.Atoi(t.String())

		replacements = append(replacements, primitives.NewConstantInt(i))

		if l > 2 {
			if err := t.Permutation(int(math.Ceil(float64(l) / 2.0))); err != nil {
				panic(err)
			}

			i, _ := strconv.Atoi(t.String())

			replacements = append(replacements, primitives.NewConstantInt(i))
		}
		if l != 1 {
			if err := t.Permutation(l); err != nil {
				panic(err)
			}

			i, _ := strconv.Atoi(t.String())

			replacements = append(replacements, primitives.NewConstantInt(i))
		}

		return replacements, nil
	}

	return nil, nil
}
