package tavor

import (
	"fmt"
	"io"
	"strings"

	"github.com/zimmski/container/list/linkedlist"

	"github.com/zimmski/tavor/log"
	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/lists"
	"github.com/zimmski/tavor/token/primitives"
)

const (
	Version = "0.1"
)

const (
	MaxRepeat = 2
)

func PrettyPrintTree(w io.Writer, root token.Token) {
	prettyPrintTreeRek(w, root, 0)
}

func prettyPrintTreeRek(w io.Writer, tok token.Token, level int) {
	fmt.Fprintf(w, "%s(%p)%#v\n", strings.Repeat("\t", level), tok, tok)

	switch t := tok.(type) {
	case token.ForwardToken:
		if v := t.Get(); v != nil {
			prettyPrintTreeRek(w, v, level+1)
		}
	case lists.List:
		for i := 0; i < t.Len(); i++ {
			c, _ := t.Get(i)

			prettyPrintTreeRek(w, c, level+1)
		}
	}
}

func PrettyPrintInternalTree(w io.Writer, root token.Token) {
	prettyPrintInternalTreeRek(w, root, 0)
}

func prettyPrintInternalTreeRek(w io.Writer, tok token.Token, level int) {
	fmt.Fprintf(w, "%s(%p)%#v\n", strings.Repeat("\t", level), tok, tok)

	switch t := tok.(type) {
	case token.ForwardToken:
		if v := t.InternalGet(); v != nil {
			prettyPrintInternalTreeRek(w, v, level+1)
		}
	case lists.List:
		for i := 0; i < t.InternalLen(); i++ {
			c, _ := t.InternalGet(i)

			prettyPrintInternalTreeRek(w, c, level+1)
		}
	}
}

func LoopExists(root token.Token) bool {
	lookup := make(map[token.Token]struct{})
	queue := linkedlist.New()

	queue.Push(root)

	for !queue.Empty() {
		v, _ := queue.Shift()
		t, _ := v.(token.Token)

		lookup[t] = struct{}{}

		switch tok := t.(type) {
		case *primitives.Pointer:
			if v := tok.InternalGet(); v != nil {
				if _, ok := lookup[v]; ok {
					log.Debugf("Found a loop through (%p)%+v", t)

					return true
				}

				queue.Push(v)
			}
		case token.ForwardToken:
			if v := tok.InternalGet(); v != nil {
				queue.Push(v)
			}
		case lists.List:
			for i := 0; i < tok.InternalLen(); i++ {
				c, _ := tok.InternalGet(i)

				queue.Push(c)
			}
		}
	}

	return false
}

func UnrollPointers(root token.Token) token.Token {
	type unrollToken struct {
		tok    token.Token
		parent *unrollToken
	}

	log.Debug("Unroll pointers by cloning them")

	checked := make(map[token.Token]token.Token)
	counters := make(map[token.Token]int)

	queue := linkedlist.New()

	queue.Push(&unrollToken{
		tok:    root,
		parent: nil,
	})

	for !queue.Empty() {
		v, _ := queue.Shift()
		iTok, _ := v.(*unrollToken)

		switch t := iTok.tok.(type) {
		case *primitives.Pointer:
			o := t.InternalGet()

			parent, ok := checked[o]
			times := 0

			if ok {
				times = counters[parent]
			} else {
				parent = o.Clone()
				checked[o] = parent
			}

			if times != MaxRepeat {
				log.Debugf("Clone (%p)%#v with parent (%p)%#v", t, t, parent, parent)

				c := parent.Clone()

				t.Set(c)

				counters[parent] = times + 1
				checked[c] = parent

				if iTok.parent != nil {
					switch tt := iTok.parent.tok.(type) {
					case token.ForwardToken:
						tt.InternalReplace(t, c)
					case lists.List:
						tt.InternalReplace(t, c)
					}
				} else {
					root = c
				}

				queue.Unshift(&unrollToken{
					tok:    c,
					parent: iTok.parent,
				})
			} else {
				log.Debugf("Reached max repeat of %d for (%p)%#v with parent (%p)%#v", MaxRepeat, t, t, parent, parent)

				t.Set(nil)

				ta := iTok.tok
				tt := iTok.parent

			REMOVE:
				for tt != nil {
					switch l := tt.tok.(type) {
					case token.ForwardToken:
						log.Debugf("Remove (%p)%#v from (%p)%#v", ta, ta, l, l)

						c := l.InternalLogicalRemove(ta)

						if c != nil {
							break REMOVE
						}

						ta = l
						tt = tt.parent
					case lists.List:
						log.Debugf("Remove (%p)%#v from (%p)%#v", ta, ta, l, l)

						c := l.InternalLogicalRemove(ta)

						if c != nil {
							break REMOVE
						}

						ta = l
						tt = tt.parent
					}
				}
			}
		case token.ForwardToken:
			if v := t.InternalGet(); v != nil {
				queue.Push(&unrollToken{
					tok:    v,
					parent: iTok,
				})
			}
		case lists.List:
			for i := 0; i < t.InternalLen(); i++ {
				c, _ := t.InternalGet(i)

				queue.Push(&unrollToken{
					tok:    c,
					parent: iTok,
				})
			}
		}
	}

	log.Debug("Done unrolling")

	return root
}
