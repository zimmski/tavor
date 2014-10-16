package filter

import (
	"math"
	"strconv"

	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/primitives"
)

// PositiveBoundaryValueAnalysisFilter implements a fuzzing filter for positive boundary value analysis.
// This filter searches the token graph for integer range tokens which will be transformed to at most three new values: The lower and higher boundary as well as the value exactly at the middle of the range. Using this filter reduces for example integer ranges of 1-100 to the integers 1, 50 and 100. Which reduces permutations dramatically. A range of 1-2 will be reduces to the integers 1 and 2. A range of 1 will be reduced to the integer 1. Resulting integers of this filter therefore do not overlap.
type PositiveBoundaryValueAnalysisFilter struct{}

// NewPositiveBoundaryValueAnalysisFilter returns a new instance of the positive boundary value analysis fuzzing filter
func NewPositiveBoundaryValueAnalysisFilter() *PositiveBoundaryValueAnalysisFilter {
	return &PositiveBoundaryValueAnalysisFilter{}
}

func init() {
	Register("PositiveBoundaryValueAnalysis", func() Filter {
		return NewPositiveBoundaryValueAnalysisFilter()
	})
}

// Apply applies the fuzzing filter onto the token and returns a replacement token, or nil if there is no replacement.
// If a fatal error is encountered the error return argument is not nil.
func (f *PositiveBoundaryValueAnalysisFilter) Apply(tok token.Token) ([]token.Token, error) {
	t, ok := tok.(*primitives.RangeInt)
	if !ok {
		return nil, nil
	}

	l := t.Permutations()

	var replacements []token.Token

	// lower boundary
	if err := t.Permutation(1); err != nil {
		panic(err)
	}

	i, _ := strconv.Atoi(t.String())

	replacements = append(replacements, primitives.NewConstantInt(i))

	// middle
	if l > 2 {
		if err := t.Permutation(uint(math.Ceil(float64(l) / 2.0))); err != nil {
			panic(err)
		}

		i, _ := strconv.Atoi(t.String())

		replacements = append(replacements, primitives.NewConstantInt(i))
	}

	// upper boundary
	if l > 1 {
		if err := t.Permutation(l); err != nil {
			panic(err)
		}

		i, _ := strconv.Atoi(t.String())

		replacements = append(replacements, primitives.NewConstantInt(i))
	}

	return replacements, nil
}
