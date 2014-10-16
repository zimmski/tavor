// Package tavor provides all general properties, constants and functions for the Tavor framework and tools
package tavor

import (
	"fmt"
	"io"
	"strings"

	"github.com/zimmski/container/list/linkedlist"

	"github.com/zimmski/tavor/log"
	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/primitives"
)

const (
	// Version of the framework and tools
	Version = "0.3"
)

// MaxRepeat determines the maximum copies in graph cycles.
var MaxRepeat = 2

// PrettyPrintTree prints the token represenation as a text tree
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
	case token.List:
		for i := 0; i < t.Len(); i++ {
			c, _ := t.Get(i)

			prettyPrintTreeRek(w, c, level+1)
		}
	}
}

// PrettyPrintInternalTree prints the internal token represenation as a text tree
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
	case token.List:
		for i := 0; i < t.InternalLen(); i++ {
			c, _ := t.InternalGet(i)

			prettyPrintInternalTreeRek(w, c, level+1)
		}
	}
}

// LoopExists determines if a cycle exists in the internal token graph
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
					log.Debugf("found a loop through (%p)%#v", t, t)

					return true
				}

				queue.Push(v)
			}
		case token.ForwardToken:
			if v := tok.InternalGet(); v != nil {
				queue.Push(v)
			}
		case token.List:
			for i := 0; i < tok.InternalLen(); i++ {
				c, _ := tok.InternalGet(i)

				queue.Push(c)
			}
		}
	}

	return false
}

// UnrollPointers unrolls pointer tokens by copying their referenced graphs.
// Pointers that lead to themselfs are unrolled at maximum MaxRepeat times.
func UnrollPointers(root token.Token) token.Token {
	type unrollToken struct {
		tok    token.Token
		parent *unrollToken
		counts map[token.Token]int
	}

	log.Debug("start unrolling pointers by cloning them")

	parents := make(map[token.Token]token.Token)
	changed := make(map[token.Token]struct{})

	originals := make(map[token.Token]token.Token)
	originalClones := make(map[token.Token]token.Token)

	queue := linkedlist.New()

	queue.Push(&unrollToken{
		tok:    root,
		parent: nil,
		counts: make(map[token.Token]int),
	})
	parents[root] = nil

	for !queue.Empty() {
		v, _ := queue.Shift()
		iTok, _ := v.(*unrollToken)

		switch t := iTok.tok.(type) {
		case *primitives.Pointer:
			child := t.InternalGet()

			if child == nil {
				log.Debugf("Child is nil")

				continue
			}

			replace := true

			if p, ok := child.(*primitives.Pointer); ok {
				checked := map[*primitives.Pointer]struct{}{
					p: struct{}{},
				}

				for {
					log.Debugf("Child (%p)%#v is a pointer lets go one further", p, p)

					child = p.InternalGet()

					p, ok = child.(*primitives.Pointer)
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

			var original token.Token
			counted := 0

			if replace {
				if o, found := originals[child]; found {
					log.Debugf("Found original (%p)%#v for child (%p)%#v", o, o, child, child)
					original = o
					counted = iTok.counts[original]

					if counted >= MaxRepeat {
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

				counts := make(map[token.Token]int)
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
					case token.ForwardToken:
						tt.InternalReplace(t, c)
					case token.List:
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
				log.Debugf("reached max repeat of %d for (%p)%#v with child (%p)%#v", MaxRepeat, t, t, child, child)

				_ = t.Set(nil)

				ta := iTok.tok
				tt := iTok.parent

				repl := func(parent token.Token, this token.Token, that token.Token) {
					log.Debugf("replace (%p)%#v by (%p)%#v", this, this, that, that)

					if parent != nil {
						changed[parent] = struct{}{}

						switch tt := parent.(type) {
						case token.ForwardToken:
							tt.InternalReplace(this, that)
						case token.List:
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
					case token.ForwardToken:
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
					case token.List:
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
		case token.ForwardToken:
			if v := t.InternalGet(); v != nil {
				queue.Push(&unrollToken{
					tok:    v,
					parent: iTok,
					counts: iTok.counts,
				})

				parents[v] = iTok.tok
			}
		case token.List:
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
		case token.ForwardToken:
			tt.InternalReplace(child, child)
		case token.List:
			tt.InternalReplace(child, child)
		}
	}

	log.Debug("finished unrolling")

	return root
}

// ResetResetTokens resets all tokens in the token graph that fullfill the ResetToken interface
func ResetResetTokens(root token.Token) {
	var queue = linkedlist.New()

	queue.Push(root)

	for !queue.Empty() {
		v, _ := queue.Shift()

		switch tok := v.(type) {
		case token.ResetToken:
			log.Debugf("reset %#v(%p)", tok, tok)

			tok.Reset()
		}

		switch tok := v.(type) {
		case token.ForwardToken:
			if v := tok.Get(); v != nil {
				queue.Push(v)
			}
		case token.List:
			for i := 0; i < tok.Len(); i++ {
				c, _ := tok.Get(i)
				queue.Push(c)
			}
		}
	}
}

// ResetScope resets all scopes of tokens in the token graph that fullfill the ScopeToken interface
func ResetScope(root token.Token) {
	SetScope(root, make(map[string]token.Token))
}

// SetScope sets all scopes of tokens in the token graph that fullfill the ScopeToken interface
func SetScope(root token.Token, scope map[string]token.Token) {
	queue := linkedlist.New()

	type set struct {
		token token.Token
		scope map[string]token.Token
	}

	queue.Push(set{
		token: root,
		scope: scope,
	})

	for !queue.Empty() {
		v, _ := queue.Shift()
		s := v.(set)

		if t, ok := s.token.(token.ScopeToken); ok {
			log.Debugf("setScope %#v(%p)", t, t)

			t.SetScope(s.scope)
		}

		nScope := make(map[string]token.Token, len(s.scope))
		for k, v := range s.scope {
			nScope[k] = v
		}

		switch t := s.token.(type) {
		case token.ForwardToken:
			if v := t.Get(); v != nil {
				queue.Push(set{
					token: v,
					scope: nScope,
				})
			}
		case token.List:
			for i := 0; i < t.Len(); i++ {
				c, _ := t.Get(i)

				queue.Push(set{
					token: c,
					scope: nScope,
				})
			}
		}
	}
}

// SetInternalScope sets all scopes of internal tokens in the token graph that fullfill the ScopeToken interface
func SetInternalScope(root token.Token, scope map[string]token.Token) {
	queue := linkedlist.New()

	queue.Push(root)

	for !queue.Empty() {
		tok, _ := queue.Shift()

		if t, ok := tok.(token.ScopeToken); ok {
			log.Debugf("setScope %#v(%p)", t, t)

			t.SetScope(scope)
		}

		switch t := tok.(type) {
		case token.ForwardToken:
			if v := t.InternalGet(); v != nil {
				queue.Push(v)
			}
		case token.List:
			for i := 0; i < t.InternalLen(); i++ {
				c, _ := t.InternalGet(i)

				queue.Push(c)
			}
		}
	}
}
