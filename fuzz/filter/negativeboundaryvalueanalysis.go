package filter

import (
	"strconv"

	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/primitives"
)

type NegativeBoundaryValueAnalysisFilter struct{}

func NewNegativeBoundaryValueAnalysisFilter() *NegativeBoundaryValueAnalysisFilter {
	return &NegativeBoundaryValueAnalysisFilter{}
}

func init() {
	Register("NegativeBoundaryValueAnalysis", func() Filter {
		return NewNegativeBoundaryValueAnalysisFilter()
	})
}

func (f *NegativeBoundaryValueAnalysisFilter) Apply(tok token.Token) ([]token.Token, error) {
	if t, ok := tok.(*primitives.RangeInt); ok {
		l := t.Permutations()

		var replacements []token.Token

		// lower boundary
		if err := t.Permutation(1); err != nil {
			panic(err)
		}

		i, _ := strconv.Atoi(t.String())

		replacements = append(replacements, primitives.NewConstantInt(i-1))

		// upper boundary
		if err := t.Permutation(l); err != nil {
			panic(err)
		}

		i, _ = strconv.Atoi(t.String())

		replacements = append(replacements, primitives.NewConstantInt(i+1))

		return replacements, nil
	}

	return nil, nil
}
