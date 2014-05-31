package token

import (
	"testing"

	"github.com/zimmski/tavor/token/primitives"
)

func assertToken(tok Token) {

}

func TestAllTokensToBeTokens(t *testing.T) {
	// primitives
	{
		// Int
		assertToken(primitives.NewConstantInt(10))
		assertToken(primitives.NewRandomInt())
		assertToken(primitives.NewRangeInt(2, 4))
	}
}
