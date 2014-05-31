package logicals

import (
	"bytes"

	"github.com/zimmski/tavor/rand"
	"github.com/zimmski/tavor/token"
)

type Many struct {
	tokens []token.Token
	value  []token.Token
}

func NewMany(toks ...token.Token) *Many {
	if len(toks) == 0 {
		panic("at least one token needed")
	}

	return &Many{
		tokens: toks,
		value:  toks,
	}
}

func (o *Many) Fuzz(r rand.Rand) {
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

func (o *Many) String() string {
	var buffer bytes.Buffer

	for _, tok := range o.value {
		buffer.WriteString(tok.String())
	}

	return buffer.String()
}
