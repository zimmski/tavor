package token

import (
	"github.com/zimmski/container/list/linkedlist"

	"github.com/zimmski/tavor"
	"github.com/zimmski/tavor/log"
)

// UnrollPointers unrolls pointer tokens by copying their referenced graphs.
// Pointers that lead to themselfs are unrolled at maximum tavor.MaxRepeat times.
func UnrollPointers(root Token) Token {
	type unrollToken struct {
		tok    Token
		parent *unrollToken
		counts map[Token]int
	}

	log.Debug("start unrolling pointers by cloning them")

	parents := make(map[Token]Token)
	changed := make(map[Token]struct{})

	originals := make(map[Token]Token)
	originalClones := make(map[Token]Token)

	queue := linkedlist.New()

	queue.Push(&unrollToken{
		tok:    root,
		parent: nil,
		counts: make(map[Token]int),
	})
	parents[root] = nil

	for !queue.Empty() {
		v, _ := queue.Shift()
		iTok, _ := v.(*unrollToken)

		switch t := iTok.tok.(type) {
		case PointerToken:
			child := t.InternalGet()

			if child == nil {
				log.Debugf("Child is nil")

				continue
			}

			replace := true

			if p, ok := child.(PointerToken); ok {
				checked := map[PointerToken]struct{}{
					p: struct{}{},
				}

				for {
					log.Debugf("Child (%p)%#v is a pointer lets go one further", p, p)

					child = p.InternalGet()

					p, ok = child.(PointerToken)
					if !ok {
						break
					}

					if _, found := checked[p]; found {
						log.Debugf("Endless pointer loop with (%p)%#v", p, p)

						replace = false

						break
					}

					checked[p] = struct{}{}
				}
			}

			var original Token
			counted := 0

			if replace {
				if o, found := originals[child]; found {
					log.Debugf("Found original (%p)%#v for child (%p)%#v", o, o, child, child)
					original = o
					counted = iTok.counts[original]

					if counted >= tavor.MaxRepeat {
						replace = false
					}
				} else {
					log.Debugf("Found no original for child (%p)%#v, must be new!", child, child)
					originals[child] = child
					original = child

					// we want to clone only original structures so we always clone the clone since the original could have been changed in the meantime
					originalClones[child] = child.Clone()
				}
			}

			if replace {
				log.Debugf("clone (%p)%#v with child (%p)%#v", t, t, child, child)

				c := originalClones[original].Clone()

				counts := make(map[Token]int)
				for k, v := range iTok.counts {
					counts[k] = v
				}

				counts[original] = counted + 1
				originals[c] = original

				log.Debugf("clone is (%p)%#v", c, c)

				if err := t.Set(c); err != nil {
					panic(err)
				}

				if iTok.parent != nil {
					log.Debugf("replace in (%p)%#v", iTok.parent.tok, iTok.parent.tok)

					changed[iTok.parent.tok] = struct{}{}

					switch tt := iTok.parent.tok.(type) {
					case ForwardToken:
						tt.InternalReplace(t, c)
					case List:
						tt.InternalReplace(t, c)
					}
				} else {
					log.Debugf("replace as root")

					root = c
				}

				queue.Unshift(&unrollToken{
					tok:    c,
					parent: iTok.parent,
					counts: counts,
				})
			} else {
				// we reached a maximum of repetition, we cut and remove dangling tokens
				log.Debugf("reached max repeat of %d for (%p)%#v with child (%p)%#v", tavor.MaxRepeat, t, t, child, child)

				_ = t.Set(nil)

				ta := iTok.tok
				tt := iTok.parent

				repl := func(parent Token, this Token, that Token) {
					log.Debugf("replace (%p)%#v by (%p)%#v", this, this, that, that)

					if parent != nil {
						changed[parent] = struct{}{}

						switch tt := parent.(type) {
						case ForwardToken:
							tt.InternalReplace(this, that)
						case List:
							tt.InternalReplace(this, that)
						}
					} else {
						log.Debugf("replace as root")

						root = that
					}
				}

			REMOVE:
				for tt != nil {
					delete(parents, tt.tok)
					delete(changed, tt.tok)

					switch l := tt.tok.(type) {
					case ForwardToken:
						log.Debugf("remove (%p)%#v from (%p)%#v", ta, ta, l, l)

						c := l.InternalLogicalRemove(ta)

						if c != nil {
							if c != l {
								repl(tt.parent.tok, l, c)
							}

							break REMOVE
						}

						ta = l
						tt = tt.parent
					case List:
						log.Debugf("remove (%p)%#v from (%p)%#v", ta, ta, l, l)

						c := l.InternalLogicalRemove(ta)

						if c != nil {
							if c != l {
								repl(tt.parent.tok, l, c)
							}

							break REMOVE
						}

						ta = l
						tt = tt.parent
					}
				}
			}
		case ForwardToken:
			if v := t.InternalGet(); v != nil {
				queue.Push(&unrollToken{
					tok:    v,
					parent: iTok,
					counts: iTok.counts,
				})

				parents[v] = iTok.tok
			}
		case List:
			for i := 0; i < t.InternalLen(); i++ {
				c, _ := t.InternalGet(i)

				queue.Push(&unrollToken{
					tok:    c,
					parent: iTok,
					counts: iTok.counts,
				})

				parents[c] = iTok.tok
			}
		}
	}

	// we need to update some tokens with the same child to regenerate clones
	for child := range changed {
		parent := parents[child]

		if parent == nil {
			continue
		}

		log.Debugf("update (%p)%#v with child (%p)%#v", parent, parent, child, child)

		switch tt := parent.(type) {
		case ForwardToken:
			tt.InternalReplace(child, child)
		case List:
			tt.InternalReplace(child, child)
		}
	}

	log.Debug("finished unrolling")

	return root
}
