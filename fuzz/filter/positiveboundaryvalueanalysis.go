package filter

import (
	"strconv"

	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/lists"
	"github.com/zimmski/tavor/token/primitives"
)

// PositiveBoundaryValueAnalysisFilter implements a fuzzing filter for positive boundary-value analysis.
// This filter searches the token graph for range tokens which will be transformed to a most 5 values: the lower and high boundaries as well as the middle values of the range. Using this filter reduces for example integer ranges of 1-100 to the integers 1, 50 and 100, which reduces permutations dramatically. A range of 1-2 will be reduces to the integers 1 and 2. A range of 1 will be reduced to the integer 1. Resulting integers of this filter therefore do not overlap. As a special case, integer ranges where the signs of the two boundaries are different are reduced to a maximum of 5 non-overlapping values. For instance, the integer range [-5, 10] is reduced to the integers -5, -1, 0, 1 and 10.
type PositiveBoundaryValueAnalysisFilter struct{}

// NewPositiveBoundaryValueAnalysisFilter returns a new instance of the positive boundary-value analysis fuzzing filter
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
func (f *PositiveBoundaryValueAnalysisFilter) Apply(tok token.Token) (token.Token, error) {
	var replacements []token.Token

	switch tok := tok.(type) {
	case *primitives.CharacterClass:
		l := tok.Permutations()

		// lower boundary
		if err := tok.Permutation(0); err != nil {
			panic(err)
		}

		replacements = append(replacements, primitives.NewConstantString(tok.String()))

		// middle
		if l > 2 {
			if err := tok.Permutation(l / 2); err != nil {
				panic(err)
			}

			replacements = append(replacements, primitives.NewConstantString(tok.String()))
		}

		// upper boundary
		if l > 1 {
			if err := tok.Permutation(l - 1); err != nil {
				panic(err)
			}

			replacements = append(replacements, primitives.NewConstantString(tok.String()))
		}
	case *primitives.RangeInt:
		l := tok.Permutations()

		// lower boundary
		if err := tok.Permutation(0); err != nil {
			panic(err)
		}

		i, _ := strconv.Atoi(tok.String())

		replacements = append(replacements, primitives.NewConstantInt(i))

		// middle
		if l > 2 {
			if tok.From() < 0 && tok.To() > 0 {
				// the boundaries are -1, 0 and 1
				if tok.From() < -1 {
					replacements = append(replacements, primitives.NewConstantInt(-1))
				}

				replacements = append(replacements, primitives.NewConstantInt(0))

				if tok.To() > 1 {
					replacements = append(replacements, primitives.NewConstantInt(1))
				}
			} else {
				// the boundary is just the middle value
				if err := tok.Permutation(l / 2); err != nil {
					panic(err)
				}

				i, _ := strconv.Atoi(tok.String())

				replacements = append(replacements, primitives.NewConstantInt(i))
			}
		}

		// upper boundary
		if l > 1 {
			if err := tok.Permutation(l - 1); err != nil {
				panic(err)
			}

			i, _ := strconv.Atoi(tok.String())

			replacements = append(replacements, primitives.NewConstantInt(i))
		}
	default:
		return nil, nil
	}

	if len(replacements) == 1 {
		return replacements[0], nil
	}
	return lists.NewOne(replacements...), nil
}
