package strategy

import (
	"fmt"
	"math"

	"github.com/zimmski/container/list/linkedlist"

	"github.com/zimmski/tavor"
	"github.com/zimmski/tavor/rand"
	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/lists"
)

type PermuteOptionalsStrategy struct {
	root token.Token
}

func NewPermuteOptionalsStrategy(tok token.Token) *PermuteOptionalsStrategy {
	s := &PermuteOptionalsStrategy{
		root: tok,
	}

	return s
}

func init() {
	Register("PermuteOptionals", func(tok token.Token) Strategy {
		return NewPermuteOptionalsStrategy(tok)
	})
}

func (s *PermuteOptionalsStrategy) findOptionals(r rand.Rand, root token.Token, fromChilds bool) ([]token.OptionalToken, map[token.ResetToken]struct{}) {
	var optionals []token.OptionalToken
	var queue = linkedlist.New()
	var resets = make(map[token.ResetToken]struct{})

	if fromChilds {
		switch t := root.(type) {
		case token.ForwardToken:
			queue.Push(t.Get())
		case lists.List:
			l := t.Len()

			for i := 0; i < l; i++ {
				c, _ := t.Get(i)
				queue.Push(c)
			}
		}
	} else {
		queue.Push(root)
	}

	for !queue.Empty() {
		tok, _ := queue.Shift()

		switch t := tok.(type) {
		case token.ResetToken:
			resets[t] = struct{}{}
		}

		switch t := tok.(type) {
		case token.OptionalToken:
			if !t.IsOptional() {
				opts, rets := s.findOptionals(r, t, true)

				if len(opts) != 0 {
					optionals = append(optionals, opts...)
				}
				if len(rets) != 0 {
					for t := range rets {
						resets[t] = struct{}{}
					}
				}

				continue
			}

			if tavor.DEBUG {
				fmt.Printf("Found optional %#v\n", t)
			}

			t.Deactivate()

			optionals = append(optionals, t)
		case token.ForwardToken:
			c := t.Get()

			c.Fuzz(r)

			queue.Push(c)
		case lists.List:
			l := t.Len()

			for i := 0; i < l; i++ {
				c, _ := t.Get(i)

				c.Fuzz(r)

				queue.Push(c)
			}
		}
	}

	return optionals, resets
}

func (s *PermuteOptionalsStrategy) resetResetTokens() {
	var queue = linkedlist.New()

	queue.Push(s.root)

	for !queue.Empty() {
		v, _ := queue.Shift()

		switch tok := v.(type) {
		case token.ResetToken:
			if tavor.DEBUG {
				fmt.Printf("Reset %#v(%p)\n", tok, tok)
			}

			tok.Reset()
		}

		switch tok := v.(type) {
		case token.ForwardToken:
			if v := tok.Get(); v != nil {
				queue.Push(v)
			}
		case lists.List:
			for i := 0; i < tok.Len(); i++ {
				c, _ := tok.Get(i)
				queue.Push(c)
			}
		}
	}
}

func (s *PermuteOptionalsStrategy) Fuzz(r rand.Rand) chan struct{} {
	continueFuzzing := make(chan struct{})

	go func() {
		if tavor.DEBUG {
			fmt.Println("Start permute optionals routine")
		}

		optionals, resets := s.findOptionals(r, s.root, false)

		if len(optionals) != 0 {
			if !s.fuzz(r, continueFuzzing, optionals, resets) {
				return
			}
		}

		s.resetResetTokens()

		if tavor.DEBUG {
			fmt.Println("Done with fuzzing step")
		}

		for t := range resets {
			t.Reset()
		}

		// done with the last fuzzing step
		continueFuzzing <- struct{}{}

		if tavor.DEBUG {
			fmt.Println("Finished fuzzing. Wait till the outside is ready to close.")
		}

		if _, ok := <-continueFuzzing; ok {
			if tavor.DEBUG {
				fmt.Println("Close fuzzing channel")
			}

			close(continueFuzzing)
		}
	}()

	return continueFuzzing
}

func (s *PermuteOptionalsStrategy) fuzz(r rand.Rand, continueFuzzing chan struct{}, optionals []token.OptionalToken, resets map[token.ResetToken]struct{}) bool {
	if tavor.DEBUG {
		fmt.Printf("Fuzzing optionals %#v\n", optionals)
	}

	// TODO make this WAYYYYYYYYYYYYYYYYYYYYYYYYYYYYYY smarter
	// since we can only fuzz 64 optionals at max
	// https://en.wikipedia.org/wiki/Steinhaus%E2%80%93Johnson%E2%80%93Trotter_algorithm
	p := 0
	maxSteps := int(math.Pow(2, float64(len(optionals))))

	for {
		if tavor.DEBUG {
			fmt.Printf("Fuzzing step %b\n", p)
		}

		ith := 1

		for i := range optionals {
			if p&ith == 0 {
				optionals[i].Deactivate()
			} else {
				optionals[i].Activate()

				childs, rets := s.findOptionals(r, optionals[i], true)

				if len(rets) != 0 {
					for t := range rets {
						resets[t] = struct{}{}
					}
				}

				if len(childs) != 0 {
					if !s.fuzz(r, continueFuzzing, childs, resets) {
						return false
					}
				}
			}

			ith = ith << 1
		}

		p++

		if p == maxSteps {
			if tavor.DEBUG {
				fmt.Println("Done with fuzzing these optionals")
			}

			break
		}

		s.resetResetTokens()

		if tavor.DEBUG {
			fmt.Println("Done with fuzzing step")
		}

		for t := range resets {
			t.Reset()
		}

		// done with this fuzzing step
		continueFuzzing <- struct{}{}

		// wait until we are allowed to continue
		if _, ok := <-continueFuzzing; !ok {
			if tavor.DEBUG {
				fmt.Println("Fuzzing channel closed from outside")
			}

			return false
		}
	}

	return true
}
