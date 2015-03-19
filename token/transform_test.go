package token_test

import (
	"testing"

	. "github.com/zimmski/tavor/test/assert"

	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/lists"
	"github.com/zimmski/tavor/token/primitives"
)

func TestUnrollPointers(t *testing.T) {
	{
		// Corner case: START = "group" +(START | "a")
		// This should not leave any pointer in the internal as well as in the external representation

		var tok *token.Token

		p := primitives.NewEmptyPointer(tok)
		s := lists.NewAll(
			primitives.NewConstantString("group"),
			lists.NewRepeat(
				lists.NewOne(
					p,
					primitives.NewConstantString("a"),
				),
				1,
				2,
			),
		)

		err := p.Set(s)
		Nil(t, err)

		unrolled, err := token.UnrollPointers(s)
		Nil(t, err)

		Nil(t, token.WalkInternal(unrolled, func(tok token.Token) error {
			if _, ok := tok.(*primitives.Pointer); ok {
				t.Fatalf("Found pointer in the internal structure %#v", tok)
			}

			return nil
		}))

		Nil(t, token.Walk(unrolled, func(tok token.Token) error {
			if _, ok := tok.(*primitives.Pointer); ok {
				t.Fatalf("Found pointer in the external structure %#v", tok)
			}

			return nil
		}))
	}
}
