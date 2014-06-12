package strategy

import (
	"github.com/zimmski/tavor/rand"
	"github.com/zimmski/tavor/token"
)

type Strategy interface {
	Fuzz(tok token.Token, r rand.Rand)
}
