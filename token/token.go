package token

import (
	"fmt"

	"github.com/zimmski/tavor/rand"
)

type Token interface {
	fmt.Stringer

	Clone() Token
	Fuzz(r rand.Rand)
	FuzzAll(r rand.Rand)
	Permutations() int
}
