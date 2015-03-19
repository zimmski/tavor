package strategy

import (
	"testing"

	. "github.com/zimmski/tavor/test/assert"

	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/constraints"
	"github.com/zimmski/tavor/token/lists"
	"github.com/zimmski/tavor/token/primitives"
)

type mockStrategy struct {
	root token.Token
}

func (s *mockStrategy) Reduce() (chan struct{}, chan<- ReduceFeedbackType, error) {
	// do nothing

	return nil, nil, nil
}

func TestStrategy(t *testing.T) {
	a := primitives.NewConstantInt(123)

	// mock is not registered
	for _, name := range List() {
		if name == "mock" {
			Fail(t, "mock should not be in the strategy list yet")
		}
	}

	stat, err := New("mock", a)
	Nil(t, stat)
	NotNil(t, err)

	// register mock
	Register("mock", func(tok token.Token) Strategy {
		return &mockStrategy{
			root: tok,
		}
	})

	// mock is registered
	found := false
	for _, name := range List() {
		if name == "mock" {
			found = true

			break
		}
	}
	True(t, found)

	stat, err = New("mock", a)
	NotNil(t, stat)
	True(t, Exactly(t, a, stat.(*mockStrategy).root))
	Nil(t, err)

	// register mock a second time
	caught := false
	func() {
		defer func() {
			if r := recover(); r != nil {
				caught = true
			}
		}()

		Register("mock", func(tok token.Token) Strategy {
			return &mockStrategy{
				root: tok,
			}
		})
	}()
	True(t, caught)

	// register nil function
	caught = false
	func() {
		defer func() {
			if r := recover(); r != nil {
				caught = true
			}
		}()

		Register("mockachino", nil)
	}()
	True(t, caught)
}

func testStrategyLoopDetection(t *testing.T, newStrategy func(root token.Token) Strategy) {
	var tok *token.Token

	{
		// allow unloopy pointers

		a := primitives.NewConstantInt(2)
		p := primitives.NewPointer(a)
		o := lists.NewAll(
			p,
			primitives.NewConstantInt(1),
		)

		s := newStrategy(o)

		contin, feedback, err := s.Reduce()
		NotNil(t, contin)
		NotNil(t, feedback)
		Nil(t, err)
	}
	{
		// check for simple loops

		p := primitives.NewEmptyPointer(tok)
		o := lists.NewAll(
			p,
			primitives.NewConstantInt(1),
		)
		Nil(t, p.Set(o))

		s := newStrategy(o)

		contin, feedback, err := s.Reduce()
		Nil(t, contin)
		Nil(t, feedback)
		Equal(t, ErrorEndlessLoopDetected, err.(*Error).Type)

		p = primitives.NewEmptyPointer(tok)
		o = lists.NewAll(
			primitives.NewConstantInt(1),
			p,
		)
		Nil(t, p.Set(o))

		s = newStrategy(o)

		contin, feedback, err = s.Reduce()
		Nil(t, contin)
		Nil(t, feedback)
		Equal(t, ErrorEndlessLoopDetected, err.(*Error).Type)
	}
	{
		// deeper loops

		p := primitives.NewEmptyPointer(tok)
		o := lists.NewAll(
			lists.NewOne(
				p,
				primitives.NewConstantInt(1),
			),
			lists.NewOne(
				p,
				primitives.NewConstantInt(2),
			),
			constraints.NewOptional(
				p,
			),
		)
		Nil(t, p.Set(o))

		s := newStrategy(o)

		contin, feedback, err := s.Reduce()
		Nil(t, contin)
		Nil(t, feedback)
		Equal(t, ErrorEndlessLoopDetected, err.(*Error).Type)
	}
}
