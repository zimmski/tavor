package logicals

import (
	"bytes"

	"github.com/zimmski/tavor/rand"
	"github.com/zimmski/tavor/token"
)

type Or struct {
	tokens []token.Token
	value  []token.Token
}

func NewOr(tokens ...token.Token) *Or {
	if len(tokens) == 0 {
		panic("at least one token needed")
	}

	return &Or{
		tokens: tokens,
		value:  tokens,
	}
}

func (o *Or) Fuzz(r rand.Rand) {
	tl := len(o.tokens)

	n := r.Intn(tl) + 1
	toks := make([]token.Token, n)
	chosen := make(map[int]struct{})

	for i, _ := range toks {
		for {
			ri := r.Intn(tl)

			if _, ok := chosen[ri]; !ok {
				toks[i] = o.value[ri]
				chosen[ri] = struct{}{}

				toks[i].Fuzz(r)

				break
			}
		}
	}

	o.value = toks
}

func (o *Or) String() string {
	var buffer bytes.Buffer

	for _, tok := range o.value {
		buffer.WriteString(tok.String())
	}

	return buffer.String()
}
