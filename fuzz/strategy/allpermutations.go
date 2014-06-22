package strategy

import (
	"fmt"

	"github.com/zimmski/container/list/linkedlist"

	"github.com/zimmski/tavor"
	"github.com/zimmski/tavor/rand"
	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/lists"
)

type allPermutationsLevel struct {
	token           token.Token
	permutation     int
	maxPermutations int

	childs []allPermutationsLevel
}

type AllPermutationsStrategy struct {
	root token.Token
}

func NewAllPermutationsStrategy(tok token.Token) *AllPermutationsStrategy {
	s := &AllPermutationsStrategy{
		root: tok,
	}

	return s
}

func init() {
	Register("AllPermutations", func(tok token.Token) Strategy {
		return NewAllPermutationsStrategy(tok)
	})
}

func (s *AllPermutationsStrategy) getTree(root token.Token, fromChilds bool) []allPermutationsLevel {
	var tree []allPermutationsLevel

	add := func(tok token.Token) {
		s.setPermutation(tok, 1)

		tree = append(tree, allPermutationsLevel{
			token:           tok,
			permutation:     1,
			maxPermutations: tok.Permutations(),

			childs: s.getTree(tok, true),
		})
	}

	if fromChilds {
		switch t := root.(type) {
		case token.ForwardToken:
			if v := t.Get(); v != nil {
				add(v)
			}
		case lists.List:
			for i := 0; i < t.Len(); i++ {
				c, _ := t.Get(i)

				add(c)
			}
		}
	} else {
		add(root)
	}

	return tree
}

func (s *AllPermutationsStrategy) setPermutation(tok token.Token, permutation int) {
	if tavor.DEBUG {
		fmt.Printf("Set %#v(%p) to permutation %d\n", tok, tok, permutation)
	}

	if err := tok.Permutation(permutation); err != nil {
		panic(err)
	}
}

func (s *AllPermutationsStrategy) Fuzz(r rand.Rand) chan struct{} {
	continueFuzzing := make(chan struct{})

	go func() {
		if tavor.DEBUG {
			fmt.Println("Start all permutations routine")
		}

		tree := s.getTree(s.root, false)

		if len(tree) != 0 {
			if tavor.DEBUG {
				fmt.Println("Start fuzzing step")
			}

			if contin, _ := s.fuzz(continueFuzzing, tree, false); !contin {
				return
			}

			if tavor.DEBUG {
				fmt.Println("Finished fuzzing.")
			}

			close(continueFuzzing)
		} else {
			s.resetResetTokens()

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
		}
	}()

	return continueFuzzing
}

func (s *AllPermutationsStrategy) fuzz(continueFuzzing chan struct{}, tree []allPermutationsLevel, justastep bool) (bool, bool) {
	if tavor.DEBUG {
		fmt.Printf("Fuzzing level %d->%#v\n", len(tree), tree)
	}

STEP:
	for {
		if justastep && len(tree[0].childs) != 0 {
			if tavor.DEBUG {
				fmt.Printf("STEP FURTHER INTO %#v\n", tree[0])
			}

			if contin, step := s.fuzz(continueFuzzing, tree[0].childs, justastep); !contin {
				return false, false
			} else if step {
				if tavor.DEBUG {
					fmt.Printf("CONTINUE after child step\n")
				}

				return true, true
			}
		} else {
			if tavor.DEBUG {
				fmt.Printf("Permute %d->%#v\n", 0, tree[0])
			}

			if tree[0].permutation != 1 {
				s.setPermutation(tree[0].token, tree[0].permutation)
				tree[0].childs = s.getTree(tree[0].token, true)

				if justastep {
					if tavor.DEBUG {
						fmt.Printf("CONTINUE after permutate\n")
					}

					return true, true
				}
			}

			if len(tree[0].childs) != 0 {
				if contin, step := s.fuzz(continueFuzzing, tree[0].childs, justastep); !contin {
					return false, false
				} else if step {
					if tavor.DEBUG {
						fmt.Printf("CONTINUE after child step\n")
					}

					return true, true
				}
			} else {
				if !justastep && (tree[0].token != s.root || tree[0].permutation <= tree[0].maxPermutations) && !s.nextStep(continueFuzzing) {
					return false, false
				}
			}
		}

		tree[0].permutation++

		if tree[0].permutation > tree[0].maxPermutations {
			/*if tavor.DEBUG {
				for i := 1; i < len(tree); i++ {
					fmt.Printf("Check %d vs %d for %#v\n", tree[i].permutation, tree[i].maxPermutations, tree[i])
				}
			}
			found := false
			for i := 1; i < len(tree); i++ {
				if tavor.DEBUG {
					fmt.Printf("%d vs %d for %#v\n", tree[i].permutation, tree[i].maxPermutations, tree[i])
				}
				if tree[i].permutation < tree[i].maxPermutations {
					found = true

					break
				}
			}
			if !found {
				if tavor.DEBUG {
					fmt.Printf("Done with fuzzing this level because %#v\n", tree)
				}

				break STEP
			}*/

			if tavor.DEBUG {
				for i := 0; i < len(tree); i++ {
					fmt.Printf("Check %d vs %d for %#v\n", tree[i].permutation, tree[i].maxPermutations, tree[i])
				}
			}

			i := 0

			for {
				if i == len(tree)-1 {
					if tavor.DEBUG {
						fmt.Printf("Done with fuzzing this level because %#v\n", tree)
					}

					break STEP
				}

				i++

				if len(tree[i].childs) != 0 {
					if tavor.DEBUG {
						fmt.Printf("CHECK children %#v\n", tree[i])
					}

					if contin, step := s.fuzz(continueFuzzing, tree[i].childs, true); !contin {
						return false, false
					} else if step {
						for j := 0; j < i; j++ {
							tree[j].permutation = 1
							s.setPermutation(tree[j].token, tree[j].permutation)
							tree[j].childs = s.getTree(tree[j].token, true)
						}

						if justastep {
							return true, true
						}

						if tavor.DEBUG {
							fmt.Printf("STEP continue\n")
						}
						continue STEP
					}
					if tavor.DEBUG {
						fmt.Printf("PERMUTATE continue\n")
					}
				}

				tree[i].permutation++

				if tree[i].permutation <= tree[i].maxPermutations {
					for j := 0; j < i; j++ {
						tree[j].permutation = 1
						s.setPermutation(tree[j].token, tree[j].permutation)
						tree[j].childs = s.getTree(tree[j].token, true)
					}

					if tavor.DEBUG {
						fmt.Printf("Permute %d->%#v\n", i, tree[i])
					}

					s.setPermutation(tree[i].token, tree[i].permutation)
					tree[i].childs = s.getTree(tree[i].token, true)

					if justastep {
						return true, true
					}

					continue STEP
				}
			}
		}
	}

	return true, false
}

