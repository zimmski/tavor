package filter

import (
	"strconv"

	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/lists"
	"github.com/zimmski/tavor/token/primitives"
)

// NegativeBoundaryValueAnalysisFilter implements a fuzzing filter for negative boundary-value analysis.
// This filter searches the token graph for integer range tokens which will be transformed to exactly two integers: The lower and higher negative boundary. Using this filter reduces for example the integer range 1-100 to the integers 0 and 101. Which reduces the range away from the model definition and therefore to an invalid data generation, which can be used for example for negative tests.
type NegativeBoundaryValueAnalysisFilter struct{}

// NewNegativeBoundaryValueAnalysisFilter returns a new instance of the negative boundary-value analysis fuzzing filter
func NewNegativeBoundaryValueAnalysisFilter() *NegativeBoundaryValueAnalysisFilter {
	return &NegativeBoundaryValueAnalysisFilter{}
}

func init() {
	Register("NegativeBoundaryValueAnalysis", func() Filter {
		return NewNegativeBoundaryValueAnalysisFilter()
	})
}

// Apply applies the fuzzing filter onto the token and returns a replacement token, or nil if there is no replacement.
// If a fatal error is encountered the error return argument is not nil.
func (f *NegativeBoundaryValueAnalysisFilter) Apply(tok token.Token) (token.Token, error) {
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

	replacements = append(replacements, primitives.NewConstantInt(i-1))

	// upper boundary
	if err := t.Permutation(l); err != nil {
		panic(err)
	}

	i, _ = strconv.Atoi(t.String())

	replacements = append(replacements, primitives.NewConstantInt(i+1))

	return lists.NewOne(replacements...), nil
}
