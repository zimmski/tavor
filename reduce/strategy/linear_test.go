package strategy

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	. "github.com/zimmski/tavor/test/assert"

	"github.com/zimmski/tavor/parser"
	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/constraints"
	"github.com/zimmski/tavor/token/lists"
	"github.com/zimmski/tavor/token/primitives"
)

func TestLinearStrategy(t *testing.T) {
	{
		root := primitives.NewConstantInt(1)

		contin, _, err := NewLinear(root)
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
		root := lists.NewConcatenation(
			primitives.NewConstantInt(1),
			c,
		)

		contin, feedback, err := NewLinear(root)
		Nil(t, err)

		_, ok := <-contin
		True(t, ok)

		Equal(t, "1", root.String())

		feedback <- Bad
		contin <- struct{}{}

		_, ok = <-contin
		False(t, ok)

		Equal(t, "12", root.String())
	}
	{
		c := constraints.NewOptional(
			primitives.NewConstantInt(2),
		)
		c.Activate()
		root := lists.NewConcatenation(
			primitives.NewConstantInt(1),
			c,
		)

		contin, feedback, err := NewLinear(root)
		Nil(t, err)

		_, ok := <-contin
		True(t, ok)

		Equal(t, "1", root.String())

		feedback <- Good
		contin <- struct{}{}

		_, ok = <-contin
		False(t, ok)

		Equal(t, "1", root.String())
	}
	{
		// Test that inputs are never changed if they cannot be reduced

		root := lists.NewRepeat(primitives.NewCharacterClass(`\w`), 10, 10)
		input := "KrOxDOj4fU"

		errs := parser.ParseInternal(root, bytes.NewBufferString(input))
		Nil(t, errs)

		Equal(t, input, root.String())

		contin, _, err := NewLinear(root)
		Nil(t, err)

		_, ok := <-contin
		False(t, ok)

		Equal(t, input, root.String())
	}
	{
		// Test that reductions grow again
		tok := lists.NewRepeat(primitives.NewConstantString("a"), 0, 100)

		validateTavorLinear(
			t,
			tok,
			"aaaaaa",
			func(out string) ReduceFeedbackType {
				if len(out) == 2 {
					return Good
				}

				return Bad
			},
			[]string{
				"",
				"a",
				"a",
				"a",
				"a",
				"a",
				"a",
				"aa",
			},
			"aa",
		)
	}
	{
		// Pathed reduction
		tok, err := parser.ParseTavor(bytes.NewBufferString(`
			START = A *(B)

			A = ?("a")
			B = "b" (C | )
			C = "c"
		`))
		Nil(t, err)

		validateTavorLinear(
			t,
			tok,
			"abc",
			func(out string) ReduceFeedbackType {
				if strings.Contains(out, "b") {
					return Good
				}

				return Bad
			},
			[]string{
				"bc",
				"",
				"b",
			},
			"b",
		)
	}

	/*
		TODO read this files in with the aag.tavor file.
		Since they have a AND loop in them, they should be not valid

		aag 3 0 0 0 2
		2 5 4
		4 3 0

		aag 3 0 0 0 3
		2 5 4
		4 6 0
		6 2 0
	*/
}

func validateTavorLinear(t *testing.T, tok token.Token, input string, feedback func(out string) ReduceFeedbackType, expected []string, final string) {
	errs := parser.ParseInternal(tok, bytes.NewBufferString(input))
	Nil(t, errs)
	if errs != nil {
		panic(errs)
	}

	Equal(t, input, tok.String(), "Generation 0")

	continueFuzzing, feedbackReducing, err := NewLinear(tok)
	if err != nil {
		panic(err)
	}

	n := 0

	for i := range continueFuzzing {
		out := tok.String()

		if n == len(expected) {
			Fail(t, fmt.Sprintf("%q is an unexpected generation at index  %d", out, n))
		} else {
			Equal(t, expected[n], out, fmt.Sprintf("Generation %d", n))
		}
		n++

		feedbackReducing <- feedback(out)

		continueFuzzing <- i
	}

	Equal(t, final, tok.String(), "Final generation")
}

func TestLinearStrategyLoopDetection(t *testing.T) {
	testStrategyLoopDetection(t, NewLinear)
}