func (s *AllPermutationsStrategy) nextStep(continueFuzzing chan struct{}) bool {
	s.resetResetTokens()

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

	if tavor.DEBUG {
		fmt.Println("Start fuzzing step")
	}

	return true
}

/*func (s *AllPermutationsStrategy) fuzz(continueFuzzing chan struct{}, tree []allPermutationsLevel) bool {
	if tavor.DEBUG {
		fmt.Printf("Fuzzing level %d->%#v\n", len(tree), tree)
	}

STEP:
	for {
		if tavor.DEBUG {
			fmt.Printf("Permute %d->%#v\n", 0, tree[0])
		}

		if tree[0].permutation != 1 {
			s.setPermutation(tree[0].token, tree[0].permutation)
			tree[0].childs = s.getTree(tree[0].token, true)
		}

		if len(tree[0].childs) != 0 {
			if !s.fuzz(continueFuzzing, tree[0].childs) {
				return false
			}
		}

		tree[0].permutation++

		if tree[0].permutation > tree[0].maxPermutations {
			found := false
			for i := 1; i < len(tree); i++ {
				if tavor.DEBUG {
					fmt.Printf("%d vs %d for %#v\n", tree[i].permutation, tree[i].maxPermutations, tree[i])
				}
				if tree[i].permutation < tree[i].maxPermutations {
					found = true

					break
				}
			}
			if !found {
				if tavor.DEBUG {
					fmt.Printf("Done with fuzzing this level because %#v\n", tree)
				}

				break STEP
			}

			i := 0

			if i < len(tree)-1 {
				s.resetResetTokens()

				if tavor.DEBUG {
					fmt.Println("Done with fuzzing step INNER")
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

				if tavor.DEBUG {
					fmt.Println("Start fuzzing step")
				}
			}

			for {
				tree[i].permutation = 1
				s.setPermutation(tree[i].token, tree[i].permutation)
				tree[i].childs = s.getTree(tree[i].token, true)

				i++

				tree[i].permutation++

				if tree[i].permutation <= tree[i].maxPermutations {
					if tavor.DEBUG {
						fmt.Printf("Permute %d->%#v\n", i, tree[i])
					}

					s.setPermutation(tree[i].token, tree[i].permutation)
					tree[i].childs = s.getTree(tree[i].token, true)

					if len(tree[i].childs) != 0 {
						if !s.fuzz(continueFuzzing, tree[i].childs) {
							return false
						}
					}

					continue STEP
				}
			}
		}

		s.resetResetTokens()

		if tavor.DEBUG {
			fmt.Println("Done with fuzzing step OUTER")
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

		if tavor.DEBUG {
			fmt.Println("Start fuzzing step")
		}
	}

	return true
}*/

func (s *AllPermutationsStrategy) resetResetTokens() {
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
