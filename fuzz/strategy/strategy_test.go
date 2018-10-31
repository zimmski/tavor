package strategy

import (
	"testing"

	. "github.com/zimmski/tavor/test/assert"

	"github.com/zimmski/tavor/rand"
	"github.com/zimmski/tavor/test"
	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/constraints"
	"github.com/zimmski/tavor/token/lists"
	"github.com/zimmski/tavor/token/primitives"
)

func newMockStrategy(root token.Token, r rand.Rand) (chan struct{}, error) {
	// do nothing

	return nil, nil
}

func TestStrategy(t *testing.T) {
	// mock is not registered
	for _, name := range List() {
		if name == "mock" {
			Fail(t, "mock should not be in the strategy list yet")
		}
	}

	strat, err := New("mock")
	Nil(t, strat)
	NotNil(t, err)

	// register mock
	Register("mock", newMockStrategy)

	// mock is registered
	found := false
	for _, name := range List() {
		if name == "mock" {
			found = true

			break
		}
	}
	True(t, found)

	strat, err = New("mock")
	NotNil(t, strat)
	Nil(t, err)

	// register mock a second time
	caught := false
	func() {
		defer func() {
			if r := recover(); r != nil {
				caught = true
			}
		}()

		Register("mock", newMockStrategy)
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

func testStrategyLoopDetection(t *testing.T, newStrategy Strategy) {
	var tok *token.Token
	r := test.NewRandTest(1)

	{
		// allow unloopy pointers

		a := primitives.NewConstantInt(2)
		p := primitives.NewPointer(a)
		o := lists.NewConcatenation(
			p,
			primitives.NewConstantInt(1),
		)

		ch, err := newStrategy(o, r)
		NotNil(t, ch)
		Nil(t, err)
	}
	{
		// check for simple loops

		p := primitives.NewEmptyPointer(tok)
		o := lists.NewConcatenation(
			p,
			primitives.NewConstantInt(1),
		)
		Nil(t, p.Set(o))

		ch, err := newStrategy(o, r)
		Nil(t, ch)
		Equal(t, ErrEndlessLoopDetected, err.(*Error).Type)

		p = primitives.NewEmptyPointer(tok)
		o = lists.NewConcatenation(
			primitives.NewConstantInt(1),
			p,
		)
		Nil(t, p.Set(o))

		ch, err = newStrategy(o, r)
		Nil(t, ch)
		Equal(t, ErrEndlessLoopDetected, err.(*Error).Type)
	}
	{
		// deeper loops

		p := primitives.NewEmptyPointer(tok)
		o := lists.NewConcatenation(
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

		ch, err := newStrategy(o, r)
		Nil(t, ch)
		Equal(t, ErrEndlessLoopDetected, err.(*Error).Type)
	}
}
