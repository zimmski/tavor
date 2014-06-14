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

type optionalLookup struct {
	token  token.OptionalToken
	childs []optionalLookup
}

type PermuteOptionalsStrategy struct {
	root token.Token

	continueFuzzing chan struct{}
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

func (s *PermuteOptionalsStrategy) findOptionals(root token.Token, fromChilds bool) []optionalLookup {
	var optionals []optionalLookup
	var queue = linkedlist.New()

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
		case token.OptionalToken:
			if !t.IsOptional() {
				opts := s.findOptionals(t, true)

				if len(opts) != 0 {
					optionals = append(optionals, opts...)
				}

				continue
			}

			if tavor.DEBUG {
				fmt.Printf("Found optional %#v\n", t)
			}

			t.Deactivate()

			optionals = append(optionals, optionalLookup{
				token:  t,
				childs: nil,
			})
		case token.ForwardToken:
			queue.Push(t.Get())
		case lists.List:
			l := t.Len()

			for i := 0; i < l; i++ {
				c, _ := t.Get(i)
				queue.Push(c)
			}
		}
	}

	return optionals
}

func (s *PermuteOptionalsStrategy) Fuzz(r rand.Rand) chan struct{} {
	continueFuzzing := make(chan struct{})

	go func() {
		if tavor.DEBUG {
			fmt.Println("Start permute optionals routine")
		}

		optionals := s.findOptionals(s.root, false)

		if len(optionals) != 0 {
			if !s.fuzz(continueFuzzing, optionals) {
				return
			}
		}

		if tavor.DEBUG {
			fmt.Println("Done with fuzzing step")
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

func (s *PermuteOptionalsStrategy) fuzz(continueFuzzing chan struct{}, optionals []optionalLookup) bool {
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
				optionals[i].token.Deactivate()
				optionals[i].childs = nil
			} else {
				optionals[i].token.Activate()
				optionals[i].childs = s.findOptionals(optionals[i].token, true)

				if len(optionals[i].childs) != 0 {
					if !s.fuzz(continueFuzzing, optionals[i].childs) {
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

		if tavor.DEBUG {
			fmt.Println("Done with fuzzing step")
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
