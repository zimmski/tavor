package strategy

import (
	"testing"

	. "github.com/stretchr/testify/assert"

	"github.com/zimmski/tavor/rand"
	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/primitives"
)

type mockStrategy struct {
	root token.Token
}

func (s *mockStrategy) Fuzz(r rand.Rand) chan struct{} {
	// do nothing

	return nil
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
